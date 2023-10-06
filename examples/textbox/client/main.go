// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"context"
	"log"

	"textbox"

	"github.com/gligneul/eggroll"
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
		textbox.Clear{},
		textbox.Append{Value: "egg"},
		textbox.Append{Value: "roll"},
	}
	for _, input := range inputs {
		log.Printf("Sending input %#v\n", input)
		Check(client.SendGeneric(ctx, input))
	}

	log.Println("Waiting for inputs to be processed")
	Check(client.WaitFor(ctx, 2))

	Check(client.Sync(ctx))
	var contract textbox.Contract
	Check(client.ReadState(&contract))
	log.Printf("Text box: '%v'\n", contract.TextBox)

	logs := Must(client.Logs(ctx))
	log.Println("Logs:")
	for _, msg := range logs {
		log.Print(">", msg)
	}
}
