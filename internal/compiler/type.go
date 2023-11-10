// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import "fmt"

var basicTypes = map[string]any{
	"bool":    typeBool{},
	"int":     typeInt{true, 256},
	"uint":    typeInt{false, 256},
	"address": typeAddress{},
	"bytes":   typeBytes{},
	"string":  typeString{},
}

func init() {
	for i := 8; i <= 256; i += 8 {
		basicTypes[fmt.Sprintf("int%v", i)] = typeInt{true, i}
		basicTypes[fmt.Sprintf("uint%v", i)] = typeInt{false, i}
	}
}

type typeBool struct{}

type typeInt struct {
	Signed bool
	Bits   int
}

type typeAddress struct{}

type typeBytes struct{}

type typeString struct{}

type typeArray struct {
	Elem any
}

type typeStructRef struct {
	Name  string
	Index int
}
