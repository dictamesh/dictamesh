// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package config

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Configuration represents a single configuration entry
type Configuration struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Hierarchy
	Environment string  `gorm:"type:varchar(50);not null;index:idx_config_env_svc" json:"environment"`
	Service     string  `gorm:"type:varchar(100);not null;index:idx_config_env_svc" json:"service"`
	Component   *string `gorm:"type:varchar(100);index:idx_config_component" json:"component,omitempty"`
	Key         string  `gorm:"type:varchar(255);not null" json:"key"`

	// Value and metadata
	Value     json.RawMessage `gorm:"type:jsonb;not null" json:"value"`
	ValueType string          `gorm:"type:varchar(50);not null" json:"value_type"` // string, number, boolean, object, array
	IsSecret  bool            `gorm:"default:false;index:idx_config_is_secret" json:"is_secret"`

	// Validation and documentation
	Schema      json.RawMessage `gorm:"type:jsonb" json:"schema,omitempty"`
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	Tags        []string        `gorm:"type:text[]" json:"tags,omitempty"`

	// Versioning
	Version int `gorm:"not null;default:1" json:"version"`

	// Audit fields
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	CreatedBy *string   `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:now();index:idx_config_updated_at" json:"updated_at"`
	UpdatedBy *string   `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
}

// TableName overrides the default table name
func (Configuration) TableName() string {
	return "dictamesh_configurations"
}

// BeforeCreate hook to set defaults
func (c *Configuration) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.Version == 0 {
		c.Version = 1
	}
	return nil
}

// ConfigurationVersion represents a historical version of a configuration
type ConfigurationVersion struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConfigID uuid.UUID `gorm:"type:uuid;not null;index:idx_config_version_config" json:"config_id"`

	// Version information
	Version  int             `gorm:"not null" json:"version"`
	Value    json.RawMessage `gorm:"type:jsonb;not null" json:"value"`
	IsSecret bool            `gorm:"default:false" json:"is_secret"`

	// Change metadata
	ChangeDescription *string `gorm:"type:text" json:"change_description,omitempty"`

	// Audit fields
	CreatedAt time.Time `gorm:"not null;default:now();index:idx_config_version_created" json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(255);not null" json:"created_by"`

	// Relation
	Configuration *Configuration `gorm:"foreignKey:ConfigID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName overrides the default table name
func (ConfigurationVersion) TableName() string {
	return "dictamesh_config_versions"
}

// BeforeCreate hook
func (cv *ConfigurationVersion) BeforeCreate(tx *gorm.DB) error {
	if cv.ID == uuid.Nil {
		cv.ID = uuid.New()
	}
	return nil
}

// ConfigAuditLog represents an audit log entry for configuration operations
type ConfigAuditLog struct {
	ID       uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConfigID *uuid.UUID `gorm:"type:uuid;index:idx_audit_config" json:"config_id,omitempty"`

	// Action information
	Action string `gorm:"type:varchar(50);not null;index:idx_audit_action" json:"action"` // CREATE, UPDATE, DELETE, ACCESS, ROLLBACK, EXPORT, IMPORT

	// Actor information
	Actor     string  `gorm:"type:varchar(255);not null;index:idx_audit_actor" json:"actor"`
	ActorType *string `gorm:"type:varchar(50)" json:"actor_type,omitempty"` // USER, SERVICE, API_KEY, SYSTEM

	// Request metadata
	IPAddress *string `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent *string `gorm:"type:text" json:"user_agent,omitempty"`
	RequestID *string `gorm:"type:varchar(100);index:idx_audit_request" json:"request_id,omitempty"`

	// Change details
	Changes  json.RawMessage `gorm:"type:jsonb" json:"changes,omitempty"`
	Metadata json.RawMessage `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Timestamp
	Timestamp time.Time `gorm:"not null;default:now();index:idx_audit_timestamp" json:"timestamp"`

	// Relation
	Configuration *Configuration `gorm:"foreignKey:ConfigID;constraint:OnDelete:SET NULL" json:"-"`
}

// TableName overrides the default table name
func (ConfigAuditLog) TableName() string {
	return "dictamesh_config_audit_logs"
}

// BeforeCreate hook
func (cal *ConfigAuditLog) BeforeCreate(tx *gorm.DB) error {
	if cal.ID == uuid.Nil {
		cal.ID = uuid.New()
	}
	return nil
}

// EncryptionKey represents an encryption key for secrets management
type EncryptionKey struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Key identification
	KeyName     string  `gorm:"type:varchar(100);unique;not null" json:"key_name"`
	KeyType     string  `gorm:"type:varchar(50);not null" json:"key_type"` // MASTER, DATA_ENCRYPTION_KEY
	Environment *string `gorm:"type:varchar(50);index:idx_encryption_key_env" json:"environment,omitempty"`

	// Key material (encrypted with KEK from environment)
	EncryptedKey []byte `gorm:"type:bytea;not null" json:"-"` // Never expose in JSON
	Algorithm    string `gorm:"type:varchar(50);not null;default:'AES-256-GCM'" json:"algorithm"`
	KeyVersion   int    `gorm:"not null;default:1" json:"key_version"`

	// Key lifecycle
	IsActive    bool       `gorm:"default:true;index:idx_encryption_key_active" json:"is_active"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	RotatedAt   *time.Time `json:"rotated_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`

	// Metadata
	Metadata json.RawMessage `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// TableName overrides the default table name
func (EncryptionKey) TableName() string {
	return "dictamesh_encryption_keys"
}

// BeforeCreate hook
func (ek *EncryptionKey) BeforeCreate(tx *gorm.DB) error {
	if ek.ID == uuid.Nil {
		ek.ID = uuid.New()
	}
	if ek.KeyVersion == 0 {
		ek.KeyVersion = 1
	}
	return nil
}

// ConfigWatcher represents an active configuration watcher for hot reload
type ConfigWatcher struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Watcher identification
	ServiceInstance string `gorm:"type:varchar(255);not null" json:"service_instance"`
	ServiceName     string `gorm:"type:varchar(100);not null;index:idx_watcher_service" json:"service_name"`

	// Watch specification
	Environment  string `gorm:"type:varchar(50);not null;index:idx_watcher_service" json:"environment"`
	WatchPattern string `gorm:"type:varchar(255);not null" json:"watch_pattern"` // e.g., "database.*", "cache.redis.*"

	// Watcher status
	LastHeartbeat time.Time `gorm:"not null;default:now();index:idx_watcher_heartbeat" json:"last_heartbeat"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`

	// Registration
	RegisteredAt time.Time       `gorm:"not null;default:now()" json:"registered_at"`
	Metadata     json.RawMessage `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// TableName overrides the default table name
