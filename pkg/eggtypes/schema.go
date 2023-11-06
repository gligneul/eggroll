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

// Log messages from a DApp contract.
type Log struct {
	Message string
}

// ID for the log message type.
// The ABI prototype is `log(string)`.
var LogID = ID([]byte{0x41, 0x30, 0x4f, 0xac})

// Encode the log fields to into binary data.
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
func (l Log) Encode() []byte {
	return EncodeLog(l.Message)
}

func _log_Decode(values []any) (any, error) {
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
	MustAddSchema(MessageSchema{
		ID:        LogID,
		Kind:      "log",
		Arguments: _abi.Methods["log"].Inputs,
		Decoder:   _log_Decode,
	})
}
