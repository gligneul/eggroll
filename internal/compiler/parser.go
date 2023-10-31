// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// In the yaml file, fields are represented as a map with a single value.
type yamlField map[string]string

const yamlFieldFmt = `
<fieldName>: <fieldType>`

// In the yaml file, structs are represented as a map with a single value.
type yamlStruct map[string][]yamlField

const yamlStructFmt = `
<structName>:
  - <field1>
  - <field2>
  - <fieldN>`

func parse(input []byte) (Ast, error) {
	var yamlAst struct {
		Structs  []yamlStruct
		Messages []yamlStruct
	}
	err := yaml.Unmarshal(input, &yamlAst)
	if err != nil {
		return Ast{}, err
	}
	var ast Ast
	ast.Structs, err = parseStructs("struct", yamlAst.Structs)
	if err != nil {
		return Ast{}, err
	}
	ast.Messages, err = parseStructs("message", yamlAst.Messages)
	if err != nil {
		return Ast{}, err
	}
	return ast, nil
}

func parseStructs(kind string, yamlStructs []yamlStruct) ([]Struct, error) {
	var structs []Struct
	for _, yamlStruct := range yamlStructs {
		struct_, err := parseStruct(kind, yamlStruct)
		if err != nil {
			return nil, err
		}
		structs = append(structs, struct_)
	}
	return structs, nil
}

func parseStruct(kind string, yamlStruct yamlStruct) (Struct, error) {
	if len(yamlStruct) != 1 {
		return Struct{}, fmt.Errorf("invalid %v syntax, expected: %v", kind, yamlStructFmt)
	}
	var name string
	for key := range yamlStruct {
		// This works because we checked the struct has only one key
		name = key
	}
	if err := checkIdentifier(name); err != nil {
		return Struct{}, fmt.Errorf("%v name: %v", kind, err)
	}
	var struct_ Struct
	struct_.Name = name
	for _, yamlField := range yamlStruct[name] {
		field, err := parseField(yamlField)
		if err != nil {
			msg := "%v %v: %v"
			return Struct{}, fmt.Errorf(msg, kind, name, err)
		}
		struct_.Fields = append(struct_.Fields, field)
	}
	return struct_, nil
}

func parseField(yamlField yamlField) (Field, error) {
	if len(yamlField) != 1 {
		return Field{}, fmt.Errorf("invalid field syntax, expected: %v", yamlFieldFmt)
	}
	var name string
	for key := range yamlField {
		// This works because we checked the struct has only one key
		name = key
	}
	if err := checkIdentifier(name); err != nil {
		return Field{}, fmt.Errorf("field name: %v", err)
	}
	type_, err := parseType(yamlField[name])
	if err != nil {
		return Field{}, fmt.Errorf("field %v: %v", name, err)
	}
	var field Field
	field.Name = name
	field.Type = type_
	return field, nil
}

func parseType(rawType string) (any, error) {
	id, isArray, err := tokenizeType(rawType)
	if err != nil {
		return nil, err
	}
	var t any = TypeName{Name: id}
	if isArray {
		t = TypeArray{Elem: t}
	}
	return t, nil
}
