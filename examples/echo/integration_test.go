// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"testing"
	"time"

	"github.com/gligneul/eggroll/pkg/eggroll"
	"github.com/gligneul/eggroll/pkg/eggtest"
	"github.com/gligneul/eggroll/pkg/eggtypes"
)

const testTimeout = 300 * time.Second

func TestTemplate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	opts := eggtest.LoadIntegrationTesterOpts()
	opts.DockerContext = "../.."
	opts.BuildTarget = "echo"
	tester := eggtest.NewIntegrationTester(ctx, opts, t)
	defer tester.Close()

	client, signer, err := eggroll.NewDevClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test advance
	inputIndex, err := client.Eth.SendInput(ctx, signer, EncodeAdvanceEcho("eggroll"))
	if err != nil {
		t.Fatalf("failed to send input: %v", err)
	}
	advanceResult, err := client.WaitFor(ctx, inputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}
	report, found := eggtypes.FindReport[EchoResponse](advanceResult.Reports, EchoResponseID)
	if !found {
		t.Fatalf("honeypot value not found")
	}
	if report.Value != "eggroll" {
		t.Fatalf("wrong report: %v", report)
	}

	// Test inspect
	inspectResult, err := client.Inspect(ctx, EncodeInspectEcho("rollegg"))
	if err != nil {
		t.Fatalf("failed to inspect: %v", err)
	}
	report, found = eggtypes.FindReport[EchoResponse](inspectResult.Reports, EchoResponseID)
	if !found {
		t.Fatalf("honeypot value not found")
	}
	if report.Value != "rollegg" {
		t.Fatalf("wrong report: %v", report)
	}
}
