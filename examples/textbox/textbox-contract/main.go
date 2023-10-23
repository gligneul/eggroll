// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	eggroll2 "github.com/gligneul/eggroll/pkg/eggroll"
)

type Contract struct {
	eggroll2.DefaultContract
	TextBox
}

func (c *Contract) Codecs() []eggroll2.Codec {
	return Codecs()
}

func (c *Contract) Advance(env eggroll2.Env) (any, error) {
	switch input := env.DecodeInput().(type) {
	case *Clear:
		env.Log("received input clear")
		c.TextBox.Value = ""
	case *Append:
		env.Logf("received input append with '%v'\n", input.Value)
		c.TextBox.Value += input.Value
	default:
		return nil, fmt.Errorf("invalid input: %v", input)
	}
	return &c.TextBox, nil
}

func main() {
	eggroll2.Roll(&Contract{})
}
