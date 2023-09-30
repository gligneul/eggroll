// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/examples/textbox"
)

// Redefine the types to make the example cleaner.
type (
	Append textbox.Append
	Clear  textbox.Clear
	State  textbox.State
)

func main() {
	dapp := eggroll.NewDApp[State]()

	eggroll.Register(dapp, func(env *eggroll.Env, state *State, _ *Clear) error {
		env.Report("received clear request")
		state.TextBox = ""
		return nil
	})

	eggroll.Register(dapp, func(env *eggroll.Env, state *State, input *Append) error {
		env.Report("received append request with '%v'", input.Value)
		state.TextBox += input.Value
		return nil
	})

	dapp.Roll()
}
