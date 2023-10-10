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

	client := eggroll.NewLocalClient()

	inputs := []any{
		Append{Value: "egg"},
		Append{Value: "roll"},
	}
	lastInputIndex := 1
	sendInputsAndVerifyState(t, client, inputs, lastInputIndex, "eggroll")

	inputs = []any{
		Clear{},
		Append{Value: "hi"},
	}
	lastInputIndex = 3
	sendInputsAndVerifyState(t, client, inputs, lastInputIndex, "hi")
}

func sendInputsAndVerifyState(
	t *testing.T, client *eggroll.Client,
	inputs []any, lastInputIndex int,
	expectedState string) {

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	for _, input := range inputs {
		err := client.SendGeneric(ctx, input)
		if err != nil {
			t.Fatalf("failed to send input: %v", err)
		}
	}

	r, err := client.WaitFor(ctx, lastInputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}

	if string(r.Result) != expectedState {
		t.Fatalf("invalid state: '%v'; expected '%v'", string(r.Result), expectedState)
	}
}
