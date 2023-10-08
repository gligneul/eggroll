// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/rollups"
	"github.com/gligneul/eggroll/wallets"
	"github.com/holiman/uint256"
)

// Interface with the Rollups API.
// We don't expose this API because calling it directly will break EggRoll assumptions.
type envRollupsAPI interface {
	SendVoucher(destination common.Address, payload []byte) (int, error)
	SendNotice(payload []byte) (int, error)
}

// Env allows the DApp contract to interact with the Rollups API.
type Env struct {
	*Reporter

	rollups     envRollupsAPI
	etherWallet *wallets.EtherWallet
	dappAddress *common.Address

	// Reset for each input.
	metadata *rollups.Metadata

	// Reset for each input.
	deposit wallets.Deposit
}

// Create a new env.
func newEnv(reporter *Reporter, rollups envRollupsAPI) *Env {
	return &Env{
		Reporter:    reporter,
		rollups:     rollups,
		etherWallet: wallets.NewEtherWallet(),
	}
}

func (e *Env) setInputData(metadata *rollups.Metadata, deposit wallets.Deposit) {
	e.metadata = metadata
	e.deposit = deposit
}

// Get the Metadata for the current input.
func (e *Env) Metadata() *rollups.Metadata {
	return e.metadata
}

// Get the deposit for the current input if it came from a portal.
func (e *Env) Deposit() wallets.Deposit {
	return e.deposit
}

// Get the original sender for the current input.
// If the input sender was a portal, this function returns the address that called the portal.
func (e *Env) Sender() common.Address {
	if e.deposit != nil {
		return e.deposit.GetSender()
	}
	return e.metadata.Sender
}

func (e *Env) setDAppAddress(address *common.Address) {
	e.dappAddress = address
}

// Get the DApp address.
// The address is initialized after the contract receives an input from the AddressRelay contract.
func (e *Env) DAppAddress() *common.Address {
	return e.dappAddress
}

// Return the list of addresses that have assets.
func (e *Env) EtherAddresses() []common.Address {
	return e.etherWallet.Addresses()
}

// Return the balance of the given address.
func (e *Env) EtherBalanceOf(address common.Address) uint256.Int {
	return e.etherWallet.BalanceOf(address)
}

// Transfer the given amount of funds from source to destination.
// Return error if the source doesn't have enough funds.
func (e *Env) EtherTransfer(src common.Address, dst common.Address, value *uint256.Int) error {
	return e.etherWallet.Transfer(src, dst, value)
}

// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
// Return the voucher index.
// Return error if the address doesn't have enough assets.
func (e *Env) EtherWithdraw(address common.Address, value *uint256.Int) (int, error) {
	if e.dappAddress == nil {
		return 0, fmt.Errorf("need dapp address to withdraw")
	}
	voucher, err := e.etherWallet.Withdraw(address, value)
	if err != nil {
		return 0, err
	}
	return e.Voucher(*e.dappAddress, voucher), nil
}

// Send a voucher. Return the voucher's index.
func (e *Env) Voucher(destination common.Address, payload []byte) int {
	index, err := e.rollups.SendVoucher(destination, payload)
	if err != nil {
		e.Fatalf("failed to send voucher: %v", err)
	}
	return index
}

// Send a notice. Return the notice's index.
func (e *Env) Notice(payload []byte) int {
	index, err := e.rollups.SendNotice(payload)
	if err != nil {
		e.Fatalf("failed to send notice: %v", err)
	}
	return index
}
