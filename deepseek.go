package deepseek

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	client = makeHTTPClient()
)

// DeepSeekConfig has the configurations required to call the APIs
type DeepSeekConfig struct {
	// Token can be obtained from: https://platform.deepseek.com/api_keys
	Token string
}

// NewDeepSeekConfig creates a DeepSeek config with token passed.
func NewDeepSeekConfig(token string) (*DeepSeekConfig, error) {
	if token == "" {
		return nil, errors.New("Token key cannot be empty")
	}

	return &DeepSeekConfig{
		Token: token,
	}, nil
}

// ChatCompletion return a model response for the given chat conversation.
func (d DeepSeekConfig) ChatCompletion(payload ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if d.Token == "" {
		return nil, errors.New("token cannot be empty in config")
	}

	url := fmt.Sprintf("%s%s", BaseURL, ChatComletion)

	// marshal the request payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := d.makeRequest(http.MethodGet, url, payloadBytes)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// decode JSON response body
	var chat ChatCompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&chat)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// makeRequest is a helper function for making API requests.
func (d DeepSeekConfig) makeRequest(method, url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.Token))

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func makeHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
