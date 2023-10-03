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

// Key that identifies the input type.
type inputKey [4]byte

// Encode the input into bytes.
func EncodeInput(input any) ([]byte, error) {
	inputType := reflect.TypeOf(input)
	key := genericInputKey(inputType)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to encode input: %v", err)
	}
	return append(key[:], inputBytes...), nil
}

// Split input payload into key and bytes.
func splitInput(payload []byte) (inputKey, []byte, error) {
	if len(payload) < 4 {
		msg := "invalid payload len (%v bytes)"
		return inputKey{}, nil, fmt.Errorf(msg, len(payload))
	}
	return inputKey(payload[:4]), payload[4:], nil
}

// Function that decodes the input given the bytes.
type inputDecoder[I any] func([]byte) (*I, error)

// Get the input key and the decoder function.
func getInputKeyDecoder[I any]() (inputKey, inputDecoder[I]) {
	inputType := reflect.TypeOf((*I)(nil))
	key := genericInputKey(inputType)
	decoder := func(inputBytes []byte) (*I, error) {
		input := reflect.New(inputType.Elem()).Interface()
		if err := json.Unmarshal(inputBytes, input); err != nil {
			return nil, fmt.Errorf("failed to decode input: %v", err)
		}
		return input.(*I), nil
	}
	return key, inputDecoder[I](decoder)
}

// Get the key for a generic input.
func genericInputKey(inputType reflect.Type) inputKey {
	// Check if inputType is pointer to struct
	if inputType.Kind() != reflect.Ptr {
		log.Panicf("input type must be a pointer; is %v\n", inputType)
	}
	inputType = inputType.Elem()
	if inputType.Kind() != reflect.Struct {
		log.Panicf("input type must be a struct pointer; is *%v\n", inputType)
	}

	// Use the first 4 bytes of the keccak of the Input type as the handler key.
	// We do this to be compatible with inputs that are ABI encoded.
	// See: https://docs.soliditylang.org/en/latest/abi-spec.html
	hash := crypto.Keccak256Hash([]byte(inputType.Name()))
	return inputKey(hash[:4])
}
