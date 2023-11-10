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

// Owner of the honeypot that can withdraw all funds.
var Owner common.Address

func init() {
	Owner = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
}

type Contract struct{}

func (c *Contract) Deposit(env eggroll.Env) error {
	switch deposit := env.Deposit().(type) {
	case *eggwallets.EtherDeposit:
		env.Log(deposit)
		if env.Sender() != Owner {
			env.EtherTransfer(env.Sender(), Owner, deposit.Value)
		}
		env.Report(EncodeCurrentBalance(env.EtherBalanceOf(Owner)))
		return nil
	default:
		return fmt.Errorf("unsupported deposit: %T", deposit)
	}
}

func (c *Contract) Withdraw(env eggroll.Env, value *big.Int) error {
	if env.Sender() != Owner {
		return fmt.Errorf("ignoring input from %v", env.Sender())
	}
	_, err := env.EtherWithdraw(Owner, value)
	if err != nil {
		return err
	}
	env.Logf("withdrawn %v\n", value)
	env.Report(EncodeCurrentBalance(env.EtherBalanceOf(Owner)))
	return nil
}

func main() {
	Roll(&Contract{})
}
