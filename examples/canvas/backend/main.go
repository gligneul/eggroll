// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/examples/canvas"
)

func advance(state eggroll.State) error {
	var input canvas.Input
	state.Input(&input)
	state.Report("got input %v", input)

	if input.X < 0 || input.X >= canvas.Width {
		return fmt.Errorf("invalid value for x axis")
	}
	if input.Y < 0 || input.Y >= canvas.Height {
		return fmt.Errorf("invalid value for y axis")
	}

	var c canvas.Canvas
	if state.Has("canvas") {
		state.Get("canvas", &c)
	}
	c[input.X][input.X] = input.Color
	state.Set("canvas", &c)

	return nil
}

func main() {
	eggroll.Roll(advance)
}
