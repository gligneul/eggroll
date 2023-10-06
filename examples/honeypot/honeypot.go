// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package honeypot

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll"
)

// Owner of the honeypot that can withdraw all funds.
var owner common.Address

func init() {
	owner = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
}

type Withdraw struct{}

type Contract struct{}

func (c *Contract) Decoders() []eggroll.Decoder {
	return []eggroll.Decoder{
		eggroll.NewGenericDecoder[Withdraw](),
	}
}

func (c *Contract) Advance(env *eggroll.Env, input any) error {
	if env.Deposit() != nil {
		env.Logf("received deposit: %v\n", env.Deposit())
		return nil
	}

	if env.Sender() != owner {
		return fmt.Errorf("ignoring input from %v", env.Sender())
	}

	if _, ok := input.(*Withdraw); !ok {
		return fmt.Errorf("ignoring input: %v", input)
	}

	for _, src := range env.EtherAddresses() {
		if src != owner {
			balance := env.EtherBalanceOf(src)
			env.EtherTransfer(src, owner, &balance)
		}
	}

	balance := env.EtherBalanceOf(owner)
	if balance.IsZero() {
		return fmt.Errorf("nothing to withdraw")
	}

	err := env.EtherWithdraw(owner, &balance)
	if err != nil {
		return err
	}
	env.Logf("withdraw %v\n", balance.ToBig().String())

	return nil
}