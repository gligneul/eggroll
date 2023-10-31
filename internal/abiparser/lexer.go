// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package abiparser

import (
	"fmt"
	"regexp"
)

func checkIdentifier(id string) error {
	pattern := "^[a-zA-Z][a-zA-Z0-9]*$"
	match, err := regexp.MatchString(pattern, id)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("%q does not match %v", id, pattern)
	}
	return nil
}

func tokenizeType(rawType string) (identifier string, isArray bool, err error) {
	pattern := `^([a-zA-Z][a-zA-Z0-9]*)(\[\])?$`
	exp := regexp.MustCompile(pattern)
	matches := exp.FindStringSubmatch(rawType)
	if len(matches) != 3 {
		return "", false, fmt.Errorf("%q does not match %v", rawType, pattern)
	}
	return matches[1], matches[2] == "[]", nil
}
