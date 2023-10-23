// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Provide testing tools for eggroll dapps.
package eggtest

import (
	"os"
	"sync"
	"testing"

	"github.com/gligneul/eggroll/internal/sunodo"
)

// Integration test options.
type IntegrationTesterOpts struct {

	// Context of the sunodo Docker (default: ".").
	Context string

	// Target for sunodo build (default: "").
	BuildTarget string

	// If set, increases the verbosity of the test (default: false).
	Verbose bool

	// If set, skip the integration test (default: false).
	Skip bool
}

// Create integration test options with the default values.
func NewIntegrationTesterOpts() *IntegrationTesterOpts {
	return &IntegrationTesterOpts{
		Context:     ".",
		BuildTarget: "",
		Verbose:     false,
		Skip:        false,
	}
}

// Load the some of the integration test opts from environment variables.
func (opts *IntegrationTesterOpts) LoadFromEnv() {
	if os.Getenv("EGGTEST_SKIP_INTEGRATION") != "" {
		opts.Skip = true
	}
	if os.Getenv("EGGTEST_VERBOSE") != "" {
		opts.Verbose = true
	}
}

// Use sunodo to run integration tests.
// The tester will build the sunodo image, if necessary.
// Then, it will start the DApp contract with sunodo run.
type IntegrationTester struct {
	*testing.T
	session *sunodo.Session
}

// Use mutex to make sure only runs one test at a time
var integrationMutex sync.Mutex

// Create a new sunodo tester.
// It is necessary to Close the tester at the end of the test.
func NewIntegrationTester(t *testing.T, opts *IntegrationTesterOpts) *IntegrationTester {
	if opts == nil {
		opts = NewIntegrationTesterOpts()
	}
	if opts.Skip {
		t.Skip("skipping integration test")
		return nil
	}
	if opts.Context != "." {
		os.Chdir(opts.Context)
	}

	t.Log("executing sunodo build")
	err := sunodo.Build(opts.BuildTarget, opts.Verbose)
	if err != nil {
		t.Fatalf("failed to execute sunodo build: %v", err)
	}

	t.Log("executing sunodo run")
	session, err := sunodo.Run(opts.Verbose)
	if err != nil {
		t.Fatalf("failed to execute sunodo run: %v", err)
	}

	tester := &IntegrationTester{
		T:       t,
		session: session,
	}
	integrationMutex.Lock()
	return tester
}

// Close the tester.
func (t *IntegrationTester) Close() error {
	if err := t.session.Close(); err != nil {
		t.Errorf("failed to close sunodo session: %v", err)
	}
	integrationMutex.Unlock()
	return nil
}
