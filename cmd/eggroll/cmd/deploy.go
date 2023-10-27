// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"github.com/spf13/cobra"
)

var deployArgs struct {
	rpc string
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a test contract",
	Long: `
Deploy a contract for testing in a local Ethereum node`,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringVar(
		&deployArgs.rpc, "rpc", "http://localhost:8545", "Ethereum node rpc endpoint")
}
