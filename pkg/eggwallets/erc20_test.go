// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggwallets

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestERC20DepositString(t *testing.T) {
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	sender := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	value := big.NewInt(123)
	deposit := &ERC20Deposit{token, sender, value}
	expectedString := "0xfafafafafafafafafafafafafafafafafafafafa deposited 123 of 0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef token"
	if deposit.String() != expectedString {
		t.Fatalf("wrong deposit string: %v", deposit.String())
	}
}

func TestERC20Tokens(t *testing.T) {
	// test zero tokens
	wallet := NewERC20Wallet()
	tokens := wallet.Tokens()
	if len(tokens) != 0 {
		t.Fatalf("expected 0 tokens; got %v", len(tokens))
	}

	// test single token
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet.setBalance(token, address, big.NewInt(1))
	tokens = wallet.Tokens()
	expectedTokens := []common.Address{
		token,
	}
	if !reflect.DeepEqual(expectedTokens, tokens) {
		t.Fatalf("wrong tokens: %+v", tokens)
	}

	// test two tokens
	anotherToken := common.HexToAddress("0xfeebfeebfeebfeebfeebfeebfeebfeebfeebfeeb")
	anotherAddress := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet.setBalance(anotherToken, anotherAddress, big.NewInt(1))
	tokens = wallet.Tokens()
	expectedTokens = []common.Address{
		token,
		anotherToken,
	}
	if !reflect.DeepEqual(expectedTokens, tokens) {
		t.Fatalf("wrong tokens: %+v", tokens)
	}
}

func TestERC20Addresses(t *testing.T) {
	// test zero addresses
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	addresses := wallet.Addresses(token)
	if len(addresses) != 0 {
		t.Fatalf("expected 0 addresses; got %v", len(addresses))
	}

	// test single address
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet.setBalance(token, address, big.NewInt(1))
	addresses = wallet.Addresses(token)
	expectedAddresses := []common.Address{
		address,
	}
	if !reflect.DeepEqual(expectedAddresses, addresses) {
		t.Fatalf("wrong addresses: %+v", addresses)
	}

	// test two addresses
	anotherAddress := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet.setBalance(token, anotherAddress, big.NewInt(1))
	addresses = wallet.Addresses(token)
	expectedAddresses = []common.Address{
		address,
		anotherAddress,
	}
	if !reflect.DeepEqual(expectedAddresses, addresses) {
		t.Fatalf("wrong addresses: %+v", addresses)
	}

	// test adding balance to another token
	anotherToken := common.HexToAddress("0xfeebfeebfeebfeebfeebfeebfeebfeebfeebfeeb")
	yetAnotherAddress := common.HexToAddress("0xbabababababababababababababababababababa")
	wallet.setBalance(anotherToken, yetAnotherAddress, big.NewInt(1))
	addresses = wallet.Addresses(token)
	if !reflect.DeepEqual(expectedAddresses, addresses) {
		t.Fatalf("wrong addresses: %+v", addresses)
	}
}

func TestERC20BalanceOf(t *testing.T) {
	// test zero balance
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	balance := wallet.BalanceOf(token, address)
	if balance.Sign() != 0 {
		t.Fatalf("expected 0 balance")
	}

	// test non-zero balance
	wallet.setBalance(token, address, big.NewInt(50))
	balance = wallet.BalanceOf(token, address)
	if balance.Cmp(big.NewInt(50)) != 0 {
		t.Fatalf("expected 50 balance")
	}

	// test adding balance to another token
	anotherToken := common.HexToAddress("0xfeebfeebfeebfeebfeebfeebfeebfeebfeebfeeb")
	wallet.setBalance(anotherToken, address, big.NewInt(100))
	balance = wallet.BalanceOf(token, address)
	if balance.Cmp(big.NewInt(50)) != 0 {
		t.Fatalf("expected 50 balance")
	}

	// test setting balance to zero
	wallet.setBalance(token, address, big.NewInt(0))
	balance = wallet.BalanceOf(token, address)
	if balance.Sign() != 0 {
		t.Fatalf("expected 0 balance")
	}
}

