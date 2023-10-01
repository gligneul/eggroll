// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/examples/textbox"
)

// Redefine the types to make the example cleaner.
type (
	InputAppend textbox.InputAppend
	InputClear  textbox.InputClear
	State       textbox.State
)

func clearHandler(env *eggroll.Env, state *State, _ *InputClear) error {
	env.Report("received input clear")
	state.TextBox = ""
	return nil
}

func appendHandler(env *eggroll.Env, state *State, input *InputAppend) error {
	env.Report("received input append with '%v'", input.Value)
	state.TextBox += input.Value
	return nil
}

func main() {
	dapp := eggroll.NewDApp[State]()
	eggroll.Register(dapp, clearHandler)
	eggroll.Register(dapp, appendHandler)
	dapp.Roll()
}
