// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package wallets

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

func TestEtherDepositString(t *testing.T) {
	sender := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	value := uint256.NewInt(100)
	deposit := &EtherDeposit{sender, *value}
	expectedString := "0xfafafafafafafafafafafafafafafafafafafafa deposited 100 wei"
	if deposit.String() != expectedString {
		t.Fatalf("wrong deposit string")
	}
}

func TestEtherAddresses(t *testing.T) {
	wallet := NewEtherWallet()
	addresses := wallet.Addresses()
	if len(addresses) != 0 {
		t.Fatalf("expected 0 addresses; got %v", len(addresses))
	}
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet.setBalance(address, uint256.NewInt(1))
	addresses = wallet.Addresses()
	if len(addresses) != 1 {
		t.Fatalf("expected 1 addresses; got %v", len(addresses))
	}
}

func TestEtherBalanceOf(t *testing.T) {
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet := NewEtherWallet()
	balance := wallet.BalanceOf(address)
	if !balance.IsZero() {
		t.Fatalf("expected 0 balance")
	}
	wallet.setBalance(address, uint256.NewInt(50))
	balance = wallet.BalanceOf(address)
	if !balance.Eq(uint256.NewInt(50)) {
		t.Fatalf("expected 50 balance")
	}
}

func TestValidEtherTransfer(t *testing.T) {
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet := NewEtherWallet()
	wallet.setBalance(src, uint256.NewInt(50))
	wallet.setBalance(dst, uint256.NewInt(50))
	err := wallet.Transfer(src, dst, uint256.NewInt(50))
	if err != nil {
		t.Fatalf("expected nil err; got %v", err)
	}
	srcBalance := wallet.BalanceOf(src)
	if !srcBalance.IsZero() {
		t.Fatalf("expected 0 balance in src")
	}
	dstBalance := wallet.BalanceOf(dst)
	if !dstBalance.Eq(uint256.NewInt(100)) {
		t.Fatalf("expected 100 balance in dst")
	}
}

func TestZeroEtherTransfer(t *testing.T) {
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet := NewEtherWallet()
	err := wallet.Transfer(src, dst, uint256.NewInt(0))
	if err != nil {
		t.Fatalf("expected nil err; got %v", err)
	}
	srcBalance := wallet.BalanceOf(src)
	if !srcBalance.IsZero() {
		t.Fatalf("expected 0 balance in src")
	}
	dstBalance := wallet.BalanceOf(dst)
	if !dstBalance.IsZero() {
		t.Fatalf("expected 0 balance in dst")
	}
}

func TestInsuficientFundsEtherTransfer(t *testing.T) {
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet := NewEtherWallet()
	wallet.setBalance(src, uint256.NewInt(50))
	err := wallet.Transfer(src, dst, uint256.NewInt(100))
	if err == nil {
		t.Fatalf("expected error; got nil")
	}
	if err.Error() != "insuficient funds" {
		t.Fatalf("wrong error message")
	}
}

func TestBalanceOverflowEtherTransfer(t *testing.T) {
	src := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	dst := common.HexToAddress("0xfefefefefefefefefefefefefefefefefefefefe")
	wallet := NewEtherWallet()
	wallet.setBalance(src, uint256.NewInt(50))
	wallet.setBalance(dst, IntMax)
	err := wallet.Transfer(src, dst, uint256.NewInt(50))
	if err == nil {
		t.Fatalf("expected error; got nil")
	}
	if err.Error() != "balance overflow" {
		t.Fatalf("wrong error message")
	}
}

func TestEtherWithdrawEncode(t *testing.T) {
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	value := uint256.NewInt(100)
	voucher := EncodeEtherWithdraw(address, value)
	expectedVoucher := common.Hex2Bytes("522f6815000000000000000000000000fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	if !reflect.DeepEqual(voucher, expectedVoucher) {
		t.Fatalf("got wrong voucher")
	}
}

func TestInsuficientFundsEtherWithdraw(t *testing.T) {
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet := NewEtherWallet()
	wallet.setBalance(address, uint256.NewInt(50))
	voucher, err := wallet.Withdraw(address, uint256.NewInt(100))
	if voucher != nil || err == nil {
		t.Fatalf("expected nil, err; got %v, %v", voucher, err)
	}
	if err.Error() != "insuficient funds" {
		t.Fatalf("wrong error message")
	}
	balance := wallet.BalanceOf(address)
	if !balance.Eq(uint256.NewInt(50)) {
		t.Fatalf("wrong balance; expected 50")
	}
}

func TestValidWithdraw(t *testing.T) {
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	wallet := NewEtherWallet()
	wallet.setBalance(address, uint256.NewInt(100))
	voucher, err := wallet.Withdraw(address, uint256.NewInt(100))
	if voucher == nil || err != nil {
		t.Fatalf("expected voucher, nil; got %v, %v", voucher, err)
	}
	balance := wallet.BalanceOf(address)
	if !balance.IsZero() {
		t.Fatalf("wrong balance; expected 0")
	}
}

func TestValidEtherDeposit(t *testing.T) {
	wallet := NewEtherWallet()
	payload := common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064deadbeef")
	deposit, input, err := wallet.deposit(payload)

	if err != nil {
		t.Fatalf("expected nil err; got %v", err)
	}

	if deposit == nil {
		t.Fatal("expected deposit; got nil")
	}
	etherDeposit := deposit.(*EtherDeposit)
	expectedAddress := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	if etherDeposit.Sender != expectedAddress {
		t.Fatal("wrong deposit address")
	}
	if !etherDeposit.Value.Eq(uint256.NewInt(100)) {
		t.Fatal("wrong deposit value; expected 100")
	}

	if input == nil {
		t.Fatal("expected input; got nil")
	}
	if common.Bytes2Hex(input) != "deadbeef" {
		t.Fatal("wrong input")
	}

	balance := wallet.BalanceOf(expectedAddress)
	if !balance.Eq(uint256.NewInt(100)) {
		t.Fatal("wrong balance; expected 100")
	}
}

func TestValidDepositWithEmptyInput(t *testing.T) {
	wallet := NewEtherWallet()
	payload := common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000000064")
	deposit, input, err := wallet.deposit(payload)

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

func TestOverflowDeposit(t *testing.T) {
	wallet := NewEtherWallet()

	// deposit int max
	payload := common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafaffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	deposit, input, err := wallet.deposit(payload)
	if deposit == nil || input == nil || err != nil {
		t.Fatalf("expected deposit, input, nil; got %v, %v, %v", deposit, input, err)
	}

	// deposit more ether
	payload = common.Hex2Bytes("fafafafafafafafafafafafafafafafafafafafa0000000000000000000000000000000000000000000000000000000000001000")
	deposit, input, err = wallet.deposit(payload)
	if deposit == nil || input == nil || err != nil {
		t.Fatalf("expected deposit, input, nil; got %v, %v, %v", deposit, input, err)
	}

	// check balance
	address := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	balance := wallet.BalanceOf(address)
	if !balance.Eq(IntMax) {
		t.Fatal("wrong balance; expected int max")
	}
}

func TestMalformedDeposit(t *testing.T) {
	wallet := NewEtherWallet()
	payload := common.Hex2Bytes("fafafa")
	deposit, input, err := wallet.deposit(payload)

	if err == nil {
		t.Fatal("expected err; got nil")
	}
	if err.Error() != "invalid eth deposit size; got 3" {
		t.Fatal("wrong err message")
	}
	if deposit != nil {
		t.Fatal("expected nil deposit; got something")
	}
	if input != nil {
		t.Fatal("expected nil input; got something")
	}
}
