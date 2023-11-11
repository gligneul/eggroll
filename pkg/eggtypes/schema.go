// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggtypes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// The Go types for message schema implement this interface.
type Encoder interface {

	// Encode a Go value into binary data.
	Encode() []byte
}

// EggRoll uses 4 bytes to identify the message kind when decoding it.
// This is based on the function selector of the Solidity ABI.
type ID [4]byte

// Decode binary data into a Go value.
type Decoder func([]any) (any, error)

// Schema for a message.
type MessageSchema struct {
	ID
	Kind string
	abi.Arguments
	Decoder
}

// Global set of schemas managed by EggRoll.
// We use a global variable to avoid passing a schema manager object around.
// Using a global is fine because schemas should be static.
var schemas struct {
	byID   map[ID]MessageSchema
	byKind map[string]MessageSchema
}

// Add a message schema to the schemas managed by EggRoll.
func AddSchema(schema MessageSchema) error {
	_, ok := schemas.byID[schema.ID]
	if ok {
		return fmt.Errorf("duplicate schema with id: %x", schema.ID)
	}
	_, ok = schemas.byKind[schema.Kind]
	if ok {
		return fmt.Errorf("duplicate schema with kind: %v", schema.Kind)
	}
	schemas.byID[schema.ID] = schema
	schemas.byKind[schema.Kind] = schema
	return nil
}

// Add a message schema to the schemas managed by EggRoll.
// Panic if an error occurs.
func MustAddSchema(schema MessageSchema) {
	err := AddSchema(schema)
	if err != nil {
		panic(err)
	}
}

// Get all schemas managed by EggRoll.
// The schemas are sorted by Kind.
func GetSchemas() (result []MessageSchema) {
	for _, s := range schemas.byKind {
		result = append(result, s)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Kind < result[j].Kind
	})
	return result
}

// Get the schema for the binary data.
func getSchema(data []byte) (MessageSchema, error) {
	if len(data) < 4 {
		return MessageSchema{}, fmt.Errorf("data doesn't contain the 4-byte ID")
	}
	id := ID(data[:4])
	if (len(data)-4)%32 != 0 {
		return MessageSchema{}, fmt.Errorf("improperly formatted data: %x", data)
	}
	schema, ok := schemas.byID[id]
	if !ok {
		return MessageSchema{}, fmt.Errorf("schema not found for ID: %x", id)
	}
	return schema, nil
}

// Decode binary data into a Go value.
func Decode(data []byte) (any, error) {
	schema, err := getSchema(data)
	if err != nil {
		return nil, err
	}
	values, err := schema.Arguments.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %v", err)
	}
	return schema.Decoder(values)
}

// Decode the binary data into a map.
// Return the schema kind.
func DecodeIntoMap(m map[string]any, data []byte) (string, error) {
	schema, err := getSchema(data)
	if err != nil {
		return "", err
	}
	err = schema.Arguments.UnpackIntoMap(m, data[4:])
	if err != nil {
		return "", err
	}
	return schema.Kind, nil
}

// Encode the JSON message into an ABI payload.
func EncodeFromMap(kind string, m map[string]any) ([]byte, error) {
	schema, ok := schemas.byKind[kind]
	if !ok {
		return nil, fmt.Errorf("schema not found for kind: %v", kind)
	}
	values := make([]any, len(schema.Arguments))
	for i, arg := range schema.Arguments {
		values[i] = m[arg.Name]
	}
	data, err := schema.Arguments.PackValues(values)
	if err != nil {
		return nil, err
	}
	return append(schema.ID[:], data...), nil
}

// Log message from a DApp contract.
type Log struct {
	Message string
}

// ID for the log message type.
var LogID ID

// Encode the log into binary data.
func EncodeLog(Message string) []byte {
	values := make([]any, 1)
	values[0] = Message
	data, err := _abi.Methods["log"].Inputs.PackValues(values)
	if err != nil {
		panic(fmt.Sprintf("failed to encode log: %v", err))
	}
	return append(LogID[:], data...)
}

// Encode the log to into binary data.
func (v Log) Encode() []byte {
	return EncodeLog(v.Message)
}

func _log_Decode(values []any) (any, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("wrong number of values")
	}
	var ok bool
	var v Log
	v.Message, ok = values[0].(string)
	if !ok {
		return nil, fmt.Errorf("failed to unpack log.message")
	}
	return v, nil
}

// Error message from the DApp contract.
type Error struct {
	Message string
}

// ID for the error message
var ErrorID ID

// Encode the error into binary data.
func EncodeError(Message string) []byte {
	values := make([]any, 1)
	values[0] = Message
	data, err := _abi.Methods["error"].Inputs.PackValues(values)
	if err != nil {
		panic(fmt.Sprintf("failed to encode error: %v", err))
	}
	return append(LogID[:], data...)
}

// Encode the error into binary data.
func (v Error) Encode() []byte {
	return EncodeError(v.Message)
}

func _error_Decode(values []any) (any, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("wrong number of values")
	}
	var ok bool
	var v Error
	v.Message, ok = values[0].(string)
	if !ok {
		return nil, fmt.Errorf("failed to unpack error.message")
	}
	return v, nil
}

const _JSON_ABI = `[
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
  },
  {
    "inputs": [
      {
	"internalType": "string",
	"name": "message",
	"type": "string"
      }
    ],
    "name": "error",
    "outputs": [],
    "stateMutability": "",
    "type": "function"
  }
]`

var _abi abi.ABI

func init() {
	schemas.byID = make(map[ID]MessageSchema)
	schemas.byKind = make(map[string]MessageSchema)

	var err error
	_abi, err = abi.JSON(strings.NewReader(_JSON_ABI))
	if err != nil {
		panic(fmt.Sprintf("failed to decode ABI: %v", err))
	}

	LogID = ID(_abi.Methods["log"].ID)
	MustAddSchema(MessageSchema{
		ID:        LogID,
		Kind:      "log",
		Arguments: _abi.Methods["log"].Inputs,
		Decoder:   _log_Decode,
	})

	ErrorID = ID(_abi.Methods["error"].ID)
	MustAddSchema(MessageSchema{
		ID:        ErrorID,
		Kind:      "error",
		Arguments: _abi.Methods["error"].Inputs,
		Decoder:   _error_Decode,
	})
}
