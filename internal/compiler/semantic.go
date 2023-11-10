// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import "fmt"

// Perform the semantic analysis of the AST.
// This function also updates the struct references in the types.
func analyze(input []byte) (astSchema, error) {
	ast, err := parse(input)
	if err != nil {
		return ast, err
	}

	structToIndex := map[string]int{}
	if err := analyzeStructs(ast.Structs, structToIndex); err != nil {
		return ast, fmt.Errorf("struct %v", err)
	}

	// Create a set to avoid naming conflicts between messages
	messageSet := map[string]bool{}
	for name := range structToIndex {
		messageSet[name] = true
	}

	if err := analyzeMessages(ast.Reports, messageSet, structToIndex); err != nil {
		return ast, fmt.Errorf("report %v", err)
	}
	if err := analyzeMessages(ast.Advances, messageSet, structToIndex); err != nil {
		return ast, fmt.Errorf("advance %v", err)
	}
	if err := analyzeMessages(ast.Inspects, messageSet, structToIndex); err != nil {
		return ast, fmt.Errorf("inspect %v", err)
	}
	return ast, nil
}

// Check for duplicates and analyze fields.
// Structs are a special kind of schema because they can be referenced as types.
// Also, structs must have at least one field.
func analyzeStructs(structs []messageSchema, structToIndex map[string]int) error {
	for i, struct_ := range structs {
		_, ok := structToIndex[struct_.Name]
		if ok {
			return fmt.Errorf("duplicate of %q", struct_.Name)
		}
		if len(struct_.Fields) == 0 {
			return fmt.Errorf("%v: must have fields", struct_.Name)
		}
		err := analyzeFields(struct_.Fields, structToIndex)
		if err != nil {
			return fmt.Errorf("%v: %v", struct_.Name, err)
		}
		structToIndex[struct_.Name] = i
	}
	return nil
}

// Check for duplicates and analyze fields.
func analyzeMessages(
	messages []messageSchema,
	messageSet map[string]bool,
	structToIndex map[string]int,
) error {
	for _, message := range messages {
		_, ok := messageSet[message.Name]
		if ok {
			return fmt.Errorf("duplicate of %q", message.Name)
		}
		err := analyzeFields(message.Fields, structToIndex)
		if err != nil {
			return fmt.Errorf("%v: %v", message.Name, err)
		}
		messageSet[message.Name] = true
	}
	return nil
}

// Analyze the type of each field.
func analyzeFields(fields []fieldSchema, structToIndex map[string]int) error {
	for i, field := range fields {
		type_, err := analyzeType(field.type_, structToIndex)
		if err != nil {
			return fmt.Errorf("field %v: %v", field.Name, err)
		}
		// Make the change directly to the slice, otherwise it
		// will be lost because field is a local copy.
		fields[i].type_ = type_
	}
	return nil
}

// Recursively analyze the type, filling up the struct references.
func analyzeType(type_ any, structToIndex map[string]int) (any, error) {
	switch type_ := type_.(type) {
	case typeArray:
		var err error
		type_.Elem, err = analyzeType(type_.Elem, structToIndex)
		if err != nil {
			return nil, err
		}
		return type_, nil
	case typeStructRef:
		var ok bool
		type_.Index, ok = structToIndex[type_.Name]
		if !ok {
			return nil, fmt.Errorf("struct %q not found", type_.Name)
		}
		return type_, nil
	default:
		// basic types are already correct
		return type_, nil
	}
}
