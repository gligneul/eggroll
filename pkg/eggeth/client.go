// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gligneul/eggroll/pkg/eggeth/bindings"
)

// Implements blockchain client for Ethereum using go-ethereum.
// This struct provides methods that are specific for the Cartesi Rollups.
type ETHClient struct {

	// Gas limit when building transactions.
	GasLimit uint64

	client              *ethclient.Client
	dappAddress         common.Address
	dappAddressRelay    *bindings.DAppAddressRelay
	erc1155BatchPortal  *bindings.ERC1155BatchPortal
	erc1155SinglePortal *bindings.ERC1155SinglePortal
	erc20Portal         *bindings.ERC20Portal
	erc721Portal        *bindings.ERC721Portal
	etherPortal         *bindings.EtherPortal
	inputBox            *bindings.InputBox
}

// Create new ETH client.
func NewETHClient(endpoint string, dappAddress common.Address) (*ETHClient, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %v", err)
	}
	dappAddressRelay, err := bindings.NewDAppAddressRelay(AddressDAppAddressRelay, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAppAddressRelaya contract: %v", err)
	}
	erc1155BatchPortal, err := bindings.NewERC1155BatchPortal(AddressERC1155BatchPortal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC1155BatchPortal contract: %v", err)
	}
	erc1155SinglePortal, err := bindings.NewERC1155SinglePortal(AddressERC1155SinglePortal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC1155SinglePortal contract: %v", err)
	}
	erc20Portal, err := bindings.NewERC20Portal(AddressERC20Portal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC20Portal contract: %v", err)
	}
	erc721Portal, err := bindings.NewERC721Portal(AddressERC721Portal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ERC721Portal contract: %v", err)
	}
	etherPortal, err := bindings.NewEtherPortal(AddressEtherPortal, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EtherPortal contract: %v", err)
	}
	inputBox, err := bindings.NewInputBox(AddressInputBox, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InputBox contract: %v", err)
	}
	ethClient := &ETHClient{
		GasLimit:            DefaultGasLimit,
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

// Send the input to the DApp contract.
// This function waits until the transaction is added to a block and return the input index.
func (c *ETHClient) SendInput(ctx context.Context, signer Signer, input []byte) (int, error) {
	receipt, err := sendTransaction(
		ctx, c.client, signer, big.NewInt(0), c.GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.inputBox.AddInput(txOpts, c.dappAddress, input)
		},
	)
	if err != nil {
		return 0, err
	}
	return c.getInputIndex(ctx, receipt)
}

// Send the DApp address to the DApp contract with the DAppAddressRelay contract.
// This function waits until the transaction is added to a block and return the input index.
func (c *ETHClient) SendDAppAddress(ctx context.Context, signer Signer) (int, error) {
	receipt, err := sendTransaction(
		ctx, c.client, signer, big.NewInt(0), c.GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.dappAddressRelay.RelayDAppAddress(txOpts, c.dappAddress)
		},
	)
	if err != nil {
		return 0, err
	}
	return c.getInputIndex(ctx, receipt)
}

// Send Ether to the Ether portal. This function also receives an optional input.
// This function waits until the transaction is added to a block and return the input index.
func (c *ETHClient) SendEther(
	ctx context.Context,
	signer Signer,
	txValue *big.Int,
	input []byte,
) (int, error) {
	receipt, err := sendTransaction(
		ctx, c.client, signer, txValue, c.GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.etherPortal.DepositEther(txOpts, c.dappAddress, input)
		},
	)
	if err != nil {
		return 0, err
	}
	return c.getInputIndex(ctx, receipt)
}

// Send an ERC20 token to the ERC20 portal. This function also receives an optional input.
// This function waits until the transaction is added to a block and return the input index.
func (c *ETHClient) SendERC20Tokens(
	ctx context.Context,
	signer Signer,
	token common.Address,
	amount *big.Int,
	input []byte,
) (int, error) {
	receipt, err := sendTransaction(
		ctx, c.client, signer, big.NewInt(0), c.GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return c.erc20Portal.DepositERC20Tokens(
				txOpts, token, c.dappAddress, amount, input)
		},
	)
	if err != nil {
		return 0, err
	}
	return c.getInputIndex(ctx, receipt)
}

// Send assets to the ERC721 portal.
// func (c *ETHClient) SendERC721Token(
// 	ctx context.Context,
// 	signer Signer,
// 	token common.Address,
// 	tokenId *big.Int,
// 	baseLayerData []byte,
// 	input []byte,
// ) (int, error) {
// 	receipt, err := sendTransaction(
// 		ctx, c.client, signer, big.NewInt(0), c.GasLimit,
// 		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
// 			return c.erc721Portal.DepositERC721Token(
// 				txOpts, token, c.dappAddress, tokenId, baseLayerData, input)
// 		},
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return c.getInputIndex(ctx, receipt)
// }

// Send assets ot the ERC1155 single portal.
// func (c *ETHClient) SendSingleERC1155Token(
// 	ctx context.Context,
// 	signer Signer,
// 	token common.Address,
// 	tokenId *big.Int,
// 	value *big.Int,
// 	baseLayerData []byte,
// 	input []byte,
// ) (int, error) {
// 	receipt, err := sendTransaction(
// 		ctx, c.client, signer, big.NewInt(0), c.GasLimit,
// 		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
// 			return c.erc1155SinglePortal.DepositSingleERC1155Token(
// 				txOpts, token, c.dappAddress, tokenId, value, baseLayerData, input)
// 		},
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return c.getInputIndex(ctx, receipt)
// }

// Send assets ot the ERC1155 batch portal.
// func (c *ETHClient) SendBatchERC1155Tokens(
// 	ctx context.Context,
// 	signer Signer,
// 	token common.Address,
// 	tokenIds []*big.Int,
// 	values []*big.Int,
// 	baseLayerData []byte,
// 	input []byte,
// ) (int, error) {
// 	// Basic sanity check before sending the transaction.
// 	if len(tokenIds) == 0 {
// 		return 0, fmt.Errorf("no token ids")
// 	}
// 	if len(tokenIds) != len(values) {
// 		return 0, fmt.Errorf("tokenIds and values mismatch")
// 	}
// 	receipt, err := sendTransaction(
// 		ctx, c.client, signer, big.NewInt(0), c.GasLimit,
// 		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
// 			return c.erc1155BatchPortal.DepositBatchERC1155Token(
// 				txOpts, token, c.dappAddress, tokenIds, values, baseLayerData, input)
// 		},
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return c.getInputIndex(ctx, receipt)
// }

// Get input index in the transaction by looking at the event logs.
func (c *ETHClient) getInputIndex(ctx context.Context, receipt *types.Receipt) (int, error) {
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
