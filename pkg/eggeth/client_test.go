// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/pkg/eggeth/bindings"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const testTimeout = 300 * time.Second

func TestClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	anvilContainer := setupEthContainer(t, ctx)
	defer func() {
		if err := anvilContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate anvil container: %v", err)
		}
	}()
	endpoint, err := anvilContainer.Endpoint(ctx, "ws")
	if err != nil {
		t.Fatalf("failed to get anvil endpoint: %v", err)
	}

	// Deploy test contracts
	erc20Token, err := DeployTestERC20(ctx, endpoint)
	if err != nil {
		t.Fatalf("failed to deploy test erc20: %v", err)
	}
	t.Logf("ERC20 token: %v", erc20Token)

	// Set up ETHClient
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
	signer, err := NewMnemonicSigner(FoundryMnemonic, 0, chainID)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}

	// Test client
	testCases := []struct {
		name   string
		do     func(context.Context, *ETHClient, Signer) (int, error)
		sender common.Address
		input  []byte
	}{
		{
			name: "SendInput",
			do: func(ctx context.Context, c *ETHClient, s Signer) (int, error) {
				return c.SendInput(ctx, s, common.Hex2Bytes("deadbeef"))
			},
			sender: common.HexToAddress("f39fd6e51aad88f6f4ce6ab8827279cfffb92266"),
			input:  common.Hex2Bytes("deadbeef"),
		},
		{
			name: "SendDAppAddress",
			do: func(ctx context.Context, c *ETHClient, s Signer) (int, error) {
				return c.SendDAppAddress(ctx, s)
			},
			sender: AddressDAppAddressRelay,
			input:  client.dappAddress[:],
		},
		{
			name: "SendEther",
			do: func(ctx context.Context, c *ETHClient, s Signer) (int, error) {
				return c.SendEther(ctx, s, big.NewInt(65535), common.Hex2Bytes("deadbeef"))
			},
			sender: AddressEtherPortal,
			input: common.Hex2Bytes("" +
				// sender address
				"f39fd6e51aad88f6f4ce6ab8827279cfffb92266" +
				// value
				"000000000000000000000000000000000000000000000000000000000000ffff" +
				// payload
				"deadbeef",
			),
		},
		{
			name: "SendERC20Tokens",
			do: func(ctx context.Context, c *ETHClient, s Signer) (int, error) {
				amount := big.NewInt(65535)
				input := common.Hex2Bytes("deadbeef")
				return c.SendERC20Tokens(ctx, s, erc20Token, amount, input)
			},
			sender: AddressERC20Portal,
			input: common.Hex2Bytes("" +
				// success
				"01" +
				// token
				erc20Token.Hex()[2:] +
				// sender address
				"f39fd6e51aad88f6f4ce6ab8827279cfffb92266" +
				// amount
				"000000000000000000000000000000000000000000000000000000000000ffff" +
				// payload
				"deadbeef",
			),
		},
	}
	for i, testCase := range testCases {
		t.Logf("testing client.%v", testCase.name)
		inputIndex, err := testCase.do(ctx, client, signer)
		if err != nil {
			logContainerOutput(t, ctx, anvilContainer)
			t.Fatalf("failed to send: %v", err)
		}
		if inputIndex != i {
			t.Fatalf("wrong input index: %v; expected: %v", inputIndex, i)
		}
		readSender, readInput, err := getInput(client.inputBox, client.dappAddress, inputIndex)
		if err != nil {
			t.Fatal(err)
		}
		if readSender != testCase.sender {
			t.Fatalf("wrong sender: %x", readSender)
		}
		if !bytes.Equal(readInput, testCase.input) {
			t.Fatalf("wrong input: %x", readInput)
		}
	}
}

// We use the sunodo devnet docker image to test the client.
// This image starts an anvil node with the Rollups contracts already deployed.
func setupEthContainer(t *testing.T, ctx context.Context) testcontainers.Container {
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
	return anvilContainer
}

func getInput(
	inputBox *bindings.InputBox, dappAddress common.Address, inputIndex int,
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

func logContainerOutput(t *testing.T, ctx context.Context, container testcontainers.Container) {
	reader, err := container.Logs(ctx)
	if err != nil {
		t.Fatalf("failed to get reader: %v", err)
	}
	bytes, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read logs: %v", err)
	}
	t.Log(string(bytes))
}
