// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

// Top level struct for the EggRoll schema.
type astSchema struct {

	// Plain structs that will be compiled into.
	Structs []messageSchema

	// Schemas that will be used as reports.
	Reports []messageSchema

	// Schemas that will be used as advance requests.
	Advances []messageSchema

	// Schemas that will be used as inspect requests.
	Inspects []messageSchema
}

// Schema for a message, that can be a plain struct, an input, or an output.
type messageSchema struct {
	Name   string
	Doc    string
	Fields []fieldSchema
}

// Schema for a field of a message.
type fieldSchema struct {
	Name string
	Doc  string

	// Raw type from the input schema; it might not be a valid type.
	Type string

	// Once the type is validated, this field is set to a type* struct.
	type_ any
}
