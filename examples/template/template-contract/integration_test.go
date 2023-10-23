// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggtest"
	"testing"
	"time"
)

const testTimeout = 300 * time.Second

func TestTemplate(t *testing.T) {
	opts := eggtest.NewIntegrationTesterOpts()
	opts.LoadFromEnv()
	opts.Context = "../../.."
	opts.BuildTarget = "template-contract"

	tester := eggtest.NewIntegrationTester(t, opts)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	inputIndex, err := client.SendInput(ctx, signer, []byte("eggroll"))
	if err != nil {
		t.Fatalf("failed to send input: %v", err)
	}

	result, err := client.WaitFor(ctx, inputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}

	return_ := string(result.RawReturn())
	if return_ != "eggroll" {
		t.Fatalf("wrong result: %v", return_)
	}

	logs := result.Logs()
	if len(logs) != 1 || logs[0] != "received: eggroll" {
		t.Fatalf("wrong logs: %v", logs)
	}
}
