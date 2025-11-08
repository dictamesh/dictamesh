// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// NotificationModel represents the database model for notifications
type NotificationModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Source tracking
	EventID    string     `gorm:"type:varchar(255);index"`
	RuleID     *uuid.UUID `gorm:"type:uuid"`
	TemplateID *uuid.UUID `gorm:"type:uuid"`

	// Recipient information
	RecipientType string `gorm:"type:varchar(50);not null;index:idx_recipient"`
	RecipientID   string `gorm:"type:varchar(255);not null;index:idx_recipient"`

	// Content
	Subject  string         `gorm:"type:text"`
	Body     string         `gorm:"type:text"`
	BodyHTML string         `gorm:"type:text"`
	Data     JSONB          `gorm:"type:jsonb"`

	// Routing
	Priority        string        `gorm:"type:varchar(20);not null"`
	Channels        StringArray   `gorm:"type:text[]"`
	SelectedChannel string        `gorm:"type:varchar(50)"`

	// Status tracking
	Status string `gorm:"type:varchar(20);not null;default:'pending';index:idx_status"`

	// Timing
	ScheduledAt time.Time  `gorm:"not null;default:now();index:idx_status"`
	SentAt      *time.Time
	DeliveredAt *time.Time
	ReadAt      *time.Time

	// Metadata
	Metadata JSONB  `gorm:"type:jsonb"`
	TraceID  string `gorm:"type:varchar(64);index"`

	// Error tracking
	Error       string     `gorm:"type:text"`
	RetryCount  int        `gorm:"default:0"`
	NextRetryAt *time.Time

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (NotificationModel) TableName() string {
	return "dictamesh_notifications"
}

// TemplateModel represents the database model for notification templates
type TemplateModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(255);not null;unique;index"`
	Description string    `gorm:"type:text"`

	// Content (JSONB for flexibility)
	Channels     JSONB `gorm:"type:jsonb;not null"`
	Translations JSONB `gorm:"type:jsonb"`

	// Template metadata
	Variables     JSONB  `gorm:"type:jsonb"`
	SchemaVersion string `gorm:"type:varchar(50);default:'1.0'"`

	// Lifecycle
	Version   string    `gorm:"type:varchar(50);default:'1.0.0'"`
	Enabled   bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
	CreatedBy string    `gorm:"type:varchar(255)"`

	// Organization
	Tags StringArray `gorm:"type:text[]"`
}

// TableName overrides the table name for GORM
func (TemplateModel) TableName() string {
	return "dictamesh_notification_templates"
}

// RuleModel represents the database model for notification rules
type RuleModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(255);not null;unique;index"`
	Description string    `gorm:"type:text"`

	// Trigger conditions
	EventPattern string      `gorm:"type:text;not null"`
	Domains      StringArray `gorm:"type:text[]"`
	EventTypes   StringArray `gorm:"type:text[]"`

	// Routing
	Priority         string      `gorm:"type:varchar(20);not null"`
	Channels         StringArray `gorm:"type:text[];not null"`
	FallbackChannels StringArray `gorm:"type:text[]"`

	// Recipients
	RecipientSelector JSONB `gorm:"type:jsonb;not null"`

	// Timing
	Schedule JSONB  `gorm:"type:jsonb"`
	Timezone string `gorm:"type:varchar(50);default:'UTC'"`

	// Batching
	BatchWindowSeconds *int `gorm:"type:integer"`
	BatchSize          *int `gorm:"type:integer"`

	// Template
	TemplateID   *uuid.UUID `gorm:"type:uuid"`
	TemplateVars JSONB      `gorm:"type:jsonb"`

	// Lifecycle
	Enabled    bool       `gorm:"default:true;index"`
	ValidFrom  time.Time  `gorm:"not null;default:now()"`
	ValidUntil *time.Time

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (RuleModel) TableName() string {
	return "dictamesh_notification_rules"
}

// DeliveryModel represents the database model for delivery attempts
type DeliveryModel struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NotificationID uuid.UUID `gorm:"type:uuid;not null;index:idx_delivery_notification"`

	// Delivery details
	Channel  string `gorm:"type:varchar(50);not null"`
	Provider string `gorm:"type:varchar(100)"`

	// Status
	Status        string `gorm:"type:varchar(20);not null"`
	AttemptNumber int    `gorm:"not null"`

	// Timing
	StartedAt   time.Time  `gorm:"not null;default:now()"`
	CompletedAt *time.Time

	// Result
	Success           bool   `gorm:"default:false"`
	Error             string `gorm:"type:text"`
	ProviderResponse  JSONB  `gorm:"type:jsonb"`
	ProviderMessageID string `gorm:"type:varchar(255);index"`

	// Metadata
	Metadata JSONB `gorm:"type:jsonb"`
}

