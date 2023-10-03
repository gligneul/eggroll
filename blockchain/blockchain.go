// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package blockchain

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Implements blockchain client for Ethereum using go-ethereum.
type ETHClient struct {
	endpoint string
}

// Create new ETH client.
func NewETHClient(endpoint string) *ETHClient {
	client := &ETHClient{
		endpoint: endpoint,
	}
	return client
}

// Send input to the blockchain.
func (c *ETHClient) SendInput(ctx context.Context, input []byte) error {

	client, err := ethclient.Dial(c.endpoint)
	if err != nil {
		log.Fatal(err)
	}

	inputBoxAddress := common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768")
	bytecode, err := client.CodeAt(ctx, inputBoxAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(bytecode) == 0 {
		log.Fatal("input box not a smart contract")
	}

	log.Println(common.Bytes2Hex(bytecode))

	//account := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	// balance, err := client.BalanceAt(ctx, account, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(balance)

	return nil
}
