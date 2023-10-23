// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	eggroll2 "github.com/gligneul/eggroll/pkg/eggroll"
	wallets2 "github.com/gligneul/eggroll/pkg/eggwallets"

	"github.com/ethereum/go-ethereum/common"
)

// Owner of the honeypot that can withdraw all funds.
var Owner common.Address

func init() {
	Owner = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
}

type Contract struct {
	eggroll2.DefaultContract
}

func (c Contract) Codecs() []eggroll2.Codec {
	return Codecs()
}

func (c *Contract) Advance(env eggroll2.Env) (any, error) {
	if deposit := env.Deposit(); deposit != nil {
		return c.handleDeposit(env, deposit)
	}

	return c.handleInput(env, env.DecodeInput())
}

func (c *Contract) handleDeposit(env eggroll2.Env, deposit wallets2.Deposit) (any, error) {
	switch deposit := env.Deposit().(type) {
	case *wallets2.EtherDeposit:
		env.Logf("received deposit: %v\n", deposit)
		if env.Sender() != Owner {
			// Transfer Ether deposits to Owner
			env.EtherTransfer(env.Sender(), Owner, &deposit.Value)
		}
		return c.getBalance(env), nil

	default:
		return nil, fmt.Errorf("unsupported deposit: %v", deposit)
	}
}

func (c *Contract) handleInput(env eggroll2.Env, input any) (any, error) {
	if env.Sender() != Owner {
		// Ignore inputs that are not from Owner
		return nil, fmt.Errorf("ignoring input from %v", env.Sender())
	}

	switch input := input.(type) {
	case *Withdraw:
		fmt.Printf(">> %#v\n", input)
		_, err := env.EtherWithdraw(Owner, input.Value)
		if err != nil {
			return nil, err
		}
		env.Logf("withdraw %v\n", input.Value.ToBig().String())
		return c.getBalance(env), nil

	default:
		return nil, fmt.Errorf("invalid input: %v", input)
	}
}

func (c *Contract) getBalance(env eggroll2.Env) *Honeypot {
	ownerBalance := env.EtherBalanceOf(Owner)
	return &Honeypot{
		Balance: &ownerBalance,
	}
}

func main() {
	eggroll2.Roll(&Contract{})
}
