// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

//go:generate abigen --abi abi/CartesiDApp.json --pkg eggeth --type CartesiDApp --out cartesidapp.go
//go:generate abigen --abi abi/DAppAddressRelay.json --pkg eggeth --type DAppAddressRelay --out dappaddressrelay.go
//go:generate abigen --abi abi/EtherPortal.json --pkg eggeth --type EtherPortal --out etherportal.go
//go:generate abigen --abi abi/InputBox.json --pkg eggeth --type InputBox --out inputbox.go

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Interface to create a signer for blockchain transactions.
type Signer interface {
	MakeTransactor() (*bind.TransactOpts, error)
	Account() common.Address
}

// Implements blockchain client for Ethereum using go-ethereum.
// This struct provides methods that are specific for the Cartesi Rollups.
type ETHClient struct {

	// Gas limit when building transactions.
	GasLimit uint64

	client           *ethclient.Client
	dappAddress      common.Address
	dappAddressRelay *DAppAddressRelay
	etherPortal      *EtherPortal
	inputBox         *InputBox
}

// Create new ETH client.
func NewETHClient(endpoint string, dappAddress common.Address) (*ETHClient, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %v", err)
	}
	dappAddressRelay, err := NewDAppAddressRelay(AddressDAppAddressRelay, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAppAddressRelaya contract: %v", err)
	}
	etherPortal, err := NewEtherPortal(AddressEtherPortal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EtherPortal contract: %v", err)
	}
	inputBox, err := NewInputBox(AddressInputBox, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InputBox contract: %v", err)
	}
	ethClient := &ETHClient{
		GasLimit:         30_000_000, // max gas
		client:           client,
		dappAddress:      dappAddress,
		dappAddressRelay: dappAddressRelay,
		etherPortal:      etherPortal,
		inputBox:         inputBox,
	}
	return ethClient, nil
}

// Get the chain ID.
func (c *ETHClient) ChainID(ctx context.Context) (*big.Int, error) {
	return c.client.ChainID(ctx)
}

// Send input to the blockchain.
func (c *ETHClient) SendInput(ctx context.Context, signer Signer, input []byte) (int, error) {
	txOpts, err := c.prepareTransaction(ctx, signer, big.NewInt(0))
	if err != nil {
		return 0, fmt.Errorf("failed to prepare transaction: %v", err)
	}
	tx, err := c.inputBox.AddInput(txOpts, c.dappAddress, input)
	if err != nil {
		return 0, fmt.Errorf("failed to add input: %v", err)
	}
	return c.getInputIndex(ctx, tx)
}

// Send dapp address with the dapp address relay contract.
func (c *ETHClient) SendDAppAddress(ctx context.Context, signer Signer) (int, error) {
	txOpts, err := c.prepareTransaction(ctx, signer, big.NewInt(0))
	if err != nil {
		return 0, fmt.Errorf("failed to prepare transaction: %v", err)
	}
	tx, err := c.dappAddressRelay.RelayDAppAddress(txOpts, c.dappAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to send dapp address: %v", err)
	}
	return c.getInputIndex(ctx, tx)
}

// Send funds to the Ether portal.
func (c *ETHClient) SendEther(ctx context.Context, signer Signer, txValue *big.Int, input []byte) (
	int, error) {

	txOpts, err := c.prepareTransaction(ctx, signer, txValue)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare transaction: %v", err)
	}
	tx, err := c.etherPortal.DepositEther(txOpts, c.dappAddress, input)
	if err != nil {
		return 0, fmt.Errorf("failed to send dapp address: %v", err)
	}
	return c.getInputIndex(ctx, tx)
}

//
// Private functions
//

// Prepare the blockchain transaction.
func (c *ETHClient) prepareTransaction(ctx context.Context, signer Signer, txValue *big.Int) (
	*bind.TransactOpts, error) {

	nonce, err := c.client.PendingNonceAt(ctx, signer.Account())
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	tx, err := signer.MakeTransactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	tx.Nonce = big.NewInt(int64(nonce))
	tx.Value = txValue
	tx.GasLimit = c.GasLimit
	tx.GasPrice = gasPrice
	return tx, nil
}

// Wait for transaction to be included in a block.
func (c *ETHClient) waitForTransaction(ctx context.Context, tx *types.Transaction) error {
	for {
		_, isPending, err := c.client.TransactionByHash(ctx, tx.Hash())
		if err != nil {
			return fmt.Errorf("fail to recover transaction: %v", err)
		}
		if !isPending {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			continue
		}
	}
	return nil
}

// Get input index in the transaction by looking at the event logs.
func (c *ETHClient) getInputIndex(ctx context.Context, tx *types.Transaction) (int, error) {
	err := c.waitForTransaction(ctx, tx)
	if err != nil {
		return 0, err
	}
	receipt, err := c.client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return 0, fmt.Errorf("failed to get receipt: %v", err)
	}
	if receipt.Status == 0 {
		return 0, fmt.Errorf("transaction failed")
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