func TestValidERC20Transfer(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet.setBalance(token, src, big.NewInt(50))
	wallet.setBalance(token, dst, big.NewInt(50))
	err := wallet.Transfer(token, src, dst, big.NewInt(50))
	if err != nil {
		t.Fatalf("expected nil err; got %v", err)
	}
	srcBalance := wallet.BalanceOf(token, src)
	if srcBalance.Sign() != 0 {
		t.Fatalf("expected 0 balance in src")
	}
	dstBalance := wallet.BalanceOf(token, dst)
	if dstBalance.Cmp(big.NewInt(100)) != 0 {
		t.Fatalf("expected 100 balance in dst")
	}
}

func TestZeroERC20Transfer(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	err := wallet.Transfer(token, src, dst, big.NewInt(0))
	if err != nil {
		t.Fatalf("expected nil err; got %v", err)
	}
	srcBalance := wallet.BalanceOf(token, src)
	if srcBalance.Sign() != 0 {
		t.Fatalf("expected 0 balance in src")
	}
	dstBalance := wallet.BalanceOf(token, dst)
	if dstBalance.Sign() != 0 {
		t.Fatalf("expected 0 balance in dst")
	}
}

func TestSelfERC20Transfer(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet.setBalance(token, src, big.NewInt(50))
	err := wallet.Transfer(token, src, src, big.NewInt(50))
	if err == nil {
		t.Fatalf("expected error; got nil")
	}
	if err.Error() != "can't transfer to self" {
		t.Fatalf("wrong error message")
	}
}

func TestInsuficientFundsERC20Transfer(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet.setBalance(token, src, big.NewInt(50))
	err := wallet.Transfer(token, src, dst, big.NewInt(100))
	if err == nil {
		t.Fatalf("expected error; got nil")
	}
	if err.Error() != "insuficient funds" {
		t.Fatalf("wrong error message")
	}
}

func TestBalanceOverflowERC20Transfer(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet.setBalance(token, src, big.NewInt(50))
	wallet.setBalance(token, dst, MaxUint256)
	err := wallet.Transfer(token, src, dst, big.NewInt(50))
	if err == nil {
		t.Fatalf("expected error; got nil")
	}
	if err.Error() != "balance overflow" {
		t.Fatalf("wrong error message")
	}
}

func TestERC20WithdrawEncode(t *testing.T) {
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	value := big.NewInt(100)
	voucher := EncodeERC20Withdraw(token, value)
	expectedVoucher := common.Hex2Bytes("a9059cbb000000000000000000000000fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	if !bytes.Equal(voucher, expectedVoucher) {
		t.Fatalf("got wrong voucher: %x", voucher)
	}
}

func TestValidERC20Withdraw(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet.setBalance(token, address, big.NewInt(100))
	voucher, err := wallet.Withdraw(token, address, big.NewInt(100))
	if voucher == nil || err != nil {
		t.Fatalf("expected voucher, nil; got %v, %v", voucher, err)
	}
	balance := wallet.BalanceOf(token, address)
	if balance.Sign() != 0 {
		t.Fatalf("wrong balance; expected 0")
	}
}

func TestInsuficientFundsERC20Withdraw(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet.setBalance(token, address, big.NewInt(50))
	voucher, err := wallet.Withdraw(token, address, big.NewInt(100))
	if voucher != nil || err == nil {
		t.Fatalf("expected nil, err; got %v, %v", voucher, err)
	}
	if err.Error() != "insuficient funds" {
		t.Fatalf("wrong error message")
	}
	balance := wallet.BalanceOf(token, address)
	if balance.Cmp(big.NewInt(50)) != 0 {
		t.Fatalf("wrong balance; expected 50")
	}
}

