// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"os"
)

const prefix = "EGGROLL_"

func loadVar(varName string, defaultValue string) string {
	varName = prefix + varName
	value := os.Getenv(varName)
	if value == "" {
		value = defaultValue
	}
	return value
}
