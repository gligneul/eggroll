// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level, opinionated, lambda-based framework for Cartesi Rollups in Go.
package eggroll

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/rollups"
	"github.com/gligneul/eggroll/wallets"
)

// Configuration for the Contract.
type ContractConfig struct {
	RollupsEndpoint string
}

// Load the config from environment variables.
func (c *ContractConfig) Load() {
	var defaultEndpoint string
	if runtime.GOARCH == "riscv64" {
		defaultEndpoint = "http://127.0.0.1:5004"
	} else {
		defaultEndpoint = "http://localhost:8080/host-runner"
	}
	c.RollupsEndpoint = loadVar("ROLLUPS_HTTP_ENDPOINT", defaultEndpoint)
}

// Key that identifies the input type.
type InputKey [4]byte

// Decodes inputs from bytes to Go types.
type Decoder interface {

	// Get the input key.
	InputKey() InputKey

	// Try to decode the given input.
	Decode(inputBytes []byte) (any, error)
}

// State of the rollups contract.
type State interface {

	// Advance the state with the given input.
	// If the call returns an error, reject the input.
	Advance(env *Env, input any) error
}

// Contract that advances the dapp state.
type Contract interface {

	// Get the decoders required by the contract.
	Decoders() []Decoder

	// Advance the state with the given input.
	// If the call returns an error, reject the input.
	Advance(env *Env, input any) error
}

// Start the Cartesi Rollups for the contract.
// This function doesn't return and exits if there is an error.
func Roll(contract Contract) {
	var config ContractConfig
	config.Load()
	RollWithConfig(contract, config)
}

// Start the Cartesi Rollups for the contract with the given config.
// This function doesn't return and exits if there is an error.
func RollWithConfig(contract Contract, config ContractConfig) {
	rollupsAPI := rollups.NewRollupsHTTP(config.RollupsEndpoint)
	decoders := makeDecoderMap(contract.Decoders())
	status := rollups.FinishStatusAccept
	wallets := wallets.NewWallets()
	dispatcher := wallets.MakeDispatcher()
	env := &Env{
		wallets: wallets,
		rollups: rollupsAPI,
	}

	for {
		var (
			payload    []byte
			inputBytes []byte
			err        error
		)
		payload, env.metadata, err = rollupsAPI.Finish(status)
		if err != nil {
			log.Fatalf("failed to send finish: %v\n", err)
		}

		env.deposit, inputBytes, err = dispatcher.Dispatch(env.metadata.Sender, payload)
		if err != nil {
			env.Logf("malformed portal input: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}
		if inputBytes == nil {
			inputBytes = payload
		}

		input := decodeInput(decoders, env, payload)
		if err = contract.Advance(env, input); err != nil {
			env.Logf("rejecting: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}

		stateSnapshot, err := json.Marshal(contract)
		if err != nil {
			log.Fatalf("failed to create state snapshot: %v\n", err)
		}
		if err = rollupsAPI.SendNotice(stateSnapshot); err != nil {
			log.Fatalf("failed to send notice: %v\n", err)
		}
		status = rollups.FinishStatusAccept
	}
}

// Map the input key to the respective decoder.
func makeDecoderMap(decoders []Decoder) map[InputKey]Decoder {
	decodersMap := make(map[InputKey]Decoder)
	for _, decoder := range decoders {
		key := decoder.InputKey()
		_, ok := decodersMap[key]
		if ok {
			// Bug in the application configuration, so it is reasonable to panic
			log.Panicf("decoder conflict: %v\n", common.Bytes2Hex(key[:]))
		}
		decodersMap[key] = decoder
	}
	return decodersMap
}

// Try to decode the input, if it fails, return the original payload.
func decodeInput(decoders map[InputKey]Decoder, env *Env, payload []byte) any {
	if len(payload) < 4 {
		return payload
	}
	key := InputKey(payload[:4])
	inputBytes := payload[4:]
	decoder, ok := decoders[key]
	if !ok {
		return payload
	}
	input, err := decoder.Decode(inputBytes)
	if err != nil {
		env.Logf("failed to decode input: %v\n", err)
		return payload
	}
	return input
}
