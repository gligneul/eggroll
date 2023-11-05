// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggtypes

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ID used to select the corresponding unpacker.
// This is the function selector in Solidity ABI.
type ID [4]byte

// Pack a Go value into ABI data.
type Packer interface {
	Pack() []byte
}

// Unpack an array of values into a Go value.
type Unpacker func([]any) (any, error)

// ABI encoding value.
type Encoding struct {
	ID
	Name string
	abi.Arguments
	Unpacker
}

var encodings struct {
	byID   map[ID]Encoding
	byName map[string]Encoding
}

// Add a new encoding to the EggRoll ABI.
func AddEncoding(encoding Encoding) error {
	_, ok := encodings.byID[encoding.ID]
	if ok {
		return fmt.Errorf("duplicate encoding with id: %x", encoding.ID)
	}
	_, ok = encodings.byName[encoding.Name]
	if ok {
		return fmt.Errorf("duplicate encoding with name: %v", encoding.Name)
	}
	encodings.byID[encoding.ID] = encoding
	encodings.byName[encoding.Name] = encoding
	return nil
}

// Add a new encoding to the EggRoll ABI.
// Panic if an error occurs.
func MustAddEncoding(encoding Encoding) {
	err := AddEncoding(encoding)
	if err != nil {
		panic(err)
	}
}

// This function should not be called directly; call the Pack method of the
// Packer value instead.
func PackValues(id ID, args ...interface{}) ([]byte, error) {
	encoding, ok := encodings.byID[id]
	if !ok {
		return nil, fmt.Errorf("encoding not found for ID: %x", id)
	}
	data, err := encoding.Arguments.Pack(args...)
	if err != nil {
		return nil, err
	}
	return append(id[:], data...), nil
}

// Unpack the data into a Go value.
func Unpack(data []byte) (any, error) {
	if (len(data)-4)%32 != 0 {
		return nil, fmt.Errorf("improperly formatted data: %x", data)
	}
	id := ID(data[:4])
	encoding, ok := encodings.byID[id]
	if !ok {
		return nil, fmt.Errorf("encoding not found for ID: %x", id)
	}
	values, err := encoding.Arguments.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack: %v", err)
	}
	return encoding.Unpacker(values)
}

// Log messages from a DApp contract.
type Log struct {
	Message string
}

// ID for the log message type.
// The ABI prototype is `log(string)`.
var LogID = ID([]byte{0x41, 0x30, 0x4f, 0xac})

// Pack the log to into an ABI payload.
func (l Log) Pack() []byte {
	payload, err := PackValues(LogID, l.Message)
	if err != nil {
		panic(fmt.Sprintf("failed to pack log: %v", err))
	}
	return payload
}

func _log_Unpack(values []any) (any, error) {
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
}

func init() {
	encodings.byID = make(map[ID]Encoding)
	encodings.byName = make(map[string]Encoding)

	const jsonAbi = `[
	  {
	    "inputs": [
	      {
		"internalType": "string",
		"name": "message",
		"type": "string"
	      }
	    ],
	    "name": "log",
	    "outputs": [],
	    "stateMutability": "",
	    "type": "function"
	  }
	]`
	abiInterface, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		panic(fmt.Sprintf("failed to decode ABI: %v", err))
	}
	MustAddEncoding(Encoding{
		ID:        LogID,
		Name:      "log",
		Arguments: abiInterface.Methods["log"].Inputs,
		Unpacker:  _log_Unpack,
	})
}
