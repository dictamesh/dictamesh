// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package notifications provides a comprehensive multi-channel notification system
// for the DictaMesh framework, supporting both infrastructure alerts and application-level notifications.
package notifications

import (
	"time"
)

// Priority defines the urgency level of a notification
type Priority string

const (
	PriorityCritical Priority = "CRITICAL" // Immediate attention required
	PriorityHigh     Priority = "HIGH"     // Important, deliver quickly
	PriorityNormal   Priority = "NORMAL"   // Standard delivery
	PriorityLow      Priority = "LOW"      // Can be batched/delayed
)

// Channel represents a notification delivery channel
type Channel string

const (
	ChannelEmail       Channel = "EMAIL"
	ChannelSMS         Channel = "SMS"
	ChannelPush        Channel = "PUSH"
	ChannelSlack       Channel = "SLACK"
	ChannelWebhook     Channel = "WEBHOOK"
	ChannelInApp       Channel = "IN_APP"
	ChannelBrowserPush Channel = "BROWSER_PUSH"
	ChannelPagerDuty   Channel = "PAGERDUTY"
)

// Status represents the current state of a notification
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusQueued    Status = "QUEUED"
	StatusSending   Status = "SENDING"
	StatusSent      Status = "SENT"
	StatusDelivered Status = "DELIVERED"
	StatusFailed    Status = "FAILED"
	StatusRetrying  Status = "RETRYING"
	StatusCancelled Status = "CANCELLED"
)

// RecipientType defines the type of notification recipient
type RecipientType string

const (
	RecipientTypeUser   RecipientType = "USER"
	RecipientTypeRole   RecipientType = "ROLE"
	RecipientTypeGroup  RecipientType = "GROUP"
	RecipientTypeSystem RecipientType = "SYSTEM"
)

