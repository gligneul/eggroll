// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

// Use forge to build the contract; then, use abigen to generate the Go binding.
// We opt to deploy the contract with Go instead of forge because we don't want
// to add foundry as a runtime dependency for the EggRoll users.
//go:generate forge build --root contracts --extra-output-files abi

// TestERC20.sol
//go:generate sh -c "jq .abi < contracts/out/TestERC20.sol/TestERC20.json > contracts/out/TestERC20.sol/TestERC20.abi"
//go:generate sh -c "jq -r .bytecode.object < contracts/out/TestERC20.sol/TestERC20.json | cut -c 3- > contracts/out/TestERC20.sol/TestERC20.bin"
//go:generate abigen --bin contracts/out/TestERC20.sol/TestERC20.bin --abi contracts/out/TestERC20.sol/TestERC20.abi --pkg bindings --type TestERC20 --out bindings/test_erc20.go

// IERC20.sol
//go:generate sh -c "jq .abi < contracts/out/IERC20.sol/IERC20.json > contracts/out/IERC20.sol/IERC20.abi"
//go:generate abigen --abi contracts/out/IERC20.sol/IERC20.abi --pkg bindings --type IERC20 --out bindings/ierc20.go

// Cartesi contracts
//go:generate abigen --abi cartesi_abi/CartesiDApp.json --pkg bindings --type CartesiDApp --out bindings/cartesidapp.go
//go:generate abigen --abi cartesi_abi/DAppAddressRelay.json --pkg bindings --type DAppAddressRelay --out bindings/dappaddressrelay.go
//go:generate abigen --abi cartesi_abi/ERC1155BatchPortal.json --pkg bindings --type ERC1155BatchPortal --out bindings/erc1155batchportal.go
//go:generate abigen --abi cartesi_abi/ERC1155SinglePortal.json --pkg bindings --type ERC1155SinglePortal --out bindings/erc1155singleportal.go
//go:generate abigen --abi cartesi_abi/ERC20Portal.json --pkg bindings --type ERC20Portal --out bindings/erc20portal.go
//go:generate abigen --abi cartesi_abi/ERC721Portal.json --pkg bindings --type ERC721Portal --out bindings/erc721portal.go
//go:generate abigen --abi cartesi_abi/EtherPortal.json --pkg bindings --type EtherPortal --out bindings/etherportal.go
//go:generate abigen --abi cartesi_abi/InputBox.json --pkg bindings --type InputBox --out bindings/inputbox.go
//go:generate abigen --abi cartesi_abi/InputBox.json --pkg bindings --type InputBox --out bindings/inputbox.go

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

// Default gas limit when sending transactions.
const DefaultGasLimit = 30_000_000

// Dev mnemonic used by Foundry/Anvil.
const FoundryMnemonic = "test test test test test test test test test test test junk"

var (
	AddressCartesiDAppFactory  common.Address
	AddressDAppAddressRelay    common.Address
	AddressERC1155BatchPortal  common.Address
	AddressERC1155SinglePortal common.Address
	AddressERC20Portal         common.Address
	AddressERC721Portal        common.Address
	AddressEtherPortal         common.Address
	AddressInputBox            common.Address
	AddressSunodoToken         common.Address
)

func init() {
	AddressCartesiDAppFactory = common.HexToAddress("0x7122cd1221C20892234186facfE8615e6743Ab02")
	AddressDAppAddressRelay = common.HexToAddress("0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE")
	AddressERC1155BatchPortal = common.HexToAddress("0xedB53860A6B52bbb7561Ad596416ee9965B055Aa")
	AddressERC1155SinglePortal = common.HexToAddress("0x7CFB0193Ca87eB6e48056885E026552c3A941FC4")
	AddressERC20Portal = common.HexToAddress("0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB")
	AddressERC721Portal = common.HexToAddress("0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87")
	AddressEtherPortal = common.HexToAddress("0xFfdbe43d4c855BF7e0f105c400A50857f53AB044")
	AddressInputBox = common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768")
	AddressSunodoToken = common.HexToAddress("0xae7f61eCf06C65405560166b259C54031428A9C4")
}

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
