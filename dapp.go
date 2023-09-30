// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
)

// Configuration for the DApp.
type DAppConfig struct {
	RollupsEndpoint string
}

// Load the config from environment variables.
func (c *DAppConfig) Load() {
	varName := "EGGROLL_ROLLUPS_HTTP_ENDPOINT"
	c.RollupsEndpoint = os.Getenv(varName)
	if c.RollupsEndpoint == "" {
		if runtime.GOARCH == "riscv64" {
			c.RollupsEndpoint = "http://localhost:5004"
		} else {
			c.RollupsEndpoint = "http://localhost:8080/host-runner"
		}
	}
	log.Printf("set %v=%v\n", varName, c.RollupsEndpoint)
}

// DApp is the back-end for the rollups.
// It dispatch the input to the corresponding handler while it advances the
// rollups state.
type DApp[S any] struct {
	rollups  rollupsApi
	handlers handlerMap[S]
}

// Create the DApp loading the config from environment variables.
func NewDApp[S any]() *DApp[S] {
	var config DAppConfig
	config.Load()
	return NewDAppFromConfig[S](config)
}

// Create the DApp with a custom config.
func NewDAppFromConfig[S any](config DAppConfig) *DApp[S] {
	rollups := &rollupsHttpApi{config.RollupsEndpoint}
	dapp := DApp[S]{
		rollups:  rollups,
		handlers: make(handlerMap[S]),
	}
	return &dapp
}

// Start the DApp back end.
// This function never returns and exits if there is an error.
func (d *DApp[S]) Roll() {
	var state S
	env := &Env{rollups: d.rollups}
	status := statusAccept

	for {
		var (
			payload []byte
			err     error
		)
		payload, env.metadata, err = d.rollups.finish(status)
		if err != nil {
			log.Fatalf("failed to send finish: %v\n", err)
		}

		if err = d.handlers.dispatch(env, &state, payload); err != nil {
			env.Report(err.Error())
			status = statusReject
			continue
		}

		stateSnapshot, err := json.Marshal(&state)
		if err != nil {
			log.Fatalf("failed to create state snapshot: %v\n", err)
		}
		if err = d.rollups.sendNotice(stateSnapshot); err != nil {
			log.Fatalf("failed to send notice: %v\n", err)
		}
		status = statusAccept
	}
}
