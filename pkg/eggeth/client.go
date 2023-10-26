// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

//go:generate abigen --abi abi/CartesiDApp.json --pkg eggeth --type CartesiDApp --out cartesidapp.go
//go:generate abigen --abi abi/DAppAddressRelay.json --pkg eggeth --type DAppAddressRelay --out dappaddressrelay.go
//go:generate abigen --abi abi/ERC1155BatchPortal.json --pkg eggeth --type ERC1155BatchPortal --out erc1155batchportal.go
//go:generate abigen --abi abi/ERC1155SinglePortal.json --pkg eggeth --type ERC1155SinglePortal --out erc1155singleportal.go
//go:generate abigen --abi abi/ERC20Portal.json --pkg eggeth --type ERC20Portal --out erc20portal.go
//go:generate abigen --abi abi/ERC721Portal.json --pkg eggeth --type ERC721Portal --out erc721portal.go
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

	client              *ethclient.Client
	dappAddress         common.Address
	dappAddressRelay    *DAppAddressRelay
	erc1155BatchPortal  *ERC1155BatchPortal
	erc1155SinglePortal *ERC1155SinglePortal
	erc20Portal         *ERC20Portal
	erc721Portal        *ERC721Portal
	etherPortal         *EtherPortal
	inputBox            *InputBox
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
	erc1155BatchPortal, err := NewERC1155BatchPortal(AddressERC1155BatchPortal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC1155BatchPortal contract: %v", err)
	}
	erc1155SinglePortal, err := NewERC1155SinglePortal(AddressERC1155SinglePortal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC1155SinglePortal contract: %v", err)
	}
	erc20Portal, err := NewERC20Portal(AddressERC20Portal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC20Portal contract: %v", err)
	}
	erc721Portal, err := NewERC721Portal(AddressERC721Portal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC721Portal contract: %v", err)
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
		GasLimit:            30_000_000, // max gas
		client:              client,
		dappAddress:         dappAddress,
		dappAddressRelay:    dappAddressRelay,
		erc1155BatchPortal:  erc1155BatchPortal,
		erc1155SinglePortal: erc1155SinglePortal,
		erc20Portal:         erc20Portal,
		erc721Portal:        erc721Portal,
		etherPortal:         etherPortal,
		inputBox:            inputBox,
	}
	return ethClient, nil
}

// Get the chain ID.
func (c *ETHClient) ChainID(ctx context.Context) (*big.Int, error) {
	return c.client.ChainID(ctx)
}

// Send input to the blockchain.
func (c *ETHClient) SendInput(ctx context.Context, signer Signer, input []byte) (int, error) {
	return c.doSend(
		ctx, signer, big.NewInt(0),
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.inputBox.AddInput(txOpts, c.dappAddress, input)
		},
	)
}

// Send dapp address with the dapp address relay contract.
func (c *ETHClient) SendDAppAddress(ctx context.Context, signer Signer) (int, error) {
	return c.doSend(
		ctx, signer, big.NewInt(0),
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.dappAddressRelay.RelayDAppAddress(txOpts, c.dappAddress)
		},
	)
}

// Send assets to the Ether portal.
func (c *ETHClient) SendEther(
	ctx context.Context,
	signer Signer,
	txValue *big.Int,
	input []byte,
) (int, error) {
	return c.doSend(
		ctx, signer, txValue,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.etherPortal.DepositEther(txOpts, c.dappAddress, input)
		},
	)
}

// Send assets to the ERC20 portal.
func (c *ETHClient) SendERC20Tokens(
	ctx context.Context,
	signer Signer,
	token common.Address,
	amount *big.Int,
	input []byte,
) (int, error) {
	return c.doSend(
		ctx, signer, big.NewInt(0),
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.erc20Portal.DepositERC20Tokens(
				txOpts, token, c.dappAddress, amount, input)
		},
	)
}

// Send assets to the ERC721 portal.
func (c *ETHClient) SendERC721Token(
	ctx context.Context,
	signer Signer,
	token common.Address,
	tokenId *big.Int,
	baseLayerData []byte,
	input []byte,
) (int, error) {
	return c.doSend(
		ctx, signer, big.NewInt(0),
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.erc721Portal.DepositERC721Token(
				txOpts, token, c.dappAddress, tokenId, baseLayerData, input)
		},
	)
}

// Send assets ot the ERC1155 single portal.
func (c *ETHClient) SendSingleERC1155Token(
	ctx context.Context,
	signer Signer,
	token common.Address,
	tokenId *big.Int,
	value *big.Int,
	baseLayerData []byte,
	input []byte,
) (int, error) {
	return c.doSend(
		ctx, signer, big.NewInt(0),
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.erc1155SinglePortal.DepositSingleERC1155Token(
				txOpts, token, c.dappAddress, tokenId, value, baseLayerData, input)
		},
	)
}

// Send assets ot the ERC1155 batch portal.
func (c *ETHClient) SendBatchERC1155Tokens(
	ctx context.Context,
	signer Signer,
	token common.Address,
	tokenIds []*big.Int,
	values []*big.Int,
	baseLayerData []byte,
	input []byte,
) (int, error) {
	// Basic sanity check before sending the transaction.
	if len(tokenIds) == 0 {
		return 0, fmt.Errorf("no token ids")
	}
	if len(tokenIds) != len(values) {
		return 0, fmt.Errorf("tokenIds and values mismatch")
	}
	return c.doSend(
		ctx, signer, big.NewInt(0),
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.erc1155BatchPortal.DepositBatchERC1155Token(
				txOpts, token, c.dappAddress, tokenIds, values, baseLayerData, input)
		},
	)
}

//
// Private functions
//

// Send a transaction, wait for it, and get the input index.
func (c *ETHClient) doSend(
	ctx context.Context, signer Signer, txValue *big.Int,
	sender func(txOpts *bind.TransactOpts) (*types.Transaction, error)) (
	int, error) {

	txOpts, err := c.prepareTransaction(ctx, signer, txValue)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare transaction: %v", err)
	}
	tx, err := sender(txOpts)
	if err != nil {
		return 0, fmt.Errorf("failed to send dapp address: %v", err)
	}
	err = c.waitForTransaction(ctx, tx)
	if err != nil {
		return 0, err
	}
	return c.getInputIndex(ctx, tx)
}

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
