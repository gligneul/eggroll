// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/json"
	"log"
	"reflect"
)

type dapp[S any] struct {
	rollups  rollupsApi
	handlers handlerMap[S]
}

func (d *dapp[S]) register(inputType reflect.Type, handler internalHandler[S]) {
	d.handlers.register(inputType, handler)
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

		if err = d.handlers.dispatch(&env, &state, payload); err != nil {
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
