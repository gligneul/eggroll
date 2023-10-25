// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"testing"
	"time"

	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggtest"
)

const testTimeout = 300 * time.Second

func TestTemplate(t *testing.T) {
	opts := eggtest.NewIntegrationTesterOpts()
	opts.LoadFromEnv()
	opts.Context = "../../.."
	opts.BuildTarget = "minimal-contract"

	tester := eggtest.NewIntegrationTester(t, opts)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test advance
	inputIndex, err := client.SendInput(ctx, signer, []byte("eggroll"))
	if err != nil {
		t.Fatalf("failed to send input: %v", err)
	}
	result, err := client.WaitFor(ctx, inputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}
	logs := result.Logs()
	if len(logs) != 1 || logs[0].Message != "received advance: eggroll" {
		t.Fatalf("wrong logs: %#v", logs)
	}
	report := string(result.Reports[1].Payload)
	if report != "eggroll" {
		t.Fatalf("wrong report: %v", report)
	}

	// Test inspect
	inspectResult, err := client.Inspect(ctx, []byte("rollegg"))
	if err != nil {
		t.Fatalf("failed to inspect: %v", err)
	}
	logs = inspectResult.Logs()
	if len(logs) != 1 || logs[0].Message != "received inspect: rollegg" {
		t.Fatalf("wrong logs: %#v", logs)
	}
	report = string(inspectResult.Reports[1].Payload)
	if report != "rollegg" {
		t.Fatalf("wrong report: %v", report)
	}
}
