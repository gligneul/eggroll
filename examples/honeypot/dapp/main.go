// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/wallets"

	"honeypot"
)

type HoneypotContract struct {
	eggroll.DefaultContract
}

func (c HoneypotContract) Decoders() []eggroll.Decoder {
	return []eggroll.Decoder{
		eggroll.NewGenericDecoder[honeypot.Withdraw](),
	}
}

func (c *HoneypotContract) Advance(env *eggroll.Env, input any) ([]byte, error) {
	// Handle Ether deposits
	if deposit, ok := env.Deposit().(*wallets.EtherDeposit); ok {
		env.Logf("received deposit: %v\n", deposit)
		if env.Sender() != honeypot.Owner {
			// Transfer Ether deposits to honeypot.Owner
			env.EtherTransfer(env.Sender(), honeypot.Owner, &deposit.Value)
		}
		ownerBalance := env.EtherBalanceOf(honeypot.Owner)
		return ownerBalance.Bytes(), nil
	}

	// Ignore inputs that are not from honeypot.Owner
	if env.Sender() != honeypot.Owner {
		return nil, fmt.Errorf("ignoring input from %v", env.Sender())
	}

	// Handle honeypot.Owner withdraw request
	switch input := input.(type) {
	case *honeypot.Withdraw:
		_, err := env.EtherWithdraw(honeypot.Owner, input.Value)
		if err != nil {
			return nil, err
		}
		env.Logf("withdraw %v\n", input.Value.ToBig().String())
		ownerBalance := env.EtherBalanceOf(honeypot.Owner)
		return ownerBalance.Bytes(), nil

	default:
		return nil, fmt.Errorf("invalid input")
	}
}

func main() {
	eggroll.Roll(&HoneypotContract{})
}
