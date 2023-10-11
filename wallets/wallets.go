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
	// After handling the deposit, return the parsed deposit, and the DApp input payload.
	Deposit(input []byte) (Deposit, []byte, error)
}
