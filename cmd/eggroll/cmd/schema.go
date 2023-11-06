// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gligneul/eggroll/internal/compiler"
	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/spf13/cobra"
)

var schemaArgs struct {
	yamlPath string
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Commands related to schema encoding and decoding",
}

// Load the schema into eggtypes and return the JSON ABI.
func schemaLoad() string {
	inputFile, err := os.Open(schemaArgs.yamlPath)
	cobra.CheckErr(err)
	defer inputFile.Close()

	jsonAbi, err := compiler.Compile(inputFile)
	cobra.CheckErr(err)

	a, err := abi.JSON(strings.NewReader(jsonAbi))
	cobra.CheckErr(err)

	for _, method := range a.Methods {
		err := eggtypes.AddSchema(eggtypes.MessageSchema{
			ID:        eggtypes.ID(method.ID),
			Kind:      method.Name,
			Arguments: method.Inputs,
		})
		cobra.CheckErr(err)
	}

	return jsonAbi
}

func init() {
	rootCmd.AddCommand(schemaCmd)

	schemaCmd.PersistentFlags().StringVar(
		&schemaArgs.yamlPath, "schema", "schema.yaml", "Yaml file that contains the schema")
}
