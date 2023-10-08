// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Key that identifies the input type.
type InputKey [4]byte

// Decodes inputs from bytes to Go types.
type Decoder interface {

	// Get the input key.
	InputKey() InputKey

	// Try to decode the given input.
	Decode(inputBytes []byte) (any, error)
}

// Map the input key to the respective decoder.
func makeDecoderMap(decoders []Decoder) map[InputKey]Decoder {
	decoderMap := make(map[InputKey]Decoder)
	for _, decoder := range decoders {
		key := decoder.InputKey()
		_, ok := decoderMap[key]
		if ok {
			// Bug in the contract configuration, so it is reasonable to panic
			panic(fmt.Sprintf("two decoders with same key '%v'", common.Bytes2Hex(key[:])))
		}
		decoderMap[key] = decoder
	}
	return decoderMap
}

// Try to decode the input, if it fails, return the original payload.
// If the decoder fails return an error.
func decodeInput(decoderMap map[InputKey]Decoder, payload []byte) (any, error) {
	if len(payload) < 4 {
		return payload, nil
	}
	key := InputKey(payload[:4])
	inputBytes := payload[4:]
	decoder, ok := decoderMap[key]
	if !ok {
		return payload, nil
	}
	input, err := decoder.Decode(inputBytes)
	if err != nil {
		return nil, err
	}
	return input, nil
}

// Generic decoder for Inputs of type I.
type GenericDecoder[I any] struct {
	inputType reflect.Type
}

// Create a new generic decoder.
func NewGenericDecoder[I any]() *GenericDecoder[I] {
	return &GenericDecoder[I]{
		inputType: reflect.TypeOf((*I)(nil)).Elem(),
	}
}

func (d *GenericDecoder[I]) InputKey() InputKey {
	return genericInputKey(d.inputType)
}

func (d *GenericDecoder[I]) Decode(inputBytes []byte) (any, error) {
	input := reflect.New(d.inputType).Interface()
	if err := json.Unmarshal(inputBytes, input); err != nil {
		return nil, fmt.Errorf("failed to decode input: %v", err)
	}
	return input, nil
}

// Encode the input into bytes.
func EncodeGenericInput(input any) ([]byte, error) {
	inputType := reflect.TypeOf(input)
	key := genericInputKey(inputType)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to encode input: %v", err)
	}
	return append(key[:], inputBytes...), nil
}

// Use the first 4 bytes of the keccak of the Input type as the handler key.
// This is inspired by the Ethereum ABI encoding.
// See: https://docs.soliditylang.org/en/latest/abi-spec.html
func genericInputKey(inputType reflect.Type) InputKey {
	// Check if inputType is struct
	if inputType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("input type must be a struct; is %v\n", inputType))
	}
	hash := crypto.Keccak256Hash([]byte(inputType.Name()))
	return InputKey(hash[:4])
}
