// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level, opinionated, lambda-based framework for Cartesi Rollups in Go.
package eggroll

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
)

// Configuration for the Contract.
type ContractConfig struct {
	RollupsEndpoint string
}

// Load the config from environment variables.
func (c *ContractConfig) Load() {
	var defaultEndpoint string
	if runtime.GOARCH == "riscv64" {
		defaultEndpoint = "http://localhost:5004"
	} else {
		defaultEndpoint = "http://localhost:8080/host-runner"
	}
	c.RollupsEndpoint = loadVar("ROLLUPS_HTTP_ENDPOINT", defaultEndpoint)
}

// Contract is the back-end for the rollups.
// It dispatch the input to the corresponding handler while it advances the
// rollups state.
type Contract[S any] struct {
	rollups  rollupsApi
	handlers handlerMap[S]
}

// Create the Contract loading the config from environment variables.
func NewContract[S any]() *Contract[S] {
	var config ContractConfig
	config.Load()
	return NewContractFromConfig[S](config)
}

// Create the Contract with a custom config.
func NewContractFromConfig[S any](config ContractConfig) *Contract[S] {
	rollups := &rollupsHttpApi{config.RollupsEndpoint}
	contract := &Contract[S]{
		rollups:  rollups,
		handlers: make(handlerMap[S]),
	}
	return contract
}

// Start the Contract back end.
// This function never returns and exits if there is an error.
func (d *Contract[S]) Roll() {
	var state S
	env := &Env{rollups: d.rollups}
	status := statusAccept

	for {
		var (
			payload []byte
			err     error
		)
		payload, env.metadata, err = d.rollups.finish(status)
		if err != nil {
			log.Fatalf("failed to send finish: %v\n", err)
		}

		if err = d.handlers.dispatch(env, &state, payload); err != nil {
			env.Logf("rejecting input: %v\n", err)
			status = statusReject
			continue
		}

		stateSnapshot, err := json.Marshal(&state)
		if err != nil {
			log.Fatalf("failed to create state snapshot: %v\n", err)
		}
		if err = d.rollups.sendNotice(stateSnapshot); err != nil {
			log.Fatalf("failed to send notice: %v\n", err)
		}
		status = statusAccept
	}
}

// Signature of the handler that advances the rollups state.
type Handler[S, I any] func(*Env, *S, *I) error

// Register a handler for a custom input to a Contract.
func Register[S, I any](d *Contract[S], handler Handler[S, I]) {
	// This function needs to be defined outside of the Contract interface
	// because Go doesn't support template parameters in methods.
	// So, it is not possible to write Contract[S].Register[I](Handler[S, I]).

	key, decoder := getInputKeyDecoder[I]()

	gHandler := func(env *Env, state *S, inputBytes []byte) error {
		input, err := decoder(inputBytes)
		if err != nil {
			return err
		}
		return handler(env, state, input)
	}

	d.handlers[key] = gHandler
}

// Internal handler that receives an encoded input.
type genericHandler[S any] func(*Env, *S, []byte) error

// Map a input key to its handler.
type handlerMap[S any] map[inputKey]genericHandler[S]

// Call the handler for the given input.
func (m handlerMap[S]) dispatch(env *Env, state *S, payload []byte) error {
	inputKey, inputBytes, err := splitInput(payload)
	if err != nil {
		return err
	}

	handler, ok := m[inputKey]
	if !ok {
		return fmt.Errorf("handler not found (%v)", hex.EncodeToString(inputKey[:]))
	}

	if err := handler(env, state, inputBytes); err != nil {
		return fmt.Errorf("input rejected: %v", err)
	}

	return nil
}
