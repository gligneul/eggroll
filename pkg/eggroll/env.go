// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"
	"os"

	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/gligneul/eggroll/pkg/eggwallets"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/internal/rollups"
	"github.com/holiman/uint256"
)

// Implementation of the Env and EnvReader interfaces.
type env struct {
	logger      *log.Logger
	rollups     *rollups.RollupsHTTP
	etherWallet *eggwallets.EtherWallet
	dappAddress *common.Address

	// The fields below should be set for each input.
	metadata *rollups.Metadata
	deposit  eggwallets.Deposit
}

//
// Internal methods
//

func newEnv(rollups *rollups.RollupsHTTP) *env {
	return &env{
		logger:      log.New(os.Stdout, "", 0),
		rollups:     rollups,
		etherWallet: eggwallets.NewEtherWallet(),
	}
}

func (e *env) setInputData(metadata *rollups.Metadata, deposit eggwallets.Deposit) {
	e.metadata = metadata
	e.deposit = deposit
}

func (e *env) setDAppAddress(address *common.Address) {
	e.dappAddress = address
}

// Log the message and send a report.
func (e *env) log(message string) {
	e.logger.Print(message)
	log := eggtypes.Log{Message: message}
	e.Report(log.Pack())
}

// Log the message, send a report, and exit.
func (e *env) fatal(message string) {
	e.log(message)
	os.Exit(1)
}

//
// Implementation of EnvReader
//

func (e *env) DAppAddress() *common.Address {
	return e.dappAddress
}

func (e *env) Report(payload []byte) {
	if err := e.rollups.SendReport(payload); err != nil {
		e.logger.Fatalf("failed to send report: %v\n", err)
	}
}

func (e *env) Logf(format string, a ...any) {
	e.log(fmt.Sprintf(format, a...))
}

func (e *env) Log(a ...any) {
	e.log(fmt.Sprint(a...))
}

func (e *env) Fatalf(format string, a ...any) {
	e.fatal(fmt.Sprintf(format, a...))
}

func (e *env) Fatal(a ...any) {
	e.fatal(fmt.Sprint(a...))
}

func (e *env) EtherAddresses() []common.Address {
	return e.etherWallet.Addresses()
}

func (e *env) EtherBalanceOf(address common.Address) uint256.Int {
	return e.etherWallet.BalanceOf(address)
}

//
// Implementation of Env
//

func (e *env) Metadata() *rollups.Metadata {
	return e.metadata
}

func (e *env) Deposit() eggwallets.Deposit {
	return e.deposit
}

func (e *env) Sender() common.Address {
	if e.deposit != nil {
		return e.deposit.GetSender()
	}
	return e.metadata.Sender
}

func (e *env) EtherTransfer(src common.Address, dst common.Address, value *uint256.Int) error {
	return e.etherWallet.Transfer(src, dst, value)
}

func (e *env) EtherWithdraw(address common.Address, value *uint256.Int) (int, error) {
	if e.dappAddress == nil {
		return 0, fmt.Errorf("need dapp address to withdraw")
	}
	voucher, err := e.etherWallet.Withdraw(address, value)
	if err != nil {
		return 0, err
	}
	return e.Voucher(*e.dappAddress, voucher), nil
}

func (e *env) Voucher(destination common.Address, payload []byte) int {
	index, err := e.rollups.SendVoucher(destination, payload)
	if err != nil {
		e.Fatalf("failed to send voucher: %v", err)
	}
	return index
}

func (e *env) Notice(payload []byte) int {
	index, err := e.rollups.SendNotice(payload)
	if err != nil {
		e.Fatalf("failed to send notice: %v", err)
	}
	return index
}
