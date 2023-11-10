// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Parse the AST from a YAML file.
func parse(input []byte) (ast astSchema, err error) {
	if err = yaml.Unmarshal(input, &ast); err != nil {
		return ast, err
	}
	if err = parseMessages(ast.Structs); err != nil {
		return ast, fmt.Errorf("struct %v", err)
	}
	if err = parseMessages(ast.Reports); err != nil {
		return ast, fmt.Errorf("report %v", err)
	}
	if err = parseMessages(ast.Advances); err != nil {
		return ast, fmt.Errorf("advance %v", err)
	}
	if err = parseMessages(ast.Inspects); err != nil {
		return ast, fmt.Errorf("inspect %v", err)
	}
	return ast, nil
}

// Validate the message and field names, and the field types.
func parseMessages(messages []messageSchema) error {
	for _, message := range messages {
		if err := checkName(message.Name); err != nil {
			return fmt.Errorf("name: %v", err)
		}
		for i, field := range message.Fields {
			if err := checkName(field.Name); err != nil {
				return fmt.Errorf("%v: field name: %v", message.Name, err)
			}
			type_, err := parseType(field.Type)
			if err != nil {
				return fmt.Errorf("%v.%v type: %v", message.Name, field.Name, err)
			}
			// Make the change directly to the slice, otherwise it
			// will be lost because field is a local copy.
			message.Fields[i].type_ = type_
		}
	}
	return nil
}

func parseType(rawType string) (any, error) {
	typeName, isArray, err := tokenizeType(rawType)
	if err != nil {
		return nil, err
	}
	type_ := basicTypes[typeName]
	if type_ == nil {
		// Assume it is a struct reference if it isn't a basic type
		type_ = typeStructRef{
			Name: typeName,
		}
	}
	if isArray {
		type_ = typeArray{Elem: type_}
	}
	return type_, nil
}
