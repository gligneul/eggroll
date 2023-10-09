// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	"textbox"

	"github.com/gligneul/eggroll"
)

type TextBoxContract struct {
	eggroll.DefaultContract
	TextBox string
}

func (c *TextBoxContract) Decoders() []eggroll.Decoder {
	return []eggroll.Decoder{
		eggroll.NewGenericDecoder[textbox.Clear](),
		eggroll.NewGenericDecoder[textbox.Append](),
	}
}

func (c *TextBoxContract) Advance(env *eggroll.Env, input any) ([]byte, error) {
	switch input := input.(type) {
	case *textbox.Clear:
		env.Logln("received input clear")
		c.TextBox = ""
	case *textbox.Append:
		env.Logf("received input append with '%v'\n", input.Value)
		c.TextBox += input.Value
	default:
		return nil, fmt.Errorf("invalid input")
	}
	return []byte(c.TextBox), nil
}

func main() {
	eggroll.Roll(&TextBoxContract{})
}
