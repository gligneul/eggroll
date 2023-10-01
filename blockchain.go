// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import "context"

// Communicate with the blockchain.
type blockchainClient interface {
	Send(ctx context.Context)
}
