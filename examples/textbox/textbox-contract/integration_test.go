// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"github.com/gligneul/eggroll/pkg/eggeth"
	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggtest"
	"testing"
	"time"
)

const testTimeout = 300 * time.Second

func TestTextBox(t *testing.T) {
	opts := eggtest.NewIntegrationTesterOpts()
	opts.LoadFromEnv()
	opts.Context = "../../.."
	opts.BuildTarget = "textbox-contract"

	tester := eggtest.NewIntegrationTester(t, opts)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, Codecs())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	inputs := []any{
		&Append{Value: "egg"},
		&Append{Value: "roll"},
	}
	sendInputsAndVerifyState(t, ctx, client, signer, inputs, "eggroll")

	inputs = []any{
		&Clear{},
		&Append{Value: "hi"},
	}
	sendInputsAndVerifyState(t, ctx, client, signer, inputs, "hi")
}

func sendInputsAndVerifyState(
	t *testing.T, ctx context.Context, client *eggroll.Client,
	signer eggeth.Signer, inputs []any, expectedState string) {

	var lastInputIndex int
	for _, input := range inputs {
		var err error
		lastInputIndex, err = client.SendInput(ctx, signer, input)
		if err != nil {
			t.Fatalf("failed to send input: %v", err)
		}
	}

	result, err := client.WaitFor(ctx, lastInputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}

	textBox, ok := client.DecodeReturn(result).(*TextBox)
	if !ok {
		t.Fatalf("expected TextBox value")
	}
	if textBox.Value != expectedState {
		t.Fatalf("invalid value %v; expected %v", textBox.Value, expectedState)
	}
}
