// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

// Status when finishing a rollups request.
type finishStatus int

const (
	statusAccept finishStatus = iota
	statusReject
)

func (status finishStatus) String() string {
	switch status {
	case statusAccept:
		return "accept"
	case statusReject:
		return "reject"
	}
	panic("invalid status")
}

// Interface with the Rollups backend API.
type rollupsApi interface {

	// Send a voucher to the Rollups API.
	sendVoucher(destination Address, payload []byte) error

	// Send a notice to the Rollups API.
	sendNotice(payload []byte) error

	// Send a report to the Rollups API.
	sendReport(payload []byte) error

	// Send a finish request to the Rollups API.
	// Return the advance payload and the metadata.
	finish(status finishStatus) ([]byte, *Metadata, error)
}
