// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"github.com/gligneul/eggroll/pkg/eggroll"
)

type Append struct {
	Value string
}

type Clear struct{}

type TextBox struct {
	Value string
}

func Codecs() []eggroll.Codec {
	return []eggroll.Codec{
		eggroll.NewJSONCodec[Clear](),
		eggroll.NewJSONCodec[Append](),
		eggroll.NewJSONCodec[TextBox](),
	}
}
