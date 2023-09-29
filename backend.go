// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

// Define the handler as a function
type AdvanceHandler func(State) error

func (f AdvanceHandler) Advance(s State) error {
	return f(s)
}

// Run the main loop that executes the DApp for a function handler
func Roll(f func(State) error) {
	RollBackend(AdvanceHandler(f))
}

// Run main loop that execute the DApp
func RollBackend(d DAppBackend) {
}
