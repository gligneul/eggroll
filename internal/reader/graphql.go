// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package reader

//go:generate go run github.com/Khan/genqlient

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
	"github.com/gligneul/eggroll/eggtypes"
)

// Error when an object is not found.
type NotFound struct {
	typeName string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("%v not found", e.typeName)
}

// Given the GraphQL error message, check whether the error should be NotFound.
func checkNotFound(typeName string, err error) error {
	if strings.HasSuffix(err.Error(), "not found\n") {
		return NotFound{typeName}
	}
	return err
}

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
func (r *GraphQLReader) AdvanceResult(ctx context.Context, index int) (
	*eggtypes.AdvanceResult, error) {

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

	status, err := convertAdvanceStatus(resp.Input.Status)
	if err != nil {
		return nil, err
	}

	var reports []eggtypes.Report
	for _, edge := range resp.Input.Reports.Edges {
		var report eggtypes.Report
		report.InputIndex = index
		report.OutputIndex = edge.Node.Index
		report.Payload, err = hexutil.Decode(edge.Node.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode report payload: %v", err)
		}
		reports = append(reports, report)
	}

	payload, err := hexutil.Decode(resp.Input.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	sender, err := hexutil.Decode(resp.Input.MsgSender)
	if err != nil {
		return nil, fmt.Errorf("failed to decode msgSender: %v", err)
	}

	blockNumber, err := strconv.ParseInt(resp.Input.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block number: %v", err)
	}

	blockTimestamp, err := strconv.ParseInt(resp.Input.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode timestmap: %v", err)
	}

	var vouchers []eggtypes.Voucher
	for _, edge := range resp.Input.Vouchers.Edges {
		var voucher eggtypes.Voucher
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
		vouchers = append(vouchers, voucher)
	}

	var notices []eggtypes.Notice
	for _, edge := range resp.Input.Notices.Edges {
		var notice eggtypes.Notice
		notice.InputIndex = index
		notice.OutputIndex = edge.Node.Index
		notice.Payload, err = hexutil.Decode(edge.Node.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode notice payload: %v", err)
		}
		notices = append(notices, notice)
	}

	result := &eggtypes.AdvanceResult{
		Result: eggtypes.Result{
			Status:  status,
			Reports: reports,
		},
		Index:          index,
		Payload:        payload,
		Sender:         common.Address(sender),
		BlockNumber:    blockNumber,
		BlockTimestamp: time.Unix(blockTimestamp, 0),
		Vouchers:       vouchers,
		Notices:        notices,
	}

	return result, nil
}

func convertAdvanceStatus(s CompletionStatus) (eggtypes.CompletionStatus, error) {
	statusMap := map[CompletionStatus]eggtypes.CompletionStatus{
		CompletionStatusUnprocessed:                eggtypes.CompletionStatusUnprocessed,
		CompletionStatusAccepted:                   eggtypes.CompletionStatusAccepted,
		CompletionStatusRejected:                   eggtypes.CompletionStatusRejected,
		CompletionStatusException:                  eggtypes.CompletionStatusException,
		CompletionStatusMachineHalted:              eggtypes.CompletionStatusMachineHalted,
		CompletionStatusCycleLimitExceeded:         eggtypes.CompletionStatusCycleLimitExceeded,
		CompletionStatusTimeLimitExceeded:          eggtypes.CompletionStatusTimeLimitExceeded,
		CompletionStatusPayloadLengthLimitExceeded: eggtypes.CompletionStatusPayloadLengthLimitExceeded,
	}
	status, ok := statusMap[s]
	if !ok {
		return status, fmt.Errorf("invalid completion status: %v", s)
	}
	return status, nil
}
