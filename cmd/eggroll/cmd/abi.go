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

var abiArgs struct {
	yamlPath string
}

var abiCmd = &cobra.Command{
	Use:   "abi",
	Short: "Commands related to ABI encoding",
}

// Load the ABI into eggtypes and return the JSON ABI.
func abiLoad() string {
	inputFile, err := os.Open(abiArgs.yamlPath)
	cobra.CheckErr(err)
	defer inputFile.Close()

	jsonAbi, err := compiler.Compile(inputFile)
	cobra.CheckErr(err)

	a, err := abi.JSON(strings.NewReader(jsonAbi))
	cobra.CheckErr(err)

	for _, method := range a.Methods {
		err := eggtypes.AddEncoding(eggtypes.Encoding{
			ID:        eggtypes.ID(method.ID),
			Name:      method.Name,
			Arguments: method.Inputs,
		})
		cobra.CheckErr(err)
	}

	return jsonAbi
}

func init() {
	rootCmd.AddCommand(abiCmd)

	abiCmd.PersistentFlags().StringVar(
		&abiArgs.yamlPath, "input", "abi.yaml", "Input file that contains the ABI yaml")
}
