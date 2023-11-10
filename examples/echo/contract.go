// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

//go:generate go run github.com/gligneul/eggroll/cmd/eggroll schema gen

import (
	"github.com/gligneul/eggroll/pkg/eggroll"
)

type Contract struct {
}

func (c *Contract) AdvanceEcho(env eggroll.Env, value string) error {
	env.Report(EncodeEchoResponse(value))
	return nil
}

func (c *Contract) InspectEcho(env eggroll.EnvReader, value string) error {
	env.Report(EncodeEchoResponse(value))
	return nil
}

func main() {
	Roll(&Contract{})
}
