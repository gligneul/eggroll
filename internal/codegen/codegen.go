// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// Generate the EggRoll Go ABI binding for the given JSON ABI.
func Gen(jsonAbi string, packageName string) (string, error) {
	var data tmplData
	data.Package = packageName
	data.JsonAbi = jsonAbi
	a, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return "", err
	}
	data.Structs, err = loadStructs(a)
	if err != nil {
		return "", err
	}
	data.Schemas, err = loadSchemas(a, data.Structs)
	if err != nil {
		return "", err
	}
	return genCode(data)
}

var keyword = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"defer":       true,
	"else":        true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"iota":        true,
	"map":         true,
	"make":        true,
	"new":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

// Convert a name to a Go identifier.
func toGoIdentifier(name string) (string, error) {
	if keyword[name] {
		return "", fmt.Errorf("%v is a keyword", name)
	}
	capitalized := strings.ToUpper(name[0:1]) + name[1:]
	return capitalized, nil
}

// Get the identifier for the given tuple type.
func structID(t abi.Type) string {
	return t.TupleRawName + t.String()
}

// Convert an ABI type to a Go type.
func toGoType(t abi.Type, structs map[string]tmplStruct) (string, error) {
	switch t.T {
	case abi.IntTy, abi.UintTy:
		parts := regexp.MustCompile(`(u)?int([0-9]*)`).FindStringSubmatch(t.String())
		switch parts[2] {
		case "8", "16", "32", "64":
			return fmt.Sprintf("%sint%s", parts[1], parts[2]), nil
		}
		return "*big.Int", nil
	case abi.BoolTy:
		return "bool", nil
	case abi.StringTy:
		return "string", nil
	case abi.SliceTy:
		elemType, err := toGoType(*t.Elem, structs)
		if err != nil {
			return "", err
		}
		return "[]" + elemType, nil
	case abi.ArrayTy:
		elemType, err := toGoType(*t.Elem, structs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%d]", t.Size) + elemType, nil
	case abi.TupleTy:
		return structs[structID(t)].Kind, nil
	case abi.AddressTy:
		return "common.Address", nil
	case abi.FixedBytesTy:
		return fmt.Sprintf("[%d]byte", t.Size), nil
	case abi.BytesTy:
		return "[]byte", nil
	case abi.FunctionTy:
		return "[24]byte", nil
	default:
		return "", fmt.Errorf("type not supported: %v", t)
	}
}

// Add structs by recursively traversing the type tree.
func recAddStructs(structs map[string]tmplStruct, t abi.Type) error {
	if t.T == abi.SliceTy || t.T == abi.ArrayTy {
		return recAddStructs(structs, *t.Elem)
	}
	if t.T != abi.TupleTy {
		return nil
	}
	for _, elem := range t.TupleElems {
		err := recAddStructs(structs, *elem)
		if err != nil {
			return err
		}
	}
	var err error
	var struct_ tmplStruct
	struct_.Kind = t.TupleRawName
	struct_.GoName, err = toGoIdentifier(t.TupleRawName)
	if err != nil {
		return err
	}
	for i, name := range t.TupleRawNames {
		var field tmplField
		field.Kind = name
		field.GoName, err = toGoIdentifier(name)
		if err != nil {
			return err
		}
		field.Type, err = toGoType(*t.TupleElems[i], structs)
		if err != nil {
			return err
		}
		struct_.Fields = append(struct_.Fields, field)
	}
	structs[structID(t)] = struct_
	return nil
}

func loadStructs(a abi.ABI) (map[string]tmplStruct, error) {
	structs := make(map[string]tmplStruct)
	for _, method := range a.Methods {
		for _, input := range method.Inputs {
			err := recAddStructs(structs, input.Type)
			if err != nil {
				return nil, err
			}
		}
	}
	return structs, nil
}

func loadSchemas(a abi.ABI, structs map[string]tmplStruct) ([]tmplSchema, error) {
	var schemas []tmplSchema
	var err error
	for name, method := range a.Methods {
		var schema tmplSchema
		schema.ID = method.ID
		schema.Kind = name
		schema.GoName, err = toGoIdentifier(name)
		if err != nil {
			return nil, err
		}
		for _, input := range method.Inputs {
			var field tmplField
			field.Kind = input.Name
			field.GoName, err = toGoIdentifier(input.Name)
			if err != nil {
				return nil, err
			}
			field.Type, err = toGoType(input.Type, structs)
			if err != nil {
				return nil, err
			}
			schema.Fields = append(schema.Fields, field)
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

func genCode(data tmplData) (string, error) {
	tmpl := template.Must(template.New("eggroll").Parse(tmplSource))
	var codeBuffer bytes.Buffer
	err := tmpl.Execute(&codeBuffer, data)
	if err != nil {
		return "", err
	}
	code, err := format.Source(codeBuffer.Bytes())
	if err != nil {
		// This should not happen
		panic(fmt.Errorf("%v\n%v", err, codeBuffer.String()))
	}
	return string(code), nil
}
