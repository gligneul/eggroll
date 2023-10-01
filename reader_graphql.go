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
	client graphql.Client
}

func NewGraphqlReader(endpoint string) *GraphqlReader {
	client := graphql.NewClient(endpoint, http.DefaultClient)
	return &GraphqlReader{
		client: client,
	}
}

func (r *GraphqlReader) Input(ctx context.Context, index int) (*Input, error) {
	_ = `# @genqlient
	  query getInput($inputIndex: Int!) {
	    input(index: $inputIndex) {
	      status
	      blockNumber
	    }
	  }
	`

	resp, err := getInput(ctx, r.client, index)
	if err != nil {
		return nil, err
	}

	blockNumber, err := strconv.ParseInt(resp.Input.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block number: %v", err)
	}

	input := &Input{
		Index:       index,
		Status:      resp.Input.Status,
		BlockNumber: blockNumber,
	}
	return input, nil
}

func (r *GraphqlReader) Notice(ctx context.Context, inputIndex int, noticeIndex int) (*Notice, error) {
	_ = `# @genqlient
	  query getNotice($inputIndex: Int!, $noticeIndex: Int!) {
	    notice(noticeIndex: $noticeIndex, inputIndex: $inputIndex) {
	      payload
	    }
	  }
	`

	resp, err := getNotice(ctx, r.client, inputIndex, noticeIndex)
	if err != nil {
		return nil, err
	}

	payload, err := decodeHex(resp.Notice.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode notice payload: %v", err)
	}

	notice := &Notice{
		InputIndex:  inputIndex,
		NoticeIndex: noticeIndex,
		Payload:     payload,
	}

	return notice, nil
}

//go:generate go run github.com/Khan/genqlient
