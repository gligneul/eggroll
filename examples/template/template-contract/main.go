// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll"
)

type Contract struct {
	eggroll.DefaultContract
}

func (c *Contract) Advance(env eggroll.Env) (any, error) {
	input := env.RawInput()
	env.Logf("received: %v", string(input))
	return input, nil
}

func main() {
	eggroll.Roll(&Contract{})
}
