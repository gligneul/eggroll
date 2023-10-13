// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package textbox

import (
	"context"
	"testing"
	"time"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/eggtest"
)

const testTimeout = 300 * time.Second

func TestTextBox(t *testing.T) {
	tester := eggtest.NewIntegrationTester(t)
	defer tester.Close()

	client, err := eggroll.NewDevClient(Codecs())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	inputs := []any{
		&Append{Value: "egg"},
		&Append{Value: "roll"},
	}
	sendInputsAndVerifyState(t, client, inputs, "eggroll")

	inputs = []any{
		&Clear{},
		&Append{Value: "hi"},
	}
	sendInputsAndVerifyState(t, client, inputs, "hi")
}

func sendInputsAndVerifyState(t *testing.T, client *eggroll.DevClient,
	inputs []any, expectedState string) {

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var lastInputIndex int
	for _, input := range inputs {
		var err error
		lastInputIndex, err = client.SendInput(ctx, input)
		if err != nil {
			t.Fatalf("failed to send input: %v", err)
		}
	}

	r, err := client.WaitFor(ctx, lastInputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}

	textBox, ok := r.DecodeReturn().(*TextBox)
	if !ok {
		t.Fatalf("expected TextBox value")
	}
	if textBox.Value != expectedState {
		t.Fatalf("invalid value %v; expected %v", textBox.Value, expectedState)
	}
}
