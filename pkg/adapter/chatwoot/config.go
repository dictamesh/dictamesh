// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package chatwoot

import (
	"errors"
	"fmt"
	"time"

	"github.com/click2-run/dictamesh/pkg/adapter"
)

// Config represents the configuration for the Chatwoot adapter
type Config struct {
	// BaseURL is the base URL of the Chatwoot instance (e.g., "https://app.chatwoot.com")
	BaseURL string

	// PlatformAPIKey is the API key for Platform API (multi-tenant operations)
	PlatformAPIKey string

	// UserAPIKey is the API key for Application/User API (account-level operations)
	UserAPIKey string

	// AccountID is the Chatwoot account ID for Application API
	AccountID int64

	// InboxIdentifier is the inbox identifier for Public API
	InboxIdentifier string

	// Timeout is the request timeout duration
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// RetryBackoff is the initial backoff duration for retries
	RetryBackoff time.Duration

	// RateLimitPerSecond is the maximum number of requests per second
	RateLimitPerSecond int

	// EnableRequestLogging enables HTTP request/response logging
	EnableRequestLogging bool

	// WebhookSecret is the secret for webhook signature verification
	WebhookSecret string

	// EnablePlatformAPI enables Platform API features
	EnablePlatformAPI bool

	// EnableApplicationAPI enables Application API features
	EnableApplicationAPI bool

	// EnablePublicAPI enables Public API features
	EnablePublicAPI bool
}

// NewConfig creates a new Chatwoot configuration from a map
func NewConfig(data map[string]interface{}) (*Config, error) {
	cfg := &Config{
		Timeout:              30 * time.Second,
		MaxRetries:           3,
		RetryBackoff:         1 * time.Second,
		RateLimitPerSecond:   10,
		EnableRequestLogging: false,
	}

	// Parse configuration from map
	if baseURL, ok := data["base_url"].(string); ok {
		cfg.BaseURL = baseURL
	}

	if platformAPIKey, ok := data["platform_api_key"].(string); ok {
		cfg.PlatformAPIKey = platformAPIKey
		cfg.EnablePlatformAPI = true
	}

	if userAPIKey, ok := data["user_api_key"].(string); ok {
		cfg.UserAPIKey = userAPIKey
		cfg.EnableApplicationAPI = true
	}

	if accountID, ok := data["account_id"].(int64); ok {
		cfg.AccountID = accountID
	} else if accountID, ok := data["account_id"].(float64); ok {
		cfg.AccountID = int64(accountID)
	}

	if inboxIdentifier, ok := data["inbox_identifier"].(string); ok {
		cfg.InboxIdentifier = inboxIdentifier
		cfg.EnablePublicAPI = true
	}

	if webhookSecret, ok := data["webhook_secret"].(string); ok {
		cfg.WebhookSecret = webhookSecret
	}

	if timeout, ok := data["timeout"].(string); ok {
		if d, err := time.ParseDuration(timeout); err == nil {
			cfg.Timeout = d
		}
	}

	if maxRetries, ok := data["max_retries"].(int); ok {
		cfg.MaxRetries = maxRetries
	} else if maxRetries, ok := data["max_retries"].(float64); ok {
		cfg.MaxRetries = int(maxRetries)
	}

	if retryBackoff, ok := data["retry_backoff"].(string); ok {
		if d, err := time.ParseDuration(retryBackoff); err == nil {
			cfg.RetryBackoff = d
		}
	}

	if rateLimit, ok := data["rate_limit_per_second"].(int); ok {
		cfg.RateLimitPerSecond = rateLimit
	} else if rateLimit, ok := data["rate_limit_per_second"].(float64); ok {
		cfg.RateLimitPerSecond = int(rateLimit)
	}

	if enableLogging, ok := data["enable_request_logging"].(bool); ok {
		cfg.EnableRequestLogging = enableLogging
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return errors.New("base_url is required")
	}

	// At least one API must be enabled
	if !c.EnablePlatformAPI && !c.EnableApplicationAPI && !c.EnablePublicAPI {
		return errors.New("at least one API must be enabled (platform, application, or public)")
	}

	// Validate Platform API configuration
	if c.EnablePlatformAPI && c.PlatformAPIKey == "" {
		return errors.New("platform_api_key is required when platform API is enabled")
	}

	// Validate Application API configuration
	if c.EnableApplicationAPI {
		if c.UserAPIKey == "" {
			return errors.New("user_api_key is required when application API is enabled")
		}
		if c.AccountID == 0 {
			return errors.New("account_id is required when application API is enabled")
		}
	}

	// Validate Public API configuration
	if c.EnablePublicAPI && c.InboxIdentifier == "" {
		return errors.New("inbox_identifier is required when public API is enabled")
	}

	// Validate timeouts and limits
	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	if c.MaxRetries < 0 {
		return errors.New("max_retries cannot be negative")
	}

	if c.RateLimitPerSecond < 0 {
		return errors.New("rate_limit_per_second cannot be negative")
	}

	return nil
}

// GetString retrieves a string configuration value
func (c *Config) GetString(key string) (string, error) {
	switch key {
	case "base_url":
		return c.BaseURL, nil
	case "platform_api_key":
		return c.PlatformAPIKey, nil
	case "user_api_key":
		return c.UserAPIKey, nil
	case "inbox_identifier":
		return c.InboxIdentifier, nil
	case "webhook_secret":
		return c.WebhookSecret, nil
	default:
		return "", fmt.Errorf("unknown configuration key: %s", key)
	}
}

// GetInt retrieves an integer configuration value
func (c *Config) GetInt(key string) (int, error) {
	switch key {
	case "max_retries":
		return c.MaxRetries, nil
	case "rate_limit_per_second":
		return c.RateLimitPerSecond, nil
	case "account_id":
		return int(c.AccountID), nil
	default:
		return 0, fmt.Errorf("unknown configuration key: %s", key)
	}
}

// GetBool retrieves a boolean configuration value
func (c *Config) GetBool(key string) (bool, error) {
	switch key {
	case "enable_request_logging":
		return c.EnableRequestLogging, nil
	case "enable_platform_api":
		return c.EnablePlatformAPI, nil
	case "enable_application_api":
		return c.EnableApplicationAPI, nil
	case "enable_public_api":
		return c.EnablePublicAPI, nil
	default:
		return false, fmt.Errorf("unknown configuration key: %s", key)
	}
}

// GetDuration retrieves a duration configuration value
func (c *Config) GetDuration(key string) (time.Duration, error) {
	switch key {
	case "timeout":
		return c.Timeout, nil
	case "retry_backoff":
		return c.RetryBackoff, nil
	default:
		return 0, fmt.Errorf("unknown configuration key: %s", key)
	}
}

// ToMap converts the configuration to a map
func (c *Config) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"base_url":               c.BaseURL,
		"platform_api_key":       c.PlatformAPIKey,
		"user_api_key":           c.UserAPIKey,
		"account_id":             c.AccountID,
		"inbox_identifier":       c.InboxIdentifier,
		"timeout":                c.Timeout.String(),
		"max_retries":            c.MaxRetries,
		"retry_backoff":          c.RetryBackoff.String(),
		"rate_limit_per_second":  c.RateLimitPerSecond,
		"enable_request_logging": c.EnableRequestLogging,
		"webhook_secret":         c.WebhookSecret,
		"enable_platform_api":    c.EnablePlatformAPI,
		"enable_application_api": c.EnableApplicationAPI,
		"enable_public_api":      c.EnablePublicAPI,
	}
}

// Ensure Config implements adapter.Config
var _ adapter.Config = (*Config)(nil)
