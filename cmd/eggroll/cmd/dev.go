// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"github.com/gligneul/eggroll/internal/dev"
	"github.com/spf13/cobra"
)

var devArgs dev.RunOpts

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start the DApp contract in dev mode",
	Long: `This command starts sunodo with no-backend mode, and the DApp in the current directory.
If the DApp exits, the command will restart it automatically.
If there is a file change in the DApp directory, the command will restart the DApp.

This command captures the standard error from sunodo and save it to the file sunodo-stderr.log.`,
	Run: func(cmd *cobra.Command, args []string) {
		dev.Run(cmd.Context(), devArgs)
	},
}

func init() {
	rootCmd.AddCommand(devCmd)

	devCmd.Flags().StringVarP(
		&devArgs.DAppPath, "path", "p", ".", "path to the DApp directory")
}
