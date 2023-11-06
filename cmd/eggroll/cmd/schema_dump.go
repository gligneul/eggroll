// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"bytes"
	"fmt"

	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/spf13/cobra"
)

var schemaDumpArgs struct {
	dumpJson bool
}

var schemaDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump the schema",
	Long: `Dump the schema based on the Yaml file.
This command also dumps the EggRoll internal message schemas, such as log.`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonAbi := schemaLoad()
		if schemaDumpArgs.dumpJson {
			fmt.Print(jsonAbi)
		} else {
			schemaDump()
		}
	},
}

func schemaDump() {
	var buffer bytes.Buffer
	schemas := eggtypes.GetSchemas()
	for _, schema := range schemas {
		fmt.Fprintf(&buffer, "%x %v(",
			schema.ID,
			schema.Kind,
		)
		for i, arg := range schema.Arguments {
			if i != 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(arg.Type.String())
		}
		buffer.WriteString(")\n")
	}
	fmt.Print(buffer.String())
}

func init() {
	schemaCmd.AddCommand(schemaDumpCmd)

	schemaDumpCmd.Flags().BoolVar(
		&schemaDumpArgs.dumpJson, "json", false, "If set, dump the JSON ABI")
}
