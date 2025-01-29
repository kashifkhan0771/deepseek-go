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
	url := fmt.Sprintf("%s%s", baseURL, chatCompletion)

	resp, err := c.doRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chat ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chat, nil
}

func (c *APIClient) ListModel() (*ListModelResponse, error) {
	url := fmt.Sprintf("%s%s", baseURL, listModel)

	resp, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var models ListModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &models, nil
}

// doRequest sends an HTTP request with the given method and payload.
func (c *APIClient) doRequest(method, url string, payload any) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		switch v := payload.(type) {
		case io.Reader:
			// Directly use an io.Reader (e.g., for files/streams)
			body = v
		default:
			// Marshal JSON for other types
			payloadBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request payload: %w", err)
			}
			body = bytes.NewReader(payloadBytes)
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil && !isStream(payload) {
		// Only set Content-Type for non-streaming JSON payloads
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close() // close body on error
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

// isStream checks if the payload is an io.Reader (streaming).
func isStream(payload any) bool {
	_, ok := payload.(io.Reader)
	return ok
}
