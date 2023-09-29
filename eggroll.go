// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Framework for Cartesi Rollups in Go
package eggroll

import (
	"github.com/ethereum/go-ethereum/common"
)

// Advance the rollups state inside the Cartesi Machine.
// The DApp developer should implement this interface.
type DAppBackend interface {

	// Process an advance state request.
	// Receives the rollups state and the input.
	// If returns an error, the input is rejected.
	Advance(s State) error
}

// Use the Rollups Node to read the DApp state from the outside.
// The eggroll framework provides an implementation for this interface.
type Inspector interface {

	// Get a value from the DApp state.
	// Returns an error if the value is not found.
	Get(key string, value any) error

	// Check whether a key is in the DApp state
	Has(key string) bool

	// Get the log report from the DApp
	Log() string
}

// Send an input to the DApp from the outside.
// The eggroll framework provides an implementation for this interface.
type InputSender interface {

	// Sends an input to the DApp
	Send(payload any) error
}

// Manipulate the state of the DApp
type State interface {

	// Get the input for the current advance
	Input(value any)

	// Get the Metadata for the current advance
	Metadata() Metadata

	// Check whether a key is in the DApp state
	Has(key string) bool

	// Get a value from the DApp state
	Get(key string, value any)

	// Set a value in the DApp state.
	// Once the DApp finishes an advance step, the value is available to the Inspector.
	Set(key string, value any)

	// Delete the key from the DApp state
	Delete(key string)

	// Send a report for debugging purposes
	Report(format string, a ...any)

	// Send a voucher
	Voucher(destination Address, payload []byte)
}

// Input to advance the DApp state
type Metadata struct {
}

// Ethereum address
type Address common.Address
