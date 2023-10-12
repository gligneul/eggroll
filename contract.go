// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level framework for Cartesi Rollups in Go.
package eggroll

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/blockchain"
	"github.com/gligneul/eggroll/rollups"
	"github.com/gligneul/eggroll/wallets"
)

// Contract that advances the dapp state.
type Contract interface {

	// Get the codecs required by the contract.
	// Codecs are used to decode inputs and encode returns.
	Codecs() []Codec

	// Advance the state with the given input, returning the advance result.
	// EggRoll uses the contract's codecs to encode the return value.
	// If the return value is []byte, return the raw bytes.
	// If the call returns an error, EggRoll rejects the input.
	Advance(env *Env) (any, error)

	// Inspect the state with the given input, returning the inspect result.
	// EggRoll uses the contract's codecs to encode the return value.
	// If the return value is []byte, return the raw bytes.
	// If the call returns an error, EggRoll rejects the input.
	Inspect(env *Env) (any, error)
}

// DefaultContract provides a default implementation for optional contract methods.
type DefaultContract struct{}

// Return empty list of codecs.
func (_ *DefaultContract) Codecs() []Codec {
	return nil
}

// Reject inspect request.
func (_ DefaultContract) Inspect(env *Env) (any, error) {
	return nil, fmt.Errorf("inspect not supported")
}

// Start the Cartesi Rollups for the contract.
// This function doesn't return and exits if there is an error.
func Roll(contract Contract) {
	rollupsAPI := rollups.NewRollupsHTTP()
	reporter := newReporter(rollupsAPI)
	codecManager := newCodecManager(contract.Codecs())
	env := newEnv(reporter, rollupsAPI, codecManager)
	walletMap := map[common.Address]wallets.Wallet{
		blockchain.AddressEtherPortal: env.etherWallet,
	}

	status := rollups.FinishStatusAccept

	for {
		input, err := rollupsAPI.Finish(status)
		if err != nil {
			env.Fatalf("failed to send finish: %v\n", err)
		}

		var return_ any
		switch input := input.(type) {
		case *rollups.AdvanceInput:
			return_, err = handleAdvance(env, contract, walletMap, input)
		case *rollups.InspectInput:
			return_, err = handleInspect(env, contract, input)
		default:
			err = fmt.Errorf("invalid input type")
		}

		returnPayload, ok := return_.([]byte)
		if !ok {
			returnPayload, err = codecManager.encode(return_)
		}

		if err != nil {
			env.Logf("rejecting: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}

		env.sendReturn(returnPayload)
		status = rollups.FinishStatusAccept
	}
}

// Handle the advance input from the rollups API.
func handleAdvance(
	env *Env,
	contract Contract,
	walletMap map[common.Address]wallets.Wallet,
	input *rollups.AdvanceInput,
) (any, error) {
	var deposit wallets.Deposit
	var rawInput []byte

	if input.Metadata.Sender == blockchain.AddressDAppAddressRelay {
		return nil, handleDAppAddressRelay(env, input.Payload)
	}

	wallet, ok := walletMap[input.Metadata.Sender]
	if ok {
		var err error
		deposit, rawInput, err = wallet.Deposit(input.Payload)
		if err != nil {
			return nil, fmt.Errorf("malformed portal input: %v", err)
		}
	} else {
		deposit = nil
		rawInput = input.Payload
	}

	env.setInputData(input.Metadata, deposit, rawInput)
	return contract.Advance(env)
}

// Handle the input from the DAppAddressRelay.
func handleDAppAddressRelay(env *Env, payload []byte) error {
	if len(payload) != 20 {
		return fmt.Errorf("invalid len from DAppAddressRelay %v", len(payload))
	}
	address := (common.Address)(payload)
	env.setDAppAddress(&address)
	env.Logf("got dapp address: %v", address)
	return nil
}

// Handle the inspect input from the rollups API.
func handleInspect(
	env *Env,
	contract Contract,
	input *rollups.InspectInput,
) (any, error) {
	env.setInputData(nil, nil, input.Payload)
	return contract.Inspect(env)
}
