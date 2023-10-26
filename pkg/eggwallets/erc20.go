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

// An ERC20 deposit that arrived in the wallet.
type ERC20Deposit struct {
	Token  common.Address
	Sender common.Address
	Amount *big.Int
}

func (d *ERC20Deposit) GetSender() common.Address {
	return d.Sender
}

func (d *ERC20Deposit) String() string {
	token := strings.ToLower(d.Token.String())
	sender := strings.ToLower(d.Sender.String())
	value := d.Amount.String()
	return fmt.Sprintf("%v deposited %v of %v token", sender, value, token)
}

// Wallet that manages ERC20 tokens.
type ERC20Wallet struct {
	balance map[common.Address]map[common.Address]big.Int
}

// Create new ERC20 Wallet.
func NewERC20Wallet() *ERC20Wallet {
	return &ERC20Wallet{
		balance: make(map[common.Address]map[common.Address]big.Int),
	}
}

// Return the list of tokens with assets.
func (w *ERC20Wallet) Tokens() []common.Address {
	var tokens []common.Address
	for t := range w.balance {
		tokens = append(tokens, t)
	}
	sortAddresses(tokens)
	return tokens
}

// Return the list of addresses that have assets for the given token.
func (w *ERC20Wallet) Addresses(token common.Address) []common.Address {
	var addresses []common.Address
	for a := range w.balance[token] {
		addresses = append(addresses, a)
	}
	sortAddresses(addresses)
	return addresses
}

func (w *ERC20Wallet) setBalance(token common.Address, address common.Address, value *big.Int) {
	if value.Sign() == 0 {
		if w.balance[token] != nil {
			delete(w.balance[token], address)
			if len(w.balance[token]) == 0 {
				delete(w.balance, token)
			}
		}
	} else {
		if w.balance[token] == nil {
			w.balance[token] = make(map[common.Address]big.Int)
		}
		w.balance[token][address] = *value
	}
}

// Return the balance of the given address for the given token.
func (w *ERC20Wallet) BalanceOf(token common.Address, address common.Address) *big.Int {
	balance := w.balance[token][address]
	return &balance
}

// Transfer the given amount of tokens from source to destination.
// Return error if the source doesn't have enough funds.
func (w *ERC20Wallet) Transfer(
	token common.Address, src common.Address, dst common.Address, value *big.Int) error {

	if src == dst {
		return fmt.Errorf("can't transfer to self")
	}

	newSrcBalance := new(big.Int).Sub(w.BalanceOf(token, src), value)
	if newSrcBalance.Sign() < 0 {
		return fmt.Errorf("insuficient funds")
	}

	newDstBalance := new(big.Int).Add(w.BalanceOf(token, dst), value)
	if newDstBalance.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("balance overflow")
	}

	// commit
	w.setBalance(token, src, newSrcBalance)
	w.setBalance(token, dst, newDstBalance)
	return nil
}

// Encode the withdraw request to the portal.
func EncodeERC20Withdraw(token common.Address, value *big.Int) []byte {
	abiJson := `[{
		"type": "function",
		"name": "transfer",
		"inputs": [
			{"type": "address"},
			{"type": "uint256"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		log.Panicf("failed to decode ABI: %v", err)
	}
	voucher, err := abiInterface.Pack("transfer", token, value)
	if err != nil {
		log.Panicf("failed to pack: %v", err)
	}
	return voucher
}

// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
// Return error if the address doesn't have enough assets.
func (w *ERC20Wallet) Withdraw(
	token common.Address, address common.Address, value *big.Int) ([]byte, error) {

	newBalance := new(big.Int).Sub(w.BalanceOf(token, address), value)
	if newBalance.Sign() < 0 {
		return nil, fmt.Errorf("insuficient funds")
	}
	w.setBalance(token, address, newBalance)
	return EncodeERC20Withdraw(address, value), nil
}

// Handle a deposit from the ERC20 portal.
func (w *ERC20Wallet) Deposit(payload []byte) (Deposit, []byte, error) {
	if len(payload) < 1+20+20+32 {
		return nil, nil, fmt.Errorf("invalid erc20 deposit size; got %v", len(payload))
	}

	if payload[0] == 0 {
		return nil, nil, fmt.Errorf("received failed erc20 transfer")
	}
	payload = payload[1:]

	token := common.BytesToAddress(payload[:20])
	payload = payload[20:]

	sender := common.BytesToAddress(payload[:20])
	payload = payload[20:]

	amount := new(big.Int).SetBytes(payload[:32])
	payload = payload[32:]

	newBalance := new(big.Int).Add(w.BalanceOf(token, sender), amount)
	if newBalance.Cmp(MaxUint256) > 0 {
		// This should not be possible in real world, but we handle it anyway.
		newBalance = MaxUint256
	}
	w.setBalance(token, sender, newBalance)

	deposit := &ERC20Deposit{token, sender, amount}
	return deposit, payload, nil
}
