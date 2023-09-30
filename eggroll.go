// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level, opinionated, lambda-based framework for Cartesi Rollups in Go
package eggroll

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

// Ethereum type
type (
	Address common.Address
	Hash    common.Hash
)

// Metadata from the input
type Metadata struct {
	Sender         Address
	BlockNumber    int64
	BlockTimestamp int64
}

// Interface to interact with the rollups environment
type Env interface {

	// Get the Metadata for the current input
	Metadata() *Metadata

	// Send a report for debugging purposes.
	// The reports will be available as logs for the front end client.
	Report(format string, a ...any)

	// Send a voucher
	Voucher(destination Address, payload []byte)
}

// Signature of the handler that advances the rollups state
type Handler[S, I any] func(Env, *S, *I) error

// Register a handler to a DApp
func Register[S, I any](dapp DApp[S], handler Handler[S, I]) {
	// This function needs to be defined outside of the DApp interface
	// because Go doesn't support template parameters in methods.
	// So, it is not possible to write DApp[S].Register[I](Handler[S, I]).
	inputType := reflect.TypeOf((*I)(nil))
	dapp.register(inputType, func(env Env, state *S, input any) error {
		concreteInput := input.(*I)
		return handler(env, state, concreteInput)
	})
}

// Interface to manage the handlers and run the DApp back end
type DApp[S any] interface {

	// Register a handler for the given input type
	register(inputType reflect.Type, handler internalHandler[S])

	// Start the DApp back end.
	// This function only returns if there is an error.
	Roll()
}

// Set up the DApp back end loading the config from environment variables
func SetupDApp[S any]() DApp[S] {
	var config DAppConfig
	config.Load()
	return SetupDAppWithConfig[S](config)
}

// Set up the DApp with a custom config
func SetupDAppWithConfig[S any](config DAppConfig) DApp[S] {
	rollups := &rollupsHttpApi{config.RollupsEndpoint}
	dapp := dapp[S]{
		rollups:  rollups,
		handlers: make(handlerMap[S]),
	}
	return &dapp
}

// Interface to interact with the DApp from the front end
type Client[S any] interface {

	// Send inputs to the DApp back end.
	// Returns an slice with each input index.
	Send(input ...any) ([]int, error)

	// Wait until the DApp back end processes a given input
	WaitFor(inputIndex int) error

	// Get a copy of the current DApp state
	State() *S

	// Get the reports for a given input
	LogsFrom(input int) (string, error)

	// Get the last n reports
	LogsTail(n int) (string, error)

	// Get the last 20 reports
	Logs() (string, error)
}

// Set up the front-end Client loading the config from environment variables
func SetupClient[S any]() Client[S] {
	var config ClientConfig
	config.Load()
	return SetupClientWithConfig[S](config)
}

// Set up the DApp with a custom config
func SetupClientWithConfig[S any](config ClientConfig) Client[S] {
	return &client[S]{}
}
