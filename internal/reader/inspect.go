// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Client for the reader node inspect API.
type InspectClient struct {
	endpoint string
}

// Create a new inspect client.
// The endpoint is for the reader node inspect HTTP API.
func NewInspectClient(endpoint string) *InspectClient {
	return &InspectClient{
		endpoint: endpoint,
	}
}

// Send a inspect request with the given payload.
func (c *InspectClient) Inspect(ctx context.Context, payload []byte) (*InspectResult, error) {

	// Prepare the request
	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req = req.WithContext(ctx)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %v", err)
		}

		if resp.StatusCode == http.StatusBadRequest &&
			strings.Contains(string(body), "concurrent call in session") {

			// This is an unpredictable error on the inspect
			// server, so we wait for a bit and retry
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(500 * time.Millisecond):
				return c.Inspect(ctx, payload)
			}
		}

		msg := "inspect error (status %v): %v"
		return nil, fmt.Errorf(msg, resp.StatusCode, string(body))
	}

	// Decode the json response
	var jsonResponse struct {
		Status  string `json:"status"`
		Reports []struct {
			Payload string `json:"payload"`
		} `json:"reports"`
		ProcessedInputCount int
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&jsonResponse)

	// Build the result
	status, err := convertInspectStatus(jsonResponse.Status)
	if err != nil {
		return nil, err
	}
	var reports []Report
	for _, jsonReport := range jsonResponse.Reports {
		payload, err := hexutil.Decode(jsonReport.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode report payload: %v", err)
		}
		report := Report{
			Payload: payload,
		}
		reports = append(reports, report)
	}
	result := &InspectResult{
		Status:              status,
		Reports:             reports,
		ProcessedInputCount: jsonResponse.ProcessedInputCount,
	}
	return result, nil
}

func convertInspectStatus(str string) (CompletionStatus, error) {
	statusMap := map[string]CompletionStatus{
		"Unprocessed":                CompletionStatusUnprocessed,
		"Accepted":                   CompletionStatusAccepted,
		"Rejected":                   CompletionStatusRejected,
		"Exception":                  CompletionStatusException,
		"Machine_halted":             CompletionStatusMachineHalted,
		"CycleLimitExceeded":         CompletionStatusCycleLimitExceeded,
		"TimeLimitExceeded":          CompletionStatusTimeLimitExceeded,
		"PayloadLengthLimitExceeded": CompletionStatusPayloadLengthLimitExceeded,
	}
	status, ok := statusMap[str]
	if !ok {
		return status, fmt.Errorf("invalid completion status: %v", str)
	}
	return status, nil
}
