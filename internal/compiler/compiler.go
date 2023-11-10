// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Package responsible for compiling EggRoll schema definition to the Go binding.
package compiler

// Compile the input into Solidity JSON ABI.
func YamlSchemaToJsonAbi(input []byte) ([]byte, error) {
	ast, err := analyze(input)
	if err != nil {
		return nil, err
	}
	return generateAbi(ast), nil
}

// Compile the input into the EggRoll Go binding.
func YamlSchemaToGoBinding(input []byte, packageName string) ([]byte, error) {
	ast, err := analyze(input)
	if err != nil {
		return nil, err
	}
	return generateGo(ast, packageName), nil
}
