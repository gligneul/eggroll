// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"testing"
	"time"

	"github.com/gligneul/eggroll/pkg/eggeth"
	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggtest"
	"github.com/gligneul/eggroll/pkg/eggtypes"
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

	client, signer, err := eggroll.NewDevClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	inputs := []eggtypes.Packer{
		&Append{Value: "egg"},
		&Append{Value: "roll"},
	}
	sendInputsAndVerifyState(t, ctx, client, signer, inputs, "eggroll")

	inputs = []eggtypes.Packer{
		&Clear{},
		&Append{Value: "hi"},
	}
	sendInputsAndVerifyState(t, ctx, client, signer, inputs, "hi")
}

func sendInputsAndVerifyState(
	t *testing.T, ctx context.Context, client *eggroll.Client,
	signer eggeth.Signer, inputs []eggtypes.Packer, expectedState string) {

	var lastInputIndex int
	for _, input := range inputs {
		var err error
		lastInputIndex, err = client.SendInput(ctx, signer, input.Pack())
		if err != nil {
			t.Fatalf("failed to send input: %v", err)
		}
	}

	result, err := client.WaitFor(ctx, lastInputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}

	textBox, found := eggtypes.FindReport[TextBox](result.Reports, TextBoxID)
	if !found {
		t.Fatalf("textbox value not found")
	}
	if textBox.Value != expectedState {
		t.Fatalf("invalid value %v; expected %v", textBox.Value, expectedState)
	}
}
