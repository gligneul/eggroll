// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"log"
	"os"
	"runtime"
)

// Configuration for the DApp
type DAppConfig struct {
	RollupsEndpoint string
}

// Load the config from environment variables
func (c *DAppConfig) Load() {
	varName := "EGGROLL_ROLLUPS_HTTP_ENDPOINT"
	c.RollupsEndpoint = os.Getenv(varName)
	if c.RollupsEndpoint == "" {
		if runtime.GOARCH == "riscv64" {
			c.RollupsEndpoint = "http://localhost:5004"
		} else {
			c.RollupsEndpoint = "http://localhost:8080/host-runner"
		}
	}
	log.Printf("set %v=%v\n", varName, c.RollupsEndpoint)
}
