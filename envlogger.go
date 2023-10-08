// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"fmt"
	"log"
	"os"
)

// Interface with the Rollups API.
type envLoggerRollupsAPI interface {
	SendReport(payload []byte) error
}

// Manages the log of the DApp contract.
type EnvLogger struct {
	rollups envLoggerRollupsAPI
	logger  *log.Logger
}

func newEnvLogger(rollups envLoggerRollupsAPI) *EnvLogger {
	return &EnvLogger{
		rollups: rollups,
		logger:  log.New(os.Stdout, "", 0),
	}
}

// Call fmt.Sprintln, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *EnvLogger) Logln(a ...any) {
	e.Log(fmt.Sprintln(a...))
}

// Call fmt.Sprintf, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *EnvLogger) Logf(format string, a ...any) {
	e.Log(fmt.Sprintf(format, a...))
}

// Call fmt.Sprint, print the log, and store the result in the rollups state.
// It is possible to retrieve this log in the DApp front end.
func (e *EnvLogger) Log(a ...any) {
	e.log(fmt.Sprint(a...))
}

// Log the message and send a report.
func (e *EnvLogger) log(message string) {
	e.logger.Print(message)
	report := encodeLogReport([]byte(message))
	if err := e.rollups.SendReport(report); err != nil {
		e.logger.Fatalf("failed to send report: %v\n", err)
	}
}

// Call fmt.Sprintln, print the log, try to store the result in the rollups state, and exit.
// It is possible to retrieve this log in the DApp front end.
func (e *EnvLogger) Fatalln(a ...any) {
	e.fatal(fmt.Sprintln(a...))
}

// Call fmt.Sprintf, print the log, try to store the result in the rollups state, and exit.
// It is possible to retrieve this log in the DApp front end.
func (e *EnvLogger) Fatalf(format string, a ...any) {
	e.fatal(fmt.Sprintf(format, a...))
}

// Call fmt.Sprint, print the log, try to store the result in the rollups state, and exit.
// It is possible to retrieve this log in the DApp front end.
func (e *EnvLogger) Fatal(a ...any) {
	e.fatal(fmt.Sprint(a...))
}

// Log the message, try to send a report, and exit.
func (e *EnvLogger) fatal(message string) {
	e.log(message)
	os.Exit(1)
}
