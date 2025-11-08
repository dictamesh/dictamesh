// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient is a wrapper around http.Client with additional features
// like retry logic, rate limiting, and request/response logging
type HTTPClient struct {
	client        *http.Client
	retryConfig   *RetryConfig
	rateLimiter   RateLimiter
	requestLogger RequestLogger
	baseURL       string
	headers       map[string]string
}

// HTTPClientConfig represents the configuration for an HTTP client
type HTTPClientConfig struct {
	// BaseURL is the base URL for all requests
	BaseURL string

	// Timeout is the request timeout
	Timeout time.Duration

	// RetryConfig configures retry behavior
	RetryConfig *RetryConfig

	// RateLimiter limits the rate of requests
	RateLimiter RateLimiter

	// RequestLogger logs requests and responses
	RequestLogger RequestLogger

	// Headers are default headers to include in all requests
	Headers map[string]string

	// Transport is the HTTP transport to use
	Transport http.RoundTripper
}

// RetryConfig configures retry behavior
type RetryConfig struct {
	// MaxRetries is the maximum number of retries
	MaxRetries int

	// InitialBackoff is the initial backoff duration
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration
	MaxBackoff time.Duration

	// BackoffMultiplier is the backoff multiplier for exponential backoff
	BackoffMultiplier float64

	// RetryableStatusCodes are HTTP status codes that should trigger a retry
	RetryableStatusCodes []int
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    1 * time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
		RetryableStatusCodes: []int{
			http.StatusTooManyRequests,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusBadGateway,
		},
	}
}

// RateLimiter is an interface for rate limiting
type RateLimiter interface {
	// Wait blocks until the rate limiter allows a request
	Wait(ctx context.Context) error

	// Allow checks if a request is allowed without blocking
	Allow() bool
}

// RequestLogger is an interface for logging requests
type RequestLogger interface {
	// LogRequest logs an HTTP request
	LogRequest(req *http.Request)

	// LogResponse logs an HTTP response
	LogResponse(req *http.Request, resp *http.Response, duration time.Duration, err error)
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(config *HTTPClientConfig) *HTTPClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.RetryConfig == nil {
		config.RetryConfig = DefaultRetryConfig()
	}

	client := &http.Client{
		Timeout:   config.Timeout,
		Transport: config.Transport,
	}

	return &HTTPClient{
		client:        client,
		retryConfig:   config.RetryConfig,
		rateLimiter:   config.RateLimiter,
		requestLogger: config.RequestLogger,
		baseURL:       config.BaseURL,
		headers:       config.Headers,
	}
}

// Do executes an HTTP request with retry logic
func (c *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Apply rate limiting if configured
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter wait failed: %w", err)
		}
	}

	// Apply default headers
	for key, value := range c.headers {
		if req.Header.Get(key) == "" {
			req.Header.Set(key, value)
		}
	}

	var resp *http.Response
	var err error
	var lastErr error

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate backoff duration
			backoff := c.calculateBackoff(attempt)

			// Wait for backoff duration
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Log request
		if c.requestLogger != nil {
			c.requestLogger.LogRequest(req)
		}

		startTime := time.Now()

		// Clone the request for retry attempts
		clonedReq := req.Clone(ctx)

		// Execute the request
		resp, err = c.client.Do(clonedReq)
		duration := time.Since(startTime)

		// Log response
		if c.requestLogger != nil {
			c.requestLogger.LogResponse(req, resp, duration, err)
		}

		// Check if we should retry
		if err != nil {
			lastErr = err
			continue
		}

		// Check status code
		if !c.isRetryableStatusCode(resp.StatusCode) {
			return resp, nil
		}

		// Store response for potential retry
		lastErr = fmt.Errorf("request failed with status %d", resp.StatusCode)
		resp.Body.Close()
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return resp, err
}

// Get performs a GET request
func (c *HTTPClient) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(ctx, req)
}

// Post performs a POST request
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(path)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(ctx, req)
}

// Put performs a PUT request
func (c *HTTPClient) Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(path)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(ctx, req)
}

// Patch performs a PATCH request
func (c *HTTPClient) Patch(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(path)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(ctx, req)
}

// Delete performs a DELETE request
func (c *HTTPClient) Delete(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(ctx, req)
}

// buildURL constructs the full URL from the base URL and path
func (c *HTTPClient) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}
	return c.baseURL + path
}

// calculateBackoff calculates the backoff duration for a retry attempt
func (c *HTTPClient) calculateBackoff(attempt int) time.Duration {
	backoff := float64(c.retryConfig.InitialBackoff) * float64(attempt) * c.retryConfig.BackoffMultiplier
	duration := time.Duration(backoff)

	if duration > c.retryConfig.MaxBackoff {
		return c.retryConfig.MaxBackoff
	}

	return duration
}

// isRetryableStatusCode checks if a status code should trigger a retry
func (c *HTTPClient) isRetryableStatusCode(statusCode int) bool {
	for _, code := range c.retryConfig.RetryableStatusCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

// ParseJSONResponse parses a JSON response
func ParseJSONResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return HTTPErrorToAdapterError(resp.StatusCode, string(body), nil)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// ReadResponseBody reads the full response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
