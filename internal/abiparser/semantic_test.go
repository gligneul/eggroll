// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package abiparser

import (
	"reflect"
	"testing"
)

func TestFailToAnalyzeDuplicateStruct(t *testing.T) {
	ast := Ast{
		Structs: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "bar", Type: TypeName{"int"}},
				},
			},
			{
				Name: "foo",
				Fields: []Field{
					{Name: "bar", Type: TypeName{"int"}},
				},
			},
		},
	}
	ast, err := analyze(ast)
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `duplicate struct: "foo"` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeStructWithNoFields(t *testing.T) {
	ast := Ast{
		Structs: []Struct{
			{
				Name: "foo",
			},
		},
	}
	ast, err := analyze(ast)
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `struct foo: must have fields` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeMessageWithStructName(t *testing.T) {
	ast := Ast{
		Structs: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "bar", Type: TypeName{"int"}},
				},
			},
		},
		Messages: []Struct{
			{
				Name: "foo",
			},
		},
	}
	ast, err := analyze(ast)
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `message with struct name: "foo"` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeDuplicateMessage(t *testing.T) {
	ast := Ast{
		Messages: []Struct{
			{
				Name: "foo",
			},
			{
				Name: "foo",
			},
		},
	}
	ast, err := analyze(ast)
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `duplicate message: "foo"` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeNoMessages(t *testing.T) {
	ast := Ast{}
	ast, err := analyze(ast)
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `no messages` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestFailToAnalyzeTypeNotFound(t *testing.T) {
	ast := Ast{
		Structs: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "bar", Type: TypeName{"foo"}},
				},
			},
		},
	}
	ast, err := analyze(ast)
	if err == nil {
		t.Fatalf("expected err; got %+v", ast)
	}
	if err.Error() != `struct foo: field bar: type not found "foo"` {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestAnalyzeSingleMessage(t *testing.T) {
	ast := Ast{
		Messages: []Struct{
			{
				Name: "foo",
			},
		},
	}
	retAst, err := analyze(ast)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !reflect.DeepEqual(retAst, ast) {
		t.Fatalf("wrong ast: %+v", retAst)
	}
}

func TestAnalyzeBasicTypes(t *testing.T) {
	ast := Ast{
		Messages: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "bool", Type: TypeName{"bool"}},
					{Name: "int", Type: TypeName{"int"}},
					{Name: "int8", Type: TypeName{"int8"}},
					{Name: "int256", Type: TypeName{"int256"}},
					{Name: "uint", Type: TypeName{"uint"}},
					{Name: "uint8", Type: TypeName{"uint8"}},
					{Name: "uint256", Type: TypeName{"uint256"}},
					{Name: "address", Type: TypeName{"address"}},
					{Name: "string", Type: TypeName{"string"}},
					{Name: "bytes", Type: TypeName{"bytes"}},
				},
			},
		},
	}
	expAst := Ast{
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
	retAst, err := analyze(ast)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !reflect.DeepEqual(retAst, expAst) {
		t.Fatalf("wrong ast: %+v", retAst)
	}
}

func TestAnalyzeArrayTypes(t *testing.T) {
	ast := Ast{
		Messages: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "i", Type: TypeArray{Elem: TypeName{"int"}}},
				},
			},
		},
	}
	expAst := Ast{
		Messages: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "i", Type: TypeArray{Elem: TypeInt{true, 256}}},
				},
			},
		},
	}
	retAst, err := analyze(ast)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !reflect.DeepEqual(retAst, expAst) {
		t.Fatalf("wrong ast: %+v", retAst)
	}
}

func TestAnalyzeStructType(t *testing.T) {
	ast := Ast{
		Structs: []Struct{
			{
				Name: "foo",
				Fields: []Field{
					{Name: "ivalue", Type: TypeName{"int"}},
				},
			},
		},
		Messages: []Struct{
			{
				Name: "bar",
				Fields: []Field{
					{Name: "foovalue", Type: TypeName{"foo"}},
				},
			},
		},
	}
	expAst := Ast{
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
	retAst, err := analyze(ast)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !reflect.DeepEqual(retAst, expAst) {
		t.Fatalf("wrong ast: %+v", retAst)
	}
}

//func TestFailToAnalyze(t *testing.T) {
