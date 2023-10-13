// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package blockchain

//go:generate abigen --abi abi/CartesiDApp.json --pkg blockchain --type CartesiDApp --out cartesidapp.go
//go:generate abigen --abi abi/DAppAddressRelay.json --pkg blockchain --type DAppAddressRelay --out dappaddressrelay.go
//go:generate abigen --abi abi/EtherPortal.json --pkg blockchain --type EtherPortal --out etherportal.go
//go:generate abigen --abi abi/InputBox.json --pkg blockchain --type InputBox --out inputbox.go

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Implements blockchain client for Ethereum using go-ethereum.
// This struct provides methods that are specific for the Cartesi Rollups.
type ETHClient struct {
	client           *ethclient.Client
	dappAddressRelay *DAppAddressRelay
	inputBox         *InputBox
}

// Create new ETH client.
func NewETHClient(endpoint string) (*ETHClient, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %v", err)
	}
	dappAddressRelay, err := NewDAppAddressRelay(AddressDAppAddressRelay, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAppAddressRelaya contract: %v", err)
	}
	inputBox, err := NewInputBox(AddressInputBox, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InputBox contract: %v", err)
	}
	ethClient := &ETHClient{
		client:           client,
		dappAddressRelay: dappAddressRelay,
		inputBox:         inputBox,
	}
	return ethClient, nil
}

// Send input to the blockchain.
func (c *ETHClient) SendInput(
	ctx context.Context, signer *bind.TransactOpts, dappAddress common.Address, input []byte,
) (
	*types.Transaction, error,
) {
	tx, err := c.inputBox.AddInput(signer, dappAddress, input)
	if err != nil {
		return nil, fmt.Errorf("failed to add input: %v", err)
	}
	return tx, nil
}

// Send dapp address with the dapp address relay contract.
func (c *ETHClient) SendDAppAddress(
	ctx context.Context, signer *bind.TransactOpts, dappAddress common.Address,
) (
	*types.Transaction, error,
) {
	tx, err := c.dappAddressRelay.RelayDAppAddress(signer, dappAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to send dapp address: %v", err)
	}
	return tx, nil
}

// Create a signer from a private key.
// This should only be used in a development environment.
func (c *ETHClient) CreateSigner(ctx context.Context, privateKey *ecdsa.PrivateKey,
) (
	*bind.TransactOpts, error,
) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	chainId, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

// Wait for transaction to be included in a block.
func (c *ETHClient) WaitForTransaction(ctx context.Context, tx *types.Transaction) error {
	for {
		_, isPending, err := c.client.TransactionByHash(ctx, tx.Hash())
		if err != nil {
			return fmt.Errorf("Fail to recover transaction: %v", err)
		}
		if !isPending {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// Get input index in the transaction by looking at the event logs.
func (c *ETHClient) GetInputIndex(ctx context.Context, tx *types.Transaction) (int, error) {
	receipt, err := c.client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return 0, fmt.Errorf("failed to get receipt: %v", err)
	}
	for _, log := range receipt.Logs {
		if log.Address != AddressInputBox {
			continue
		}
		inputAdded, err := c.inputBox.ParseInputAdded(*log)
		if err != nil {
			return 0, fmt.Errorf("failed to parse input added event: %v", err)
		}
		// We assume that int will fit all dapp inputs
		inputIndex := int(inputAdded.InputIndex.Int64())
		return inputIndex, nil
	}
	return 0, fmt.Errorf("input index not found")
}

// Get the number of inputs from the dapp.
func (c *ETHClient) GetNumberOfInputs(ctx context.Context, dappAddress common.Address) (int, error) {
	numInputs, err := c.inputBox.GetNumberOfInputs(nil, dappAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to get num inputs: %v", err)
	}
	return int(numInputs.Int64()), nil
}
