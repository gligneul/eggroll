// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/eggeth"
	"github.com/gligneul/eggroll/eggtypes"
	"github.com/gligneul/eggroll/internal/reader"
	"github.com/gligneul/eggroll/internal/sunodo"
)

// Configuration for the client struct.
type ClientConfig struct {
	Codecs           []Codec
	DAppAddress      common.Address
	GraphqlEndpoint  string
	InspectEndpoint  string
	ProviderEndpoint string
}

// The client interacts with the DApp contract off-chain.
type Client struct {
	ClientConfig
	codecManager *codecManager
	reader       *reader.GraphQLReader
	inspect      *reader.InspectClient
	eth          *eggeth.ETHClient
}

// Create a new client with the given config.
func NewClient(config ClientConfig) (*Client, error) {
	ethClient, err := eggeth.NewETHClient(config.ProviderEndpoint, config.DAppAddress)
	if err != nil {
		return nil, err
	}
	client := &Client{
		ClientConfig: config,
		codecManager: newCodecManager(config.Codecs),
		reader:       reader.NewGraphQLReader(config.GraphqlEndpoint),
		inspect:      reader.NewInspectClient(config.InspectEndpoint),
		eth:          ethClient,
	}
	return client, nil
}

// Create a new client for local development.
// Connects to the Rollups Node and the Ethereum Node setup by sunodo.
// Return a signer that uses the Foundry's test mnemonic to send transactions.
func NewDevClient(ctx context.Context, codecs []Codec) (*Client, eggeth.Signer, error) {
	dappAddress, err := sunodo.GetDAppAddress()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get DApp address: %v", err)
	}
	config := ClientConfig{
		Codecs:           codecs,
		DAppAddress:      dappAddress,
		GraphqlEndpoint:  "http://localhost:8080/graphql",
		InspectEndpoint:  "http://localhost:8080/inspect",
		ProviderEndpoint: "ws://localhost:8545",
	}
	client, err := NewClient(config)
	if err != nil {
		return nil, nil, err
	}
	chainId, err := client.eth.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain id: %v", err)
	}
	signer, err := eggeth.NewMnemonicSigner(
		"test test test test test test test test test test test junk", 0, chainId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create signer: %v", err)
	}
	return client, signer, nil
}

//
// Send functions
//

// Send the input to the DApp contract.
// If the input has type []byte send it as raw bytes; otherwise, use codecs to encode it.
// This function waits until the transaction is added to a block and return the input index.
func (c *Client) SendInput(ctx context.Context, signer eggeth.Signer, input any) (int, error) {
	inputBytes, err := c.encodeInput(input)
	if err != nil {
		return 0, err
	}
	return c.eth.SendInput(ctx, signer, inputBytes)
}

// Send the DApp address to the DApp contract with the DAppAddressRelay contract.
// This function waits until the transaction is added to a block and return the input index.
func (c *Client) SendDAppAddress(ctx context.Context, signer eggeth.Signer) (int, error) {
	return c.eth.SendDAppAddress(ctx, signer)
}

// Send Ether to the Ether portal. This function also receives an optional input.
// If the input has type []byte send it as raw bytes; otherwise, use codecs to encode it.
// This function waits until the transaction is added to a block and return the input index.
func (c *Client) SendEther(ctx context.Context, signer eggeth.Signer, txValue *big.Int, input any) (
	int, error) {

	inputBytes, err := c.encodeInput(input)
	if err != nil {
		return 0, err
	}
	return c.eth.SendEther(ctx, signer, txValue, inputBytes)
}

//
// Reader functions
//

// Wait until the DApp contract processes a given input.
// Returns the advance result of that input.
func (c *Client) WaitFor(ctx context.Context, inputIndex int) (*eggtypes.AdvanceResult, error) {
	for {
		result, err := c.reader.AdvanceResult(ctx, inputIndex)
		if err != nil {
			if _, ok := err.(reader.NotFound); ok {
				goto wait
			}
			return nil, fmt.Errorf("faild to read result: %v", err)
		}
		if result.Status != eggtypes.CompletionStatusUnprocessed {
			return result, nil
		}
	wait:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			continue
		}
	}
}

// Get the results starting from the given input index.
func (c *Client) GetResults(ctx context.Context, inputIndex int) (
	[]*eggtypes.AdvanceResult, error) {

	var results []*eggtypes.AdvanceResult
	for {
		result, err := c.reader.AdvanceResult(ctx, inputIndex)
		if err != nil {
			if _, ok := err.(reader.NotFound); ok {
				break
			}
			return nil, fmt.Errorf("failed to get result: %v", err)
		}
		if result.Status == eggtypes.CompletionStatusUnprocessed {
			break
		}
		results = append(results, result)
		inputIndex++
	}
	return results, nil
}

// Decode the return from a result.
// The result can be either an advance result or an inspect result.
func (c *Client) DecodeReturn(result interface{ RawReturn() []byte }) any {
	return_ := result.RawReturn()
	if return_ == nil {
		return nil
	}
	return c.codecManager.decode(return_)
}

// Send an inspect request.
func (c *Client) Inspect(ctx context.Context, input any) (*eggtypes.InspectResult, error) {
	inputBytes, err := c.encodeInput(input)
	if err != nil {
		return nil, err
	}
	return c.inspect.Inspect(ctx, inputBytes)
}

//
// Private functions
//

func (c *Client) encodeInput(input any) ([]byte, error) {
	if input == nil {
		return nil, nil
	}
	inputBytes, ok := input.([]byte)
	if ok {
		return inputBytes, nil
	}
	return c.codecManager.encode(input)
}
