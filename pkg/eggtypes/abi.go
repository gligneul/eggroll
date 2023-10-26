// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggtypes

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// Pack a Go value into ABI data.
type Packer interface {
	Pack() []byte
}

// Receive an array of values from the ABI.Unpack method and return the
// corresponding Go value.
type Unpacker func([]any) (any, error)

var globalMetadata struct {
	abi       abi.ABI
	unpackers map[string]Unpacker
}

// Add a new ABI method with an unpacker to the EggRoll ABI.
func AddMethod(method abi.Method, unpacker Unpacker) {
	if globalMetadata.unpackers == nil {
		globalMetadata.unpackers = make(map[string]Unpacker)
	}
	if globalMetadata.abi.Methods == nil {
		globalMetadata.abi.Methods = make(map[string]abi.Method)
	}
	found, _ := globalMetadata.abi.MethodById(method.ID)
	if found != nil {
		panic(fmt.Errorf("method already registred: %#v", method))
	}
	globalMetadata.abi.Methods[method.Name] = method
	globalMetadata.unpackers[method.Name] = unpacker
}

// Pack the given method name to conform the ABI.
// For more info, see abi.ABI.Pack.
func Pack(name string, args ...interface{}) ([]byte, error) {
	return globalMetadata.abi.Pack(name, args...)
}

// Unpack the data into a Go value.
func Unpack(data []byte) (any, error) {
	if (len(data)-4)%32 != 0 {
		return nil, fmt.Errorf("improperly formatted output: %x", data)
	}
	method, err := globalMetadata.abi.MethodById(data)
	if err != nil {
		return nil, err
	}
	values, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack: %v", err)
	}
	unpacker := globalMetadata.unpackers[method.Name]
	return unpacker(values)
}

// Log messages from a DApp contract.
type Log struct {
	Message string
}

func (l Log) Pack() []byte {
	payload, err := Pack("Log", l.Message)
	if err != nil {
		panic(fmt.Sprintf("failed to pack log: %v", err))
	}
	return payload
}

var LogID [4]byte

func init() {
	abiJson := `
	[
	  {
	    "inputs": [
	      {
		"internalType": "string",
		"name": "message",
		"type": "string"
	      }
	    ],
	    "name": "Log",
	    "outputs": [],
	    "stateMutability": "",
	    "type": "function"
	  }
	]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		panic(fmt.Sprintf("failed to decode ABI: %v", err))
	}
	LogID = [4]byte(abiInterface.Methods["Log"].ID)
	AddMethod(abiInterface.Methods["Log"], func(values []any) (any, error) {
		if len(values) != 1 {
			return nil, fmt.Errorf("wrong number of values")
		}
		var ok bool
		var log Log
		log.Message, ok = values[0].(string)
		if !ok {
			return nil, fmt.Errorf("failed to unpack log.Payload")
		}
		return log, nil
	})
}
