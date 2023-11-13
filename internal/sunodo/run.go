// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package sunodo

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Options to run sunodo
type RunOpts struct {

	// Run in no backend mode
	NoBackend bool

	// Stdout writer
	Stdout io.Writer

	// Stderr writer
	Stderr io.Writer

	// Send a signal when sunodo is ready to receive inputs
	Ready chan<- struct{}
}

// Execute the sunodo run command.
// This function blocks until sunodo exits.
func Run(ctx context.Context, opts RunOpts) error {
	args := []string{"run"}
	if opts.NoBackend {
		args = append(args, "--no-backend")
	}

	cmd := exec.CommandContext(ctx, "sunodo", args...)
	cmd.Stderr = opts.Stderr
	cmd.Cancel = func() error {
		return stopSunodo(cmd.Process.Pid)
	}

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(outPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if opts.Stdout != nil {
				fmt.Fprintln(opts.Stdout, line)
			}
			if strings.Contains(line, "Press Ctrl+C to stop the node") {
				opts.Ready <- struct{}{}
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}

	return cmd.Wait()
}

// Sending sigint directly to sunodo doesn't work.
// So, we get the PID of the docker-compose child process and kill it.
func stopSunodo(pid int) error {
	dockerPid, err := procGetChild(pid)
	if err != nil {
		return fmt.Errorf("failed to get docker pid: %v", err)
	}
	dockerComposePid, err := procGetChild(dockerPid)
	if err != nil {
		return fmt.Errorf("failed to get docker compose pid: %v", err)
	}
	err = procInterrupt(dockerComposePid)
	if err != nil {
		return fmt.Errorf("failed to send interrupt to docker compose: %v", err)
	}
	err = procWait(dockerComposePid, 30)
	if err != nil {
		return fmt.Errorf("failed to wait for docker compose: %v", err)
	}
	return nil
}
