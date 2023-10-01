// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Configuration for the Client.
type ClientConfig struct {
	GraphqlEndpoint string
}

// Load the config from environment variables.
func (c *ClientConfig) Load() {
	varName := "EGGROLL_GRAPHQL_ENDPOINT"
	c.GraphqlEndpoint = os.Getenv(varName)
	if c.GraphqlEndpoint == "" {
		c.GraphqlEndpoint = "http://localhost:8080/graphql"
	}
	log.Printf("set %v=%v\n", varName, c.GraphqlEndpoint)
}

// The Client interacts with the DApp from the outside.
type Client[S any] struct {
	Reader    Reader
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
	return &Client[S]{
		Reader:    reader,
		nextInput: 0,
	}
}

// Send inputs to the DApp back end.
// Returns an slice with each input index.
func (c *Client[S]) Send(input ...any) ([]int, error) {
	return nil, nil
}

// Wait until the DApp back end processes a given input.
func (c *Client[S]) WaitFor(ctx context.Context, inputIndex int) error {
	for {
		input, err := c.Reader.Input(ctx, inputIndex)
		if err != nil {
			if strings.HasSuffix(err.Error(), "not found\n") {
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
		input, err := c.Reader.Input(ctx, c.nextInput)
		if err != nil {
			if strings.HasSuffix(err.Error(), "not found\n") {
				break
			}
			return nil, fmt.Errorf("failed to read input: %v", err)
		}
		if input.Status == CompletionStatusUnprocessed {
			break
		}
		if input.Status == CompletionStatusAccepted {
			notice, err := c.Reader.Notice(ctx, c.nextInput, 0)
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

// Get the logs for a given input.
func (c *Client[S]) LogsFrom(input int) (string, error) {
	return "", nil
}

// Get the last n lines of logs from the DApp.
func (c *Client[S]) LogsTail(n int) (string, error) {
	return "", nil
}

// Get the last 20 lines of log.
func (c *Client[S]) Logs() (string, error) {
	return c.LogsTail(20)
}
