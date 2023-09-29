// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Example DApp that stores a canvas in the blockchain.
// Anyone can send an input to change a pixel of the canvas.
package canvas

import "image/color"

const (
	Width  = 800
	Height = 600
)

// Canvas representation
type Canvas [Width][Height]color.RGBA

// Request a change to the canvas
type Input struct {
	X     int
	Y     int
	Color color.RGBA
}
