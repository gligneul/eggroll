// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"log"
	"os"
	"sync"
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
	Reader Reader
	state  clientState[S]
}

// Create the Client loading the config from environment variables.
func NewClient[S any]() *Client[S] {
	var config ClientConfig
	config.Load()
	return NewClientFromConfig[S](config)
}

// Create the Client with a custom config.
func NewClientFromConfig[S any](config ClientConfig) *Client[S] {
	reader := NewGraphqlReader(context.Background(), config.GraphqlEndpoint)
	return &Client[S]{
		Reader: reader,
	}
}

// Send inputs to the DApp back end.
// Returns an slice with each input index.
func (c *Client[S]) Send(input ...any) ([]int, error) {
	return nil, nil
}

// Wait until the DApp back end processes a given input.
func (c *Client[S]) WaitFor(inputIndex int) error {
	return nil
}

// Get a copy of the current DApp state.
func (c *Client[S]) State() *S {
	return nil
}

// Get the reports for a given input.
func (c *Client[S]) LogsFrom(input int) (string, error) {
	return "", nil
}

// Get the last n reports.
func (c *Client[S]) LogsTail(n int) (string, error) {
	return "", nil
}

// Get the last 20 reports.
func (c *Client[S]) Logs() (string, error) {
	return c.LogsTail(20)
}

// Internal state of the client.
type clientState[S any] struct {
	sync.RWMutex
	state S
}
