// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll"

	"honeypot"
)

func main() {
	eggroll.Roll(&honeypot.Contract{})
}
