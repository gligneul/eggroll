// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level framework for Cartesi rollups in Go.
package eggroll

import (
	"fmt"
	"math/big"

	"github.com/gligneul/eggroll/pkg/eggeth"
	"github.com/gligneul/eggroll/pkg/eggwallets"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/internal/rollups"
)

// Read from the rollups environment.
// This interface will be passed to the contract inspect method.
type EnvReader interface {

	// Get the DApp address.
	// The address is initialized after the contract receives an input from
	// the AddressRelay contract.
	DAppAddress() *common.Address

	// Send a report to the Rollups API.
	// Reports can be any array of bytes, and are save even when the
	// contract revers the input.
	Report(payload []byte)

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
	EtherBalanceOf(address common.Address) *big.Int

	// Return the list of tokens with assets.
	ERC20Tokens() []common.Address

	// Return the list of addresses that have assets for the given token.
	ERC20Addresses(token common.Address) []common.Address

	// Return the balance of the given address for the given token.
	ERC20BalanceOf(token common.Address, address common.Address) *big.Int
}

// Read and write the rollups environment.
// This interface will be passed to the contract advance method.
type Env interface {
	EnvReader

	// Get the Metadata for the current input.
	Metadata() *rollups.Metadata

	// Get the deposit for the current input if it came from a portal.
	Deposit() eggwallets.Deposit

	// Get the original sender for the current input.
	// If the input sender was a portal, this function returns the address that called the portal.
	Sender() common.Address

	// Send a voucher. Return the voucher's index.
	Voucher(destination common.Address, payload []byte) int

	// Send a notice. Return the notice's index.
	Notice(payload []byte) int

	// Transfer the given amount of funds from source to destination.
	// Return error if the source doesn't have enough funds.
	EtherTransfer(src common.Address, dst common.Address, value *big.Int) error

	// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
	// Return the voucher index.
	// Return error if the address doesn't have enough assets.
	EtherWithdraw(address common.Address, value *big.Int) (int, error)

	// Transfer the given amount of tokens from source to destination.
	// Return error if the source doesn't have enough funds.
	ERC20Transfer(token common.Address, src common.Address, dst common.Address, value *big.Int) error

	// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
	// Return error if the address doesn't have enough assets.
	ERC20Withdraw(token common.Address, address common.Address, value *big.Int) (int, error)
}

// The MiddlewareContract is the on-chain part of a rollups DApp.
// EggRoll uses the contract's codecs to encode the input and return values.
// For the advance and inspect methods, if the return value is []byte, return
// the raw bytes.
// If the call returns an error, EggRoll rejects the input.
type MiddlewareContract interface {

	// Advance the contract state.
	Advance(env Env, input []byte) error

	// Inspect the contract state.
	Inspect(env EnvReader, input []byte) error
}

// Start the Cartesi rollups for the contract.
// This function doesn't return and exits if there is an error.
func Roll(contract MiddlewareContract) {
	rollupsAPI := rollups.NewRollupsHTTP()
	env := newEnv(rollupsAPI)
	status := rollups.FinishStatusAccept
	for {
		input, err := rollupsAPI.Finish(status)
		if err != nil {
			env.Fatalf("failed to send finish: %v\n", err)
		}

		switch input := input.(type) {
		case *rollups.AdvanceInput:
			err = handleAdvance(env, contract, input)
		case *rollups.InspectInput:
			err = handleInspect(env, contract, input)
		default:
			// impossible
			panic("invalid input type")
		}

		if err != nil {
			env.Logf("rejecting: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}

		status = rollups.FinishStatusAccept
	}
}

func handleAdvance(
	env *env,
	contract MiddlewareContract,
	input *rollups.AdvanceInput,
) error {
	var deposit eggwallets.Deposit
	var rawInput []byte

	if input.Metadata.Sender == eggeth.AddressDAppAddressRelay {
		return handleDAppAddressRelay(env, input.Payload)
	}

	wallet, ok := env.walletMap[input.Metadata.Sender]
	if ok {
		var err error
		deposit, rawInput, err = wallet.Deposit(input.Payload)
		if err != nil {
			return fmt.Errorf("malformed portal input: %v", err)
		}
	} else {
		deposit = nil
		rawInput = input.Payload
	}

	env.setInputData(input.Metadata, deposit)
	return contract.Advance(env, rawInput)
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
	contract MiddlewareContract,
	input *rollups.InspectInput,
) error {
	env.setInputData(nil, nil)
	return contract.Inspect(env, input.Payload)
}
