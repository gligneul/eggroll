// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"encoding/json"
	"fmt"
)

type jsonAbiArg struct {
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	InternalType string       `json:"internalType"`
	Components   []jsonAbiArg `json:"components"`
}

type jsonAbiMethod struct {
	Name            string       `json:"name"`
	Type            string       `json:"type"`
	StateMutability string       `json:"stateMutability"`
	Inputs          []jsonAbiArg `json:"inputs"`
	Outputs         []jsonAbiArg `json:"outputs"`
}

func generate(ast Ast) string {
	var methods []jsonAbiMethod
	for _, message := range ast.Messages {
		var method jsonAbiMethod
		method.Name = message.Name
		method.Type = "function"
		method.StateMutability = "nonpayable"
		for _, field := range message.Fields {
			arg := generateArg(ast, field.Name, field.Type)
			method.Inputs = append(method.Inputs, arg)
		}
		methods = append(methods, method)
	}
	result, err := json.MarshalIndent(methods, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("json marshal error: %v", err))
	}
	return string(result)
}

func generateArg(ast Ast, name string, type_ any) jsonAbiArg {
	switch type_ := type_.(type) {
	case TypeBool:
		return jsonAbiArg{
			Name:         name,
			Type:         "bool",
			InternalType: "bool",
		}
	case TypeInt:
		prefix := ""
		if !type_.Signed {
			prefix = "u"
		}
		typeName := fmt.Sprintf("%vint%v", prefix, type_.Bits)
		return jsonAbiArg{
			Name:         name,
			Type:         typeName,
			InternalType: typeName,
		}
	case TypeAddress:
		return jsonAbiArg{
			Name:         name,
			Type:         "address",
			InternalType: "address",
		}
	case TypeBytes:
		return jsonAbiArg{
			Name:         name,
			Type:         "bytes",
			InternalType: "bytes",
		}
	case TypeString:
		return jsonAbiArg{
			Name:         name,
			Type:         "string",
			InternalType: "string",
		}
	case TypeArray:
		elemType := generateArg(ast, name, type_.Elem)
		elemType.Type += "[]"
		elemType.InternalType += "[]"
		return elemType
	case TypeStructRef:
		struct_ := ast.Structs[type_.Index]
		var components []jsonAbiArg
		for _, field := range struct_.Fields {
			arg := generateArg(ast, field.Name, field.Type)
			components = append(components, arg)
		}
		return jsonAbiArg{
			Name:         name,
			Type:         "tuple",
			InternalType: "struct " + struct_.Name,
			Components:   components,
		}
	default:
		// This should not happen
		panic(fmt.Errorf("invalid type: %T", type_))
	}
}
