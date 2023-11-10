// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

//go:generate go run github.com/gligneul/eggroll/cmd/eggroll schema gen

import (
	"fmt"
	"math/big"

	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggwallets"

	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	owner common.Address
}

func (c *Contract) Deposit(env eggroll.Env) error {
	switch deposit := env.Deposit().(type) {
	case *eggwallets.EtherDeposit:
		env.Log(deposit)
		if env.Sender() != c.owner {
			env.EtherTransfer(env.Sender(), c.owner, deposit.Value)
		}
		env.Report(EncodeCurrentBalance(env.EtherBalanceOf(c.owner)))
		return nil
	default:
		return fmt.Errorf("unsupported deposit: %T", deposit)
	}
}

func (c *Contract) Withdraw(env eggroll.Env, value *big.Int) error {
	if env.Sender() != c.owner {
		return fmt.Errorf("ignoring input from %v", env.Sender())
	}
	_, err := env.EtherWithdraw(c.owner, value)
	if err != nil {
		return err
	}
	env.Logf("withdrawn %v\n", value)
	env.Report(EncodeCurrentBalance(env.EtherBalanceOf(c.owner)))
	return nil
}

func main() {
	Roll(&Contract{common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")})
}
