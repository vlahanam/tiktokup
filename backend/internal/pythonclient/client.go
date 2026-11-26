package pythonclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tiktokup/backend/pkg/response"
)

// Client communicates with the internal Dockerized Python FastAPI service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// ServiceInfo captures runtime metadata returned by the Python service /info endpoint.
type ServiceInfo struct {
	Service       string `json:"service"`
	PythonVersion string `json:"python_version"`
}

// New creates a Client targeting baseURL (e.g. "http://python:8000").
func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

// Info queries GET /info on the Python service and propagates X-Request-ID.
func (c *Client) Info(ctx context.Context) (*ServiceInfo, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("python service url is not configured")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/info", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if reqID := response.RequestIDFromContext(ctx); reqID != "" {
		req.Header.Set("X-Request-ID", reqID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python service returned status %d", resp.StatusCode)
	}

	var envelope struct {
		Success bool        `json:"success"`
		Data    ServiceInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}
	if !envelope.Success {
		return nil, fmt.Errorf("python service reported unsuccessful response")
	}
	return &envelope.Data, nil
}
