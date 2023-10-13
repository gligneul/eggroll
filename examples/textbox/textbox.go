// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package textbox

import "github.com/gligneul/eggroll"

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
