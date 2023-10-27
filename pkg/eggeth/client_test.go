// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const testTimeout = 300 * time.Second

func TestClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// We use the sunodo devnet docker image to test the client.
	// This image starts an anvil node with the Rollups contracts already deployed.
	req := testcontainers.ContainerRequest{
		Image: "sunodo/devnet:1.1.1",
		Cmd: []string{
			"anvil",
			"--block-time",
			"1",
			"--load-state",
			"/usr/share/sunodo/anvil_state.json",
		},
		Env: map[string]string{
			"ANVIL_IP_ADDR": "0.0.0.0",
		},
		ExposedPorts: []string{"8545/tcp"},
		WaitingFor:   wait.ForLog("Listening on 0.0.0.0:8545"),
	}
	anvilContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start anvil container: %v", err)
	}
	defer func() {
		if err := anvilContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate anvil container: %v", err)
		}
	}()

	// Set up ETHClient
	endpoint, err := anvilContainer.Endpoint(ctx, "ws")
	if err != nil {
		t.Fatalf("failed to get anvil endpoint: %v", err)
	}
	dappAddress := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	client, err := NewETHClient(endpoint, dappAddress)

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatalf("failed to get chain id: %v", err)
	}
	if chainID.Cmp(big.NewInt(31337)) != 0 {
		t.Fatalf("wrong chain id: %v", chainID)
	}

	// Set up signer
	mnemonic := "test test test test test test test test test test test junk"
	signer, err := NewMnemonicSigner(mnemonic, 0, chainID)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}

	// Test SendInput
	payload := common.Hex2Bytes("deadbeef")
	inputIndex, err := client.SendInput(ctx, signer, payload)
	if err != nil {
		t.Fatalf("failed to send input: %v", err)
	}
	if inputIndex != 0 {
		t.Fatalf("wrong input index: %v", inputIndex)
	}
	readSender, readPayload, err := getInput(client.inputBox, dappAddress, inputIndex)
	if err != nil {
		t.Fatal(err)
	}
	if readSender != common.HexToAddress("f39fd6e51aad88f6f4ce6ab8827279cfffb92266") {
		t.Fatalf("wrong sender: %x", readSender)
	}
	if !bytes.Equal(readPayload, payload) {
		t.Fatalf("wrong payload: %x", readPayload)
	}
}

func getInput(
	inputBox *InputBox, dappAddress common.Address, inputIndex int,
) (common.Address, []byte, error) {
	it, err := inputBox.FilterInputAdded(
		nil,
		[]common.Address{dappAddress},
		[]*big.Int{big.NewInt(int64(inputIndex))},
	)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to filter input added: %v", err)
	}
	defer it.Close()
	if !it.Next() {
		return common.Address{}, nil, fmt.Errorf("event not found")
	}
	return it.Event.Sender, it.Event.Input, nil
}
