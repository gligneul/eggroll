// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gligneul/eggroll/pkg/eggtypes"
)

type Deposit struct {
}

type Withdraw struct {
	Value *big.Int
}

type Honeypot struct {
	Balance *big.Int
}

func (v Deposit) Pack() []byte {
	payload, err := eggtypes.Pack("Deposit")
	if err != nil {
		panic(fmt.Sprintf("failed to pack Withdraw: %v", err))
	}
	return payload
}

func (v Withdraw) Pack() []byte {
	payload, err := eggtypes.Pack("Withdraw", v.Value)
	if err != nil {
		panic(fmt.Sprintf("failed to pack Withdraw: %v", err))
	}
	return payload
}

func (v Honeypot) Pack() []byte {
	payload, err := eggtypes.Pack("Honeypot", v.Balance)
	if err != nil {
		panic(fmt.Sprintf("failed to pack Honeypot: %v", err))
	}
	return payload
}

var DepositID [4]byte
var WithdrawID [4]byte
var HoneypotID [4]byte

func init() {
	abiJson := `
	[
	  {
	    "inputs": [],
	    "name": "Deposit",
	    "outputs": [],
	    "stateMutability": "",
	    "type": "function"
	  },
	  {
	    "inputs": [
	      {
		"internalType": "uint256",
		"name": "Value",
		"type": "uint256"
	      }
	    ],
	    "name": "Withdraw",
	    "outputs": [
	      {
		"internalType": "uint256",
		"name": "Value",
		"type": "uint256"
	      }
	    ],
	    "stateMutability": "",
	    "type": "function"
	  },
	  {
	    "inputs": [
	      {
		"internalType": "uint256",
		"name": "Balance",
		"type": "uint256"
	      }
	    ],
	    "name": "Honeypot",
	    "outputs": [
	      {
		"internalType": "uint256",
		"name": "Balance",
		"type": "uint256"
	      }
	    ],
	    "stateMutability": "",
	    "type": "function"
	  }
	]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		panic(fmt.Sprintf("failed to decode ABI: %v", err))
	}

	DepositID = [4]byte(abiInterface.Methods["Deposit"].ID)
	WithdrawID = [4]byte(abiInterface.Methods["Withdraw"].ID)
	HoneypotID = [4]byte(abiInterface.Methods["Honeypot"].ID)

	eggtypes.AddMethod(abiInterface.Methods["Deposit"], func(values []any) (any, error) {
		if len(values) != 0 {
			return nil, fmt.Errorf("wrong number of values")
		}
		var v Deposit
		return v, nil
	})
	eggtypes.AddMethod(abiInterface.Methods["Withdraw"], func(values []any) (any, error) {
		if len(values) != 1 {
			return nil, fmt.Errorf("wrong number of values")
		}
		var ok bool
		var v Withdraw
		v.Value, ok = values[0].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("failed to unpack Withdraw.Value")
		}
		return v, nil
	})
	eggtypes.AddMethod(abiInterface.Methods["Honeypot"], func(values []any) (any, error) {
		if len(values) != 1 {
			return nil, fmt.Errorf("wrong number of values")
		}
		var ok bool
		var v Honeypot
		v.Balance, ok = values[0].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("failed to unpack Honeypot.Balance")
		}
		return v, nil
	})
}
