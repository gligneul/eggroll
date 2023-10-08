// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import "fmt"

// The first byte of a report has a tag to identify its semantic.
type reportTag byte

const (
	reportTagLog reportTag = iota
	reportTagResult
	reportTagLen
)

func encodeLogReport(payload []byte) []byte {
	return append([]byte{byte(reportTagLog)}, payload...)
}

func encodeResultReport(payload []byte) []byte {
	return append([]byte{byte(reportTagResult)}, payload...)
}

func decodeReport(payload []byte) (reportTag, []byte, error) {
	if len(payload) == 0 {
		return 0, payload, fmt.Errorf("invalid report")
	}
	tag := payload[0]
	if tag >= byte(reportTagLen) {
		return 0, payload, fmt.Errorf("invalid report tag %v", tag)
	}
	return reportTag(tag), payload[1:], nil
}
