// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package dev

import (
	"context"
	"sync"
)

// Start the service, blocking util it ends.
// The service should exit when the context is done.
type service func(ctx context.Context, ready chan<- struct{}) error

// Start each service in one gorouting, waiting until it is ready to start the next one.
// When a service exits, send a cancel signal to all of them and wait them to finish.
func startServices(ctx context.Context, services []service) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, s := range services {
		wg.Add(1)
		ready := make(chan struct{})
		go func() {
			defer wg.Done()
			err := s(ctx, ready)
			if err != nil {
				logger.Printf("%v", err)
			}
			cancel()
		}()
		select {
		case <-ready:
		case <-ctx.Done():
			break
		}
	}

	<-ctx.Done()
	wg.Wait()
}
