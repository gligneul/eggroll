// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/spf13/cobra"
)

var schemaEncodeArgs struct {
	kind     string
	args     string
	readable bool
}

var schemaEncodeCmd = &cobra.Command{
	Use:     "encode",
	Short:   "Encode message into bytes",
	Example: `eggroll schema encode --kind log --args '{"message": "hello"}'`,
	Long:    `Encode ABI bindings into JSON`,
	Run: func(cmd *cobra.Command, _args []string) {
		schemaLoad()

		args := make(map[string]any)
		err := json.Unmarshal([]byte(schemaEncodeArgs.args), &args)
		cobra.CheckErr(err)

		payload, err := eggtypes.EncodeFromMap(schemaEncodeArgs.kind, args)
		cobra.CheckErr(err)

		if schemaEncodeArgs.readable {
			fmt.Printf("  -4 %x\n", payload[:4])
			for i := 4; i < len(payload); i += 32 {
				fmt.Printf("%04x %x\n", i-4, payload[i:i+32])
			}
		} else {
			fmt.Printf("0x%x\n", payload)
		}
	},
}

func init() {
	schemaCmd.AddCommand(schemaEncodeCmd)

	schemaEncodeCmd.Flags().StringVar(
		&schemaEncodeArgs.kind, "kind", "", "Message kind")
	schemaEncodeCmd.MarkFlagRequired("kind")

	schemaEncodeCmd.Flags().StringVar(
		&schemaEncodeArgs.args, "args", "", "Message args encoded as JSON")
	schemaEncodeCmd.MarkFlagRequired("args")

	schemaEncodeCmd.Flags().BoolVarP(
		&schemaEncodeArgs.readable, "readable", "r", false, "If set, prints in human-readable format")
}
