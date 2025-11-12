// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package chatwoot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/click2-run/dictamesh/pkg/adapter"
)

// ApplicationClient handles Application/User API operations (account-specific)
type ApplicationClient struct {
	httpClient *adapter.HTTPClient
	accountID  int64
	apiKey     string
}

// NewApplicationClient creates a new Application API client
func NewApplicationClient(config *Config) *ApplicationClient {
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
			"api_access_token": config.UserAPIKey,
			"Content-Type":     "application/json",
		},
	})

	return &ApplicationClient{
		httpClient: httpClient,
		accountID:  config.AccountID,
		apiKey:     config.UserAPIKey,
	}
}

// Close closes the client
func (c *ApplicationClient) Close() error {
	return nil
}

// Ping checks if the Application API is accessible
func (c *ApplicationClient) Ping(ctx context.Context) error {
	_, err := c.GetAccount(ctx)
	return err
}

// Account Management

// GetAccount retrieves the current account
func (c *ApplicationClient) GetAccount(ctx context.Context) (*Account, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d", c.accountID)

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

// UpdateAccount updates the current account
func (c *ApplicationClient) UpdateAccount(ctx context.Context, account *Account) (*Account, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d", c.accountID)

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

// Agent Management

// ListAgents lists all agents in the account
func (c *ApplicationClient) ListAgents(ctx context.Context) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []Agent
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// AddAgent adds an agent to the account
func (c *ApplicationClient) AddAgent(ctx context.Context, agent *Agent) (*Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, agent, nil)
	if err != nil {
		return nil, err
	}

	var result Agent
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateAgent updates an agent
func (c *ApplicationClient) UpdateAgent(ctx context.Context, agentID int64, agent *Agent) (*Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents/%d", c.accountID, agentID)

	resp, err := c.httpClient.Patch(ctx, path, agent, nil)
	if err != nil {
		return nil, err
	}

	var result Agent
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RemoveAgent removes an agent from the account
func (c *ApplicationClient) RemoveAgent(ctx context.Context, agentID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents/%d", c.accountID, agentID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to remove agent", nil)
	}

	return nil
}

// Contact Management

// ContactListOptions represents options for listing contacts
type ContactListOptions struct {
	Page     int    `json:"page,omitempty"`
	Sort     string `json:"sort,omitempty"`
}

// ListContacts lists all contacts
func (c *ApplicationClient) ListContacts(ctx context.Context, opts *ContactListOptions) (*ListResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts", c.accountID)

	if opts != nil && opts.Page > 0 {
		params := url.Values{}
		params.Set("page", strconv.Itoa(opts.Page))
		if opts.Sort != "" {
			params.Set("sort", opts.Sort)
		}
		path = path + "?" + params.Encode()
	}

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetContact retrieves a contact by ID
func (c *ApplicationClient) GetContact(ctx context.Context, contactID int64) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d", c.accountID, contactID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload Contact `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// CreateContact creates a new contact
func (c *ApplicationClient) CreateContact(ctx context.Context, contact *Contact) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, contact, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload Contact `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// UpdateContact updates a contact
func (c *ApplicationClient) UpdateContact(ctx context.Context, contactID int64, contact *Contact) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d", c.accountID, contactID)

	resp, err := c.httpClient.Put(ctx, path, contact, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload Contact `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// DeleteContact deletes a contact
func (c *ApplicationClient) DeleteContact(ctx context.Context, contactID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d", c.accountID, contactID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete contact", nil)
	}

	return nil
}

// SearchContacts searches for contacts
func (c *ApplicationClient) SearchContacts(ctx context.Context, query string, page int) (*ListResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/search", c.accountID)

	params := url.Values{}
	params.Set("q", query)
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	path = path + "?" + params.Encode()

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// FilterContacts filters contacts based on criteria
func (c *ApplicationClient) FilterContacts(ctx context.Context, filter map[string]interface{}, page int) (*ListResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/filter", c.accountID)

	if page > 0 {
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		path = path + "?" + params.Encode()
	}

	resp, err := c.httpClient.Post(ctx, path, filter, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListContactConversations lists all conversations for a contact
func (c *ApplicationClient) ListContactConversations(ctx context.Context, contactID int64) ([]Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d/conversations", c.accountID, contactID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []Conversation `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// Conversation Management

// ConversationListOptions represents options for listing conversations
type ConversationListOptions struct {
	AssigneeType string `json:"assignee_type,omitempty"`
	Status       string `json:"status,omitempty"`
	Page         int    `json:"page,omitempty"`
	InboxID      int64  `json:"inbox_id,omitempty"`
	Labels       []string `json:"labels,omitempty"`
}

// ListConversations lists all conversations
func (c *ApplicationClient) ListConversations(ctx context.Context, opts *ConversationListOptions) (*ListResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations", c.accountID)

	if opts != nil {
		params := url.Values{}
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if opts.AssigneeType != "" {
			params.Set("assignee_type", opts.AssigneeType)
		}
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.InboxID > 0 {
			params.Set("inbox_id", strconv.FormatInt(opts.InboxID, 10))
		}
		for _, label := range opts.Labels {
			params.Add("labels[]", label)
		}
		if len(params) > 0 {
			path = path + "?" + params.Encode()
		}
	}

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetConversation retrieves a conversation by ID
func (c *ApplicationClient) GetConversation(ctx context.Context, conversationID int64) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d", c.accountID, conversationID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result Conversation
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateConversation updates a conversation
func (c *ApplicationClient) UpdateConversation(ctx context.Context, conversationID int64, conversation *Conversation) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d", c.accountID, conversationID)

	resp, err := c.httpClient.Patch(ctx, path, conversation, nil)
	if err != nil {
		return nil, err
	}

	var result Conversation
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AssignConversation assigns a conversation to an agent
func (c *ApplicationClient) AssignConversation(ctx context.Context, conversationID int64, assigneeID int64, teamID *int64) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/assignments", c.accountID, conversationID)

	payload := map[string]interface{}{
		"assignee_id": assigneeID,
	}
	if teamID != nil {
		payload["team_id"] = *teamID
	}

	resp, err := c.httpClient.Post(ctx, path, payload, nil)
	if err != nil {
		return nil, err
	}

	var result Conversation
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ToggleConversationStatus toggles the status of a conversation (open/resolved)
func (c *ApplicationClient) ToggleConversationStatus(ctx context.Context, conversationID int64, status string) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/toggle_status", c.accountID, conversationID)

	payload := map[string]interface{}{
		"status": status,
	}

	resp, err := c.httpClient.Post(ctx, path, payload, nil)
	if err != nil {
		return nil, err
	}

	var result Conversation
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Message Management

// ListMessages lists all messages in a conversation
func (c *ApplicationClient) ListMessages(ctx context.Context, conversationID int64) ([]Message, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages", c.accountID, conversationID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []Message `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// CreateMessage creates a new message in a conversation
func (c *ApplicationClient) CreateMessage(ctx context.Context, conversationID int64, message *Message) (*Message, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages", c.accountID, conversationID)

	resp, err := c.httpClient.Post(ctx, path, message, nil)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateMessage updates a message
func (c *ApplicationClient) UpdateMessage(ctx context.Context, conversationID int64, messageID int64, message *Message) (*Message, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages/%d", c.accountID, conversationID, messageID)

	resp, err := c.httpClient.Patch(ctx, path, message, nil)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteMessage deletes a message
func (c *ApplicationClient) DeleteMessage(ctx context.Context, conversationID int64, messageID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages/%d", c.accountID, conversationID, messageID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete message", nil)
	}

	return nil
}

// Continued in next message due to length...
