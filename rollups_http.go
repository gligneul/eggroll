// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Implement the Rollups API using the Rollups HTTP server.
type rollupsHttpApi struct {
	endpoint string
}

// Send a post request and return the http response.
func (r *rollupsHttpApi) sendPost(route string, data []byte) (*http.Response, error) {
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

func (r *rollupsHttpApi) sendVoucher(destination Address, payload []byte) error {
	request := struct {
		Destination string `json:"destination"`
		Payload     string `json:"payload"`
	}{
		Destination: hexutil.Encode(destination[:]),
		Payload:     hexutil.Encode(payload),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %v", err)
	}

	resp, err := r.sendPost("voucher", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = checkStatusOk(resp); err != nil {
		return err
	}

	return nil
}

func (r *rollupsHttpApi) sendNotice(payload []byte) error {
	request := struct {
		Payload string `json:"payload"`
	}{
		Payload: hexutil.Encode(payload),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %v", err)
	}

	resp, err := r.sendPost("notice", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = checkStatusOk(resp); err != nil {
		return err
	}

	return nil
}

func (r *rollupsHttpApi) sendReport(payload []byte) error {
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

func (r *rollupsHttpApi) finish(status finishStatus) ([]byte, *Metadata, error) {
	request := struct {
		Status string `json:"status"`
	}{
		Status: status.String(),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize request: %v", err)
	}

	resp, err := r.sendPost("finish", body)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		// got StatusAccepted, trying again
		return r.finish(status)
	}

	if err = checkStatusOk(resp); err != nil {
		return nil, nil, err
	}

	var finishResp struct {
		RequestType string          `json:"request_type"`
		Data        json.RawMessage `json:"data"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&finishResp); err != nil {
		return nil, nil, fmt.Errorf("failed to decode finish response: %v", err)
	}

	if finishResp.RequestType != "advance_state" {
		log.Printf("rejecting %v", finishResp.RequestType)
		return r.finish(statusReject)
	}

	var advanceRequest struct {
		Payload  string `json:"payload"`
		Metadata struct {
			MsgSender   string `json:"msg_sender"`
			EpochIndex  int64  `json:"epoch_index"`
			InputIndex  int64  `json:"input_index"`
			BlockNumber int64  `json:"block_number"`
			Timestamp   int64  `json:"timestamp"`
		}
	}
	if err = json.Unmarshal(finishResp.Data, &advanceRequest); err != nil {
		return nil, nil, fmt.Errorf("failed to decode advance request: %v", err)
	}

	payload, err := hexutil.Decode(advanceRequest.Payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode advance payload: %v", err)
	}

	sender, err := hexutil.Decode(advanceRequest.Metadata.MsgSender)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode advance metadata sender: %v", err)
	}

	metadata := &Metadata{
		Sender:         Address(sender),
		BlockNumber:    advanceRequest.Metadata.BlockNumber,
		BlockTimestamp: advanceRequest.Metadata.Timestamp,
	}

	return payload, metadata, nil
}

// Check the whether the status code is Ok, if not return an error.
func checkStatusOk(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %v", err)
		}
		return fmt.Errorf("got invalid status %v: %v\n",
			resp.StatusCode, string(body))
	}
	return nil
}
