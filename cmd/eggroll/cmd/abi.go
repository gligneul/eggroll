// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"github.com/spf13/cobra"
)

var abiArgs struct {
	yamlPath string
}

var abiCmd = &cobra.Command{
	Use:   "abi",
	Short: "Commands related to ABI encoding",
}

func init() {
	rootCmd.AddCommand(abiCmd)

	abiCmd.PersistentFlags().StringVar(
		&abiArgs.yamlPath, "input", "abi.yaml", "Input file that contains the ABI yaml")
}
