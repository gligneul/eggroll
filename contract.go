// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level framework for Cartesi rollups in Go.
package eggroll

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/eggeth"
	"github.com/gligneul/eggroll/internal/rollups"
	"github.com/gligneul/eggroll/wallets"
	"github.com/holiman/uint256"
)

// Read from the rollups environment.
// This interface will be passed to the contract inspect method.
type EnvReader interface {

	// Get the raw input bytes.
	RawInput() []byte

	// Decode the input using the codecs.
	// If fails, return the error in the place of the value.
	DecodeInput() any

	// Get the DApp address.
	// The address is initialized after the contract receives an input from
	// the AddressRelay contract.
	DAppAddress() *common.Address

	// Call fmt.Sprintf, print the log, and store the result in the rollups state.
	// It is possible to retrieve this log in the DApp client.
	Logf(format string, a ...any)

	// Call fmt.Sprint, print the log, and store the result in the rollups state.
	// It is possible to retrieve this log in the DApp client.
	Log(a ...any)

	// Call fmt.Sprintf, print the log, try to store the result in the rollups state, and exit.
	// It is possible to retrieve this log in the DApp client.
	Fatalf(format string, a ...any)

	// Call fmt.Sprint, print the log, try to store the result in the rollups state, and exit.
	// It is possible to retrieve this log in the DApp client.
	Fatal(a ...any)

	// Return the list of addresses that have assets.
	EtherAddresses() []common.Address

	// Return the balance of the given address.
	EtherBalanceOf(address common.Address) uint256.Int
}

// Read and write the rollups environment.
// This interface will be passed to the contract advance method.
type Env interface {
	EnvReader

	// Get the Metadata for the current input.
	Metadata() *rollups.Metadata

	// Get the deposit for the current input if it came from a portal.
	Deposit() wallets.Deposit

	// Get the original sender for the current input.
	// If the input sender was a portal, this function returns the address that called the portal.
	Sender() common.Address

	// Transfer the given amount of funds from source to destination.
	// Return error if the source doesn't have enough funds.
	EtherTransfer(src common.Address, dst common.Address, value *uint256.Int) error

	// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
	// Return the voucher index.
	// Return error if the address doesn't have enough assets.
	EtherWithdraw(address common.Address, value *uint256.Int) (int, error)

	// Send a voucher. Return the voucher's index.
	Voucher(destination common.Address, payload []byte) int

	// Send a notice. Return the notice's index.
	Notice(payload []byte) int
}

// The Contract is the on-chain part of a rollups DApp.
// EggRoll uses the contract's codecs to encode the input and return values.
// For the advance and inspect methods, if the return value is []byte, return
// the raw bytes.
// If the call returns an error, EggRoll rejects the input.
type Contract interface {

	// Advance the contract state.
	Advance(env Env) (any, error)

	// Inspect the contract state.
	Inspect(env EnvReader) (any, error)

	// Get the codecs required by the contract.
	Codecs() []Codec
}

// DefaultContract provides a default implementation for optional contract methods.
type DefaultContract struct{}

// Reject inspect request.
func (_ DefaultContract) Inspect(env EnvReader) (any, error) {
	return nil, fmt.Errorf("inspect not supported")
}

// Return empty list of codecs.
func (_ *DefaultContract) Codecs() []Codec {
	return nil
}

// Start the Cartesi rollups for the contract.
// This function doesn't return and exits if there is an error.
func Roll(contract Contract) {
	rollupsAPI := rollups.NewRollupsHTTP()
	codecManager := newCodecManager(contract.Codecs())
	env := newEnv(rollupsAPI, codecManager)
	walletMap := map[common.Address]wallets.Wallet{
		eggeth.AddressEtherPortal: env.etherWallet,
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

		var returnPayload []byte
		if return_ != nil {
			var ok bool
			returnPayload, ok = return_.([]byte)
			if !ok {
				returnPayload, err = codecManager.encode(return_)
			}
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

func handleAdvance(
	env *env,
	contract Contract,
	walletMap map[common.Address]wallets.Wallet,
	input *rollups.AdvanceInput,
) (any, error) {
	var deposit wallets.Deposit
	var rawInput []byte

	if input.Metadata.Sender == eggeth.AddressDAppAddressRelay {
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

func handleDAppAddressRelay(env *env, payload []byte) error {
	if len(payload) != 20 {
		return fmt.Errorf("invalid len from DAppAddressRelay %v", len(payload))
	}
	address := (common.Address)(payload)
	env.setDAppAddress(&address)
	env.Logf("got dapp address: %v", address)
	return nil
}

func handleInspect(
	env *env,
	contract Contract,
	input *rollups.InspectInput,
) (any, error) {
	env.setInputData(nil, nil, input.Payload)
	return contract.Inspect(env)
}
