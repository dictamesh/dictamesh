// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package models provides database models for the metadata catalog
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB represents a PostgreSQL JSONB column
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
		return errors.New("failed to scan JSONB: value is not []byte")
	}

	return json.Unmarshal(bytes, j)
}

// EntityCatalog represents an entity in the catalog
type EntityCatalog struct {
	ID               string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EntityType       string    `gorm:"type:varchar(100);not null;index:idx_entity_type"`
	Domain           string    `gorm:"type:varchar(100);not null;index:idx_domain"`
	SourceSystem     string    `gorm:"type:varchar(100);not null;index:idx_source_system"`
	SourceEntityID   string    `gorm:"type:varchar(255);not null"`
	APIBaseURL       string    `gorm:"type:text;not null"`
	APIPathTemplate  string    `gorm:"type:text;not null"`
	APIMethod        string    `gorm:"type:varchar(10);default:'GET'"`
	APIAuthType      string    `gorm:"type:varchar(50)"`
	SchemaID         *string   `gorm:"type:uuid"`
	SchemaVersion    *string   `gorm:"type:varchar(50)"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	LastSeenAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status           string    `gorm:"type:varchar(50);default:'active';index:idx_status"`
	AvailabilitySLA  *float64  `gorm:"type:decimal(5,4)"`
	LatencyP99Ms     *int      `gorm:"type:integer"`
	FreshnessSLA     *int      `gorm:"type:integer"`
	ContainsPII      bool      `gorm:"default:false;index:idx_contains_pii,where:contains_pii = true"`
	DataClassification *string `gorm:"type:varchar(50);index:idx_data_classification"`
	RetentionDays    *int      `gorm:"type:integer"`
}

// TableName returns the table name
func (EntityCatalog) TableName() string {
	return "dictamesh_entity_catalog"
}

// EntityRelationship represents a relationship between entities
type EntityRelationship struct {
	ID                     string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SubjectCatalogID       string    `gorm:"type:uuid;not null"`
	SubjectEntityType      string    `gorm:"type:varchar(100);not null;index:idx_subject"`
	SubjectEntityID        string    `gorm:"type:varchar(255);not null;index:idx_subject"`
	RelationshipType       string    `gorm:"type:varchar(100);not null;index:idx_relationship_type"`
	RelationshipCardinality *string  `gorm:"type:varchar(20)"`
	ObjectCatalogID        string    `gorm:"type:uuid;not null"`
	ObjectEntityType       string    `gorm:"type:varchar(100);not null;index:idx_object"`
	ObjectEntityID         string    `gorm:"type:varchar(255);not null;index:idx_object"`
	ValidFrom              time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_temporal,where:valid_to IS NULL"`
	ValidTo                *time.Time
	SubjectDisplayName     *string `gorm:"type:varchar(255)"`
	ObjectDisplayName      *string `gorm:"type:varchar(255)"`
	RelationshipMetadata   JSONB   `gorm:"type:jsonb"`
	CreatedByEventID       *string `gorm:"type:varchar(255)"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	SubjectCatalog *EntityCatalog `gorm:"foreignKey:SubjectCatalogID"`
	ObjectCatalog  *EntityCatalog `gorm:"foreignKey:ObjectCatalogID"`
}

// TableName returns the table name
func (EntityRelationship) TableName() string {
	return "dictamesh_entity_relationships"
}

// Schema represents a versioned entity schema
type Schema struct {
	ID                 string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EntityType         string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_entity_version"`
	Version            string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_entity_version"`
	SchemaFormat       string    `gorm:"type:varchar(50);not null"`
	SchemaDefinition   JSONB     `gorm:"type:jsonb;not null"`
	BackwardCompatible bool      `gorm:"default:true"`
	ForwardCompatible  bool      `gorm:"default:false"`
	PublishedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeprecatedAt       *time.Time
	RetiredAt          *time.Time
}

// TableName returns the table name
func (Schema) TableName() string {
	return "dictamesh_schemas"
}

// EventLog represents an event in the audit log
type EventLog struct {
	ID             string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID        string    `gorm:"type:varchar(255);unique;not null"`
	EventType      string    `gorm:"type:varchar(100);not null;index:idx_event_type"`
	CatalogID      *string   `gorm:"type:uuid;index:idx_event_catalog"`
	EntityType     *string   `gorm:"type:varchar(100);index:idx_event_type"`
	EntityID       *string   `gorm:"type:varchar(255);index:idx_event_type"`
	ChangedFields  []string  `gorm:"type:text[]"`
	EventPayload   JSONB     `gorm:"type:jsonb"`
	TraceID        *string   `gorm:"type:varchar(64);index:idx_trace"`
	SpanID         *string   `gorm:"type:varchar(16)"`
	EventTimestamp time.Time `gorm:"not null;index:idx_event_timestamp"`
	IngestedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Catalog *EntityCatalog `gorm:"foreignKey:CatalogID"`
}

// TableName returns the table name
func (EventLog) TableName() string {
	return "dictamesh_event_log"
}

// DataLineage represents data lineage tracking
type DataLineage struct {
	ID                  string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UpstreamCatalogID   string    `gorm:"type:uuid;not null;index:idx_lineage_upstream"`
	UpstreamSystem      string    `gorm:"type:varchar(100)"`
	DownstreamCatalogID string    `gorm:"type:uuid;not null;index:idx_lineage_downstream"`
	DownstreamSystem    string    `gorm:"type:varchar(100)"`
	TransformationType  *string   `gorm:"type:varchar(50)"`
	TransformationLogic *string   `gorm:"type:text"`
	DataFlowActive      bool      `gorm:"default:true;index:idx_lineage_active,where:data_flow_active = true"`
	LastFlowAt          *time.Time
	AverageLatencyMs    *int      `gorm:"type:integer"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	UpstreamCatalog   *EntityCatalog `gorm:"foreignKey:UpstreamCatalogID"`
	DownstreamCatalog *EntityCatalog `gorm:"foreignKey:DownstreamCatalogID"`
}

// TableName returns the table name
func (DataLineage) TableName() string {
	return "dictamesh_data_lineage"
}

// CacheStatus represents cache status tracking
type CacheStatus struct {
	EntityCatalogID string    `gorm:"type:uuid;not null;primaryKey"`
	EntityID        string    `gorm:"type:varchar(255);not null;primaryKey"`
	CacheLayer      string    `gorm:"type:varchar(50);not null;primaryKey;index:idx_cache_layer"`
	CachedAt        time.Time `gorm:"not null"`
	ExpiresAt       *time.Time `gorm:"index:idx_cache_expiry"`
	CacheKey        *string   `gorm:"type:varchar(500)"`
	HitCount        int       `gorm:"default:0"`

	// Relations
	Catalog *EntityCatalog `gorm:"foreignKey:EntityCatalogID"`
}

// TableName returns the table name
func (CacheStatus) TableName() string {
	return "dictamesh_cache_status"
}
