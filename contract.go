// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// A high-level, opinionated, lambda-based framework for Cartesi Rollups in Go.
package eggroll

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/rollups"
)

// Configuration for the Contract.
type ContractConfig struct {
	RollupsEndpoint string
}

// Load the config from environment variables.
func (c *ContractConfig) Load() {
	var defaultEndpoint string
	if runtime.GOARCH == "riscv64" {
		defaultEndpoint = "http://127.0.0.1:5004"
	} else {
		defaultEndpoint = "http://localhost:8080/host-runner"
	}
	c.RollupsEndpoint = loadVar("ROLLUPS_HTTP_ENDPOINT", defaultEndpoint)
}

// Interface with the Rollups API.
// We don't expose this API because calling it directly will break EggRoll assumptions.
type rollupsAPI interface {
	SendVoucher(destination common.Address, payload []byte) error
	SendNotice(payload []byte) error
	SendReport(payload []byte) error
	Finish(status rollups.FinishStatus) ([]byte, *rollups.Metadata, error)
}

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

// Contract is the back-end for the rollups.
// It dispatch the input to the corresponding handler while it advances the
// rollups state.
type Contract[S any] struct {
	rollups  rollupsAPI
	handlers handlerMap[S]
}

// Create the Contract loading the config from environment variables.
func NewContract[S any]() *Contract[S] {
	var config ContractConfig
	config.Load()
	return NewContractFromConfig[S](config)
}

// Create the Contract with a custom config.
func NewContractFromConfig[S any](config ContractConfig) *Contract[S] {
	rollups := rollups.NewRollupsHTTP(config.RollupsEndpoint)
	contract := &Contract[S]{
		rollups:  rollups,
		handlers: make(handlerMap[S]),
	}
	return contract
}

// Start the Contract back end.
// This function never returns and exits if there is an error.
func (d *Contract[S]) Roll() {
	var state S
	env := &Env{rollups: d.rollups}
	status := rollups.FinishStatusAccept

	for {
		var (
			payload []byte
			err     error
		)
		payload, env.metadata, err = d.rollups.Finish(status)
		if err != nil {
			log.Fatalf("failed to send finish: %v\n", err)
		}

		if err = d.handlers.dispatch(env, &state, payload); err != nil {
			env.Logf("rejecting input: %v\n", err)
			status = rollups.FinishStatusReject
			continue
		}

		stateSnapshot, err := json.Marshal(&state)
		if err != nil {
			log.Fatalf("failed to create state snapshot: %v\n", err)
		}
		if err = d.rollups.SendNotice(stateSnapshot); err != nil {
			log.Fatalf("failed to send notice: %v\n", err)
		}
		status = rollups.FinishStatusAccept
	}
}

// Signature of the handler that advances the rollups state.
type Handler[S, I any] func(*Env, *S, *I) error

// Register a handler for a custom input to a Contract.
func Register[S, I any](d *Contract[S], handler Handler[S, I]) {
	// This function needs to be defined outside of the Contract interface
	// because Go doesn't support template parameters in methods.
	// So, it is not possible to write Contract[S].Register[I](Handler[S, I]).

	key, decoder := getInputKeyDecoder[I]()

	gHandler := func(env *Env, state *S, inputBytes []byte) error {
		input, err := decoder(inputBytes)
		if err != nil {
			return err
		}
		return handler(env, state, input)
	}

	d.handlers[key] = gHandler
}

// Internal handler that receives an encoded input.
type genericHandler[S any] func(*Env, *S, []byte) error

// Map a input key to its handler.
type handlerMap[S any] map[inputKey]genericHandler[S]

// Call the handler for the given input.
func (m handlerMap[S]) dispatch(env *Env, state *S, payload []byte) error {
	inputKey, inputBytes, err := splitInput(payload)
	if err != nil {
		return err
	}

	handler, ok := m[inputKey]
	if !ok {
		return fmt.Errorf("handler not found (%v)", hex.EncodeToString(inputKey[:]))
	}

	if err := handler(env, state, inputBytes); err != nil {
		return fmt.Errorf("input rejected: %v", err)
	}

	return nil
}
