// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
)

const defaultRollupsEndpoint string = "http://localhost:8080/host-runner"

// Set up the DApp backend.
//
// Load the Rollups HTTP server endpoint from the ROLLUP_HTTP_SERVER_URL variable.
func SetupDApp[S any]() DApp[S] {
	endpoint := os.Getenv("ROLLUP_HTTP_SERVER_URL")
	if endpoint == "" {
		endpoint = defaultRollupsEndpoint
	}

	log.Printf("setting rollups http endpoint to %v", endpoint)

	rollups := &rollupsHttpApi{endpoint}
	dapp := dapp[S]{
		rollups:  rollups,
		handlers: make(handlerMap[S]),
	}

	return &dapp
}

type dapp[S any] struct {
	rollups  rollupsApi
	handlers handlerMap[S]
}

// dapp[S] implements the DApp[S] interface
var _ DApp[struct{}] = (*dapp[struct{}])(nil)

func (d *dapp[S]) register(inputType reflect.Type, handler internalHandler[S]) {
	key := computeHandlerKey(inputType)
	d.handlers[key] = handlerEntry[S]{inputType, handler}
}

func (d *dapp[S]) Roll() {
	var state S
	env := rollupsEnv{rollups: d.rollups}
	status := statusAccept

	for {
		payload, metadata, err := d.rollups.finish(status)
		if err != nil {
			log.Fatalf("failed to send finish: %v\n", err)
		}

		env.metadata = metadata

		if len(payload) < 4 {
			env.Report("invalid payload len")
			status = statusReject
			continue
		}

		key := handlerKey(payload[:4])
		entry, ok := d.handlers[key]
		if !ok {
			env.Report("handler not found")
			status = statusReject
			continue
		}

		input := reflect.New(entry.inputType.Elem())
		inputPointer := input.Interface()
		if err = json.Unmarshal(payload[4:], inputPointer); err != nil {
			env.Report("failed to decode input: %v", err)
			status = statusReject
			continue
		}

		err = entry.handler(&env, &state, inputPointer)
		if err != nil {
			env.Report("rejecting: %v", err)
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
