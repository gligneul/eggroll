// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

type (
	// Ethereum Address.
	Address common.Address

	// Ethereum Hash.
	Hash common.Hash
)

// Metadata from the input.
type Metadata struct {
	Sender         Address
	BlockNumber    int64
	BlockTimestamp int64
}

// Env allows the DApp backend to interact with the Rollups API.
type Env struct {
	rollups  rollupsApi
	metadata *Metadata
}

// Get the Metadata for the current input.
func (e *Env) Metadata() *Metadata {
	return e.metadata
}

// Send a report for debugging purposes.
// The reports will be available as logs for the front-end client.
// It is not necessary to add a `\n` at the end of the report.
func (e *Env) Report(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	log.Println(message)
	if err := e.rollups.sendReport([]byte(message)); err != nil {
		log.Fatalf("failed to send report: %v\n", err)
	}
}

// Send a voucher.
func (e *Env) Voucher(destination Address, payload []byte) {
	if err := e.rollups.sendVoucher(destination, payload); err != nil {
		log.Fatalf("failed to send voucher: %v\n", err)
	}
}
