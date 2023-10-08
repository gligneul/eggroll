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
	SendVoucher(destination common.Address, payload []byte) error
	SendNotice(payload []byte) error
	SendReport(payload []byte) error
	Finish(status rollups.FinishStatus) ([]byte, *rollups.Metadata, error)
}

// Env allows the DApp contract to interact with the Rollups API.
type Env struct {
	*EnvLogger

	rollups     envRollupsAPI
	etherWallet *wallets.EtherWallet
	dappAddress *common.Address

	// Reset for each input.
	metadata *rollups.Metadata

	// Reset for each input.
	deposit wallets.Deposit
}

// Create a new env.
func newEnv(rollups envRollupsAPI) *Env {
	return &Env{
		EnvLogger:   newEnvLogger(rollups),
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
// Return error if the address doesn't have enough assets.
func (e *Env) EtherWithdraw(address common.Address, value *uint256.Int) error {
	if e.dappAddress == nil {
		return fmt.Errorf("need dapp address to withdraw")
	}
	voucher, err := e.etherWallet.Withdraw(address, value)
	if err != nil {
		return err
	}
	e.Voucher(*e.dappAddress, voucher)
	return nil
}

// Send a voucher.
func (e *Env) Voucher(destination common.Address, payload []byte) {
	if err := e.rollups.SendVoucher(destination, payload); err != nil {
		e.Fatalf("failed to send voucher: %v", err)
	}
}
