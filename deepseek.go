package deepseek

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient handles communication with the DeepSeek APIs.
type APIClient struct {
	// token can be obtained from: https://platform.deepseek.com/api_keys
	token  string
	client *http.Client
}

// NewClient initializes an API client.
func NewClient(token string) (*APIClient, error) {
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	return &APIClient{
		token: token,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

// ChatCompletion sends a request to the chat API and returns the response.
func (c *APIClient) ChatCompletion(payload ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if c.token == "" {
		return nil, errors.New("token cannot be empty")
	}

	url := fmt.Sprintf("%s%s", BaseURL, ChatComletion)

	resp, err := c.doRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	var chat ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chat, nil
}

// doRequest sends an HTTP request with the given method and payload.
func (c *APIClient) doRequest(method, url string, payload any) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
