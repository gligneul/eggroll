// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package dev

import (
	"context"
	"os"
	"os/exec"
	"time"
)

// Run the dapp every time it receives a signal.
type runner struct {
	path string
	run  <-chan struct{}
}

func (r *runner) start(ctx context.Context, ready chan<- struct{}) error {
	ready <- struct{}{}
	for {
		dappCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			select {
			case <-dappCtx.Done():
				return
			case <-r.run:
				cancel()
			}
		}()

		logger.Print("dapp: starting...")
		cmd := exec.CommandContext(dappCtx, "go", "run", r.path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		logger.Print("dapp: done")

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(3 * time.Second):
		}
	}
}
