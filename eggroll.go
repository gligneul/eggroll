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
	// This function needs to be defined outside of the register interface
	// because Go doesn't template parameters in methods.
	// It is not possible to have DApp[S].Register[I](Handler[S, I]).
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

	// Start the DApp backend.
	// This function only returns if there is an error.
	Roll()
}

// Interface to interact with the DApp from the front end
type Client interface {

	// Send inputs to the DApp backend.
	// Returns an slice with each input index.
	Send(input ...any) ([]int, error)

	// Wait until the DApp backend processes a given input
	WaitFor(inputIndex int) error

	// Read the DApp state
	Read(state any) error

	// Get the reports in form of a log
	Log() (string, error)
}
