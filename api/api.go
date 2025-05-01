package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"library-client/variables"
	"net/http"
	"time"
)

// var ApiBaseUrl = variables.ApiBaseUrl

var Client = &http.Client{
	Timeout: 10 * time.Second,
}

func RefreshToken() (*http.Response, error) {
	var Headers = map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer " + variables.RefreshToken,
	}

	endpoint := "users/refresh-token"
	var reqBody io.Reader

	ApiBaseUrl := variables.ApiBaseUrl
	fmt.Println("Loading data from server -->", ApiBaseUrl+endpoint)

	req, err := http.NewRequest(http.MethodGet, ApiBaseUrl+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// fmt.Println("Headers: ", Headers)
	for key, value := range Headers {
		req.Header.Set(key, value)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

func SendRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var Headers = map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer " + variables.AccessToken,
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	ApiBaseUrl := variables.ApiBaseUrl
	fmt.Println("Loading data from server -->", ApiBaseUrl+endpoint)

	req, err := http.NewRequest(method, ApiBaseUrl+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// fmt.Println("Headers: ", Headers)
	for key, value := range Headers {
		req.Header.Set(key, value)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}
