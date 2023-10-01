// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"fmt"
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

func main() {
	var err error
	ctx := context.Background()
	client := eggroll.NewClient[State]()

	// indices, err := client.Send(
	// 	&InputClear{},
	// 	&InputAppend{Value: "egg"},
	// 	&InputAppend{Value: "roll"},
	// )
	// if err != nil {
	// 	log.Panic(err)
	// }
	//
	// lastInput := indices[len(indices)-1]

	lastInput := 2
	if err = client.WaitFor(ctx, lastInput); err != nil {
		log.Fatal(err)
	}

	var state *State
	if state, err = client.State(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println(state.TextBox) // -> eggroll
}
