// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Shared types for the TextBox DApp.
package textbox

import (
	"fmt"

	"github.com/gligneul/eggroll"
)

// @cut

type (
	Append struct {
		Value string
	}
	Clear struct{}
)

type Contract struct {
	TextBox string
}

func (c *Contract) Decoders() []eggroll.Decoder {
	return []eggroll.Decoder{
		eggroll.NewGenericDecoder[Clear](),
		eggroll.NewGenericDecoder[Append](),
	}
}

func (c *Contract) Advance(env *eggroll.Env, input any) error {
	switch input := input.(type) {
	case *Clear:
		env.Logln("received input clear")
		c.TextBox = ""
	case *Append:
		env.Logf("received input append with '%v'\n", input.Value)
		c.TextBox += input.Value
	default:
		return fmt.Errorf("invalid input")
	}
	return nil
}
