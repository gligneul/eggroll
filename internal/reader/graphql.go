// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package reader

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Read the rollups state by connecting to the rollups node GraphQL API.
type GraphQLReader struct {
	client graphql.Client
}

// Create a new GraphQL reader given the endpoint.
func NewGraphQLReader(endpoint string) *GraphQLReader {
	client := graphql.NewClient(endpoint, http.DefaultClient)
	return &GraphQLReader{
		client: client,
	}
}

// Get an input from the rollups node.
func (r *GraphQLReader) Input(ctx context.Context, index int) (*Input, error) {
	_ = `# @genqlient
	query getInput($inputIndex: Int!) {
	  input(index: $inputIndex) {
	    status
	    payload
	    msgSender
	    timestamp
	    blockNumber
	    vouchers {
	      edges {
		node {
		  index
		  destination
		  payload
		}
	      }
	    }
	    notices {
	      edges {
		node {
		  index
		  payload
		}
	      }
	    }
	    reports {
	      edges {
		node {
		  index
		  payload
		}
	      }
	    }
	  }
	}`

	resp, err := getInput(ctx, r.client, index)
	if err != nil {
		return nil, checkNotFound("input", err)
	}

	var input Input
	input.Index = index
	input.Status = resp.Input.Status

	input.Payload, err = hexutil.Decode(resp.Input.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	sender, err := hexutil.Decode(resp.Input.MsgSender)
	if err != nil {
		return nil, fmt.Errorf("failed to decode msgSender: %v", err)
	}
	input.Sender = common.Address(sender)

	input.BlockNumber, err = strconv.ParseInt(resp.Input.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block number: %v", err)
	}

	timestamp, err := strconv.ParseInt(resp.Input.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode timestmap: %v", err)
	}
	input.BlockTimestamp = time.Unix(timestamp, 0)

	for _, edge := range resp.Input.Vouchers.Edges {
		var voucher Voucher
		voucher.InputIndex = index
		voucher.OutputIndex = edge.Node.Index
		destination, err := hexutil.Decode(edge.Node.Destination)
		if err != nil {
			return nil, fmt.Errorf("failed to decode voucher destination: %v", err)
		}
		voucher.Destination = common.Address(destination)
		voucher.Payload, err = hexutil.Decode(edge.Node.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode voucher payload: %v", err)
		}
		input.Vouchers = append(input.Vouchers, voucher)
	}

	for _, edge := range resp.Input.Notices.Edges {
		var notice Notice
		notice.InputIndex = index
		notice.OutputIndex = edge.Node.Index
		notice.Payload, err = hexutil.Decode(edge.Node.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode notice payload: %v", err)
		}
		input.Notices = append(input.Notices, notice)
	}

	for _, edge := range resp.Input.Reports.Edges {
		var report Report
		report.InputIndex = index
		report.OutputIndex = edge.Node.Index
		report.Payload, err = hexutil.Decode(edge.Node.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode report payload: %v", err)
		}
		input.Reports = append(input.Reports, report)
	}

	return &input, nil
}

// Given the GraphQL error message, check whether the error should be NotFound.
func checkNotFound(typeName string, err error) error {
	if strings.HasSuffix(err.Error(), "not found\n") {
		return NotFound{typeName}
	}
	return err
}

//go:generate go run github.com/Khan/genqlient
