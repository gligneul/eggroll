// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/rollups"
	"github.com/gligneul/eggroll/wallets"
	"github.com/holiman/uint256"
)

// Env allows the DApp contract to interact with the Rollups API.
type Env struct {
	rollups     rollupsAPI
	wallets     *wallets.Wallets
	dappAddress common.Address
	metadata    *rollups.Metadata
	deposit     wallets.Deposit
}

// Get the Metadata for the current input.
func (e *Env) Metadata() *rollups.Metadata {
	return e.metadata
}

// Get the original sender for the current input.
// If the input sender was a portal, this function returns the address that called the portal.
func (e *Env) Sender() common.Address {
	if e.deposit != nil {
		return e.deposit.GetSender()
	}
	return e.metadata.Sender
}

// Get the deposit for the current input if it came from a portal.
func (e *Env) Deposit() wallets.Deposit {
	return e.deposit
}

// Return the list of addresses that have assets.
func (e *Env) EtherAddresses() []common.Address {
	return e.wallets.Ether.Addresses()
}

// Return the balance of the given address.
func (e *Env) EtherBalanceOf(address common.Address) uint256.Int {
	return e.wallets.Ether.BalanceOf(address)
}

// Transfer the given amount of funds from source to destination.
// Return error if the source doesn't have enough funds.
func (e *Env) EtherTransfer(src common.Address, dst common.Address, value *uint256.Int) error {
	return e.wallets.Ether.Transfer(src, dst, value)
}

// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
// Return error if the address doesn't have enough assets.
func (e *Env) EtherWithdraw(address common.Address, value *uint256.Int) error {
	voucher, err := e.wallets.Ether.Withdraw(address, value)
	if err != nil {
		return err
	}
	e.Voucher(e.dappAddress, voucher)
	return nil
}

// Call fmt.Sprintln, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *Env) Logln(a any) {
	e.Log(fmt.Sprintln(a))
}

// Call fmt.Sprintf, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *Env) Logf(format string, a ...any) {
	e.Log(fmt.Sprintf(format, a...))
}

// Call fmt.Sprint, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *Env) Log(a any) {
	e.log(fmt.Sprint(a))
}

// Log the message and send a report.
func (e *Env) log(message string) {
	log.Print(message)
	if err := e.rollups.SendReport([]byte(message)); err != nil {
		log.Fatalf("failed to send report: %v\n", err)
	}
}

// Send a voucher.
func (e *Env) Voucher(destination common.Address, payload []byte) {
	if err := e.rollups.SendVoucher(destination, payload); err != nil {
		log.Fatalf("failed to send voucher: %v\n", err)
	}
}
