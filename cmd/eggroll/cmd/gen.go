// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"fmt"
	"os"

	"github.com/gligneul/eggroll/internal/compiler"
	"github.com/spf13/cobra"
)

var genArgs struct {
	inputPath string
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate ABI bindings",
	Long:  `Generate the Go bindings for the given ABI yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(genArgs.inputPath)
		cobra.CheckErr(err)
		abiJson, err := compiler.Compile(f)
		cobra.CheckErr(err)
		fmt.Println(abiJson)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVar(
		&genArgs.inputPath, "input", "abi.yaml", "input file that contains the ABI specification")
}
