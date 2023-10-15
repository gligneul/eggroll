// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func testKey(t *testing.T, mnemonic string, account uint32, expected *ecdsa.PrivateKey) {
	key, err := mnemonicToPrivateKey(mnemonic, account)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !key.Equal(expected) {
		t.Fatalf("wrong key; expected %v; got %v", expected, key)
	}
}

func TestMnemonicToPrivateKey(t *testing.T) {
	mnemonic := "test test test test test test test test test test test junk"

	expected, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	testKey(t, mnemonic, 0, expected)

	expected, _ = crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	testKey(t, mnemonic, 1, expected)
}
