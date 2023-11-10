// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Generate the test bindings.
package main

import (
	"io"
	"os"

	"github.com/gligneul/eggroll/internal/compiler"
	"github.com/spf13/cobra"
)

const packageName = "testbinding"
const inputPath = packageName + "/schema.yaml"
const outputPath = packageName + "/schema.go"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	inputFile, err := os.Open(inputPath)
	checkErr(err)
	defer inputFile.Close()

	input, err := io.ReadAll(inputFile)
	checkErr(err)

	output, err := compiler.YamlSchemaToGoBinding(input, packageName)
	checkErr(err)

	outputFile, err := os.Create(outputPath)
	cobra.CheckErr(err)
	defer outputFile.Close()

	_, err = outputFile.Write(output)
	checkErr(err)
}