// TableName overrides the table name for GORM
func (DeliveryModel) TableName() string {
	return "dictamesh_notification_delivery"
}

// PreferencesModel represents the database model for user preferences
type PreferencesModel struct {
	UserID string `gorm:"type:varchar(255);primary_key"`

	// Global settings
	Enabled  bool   `gorm:"default:true"`
	Timezone string `gorm:"type:varchar(50);default:'UTC'"`
	Locale   string `gorm:"type:varchar(10);default:'en'"`

	// Channel addresses
	Email      string `gorm:"type:varchar(255);index"`
	Phone      string `gorm:"type:varchar(20);index"`
	PushTokens JSONB  `gorm:"type:jsonb"`

	// Channel preferences
	ChannelPrefs JSONB `gorm:"type:jsonb;default:'{}'"`

	// Quiet hours
	QuietHoursEnabled      bool `gorm:"default:false"`
	QuietHoursStart        *time.Time
	QuietHoursEnd          *time.Time
	QuietHoursAllowCritical bool `gorm:"default:true"`

	// Category preferences
	CategoryPrefs JSONB `gorm:"type:jsonb;default:'{}'"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (PreferencesModel) TableName() string {
	return "dictamesh_notification_preferences"
}

// BatchModel represents the database model for notification batches
type BatchModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Batch config
	RuleID   *uuid.UUID `gorm:"type:uuid"`
	BatchKey string     `gorm:"type:varchar(255);not null;index:idx_batch_key_window"`

	// Timing
	WindowStart time.Time  `gorm:"not null"`
	WindowEnd   time.Time  `gorm:"not null;index:idx_batch_key_window"`
	ScheduledAt time.Time  `gorm:"not null;index:idx_batch_scheduled"`
	SentAt      *time.Time

	// Content
	NotificationIDs UUIDArray `gorm:"type:uuid[]"`
	Count           int       `gorm:"not null"`

	// Status
	Status string `gorm:"type:varchar(20);default:'pending';index:idx_batch_scheduled"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (BatchModel) TableName() string {
	return "dictamesh_notification_batches"
}

// RateLimitModel represents the database model for rate limit configuration
type RateLimitModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	Scope    string  `gorm:"type:varchar(50);not null"`
	ScopeID  *string `gorm:"type:varchar(255)"`
	Channel  string  `gorm:"type:varchar(50);not null"`

	// Limit definition
	MaxCount      int `gorm:"not null"`
	WindowSeconds int `gorm:"not null"`

	// Metadata
	Enabled bool `gorm:"default:true"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (RateLimitModel) TableName() string {
	return "dictamesh_notification_rate_limits"
}

// AuditModel represents the database model for audit logs
type AuditModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NotificationID *uuid.UUID `gorm:"type:uuid;index:idx_audit_notification"`

	// Event details
	EventType string  `gorm:"type:varchar(100);not null;index:idx_audit_type"`
	ActorType *string `gorm:"type:varchar(50)"`
	ActorID   *string `gorm:"type:varchar(255)"`

	// Details
	Details JSONB `gorm:"type:jsonb"`

	// Timing
	Timestamp time.Time `gorm:"not null;default:now();index:idx_audit_notification;index:idx_audit_type"`
	TraceID   string    `gorm:"type:varchar(64)"`
}

// TableName overrides the table name for GORM
func (AuditModel) TableName() string {
	return "dictamesh_notification_audit"
}

// JSONB is a custom type for JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// StringArray is a custom type for text[] columns
type StringArray []string

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, s)
}

// UUIDArray is a custom type for uuid[] columns
type UUIDArray []uuid.UUID

// Value implements the driver.Valuer interface
func (u UUIDArray) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}

	strArray := make([]string, len(u))
	for i, id := range u {
		strArray[i] = id.String()
	}

	return json.Marshal(strArray)
}

// Scan implements the sql.Scanner interface
func (u *UUIDArray) Scan(value interface{}) error {
	if value == nil {
		*u = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var strArray []string
	if err := json.Unmarshal(bytes, &strArray); err != nil {
		return err
	}

	uuids := make([]uuid.UUID, len(strArray))
	for i, str := range strArray {
		id, err := uuid.Parse(str)
		if err != nil {
			return err
		}
		uuids[i] = id
	}

	*u = uuids
	return nil
}
