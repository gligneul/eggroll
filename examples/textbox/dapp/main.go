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
	textbox.TextBox
}

func (c *TextBoxContract) Codecs() []eggroll.Codec {
	return textbox.Codecs()
}

func (c *TextBoxContract) Advance(env *eggroll.Env) (any, error) {
	switch input := env.DecodeInput().(type) {
	case *textbox.Clear:
		env.Logln("received input clear")
		c.TextBox.Value = ""
	case *textbox.Append:
		env.Logf("received input append with '%v'\n", input.Value)
		c.TextBox.Value += input.Value
	default:
		return nil, fmt.Errorf("invalid input: %v", input)
	}
	return &c.TextBox, nil
}

func main() {
	eggroll.Roll(&TextBoxContract{})
}
