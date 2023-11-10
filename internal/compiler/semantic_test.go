// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"reflect"
	"testing"
)

func TestFailToAnalyzeDuplicateStruct(t *testing.T) {
	ast, err := analyze([]byte(`
structs:
  - name: foo
    fields:
      - name: bar
        type: int
  - name: foo
    fields:
      - name: bar
        type: int
`))
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `struct duplicate of "foo"` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeStructWithNoFields(t *testing.T) {
	ast, err := analyze([]byte(`
structs:
  - name: foo
`))
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `struct foo: must have fields` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeMessageWithStructName(t *testing.T) {
	ast, err := analyze([]byte(`
structs:
  - name: foo
    fields:
      - name: bar
        type: int
reports:
  - name: foo
`))
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `report duplicate of "foo"` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeDuplicateMessage(t *testing.T) {
	ast, err := analyze([]byte(`
reports:
  - name: foo
advances:
  - name: foo
`))
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `advance duplicate of "foo"` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeStructNotFound(t *testing.T) {
	ast, err := analyze([]byte(`
structs:
  - name: foo
    fields:
      - name: bar
        type: foo
`))
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `struct foo: field bar: struct "foo" not found` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestAnalyzeEmpty(t *testing.T) {
	_, err := analyze([]byte(``))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestAnalyzeSingleMessage(t *testing.T) {
	_, err := analyze([]byte(`
reports:
  - name: foo
`))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestAnalyzeTypes(t *testing.T) {
	ast, err := analyze([]byte(`
structs:
  - name: foo
    fields:
      - name: x
        type: bool
  - name: bar
    fields:
      - name: y
        type: bool
reports:
  - name: report
    fields:
      - name: foo
        type: foo
      - name: barArray
        type: bar[]
`))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	expected := astSchema{
		Structs: []messageSchema{
			{
				Name: "foo",
				Fields: []fieldSchema{
					{Name: "x", Type: "bool", type_: typeBool{}},
				},
			},
			{
				Name: "bar",
				Fields: []fieldSchema{
					{Name: "y", Type: "bool", type_: typeBool{}},
				},
			},
		},
		Reports: []messageSchema{
			{
				Name: "report",
				Fields: []fieldSchema{
					{
						Name: "foo",
						Type: "foo",
						type_: typeStructRef{
							Name:  "foo",
							Index: 0,
						},
					},
					{
						Name: "barArray",
						Type: "bar[]",
						type_: typeArray{
							Elem: typeStructRef{
								Name:  "bar",
								Index: 1,
							},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(ast, expected) {
		t.Fatalf("wrong AST: %#v", ast)
	}
}
