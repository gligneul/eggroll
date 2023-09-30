// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"golang.org/x/crypto/sha3"
)

// Internal handler that doesn't require the input type
type internalHandler[S any] func(Env, *S, any) error

// Use the first 4 bytes of the keccak of the Input type as the handler key.
// We do this to be compatible with inputs that are ABI encoded.
// See: https://docs.soliditylang.org/en/latest/abi-spec.html
type handlerKey [4]byte

type handlerEntry[S any] struct {
	inputType reflect.Type
	handler   internalHandler[S]
}

type handlerMap[S any] map[handlerKey]handlerEntry[S]

// Register the handler for the given input
func (m handlerMap[S]) register(inputType reflect.Type, handler internalHandler[S]) {
	key := computeHandlerKey(inputType)
	m[key] = handlerEntry[S]{inputType, handler}
}

// Call the handler for the given input
func (m handlerMap[S]) dispatch(env Env, state *S, payload []byte) error {
	if len(payload) < 4 {
		return fmt.Errorf("invalid payload len (%v bytes)", len(payload))
	}

	// The first 4 bytes are the handler key, the rest is the input
	keyBytes := payload[:4]
	inputBytes := payload[4:]

	entry, ok := m[handlerKey(keyBytes)]
	if !ok {
		return fmt.Errorf("handler not found (%v)", hex.EncodeToString(keyBytes))
	}

	input := reflect.New(entry.inputType.Elem()).Interface()
	if err := json.Unmarshal(inputBytes, input); err != nil {
		return fmt.Errorf("failed to decode input: %v", err)
	}

	if err := entry.handler(env, state, input); err != nil {
		return fmt.Errorf("input rejected: %v", err)
	}

	return nil
}

// Get the input type key
func computeHandlerKey(inputType reflect.Type) handlerKey {
	// Check if inputType is pointer to struct
	if inputType.Kind() != reflect.Ptr {
		log.Panicf("input type must be a pointer; is %v\n", inputType)
	}
	inputType = inputType.Elem()
	if inputType.Kind() != reflect.Struct {
		log.Panicf("input type must be a struct pointer; is *%v\n", inputType)
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(inputType.Name()))
	hash := hasher.Sum(nil)

	return handlerKey(hash)
}
