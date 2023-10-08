// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level, opinionated, lambda-based framework for Cartesi Rollups in Go.
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

	// Get the decoders required by the contract.
	Decoders() []Decoder

	// Advance the state with the given input, returning the advance result.
	// EggRoll uses the contract decoders to decode the incoming input.
	// If the input doesn't match any decoder, the type of the input will be []byte.
	// The advance result is limited to 1000 Kb.
	// If the call returns an error, EggRoll rejects the input.
	Advance(env *Env, input any) ([]byte, error)

	// Inspect the state with the given input, returning the inspect result.
	// EggRoll uses the contract decoders to decode the incoming input.
	// If the input doesn't match any decoder, the type of the input will be []byte.
	// The inspect result is limited to 1000 Kb.
	// The inspect call mustn't change the state of the application.
	// If the call returns an error, EggRoll rejects the input.
	Inspect(env *Env, input any) ([]byte, error)
}

// DefaultContract provides a default implementation for optional contract methods.
type DefaultContract struct{}

func (_ *DefaultContract) Decoders() []Decoder {
	return nil
}

func (_ DefaultContract) Inspect(env *Env, input any) ([]byte, error) {
	return nil, fmt.Errorf("inspect not supported")
}

// Start the Cartesi Rollups for the contract.
// This function doesn't return and exits if there is an error.
func Roll(contract Contract) {
	rollupsAPI := rollups.NewRollupsHTTP()
	reporter := newReporter(rollupsAPI)
	env := newEnv(reporter, rollupsAPI)
	decoderMap := makeDecoderMap(contract.Decoders())
	walletMap := map[common.Address]wallets.Wallet{
		blockchain.AddressEtherPortal: env.etherWallet,
	}

	status := rollups.FinishStatusAccept

	for {
		input, err := rollupsAPI.Finish(status)
		if err != nil {
			env.Fatalf("failed to send finish: %v\n", err)
		}

		var result []byte
		switch input := input.(type) {
		case *rollups.AdvanceInput:
			result, err = handleAdvance(env, contract, decoderMap, walletMap, input)
		case *rollups.InspectInput:
			err = fmt.Errorf("inspect")
		default:
			err = fmt.Errorf("invalid input type")
		}

		if err != nil {
			env.Logf("rejecting: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}

		env.sendResult(result)
		status = rollups.FinishStatusAccept
	}
}

// Handle the advance input from the rollups API.
func handleAdvance(
	env *Env,
	contract Contract,
	decoderMap map[InputKey]Decoder,
	walletMap map[common.Address]wallets.Wallet,
	input *rollups.AdvanceInput,
) (
	result []byte,
	err error,
) {
	var deposit wallets.Deposit
	var inputBytes []byte

	if input.Metadata.Sender == blockchain.AddressDAppAddressRelay {
		return nil, handleDAppAddressRelay(env, input.Payload)
	}

	wallet, ok := walletMap[input.Metadata.Sender]
	if ok {
		deposit, inputBytes, err = wallet.Deposit(input.Payload)
		if err != nil {
			return nil, fmt.Errorf("malformed portal input: %v", err)
		}
	} else {
		deposit = nil
		inputBytes = input.Payload
	}

	decodedInput, err := decodeInput(decoderMap, inputBytes)
	if err != nil {
		return nil, err
	}

	env.setInputData(input.Metadata, deposit)
	return contract.Advance(env, decodedInput)
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
	decoderMap map[InputKey]Decoder,
	input *rollups.InspectInput,
) (
	result []byte,
	err error,
) {
	decodedInput, err := decodeInput(decoderMap, input.Payload)
	if err != nil {
		return nil, err
	}
	env.setInputData(nil, nil)
	return contract.Inspect(env, decodedInput)
}
