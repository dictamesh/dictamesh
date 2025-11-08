// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package adapter provides the core adapter interface and base implementation
// for building data product adapters in the DictaMesh framework.
package adapter

import (
	"context"
	"time"
)

// DataProductAdapter is the core interface that all adapters must implement.
// It defines the contract for integrating external data sources into the data mesh.
type DataProductAdapter interface {
	// Lifecycle methods
	Initialize(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Health(ctx context.Context) error

	// Metadata methods
	Name() string
	Version() string
	Description() string
	SourceSystem() string
	Domain() string

	// Entity operations
	GetEntity(ctx context.Context, entityType, id string) (*Entity, error)
	ListEntities(ctx context.Context, entityType string, opts ListOptions) ([]*Entity, error)
	CreateEntity(ctx context.Context, entityType string, data map[string]interface{}) (*Entity, error)
	UpdateEntity(ctx context.Context, entityType, id string, data map[string]interface{}) (*Entity, error)
	DeleteEntity(ctx context.Context, entityType, id string) error

	// Relationship operations
	GetRelationships(ctx context.Context, entityType, id string) ([]*Relationship, error)

	// Schema operations
	GetSchema(ctx context.Context, entityType string) (*Schema, error)
	ListSchemas(ctx context.Context) ([]*Schema, error)

	// Cache operations
	InvalidateCache(ctx context.Context, entityType, id string) error
}

// Entity represents a data mesh entity
type Entity struct {
	// Identity
	ID         string `json:"id"`
	Type       string `json:"type"`
	Domain     string `json:"domain"`
	Source     string `json:"source"`

	// Data
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]string      `json:"metadata"`

	// Versioning
	Version      int       `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Caching
	CacheKey     string    `json:"-"`
	CacheTTL     int       `json:"-"`
}

// Relationship represents a relationship between entities
type Relationship struct {
	FromEntityType string                 `json:"from_entity_type"`
	FromEntityID   string                 `json:"from_entity_id"`
	ToEntityType   string                 `json:"to_entity_type"`
	ToEntityID     string                 `json:"to_entity_id"`
	Type           string                 `json:"type"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// Schema represents an entity schema
type Schema struct {
	EntityType  string                 `json:"entity_type"`
	Version     int                    `json:"version"`
	Format      string                 `json:"format"` // "avro", "json", "protobuf"
	Definition  map[string]interface{} `json:"definition"`
	Description string                 `json:"description"`
}

// ListOptions provides options for listing entities
type ListOptions struct {
	Limit      int
	Offset     int
	Filter     map[string]interface{}
	Sort       string
	SortOrder  string // "asc" or "desc"
}

// AdapterCapabilities describes what an adapter can do
type AdapterCapabilities struct {
	SupportsRead   bool
	SupportsWrite  bool
	SupportsDelete bool
	SupportsCache  bool
	SupportsEvents bool
	SupportsStream bool
}

// AdapterMetrics holds adapter performance metrics
type AdapterMetrics struct {
	RequestsTotal     int64
	RequestsSucceeded int64
	RequestsFailed    int64
	CacheHits         int64
	CacheMisses       int64
	AvgLatencyMs      float64
	LastError         error
	LastErrorTime     time.Time
}

// AdapterStatus represents the current status of an adapter
type AdapterStatus string

const (
	StatusUninitialized AdapterStatus = "uninitialized"
	StatusInitializing  AdapterStatus = "initializing"
	StatusReady         AdapterStatus = "ready"
	StatusDegraded      AdapterStatus = "degraded"
	StatusFailed        AdapterStatus = "failed"
	StatusStopped       AdapterStatus = "stopped"
)
