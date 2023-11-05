// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggtypes

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestLogABI(t *testing.T) {
	if bytes.Compare(LogID[:], common.Hex2Bytes("41304fac")) != 0 {
		t.Fatalf("wrong log id: %x", LogID[:])
	}

	log := &Log{
		Message: "hello",
	}
	expectedData := common.Hex2Bytes("41304fac" +
		// offset
		"0000000000000000000000000000000000000000000000000000000000000020" +
		// num bytes
		"0000000000000000000000000000000000000000000000000000000000000005" +
		// value
		"68656c6c6f000000000000000000000000000000000000000000000000000000")

	// Test pack
	packData := log.Pack()
	if !bytes.Equal(packData, expectedData) {
		t.Fatalf("wrong pack return; got %x", packData)
	}

	// Test unpack
	value, err := Unpack(packData)
	if err != nil {
		t.Fatalf("failed to decode log: %v", err)
	}
	if value.(Log).Message != "hello" {
		t.Fatalf("wrong payload")
	}
}
