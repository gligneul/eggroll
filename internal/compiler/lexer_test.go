// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import "testing"

func TestCheckIdentifier(t *testing.T) {
	testCases := []struct {
		id    string
		isErr bool
	}{
		{"f", false},
		{"foo", false},
		{"fooBar", false},
		{"F", false},
		{"Foo", false},
		{"FB", false},
		{"FooBar", false},
		{"FooBar123", false},
		{"", true},
		{"_", true},
		{"_foo", true},
		{"1foo", true},
		{"Foo_", true},
		{"var", true},
	}
	for _, testCase := range testCases {
		err := checkName(testCase.id)
		if !testCase.isErr && err != nil {
			t.Fatalf("unexpected err for %v: %v", testCase.id, err)
		}
		if testCase.isErr {
			t.Log(err)
			if err == nil {
				t.Fatalf("expected err for %v", testCase.id)
			}
		}
	}
}

func TestTokenizeType(t *testing.T) {
	testCases := []struct {
		rawType string
		id      string
		isArray bool
		isErr   bool
	}{
		// valid ids
		{"f", "f", false, false},
		{"foo", "foo", false, false},
		{"fooBar", "fooBar", false, false},
		{"F", "F", false, false},
		{"Foo", "Foo", false, false},
		{"FB", "FB", false, false},
		{"FooBar", "FooBar", false, false},
		{"FooBar123", "FooBar123", false, false},

		// valid arrays
		{"f[]", "f", true, false},
		{"fooBar[]", "fooBar", true, false},
		{"foo123[]", "foo123", true, false},

		// error cases
		{"[]", "", false, true},
		{"foo[", "", false, true},
		{"fo[o]", "", false, true},
		{"foo]", "", false, true},
		{"foo_bar", "", false, true},
		{"_", "", false, true},
		{"_foo", "", false, true},
		{"1foo", "", false, true},
	}
	for _, testCase := range testCases {
		id, isArray, err := tokenizeType(testCase.rawType)
		if !testCase.isErr {
			if err != nil {
				t.Fatalf("unexpected err for %q: %v", testCase.rawType, err)
			}
			if id != testCase.id || isArray != testCase.isArray {
				t.Fatalf("expected %q, %v for %q; got %q, %v",
					testCase.id, testCase.isArray, testCase.rawType, id, isArray)
			}
		}
		if testCase.isErr {
			t.Log(err)
			if err == nil {
				t.Fatalf("expected err for %v", testCase.id)
			}
		}
	}
}
