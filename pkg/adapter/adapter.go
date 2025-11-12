// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package adapter provides the base interfaces and types for DictaMesh adapters.
// This package defines the contracts that all third-party adapters must implement
// to integrate with the DictaMesh framework.
package adapter

import (
	"context"
	"time"
)

// Adapter is the core interface that all third-party adapters must implement.
// It provides the foundation for integrating external systems into the DictaMesh framework.
type Adapter interface {
	// Name returns the unique identifier for this adapter (e.g., "chatwoot", "salesforce")
	Name() string

	// Version returns the adapter version following semantic versioning
	Version() string

	// Initialize sets up the adapter with the provided configuration
	Initialize(ctx context.Context, config Config) error

	// Health checks if the adapter and its connection to the external system are healthy
	Health(ctx context.Context) (*HealthStatus, error)

	// Shutdown gracefully closes all connections and cleans up resources
	Shutdown(ctx context.Context) error

	// GetCapabilities returns the list of capabilities supported by this adapter
	GetCapabilities() []Capability
}

// Config represents the configuration for an adapter
type Config interface {
	// GetString retrieves a string configuration value
	GetString(key string) (string, error)

	// GetInt retrieves an integer configuration value
	GetInt(key string) (int, error)

	// GetBool retrieves a boolean configuration value
	GetBool(key string) (bool, error)

	// GetDuration retrieves a duration configuration value
	GetDuration(key string) (time.Duration, error)

	// Validate checks if the configuration is valid
	Validate() error
}

// HealthStatus represents the health status of an adapter
type HealthStatus struct {
	// Status indicates if the adapter is healthy
	Status HealthStatusType `json:"status"`

	// Message provides additional context about the health status
	Message string `json:"message,omitempty"`

	// Timestamp when the health check was performed
	Timestamp time.Time `json:"timestamp"`

	// Details contains adapter-specific health information
	Details map[string]interface{} `json:"details,omitempty"`

	// Latency of the health check operation
	Latency time.Duration `json:"latency"`
}

// HealthStatusType represents the health status type
type HealthStatusType string

const (
	// HealthStatusHealthy indicates the adapter is fully operational
	HealthStatusHealthy HealthStatusType = "healthy"

	// HealthStatusDegraded indicates the adapter is operational but with reduced functionality
	HealthStatusDegraded HealthStatusType = "degraded"

	// HealthStatusUnhealthy indicates the adapter is not operational
	HealthStatusUnhealthy HealthStatusType = "unhealthy"
)

// Capability represents a capability that an adapter supports
type Capability string

const (
	// CapabilityRead indicates the adapter can read data from the external system
	CapabilityRead Capability = "read"

	// CapabilityWrite indicates the adapter can write data to the external system
	CapabilityWrite Capability = "write"

	// CapabilityStream indicates the adapter can stream data changes in real-time
	CapabilityStream Capability = "stream"

	// CapabilityWebhooks indicates the adapter can receive webhooks from the external system
	CapabilityWebhooks Capability = "webhooks"

	// CapabilityBatch indicates the adapter supports batch operations
	CapabilityBatch Capability = "batch"

	// CapabilitySearch indicates the adapter supports search operations
	CapabilitySearch Capability = "search"

	// CapabilityPagination indicates the adapter supports paginated queries
	CapabilityPagination Capability = "pagination"
)

// ResourceAdapter is an interface for adapters that manage external resources
type ResourceAdapter interface {
	Adapter

	// ListResources lists all available resources from the external system
	ListResources(ctx context.Context, opts *ListOptions) (*ResourceList, error)

	// GetResource retrieves a specific resource by ID
	GetResource(ctx context.Context, resourceType, resourceID string) (*Resource, error)

	// CreateResource creates a new resource in the external system
	CreateResource(ctx context.Context, resource *Resource) (*Resource, error)

	// UpdateResource updates an existing resource
	UpdateResource(ctx context.Context, resource *Resource) (*Resource, error)

	// DeleteResource deletes a resource from the external system
	DeleteResource(ctx context.Context, resourceType, resourceID string) error
}

// ListOptions represents options for listing resources
type ListOptions struct {
	// Page is the page number for pagination (1-indexed)
	Page int `json:"page"`

	// PageSize is the number of items per page
	PageSize int `json:"page_size"`

	// Filter contains filter criteria
	Filter map[string]interface{} `json:"filter,omitempty"`

	// Sort specifies the sort field and direction
	Sort string `json:"sort,omitempty"`

	// Fields specifies which fields to include in the response
	Fields []string `json:"fields,omitempty"`
}

// ResourceList represents a paginated list of resources
type ResourceList struct {
	// Items contains the resources
	Items []*Resource `json:"items"`

	// Total is the total number of items available
	Total int `json:"total"`

	// Page is the current page number
	Page int `json:"page"`

	// PageSize is the number of items per page
	PageSize int `json:"page_size"`

	// HasMore indicates if there are more pages available
	HasMore bool `json:"has_more"`
}