func (ConfigWatcher) TableName() string {
	return "dictamesh_config_watchers"
}

// BeforeCreate hook
func (cw *ConfigWatcher) BeforeCreate(tx *gorm.DB) error {
	if cw.ID == uuid.Nil {
		cw.ID = uuid.New()
	}
	return nil
}

// ValueType constants
const (
	ValueTypeString  = "string"
	ValueTypeNumber  = "number"
	ValueTypeBoolean = "boolean"
	ValueTypeObject  = "object"
	ValueTypeArray   = "array"
)

// Action constants for audit logs
const (
	ActionCreate   = "CREATE"
	ActionUpdate   = "UPDATE"
	ActionDelete   = "DELETE"
	ActionAccess   = "ACCESS"
	ActionRollback = "ROLLBACK"
	ActionExport   = "EXPORT"
	ActionImport   = "IMPORT"
)

// ActorType constants
const (
	ActorTypeUser    = "USER"
	ActorTypeService = "SERVICE"
	ActorTypeAPIKey  = "API_KEY"
	ActorTypeSystem  = "SYSTEM"
)

// KeyType constants
const (
	KeyTypeMaster             = "MASTER"
	KeyTypeDataEncryptionKey  = "DATA_ENCRYPTION_KEY"
)

// Environment constants
const (
	EnvironmentDev        = "dev"
	EnvironmentDevelopment = "development"
	EnvironmentStaging    = "staging"
	EnvironmentProduction = "production"
	EnvironmentTest       = "test"
)
