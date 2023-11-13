// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Provide testing tools for eggroll dapps.
package eggtest

import (
	"context"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/gligneul/eggroll/internal/sunodo"
)

// Integration test options.
type IntegrationTesterOpts struct {

	// Context of the sunodo Docker.
	DockerContext string

	// Target for sunodo build.
	BuildTarget string

	// If set, print sunodo Stderr.
	Verbose bool

	// If set, skip the integration test.
	Skip bool
}

// Return the default values for the options.
func MakeIntegrationTesterOpts() IntegrationTesterOpts {
	return IntegrationTesterOpts{
		DockerContext: ".",
		BuildTarget:   "",
		Verbose:       false,
		Skip:          true,
	}
}

// Load the some of the integration test opts from environment variables.
func LoadIntegrationTesterOpts() IntegrationTesterOpts {
	opts := MakeIntegrationTesterOpts()
	opts.Skip = os.Getenv("EGGTEST_RUN_INTEGRATION") == ""
	opts.Verbose = os.Getenv("EGGTEST_VERBOSE") != ""
	return opts
}

// Use sunodo to run integration tests.
// The tester will build the sunodo image, if necessary.
// Then, it will start the DApp contract with sunodo run.
type IntegrationTester struct {
	cancel context.CancelFunc
	done   <-chan struct{}
}

// Use mutex to make sure only runs one test at a time
var integrationMutex sync.Mutex

// Create a new sunodo tester.
// It is necessary to Close the tester at the end of the test.
func NewIntegrationTester(
	ctx context.Context, opts IntegrationTesterOpts, t *testing.T,
) *IntegrationTester {

	// Skip tests if set
	if opts.Skip {
		t.Skip("skipping integration test")
		return nil
	}

	// Check if sunodo is already running
	running, err := sunodo.IsRunning()
	if err != nil {
		t.Fatalf("failed to check if sunodo is running: %v", err)
	}
	if running {
		t.Fatalf("sunodo already running")
	}

	// Change current directly if necessary
	if opts.DockerContext != "." {
		err := os.Chdir(opts.DockerContext)
		if err != nil {
			t.Fatalf("change dir failed: %v", err)
		}
	}

	integrationMutex.Lock()

	t.Log("executing sunodo build")
	err = sunodo.Build(opts.BuildTarget, opts.Verbose)
	if err != nil {
		t.Fatalf("failed to execute sunodo build: %v", err)
	}

	// Setup sunodo config
	ready := make(chan struct{})
	var stderr io.Writer
	if opts.Verbose {
		stderr = os.Stderr
	}
	sunodoOpts := sunodo.RunOpts{
		NoBackend: false,
		Stdout:    os.Stdout,
		Stderr:    stderr,
		Ready:     ready,
	}

	// Start sunodo
	sunodoCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		t.Log("executing sunodo run")
		err := sunodo.Run(sunodoCtx, sunodoOpts)
		if err != nil && err != context.Canceled {
			t.Errorf("failed to execute sunodo run: %v", err)
		}
		done <- struct{}{}
	}()

	// Wait until it is ready
	select {
	case <-ready:
	case <-ctx.Done():
	}

	tester := &IntegrationTester{
		done:   done,
		cancel: cancel,
	}
	return tester
}

// Close the tester.
func (t *IntegrationTester) Close() error {
	t.cancel()
	<-t.done
	integrationMutex.Unlock()
	return nil
}
