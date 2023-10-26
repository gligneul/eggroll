// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"

	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/gligneul/eggroll/pkg/eggwallets"

	"github.com/ethereum/go-ethereum/common"
)

// Owner of the honeypot that can withdraw all funds.
var Owner common.Address

func init() {
	Owner = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
}

type Contract struct{}

func (c *Contract) Advance(env eggroll.Env, input []byte) error {
	unpacked, err := eggtypes.Unpack(input)
	if err != nil {
		return err
	}
	switch input := unpacked.(type) {
	case Deposit:
		switch deposit := env.Deposit().(type) {
		case *eggwallets.EtherDeposit:
			env.Log(deposit)
			if env.Sender() != Owner {
				// Transfer Ether deposits to Owner
				env.EtherTransfer(env.Sender(), Owner, deposit.Value)
			}
			reportHoneypot(env)
			return nil

		default:
			return fmt.Errorf("unsupported deposit: %T", deposit)
		}

	case Withdraw:
		if env.Sender() != Owner {
			// Ignore inputs that are not from Owner
			return fmt.Errorf("ignoring input from %v", env.Sender())
		}
		_, err := env.EtherWithdraw(Owner, input.Value)
		if err != nil {
			return err
		}
		env.Logf("withdraw %v\n", input.Value)
		reportHoneypot(env)
		return nil

	default:
		return fmt.Errorf("unknown input: %T", input)
	}
}

func (c *Contract) Inspect(env eggroll.EnvReader, input []byte) error {
	reportHoneypot(env)
	return nil
}

func reportHoneypot(env eggroll.EnvReader) {
	honeypot := Honeypot{
		Balance: env.EtherBalanceOf(Owner),
	}
	env.Report(honeypot.Pack())
}

func main() {
	eggroll.Roll(&Contract{})
}
