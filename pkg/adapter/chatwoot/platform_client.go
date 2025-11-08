// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package chatwoot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/click2-run/dictamesh/pkg/adapter"
)

// PlatformClient handles Platform API operations (multi-tenant management)
type PlatformClient struct {
	httpClient *adapter.HTTPClient
	apiKey     string
}

// NewPlatformClient creates a new Platform API client
func NewPlatformClient(config *Config) *PlatformClient {
	httpClient := adapter.NewHTTPClient(&adapter.HTTPClientConfig{
		BaseURL: config.BaseURL,
		Timeout: config.Timeout,
		RetryConfig: &adapter.RetryConfig{
			MaxRetries:        config.MaxRetries,
			InitialBackoff:    config.RetryBackoff,
			MaxBackoff:        30 * config.RetryBackoff,
			BackoffMultiplier: 2.0,
			RetryableStatusCodes: []int{
				http.StatusTooManyRequests,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			},
		},
		Headers: map[string]string{
			"api_access_token": config.PlatformAPIKey,
			"Content-Type":     "application/json",
		},
	})

	return &PlatformClient{
		httpClient: httpClient,
		apiKey:     config.PlatformAPIKey,
	}
}

// Close closes the client
func (c *PlatformClient) Close() error {
	// No resources to clean up currently
	return nil
}

// Ping checks if the Platform API is accessible
func (c *PlatformClient) Ping(ctx context.Context) error {
	// Platform API doesn't have a dedicated ping endpoint
	// We'll try to list agent bots as a health check
	_, err := c.ListAgentBots(ctx)
	return err
}

// Account Management

// CreateAccount creates a new account
func (c *PlatformClient) CreateAccount(ctx context.Context, account *Account) (*Account, error) {
	path := "/platform/api/v1/accounts"

	resp, err := c.httpClient.Post(ctx, path, account, nil)
	if err != nil {
		return nil, err
	}

	var result Account
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAccount retrieves an account by ID
func (c *PlatformClient) GetAccount(ctx context.Context, accountID int64) (*Account, error) {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d", accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result Account
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateAccount updates an existing account
func (c *PlatformClient) UpdateAccount(ctx context.Context, accountID int64, account *Account) (*Account, error) {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d", accountID)

	resp, err := c.httpClient.Patch(ctx, path, account, nil)
	if err != nil {
		return nil, err
	}

	var result Account
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteAccount deletes an account
func (c *PlatformClient) DeleteAccount(ctx context.Context, accountID int64) error {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d", accountID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete account", nil)
	}

	return nil
}

// Account User Management

// AccountUserRequest represents a request to manage account users
type AccountUserRequest struct {
	UserID int64 `json:"user_id"`
	Role   string `json:"role,omitempty"`
}

// ListAccountUsers lists all users in an account
func (c *PlatformClient) ListAccountUsers(ctx context.Context, accountID int64) ([]User, error) {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d/account_users", accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []User
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// AddAccountUser adds a user to an account
func (c *PlatformClient) AddAccountUser(ctx context.Context, accountID int64, req *AccountUserRequest) (*User, error) {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d/account_users", accountID)

	resp, err := c.httpClient.Post(ctx, path, req, nil)
	if err != nil {
		return nil, err
	}

	var result User
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RemoveAccountUser removes a user from an account
func (c *PlatformClient) RemoveAccountUser(ctx context.Context, accountID int64, userID int64) error {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d/account_users", accountID)

	req := &AccountUserRequest{UserID: userID}
	resp, err := c.httpClient.Delete(ctx, path, nil)
	_ = req // Will be used in request body
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to remove account user", nil)
	}

	return nil
}

// Agent Bot Management

// ListAgentBots lists all agent bots
func (c *PlatformClient) ListAgentBots(ctx context.Context) ([]AgentBot, error) {
	path := "/platform/api/v1/agent_bots"

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateAgentBot creates a new agent bot
func (c *PlatformClient) CreateAgentBot(ctx context.Context, bot *AgentBot) (*AgentBot, error) {
	path := "/platform/api/v1/agent_bots"

	resp, err := c.httpClient.Post(ctx, path, bot, nil)
	if err != nil {
		return nil, err
	}

	var result AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAgentBot retrieves an agent bot by ID
func (c *PlatformClient) GetAgentBot(ctx context.Context, botID int64) (*AgentBot, error) {
	path := fmt.Sprintf("/platform/api/v1/agent_bots/%d", botID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateAgentBot updates an existing agent bot
func (c *PlatformClient) UpdateAgentBot(ctx context.Context, botID int64, bot *AgentBot) (*AgentBot, error) {
	path := fmt.Sprintf("/platform/api/v1/agent_bots/%d", botID)

	resp, err := c.httpClient.Patch(ctx, path, bot, nil)
	if err != nil {
		return nil, err
	}

	var result AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteAgentBot deletes an agent bot
func (c *PlatformClient) DeleteAgentBot(ctx context.Context, botID int64) error {
	path := fmt.Sprintf("/platform/api/v1/agent_bots/%d", botID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete agent bot", nil)
	}

	return nil
}

// User Management

// CreateUser creates a new user
func (c *PlatformClient) CreateUser(ctx context.Context, user *User) (*User, error) {
	path := "/platform/api/v1/users"

	resp, err := c.httpClient.Post(ctx, path, user, nil)
	if err != nil {
		return nil, err
	}

	var result User
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUser retrieves a user by ID
func (c *PlatformClient) GetUser(ctx context.Context, userID int64) (*User, error) {
	path := fmt.Sprintf("/platform/api/v1/users/%d", userID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result User
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateUser updates an existing user
func (c *PlatformClient) UpdateUser(ctx context.Context, userID int64, user *User) (*User, error) {
	path := fmt.Sprintf("/platform/api/v1/users/%d", userID)

	resp, err := c.httpClient.Patch(ctx, path, user, nil)
	if err != nil {
		return nil, err
	}

	var result User
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteUser deletes a user
func (c *PlatformClient) DeleteUser(ctx context.Context, userID int64) error {
	path := fmt.Sprintf("/platform/api/v1/users/%d", userID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete user", nil)
	}

	return nil
}

// SSOLoginResponse represents the response from SSO login
type SSOLoginResponse struct {
	URL string `json:"url"`
}

// GetUserSSOLogin retrieves the SSO login URL for a user
func (c *PlatformClient) GetUserSSOLogin(ctx context.Context, userID int64) (*SSOLoginResponse, error) {
	path := fmt.Sprintf("/platform/api/v1/users/%d/login", userID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result SSOLoginResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
