// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package dev

import (
	"context"
	"fmt"
	"os"

	"github.com/gligneul/eggroll/internal/sunodo"
)

func startSunodo(ctx context.Context, ready chan<- struct{}) error {
	isRunning, err := sunodo.IsRunning()
	if err != nil {
		return err
	}
	if isRunning {
		return fmt.Errorf("sunodo is already running")
	}

	sunodoStderr, err := os.Create("sunodo-stderr.log")
	if err != nil {
		return err
	}
	defer sunodoStderr.Close()

	logger.Print("sunodo: starting...")
	sunodoOpts := sunodo.RunOpts{
		NoBackend: true,
		Stdout:    logWritter{"sunodo"},
		Stderr:    sunodoStderr,
		Ready:     ready,
	}
	return sunodo.Run(ctx, sunodoOpts)
}
