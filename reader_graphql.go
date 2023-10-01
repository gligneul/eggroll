// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Khan/genqlient/graphql"
)

// RollupsReader implementation that connects to the Rollups Node GraphQL API.
type GraphqlReader struct {
	ctx    context.Context
	client graphql.Client
}

func NewGraphqlReader(ctx context.Context, endpoint string) *GraphqlReader {
	client := graphql.NewClient(endpoint, http.DefaultClient)
	return &GraphqlReader{
		ctx:    ctx,
		client: client,
	}
}

func (r *GraphqlReader) Input(index int) (*Input, error) {
	_ = `# @genqlient
	  query getInput($inputIndex: Int!) {
	    input(index: $inputIndex) {
	      index
	      status
	      blockNumber
	    }
	  }
	`

	resp, err := getInput(r.ctx, r.client, index)
	if err != nil {
		return nil, err
	}

	blockNumber, err := strconv.ParseInt(resp.Input.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block number: %v", err)
	}

	input := &Input{
		Index:       resp.Input.Index,
		Status:      resp.Input.Status,
		BlockNumber: blockNumber,
	}
	return input, nil
}
