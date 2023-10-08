// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gligneul/eggroll/blockchain"
	"github.com/gligneul/eggroll/reader"
)

const envPrefix = "EGGROLL_"

// Load variable from env.
func loadVar(varName string, defaultValue string) string {
	varName = envPrefix + varName
	value := os.Getenv(varName)
	if value == "" {
		value = defaultValue
	}
	return value
}

// Configuration for the Client.
type ClientConfig struct {
	GraphqlEndpoint  string
	ProviderEndpoint string
}

// Load the config from environment variables.
func (c *ClientConfig) Load() {
	c.GraphqlEndpoint = loadVar("GRAPHQL_ENDPOINT", "http://localhost:8080/graphql")
	c.ProviderEndpoint = loadVar("ETH_RPC_ENDPOINT", "http://localhost:8545")
}

// Read the rollups state off chain.
// For more details, see the eggroll/reader package.
type readerAPI interface {
	Input(ctx context.Context, index int) (*reader.Input, error)
	Notice(ctx context.Context, inputIndex int, noticeIndex int) (*reader.Notice, error)
	Report(ctx context.Context, inputIndex int, reportIndex int) (*reader.Report, error)
	LastReports(ctx context.Context, last int) (*reader.Page[reader.Report], error)
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
	nextInput  int
}

// Create the Client loading the config from environment variables.
func NewClient() *Client {
	var config ClientConfig
	config.Load()
	return NewClientFromConfig(config)
}

// Create the Client with a custom config.
func NewClientFromConfig(config ClientConfig) *Client {
	return &Client{
		reader:     reader.NewGraphQLReader(config.GraphqlEndpoint),
		blockchain: blockchain.NewETHClient(config.ProviderEndpoint),
		nextInput:  0,
	}
}

// Send a generic inputs to the DApp contract.
// Returns an slice with each input index.
func (c *Client) SendGeneric(ctx context.Context, input any) error {
	inputBytes, err := EncodeGenericInput(input)
	if err != nil {
		return err
	}
	return c.blockchain.SendInput(ctx, inputBytes)
}

// Wait until the DApp back end processes a given input.
func (c *Client) WaitFor(ctx context.Context, inputIndex int) error {
	for {
		input, err := c.reader.Input(ctx, inputIndex)
		if err != nil {
			if _, ok := err.(reader.NotFound); ok {
				goto wait
			}
			return fmt.Errorf("faild to read input: %v", err)
		}
		if input.Status != reader.CompletionStatusUnprocessed {
			return nil
		}
	wait:
		time.Sleep(time.Second)
	}
}

// Sync to the latest Dapp state.
func (c *Client) Sync(ctx context.Context) error {
	for {
		input, err := c.reader.Input(ctx, c.nextInput)
		if err != nil {
			if _, ok := err.(reader.NotFound); ok {
				break
			}
			return fmt.Errorf("failed to read input: %v", err)
		}
		if input.Status == reader.CompletionStatusUnprocessed {
			break
		}
		if input.Status == reader.CompletionStatusAccepted {
			notice, err := c.reader.Notice(ctx, c.nextInput, 0)
			if err != nil {
				return fmt.Errorf("failed to read notice: %v", err)
			}
			c.state = notice.Payload
		}
		c.nextInput++
	}
	return nil
}

// Get a copy of the current DApp state.
func (c *Client) ReadState(state any) error {
	err := json.Unmarshal(c.state, &state)
	if err != nil {
		return fmt.Errorf("failed to parse state: %v", err)
	}
	return nil
}

// Get the last 20 entries of log from the DApp.
func (c *Client) Logs(ctx context.Context) ([]string, error) {
	return c.LogsTail(ctx, 20)
}

// Get the last N entries of logs from the DApp.
func (c *Client) LogsTail(ctx context.Context, n int) ([]string, error) {
	page, err := c.reader.LastReports(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %v", err)
	}

	var logs []string
	for _, report := range page.Nodes {
		logs = append(logs, string(report.Payload))
	}

	return logs, nil
}
