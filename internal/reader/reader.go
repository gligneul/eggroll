// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package reader

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Input struct {
	// Advance input index.
	Index int

	// Completion status of the Input.
	Status CompletionStatus

	// Payload of the input.
	Payload []byte

	// Input sender.
	Sender common.Address

	// Number of the block the input was mined.
	BlockNumber int64

	// Time of the block the input was mined.
	BlockTimestamp time.Time

	// Resulting vouchers of the input.
	Vouchers []Voucher

	// Resulting notices of the input.
	Notices []Notice

	// Resulting reports of the input.
	Reports []Report
}

type Voucher struct {
	InputIndex  int
	OutputIndex int
	Destination common.Address
	Payload     []byte
}

type Notice struct {
	InputIndex  int
	OutputIndex int
	Payload     []byte
}

type Report struct {
	InputIndex  int
	OutputIndex int
	Payload     []byte
}

type InspectResult struct {
	Status              CompletionStatus
	Reports             []Report
	Exception           []byte
	ProcessedInputCount int
}

// Error when an object is not found.
type NotFound struct {
	typeName string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("%v not found", e.typeName)
}
