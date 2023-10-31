// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestCompile(t *testing.T) {
	input := `---
structs:
  - Chapter:
    - title: string

  - Book:
    - id: int256
    - title: string
    - author: string
    - chapters: Chapter[]

messages:
  - AddBook:
    - book: Book
`
	reader := strings.NewReader(input)
	json, err := Compile(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	abi, err := abi.JSON(strings.NewReader(json))
	if err != nil {
		panic(fmt.Errorf("failed to generate abi: %v", err))
	}
	method, ok := abi.Methods["AddBook"]
	if !ok {
		t.Fatalf("method not found")
	}
	if method.String() != `function AddBook((int256,string,string,()[]) book) returns()` {
		t.Fatalf("wrong method: %v", method.String())
	}
}
