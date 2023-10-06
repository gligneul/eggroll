// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package wallets

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// An Ether deposit that arrived in the wallet.
type EtherDeposit struct {
	Sender common.Address
	Value  uint256.Int
}

func (e *EtherDeposit) GetSender() common.Address {
	return e.Sender
}

func (e *EtherDeposit) String() string {
	sender := strings.ToLower(e.Sender.String())
	value := e.Value.ToBig().String()
	return fmt.Sprintf("%v deposited %v wei", sender, value)
}

// Wallet that manages Ether.
type EtherWallet struct {
	balance map[common.Address]uint256.Int
}

// Create new Ether Wallet.
func NewEtherWallet() *EtherWallet {
	return &EtherWallet{
		balance: make(map[common.Address]uint256.Int),
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

func (w *EtherWallet) setBalance(address common.Address, value *uint256.Int) {
	if value.IsZero() {
		delete(w.balance, address)
	} else {
		w.balance[address] = *value
	}
}

// Return the balance of the given address.
func (w *EtherWallet) BalanceOf(address common.Address) uint256.Int {
	return w.balance[address]
}

// Transfer the given amount of funds from source to destination.
// Return error if the source doesn't have enough funds.
func (w *EtherWallet) Transfer(src common.Address, dst common.Address, value *uint256.Int) error {
	if src == dst {
		return fmt.Errorf("can't transfer to self")
	}

	srcBalance := w.balance[src]
	newSrcBalance, underflow := new(uint256.Int).SubOverflow(&srcBalance, value)
	if underflow {
		return fmt.Errorf("insuficient funds")
	}

	dstBalance := w.balance[dst]
	newDstBalance, overflow := new(uint256.Int).AddOverflow(&dstBalance, value)
	if overflow {
		return fmt.Errorf("balance overflow")
	}

	// commit
	w.setBalance(src, newSrcBalance)
	w.setBalance(dst, newDstBalance)
	return nil
}

// Encode the withdraw request to the portal.
func EncodeEtherWithdraw(address common.Address, value *uint256.Int) []byte {
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
	voucher, err := abiInterface.Pack("withdrawEther", address, value.ToBig())
	if err != nil {
		log.Panicf("failed to pack: %v", err)
	}
	return voucher
}

// Withdraw the asset from the wallet and generate the voucher to withdraw from the portal.
// Return error if the address doesn't have enough assets.
func (w *EtherWallet) Withdraw(address common.Address, value *uint256.Int) ([]byte, error) {
	balance := w.balance[address]
	newBalance, underflow := new(uint256.Int).SubOverflow(&balance, value)
	if underflow {
		return nil, fmt.Errorf("insuficient funds")
	}
	w.setBalance(address, newBalance)
	return EncodeEtherWithdraw(address, value), nil
}

func (w *EtherWallet) deposit(payload []byte) (Deposit, []byte, error) {
	if len(payload) < 20+32 {
		return nil, nil, fmt.Errorf("invalid eth deposit size; got %v", len(payload))
	}

	deposit := &EtherDeposit{
		Sender: common.BytesToAddress(payload[0:20]),
		Value:  *new(uint256.Int).SetBytes(payload[20 : 20+32]),
	}
	input := payload[20+32:]

	oldBalance := w.balance[deposit.Sender]
	newBalance, overflow := new(uint256.Int).AddOverflow(&oldBalance, &deposit.Value)
	if overflow {
		newBalance = IntMax
	}
	w.setBalance(deposit.Sender, newBalance)

	return deposit, input, nil
}

func (_ *EtherWallet) portalAddress() common.Address {
	return common.HexToAddress("0xffdbe43d4c855bf7e0f105c400a50857f53ab044")
}
