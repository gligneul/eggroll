// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggtypes

import (
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

// Get the raw return of the result.
func (r *Result) RawReturn() []byte {
	for _, report := range r.Reports {
		tag, payload, err := DecodeReport(report.Payload)
		if err == nil && tag == ReportTagReturn {
			return payload
		}
	}
	return nil
}

// Get the logs of the result.
func (r *Result) Logs() []string {
	var logs []string
	for _, report := range r.Reports {
		tag, payload, err := DecodeReport(report.Payload)
		if err == nil && tag == ReportTagLog {
			logs = append(logs, string(payload))
		}
	}
	return logs
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

// The first byte of a report has a tag to identify its semantic.
type ReportTag byte

const (
	ReportTagLog ReportTag = iota
	ReportTagReturn
	ReportTagLen
)

// Encode the tag and the report into a single payload.
func EncodeReport(tag ReportTag, payload []byte) ([]byte, error) {
	if len(payload) >= 1000<<10 { // 1000 Kb
		return nil, fmt.Errorf("payload too large (%v bytes)", len(payload))
	}
	return append([]byte{byte(tag)}, payload...), nil
}

// Decode the report into tag and payload.
func DecodeReport(payload []byte) (ReportTag, []byte, error) {
	if len(payload) == 0 {
		return 0, payload, fmt.Errorf("invalid report")
	}
	tag := payload[0]
	if tag >= byte(ReportTagLen) {
		return 0, payload, fmt.Errorf("invalid report tag %v", tag)
	}
	return ReportTag(tag), payload[1:], nil
}
