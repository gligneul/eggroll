// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
)

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
		log.Panicf("input type must be a struct; is %v\n", inputType)
	}
	hash := crypto.Keccak256Hash([]byte(inputType.Name()))
	return InputKey(hash[:4])
}
