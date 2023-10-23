// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll"
	"github.com/holiman/uint256"
)

type Withdraw struct {
	Value *uint256.Int
}

type Honeypot struct {
	Balance *uint256.Int
}

func Codecs() []eggroll.Codec {
	return []eggroll.Codec{
		eggroll.NewJSONCodec[Withdraw](),
		eggroll.NewJSONCodec[Honeypot](),
	}
}
