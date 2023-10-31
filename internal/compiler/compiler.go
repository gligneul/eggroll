// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Package responsible for compiling EggRoll ABI to Solidity JSON ABI.
package compiler

import (
	"fmt"
	"io"
)

func Compile(input io.Reader) (string, error) {
	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %v", err)
	}
	ast, err := parse(inputBytes)
	if err != nil {
		return "", err
	}
	ast, err = analyze(ast)
	if err != nil {
		return "", err
	}
	return generate(ast), nil
}
