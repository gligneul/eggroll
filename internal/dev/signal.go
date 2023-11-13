// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package dev

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func startSignalListener(ctx context.Context, ready chan<- struct{}) error {
	logger.Print("Press Ctrl+C to exit")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ready <- struct{}{}
	sig := <-sigs
	return fmt.Errorf("received signal: %v", sig)
}
