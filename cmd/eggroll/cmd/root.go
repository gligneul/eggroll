// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootArgs struct {
	timeout int
}

var rootCmd = &cobra.Command{
	Use:   "eggroll",
	Short: "Command line tool for EggRoll",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(
		&rootArgs.timeout, "timeout", "t", 30, "Timeout in secs when executing a command")
}

func contextFromTimeout() (context.Context, context.CancelFunc) {
	timeout := time.Second * time.Duration(rootArgs.timeout)
	return context.WithTimeout(context.Background(), timeout)
}
