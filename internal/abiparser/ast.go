// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package abiparser

type Ast struct {
	// Structs will be compiled to ABI tuples.
	Structs []Struct

	// Messages will be compiled to ABI functions.
	// Messages have the same structure as Structs, so we use the same
	// underlying type.
	Messages []Struct
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type any
}

type TypeName struct {
	Name string
}

type TypeArray struct {
	Elem any
}
