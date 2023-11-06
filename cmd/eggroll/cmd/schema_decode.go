// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/spf13/cobra"
)

var schemaDecodeArgs struct {
	payload string
}

var schemaDecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode ABI bindings",
	Example: `eggroll schema decode --payload 0x41304fac\
0000000000000000000000000000000000000000000000000000000000000020\
0000000000000000000000000000000000000000000000000000000000000005\
68656c6c6f000000000000000000000000000000000000000000000000000000`,
	Long: `Decode ABI bindings into JSON`,
	Run: func(cmd *cobra.Command, _args []string) {
		schemaLoad()

		payload, err := hexutil.Decode(schemaDecodeArgs.payload)
		cobra.CheckErr(err)

		args := make(map[string]any)
		kind, err := eggtypes.DecodeIntoMap(args, payload)
		cobra.CheckErr(err)

		jsonArgs, err := json.MarshalIndent(args, "", "  ")
		cobra.CheckErr(err)
		fmt.Printf("%v %v\n", kind, string(jsonArgs))
	},
}

func init() {
	schemaCmd.AddCommand(schemaDecodeCmd)

	schemaDecodeCmd.Flags().StringVar(
		&schemaDecodeArgs.payload, "payload", "", "Payload in hex starting with 0x")
	schemaDecodeCmd.MarkFlagRequired("payload")
}
