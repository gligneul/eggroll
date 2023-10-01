// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

// Shared types for the TextBox DApp.
package textbox

// Append a value to the text box.
type InputAppend struct {
	Value string
}

// Clear the text box.
type InputClear struct {
}

// Text box shared state.
type State struct {
	TextBox string
}
