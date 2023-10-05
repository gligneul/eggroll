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

// Interface with the Rollups API.
// We don't expose this API because calling it directly will break EggRoll assumptions.
type rollupsAPI interface {
	SendVoucher(destination common.Address, payload []byte) error
	SendNotice(payload []byte) error
	SendReport(payload []byte) error
	Finish(status rollups.FinishStatus) ([]byte, *rollups.Metadata, error)
}

// State of the rollups contract.
type State interface {

	// Advance the state with the given input.
	// If the call returns an error, reject the input.
	Advance(env *Env, input any) error
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

// Receive inputs and advances the rollups state.
type Contract struct {
	rollups  rollupsAPI
	state    State
	decoders map[InputKey]Decoder
}

// Create the Contract loading the config from environment variables.
func NewContract(state State) *Contract {
	var config ContractConfig
	config.Load()
	return NewContractFromConfig(state, config)
}

// Create the Contract with a custom config.
func NewContractFromConfig(state State, config ContractConfig) *Contract {
	rollups := rollups.NewRollupsHTTP(config.RollupsEndpoint)
	contract := &Contract{
		rollups,
		state,
		make(map[InputKey]Decoder),
	}
	return contract
}

// Add a decoder to the contract.
func (c *Contract) AddDecoder(decoder Decoder) {
	key := decoder.InputKey()
	_, ok := c.decoders[key]
	if ok {
		log.Panicf("decoder conflict: %v\n", common.Bytes2Hex(key[:]))
	}
	c.decoders[key] = decoder
}

// Start to advance the rollups state.
// This function never returns and exits if there is an error.
func (c *Contract) Roll() {
	status := rollups.FinishStatusAccept
	env := &Env{rollups: c.rollups}

	// TODO set dapp address in env

	env.wallets = wallets.NewWallets()
	dispatcher := env.wallets.MakeDispatcher()

	for {
		var (
			payload    []byte
			inputBytes []byte
			err        error
		)
		payload, env.metadata, err = c.rollups.Finish(status)
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

		input := c.decode(env, payload)
		if err = c.state.Advance(env, input); err != nil {
			env.Logf("rejecting: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}

		stateSnapshot, err := json.Marshal(&c.state)
		if err != nil {
			log.Fatalf("failed to create state snapshot: %v\n", err)
		}
		if err = c.rollups.SendNotice(stateSnapshot); err != nil {
			log.Fatalf("failed to send notice: %v\n", err)
		}
		status = rollups.FinishStatusAccept
	}
}

// Try to decode the input, if it fails, return the original payload.
func (c *Contract) decode(env *Env, payload []byte) any {
	if len(payload) < 4 {
		return payload
	}
	key := InputKey(payload[:4])
	inputBytes := payload[4:]
	decoder, ok := c.decoders[key]
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
