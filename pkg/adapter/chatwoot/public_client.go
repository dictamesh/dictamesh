// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package chatwoot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/click2-run/dictamesh/pkg/adapter"
)

// PublicClient handles Public API operations (client-side integrations)
type PublicClient struct {
	httpClient      *adapter.HTTPClient
	inboxIdentifier string
}

// NewPublicClient creates a new Public API client
func NewPublicClient(config *Config) *PublicClient {
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
			"Content-Type": "application/json",
		},
	})

	return &PublicClient{
		httpClient:      httpClient,
		inboxIdentifier: config.InboxIdentifier,
	}
}

// Close closes the client
func (c *PublicClient) Close() error {
	return nil
}

// Inbox Operations

// GetInbox retrieves inbox information
func (c *PublicClient) GetInbox(ctx context.Context) (*Inbox, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s", c.inboxIdentifier)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result Inbox
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Contact Operations

// CreateContact creates a new contact in the inbox
func (c *PublicClient) CreateContact(ctx context.Context, contact *Contact) (*Contact, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts", c.inboxIdentifier)

	resp, err := c.httpClient.Post(ctx, path, contact, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Pubsub_token string  `json:"pubsub_token"`
		Contact      Contact `json:"contact"`
		SourceID     string  `json:"source_id"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Contact, nil
}

// GetContact retrieves a contact by identifier
func (c *PublicClient) GetContact(ctx context.Context, contactIdentifier string) (*Contact, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s", c.inboxIdentifier, contactIdentifier)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result Contact
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateContact updates a contact
func (c *PublicClient) UpdateContact(ctx context.Context, contactIdentifier string, contact *Contact) (*Contact, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s", c.inboxIdentifier, contactIdentifier)

	resp, err := c.httpClient.Patch(ctx, path, contact, nil)
	if err != nil {
		return nil, err
	}

	var result Contact
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Conversation Operations

// ConversationRequest represents a request to create a conversation
type ConversationRequest struct {
	SourceID            string                 `json:"source_id,omitempty"`
	ContactID           string                 `json:"contact_id,omitempty"`
	AdditionalAttributes map[string]interface{} `json:"additional_attributes,omitempty"`
	CustomAttributes    map[string]interface{} `json:"custom_attributes,omitempty"`
	Message             string                 `json:"message,omitempty"`
}

// CreateConversation creates a new conversation
func (c *PublicClient) CreateConversation(ctx context.Context, contactIdentifier string, req *ConversationRequest) (*Conversation, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations", c.inboxIdentifier, contactIdentifier)

	resp, err := c.httpClient.Post(ctx, path, req, nil)
	if err != nil {
		return nil, err
	}

	var result Conversation
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListConversations lists all conversations for a contact
func (c *PublicClient) ListConversations(ctx context.Context, contactIdentifier string) ([]Conversation, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations", c.inboxIdentifier, contactIdentifier)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []Conversation
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetConversation retrieves a conversation by ID
func (c *PublicClient) GetConversation(ctx context.Context, contactIdentifier string, conversationID int64) (*Conversation, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d", c.inboxIdentifier, contactIdentifier, conversationID)

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

// ResolveConversation marks a conversation as resolved
func (c *PublicClient) ResolveConversation(ctx context.Context, contactIdentifier string, conversationID int64) (*Conversation, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/toggle_status", c.inboxIdentifier, contactIdentifier, conversationID)

	resp, err := c.httpClient.Post(ctx, path, map[string]interface{}{"status": "resolved"}, nil)
	if err != nil {
		return nil, err
	}

	var result Conversation
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ToggleTypingStatus toggles the typing status in a conversation
func (c *PublicClient) ToggleTypingStatus(ctx context.Context, contactIdentifier string, conversationID int64, typing bool) error {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/toggle_typing", c.inboxIdentifier, contactIdentifier, conversationID)

	payload := map[string]interface{}{
		"typing_status": "on",
	}
	if !typing {
		payload["typing_status"] = "off"
	}

	resp, err := c.httpClient.Post(ctx, path, payload, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to toggle typing status", nil)
	}

	return nil
}

// UpdateLastSeen updates the last seen timestamp for a conversation
func (c *PublicClient) UpdateLastSeen(ctx context.Context, contactIdentifier string, conversationID int64) error {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/update_last_seen", c.inboxIdentifier, contactIdentifier, conversationID)

	resp, err := c.httpClient.Post(ctx, path, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to update last seen", nil)
	}

	return nil
}

// Message Operations

// MessageRequest represents a request to create a message
type MessageRequest struct {
	Content           string                 `json:"content"`
	MessageType       string                 `json:"message_type,omitempty"`
	Private           bool                   `json:"private,omitempty"`
	ContentAttributes map[string]interface{} `json:"content_attributes,omitempty"`
	Attachments       []Attachment           `json:"attachments,omitempty"`
}

// CreateMessage creates a new message in a conversation
func (c *PublicClient) CreateMessage(ctx context.Context, contactIdentifier string, conversationID int64, message *MessageRequest) (*Message, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/messages", c.inboxIdentifier, contactIdentifier, conversationID)

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

// ListMessages lists all messages in a conversation
func (c *PublicClient) ListMessages(ctx context.Context, contactIdentifier string, conversationID int64) ([]Message, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/messages", c.inboxIdentifier, contactIdentifier, conversationID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []Message
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateMessage updates a message
func (c *PublicClient) UpdateMessage(ctx context.Context, contactIdentifier string, conversationID int64, messageID int64, message *MessageRequest) (*Message, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/messages/%d", c.inboxIdentifier, contactIdentifier, conversationID, messageID)

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

// SetInboxIdentifier sets the inbox identifier for the client
// This allows switching between different inboxes dynamically
func (c *PublicClient) SetInboxIdentifier(inboxIdentifier string) {
	c.inboxIdentifier = inboxIdentifier
}

// GetInboxIdentifier returns the current inbox identifier
func (c *PublicClient) GetInboxIdentifier() string {
	return c.inboxIdentifier
}
