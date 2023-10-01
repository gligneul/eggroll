// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"encoding/hex"
	"fmt"
)

// Encode bytes to hex string.
func encodeHex(payload []byte) string {
	return "0x" + hex.EncodeToString(payload)
}

// Decode hex string to bytes.
func decodeHex(payload string) ([]byte, error) {
	if len(payload) < 2 {
		return nil, fmt.Errorf("invalid hex string '%v'", payload)
	}
	return hex.DecodeString(payload[2:])
}