// Notification represents a notification instance
type Notification struct {
	ID string

	// Source tracking
	EventID    string
	RuleID     string
	TemplateID string

	// Recipient information
	RecipientType RecipientType
	RecipientID   string

	// Content
	Subject  string
	Body     string
	BodyHTML string
	Data     map[string]interface{}

	// Routing
	Priority        Priority
	Channels        []Channel
	SelectedChannel Channel

	// Status tracking
	Status Status

	// Timing
	ScheduledAt time.Time
	SentAt      *time.Time
	DeliveredAt *time.Time
	ReadAt      *time.Time

	// Error handling
	Error      string
	RetryCount int
	NextRetry  *time.Time

	// Metadata
	Metadata map[string]interface{}
	TraceID  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NotificationTemplate defines a reusable notification template
type NotificationTemplate struct {
	ID          string
	Name        string
	Description string

	// Multi-channel content
	Channels map[Channel]ChannelTemplate

	// Localization support
	Translations map[string]LocalizedTemplate

	// Template metadata
	Variables     []string
	SchemaVersion string

	// Lifecycle
	Version   string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string

	// Organization
	Tags []string
}

// ChannelTemplate defines channel-specific template content
type ChannelTemplate struct {
	Subject  string // For email, push title
	Body     string // Main content (supports templates)
	BodyHTML string // HTML version (email)
	Data     map[string]interface{}
}

// LocalizedTemplate provides localized versions of templates
type LocalizedTemplate struct {
	Subject  string
	Body     string
	BodyHTML string
}

// NotificationRule defines when and how to trigger notifications
type NotificationRule struct {
	ID          string
	Name        string
	Description string

	// Trigger conditions
	EventPattern string   // CEL expression
	Domains      []string
	EventTypes   []string

	// Routing configuration
	Priority         Priority
	Channels         []Channel
	FallbackChannels []Channel

	// Recipient selection
	RecipientSelector RecipientSelector

	// Timing
	Schedule *Schedule
	Timezone string

	// Batching configuration
	BatchWindow time.Duration
	BatchSize   int

	// Template reference
	TemplateID   string
	TemplateVars map[string]interface{}

	// Lifecycle
	Enabled    bool
	ValidFrom  time.Time
	ValidUntil *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// RecipientSelector defines how to select notification recipients
type RecipientSelector struct {
	Type string // role | user | group | dynamic

	// Static recipients
	UserIDs []string
	Roles   []string
	Groups  []string

	// Dynamic recipients (evaluated at runtime using CEL)
	Expression string
}

// Schedule defines when notifications should be delivered
type Schedule struct {
	Type string // immediate | cron | interval | time

	// Cron schedule (for recurring notifications)
	Cron string

	// Interval (for periodic notifications)
	Interval time.Duration

	// Specific time (for one-time scheduled notifications)
	Time *time.Time
}

// UserPreferences stores user notification preferences
type UserPreferences struct {
	UserID string

	// Global settings
	Enabled  bool
	Timezone string
	Locale   string

	// Contact information
	Email      string
	Phone      string
	PushTokens []string

	// Channel preferences
	ChannelPrefs map[Channel]ChannelPreference

	// Quiet hours configuration
	QuietHours QuietHours

	// Category-specific preferences
	CategoryPrefs map[string]CategoryPreference

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChannelPreference defines per-channel preferences
type ChannelPreference struct {
	Enabled bool
	Address string // Email, phone, or other channel-specific address
}

// QuietHours defines do-not-disturb periods
type QuietHours struct {
	Enabled       bool
	StartTime     string // Format: "HH:MM"
	EndTime       string // Format: "HH:MM"
	Timezone      string
	AllowCritical bool // Allow CRITICAL priority notifications
}

// CategoryPreference defines preferences for notification categories
type CategoryPreference struct {
	Enabled     bool
	Channels    []Channel
	MinPriority Priority
}

// DeliveryAttempt tracks individual delivery attempts
type DeliveryAttempt struct {
	ID             string
	NotificationID string

	// Delivery details
	Channel  Channel
	Provider string

	// Status
	Status        Status
	AttemptNumber int

	// Timing
	StartedAt   time.Time
	CompletedAt *time.Time

	// Result
	Success           bool
	Error             string
	ProviderResponse  map[string]interface{}
	ProviderMessageID string

	// Metadata
	Metadata map[string]interface{}
}

// NotificationBatch groups multiple notifications for batch delivery
type NotificationBatch struct {
	ID string

	// Batch configuration
	RuleID   string
	BatchKey string

	// Timing window
	WindowStart time.Time
	WindowEnd   time.Time
	ScheduledAt time.Time
	SentAt      *time.Time

	// Content
	NotificationIDs []string
	Count           int

	// Status
	Status Status

	CreatedAt time.Time
}

// RateLimit defines rate limiting configuration
type RateLimit struct {
	ID string

	Scope    string // user | system | category
	ScopeID  string
	Channel  Channel

	// Limit definition
	MaxCount      int
	WindowSeconds int

	// State
	Enabled bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NotificationEvent represents an event that triggers notifications
type NotificationEvent struct {
	EventID   string
	EventType string
	Timestamp time.Time

	// Event source
	Domain       string
	SourceSystem string

	// Event data
	Data map[string]interface{}

	// Tracing
	TraceID string
	SpanID  string
}

// SendNotificationRequest represents a request to send a notification
type SendNotificationRequest struct {
	// Recipient
	RecipientType RecipientType
	RecipientID   string

	// Priority and channels
	Priority Priority
	Channels []Channel

	// Content (either direct or via template)
	TemplateID   string
	TemplateVars map[string]interface{}
	Subject      string
	Body         string
	BodyHTML     string

	// Scheduling
	ScheduledAt *time.Time

	// Metadata
	Metadata map[string]interface{}
	TraceID  string
}

// BulkSendRequest represents a bulk notification request
type BulkSendRequest struct {
	Notifications []SendNotificationRequest
}

// BulkSendResponse represents the response for bulk send
type BulkSendResponse struct {
	TotalRequested int
	TotalAccepted  int
	TotalRejected  int
	Notifications  []Notification
	Errors         []error
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalSent        int64
	TotalDelivered   int64
	TotalFailed      int64
	ByChannel        map[Channel]ChannelStats
	ByPriority       map[Priority]PriorityStats
	AvgDeliveryTime  time.Duration
	SuccessRate      float64
	TimeRange        TimeRange
}

// ChannelStats represents statistics for a specific channel
type ChannelStats struct {
	Sent         int64
	Delivered    int64
	Failed       int64
	AvgLatency   time.Duration
	SuccessRate  float64
}

// PriorityStats represents statistics for a specific priority level
type PriorityStats struct {
	Sent         int64
	Delivered    int64
	Failed       int64
	AvgLatency   time.Duration
}

// TimeRange represents a time range for statistics
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// AuditEvent represents an audit log entry
type AuditEvent struct {
	ID             string
	NotificationID string

	// Event details
	EventType string
	ActorType string
	ActorID   string

	// Details
	Details map[string]interface{}

	// Timing
	Timestamp time.Time
	TraceID   string
}
