// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package abiparser

import "fmt"

var basicTypes = map[string]any{
	"bool":    TypeBool{},
	"int":     TypeInt{true, 256},
	"uint":    TypeInt{false, 256},
	"address": TypeAddress{},
	"bytes":   TypeBytes{},
	"string":  TypeString{},
}

func init() {
	for i := 8; i <= 256; i += 8 {
		basicTypes[fmt.Sprintf("int%v", i)] = TypeInt{true, i}
		basicTypes[fmt.Sprintf("uint%v", i)] = TypeInt{false, i}
	}
}

func analyze(ast Ast) (Ast, error) {
	structs := map[string]int{}
	for i, struct_ := range ast.Structs {
		_, ok := structs[struct_.Name]
		if ok {
			return Ast{}, fmt.Errorf("duplicate struct: %q", struct_.Name)
		}
		if len(struct_.Fields) == 0 {
			return Ast{}, fmt.Errorf("struct %v: must have fields", struct_.Name)
		}
		err := analyzeFields(structs, struct_.Fields)
		if err != nil {
			return Ast{}, fmt.Errorf("struct %v: %v", struct_.Name, err)
		}
		structs[struct_.Name] = i
	}
	messages := map[string]bool{}
	if len(ast.Messages) == 0 {
		return Ast{}, fmt.Errorf("no messages")
	}
	for _, message := range ast.Messages {
		_, ok := structs[message.Name]
		if ok {
			return Ast{}, fmt.Errorf("message with struct name: %q", message.Name)
		}
		ok = messages[message.Name]
		if ok {
			return Ast{}, fmt.Errorf("duplicate message: %q", message.Name)
		}
		err := analyzeFields(structs, message.Fields)
		if err != nil {
			return Ast{}, fmt.Errorf("message %v: %v", message.Name, err)
		}
		messages[message.Name] = true
	}
	return ast, nil
}

func analyzeFields(structs map[string]int, fields []Field) error {
	for i, field := range fields {
		type_, err := analyzeType(structs, field.Type)
		if err != nil {
			return fmt.Errorf("field %v: %v", field.Name, err)
		}
		// Update slice instead of the local copy
		fields[i].Type = type_
	}
	return nil
}

func analyzeType(structs map[string]int, type_ any) (any, error) {
	switch type_ := type_.(type) {
	case TypeName:
		basicType, ok := basicTypes[type_.Name]
		if ok {
			return basicType, nil
		}
		structIndex, ok := structs[type_.Name]
		if ok {
			return TypeStructRef{Index: structIndex}, nil
		}
		return nil, fmt.Errorf("type not found %q", type_.Name)
	case TypeArray:
		elemType, err := analyzeType(structs, type_.Elem)
		if err != nil {
			return nil, err
		}
		return TypeArray{Elem: elemType}, nil
	default:
		// This should not happen
		panic(fmt.Errorf("invalid type: %T", type_))
	}
}
