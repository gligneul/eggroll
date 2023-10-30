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

// Deploy a ERC20 contract for local testing.
// This function uses the account 0 from the foundry mnemonic to deploy the contract.
// It returns the the address of the ERC20 contract.
func DeployTestERC20(ctx context.Context, endpoint string) (common.Address, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to connect to Ethereum: %v", err)
	}

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get chain id: %v", err)
	}

	signer, err := NewMnemonicSigner(FoundryMnemonic, 0, chainId)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create signer: %v", err)
	}

	var address common.Address
	_, err = sendTransaction(
		ctx, client, signer, big.NewInt(0), DefaultGasLimit,
		func(txOpts *bind.TransactOpts) (tx *types.Transaction, err error) {
			address, tx, _, err = bindings.DeployTestERC20(
				txOpts, client, signer.Account())
			return tx, err
		},
	)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}
