// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import "sync"

type clientState[S any] struct {
	sync.RWMutex
	state S
}

type client[S any] struct {
	clientState[S]
}

func (c *client[S]) Send(input ...any) ([]int, error) {
	return nil, nil
}

func (c *client[S]) WaitFor(inputIndex int) error {
	return nil
}

func (c *client[S]) State() *S {
	return nil
}

func (c *client[S]) LogsFrom(input int) (string, error) {
	return "", nil
}

func (c *client[S]) LogsTail(n int) (string, error) {
	return "", nil
}

func (c *client[S]) Logs() (string, error) {
	return "", nil
}
