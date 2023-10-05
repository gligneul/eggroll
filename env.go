// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/rollups"
)

// Env allows the DApp contract to interact with the Rollups API.
type Env struct {
	rollups  rollupsAPI
	metadata *rollups.Metadata
}

// Get the Metadata for the current input.
func (e *Env) Metadata() *rollups.Metadata {
	return e.metadata
}

// Call fmt.Sprintln, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *Env) Logln(a any) {
	e.Log(fmt.Sprintln(a))
}

// Call fmt.Sprintf, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *Env) Logf(format string, a ...any) {
	e.Log(fmt.Sprintf(format, a...))
}

// Call fmt.Sprint, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *Env) Log(a any) {
	e.log(fmt.Sprint(a))
}

// Log the message and send a report.
func (e *Env) log(message string) {
	log.Print(message)
	if err := e.rollups.SendReport([]byte(message)); err != nil {
		log.Fatalf("failed to send report: %v\n", err)
	}
}

// Send a voucher.
func (e *Env) Voucher(destination common.Address, payload []byte) {
	if err := e.rollups.SendVoucher(destination, payload); err != nil {
		log.Fatalf("failed to send voucher: %v\n", err)
	}
}
