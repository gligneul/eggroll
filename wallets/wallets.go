// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Wallets that manage assets in Cartesi Rollups Portals.
package wallets

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// Max value for uint256.
var IntMax *uint256.Int

func init() {
	IntMax = uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
}

// Asset deposit to a portal.
type Deposit interface {
	fmt.Stringer

	// Get the deposit sender.
	GetSender() common.Address
}

// Manages an asset in a portal.
type Wallet interface {

	// Handle the raw bytes of a deposit that came from the portal.
	// After handling the deposit, return the parsed deposit, and the DApp input paylaod.
	deposit(input []byte) (Deposit, []byte, error)

	// Get the portal address.
	portalAddress() common.Address
}

// Map portals to wallets.
type WalletDispatcher struct {
	handlers map[common.Address]Wallet
}

// If there is a wallet for the given portal, dispatch the input to it, and return the parsed
// deposit, and the DApp input payload.
// If there isn't a corresponding portal, return nil.
func (d *WalletDispatcher) Dispatch(portal common.Address, input []byte) (Deposit, []byte, error) {
	wallet, ok := d.handlers[portal]
	if ok {
		return wallet.deposit(input)
	}
	return nil, nil, nil
}

// The Wallets struct store all available wallets.
// Each type of asset has its wallet.
type Wallets struct {
	Ether *EtherWallet
}

// Create all available wallets.
func NewWallets() *Wallets {
	return &Wallets{
		Ether: NewEtherWallet(),
	}
}

// Create the wallet dispatcher from the wallets.
func (w *Wallets) MakeDispatcher() *WalletDispatcher {
	handlers := map[common.Address]Wallet{
		w.Ether.portalAddress(): w.Ether,
	}
	return &WalletDispatcher{handlers}
}
