// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"testing"
)

func TestGenerateBasicTypes(t *testing.T) {
	ast := Ast{
		Messages: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "bool", Type: TypeBool{}},
					{Name: "int", Type: TypeInt{true, 256}},
					{Name: "int8", Type: TypeInt{true, 8}},
					{Name: "int256", Type: TypeInt{true, 256}},
					{Name: "uint", Type: TypeInt{false, 256}},
					{Name: "uint8", Type: TypeInt{false, 8}},
					{Name: "uint256", Type: TypeInt{false, 256}},
					{Name: "address", Type: TypeAddress{}},
					{Name: "string", Type: TypeString{}},
					{Name: "bytes", Type: TypeBytes{}},
				},
			},
		},
	}
	expected := `[
  {
    "name": "foo",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": [
      {
        "name": "bool",
        "type": "bool",
        "internalType": "bool",
        "components": null
      },
      {
        "name": "int",
        "type": "int256",
        "internalType": "int256",
        "components": null
      },
      {
        "name": "int8",
        "type": "int8",
        "internalType": "int8",
        "components": null
      },
      {
        "name": "int256",
        "type": "int256",
        "internalType": "int256",
        "components": null
      },
      {
        "name": "uint",
        "type": "uint256",
        "internalType": "uint256",
        "components": null
      },
      {
        "name": "uint8",
        "type": "uint8",
        "internalType": "uint8",
        "components": null
      },
      {
        "name": "uint256",
        "type": "uint256",
        "internalType": "uint256",
        "components": null
      },
      {
        "name": "address",
        "type": "address",
        "internalType": "address",
        "components": null
      },
      {
        "name": "string",
        "type": "string",
        "internalType": "string",
        "components": null
      },
      {
        "name": "bytes",
        "type": "bytes",
        "internalType": "bytes",
        "components": null
      }
    ],
    "outputs": null
  }
]`
	json := generate(ast)
	if json != expected {
		t.Fatalf("wrong json: %v", json)
	}
}

func TestGenerateArrayType(t *testing.T) {
	ast := Ast{
		Messages: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "i", Type: TypeArray{Elem: TypeInt{true, 256}}},
				},
			},
		},
	}
	expected := `[
  {
    "name": "foo",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": [
      {
        "name": "i",
        "type": "int256[]",
        "internalType": "int256[]",
        "components": null
      }
    ],
    "outputs": null
  }
]`
	json := generate(ast)
	if json != expected {
		t.Fatalf("wrong json: %v", json)
	}
}

func TestGenerateStructType(t *testing.T) {
	ast := Ast{
		Structs: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "ivalue", Type: TypeInt{true, 256}},
				},
			},
		},
		Messages: []Struct{
			{
				Name: "bar",
				Fields: []Field{
					{Name: "foovalue", Type: TypeStructRef{0}},
				},
			},
		},
	}
	expected := `[
  {
    "name": "bar",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": [
      {
        "name": "foovalue",
        "type": "tuple",
        "internalType": "struct foo",
        "components": [
          {
            "name": "ivalue",
            "type": "int256",
            "internalType": "int256",
            "components": null
          }
        ]
      }
    ],
    "outputs": null
  }
]`
	json := generate(ast)
	if json != expected {
		t.Fatalf("wrong json: %v", json)
	}
}
