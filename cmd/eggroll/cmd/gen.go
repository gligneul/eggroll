// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"os"

	"github.com/gligneul/eggroll/internal/codegen"
	"github.com/gligneul/eggroll/internal/compiler"
	"github.com/spf13/cobra"
)

var genArgs struct {
	inputPath   string
	outputPath  string
	packageName string
}

var genCmd = &cobra.Command{
	Use:   "abi-gen",
	Short: "Generate ABI bindings",
	Long:  `Generate the Go bindings for the given ABI yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, err := os.Open(genArgs.inputPath)
		cobra.CheckErr(err)
		defer inputFile.Close()

		abiJson, err := compiler.Compile(inputFile)
		cobra.CheckErr(err)

		code, err := codegen.Gen(abiJson, genArgs.packageName)
		cobra.CheckErr(err)

		outputFile, err := os.Create(genArgs.outputPath)
		cobra.CheckErr(err)
		defer outputFile.Close()
		outputFile.Write([]byte(code))
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVar(
		&genArgs.inputPath, "input", "abi.yaml", "Input file that contains the ABI yaml")

	genCmd.Flags().StringVar(
		&genArgs.outputPath, "output", "abi.go", "Output file for the generated Go binding")

	genCmd.Flags().StringVar(
		&genArgs.packageName, "package", "main", "Name of the generated package")
}
