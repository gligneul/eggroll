// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gligneul/eggroll/blockchain"
	"github.com/gligneul/eggroll/reader"
)

var logger *log.Logger

func init() {
	flags := log.LstdFlags | log.Lmsgprefix
	logger = log.New(os.Stdout, "eggroll: ", flags)
}

// Result of an advance input.
type AdvanceResult struct {
	*reader.Input

	// Result returned by the contract advance method.
	Result []byte

	// Logs generated during the advance method.
	Logs []string
}

func newAdvanceResult(input *reader.Input) *AdvanceResult {
	var result AdvanceResult
	result.Input = input
	for _, report := range input.Reports {
		tag, payload, err := decodeReport(report.Payload)
		if err != nil {
			logger.Printf("failed to decode report: %v", err)
			continue
		}
		switch tag {
		case reportTagResult:
			result.Result = payload
		case reportTagLog:
			result.Logs = append(result.Logs, string(payload))
		}
	}
	return &result
}

// Read the rollups state off chain.
// For more details, see the eggroll/reader package.
type readerAPI interface {
	Input(ctx context.Context, index int) (*reader.Input, error)
}

// Communicate with the blockchain.
type blockchainAPI interface {
	SendInput(ctx context.Context, input []byte) error
}

// The Client interacts with the DApp contract off chain.
type Client struct {
	reader     readerAPI
	blockchain blockchainAPI
	state      []byte
}

// Configuration for the Client.
type ClientConfig struct {
	GraphqlEndpoint  string
	ProviderEndpoint string
}

// Create the Client with a custom config.
func NewClient(config ClientConfig) *Client {
	return &Client{
		reader:     reader.NewGraphQLReader(config.GraphqlEndpoint),
		blockchain: blockchain.NewETHClient(config.ProviderEndpoint),
	}
}

// Create the Client loading the config from environment variables.
func NewLocalClient() *Client {
	config := ClientConfig{
		GraphqlEndpoint:  "http://localhost:8080/graphql",
		ProviderEndpoint: "http://localhost:8545",
	}
	return NewClient(config)
}

//
// Send functions
//

// Send the input as bytes to the DApp contract.
func (c *Client) SendBytes(ctx context.Context, inputBytes []byte) error {
	return c.blockchain.SendInput(ctx, inputBytes)
}

// Send a generic input to the DApp contract.
func (c *Client) SendGeneric(ctx context.Context, input any) error {
	inputBytes, err := EncodeGenericInput(input)
	if err != nil {
		return err
	}
	return c.SendBytes(ctx, inputBytes)
}

//
// Reader functions
//

// Wait until the DApp contract processes a given input.
// Returns the advance result of that input.
func (c *Client) WaitFor(ctx context.Context, inputIndex int) (*AdvanceResult, error) {
	for {
		input, err := c.reader.Input(ctx, inputIndex)
		if err != nil {
			if _, ok := err.(reader.NotFound); ok {
				goto wait
			}
			return nil, fmt.Errorf("faild to read input: %v", err)
		}
		if input.Status != reader.CompletionStatusUnprocessed {
			return newAdvanceResult(input), nil
		}
	wait:
		time.Sleep(time.Second)
	}
}

// Sync to the latest Dapp state.
// Return the updated slice of Advance results.
func (c *Client) Sync(ctx context.Context, results []*AdvanceResult) ([]*AdvanceResult, error) {
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
		results = append(results, newAdvanceResult(input))
		inputIndex++
	}
	return results, nil
}
