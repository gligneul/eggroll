// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Package responsible for parsing EggRoll ABI spec.
package abiparser

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func Parse(input io.Reader) (abi.ABI, error) {
	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read input: %v", err)
	}
	ast, err := parse(inputBytes)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("parse error: %v", err)
	}
	_ = ast
	return abi.ABI{}, nil
}
