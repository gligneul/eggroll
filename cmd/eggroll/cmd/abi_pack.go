// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"fmt"

	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/spf13/cobra"
)

var abiPackArgs struct {
	message string
}

var abiPackCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack ABI bindings",
	Example: `eggroll abi pack --message '
{
  "name": "log",
  "args": {
    "message": "hello"
  }
}'

output:

0x41304fac
0000000000000000000000000000000000000000000000000000000000000020
0000000000000000000000000000000000000000000000000000000000000005
68656c6c6f000000000000000000000000000000000000000000000000000000
`,
	Long: `Pack ABI bindings into JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		abiLoad()

		payload, err := eggtypes.PackFromJson([]byte(abiPackArgs.message))
		cobra.CheckErr(err)

		fmt.Printf("0x%x\n", payload[:4])
		for i := 4; i < len(payload); i += 32 {
			fmt.Printf("%x\n", payload[i:i+32])
		}
	},
}

func init() {
	abiCmd.AddCommand(abiPackCmd)

	abiPackCmd.Flags().StringVar(
		&abiPackArgs.message, "message", "", "JSON message that will be converted to hex")
	abiPackCmd.MarkFlagRequired("message")
}
