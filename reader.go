// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/json"
	"fmt"
)

// Read the rollups state from the outside.
type Reader interface {

	// Get a completion status for the given input.
	Input(index int) (*Input, error)
}

// Status of a rollups advance.
type CompletionStatus int

const (
	CompletionUnprocessed CompletionStatus = iota
	CompletionAccepted
	CompletionRejected
	CompletionException
	CompletionMachineHalted
	CompletionCycleLimitExceeded
	CompletionTimeLimitExceeded
	CompletionPayloadLengthLimitExceeded
)

func (c CompletionStatus) MarshalJSON() ([]byte, error) {
	toString := map[CompletionStatus]string{
		CompletionUnprocessed:                "UNPROCESSED",
		CompletionAccepted:                   "ACCEPTED",
		CompletionRejected:                   "REJECTED",
		CompletionException:                  "EXCEPTION",
		CompletionMachineHalted:              "MACHINE_HALTED",
		CompletionCycleLimitExceeded:         "CYCLE_LIMIT_EXCEEDED",
		CompletionTimeLimitExceeded:          "TIME_LIMIT_EXCEEDED",
		CompletionPayloadLengthLimitExceeded: "PAYLOAD_LENGTH_LIMIT_EXCEEDED",
	}
	return json.Marshal(toString[c])
}

func (c *CompletionStatus) UnmarshalJSON(data []byte) error {
	toCompletion := map[string]CompletionStatus{
		"UNPROCESSED":                   CompletionUnprocessed,
		"ACCEPTED":                      CompletionAccepted,
		"REJECTED":                      CompletionRejected,
		"EXCEPTION":                     CompletionException,
		"MACHINE_HALTED":                CompletionMachineHalted,
		"CYCLE_LIMIT_EXCEEDED":          CompletionCycleLimitExceeded,
		"TIME_LIMIT_EXCEEDED":           CompletionTimeLimitExceeded,
		"PAYLOAD_LENGTH_LIMIT_EXCEEDED": CompletionPayloadLengthLimitExceeded,
	}

	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	var ok bool
	if *c, ok = toCompletion[dataStr]; !ok {
		return fmt.Errorf("invalid completion %v", dataStr)
	}

	return nil
}

// Rollups input from the Reader API.
type Input struct {
	Index       int
	Status      CompletionStatus
	BlockNumber string
}
