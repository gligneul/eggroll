// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import "context"

// Rollups input from the Reader API.
type Input struct {
	Index       int
	Status      CompletionStatus
	BlockNumber int64
}

// Rollups notice from the Reader API.
type Notice struct {
	InputIndex  int
	NoticeIndex int
	Payload     []byte
}

// Read the rollups state from the outside.
type Reader interface {

	// Get an input.
	Input(ctx context.Context, index int) (*Input, error)

	// Get a notice notice.
	Notice(ctx context.Context, inputIndex int, noticeIndex int) (*Notice, error)
}
