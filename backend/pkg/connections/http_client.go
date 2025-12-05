package connections

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPComputeClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPComputeClient(cfg Config) (*HTTPComputeClient, error) {
	timeout := 30
	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}

	return &HTTPComputeClient{
		baseURL: cfg.BaseURL,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

func (c *HTTPComputeClient) StartGame(ctx context.Context, req *StartGameRequest) (*StartGameResponse, error) {
	var resp StartGameResponse
	if err := c.call(ctx, "POST", "/start-game", req, &resp); err != nil {
		return nil, fmt.Errorf("start game failed: %w", err)
	}
	return &resp, nil
}

func (c *HTTPComputeClient) MakeGuess(ctx context.Context, req *GuessRequest) (*GuessResponse, error) {
	var resp GuessResponse
	if err := c.call(ctx, "POST", "/guess", req, &resp); err != nil {
		return nil, fmt.Errorf("make guess failed: %w", err)
	}
	return &resp, nil
}

func (c *HTTPComputeClient) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	var resp HealthResponse
	if err := c.call(ctx, "GET", "/health", nil, &resp); err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}
	return &resp, nil
}

func (c *HTTPComputeClient) call(ctx context.Context, method, path string, body, response interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *HTTPComputeClient) Close() error {
	return nil
}
