// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"log"

	"textbox"

	"github.com/gligneul/eggroll"
)

type (
	InputAppend textbox.InputAppend
	InputClear  textbox.InputClear
	State       textbox.State
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Must[T any](obj T, err error) T {
	Check(err)
	return obj
}

// @cut

func main() {
	ctx := context.Background()
	client := eggroll.NewClient()

	inputs := []any{
		InputClear{},
		InputAppend{Value: "egg"},
		InputAppend{Value: "roll"},
	}
	for _, input := range inputs {
		log.Printf("Sending input %#v\n", input)
		Check(client.SendGeneric(ctx, input))
	}

	log.Println("Waiting for inputs to be processed")
	Check(client.WaitFor(ctx, 2))

	Check(client.Sync(ctx))
	var state State
	Check(client.ReadState(&state))
	log.Printf("Text box: '%v'\n", state.TextBox)

	logs := Must(client.Logs(ctx))
	log.Println("Logs:")
	for _, msg := range logs {
		log.Print(">", msg)
	}
}
