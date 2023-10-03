// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package blockchain

import (
	"context"
	"os/exec"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
	inputHex := hexutil.Encode(input)

	cmd := exec.Command(
		"sunodo",
		"send",
		"generic",
		"--dapp=0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C",
		"--chain-id=31337",
		"--rpc-url=http://127.0.0.1:8545",
		"--mnemonic-index=0",
		"--mnemonic-passphrase=test test test test test test test test test test test junk",
		"--input="+inputHex)

	return cmd.Run()
}
