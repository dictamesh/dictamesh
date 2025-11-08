// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package chatwoot provides a comprehensive adapter for integrating Chatwoot
// customer engagement platform with the DictaMesh framework.
//
// This adapter supports all three Chatwoot API types:
//   - Platform API: Multi-tenant account management
//   - Application API: Account-specific operations
//   - Public API: Client-side integrations
//
// For detailed documentation on usage patterns and best practices,
// see the README.md file in this package.
package chatwoot

import "time"

// Account represents a Chatwoot account
type Account struct {
	ID               int64                  `json:"id"`
	Name             string                 `json:"name"`
	Locale           string                 `json:"locale,omitempty"`
	Domain           string                 `json:"domain,omitempty"`
	SupportEmail     string                 `json:"support_email,omitempty"`
	Status           string                 `json:"status,omitempty"`
	AutoResolve      bool                   `json:"auto_resolve,omitempty"`
	CustomAttributes map[string]interface{} `json:"custom_attributes,omitempty"`
	CreatedAt        time.Time              `json:"created_at,omitempty"`
	UpdatedAt        time.Time              `json:"updated_at,omitempty"`
}

// User represents a Chatwoot user
type User struct {
	ID                 int64     `json:"id"`
	ProviderID         string    `json:"provider_id,omitempty"`
	UID                string    `json:"uid,omitempty"`
	Name               string    `json:"name"`
	DisplayName        string    `json:"display_name,omitempty"`
	Email              string    `json:"email"`
	AccountID          int64     `json:"account_id,omitempty"`
	Role               string    `json:"role,omitempty"`
	Confirmed          bool      `json:"confirmed,omitempty"`
	CustomAttributes   map[string]interface{} `json:"custom_attributes,omitempty"`
	AvailabilityStatus string    `json:"availability_status,omitempty"`
	AvatarURL          string    `json:"avatar_url,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// Agent represents a Chatwoot agent (user with agent role)
type Agent struct {
	ID                 int64                  `json:"id"`
	Name               string                 `json:"name"`
	Email              string                 `json:"email"`
	Role               string                 `json:"role"`
	AvailabilityStatus string                 `json:"availability_status"`
	AvatarURL          string                 `json:"avatar_url,omitempty"`
	Confirmed          bool                   `json:"confirmed"`
	CustomAttributes   map[string]interface{} `json:"custom_attributes,omitempty"`
}

// AgentBot represents a Chatwoot agent bot
type AgentBot struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	OutgoingURL      string    `json:"outgoing_url,omitempty"`
	BotType          string    `json:"bot_type,omitempty"`
	BotConfig        map[string]interface{} `json:"bot_config,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

// Contact represents a Chatwoot contact
type Contact struct {
	ID                  int64                  `json:"id"`
	Email               string                 `json:"email,omitempty"`
	Name                string                 `json:"name,omitempty"`
	PhoneNumber         string                 `json:"phone_number,omitempty"`
	Identifier          string                 `json:"identifier,omitempty"`
	Thumbnail           string                 `json:"thumbnail,omitempty"`
	AdditionalAttributes map[string]interface{} `json:"additional_attributes,omitempty"`
	CustomAttributes    map[string]interface{} `json:"custom_attributes,omitempty"`
	ContactInboxes      []ContactInbox         `json:"contact_inboxes,omitempty"`
	LastActivityAt      *time.Time             `json:"last_activity_at,omitempty"`
	CreatedAt           time.Time              `json:"created_at,omitempty"`
	UpdatedAt           time.Time              `json:"updated_at,omitempty"`
}

// ContactInbox represents the association between a contact and an inbox
type ContactInbox struct {
	ID          int64     `json:"id"`
	ContactID   int64     `json:"contact_id"`
	InboxID     int64     `json:"inbox_id"`
	SourceID    string    `json:"source_id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// Conversation represents a Chatwoot conversation
type Conversation struct {
	ID                 int64                  `json:"id"`
	AccountID          int64                  `json:"account_id"`
	InboxID            int64                  `json:"inbox_id"`
	ContactID          int64                  `json:"contact_id,omitempty"`
	AssigneeID         *int64                 `json:"assignee_id,omitempty"`
	TeamID             *int64                 `json:"team_id,omitempty"`
	Status             string                 `json:"status"`
	Priority           *string                `json:"priority,omitempty"`
	Channel            string                 `json:"channel,omitempty"`
	UUID               string                 `json:"uuid,omitempty"`
	Identifier         string                 `json:"identifier,omitempty"`
	AgentLastSeenAt    *time.Time             `json:"agent_last_seen_at,omitempty"`
	ContactLastSeenAt  *time.Time             `json:"contact_last_seen_at,omitempty"`
	Timestamp          int64                  `json:"timestamp,omitempty"`
	UnreadCount        int                    `json:"unread_count,omitempty"`
	AdditionalAttributes map[string]interface{} `json:"additional_attributes,omitempty"`
	CustomAttributes   map[string]interface{} `json:"custom_attributes,omitempty"`
	Labels             []string               `json:"labels,omitempty"`
	Messages           []Message              `json:"messages,omitempty"`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
	UpdatedAt          time.Time              `json:"updated_at,omitempty"`
}

// Message represents a Chatwoot message
type Message struct {
	ID             int64                  `json:"id"`
	Content        string                 `json:"content"`
	AccountID      int64                  `json:"account_id"`
	InboxID        int64                  `json:"inbox_id"`
	ConversationID int64                  `json:"conversation_id"`
	MessageType    int                    `json:"message_type"`
	ContentType    string                 `json:"content_type,omitempty"`
	ContentAttributes map[string]interface{} `json:"content_attributes,omitempty"`
	Status         string                 `json:"status,omitempty"`
	Private        bool                   `json:"private,omitempty"`
	SourceID       string                 `json:"source_id,omitempty"`
	SenderID       *int64                 `json:"sender_id,omitempty"`
	SenderType     string                 `json:"sender_type,omitempty"`
	Sender         interface{}            `json:"sender,omitempty"`
	Attachments    []Attachment           `json:"attachments,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty"`
	UpdatedAt      time.Time              `json:"updated_at,omitempty"`
}

// Attachment represents a file attachment
type Attachment struct {
	ID           int64     `json:"id"`
	MessageID    int64     `json:"message_id,omitempty"`
	FileType     string    `json:"file_type"`
	AccountID    int64     `json:"account_id,omitempty"`
	Extension    string    `json:"extension,omitempty"`
	DataURL      string    `json:"data_url"`
	ThumbURL     string    `json:"thumb_url,omitempty"`
	FileSize     int64     `json:"file_size,omitempty"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// Inbox represents a Chatwoot inbox
type Inbox struct {
	ID                   int64                  `json:"id"`
	Name                 string                 `json:"name"`
	ChannelID            int64                  `json:"channel_id"`
	ChannelType          string                 `json:"channel_type"`
	AvatarURL            string                 `json:"avatar_url,omitempty"`
	GreetingEnabled      bool                   `json:"greeting_enabled,omitempty"`
	GreetingMessage      string                 `json:"greeting_message,omitempty"`
	WorkingHoursEnabled  bool                   `json:"working_hours_enabled,omitempty"`
	EnableAutoAssignment bool                   `json:"enable_auto_assignment,omitempty"`
	OutOfOfficeMessage   string                 `json:"out_of_office_message,omitempty"`
	Timezone             string                 `json:"timezone,omitempty"`
	AllowMessagesAfterResolved bool             `json:"allow_messages_after_resolved,omitempty"`
	WebWidget            map[string]interface{} `json:"web_widget,omitempty"`
	CreatedAt            time.Time              `json:"created_at,omitempty"`
	UpdatedAt            time.Time              `json:"updated_at,omitempty"`
}

// Team represents a Chatwoot team
type Team struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description,omitempty"`
	AllowAutoAssign     bool   `json:"allow_auto_assign,omitempty"`
	AccountID           int64  `json:"account_id"`
	IsPrivate           bool   `json:"is_private,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}

// Label represents a Chatwoot label
type Label struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color"`
	ShowOnSidebar bool    `json:"show_on_sidebar,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// CannedResponse represents a pre-defined response template
type CannedResponse struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	ShortCode string    `json:"short_code"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CustomAttributeDefinition represents a custom attribute schema
type CustomAttributeDefinition struct {
	ID                 int64                  `json:"id"`
	AttributeKey       string                 `json:"attribute_key"`
	AttributeDisplayName string               `json:"attribute_display_name"`
	AttributeDisplayType string               `json:"attribute_display_type"`
	AttributeDescription string               `json:"attribute_description,omitempty"`
	AttributeModel     string                 `json:"attribute_model"`
	AttributeValues    []string               `json:"attribute_values,omitempty"`
	DefaultValue       interface{}            `json:"default_value,omitempty"`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
	UpdatedAt          time.Time              `json:"updated_at,omitempty"`
}

// AutomationRule represents an automation rule
type AutomationRule struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	EventName   string                   `json:"event_name"`
	AccountID   int64                    `json:"account_id"`
	Conditions  []AutomationCondition    `json:"conditions"`
	Actions     []AutomationAction       `json:"actions"`
	Active      bool                     `json:"active,omitempty"`
	CreatedAt   time.Time                `json:"created_at,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at,omitempty"`
}

// AutomationCondition represents a condition in an automation rule
type AutomationCondition struct {
	AttributeKey string                 `json:"attribute_key"`
	FilterOperator string               `json:"filter_operator"`
	Values       []interface{}          `json:"values"`
	QueryOperator string                `json:"query_operator,omitempty"`
	CustomAttributeType string          `json:"custom_attribute_type,omitempty"`
}

// AutomationAction represents an action in an automation rule
type AutomationAction struct {
	ActionName   string        `json:"action_name"`
	ActionParams []interface{} `json:"action_params"`
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID           int64     `json:"id"`
	AccountID    int64     `json:"account_id"`
	InboxID      *int64    `json:"inbox_id,omitempty"`
	URL          string    `json:"url"`
	WebhookType  string    `json:"webhook_type,omitempty"`
	Subscriptions []string `json:"subscriptions,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// IntegrationApp represents an available integration
type IntegrationApp struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Logo        string                 `json:"logo,omitempty"`
	Enabled     bool                   `json:"enabled,omitempty"`
	Hooks       []IntegrationHook      `json:"hooks,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// IntegrationHook represents a hook for an integration
type IntegrationHook struct {
	ID        int64                  `json:"id"`
	AppID     string                 `json:"app_id"`
	InboxID   *int64                 `json:"inbox_id,omitempty"`
	AccountID int64                  `json:"account_id"`
	Status    string                 `json:"status,omitempty"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitempty"`
	UpdatedAt time.Time              `json:"updated_at,omitempty"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID            int64                  `json:"id"`
	AccountID     int64                  `json:"account_id"`
	UserID        *int64                 `json:"user_id,omitempty"`
	Activity      string                 `json:"activity"`
	AuditableType string                 `json:"auditable_type,omitempty"`
	AuditableID   *int64                 `json:"auditable_id,omitempty"`
	AssociatedType string                `json:"associated_type,omitempty"`
	AssociatedID  *int64                 `json:"associated_id,omitempty"`
	RemoteAddress string                 `json:"remote_address,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at,omitempty"`
}

// Report represents analytics report data
type Report struct {
	Metric    string                   `json:"metric"`
	Value     interface{}              `json:"value"`
	Timestamp time.Time                `json:"timestamp,omitempty"`
	Breakdown map[string]interface{}   `json:"breakdown,omitempty"`
}

// ConversationMetrics represents conversation analytics
type ConversationMetrics struct {
	Open          int     `json:"open"`
	Unattended    int     `json:"unattended"`
	Resolved      int     `json:"resolved"`
	TotalCount    int     `json:"total_count"`
	AverageResolutionTime float64 `json:"avg_resolution_time,omitempty"`
	AverageFirstResponseTime float64 `json:"avg_first_response_time,omitempty"`
}

// PaginationMeta represents pagination metadata in API responses
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	Count       int   `json:"count"`
	TotalCount  int64 `json:"total_count,omitempty"`
	TotalPages  int   `json:"total_pages,omitempty"`
}

// ListResponse represents a generic paginated list response
type ListResponse struct {
	Payload  interface{}     `json:"payload"`
	Meta     *PaginationMeta `json:"meta,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

// WebhookPayload represents the payload received via webhook
type WebhookPayload struct {
	Event         string                 `json:"event"`
	ID            interface{}            `json:"id,omitempty"`
	MessageType   string                 `json:"message_type,omitempty"`
	Conversation  *Conversation          `json:"conversation,omitempty"`
	Message       *Message               `json:"message,omitempty"`
	Account       *Account               `json:"account,omitempty"`
	Contact       *Contact               `json:"contact,omitempty"`
	Inbox         *Inbox                 `json:"inbox,omitempty"`
	AdditionalAttributes map[string]interface{} `json:"additional_attributes,omitempty"`
	CreatedAt     time.Time              `json:"created_at,omitempty"`
}
