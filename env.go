// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"
)

type rollupsEnv struct {
	rollups  rollupsApi
	metadata *Metadata
}

// rollupsEnv implements the Env interface
var _ Env = (*rollupsEnv)(nil)

func (e *rollupsEnv) Metadata() *Metadata {
	return e.metadata
}

func (e *rollupsEnv) Report(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	log.Println(message)
	if err := e.rollups.sendReport([]byte(message)); err != nil {
		log.Fatalf("failed to send report: %v\n", err)
	}
}

func (e *rollupsEnv) Voucher(destination Address, payload []byte) {
	if err := e.rollups.sendVoucher(destination, payload); err != nil {
		log.Fatalf("failed to send voucher: %v\n", err)
	}
}
