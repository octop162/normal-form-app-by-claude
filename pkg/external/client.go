// Package external provides HTTP client functionality for external API integrations.
package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	defaultTimeout     = 30 * time.Second
	defaultMaxRetries  = 3
	defaultRetryDelay  = 1 * time.Second
	contentTypeJSON    = "application/json"
	headerContentType  = "Content-Type"
	headerUserAgent    = "User-Agent"
	userAgentValue     = "normal-form-app/1.0"
)

// HTTPClient defines the interface for HTTP operations
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a configurable HTTP client for external API calls
type Client struct {
	httpClient HTTPClient
	baseURL    string
	timeout    time.Duration
	maxRetries int
	retryDelay time.Duration
	log        *logger.Logger
}

// Config holds configuration for the external API client
type Config struct {
	BaseURL    string        `json:"base_url"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`
}

// NewClient creates a new external API client with the provided configuration
func NewClient(config *Config, log *logger.Logger) *Client {
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = defaultMaxRetries
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = defaultRetryDelay
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		timeout:    config.Timeout,
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryDelay,
		log:        log,
	}
}

// APIResponse represents a standard external API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PostJSON performs a POST request with JSON payload and returns the response
func (c *Client) PostJSON(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	url := c.baseURL + endpoint
	
	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		c.log.WithError(err).WithField("endpoint", endpoint).Error("Failed to marshal request payload")
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			c.log.WithField("attempt", attempt).WithField("endpoint", endpoint).Info("Retrying API call")
			time.Sleep(c.retryDelay)
		}

		// Create HTTP request
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		// Set headers
		req.Header.Set(headerContentType, contentTypeJSON)
		req.Header.Set(headerUserAgent, userAgentValue)

		// Execute request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			c.log.WithError(err).WithField("endpoint", endpoint).WithField("attempt", attempt).Warn("HTTP request failed")
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			continue
		}

		// Process response
		err = c.processResponse(resp, result)
		if err != nil {
			c.log.WithError(err).WithField("endpoint", endpoint).WithField("status", resp.StatusCode).Warn("Failed to process response")
			lastErr = err
			
			// Don't retry on client errors (4xx)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				break
			}
			continue
		}

		// Success
		c.log.WithField("endpoint", endpoint).WithField("attempt", attempt).Debug("API call successful")
		return nil
	}

	c.log.WithError(lastErr).WithField("endpoint", endpoint).WithField("max_retries", c.maxRetries).Error("API call failed after all retries")
	return fmt.Errorf("API call failed after %d retries: %w", c.maxRetries, lastErr)
}

// GetJSON performs a GET request and returns the response
func (c *Client) GetJSON(ctx context.Context, endpoint string, result interface{}) error {
	url := c.baseURL + endpoint

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			c.log.WithField("attempt", attempt).WithField("endpoint", endpoint).Info("Retrying API call")
			time.Sleep(c.retryDelay)
		}

		// Create HTTP request
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		// Set headers
		req.Header.Set(headerUserAgent, userAgentValue)

		// Execute request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			c.log.WithError(err).WithField("endpoint", endpoint).WithField("attempt", attempt).Warn("HTTP request failed")
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			continue
		}

		// Process response
		err = c.processResponse(resp, result)
		if err != nil {
			c.log.WithError(err).WithField("endpoint", endpoint).WithField("status", resp.StatusCode).Warn("Failed to process response")
			lastErr = err
			
			// Don't retry on client errors (4xx)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				break
			}
			continue
		}

		// Success
		c.log.WithField("endpoint", endpoint).WithField("attempt", attempt).Debug("API call successful")
		return nil
	}

	c.log.WithError(lastErr).WithField("endpoint", endpoint).WithField("max_retries", c.maxRetries).Error("API call failed after all retries")
	return fmt.Errorf("API call failed after %d retries: %w", c.maxRetries, lastErr)
}

// processResponse handles the HTTP response and unmarshals it into the result
func (c *Client) processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}