// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggwallets

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Generate a string converting Wei to Ether.
func EtherString(wei *big.Int) string {
	weiFloat := new(big.Float).SetInt(wei)
	tenToEighteen := new(big.Float).SetFloat64(1e18)
	etherFloat := new(big.Float).Quo(weiFloat, tenToEighteen)
	return etherFloat.Text('f', 18)
}

// An Ether deposit that arrived in the wallet.
type EtherDeposit struct {
	Sender common.Address
	Value  *big.Int
}

func (e *EtherDeposit) GetSender() common.Address {
	return e.Sender
}

func (e *EtherDeposit) String() string {
	sender := strings.ToLower(e.Sender.String())
	value := EtherString(e.Value)
	return fmt.Sprintf("%v deposited %v Ether", sender, value)
}

// Wallet that manages Ether.
type EtherWallet struct {
	balance map[common.Address]big.Int
}

// Create new Ether Wallet.
func NewEtherWallet() *EtherWallet {
	return &EtherWallet{
		balance: make(map[common.Address]big.Int),
	}
}

// Return the list of addresses that have assets.
func (w *EtherWallet) Addresses() []common.Address {
	var addresses []common.Address
	for address := range w.balance {
		addresses = append(addresses, address)
	}
	return addresses
}

func (w *EtherWallet) setBalance(address common.Address, value *big.Int) {
	if value.Sign() == 0 {
		delete(w.balance, address)
	} else {
		w.balance[address] = *value
	}
}

// Return the balance of the given address.
func (w *EtherWallet) BalanceOf(address common.Address) *big.Int {
	balance := w.balance[address]
	return &balance
}

// Transfer the given amount of funds from source to destination.
// Return error if the source doesn't have enough funds.
func (w *EtherWallet) Transfer(src common.Address, dst common.Address, value *big.Int) error {
	if src == dst {
		return fmt.Errorf("can't transfer to self")
	}

	newSrcBalance := new(big.Int).Sub(w.BalanceOf(src), value)
	if newSrcBalance.Sign() < 0 {
		return fmt.Errorf("insuficient funds")
	}

	newDstBalance := new(big.Int).Add(w.BalanceOf(dst), value)
	if newDstBalance.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("balance overflow")
	}

	// commit
	w.setBalance(src, newSrcBalance)
	w.setBalance(dst, newDstBalance)
	return nil
}

// Encode the withdraw request to the portal.
func EncodeEtherWithdraw(address common.Address, value *big.Int) []byte {
	abiJson := `[{
		"type": "function",
		"name": "withdrawEther",
		"inputs": [
			{"type": "address"},
			{"type": "uint256"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		log.Panicf("failed to decode ABI: %v", err)
	}
	voucher, err := abiInterface.Pack("withdrawEther", address, value)
	if err != nil {
		log.Panicf("failed to pack: %v", err)
	}
	return voucher
}

// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
// Return error if the address doesn't have enough assets.
func (w *EtherWallet) Withdraw(address common.Address, value *big.Int) ([]byte, error) {
	newBalance := new(big.Int).Sub(w.BalanceOf(address), value)
	if newBalance.Sign() < 0 {
		return nil, fmt.Errorf("insuficient funds")
	}
	w.setBalance(address, newBalance)
	return EncodeEtherWithdraw(address, value), nil
}

// Handle a deposit from the Ether portal.
func (w *EtherWallet) Deposit(payload []byte) (Deposit, []byte, error) {
	if len(payload) < 20+32 {
		return nil, nil, fmt.Errorf("invalid eth deposit size; got %v", len(payload))
	}

	sender := common.BytesToAddress(payload[:20])
	payload = payload[20:]

	value := new(big.Int).SetBytes(payload[:32])
	payload = payload[32:]

	newBalance := new(big.Int).Add(w.BalanceOf(sender), value)
	if newBalance.Cmp(MaxUint256) > 0 {
		// This should not be possible in real world, but we handle it anyway.
		newBalance = MaxUint256
	}
	w.setBalance(sender, newBalance)

	deposit := &EtherDeposit{sender, value}
	return deposit, payload, nil
}
