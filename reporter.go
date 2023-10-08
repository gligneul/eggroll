// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"
	"os"
)

// The first byte of a report has a tag to identify its semantic.
type reportTag byte

const (
	reportTagLog reportTag = iota
	reportTagResult
	reportTagLen
)

func encodeReport(tag reportTag, payload []byte) ([]byte, error) {
	if len(payload) >= 1000<<10 { // 1000 Kb
		return nil, fmt.Errorf("payload too large (%v bytes)", len(payload))
	}
	return append([]byte{byte(tag)}, payload...), nil
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

// Interface with the Rollups API.
type reporterRollupsAPI interface {
	SendReport(payload []byte) error
}

// Manages the reports of the DApp contract.
type Reporter struct {
	rollups reporterRollupsAPI
	logger  *log.Logger
}

func newReporter(rollups reporterRollupsAPI) *Reporter {
	return &Reporter{
		rollups: rollups,
		logger:  log.New(os.Stdout, "", 0),
	}
}

// Call fmt.Sprintln, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp client.
func (e *Reporter) Logln(a ...any) {
	e.Log(fmt.Sprintln(a...))
}

// Call fmt.Sprintf, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp client.
func (e *Reporter) Logf(format string, a ...any) {
	e.Log(fmt.Sprintf(format, a...))
}

// Call fmt.Sprint, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp client.
func (e *Reporter) Log(a ...any) {
	e.log(fmt.Sprint(a...))
}

// Call fmt.Sprintln, print the log, try to store the result in the rollups state, and exit.
// It is possible to retrieve this log in the DApp client.
func (e *Reporter) Fatalln(a ...any) {
	e.fatal(fmt.Sprintln(a...))
}

// Call fmt.Sprintf, print the log, try to store the result in the rollups state, and exit.
// It is possible to retrieve this log in the DApp client.
func (e *Reporter) Fatalf(format string, a ...any) {
	e.fatal(fmt.Sprintf(format, a...))
}

// Call fmt.Sprint, print the log, try to store the result in the rollups state, and exit.
// It is possible to retrieve this log in the DApp client.
func (e *Reporter) Fatal(a ...any) {
	e.fatal(fmt.Sprint(a...))
}

// Send a result advance or inspect result.
func (e *Reporter) sendResult(payload []byte) {
	e.sendReport(reportTagResult, payload)
}

// Send a report.
func (e *Reporter) sendReport(tag reportTag, payload []byte) {
	payload, err := encodeReport(tag, payload)
	if err != nil {
		e.logger.Fatalf("failed to encode report: %v", err)
	}
	if err := e.rollups.SendReport(payload); err != nil {
		e.logger.Fatalf("failed to send report: %v\n", err)
	}
}

// Log the message and send a report.
func (e *Reporter) log(message string) {
	e.logger.Print(message)
	e.sendReport(reportTagLog, []byte(message))
}

// Log the message, try to send a report, and exit.
func (e *Reporter) fatal(message string) {
	e.log(message)
	os.Exit(1)
}
