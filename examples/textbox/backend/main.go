// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"log"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/examples/textbox"
)

// Redefine the types to make the example cleaner
type (
	Append textbox.Append
	Clear  textbox.Clear
	State  textbox.State
)

func main() {
	dapp := eggroll.SetupDApp[State]()

	eggroll.Register(dapp, func(_ eggroll.Env, state *State, _ *Clear) error {
		state.TextBox = ""
		return nil
	})

	eggroll.Register(dapp, func(_ eggroll.Env, state *State, input *Append) error {
		state.TextBox += input.Value
		return nil
	})

	log.Fatal(dapp.Roll())
}
