// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
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

// Get the input type key
func computeHandlerKey(inputType reflect.Type) handlerKey {
	// Check if inputType is pointer to struct
	if inputType.Kind() != reflect.Ptr {
		log.Panicf("inputType must be a pointer; is %v\n", inputType)
	}
	inputType = inputType.Elem()
	if inputType.Kind() != reflect.Struct {
		log.Panicf("*inputType must be a struct; is %v\n", inputType)
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(inputType.Name()))
	hash := hasher.Sum(nil)

	return handlerKey(hash)
}
