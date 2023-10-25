// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggtypes

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestLogABI(t *testing.T) {
	if bytes.Compare(LogID[:], common.Hex2Bytes("cf34ef53")) != 0 {
		t.Fatalf("wrong log id: %x", LogID[:])
	}

	log := &Log{
		Message: "hello",
	}
	data := common.Hex2Bytes("cf34ef530000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000568656c6c6f000000000000000000000000000000000000000000000000000000")

	// Test pack
	if !bytes.Equal(log.Pack(), data) {
		t.Fatalf("pack failed")
	}

	// Test unpack
	value, err := Unpack(data)
	if err != nil {
		t.Fatalf("failed to decode log: %v", err)
	}
	if value.(Log).Message != "hello" {
		t.Fatalf("wrong payload")
	}
}
