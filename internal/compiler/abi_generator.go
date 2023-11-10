// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"encoding/json"
	"fmt"
)

type jsonAbiMethod struct {
	Name            string       `json:"name"`
	Type            string       `json:"type"`
	StateMutability string       `json:"stateMutability"`
	Inputs          []jsonAbiArg `json:"inputs"`
	Outputs         []jsonAbiArg `json:"outputs"`
}

type jsonAbiArg struct {
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	InternalType string       `json:"internalType"`
	Components   []jsonAbiArg `json:"components"`
}

// Generate the JSON ABI for the AST.
// This function converts struct to tuples and message schemas to solidity functions.
func generateAbi(ast astSchema) []byte {
	methods := generateAbiMethods(nil, ast.Reports, ast.Structs)
	methods = generateAbiMethods(methods, ast.Advances, ast.Structs)
	methods = generateAbiMethods(methods, ast.Inspects, ast.Structs)
	result, err := json.MarshalIndent(methods, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("json marshal error: %v", err))
	}
	return result
}

// Generate the methods and append them to the slice
func generateAbiMethods(
	methods []jsonAbiMethod,
	messages []messageSchema,
	structs []messageSchema,
) []jsonAbiMethod {
	for _, message := range messages {
		var method jsonAbiMethod
		method.Name = message.Name
		method.Type = "function"
		method.StateMutability = "nonpayable"
		for _, field := range message.Fields {
			arg := generateAbiArg(field.Name, field.type_, structs)
			method.Inputs = append(method.Inputs, arg)
		}
		methods = append(methods, method)
	}
	return methods
}

// Recursively generate the ABI argument given the type.
func generateAbiArg(name string, type_ any, structs []messageSchema) jsonAbiArg {
	switch type_ := type_.(type) {
	case typeBool:
		return jsonAbiArg{
			Name:         name,
			Type:         "bool",
			InternalType: "bool",
		}
	case typeInt:
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
	case typeAddress:
		return jsonAbiArg{
			Name:         name,
			Type:         "address",
			InternalType: "address",
		}
	case typeBytes:
		return jsonAbiArg{
			Name:         name,
			Type:         "bytes",
			InternalType: "bytes",
		}
	case typeString:
		return jsonAbiArg{
			Name:         name,
			Type:         "string",
			InternalType: "string",
		}
	case typeArray:
		elemType := generateAbiArg(name, type_.Elem, structs)
		elemType.Type += "[]"
		elemType.InternalType += "[]"
		return elemType
	case typeStructRef:
		struct_ := structs[type_.Index]
		var components []jsonAbiArg
		for _, field := range struct_.Fields {
			arg := generateAbiArg(field.Name, field.type_, structs)
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
