// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	"textbox"

	"github.com/gligneul/eggroll"
)

type (
	State       textbox.State
	InputClear  textbox.InputClear
	InputAppend textbox.InputAppend
)

// @cut

func (s *State) Advance(env *eggroll.Env, input any) error {
	switch input := input.(type) {
	case *InputClear:
		env.Logln("received input clear")
		s.TextBox = ""
	case *InputAppend:
		env.Logf("received input append with '%v'\n", input.Value)
		s.TextBox += input.Value
	default:
		return fmt.Errorf("invalid input")
	}
	return nil
}

func main() {
	contract := eggroll.NewContract(&State{})
	contract.AddDecoder(eggroll.NewGenericDecoder[InputClear]())
	contract.AddDecoder(eggroll.NewGenericDecoder[InputAppend]())
	contract.Roll()
}