// Resource represents a generic resource from an external system
type Resource struct {
	// ID is the unique identifier of the resource in the external system
	ID string `json:"id"`

	// Type is the resource type (e.g., "contact", "conversation", "inbox")
	Type string `json:"type"`

	// Attributes contains the resource attributes
	Attributes map[string]interface{} `json:"attributes"`

	// Relationships contains related resources
	Relationships map[string]interface{} `json:"relationships,omitempty"`

	// Metadata contains additional metadata about the resource
	Metadata *ResourceMetadata `json:"metadata"`

	// Raw contains the original raw data from the external system
	Raw interface{} `json:"raw,omitempty"`
}

// ResourceMetadata contains metadata about a resource
type ResourceMetadata struct {
	// CreatedAt is when the resource was created
	CreatedAt time.Time `json:"created_at,omitempty"`

	// UpdatedAt is when the resource was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// Source identifies the external system
	Source string `json:"source"`

	// SourceURL is the URL to access the resource in the external system
	SourceURL string `json:"source_url,omitempty"`

	// Version is the resource version for optimistic locking
	Version string `json:"version,omitempty"`

	// Etag is the entity tag for caching
	Etag string `json:"etag,omitempty"`

	// ExtraFields for adapter-specific metadata
	ExtraFields map[string]interface{} `json:"extra_fields,omitempty"`
}

// StreamingAdapter is an interface for adapters that support real-time data streaming
type StreamingAdapter interface {
	Adapter

	// Subscribe creates a subscription to stream events from the external system
	Subscribe(ctx context.Context, opts *SubscriptionOptions) (<-chan *Event, error)

	// Unsubscribe cancels an active subscription
	Unsubscribe(ctx context.Context, subscriptionID string) error
}

// SubscriptionOptions represents options for creating a subscription
type SubscriptionOptions struct {
	// ResourceTypes specifies which resource types to subscribe to
	ResourceTypes []string `json:"resource_types,omitempty"`

	// EventTypes specifies which event types to subscribe to
	EventTypes []EventType `json:"event_types,omitempty"`

	// Filter contains filter criteria for events
	Filter map[string]interface{} `json:"filter,omitempty"`

	// BufferSize is the size of the event channel buffer
	BufferSize int `json:"buffer_size,omitempty"`
}

// Event represents an event from an external system
type Event struct {
	// ID is the unique event identifier
	ID string `json:"id"`

	// Type is the event type
	Type EventType `json:"type"`

	// ResourceType is the type of resource that changed
	ResourceType string `json:"resource_type"`

	// ResourceID is the ID of the resource that changed
	ResourceID string `json:"resource_id"`

	// Timestamp when the event occurred
	Timestamp time.Time `json:"timestamp"`

	// Data contains the event data
	Data interface{} `json:"data,omitempty"`

	// Metadata contains additional event metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EventType represents the type of event
type EventType string

const (
	// EventTypeCreated indicates a resource was created
	EventTypeCreated EventType = "created"

	// EventTypeUpdated indicates a resource was updated
	EventTypeUpdated EventType = "updated"

	// EventTypeDeleted indicates a resource was deleted
	EventTypeDeleted EventType = "deleted"

	// EventTypeStatusChanged indicates a resource status changed
	EventTypeStatusChanged EventType = "status_changed"
)

// WebhookAdapter is an interface for adapters that support webhooks
type WebhookAdapter interface {
	Adapter

	// RegisterWebhook registers a webhook with the external system
	RegisterWebhook(ctx context.Context, webhook *WebhookConfig) (*Webhook, error)

	// UnregisterWebhook removes a webhook from the external system
	UnregisterWebhook(ctx context.Context, webhookID string) error

	// ListWebhooks lists all registered webhooks
	ListWebhooks(ctx context.Context) ([]*Webhook, error)

	// HandleWebhook processes an incoming webhook payload
	HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*Event, error)
}

// WebhookConfig represents the configuration for a webhook
type WebhookConfig struct {
	// URL is the endpoint URL to receive webhook events
	URL string `json:"url"`

	// Events specifies which events to receive
	Events []EventType `json:"events"`

	// Secret is the secret key for webhook signature verification
	Secret string `json:"secret,omitempty"`

	// Headers are custom headers to include in webhook requests
	Headers map[string]string `json:"headers,omitempty"`
}

// Webhook represents a registered webhook
type Webhook struct {
	// ID is the webhook identifier
	ID string `json:"id"`

	// Config is the webhook configuration
	Config *WebhookConfig `json:"config"`

	// Status is the webhook status
	Status WebhookStatus `json:"status"`

	// CreatedAt is when the webhook was created
	CreatedAt time.Time `json:"created_at"`

	// LastTriggeredAt is when the webhook was last triggered
	LastTriggeredAt *time.Time `json:"last_triggered_at,omitempty"`
}

// WebhookStatus represents the status of a webhook
type WebhookStatus string

const (
	// WebhookStatusActive indicates the webhook is active
	WebhookStatusActive WebhookStatus = "active"

	// WebhookStatusInactive indicates the webhook is inactive
	WebhookStatusInactive WebhookStatus = "inactive"

	// WebhookStatusFailed indicates the webhook has failed
	WebhookStatusFailed WebhookStatus = "failed"
)
