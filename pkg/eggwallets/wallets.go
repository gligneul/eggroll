// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Wallets that manage assets in Cartesi Rollups Portals.
package eggwallets

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Max value for uint256.
var MaxUint256 *big.Int

func init() {
	one := big.NewInt(1)
	// Left shift by 256 bits and then subtract 1 to get the max value of uint256.
	MaxUint256 = new(big.Int).Sub(new(big.Int).Lsh(one, 256), one)
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
