// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package abiparser

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

func TestFailToParseMessage(t *testing.T) {
	// We only do a single test for messages because the rest of the code
	// is shared with the struct parsing.
	ast, err := parse([]byte(`---
messages:
  -
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), "invalid message syntax") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseStructWithNoEntries(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  -
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), "invalid struct syntax") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseStructWithMultipleEntries(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - foo:
    bar:
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), "invalid struct syntax") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseStructWithInvalidName(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - invalid_name:
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), `struct name: "invalid_name" does not match`) {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseFieldWithNoEntries(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - foo:
    -
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), "invalid field syntax") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseFieldWithMultipleEntries(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - foo:
    - foo:
      bar:
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), "invalid field syntax") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseFieldWithInvalidName(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - foo:
    - invalid_name:
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), `struct foo: field name: "invalid_name" does not match`) {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToParseInvalidType(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - foo:
    - bar:
`))
	if err == nil {
		t.Fatalf("expected error; got %+v", ast)
	}
	if !strings.Contains(err.Error(), `struct foo: field bar: "" does not match`) {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestParseEmpty(t *testing.T) {
	ast, err := parse([]byte(""))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	expectedAst := Ast{}
	if !reflect.DeepEqual(ast, expectedAst) {
		t.Fatalf("wrong AST: %+v", ast)
	}
}

func TestParseValid(t *testing.T) {
	ast, err := parse([]byte(`---
structs:
  - Chapter:
    - title: string

  - Book:
    - id: uint
    - title: book
    - author: string
    - chapters: Chapter[]

messages:
  - AddBook:
    - book: Book
`))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	expectedAst := Ast{
		Structs: []Struct{
			{
				Name: "Chapter",
				Fields: []Field{
					{
						Name: "title",
						Type: TypeName{Name: "string"},
					},
				},
			},
			{
				Name: "Book",
				Fields: []Field{
					{
						Name: "id",
						Type: TypeName{Name: "uint"},
					},
					{
						Name: "title",
						Type: TypeName{Name: "book"},
					},
					{
						Name: "author",
						Type: TypeName{Name: "string"},
					},
					{
						Name: "chapters",
						Type: TypeArray{Elem: TypeName{Name: "Chapter"}},
					},
				},
			},
		},
		Messages: []Struct{
			{
				Name: "AddBook",
				Fields: []Field{
					{
						Name: "book",
						Type: TypeName{Name: "Book"},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(ast, expectedAst) {
		t.Fatalf("wrong AST: %#v", ast)
	}
}
