// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"log"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/examples/textbox"
)

// Redefine the types to make the example cleaner
type (
	InputAppend textbox.InputAppend
	InputClear  textbox.InputClear
	State       textbox.State
)

func printInput(input any) {
	bytes, err := eggroll.EncodeInput(input)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(eggroll.EncodeHex(bytes))
}

func main() {
	var err error
	ctx := context.Background()
	client := eggroll.NewClient[State]()

	printInput(&InputClear{})
	printInput(&InputAppend{Value: "egg"})
	printInput(&InputAppend{Value: "roll"})

	if err = client.WaitFor(ctx, 2); err != nil {
		log.Fatal(err)
	}

	var state *State
	if state, err = client.State(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println(state.TextBox) // -> eggroll
}
