// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"github.com/spf13/cobra"
)

// abiCmd represents the abi command
var abiCmd = &cobra.Command{
	Use:   "abi",
	Short: "Commands related to ABI encoding",
}

func init() {
	rootCmd.AddCommand(abiCmd)
}
