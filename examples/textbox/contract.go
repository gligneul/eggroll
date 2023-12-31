// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

//go:generate go run github.com/gligneul/eggroll/cmd/eggroll schema gen

import (
	"github.com/gligneul/eggroll/pkg/eggroll"
)

type Contract struct {
	state string
}

func (c *Contract) Clear(env eggroll.Env) error {
	c.state = ""
	env.Report(EncodeCurrentState(c.state))
	return nil
}

func (c *Contract) Append(env eggroll.Env, value string) error {
	c.state += value
	env.Report(EncodeCurrentState(c.state))
	return nil
}

func main() {
	Roll(&Contract{})
}
