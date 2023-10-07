// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level, opinionated, lambda-based framework for Cartesi Rollups in Go.
package eggroll

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/blockchain"
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

// Contract that advances the dapp state.
type Contract interface {

	// Get the decoders required by the contract.
	Decoders() []Decoder

	// Advance the state with the given input.
	// If the call returns an error, reject the input.
	Advance(env *Env, input any) error
}

// DefaultContract provides a default implementation for optional contract methods.
type DefaultContract struct{}

func (_ *DefaultContract) Decoders() []Decoder {
	return nil
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
	env := &Env{
		etherWallet: wallets.NewEtherWallet(),
		rollups:     rollups.NewRollupsHTTP(config.RollupsEndpoint),
	}

	decoderMap := makeDecoderMap(contract.Decoders())
	walletMap := map[common.Address]wallets.Wallet{
		blockchain.AddressEtherPortal: env.etherWallet,
	}

	status := rollups.FinishStatusAccept

	for {
		payload, metadata, err := env.rollups.Finish(status)
		if err != nil {
			log.Fatalf("failed to send finish: %v\n", err)
		}

		err = handleAdvance(env, contract, decoderMap, walletMap, payload, metadata)
		if err != nil {
			env.Logf("rejecting: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}

		stateSnapshot, err := json.Marshal(contract)
		if err != nil {
			log.Fatalf("failed to create state snapshot: %v\n", err)
		}
		if err = env.rollups.SendNotice(stateSnapshot); err != nil {
			log.Fatalf("failed to send notice: %v\n", err)
		}
		status = rollups.FinishStatusAccept
	}
}

// Handle the advance request from the rollups API.
func handleAdvance(
	env *Env,
	contract Contract,
	decoderMap map[InputKey]Decoder,
	walletMap map[common.Address]wallets.Wallet,
	payload []byte,
	metadata *rollups.Metadata,
) (err error) {
	var deposit wallets.Deposit
	var inputBytes []byte

	if metadata.Sender == blockchain.AddressDAppAddressRelay {
		return handleDAppAddressRelay(env, payload)
	}

	wallet, ok := walletMap[metadata.Sender]
	if ok {
		deposit, inputBytes, err = wallet.Deposit(payload)
		if err != nil {
			return fmt.Errorf("malformed portal input: %v", err)
		}
	} else {
		deposit = nil
		inputBytes = payload
	}

	input, err := decodeInput(decoderMap, inputBytes)
	if err != nil {
		return err
	}

	// set env variables before decoding calling contract.Advance
	env.metadata = metadata
	env.deposit = deposit

	if err = contract.Advance(env, input); err != nil {
		return err
	}

	return nil
}

// Handle the input from the DAppAddressRelay.
func handleDAppAddressRelay(env *Env, payload []byte) error {
	if len(payload) != 20 {
		return fmt.Errorf("invalid len from DAppAddressRelay %v", len(payload))
	}
	address := (common.Address)(payload)
	env.dappAddress = &address
	env.Logf("got dapp address: %v", address)
	return nil
}
