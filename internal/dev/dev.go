// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Implementation of eggroll dev mode.
// See cmd/eggroll/cmd/dev.go for mode details about de dev command.
package dev

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type RunOpts struct {
	DAppPath string
}

// Run the DApp in dev mode.
func Run(ctx context.Context, opts RunOpts) {
	checkDApp(opts.DAppPath)

	run := make(chan struct{}, 1)
	watcher := watcher{opts.DAppPath, run}
	runner := runner{opts.DAppPath, run}
	services := []service{
		startSignalListener,
		startSunodo,
		watcher.start,
		runner.start,
	}
	startServices(ctx, services)
}

// Check if there is Go main pkg in the current dir.
func checkDApp(path string) {
	cmd := exec.Command("go", "list", "-f", "{{.Name}}", path)
	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes[:len(outputBytes)-1]) // trim \n
	if err != nil {
		exit(output)
	}
	if output != "main" {
		exit("the working directory does not contain a main package")
	}
}

func exit(msg any) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}

func check(err any) {
	if err != nil {
		exit(err)
	}
}
