// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package rollups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Metadata from the input.
type Metadata struct {
	InputIndex     int
	Sender         common.Address
	BlockNumber    int64
	BlockTimestamp int64
}

// Represent an advance input from finish.
type AdvanceInput struct {
	Metadata *Metadata
	Payload  []byte
}

// Represent an inspect input from finish.
type InspectInput struct {
	Payload []byte
}

// Status when finishing a rollups request.
type FinishStatus int

const (
	FinishStatusAccept FinishStatus = iota
	FinishStatusReject
)

func (status FinishStatus) String() string {
	toString := map[FinishStatus]string{
		FinishStatusAccept: "accept",
		FinishStatusReject: "reject",
	}
	return toString[status]
}

// Implement the Rollups API using the Rollups HTTP server.
type RollupsHTTP struct {
	endpoint string
}

// Create a new Rollups HTTP client.
// Load the ROLLUP_HTTP_SERVER_URL from an environment variable.
func NewRollupsHTTP() *RollupsHTTP {
	endpoint := os.Getenv("ROLLUP_HTTP_SERVER_URL")
	if endpoint == "" {
		if runtime.GOARCH == "riscv64" {
			endpoint = "http://127.0.0.1:5004"
		} else {
			endpoint = "http://localhost:8080/host-runner"
		}
	}
	return &RollupsHTTP{endpoint}
}

// Send a post request and return the http response.
func (r *RollupsHTTP) sendPost(route string, data []byte) (*http.Response, error) {
	endpoint := r.endpoint + "/" + route
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}

	return resp, nil
}

// Send a voucher to the Rollups API. Return the voucher index.
func (r *RollupsHTTP) SendVoucher(destination common.Address, payload []byte) (int, error) {
	request := struct {
		Destination string `json:"destination"`
		Payload     string `json:"payload"`
	}{
		Destination: hexutil.Encode(destination[:]),
		Payload:     hexutil.Encode(payload),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize request: %v", err)
	}

	resp, err := r.sendPost("voucher", body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err = checkStatusOk(resp); err != nil {
		return 0, err
	}

	return parseOutputIndex(resp.Body)
}

// Send a notice to the Rollups API. Return the notice index.
func (r *RollupsHTTP) SendNotice(payload []byte) (int, error) {
	request := struct {
		Payload string `json:"payload"`
	}{
		Payload: hexutil.Encode(payload),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize request: %v", err)
	}

	resp, err := r.sendPost("notice", body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err = checkStatusOk(resp); err != nil {
		return 0, err
	}

	return parseOutputIndex(resp.Body)
}

// Send a report to the Rollups API.
func (r *RollupsHTTP) SendReport(payload []byte) error {
	request := struct {
		Payload string `json:"payload"`
	}{
		Payload: hexutil.Encode(payload),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %v", err)
	}

	resp, err := r.sendPost("report", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = checkStatusOk(resp); err != nil {
		return err
	}

	return nil
}

// Send a finish request to the Rollups API.
// If there is no error, return an AdvanceInput or an InspectInput.
func (r *RollupsHTTP) Finish(status FinishStatus) (any, error) {
	request := struct {
		Status string `json:"status"`
	}{
		Status: status.String(),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %v", err)
	}

	resp, err := r.sendPost("finish", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		// got StatusAccepted, trying again
		return r.Finish(status)
	}

	if err = checkStatusOk(resp); err != nil {
		return nil, err
	}

	var finishResp struct {
		RequestType string          `json:"request_type"`
		Data        json.RawMessage `json:"data"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&finishResp); err != nil {
		return nil, fmt.Errorf("failed to decode finish response: %v", err)
	}

	switch finishResp.RequestType {
	case "advance_state":
		return parseAdvanceInput(finishResp.Data)
	case "inspect_state":
		return parseInspectInput(finishResp.Data)
	default:
		return nil, fmt.Errorf("invalid request type: %v", finishResp.RequestType)
	}
}

func parseOutputIndex(r io.Reader) (int, error) {
	var outputResp struct {
		Index int `json:"index"`
	}
	if err := json.NewDecoder(r).Decode(&outputResp); err != nil {
		return 0, fmt.Errorf("failed to decode finish response: %v", err)
	}
	return outputResp.Index, nil
}

func parseAdvanceInput(data json.RawMessage) (any, error) {
	var advanceRequest struct {
		Payload  string `json:"payload"`
		Metadata struct {
			MsgSender   string `json:"msg_sender"`
			EpochIndex  int    `json:"epoch_index"`
			InputIndex  int    `json:"input_index"`
			BlockNumber int64  `json:"block_number"`
			Timestamp   int64  `json:"timestamp"`
		}
	}

	if err := json.Unmarshal(data, &advanceRequest); err != nil {
		return nil, fmt.Errorf("failed to decode advance request: %v", err)
	}

	payload, err := hexutil.Decode(advanceRequest.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode advance payload: %v", err)
	}

	sender, err := hexutil.Decode(advanceRequest.Metadata.MsgSender)
	if err != nil {
		return nil, fmt.Errorf("failed to decode advance metadata sender: %v", err)
	}

	metadata := &Metadata{
		InputIndex:     advanceRequest.Metadata.InputIndex,
		Sender:         common.Address(sender),
		BlockNumber:    advanceRequest.Metadata.BlockNumber,
		BlockTimestamp: advanceRequest.Metadata.Timestamp,
	}

	input := &AdvanceInput{
		Metadata: metadata,
		Payload:  payload,
	}

	return input, nil
}

func parseInspectInput(data json.RawMessage) (any, error) {
	var inspectRequest struct {
		Payload string `json:"payload"`
	}

	if err := json.Unmarshal(data, &inspectRequest); err != nil {
		return nil, fmt.Errorf("failed to decode advance request: %v", err)
	}

	payload, err := hexutil.Decode(inspectRequest.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode advance payload: %v", err)
	}

	input := &InspectInput{
		Payload: payload,
	}

	return input, nil
}

// Check the whether the status code is Ok, if not return an error.
func checkStatusOk(resp *http.Response) error {
	statusOk := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOk {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %v", err)
		}
		return fmt.Errorf("got invalid status %v: %v\n",
			resp.StatusCode, string(body))
	}
	return nil
}
