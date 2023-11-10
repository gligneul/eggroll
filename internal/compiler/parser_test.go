// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"reflect"
	"strings"
	"testing"
)

func TestFailToParseInvalidYamlFile(t *testing.T) {
	ast, err := parse([]byte(`---
invalid file
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), "yaml: unmarshal errors") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseStructWithInvalidName(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - name: invalid_name
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if err.Error() != `struct name: invalid rune '_'` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseFieldWithInvalidName(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - name: foo
    fields:
    - name: invalid_name
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if err.Error() != `struct foo: field name: invalid rune '_'` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseInvalidType(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - name: foo
    fields:
    - name: bar
      type: invalid_type
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if err.Error() != `struct foo.bar type: invalid rune '_'` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestParseEmpty(t *testing.T) {
	parsedAst, err := parse([]byte(``))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	expectedAst := astSchema{}
	if !reflect.DeepEqual(parsedAst, expectedAst) {
		t.Fatalf("wrong AST: %+v", parsedAst)
	}
}

func TestParseTypes(t *testing.T) {
	ast, err := parse([]byte(`---
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
      - name: array
        type: bool[]
      - name: structRef
        type: bar
`))
	expected := astSchema{
		Reports: []messageSchema{
			{
				Name: "foo",
				Fields: []fieldSchema{
					{Name: "bool", Type: "bool", type_: typeBool{}},
					{Name: "int", Type: "int", type_: typeInt{true, 256}},
					{Name: "int8", Type: "int8", type_: typeInt{true, 8}},
					{Name: "int256", Type: "int256", type_: typeInt{true, 256}},
					{Name: "uint", Type: "uint", type_: typeInt{false, 256}},
					{Name: "uint8", Type: "uint8", type_: typeInt{false, 8}},
					{Name: "uint256", Type: "uint256", type_: typeInt{false, 256}},
					{Name: "address", Type: "address", type_: typeAddress{}},
					{Name: "string", Type: "string", type_: typeString{}},
					{Name: "bytes", Type: "bytes", type_: typeBytes{}},
					{Name: "array", Type: "bool[]", type_: typeArray{Elem: typeBool{}}},
					{Name: "structRef", Type: "bar", type_: typeStructRef{Name: "bar"}},
				},
			},
		},
	}
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !reflect.DeepEqual(expected, ast) {
		t.Fatalf("wrong ast: %#v", ast)
	}
}

func TestParseValid(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - name: Chapter
    doc: A chapter of a book.
    fields:
      - name: title
        doc: Title of the chapter.
        type: string
      - name: text
        doc: Text of the chapter.
        type: string

  - name: Book
    doc: Representation of a book in the blockchain.
    fields:
      - name: id
        doc: Identifier of the book.
        type: uint256
      - name: title
        doc: Title of the book.
        type: string
      - name: author
        doc: Address of the author of the book.
        type: address
      - name: chapters
        doc: Chapters of the book.
        type: Chapter[]

reports:
  - name: bookAdded
    doc: Report emitted when a book is added.
    fields:
      - name: book
        type: Book

advances:
  - name: addBook
    doc: Add a new book to the contract
    fields:
      - name: book
        type: Book

inspects:
  - name: getBook
    doc: Get the book by ID. If found, emit an bookAdded event.
    fields:
      - name: id
        type: uint256
`))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	expectedAst := astSchema{
		Structs: []messageSchema{
			{
				Name: "Chapter",
				Doc:  "A chapter of a book.",
				Fields: []fieldSchema{
					{
						Name:  "title",
						Doc:   "Title of the chapter.",
						Type:  "string",
						type_: typeString{},
					},
					{
						Name:  "text",
						Doc:   "Text of the chapter.",
						Type:  "string",
						type_: typeString{},
					},
				},
			},
			{
				Name: "Book",
				Doc:  "Representation of a book in the blockchain.",
				Fields: []fieldSchema{
					{
						Name:  "id",
						Doc:   "Identifier of the book.",
						Type:  "uint256",
						type_: typeInt{Signed: false, Bits: 256},
					},
					{
						Name:  "title",
						Doc:   "Title of the book.",
						Type:  "string",
						type_: typeString{},
					},
					{
						Name:  "author",
						Doc:   "Address of the author of the book.",
						Type:  "address",
						type_: typeAddress{},
					},
					{
						Name: "chapters",
						Doc:  "Chapters of the book.",
						Type: "Chapter[]",
						type_: typeArray{
							Elem: typeStructRef{
								Name:  "Chapter",
								Index: 0,
							},
						},
					},
				},
			},
		},
		Reports: []messageSchema{
			{
				Name: "bookAdded",
				Doc:  "Report emitted when a book is added.",
				Fields: []fieldSchema{
					{
						Name:  "book",
						Doc:   "",
						Type:  "Book",
						type_: typeStructRef{Name: "Book", Index: 0},
					},
				},
			},
		},
		Advances: []messageSchema{
			{
				Name: "addBook",
				Doc:  "Add a new book to the contract",
				Fields: []fieldSchema{
					{
						Name:  "book",
						Doc:   "",
						Type:  "Book",
						type_: typeStructRef{Name: "Book", Index: 0},
					},
				},
			},
		},
		Inspects: []messageSchema{
			{
				Name: "getBook",
				Doc:  "Get the book by ID. If found, emit an bookAdded event.",
				Fields: []fieldSchema{
					{
						Name:  "id",
						Doc:   "",
						Type:  "uint256",
						type_: typeInt{Signed: false, Bits: 256},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(ast, expectedAst) {
		t.Fatalf("wrong AST: %#v", ast)
	}
}
