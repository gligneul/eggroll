// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll/pkg/eggroll"
)

type Contract struct{}

func (c *Contract) Advance(env eggroll.Env, input []byte) error {
	env.Logf("received advance: %v", string(input))
	env.Report(input)
	return nil
}

func (c *Contract) Inspect(env eggroll.EnvReader, input []byte) error {
	env.Logf("received inspect: %v", string(input))
	env.Report(input)
	return nil
}

func main() {
	eggroll.Roll(&Contract{})
}
