// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"fmt"
	"regexp"
	"strings"
)

func checkEmpty(name string) error {
	if len(name) < 1 {
		return fmt.Errorf("empty name")
	}
	return nil
}

func checkFirstChar(name string) error {
	pattern := "^[a-zA-Z].*"
	regexpPattern := regexp.MustCompile(pattern)
	if !regexpPattern.MatchString(name) {
		return fmt.Errorf("invalid first rune '%c'", name[0])
	}
	return nil
}

func checkChars(name string) error {
	pattern := "[^a-zA-Z0-9]"
	regexpPattern := regexp.MustCompile(pattern)
	match := regexpPattern.FindString(name)
	if match != "" {
		return fmt.Errorf("invalid rune '%v'", match)
	}
	return nil
}

func checkKeyword(name string) error {
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
	_, ok := keyword[name]
	if ok {
		return fmt.Errorf("%s is a Go keyword", name)
	}
	return nil
}

func checkName(name string) error {
	checks := []func(string) error{
		checkEmpty,
		checkFirstChar,
		checkChars,
		checkKeyword,
	}
	for _, check := range checks {
		if err := check(name); err != nil {
			return err
		}
	}
	return nil
}

func tokenizeType(rawType string) (name string, isArray bool, err error) {
	openBracketIndex := strings.IndexRune(rawType, '[')
	if openBracketIndex != -1 {
		isArray = true
		if len(rawType) != (openBracketIndex+2) || rawType[openBracketIndex+1] != ']' {
			return "", false, fmt.Errorf("invalid array; only [] is supported")
		}
		name = rawType[:openBracketIndex]
	} else {
		name = rawType
	}
	if err := checkName(name); err != nil {
		return "", false, err
	}
	return name, isArray, nil
}
