// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

// Rollups input from the Reader API.
type Input struct {
	Index       int
	Status      CompletionStatus
	BlockNumber int64
}

// Read the rollups state from the outside.
type Reader interface {

	// Get a completion status for the given input.
	Input(index int) (*Input, error)
}
