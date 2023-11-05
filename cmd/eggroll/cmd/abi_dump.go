// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gligneul/eggroll/internal/compiler"
	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/spf13/cobra"
)

var abiDumpArgs struct {
	dumpJson bool
}

var abiDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump ABI bindings",
	Example: `eggroll abi dump

output:
41304fac log(string)`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, err := os.Open(abiArgs.yamlPath)
		cobra.CheckErr(err)
		defer inputFile.Close()

		jsonAbi, err := compiler.Compile(inputFile)
		cobra.CheckErr(err)

		a, err := abi.JSON(strings.NewReader(jsonAbi))
		cobra.CheckErr(err)

		if abiDumpArgs.dumpJson {
			fmt.Print(jsonAbi)
		} else {
			abiDump(a)
		}
	},
}

// Add the ABI methods to Eggtypes to include the internal encodings.
func abiDump(a abi.ABI) {
	for _, method := range a.Methods {
		err := eggtypes.AddEncoding(eggtypes.Encoding{
			ID:        eggtypes.ID(method.ID),
			Name:      method.Name,
			Arguments: method.Inputs,
		})
		cobra.CheckErr(err)
	}
	var buffer bytes.Buffer
	encodings := eggtypes.GetEncodings()
	for _, encoding := range encodings {
		fmt.Fprintf(&buffer, "%x %v(",
			encoding.ID,
			encoding.Name,
		)
		for i, a := range encoding.Arguments {
			if i != 0 {
				buffer.WriteString(", ")
			}
			fmt.Fprintf(&buffer, "%v", a.Type)
		}
		buffer.WriteString(")\n")
	}
	fmt.Print(buffer.String())
}

func init() {
	abiCmd.AddCommand(abiDumpCmd)

	abiDumpCmd.Flags().BoolVar(
		&abiDumpArgs.dumpJson, "json", false, "If set, dump the JSON ABI")
}
