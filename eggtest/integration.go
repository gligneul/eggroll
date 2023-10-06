// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Provide testing tools for eggroll dapps.
package eggtest

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

// Use mutex to make sure only runs one test at a time
var integrationMutex sync.Mutex

// Use sunodo to run integration tests.
// The tester will build the sunodo image, if necessary.
// Then, it will start the DApp contract with sunodo run.
type IntegrationTester struct {
	*testing.T
	cmd *exec.Cmd
}

// Create a new sunodo tester.
// It is necessary to Close the tester at the end of the test.
func NewIntegrationTester(t *testing.T) *IntegrationTester {
	tester := &IntegrationTester{T: t}
	if tester.imageExist() {
		t.Log("image already exist; not building it again")
	} else {
		tester.sunodoBuild()
	}
	tester.sunodoRun()
	integrationMutex.Lock()
	return tester
}

// Close the tester.
func (t *IntegrationTester) Close() error {
	// Sending sigint directly to sunodo doesn't work.
	// So, we get the PID of the docker-compose child process and kill it.
	output, err := exec.Command("ps", "-o", "pid", "-C", "docker-compose").Output()
	if err != nil {
		t.Fatalf("failed to run ps: %v", err)
	}

	fields := strings.Fields(string(output))
	if len(fields) < 2 {
		t.Fatalf("failed to get docker-compose pid: %v", err)
	}

	_, err = exec.Command("kill", "-2", fields[1]).Output()
	if err != nil {
		t.Fatalf("failed to kill docker-compose")
	}

	t.cmd.Wait()
	integrationMutex.Unlock()
	return nil
}

func (t *IntegrationTester) imageExist() bool {
	path := "./.sunodo/image/hash"
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	t.Fatalf("unexpected error: %v", err)
	return false
}

func (t *IntegrationTester) sunodoBuild() {
	t.Log("running sunodo build")

	output, err := exec.Command("sunodo", "build").CombinedOutput()
	if err != nil {
		t.Logf("failed to run sunodo build: %v", err)
		t.Fatalf(string(output))
	}
}

func (t *IntegrationTester) sunodoRun() {
	t.Log("starting sunodo run")

	t.cmd = exec.Command("sunodo", "run")

	outPipe, err := t.cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	errPipe, err := t.cmd.StderrPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(errPipe)
		for scanner.Scan() {
			line := scanner.Text()
			t.Log(line)
		}
	}()

	ready := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(outPipe)
		for scanner.Scan() {
			line := scanner.Text()
			t.Log(line)
			if strings.Contains(line, "Press Ctrl+C to stop the node") {
				ready <- struct{}{}
			}
		}
	}()

	if err := t.cmd.Start(); err != nil {
		t.Fatalf("failed to start command: %v", err)
	}

	<-ready
}
