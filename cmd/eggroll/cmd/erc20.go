// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"fmt"
	"log"

	"github.com/gligneul/eggroll/pkg/eggeth"
	"github.com/spf13/cobra"
)

var erc20Cmd = &cobra.Command{
	Use:   "erc20",
	Short: "Deploy ERC20 contract",
	Long: `
Deploy an ERC20 contract for testing in a local Ethereum node`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := contextFromTimeout()
		defer cancel()
		address, err := eggeth.DeployTestERC20(ctx, deployArgs.rpc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(address)
	},
}

func init() {
	deployCmd.AddCommand(erc20Cmd)
}
