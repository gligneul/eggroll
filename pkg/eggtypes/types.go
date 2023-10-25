// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggtypes

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Completion of an advance or inspect request.
type CompletionStatus int

const (
	CompletionStatusUnprocessed CompletionStatus = iota
	CompletionStatusAccepted
	CompletionStatusRejected
	CompletionStatusException
	CompletionStatusMachineHalted
	CompletionStatusCycleLimitExceeded
	CompletionStatusTimeLimitExceeded
	CompletionStatusPayloadLengthLimitExceeded
)

// Result of an request.
type Result struct {

	// Completion status of the request.
	Status CompletionStatus

	// Resulting reports of the request.
	Reports []Report
}

// Get logs from the result.
func (r *Result) Logs() []Log {
	return FilterReports[Log](r.Reports, LogID)
}

// Result of an advance request.
type AdvanceResult struct {
	Result

	// Advance input index.
	Index int

	// Payload of the input.
	Payload []byte

	// Input sender.
	Sender common.Address

	// Number of the block the input was mined.
	BlockNumber int64

	// Time of the block the input was mined.
	BlockTimestamp time.Time

	// Resulting vouchers of the request.
	Vouchers []Voucher

	// Resulting notices of the request.
	Notices []Notice
}

// Result of an inspect request.
type InspectResult struct {
	Result

	// Number of processed advance inputs when the inspect was made.
	ProcessedInputCount int
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

// Filter the reports with the given id and unpack it into T.
func FilterReports[T any](reports []Report, id [4]byte) []T {
	var values []T
	for _, r := range reports {
		if bytes.HasPrefix(r.Payload, id[:]) {
			v, err := Unpack(r.Payload)
			if err != nil {
				// This should never happen because the callee
				// requested for an specific id.
				panic(fmt.Errorf("error unpacking: %v", err))
			}
			values = append(values, v.(T))
		}
	}
	return values
}

// Find the report with the given id and unpack it into T.
func FindReport[T any](reports []Report, id [4]byte) (empty T, found bool) {
	for _, r := range reports {
		if bytes.HasPrefix(r.Payload, id[:]) {
			v, err := Unpack(r.Payload)
			if err != nil {
				// This should never happen because the callee
				// requested for an specific id.
				panic(fmt.Errorf("error unpacking: %v", err))
			}
			return v.(T), true
		}
	}
	return empty, false
}
