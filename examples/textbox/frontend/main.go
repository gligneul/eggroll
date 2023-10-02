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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func must[T any](obj T, err error) T {
	check(err)
	return obj
}

func main() {
	ctx := context.Background()
	client := eggroll.NewClient[State]()

	printInput(&InputClear{})
	printInput(&InputAppend{Value: "egg"})
	printInput(&InputAppend{Value: "roll"})

	check(client.WaitFor(ctx, 2))

	state := must(client.State(ctx))
	log.Printf("Text box: '%v'\n", state.TextBox)

	logs := must(client.Logs(ctx))
	log.Println("Logs:")
	for _, msg := range logs {
		log.Print(">", msg)
	}
}