func TestValidERC20Deposit(t *testing.T) {
	wallet := NewERC20Wallet()
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	payload := common.Hex2Bytes("01beefbeefbeefbeefbeefbeefbeefbeefbeefbeeffafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064deadbeef")
	deposit, input, err := wallet.Deposit(payload)
	if err != nil {
		t.Fatalf("expected nil err; got %v", err)
	}
	if deposit == nil {
		t.Fatal("expected deposit; got nil")
	}

	// check deposit
	erc20Deposit := deposit.(*ERC20Deposit)
	if erc20Deposit.Token != token {
		t.Fatalf("wrong token: %v", erc20Deposit.Token)
	}
	if erc20Deposit.Sender != address {
		t.Fatalf("wrong sender: %v", erc20Deposit.Sender)
	}
	if erc20Deposit.Amount.Cmp(big.NewInt(100)) != 0 {
		t.Fatal("wrong amount; expected 100")
	}

	// check input data
	if input == nil {
		t.Fatal("expected input; got nil")
	}
	if common.Bytes2Hex(input) != "deadbeef" {
		t.Fatal("wrong input")
	}

	// check balance
	balance := wallet.BalanceOf(token, address)
	if balance.Cmp(big.NewInt(100)) != 0 {
		t.Fatal("wrong balance; expected 100")
	}
}

func TestValidERC20DepositWithEmptyInput(t *testing.T) {
	wallet := NewERC20Wallet()
	payload := common.Hex2Bytes("01beefbeefbeefbeefbeefbeefbeefbeefbeefbeeffafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	deposit, input, err := wallet.Deposit(payload)
	if err != nil {
		t.Fatalf("expected nil err; got %v", err)
	}
	if deposit == nil {
		t.Fatal("expected deposit; got nil")
	}
	if len(input) != 0 {
		t.Fatalf("expected empty input; got %v", len(input))
	}
}

func TestOverflowERC20Deposit(t *testing.T) {
	wallet := NewERC20Wallet()

	// deposit int max
	payload := common.Hex2Bytes("01beefbeefbeefbeefbeefbeefbeefbeefbeefbeeffafafafafafafafafafafafafafafafafafafafaffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	deposit, input, err := wallet.Deposit(payload)
	if deposit == nil || input == nil || err != nil {
		t.Fatalf("expected deposit, input, nil; got %v, %v, %v", deposit, input, err)
	}

	// deposit more ether
	payload = common.Hex2Bytes("01beefbeefbeefbeefbeefbeefbeefbeefbeefbeeffafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000001000")
	deposit, input, err = wallet.Deposit(payload)
	if deposit == nil || input == nil || err != nil {
		t.Fatalf("expected deposit, input, nil; got %v, %v, %v", deposit, input, err)
	}

	// check balance
	token := common.HexToAddress("0xbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef")
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	balance := wallet.BalanceOf(token, address)
	if balance.Cmp(MaxUint256) != 0 {
		t.Fatal("wrong balance; expected int max")
	}
}

func TestMalformedERC20Deposit(t *testing.T) {
	wallet := NewERC20Wallet()
	payload := common.Hex2Bytes("fafafa")
	deposit, input, err := wallet.Deposit(payload)
	if err == nil {
		t.Fatal("expected err; got nil")
	}
	if err.Error() != "invalid erc20 deposit size; got 3" {
		t.Fatalf("wrong err message: %v", err)
	}
	if deposit != nil {
		t.Fatal("expected nil deposit; got something")
	}
	if input != nil {
		t.Fatal("expected nil input; got something")
	}
}

func TestFailedERC20Deposit(t *testing.T) {
	wallet := NewERC20Wallet()
	payload := common.Hex2Bytes("00beefbeefbeefbeefbeefbeefbeefbeefbeefbeeffafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	deposit, input, err := wallet.Deposit(payload)
	if err == nil {
		t.Fatal("expected err; got nil")
	}
	if err.Error() != "received failed erc20 transfer" {
		t.Fatalf("wrong err message: %v", err)
	}
	if deposit != nil {
		t.Fatal("expected nil deposit; got something")
	}
	if input != nil {
		t.Fatal("expected nil input; got something")
	}
}
