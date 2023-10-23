// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// Create a signer by using a mnemonic.
type MnemonicSigner struct {
	privateKey *ecdsa.PrivateKey
	chainId    *big.Int
	mnemonic   string
}

// Create a new mnemonic signer.
func NewMnemonicSigner(mnemonic string, accountIndex uint32, chainId *big.Int) (
	s *MnemonicSigner, err error) {

	signer := &MnemonicSigner{
		privateKey: nil,
		chainId:    chainId,
		mnemonic:   mnemonic,
	}
	if err = signer.SetAccount(accountIndex); err != nil {
		return nil, err
	}
	return signer, nil
}

func (s *MnemonicSigner) SetAccount(accountIndex uint32) error {
	privateKey, err := mnemonicToPrivateKey(s.mnemonic, accountIndex)
	if err != nil {
		return err
	}
	s.privateKey = privateKey
	return nil
}

func (s *MnemonicSigner) MakeTransactor() (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainId)
}

func (s *MnemonicSigner) Account() common.Address {
	publicKey := s.privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// Create the private key from mnemonic and account index based on the BIP44 standard.
// For more info on BIP44, see https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
func mnemonicToPrivateKey(mnemonic string, accountIndex uint32) (*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key: %v", err)
	}

	// get key at path m/44'/60'/0'/0/account
	const hardenedKeyStart uint32 = 0x80000000
	levels := []uint32{
		hardenedKeyStart + 44,
		hardenedKeyStart + 60,
		hardenedKeyStart + 0,
		0,
		accountIndex,
	}
	key := masterKey
	for i, level := range levels {
		key, err = key.NewChildKey(level)
		if err != nil {
			return nil, fmt.Errorf("failed to get child %v: %v", i, err)
		}
	}

	return crypto.ToECDSA(key.Key)
}
