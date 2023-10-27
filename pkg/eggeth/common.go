// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const DefaultGasLimit = 30_000_000

// Interface to create a signer for blockchain transactions.
type Signer interface {
	MakeTransactor() (*bind.TransactOpts, error)
	Account() common.Address
}

// Prepare the blockchain transaction.
func prepareTransaction(
	ctx context.Context,
	client *ethclient.Client,
	signer Signer,
	txValue *big.Int,
	gasLimit uint64,
) (*bind.TransactOpts, error) {

	nonce, err := client.PendingNonceAt(ctx, signer.Account())
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	tx, err := signer.MakeTransactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	tx.Nonce = big.NewInt(int64(nonce))
	tx.Value = txValue
	tx.GasLimit = gasLimit
	tx.GasPrice = gasPrice
	return tx, nil
}

// Wait for transaction to be included in a block.
func waitForTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction,
) error {
	for {
		_, isPending, err := client.TransactionByHash(ctx, tx.Hash())
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

func traceTransaction(ctx context.Context, client *ethclient.Client, hash common.Hash,
) (string, error) {
	// We make a call using the rpc client directly because this function
	// is not present in the ethclient struct. More details in:
	// https://github.com/ethereum/go-ethereum/issues/17341
	var result json.RawMessage
	err := client.Client().Call(&result, "debug_traceTransaction", hash)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
