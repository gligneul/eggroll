// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// RollupsReader implementation that connects to the Rollups Node GraphQL API.
type GraphqlReader struct {
	Endpoint string
}

func (r *GraphqlReader) Input(index int) (*Input, error) {
	query := `query ($inputIndex: Int!) {
		input(index: $inputIndex) {
			index
			status
			blockNumber
		}
	}`

	variables := make(map[string]any)
	variables["inputIndex"] = index

	reqData := make(map[string]any)
	reqData["query"] = query
	reqData["variables"] = variables

	dataJson, err := json.Marshal(&reqData)
	if err != nil {
		log.Fatalf("failed to encode json: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, r.Endpoint, bytes.NewBuffer(dataJson))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

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
		return nil, fmt.Errorf("got invalid status %v: %v\n",
			resp.StatusCode, string(body))
	}

	graphqlResp := make(map[string]json.RawMessage)
	if err = json.NewDecoder(resp.Body).Decode(&graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode GraphQL response: %v", err)
	}

	if rawErrors, ok := graphqlResp["errors"]; ok {
		return nil, handleGraphqlErrors(rawErrors)
	}

	rawData, ok := graphqlResp["data"]
	if !ok {
		return nil, fmt.Errorf("graphql: data not found")
	}

	log.Printf(">> %v\n", string(rawData))

	var data struct {
		Input *Input
	}

	if err = json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to decode GraphQL response: %v", err)
	}

	return data.Input, nil
}

// Format the JSON with graphql errors into a error message
func handleGraphqlErrors(rawErrors json.RawMessage) error {
	var graphqlErrors []struct {
		Message string
	}

	if err := json.Unmarshal(rawErrors, &graphqlErrors); err != nil {
		return fmt.Errorf("failed to decode GraphQL errors: %v", err)
	}

	var sb strings.Builder
	sb.WriteString("graphql: ")
	for i, graphqlErr := range graphqlErrors {
		if i != 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(graphqlErr.Message)
	}
	return errors.New(sb.String())
}
