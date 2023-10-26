// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gligneul/eggroll/pkg/eggtypes"
)

type Append struct {
	Value string
}

type Clear struct{}

type TextBox struct {
	Value string
}

func (v Append) Pack() []byte {
	payload, err := eggtypes.Pack("Append", v.Value)
	if err != nil {
		panic(fmt.Sprintf("failed to pack Append: %v", err))
	}
	return payload
}

func (v Clear) Pack() []byte {
	payload, err := eggtypes.Pack("Clear")
	if err != nil {
		panic(fmt.Sprintf("failed to pack Clear: %v", err))
	}
	return payload
}

func (v TextBox) Pack() []byte {
	payload, err := eggtypes.Pack("TextBox", v.Value)
	if err != nil {
		panic(fmt.Sprintf("failed to pack TextBox: %v", err))
	}
	return payload
}

var AppendID [4]byte
var ClearID [4]byte
var TextBoxID [4]byte

func init() {
	abiJson := `
	[
	  {
	    "inputs": [
	      {
		"internalType": "string",
		"name": "Value",
		"type": "string"
	      }
	    ],
	    "name": "Append",
	    "outputs": [],
	    "stateMutability": "",
	    "type": "function"
	  },
	  {
	    "inputs": [],
	    "name": "Clear",
	    "outputs": [],
	    "stateMutability": "",
	    "type": "function"
	  },
	  {
	    "inputs": [
	      {
		"internalType": "string",
		"name": "Value",
		"type": "string"
	      }
	    ],
	    "name": "TextBox",
	    "outputs": [],
	    "stateMutability": "",
	    "type": "function"
	  }
	]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		panic(fmt.Sprintf("failed to decode ABI: %v", err))
	}

	AppendID = [4]byte(abiInterface.Methods["Append"].ID)
	ClearID = [4]byte(abiInterface.Methods["Clear"].ID)
	TextBoxID = [4]byte(abiInterface.Methods["TextBox"].ID)

	eggtypes.AddMethod(abiInterface.Methods["Append"], func(values []any) (any, error) {
		if len(values) != 1 {
			return nil, fmt.Errorf("wrong number of values")
		}
		var ok bool
		var v Append
		v.Value, ok = values[0].(string)
		if !ok {
			return nil, fmt.Errorf("failed to unpack Append.Value")
		}
		return v, nil
	})
	eggtypes.AddMethod(abiInterface.Methods["Clear"], func(values []any) (any, error) {
		if len(values) != 0 {
			return nil, fmt.Errorf("wrong number of values")
		}
		var v Clear
		return v, nil
	})
	eggtypes.AddMethod(abiInterface.Methods["TextBox"], func(values []any) (any, error) {
		if len(values) != 1 {
			return nil, fmt.Errorf("wrong number of values")
		}
		var ok bool
		var v TextBox
		v.Value, ok = values[0].(string)
		if !ok {
			return nil, fmt.Errorf("failed to unpack TextBox.Value")
		}
		return v, nil
	})
}
