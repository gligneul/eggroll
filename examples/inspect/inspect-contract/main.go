// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll/pkg/eggroll"
)

type Contract struct {
	eggroll.DefaultContract
}

func (c *Contract) Advance(env eggroll.Env) (any, error) {
	env.Logf("advance: %v", string(env.RawInput()))
	return env.RawInput(), nil
}

func (c *Contract) Inspect(env eggroll.EnvReader) (any, error) {
	env.Logf("inspect: %v", string(env.RawInput()))
	return env.RawInput(), nil
}

func main() {
	eggroll.Roll(&Contract{})
}
