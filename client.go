// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Configuration for the Client.
type ClientConfig struct {
	GraphqlEndpoint  string
	ProviderEndpoint string
}

// Load the config from environment variables.
func (c *ClientConfig) Load() {
	c.GraphqlEndpoint = loadVar("GRAPHQL_ENDPOINT", "http://localhost:8080/graphql")
	c.ProviderEndpoint = loadVar("ETH_ENDPOINT", "http://localhost:8545")
}

// Read the rollups state from the outside.
type readerApi interface {
	Input(ctx context.Context, index int) (*Input, error)
	Notice(ctx context.Context, inputIndex int, noticeIndex int) (*Notice, error)
	Report(ctx context.Context, inputIndex int, reportIndex int) (*Report, error)
	LastReports(ctx context.Context, last int) (*Page[Report], error)
}

// The Client interacts with the DApp from the outside.
type Client[S any] struct {
	reader    readerApi
	client    blockchainClient
	state     S
	nextInput int
}

// Create the Client loading the config from environment variables.
func NewClient[S any]() *Client[S] {
	var config ClientConfig
	config.Load()
	return NewClientFromConfig[S](config)
}

// Create the Client with a custom config.
func NewClientFromConfig[S any](config ClientConfig) *Client[S] {
	reader := NewGraphqlReader(config.GraphqlEndpoint)
	client := newEthClient(config.ProviderEndpoint)
	return &Client[S]{
		reader:    reader,
		client:    client,
		nextInput: 0,
	}
}

// Send inputs to the DApp back end.
// Returns an slice with each input index.
func (c *Client[S]) Send(ctx context.Context, input ...any) error {
	c.client.Send(ctx)
	return nil
}

// Wait until the DApp back end processes a given input.
func (c *Client[S]) WaitFor(ctx context.Context, inputIndex int) error {
	for {
		input, err := c.reader.Input(ctx, inputIndex)
		if err != nil {
			if _, ok := err.(NotFound); ok {
				goto wait
			}
			return fmt.Errorf("faild to read input: %v", err)
		}
		if input.Status != CompletionStatusUnprocessed {
			return nil
		}
	wait:
		time.Sleep(time.Second)
	}
}

// Get a copy of the current DApp state.
func (c *Client[S]) State(ctx context.Context) (*S, error) {
	// Sync to the latest state
	for {
		input, err := c.reader.Input(ctx, c.nextInput)
		if err != nil {
			if _, ok := err.(NotFound); ok {
				break
			}
			return nil, fmt.Errorf("failed to read input: %v", err)
		}
		if input.Status == CompletionStatusUnprocessed {
			break
		}
		if input.Status == CompletionStatusAccepted {
			notice, err := c.reader.Notice(ctx, c.nextInput, 0)
			if err != nil {
				return nil, fmt.Errorf("failed to read notice: %v", err)
			}
			err = json.Unmarshal(notice.Payload, &c.state)
			if err != nil {
				return nil, fmt.Errorf("failed to parse state: %v", err)
			}
		}
		c.nextInput++
	}

	return &c.state, nil
}

// Get the last 20 entries of log from the DApp.
func (c *Client[S]) Logs(ctx context.Context) ([]string, error) {
	return c.LogsTail(ctx, 20)
}

// Get the last N entries of logs from the DApp.
func (c *Client[S]) LogsTail(ctx context.Context, n int) ([]string, error) {
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
