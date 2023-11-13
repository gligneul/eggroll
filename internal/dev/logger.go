// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package dev

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
}

type logWritter struct {
	prefix string
}

func (w logWritter) Write(bytes []byte) (int, error) {
	logger.Print(w.prefix, ": ", string(bytes))
	return len(bytes), nil
}
