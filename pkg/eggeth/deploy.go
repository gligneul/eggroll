// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Use forge to build the contract; then, use abigen to generate the Go binding.
// We opt to deploy the contract with Go instead of forge because we don't want
// to add foundry as a runtime dependency for the EggRoll users.

//go:generate forge build --root contracts --extra-output-files abi
//go:generate sh -c "jq .abi < contracts/out/TestERC20.sol/TestERC20.json > contracts/out/TestERC20.sol/TestERC20.abi"
//go:generate sh -c "jq -r .bytecode.object < contracts/out/TestERC20.sol/TestERC20.json | cut -c 3- > contracts/out/TestERC20.sol/TestERC20.bin"
//go:generate abigen --bin contracts/out/TestERC20.sol/TestERC20.bin --abi contracts/out/TestERC20.sol/TestERC20.abi --pkg bindings --type TestERC20 --out bindings/test_erc20.go

package eggeth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gligneul/eggroll/pkg/eggeth/bindings"
	"github.com/gligneul/eggroll/pkg/eggwallets"
)

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

	// deploy
	txOpts, err := prepareTransaction(ctx, client, signer, big.NewInt(0), DefaultGasLimit)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to prepare transaction: %v", err)
	}
	address, tx, erc20, err := bindings.DeployTestERC20(txOpts, client, signer.Account())
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy: %v", err)
	}
	err = waitForTransaction(ctx, client, tx)
	if err != nil {
		return common.Address{}, err
	}
	receipt, err := client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get receipt: %v", err)
	}
	if receipt.Status == 0 {
		reason, err := traceTransaction(ctx, client, tx.Hash())
		if err != nil {
			return common.Address{}, fmt.Errorf("transaction failed; failed to get reason: %v", err)
		}
		return common.Address{}, fmt.Errorf("transaction failed: %v", reason)
	}

	// allowance to portal
	txOpts, err = prepareTransaction(ctx, client, signer, big.NewInt(0), DefaultGasLimit)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to prepare transaction: %v", err)
	}
	// TODO do not use eggwallets
	tx, err = erc20.Approve(txOpts, AddressERC20Portal, eggwallets.MaxUint256)
	err = waitForTransaction(ctx, client, tx)
	if err != nil {
		return common.Address{}, err
	}
	receipt, err = client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get receipt: %v", err)
	}
	if receipt.Status == 0 {
		reason, err := traceTransaction(ctx, client, tx.Hash())
		if err != nil {
			return common.Address{}, fmt.Errorf("transaction failed; failed to get reason: %v", err)
		}
		return common.Address{}, fmt.Errorf("transaction failed: %v", reason)
	}

	return address, nil
}
