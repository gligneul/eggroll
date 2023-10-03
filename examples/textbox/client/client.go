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

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Must[T any](obj T, err error) T {
	Check(err)
	return obj
}

func main() {
	ctx := context.Background()
	client := eggroll.NewClient[State]()

	printInput(&InputClear{})
	printInput(&InputAppend{Value: "egg"})
	printInput(&InputAppend{Value: "roll"})

	Check(client.WaitFor(ctx, 3))

	state := Must(client.State(ctx))
	log.Printf("Text box: '%v'\n", state.TextBox)

	logs := Must(client.Logs(ctx))
	log.Println("Logs:")
	for _, msg := range logs {
		log.Print(">", msg)
	}
}
