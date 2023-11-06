// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

//go:generate go run github.com/gligneul/eggroll/cmd/eggroll schema gen

import (
	"fmt"

	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggtypes"
)

type Contract struct {
	TextBox
}

func (c *Contract) Advance(env eggroll.Env, input []byte) error {
	unpacked, err := eggtypes.Decode(input)
	if err != nil {
		return err
	}
	switch input := unpacked.(type) {
	case Clear:
		env.Log("received input clear")
		c.TextBox.Value = ""
		env.Report(c.TextBox.Encode())
		return nil
	case Append:
		env.Logf("received input append with '%v'\n", input.Value)
		c.TextBox.Value += input.Value
		env.Report(c.TextBox.Encode())
		return nil
	default:
		return fmt.Errorf("unknown input: %T", input)
	}
}

func (c *Contract) Inspect(env eggroll.EnvReader, input []byte) error {
	env.Report(c.TextBox.Encode())
	return nil
}

func main() {
	eggroll.Roll(&Contract{})
}
