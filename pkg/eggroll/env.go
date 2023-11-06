// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/gligneul/eggroll/pkg/eggeth"
	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/gligneul/eggroll/pkg/eggwallets"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/internal/rollups"
)

// Implementation of the Env and EnvReader interfaces.
type env struct {
	logger      *log.Logger
	rollups     *rollups.RollupsHTTP
	etherWallet *eggwallets.EtherWallet
	erc20Wallet *eggwallets.ERC20Wallet
	dappAddress *common.Address
	walletMap   map[common.Address]eggwallets.Wallet

	// The fields below should be set for each input.
	metadata *rollups.Metadata
	deposit  eggwallets.Deposit
}

//
// Internal methods
//

func newEnv(rollups *rollups.RollupsHTTP) *env {
	etherWallet := eggwallets.NewEtherWallet()
	erc20Wallet := eggwallets.NewERC20Wallet()
	walletMap := map[common.Address]eggwallets.Wallet{
		eggeth.AddressEtherPortal: etherWallet,
		eggeth.AddressERC20Portal: erc20Wallet,
	}
	return &env{
		logger:      log.New(os.Stdout, "", 0),
		rollups:     rollups,
		etherWallet: etherWallet,
		erc20Wallet: erc20Wallet,
		walletMap:   walletMap,
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
	e.Report(eggtypes.EncodeLog(message))
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

func (e *env) EtherBalanceOf(address common.Address) *big.Int {
	return e.etherWallet.BalanceOf(address)
}

func (e *env) ERC20Tokens() []common.Address {
	return e.erc20Wallet.Tokens()
}

func (e *env) ERC20Addresses(token common.Address) []common.Address {
	return e.erc20Wallet.Addresses(token)
}

func (e *env) ERC20BalanceOf(token common.Address, address common.Address) *big.Int {
	return e.erc20Wallet.BalanceOf(token, address)
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

func (e *env) EtherTransfer(src common.Address, dst common.Address, value *big.Int) error {
	return e.etherWallet.Transfer(src, dst, value)
}

func (e *env) EtherWithdraw(address common.Address, value *big.Int) (int, error) {
	if e.dappAddress == nil {
		return 0, fmt.Errorf("need dapp address to withdraw")
	}
	voucher, err := e.etherWallet.Withdraw(address, value)
	if err != nil {
		return 0, err
	}
	return e.Voucher(*e.dappAddress, voucher), nil
}

func (e *env) ERC20Transfer(token common.Address, src common.Address, dst common.Address, value *big.Int) error {
	return e.erc20Wallet.Transfer(token, src, dst, value)
}

func (e *env) ERC20Withdraw(token common.Address, address common.Address, value *big.Int) (int, error) {
	voucher, err := e.erc20Wallet.Withdraw(token, address, value)
	if err != nil {
		return 0, err
	}
	return e.Voucher(eggeth.AddressERC20Portal, voucher), nil
}
