// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package wallets

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type testDeposit struct {
	sender common.Address
}

func (d *testDeposit) GetSender() common.Address {
	return d.sender
}

func (w *testDeposit) String() string {
	return "testDeposit"
}

type testWallet struct {
	portal common.Address
	sender common.Address
	err    error
}

func (w *testWallet) deposit(input []byte) (Deposit, []byte, error) {
	return &testDeposit{w.sender}, input, w.err
}

func (w *testWallet) portalAddress() common.Address {
	return w.portal
}

func TestWalletDispatcher(t *testing.T) {
	testSender := common.HexToAddress("0xFAFAFAFAFAFAFAFAFAFAFAFAFAFAFAFAFAFAFAFA")
	testPortal := common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	testErr := errors.New("error")
	testWallet := &testWallet{
		testPortal,
		testSender,
		testErr,
	}

	d := &WalletDispatcher{
		handlers: map[common.Address]Wallet{
			testPortal: testWallet,
		},
	}

	// call dispatch for existing wallet
	testInput := []byte{1, 2, 3}
	deposit, input, err := d.Dispatch(testPortal, testInput)
	if deposit == nil || deposit.GetSender() != testSender {
		t.Fatalf("wrong deposit from dispatch")
	}
	if !reflect.DeepEqual(input, testInput) {
		t.Fatal("wrong input from dispatch")
	}
	if err == nil {
		t.Fatal("wrong err from dispatch")
	}

	// call dispatch for non-registered portal
	deposit, input, err = d.Dispatch(common.Address{}, testInput)
	if deposit != nil || input != nil || err != nil {
		t.Fatal("wrong return from dispatch, expected nil")
	}
}
