// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package reader

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type CompletionStatus string

const (
	CompletionStatusUnprocessed                CompletionStatus = "UNPROCESSED"
	CompletionStatusAccepted                   CompletionStatus = "ACCEPTED"
	CompletionStatusRejected                   CompletionStatus = "REJECTED"
	CompletionStatusException                  CompletionStatus = "EXCEPTION"
	CompletionStatusMachineHalted              CompletionStatus = "MACHINE_HALTED"
	CompletionStatusCycleLimitExceeded         CompletionStatus = "CYCLE_LIMIT_EXCEEDED"
	CompletionStatusTimeLimitExceeded          CompletionStatus = "TIME_LIMIT_EXCEEDED"
	CompletionStatusPayloadLengthLimitExceeded CompletionStatus = "PAYLOAD_LENGTH_LIMIT_EXCEEDED"
)

// Rollups input from the Reader API.
type Input struct {
	Index       int
	Status      CompletionStatus
	BlockNumber int64
}

// Rollups notice from the Reader API.
type Notice struct {
	InputIndex  int
	NoticeIndex int
	Payload     []byte
}

// Rollups report from the Reader API.
type Report struct {
	InputIndex  int
	ReportIndex int
	Payload     []byte
}

// Result of a paginated query.
type Page[T any] struct {
	Nodes           []T
	TotalCount      int
	StartCursor     string
	EndCursor       string
	HasNextPage     bool
	HasPreviousPage bool
}

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
func (r *GraphQLReader) Input(ctx context.Context, index int) (*Input, error) {
	_ = `# @genqlient
	query getInput($inputIndex: Int!) {
	  input(index: $inputIndex) {
	    blockNumber
	    reports {
	      totalCount
	    }
	    notices {
	      totalCount
	    }
	  }
	}`

	resp, err := getInput(ctx, r.client, index)
	if err != nil {
		return nil, checkNotFound("input", err)
	}

	blockNumber, err := strconv.ParseInt(resp.Input.BlockNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block number: %v", err)
	}

	status := CompletionStatusUnprocessed
	if resp.Input.Reports.TotalCount != 0 {
		if resp.Input.Notices.TotalCount != 0 {
			status = CompletionStatusAccepted
		} else {
			status = CompletionStatusRejected
		}
	}

	input := &Input{
		Index:       index,
		Status:      status,
		BlockNumber: blockNumber,
	}
	return input, nil
}

// Get a notice from the rollups node.
// If the notice doesn't exist, return NotFound error.
func (r *GraphQLReader) Notice(ctx context.Context, inputIndex int, noticeIndex int) (*Notice, error) {
	_ = `# @genqlient
	query getNotice($inputIndex: Int!, $noticeIndex: Int!) {
	  notice(noticeIndex: $noticeIndex, inputIndex: $inputIndex) {
	    payload
	  }
	}`

	resp, err := getNotice(ctx, r.client, inputIndex, noticeIndex)
	if err != nil {
		return nil, checkNotFound("notice", err)
	}

	payload, err := hexutil.Decode(resp.Notice.Payload)
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

// Get a report from the rollups node.
// If the report doesn't exist, return NotFound error.
func (r *GraphQLReader) Report(ctx context.Context, inputIndex int, reportIndex int) (*Report, error) {
	_ = `# @genqlient
	query getReport($inputIndex: Int!, $reportIndex: Int!) {
	  report(reportIndex: $reportIndex, inputIndex: $inputIndex) {
	    payload
	  }
	}`

	resp, err := getReport(ctx, r.client, inputIndex, reportIndex)
	if err != nil {
		return nil, checkNotFound("report", err)
	}

	payload, err := hexutil.Decode(resp.Report.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode report payload: %v", err)
	}

	report := &Report{
		InputIndex:  inputIndex,
		ReportIndex: reportIndex,
		Payload:     payload,
	}

	return report, nil
}

// Get a page of reports from the rollups node.
func (r *GraphQLReader) LastReports(ctx context.Context, last int) (*Page[Report], error) {
	_ = `# @genqlient
	query getLastReports($last: Int) {
	  reports(last: $last) {
	    totalCount
	    pageInfo {
	      startCursor
	      endCursor
	      hasNextPage
	      hasPreviousPage
	    }
	    edges {
	      node {
	        index
	        input {
	          index
	        }
	        payload
	      }
	      cursor
	    }
	  }
	}`

	resp, err := getLastReports(ctx, r.client, last)
	if err != nil {
		return nil, err
	}

	var reports []Report
	for _, edge := range resp.Reports.Edges {
		payload, err := hexutil.Decode(edge.Node.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode report payload %v", err)
		}
		reports = append(reports, Report{
			InputIndex:  edge.Node.Input.Index,
			ReportIndex: edge.Node.Index,
			Payload:     payload,
		})
	}

	page := &Page[Report]{
		Nodes:           reports,
		TotalCount:      resp.Reports.TotalCount,
		StartCursor:     resp.Reports.PageInfo.StartCursor,
		EndCursor:       resp.Reports.PageInfo.EndCursor,
		HasNextPage:     resp.Reports.PageInfo.HasNextPage,
		HasPreviousPage: resp.Reports.PageInfo.HasPreviousPage,
	}

	return page, nil
}

//go:generate go run github.com/Khan/genqlient
