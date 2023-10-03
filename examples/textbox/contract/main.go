// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"textbox"

	"github.com/gligneul/eggroll"
)

// Redefine the types to make the example cleaner.
type (
	InputAppend textbox.InputAppend
	InputClear  textbox.InputClear
	State       textbox.State
)

func clearHandler(env *eggroll.Env, state *State, _ *InputClear) error {
	env.Logln("received input clear")
	state.TextBox = ""
	return nil
}

func appendHandler(env *eggroll.Env, state *State, input *InputAppend) error {
	env.Logf("received input append with '%v'\n", input.Value)
	state.TextBox += input.Value
	return nil
}

func main() {
	contract := eggroll.NewContract[State]()
	eggroll.Register(contract, clearHandler)
	eggroll.Register(contract, appendHandler)
	contract.Roll()
}
