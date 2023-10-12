// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/blockchain"
	"github.com/gligneul/eggroll/internal/sunodo"
	"github.com/gligneul/eggroll/reader"
)

// Result of an advance input.
type AdvanceResult struct {
	*reader.Input

	// Value returned by the contract advance method.
	RawReturn []byte

	// Logs generated during the advance method.
	Logs []string

	codecManager *codecManager
}

func newAdvanceResult(input *reader.Input, codecManager *codecManager) *AdvanceResult {
	result := &AdvanceResult{
		Input:        input,
		codecManager: codecManager,
	}
	for _, report := range input.Reports {
		tag, payload, err := decodeReport(report.Payload)
		if err != nil {
			// TODO how do we report this?
			continue
		}
		switch tag {
		case reportTagReturn:
			result.RawReturn = payload
		case reportTagLog:
			result.Logs = append(result.Logs, string(payload))
		}
	}
	return result
}

// Decode the return value using the codecs.
// If fails, return the error in the place of the value.
func (r *AdvanceResult) DecodeReturn() any {
	return r.codecManager.decode(r.RawReturn)
}

// Configuration for the DevClient.
type DevClientConfig struct {
	Codecs           []Codec
	DAppAddress      common.Address
	GraphqlEndpoint  string
	ProviderEndpoint string
	Mnemonic         string
	AccountIndex     uint32
}

// The DevClient interacts with the DApp contract off chain.
type DevClient struct {
	DevClientConfig
	codecManager *codecManager
	reader       *reader.GraphQLReader
	blockchain   *blockchain.ETHClient
}

// Create the DevClient with a custom config.
func NewDevClientWithConfig(config DevClientConfig) (*DevClient, error) {
	blockchainAPI, err := blockchain.NewETHClient(config.ProviderEndpoint)
	if err != nil {
		return nil, err
	}
	client := &DevClient{
		DevClientConfig: config,
		codecManager:    newCodecManager(config.Codecs),
		reader:          reader.NewGraphQLReader(config.GraphqlEndpoint),
		blockchain:      blockchainAPI,
	}
	return client, nil
}

// Create the DevClient for local development.
// Connects to the Rollups Node and the Ethereum Node setup by sunodo.
// This client will use the Foundry's test mnemonic to send transactions.
func NewDevClient(codecs []Codec) (*DevClient, error) {
	dappAddress, err := sunodo.GetDAppAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get DApp address: %v", err)
	}

	config := DevClientConfig{
		Codecs:           codecs,
		DAppAddress:      dappAddress,
		GraphqlEndpoint:  "http://localhost:8080/graphql",
		ProviderEndpoint: "ws://localhost:8545",
		Mnemonic:         "test test test test test test test test test test test junk",
		AccountIndex:     0,
	}
	return NewDevClientWithConfig(config)
}

//
// Send functions
//

// Send the input to the DApp contract.
// If the input has type []byte send it as raw bytes; otherwise, use codecs to encode it.
// This function waits until the transaction is added to a block and return the input index.
func (c *DevClient) SendInput(ctx context.Context, input any) (int, error) {
	inputBytes, ok := input.([]byte)
	if !ok {
		var err error
		inputBytes, err = c.codecManager.encode(input)
		if err != nil {
			return 0, err
		}
	}
	privateKey, err := blockchain.MnemonicToPrivateKey(c.Mnemonic, c.AccountIndex)
	if err != nil {
		return 0, err
	}
	signer, err := c.blockchain.CreateSigner(ctx, privateKey)
	if err != nil {
		return 0, err
	}
	tx, err := c.blockchain.SendInput(ctx, signer, c.DAppAddress, inputBytes)
	if err != nil {
		return 0, err
	}
	err = c.blockchain.WaitForTransaction(ctx, tx)
	if err != nil {
		return 0, err
	}
	inputIndex, err := c.blockchain.GetInputIndex(ctx, tx)
	if err != nil {
		return 0, err
	}
	return inputIndex, nil
}

//
// Reader functions
//

// Wait until the DApp contract processes a given input.
// Returns the advance result of that input.
func (c *DevClient) WaitFor(ctx context.Context, inputIndex int) (*AdvanceResult, error) {
	for {
		input, err := c.reader.Input(ctx, inputIndex)
		if err != nil {
			if _, ok := err.(reader.NotFound); ok {
				goto wait
			}
			return nil, fmt.Errorf("faild to read input: %v", err)
		}
		if input.Status != reader.CompletionStatusUnprocessed {
			return newAdvanceResult(input, c.codecManager), nil
		}
	wait:
		time.Sleep(time.Second)
	}
}

// Sync to the latest Dapp state.
// Return the updated slice of Advance results.
func (c *DevClient) Sync(ctx context.Context, results []*AdvanceResult) ([]*AdvanceResult, error) {
	inputIndex := 0
	if len(results) != 0 {
		inputIndex = results[len(results)-1].Index
	}
	for {
		input, err := c.reader.Input(ctx, inputIndex)
		if err != nil {
			if _, ok := err.(reader.NotFound); ok {
				break
			}
			return nil, fmt.Errorf("failed to read input: %v", err)
		}
		if input.Status == reader.CompletionStatusUnprocessed {
			break
		}
		results = append(results, newAdvanceResult(input, c.codecManager))
		inputIndex++
	}
	return results, nil
}
