// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"testing"
	"time"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/eggtest"
)

const testTimeout = 300 * time.Second

func TestInspect(t *testing.T) {
	opts := eggtest.NewIntegrationTesterOpts()
	opts.LoadFromEnv()
	opts.Context = "../../.."
	opts.BuildTarget = "inspect-contract"

	tester := eggtest.NewIntegrationTester(t, opts)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	client, _, err := eggroll.NewDevClient(ctx, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	response, err := client.Inspect(ctx, []byte("eggroll"))
	if err != nil {
		t.Fatalf("failed to inspect: %v", err)
	}
	return_ := string(response.RawReturn())
	if return_ != "eggroll" {
		t.Fatalf("wrong return: %v", return_)
	}
	if response.ProcessedInputCount != 0 {
		t.Fatalf("wrong input count: %v", response.ProcessedInputCount)
	}
}
