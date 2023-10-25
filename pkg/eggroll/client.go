// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/gligneul/eggroll/pkg/eggeth"
	"github.com/gligneul/eggroll/pkg/eggtypes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/internal/reader"
	"github.com/gligneul/eggroll/internal/sunodo"
)

// Configuration for the client struct.
type ClientConfig struct {
	DAppAddress      common.Address
	GraphqlEndpoint  string
	InspectEndpoint  string
	ProviderEndpoint string
}

// The client interacts with the DApp contract off-chain.
type Client struct {
	ClientConfig
	reader  *reader.GraphQLReader
	inspect *reader.InspectClient
	eth     *eggeth.ETHClient
}

// Create a new client with the given config.
func NewClient(config ClientConfig) (*Client, error) {
	ethClient, err := eggeth.NewETHClient(config.ProviderEndpoint, config.DAppAddress)
	if err != nil {
		return nil, err
	}
	client := &Client{
		ClientConfig: config,
		reader:       reader.NewGraphQLReader(config.GraphqlEndpoint),
		inspect:      reader.NewInspectClient(config.InspectEndpoint),
		eth:          ethClient,
	}
	return client, nil
}

// Create a new client for local development.
// Connects to the Rollups Node and the Ethereum Node setup by sunodo.
// Return a signer that uses the Foundry's test mnemonic to send transactions.
func NewDevClient(ctx context.Context) (*Client, eggeth.Signer, error) {
	dappAddress, err := sunodo.GetDAppAddress()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get DApp address: %v", err)
	}
	config := ClientConfig{
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
// This function waits until the transaction is added to a block and return the input index.
func (c *Client) SendInput(ctx context.Context, signer eggeth.Signer, payload []byte) (int, error) {
	return c.eth.SendInput(ctx, signer, payload)
}

// Send the DApp address to the DApp contract with the DAppAddressRelay contract.
// This function waits until the transaction is added to a block and return the input index.
func (c *Client) SendDAppAddress(ctx context.Context, signer eggeth.Signer) (int, error) {
	return c.eth.SendDAppAddress(ctx, signer)
}

// Send Ether to the Ether portal. This function also receives an optional input.
// If the input has type []byte send it as raw bytes; otherwise, use codecs to encode it.
// This function waits until the transaction is added to a block and return the input index.
func (c *Client) SendEther(
	ctx context.Context, signer eggeth.Signer, txValue *big.Int, payload []byte) (int, error) {

	return c.eth.SendEther(ctx, signer, txValue, payload)
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
// func (c *Client) GetResults(ctx context.Context, inputIndex int) (
// 	[]*eggtypes.AdvanceResult, error) {
//
// 	var results []*eggtypes.AdvanceResult
// 	for {
// 		result, err := c.reader.AdvanceResult(ctx, inputIndex)
// 		if err != nil {
// 			if _, ok := err.(reader.NotFound); ok {
// 				break
// 			}
// 			return nil, fmt.Errorf("failed to get result: %v", err)
// 		}
// 		if result.Status == eggtypes.CompletionStatusUnprocessed {
// 			break
// 		}
// 		results = append(results, result)
// 		inputIndex++
// 	}
// 	return results, nil
// }

// Send an inspect request.
func (c *Client) Inspect(ctx context.Context, payload []byte) (*eggtypes.InspectResult, error) {
	return c.inspect.Inspect(ctx, payload)
}
