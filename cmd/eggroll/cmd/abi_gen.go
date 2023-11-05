// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"os"

	"github.com/gligneul/eggroll/internal/codegen"
	"github.com/gligneul/eggroll/internal/compiler"
	"github.com/spf13/cobra"
)

var abiGenArgs struct {
	outputPath  string
	packageName string
}

var abiGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate ABI bindings",
	Long:  `Generate the Go bindings for the given ABI yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, err := os.Open(abiArgs.yamlPath)
		cobra.CheckErr(err)
		defer inputFile.Close()

		abiJson, err := compiler.Compile(inputFile)
		cobra.CheckErr(err)

		code, err := codegen.Gen(abiJson, abiGenArgs.packageName)
		cobra.CheckErr(err)

		outputFile, err := os.Create(abiGenArgs.outputPath)
		cobra.CheckErr(err)
		defer outputFile.Close()
		outputFile.Write([]byte(code))
	},
}

func init() {
	abiCmd.AddCommand(abiGenCmd)

	abiGenCmd.Flags().StringVar(
		&abiGenArgs.outputPath, "output", "abi.go", "Output file for the generated Go binding")

	abiGenCmd.Flags().StringVar(
		&abiGenArgs.packageName, "package", "main", "Name of the generated package")
}
