// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"testing"
)

func testGenerateAbi(t *testing.T, input string, expected string) {
	ast, err := analyze([]byte(input))
	if err != nil {
		t.Fatalf("failed to analyze: %v", err)
	}
	generated := string(generateAbi(ast))
	if generated != expected {
		t.Fatalf("wrong json: %v", generated)
	}
}

func TestGenerateAbiForAllMessageSchemas(t *testing.T) {
	input := `
reports:
  - name: reportMessage

advances:
  - name: advanceMessage

inspects:
  - name: inspectMessage
`
	expected := `[
  {
    "name": "reportMessage",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": null,
    "outputs": null
  },
  {
    "name": "advanceMessage",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": null,
    "outputs": null
  },
  {
    "name": "inspectMessage",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": null,
    "outputs": null
  }
]`
	testGenerateAbi(t, input, expected)
}

func TestGenerateAbiBasicTypes(t *testing.T) {
	input := `
reports:
  - name: foo
    fields:
      - name: bool
        type: bool
      - name: int
        type: int
      - name: int8
        type: int8
      - name: int256
        type: int256
      - name: uint
        type: uint
      - name: uint8
        type: uint8
      - name: uint256
        type: uint256
      - name: address
        type: address
      - name: string
        type: string
      - name: bytes
        type: bytes
`
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
	testGenerateAbi(t, input, expected)
}

func TestGenerateAbiArrayType(t *testing.T) {
	input := `
reports:
  - name: foo
    fields:
      - name: i
        type: int[]
`
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
	testGenerateAbi(t, input, expected)
}

func TestGenerateAbiStructType(t *testing.T) {
	input := `
structs:
  - name: foo
    fields:
      - name: ivalue
        type: int

reports:
  - name: bar
    fields:
      - name: foovalue
        type: foo
`
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
	testGenerateAbi(t, input, expected)
}
