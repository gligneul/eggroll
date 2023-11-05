// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/spf13/cobra"
)

var abiUnpackArgs struct {
	payload string
}

var abiUnpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack ABI bindings",
	Example: `eggroll abi unpack --payload 0x41304fac\
0000000000000000000000000000000000000000000000000000000000000020\
0000000000000000000000000000000000000000000000000000000000000005\
68656c6c6f000000000000000000000000000000000000000000000000000000

output:
{
  "name": "log",
  "args": {
    "message": "hello"
  }
}
`,
	Long: `Unpack ABI bindings into JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		abiLoad()

		payload, err := hexutil.Decode(abiUnpackArgs.payload)
		cobra.CheckErr(err)

		jsonData, err := eggtypes.UnpackToJson(payload)
		cobra.CheckErr(err)

		fmt.Println(string(jsonData))
	},
}

func init() {
	abiCmd.AddCommand(abiUnpackCmd)

	abiUnpackCmd.Flags().StringVar(
		&abiUnpackArgs.payload, "payload", "", "Payload in hex starting with 0x")
	abiUnpackCmd.MarkFlagRequired("payload")
}
