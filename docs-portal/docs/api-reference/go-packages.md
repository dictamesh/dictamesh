<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 3
---

# Go Packages Reference

Complete reference for DictaMesh Go packages for building adapters and services.

## Installation

```bash
go get github.com/click2-run/dictamesh@latest
```

## Package Overview

DictaMesh provides the following core packages:

| Package | Description |
|---------|-------------|
| `pkg/adapter` | Adapter interface and base implementations |
| `pkg/catalog` | Metadata catalog client |
| `pkg/database` | Database models and repository |
| `pkg/events` | Event bus integration (Kafka) |
| `pkg/gateway` | GraphQL gateway components |
| `pkg/observability` | Tracing, metrics, and logging |
| `pkg/governance` | Policy enforcement and audit |
| `pkg/saga` | Saga pattern for distributed transactions |

## pkg/adapter

Core adapter interface and base implementations.

### Interface

```go
package adapter

import "context"

// DataProductAdapter defines the contract for all adapters
type DataProductAdapter interface {
    // GetEntity retrieves an entity by ID
    GetEntity(ctx context.Context, id string) (*Entity, error)

    // ListEntities lists entities with pagination
    ListEntities(ctx context.Context, opts ListOptions) ([]*Entity, error)

    // GetEntityType returns the entity type this adapter handles
    GetEntityType() string

    // GetDomain returns the domain this adapter belongs to
    GetDomain() string

    // GetSourceSystem returns the source system identifier
    GetSourceSystem() string

    // Health checks adapter health
    Health(ctx context.Context) error

    // Close cleanly shuts down the adapter
    Close() error
}

// Entity represents a domain entity
type Entity struct {
    ID           string
    Type         string
    Domain       string
    SourceSystem string
    SourceID     string
    Attributes   map[string]interface{}
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// ListOptions defines pagination and filtering
type ListOptions struct {
    Limit      int
    Offset     int
    Filters    map[string]interface{}
    SortBy     string
    SortOrder  string
}
```

### Base Adapter

```go
package adapter

import (
    "context"
    "github.com/click2-run/dictamesh/pkg/catalog"
    "github.com/click2-run/dictamesh/pkg/events"
)

// BaseAdapter provides common functionality
type BaseAdapter struct {
    entityType   string
    domain       string
    sourceSystem string
    catalog      *catalog.Client
    events       *events.Producer
    cache        Cache
    metrics      Metrics
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(config BaseConfig) *BaseAdapter {
    return &BaseAdapter{
        entityType:   config.EntityType,
        domain:       config.Domain,
        sourceSystem: config.SourceSystem,
        catalog:      config.Catalog,
        events:       config.Events,
        cache:        config.Cache,
        metrics:      config.Metrics,
    }
}

// GetEntityType returns the entity type
func (a *BaseAdapter) GetEntityType() string {
    return a.entityType
}

// GetDomain returns the domain
func (a *BaseAdapter) GetDomain() string {
    return a.domain
}

// GetSourceSystem returns the source system
func (a *BaseAdapter) GetSourceSystem() string {
    return a.sourceSystem
}

// PublishEvent publishes an event
func (a *BaseAdapter) PublishEvent(ctx context.Context, eventType string, entity *Entity) error {
    event := &events.Event{
        Type:      eventType,
        EntityID:  entity.ID,
        EntityType: entity.Type,
        Payload:   entity.Attributes,
        Timestamp: time.Now(),
    }
    return a.events.Publish(ctx, event)
}

// RegisterInCatalog registers entity in metadata catalog
func (a *BaseAdapter) RegisterInCatalog(ctx context.Context, entity *Entity) error {
    return a.catalog.RegisterEntity(ctx, &catalog.EntityRegistration{
        EntityType:   entity.Type,
        Domain:       entity.Domain,
        SourceSystem: entity.SourceSystem,
        SourceEntityID: entity.SourceID,
    })
}
```

### Example Usage

```go
package main

import (
    "context"
    "github.com/click2-run/dictamesh/pkg/adapter"
)

// ProductAdapter implements DataProductAdapter
type ProductAdapter struct {
    *adapter.BaseAdapter
    client *ShopifyClient
}

func NewProductAdapter(cfg Config) (*ProductAdapter, error) {
    base := adapter.NewBaseAdapter(adapter.BaseConfig{
        EntityType:   "product",
        Domain:       "ecommerce",
        SourceSystem: "shopify",
        Catalog:      cfg.Catalog,
        Events:       cfg.Events,
    })

    return &ProductAdapter{
        BaseAdapter: base,
        client:      NewShopifyClient(cfg.ShopifyURL, cfg.ShopifyToken),
    }, nil
}

func (a *ProductAdapter) GetEntity(ctx context.Context, id string) (*adapter.Entity, error) {
    // Fetch from Shopify
    product, err := a.client.GetProduct(ctx, id)
    if err != nil {
        return nil, err
    }

    // Transform to entity
    entity := &adapter.Entity{
        ID:           product.ID,
        Type:         "product",
        Domain:       "ecommerce",
        SourceSystem: "shopify",
        SourceID:     product.ID,
        Attributes: map[string]interface{}{
            "name":        product.Title,
            "description": product.Description,
            "price":       product.Price,
            "sku":         product.SKU,
        },
        CreatedAt: product.CreatedAt,
        UpdatedAt: product.UpdatedAt,
    }

    // Register in catalog
    if err := a.RegisterInCatalog(ctx, entity); err != nil {
        return nil, err
    }

    // Publish event
    if err := a.PublishEvent(ctx, "entity.fetched", entity); err != nil {
        return nil, err
    }

    return entity, nil
}

func (a *ProductAdapter) ListEntities(ctx context.Context, opts adapter.ListOptions) ([]*adapter.Entity, error) {
    // Implementation
    return nil, nil
}

func (a *ProductAdapter) Health(ctx context.Context) error {
    return a.client.Ping(ctx)
}

func (a *ProductAdapter) Close() error {
    return a.client.Close()
}
```

## pkg/catalog

Metadata catalog client for entity registration and discovery.

### Client

```go
package catalog

import "context"

// Client provides access to metadata catalog
type Client struct {
    baseURL string
    token   string
}

// NewClient creates a new catalog client
func NewClient(baseURL, token string) *Client {
    return &Client{
        baseURL: baseURL,
        token:   token,
    }
}

// RegisterEntity registers an entity in the catalog
func (c *Client) RegisterEntity(ctx context.Context, reg *EntityRegistration) error

// GetEntity retrieves entity metadata
func (c *Client) GetEntity(ctx context.Context, id string) (*EntityMetadata, error)

// FindBySource finds entity by source system and ID
func (c *Client) FindBySource(ctx context.Context, sourceSystem, sourceID, entityType string) (*EntityMetadata, error)

// ListEntities lists entities with filters
func (c *Client) ListEntities(ctx context.Context, filters *EntityFilters) ([]*EntityMetadata, error)

// CreateRelationship creates a relationship between entities
func (c *Client) CreateRelationship(ctx context.Context, rel *Relationship) error

// GetRelationships gets relationships for an entity
func (c *Client) GetRelationships(ctx context.Context, entityID string, direction RelationshipDirection) ([]*Relationship, error)
```

### Types

```go
// EntityRegistration represents entity registration request
type EntityRegistration struct {
    EntityType         string
    Domain             string
    SourceSystem       string
    SourceEntityID     string
    APIBaseURL         string
    APIPathTemplate    string
    APIMethod          string
    SchemaID           string
    SchemaVersion      string
    ContainsPII        bool
    DataClassification string
}

// EntityMetadata represents entity metadata
type EntityMetadata struct {
    ID             string
    EntityType     string
    Domain         string
    SourceSystem   string
    SourceEntityID string
    Status         string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// Relationship represents entity relationship
type Relationship struct {
    SubjectID        string
    SubjectType      string
    RelationshipType string
    ObjectID         string
    ObjectType       string
    Metadata         map[string]interface{}
}

type RelationshipDirection string

const (
    DirectionOutgoing RelationshipDirection = "outgoing"
    DirectionIncoming RelationshipDirection = "incoming"
    DirectionAll      RelationshipDirection = "all"
)
```

## pkg/database

Database models and repository patterns.

### Models

```go
package models

// EntityCatalog represents entity catalog entry
type EntityCatalog struct {
    ID                 string
    EntityType         string
    Domain             string
    SourceSystem       string
    SourceEntityID     string
    APIBaseURL         string
    APIPathTemplate    string
    APIMethod          string
    SchemaID           *string
    SchemaVersion      *string
    CreatedAt          time.Time
    UpdatedAt          time.Time
    Status             string
    AvailabilitySLA    *float64
    LatencyP99Ms       *int
    ContainsPII        bool
    DataClassification *string
}

// EntityRelationship represents entity relationship
type EntityRelationship struct {
    ID                      string
    SubjectCatalogID        string
    SubjectEntityType       string
    SubjectEntityID         string
    RelationshipType        string
    RelationshipCardinality *string
    ObjectCatalogID         string
    ObjectEntityType        string
    ObjectEntityID          string
    ValidFrom               time.Time
    ValidTo                 *time.Time
    RelationshipMetadata    JSONB
    CreatedAt               time.Time
}

// Schema represents versioned schema
type Schema struct {
    ID                 string
    EntityType         string
    Version            string
    SchemaFormat       string
    SchemaDefinition   JSONB
    BackwardCompatible bool
    ForwardCompatible  bool
    PublishedAt        time.Time
}
```

### Repository

```go
package repository

import "context"

// CatalogRepository provides catalog access
type CatalogRepository struct {
    db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) *CatalogRepository

func (r *CatalogRepository) Create(ctx context.Context, entity *models.EntityCatalog) error

func (r *CatalogRepository) FindByID(ctx context.Context, id string) (*models.EntityCatalog, error)

func (r *CatalogRepository) FindBySource(ctx context.Context, sourceSystem, sourceEntityID, entityType string) (*models.EntityCatalog, error)

func (r *CatalogRepository) List(ctx context.Context, filters *CatalogFilters) ([]models.EntityCatalog, error)

func (r *CatalogRepository) Update(ctx context.Context, entity *models.EntityCatalog) error

func (r *CatalogRepository) Delete(ctx context.Context, id string) error

// RelationshipRepository provides relationship access
type RelationshipRepository struct {
    db *gorm.DB
}

func NewRelationshipRepository(db *gorm.DB) *RelationshipRepository

func (r *RelationshipRepository) Create(ctx context.Context, rel *models.EntityRelationship) error

func (r *RelationshipRepository) FindBySubject(ctx context.Context, entityType, entityID string) ([]models.EntityRelationship, error)

func (r *RelationshipRepository) FindByObject(ctx context.Context, entityType, entityID string) ([]models.EntityRelationship, error)
```

## pkg/events

Kafka event bus integration.

### Producer

```go
package events

import "context"

// Producer publishes events to Kafka
type Producer struct {
    kafka    *kafka.Writer
    schemaRegistry *SchemaRegistry
}

// NewProducer creates a new event producer
func NewProducer(config ProducerConfig) (*Producer, error)

// Publish publishes an event
func (p *Producer) Publish(ctx context.Context, event *Event) error

// PublishBatch publishes multiple events
func (p *Producer) PublishBatch(ctx context.Context, events []*Event) error

// Close closes the producer
func (p *Producer) Close() error

// Event represents an event
type Event struct {
    Type       string
    EntityID   string
    EntityType string
    Payload    map[string]interface{}
    Timestamp  time.Time
    TraceID    string
    SpanID     string
}
```

### Consumer

```go
package events

// Consumer consumes events from Kafka
type Consumer struct {
    kafka    *kafka.Reader
    handlers map[string]EventHandler
}

// NewConsumer creates a new event consumer
func NewConsumer(config ConsumerConfig) (*Consumer, error)

// Subscribe registers handler for event type
func (c *Consumer) Subscribe(eventType string, handler EventHandler) error

// Start starts consuming events
func (c *Consumer) Start(ctx context.Context) error

// Close closes the consumer
func (c *Consumer) Close() error

// EventHandler handles events
type EventHandler func(ctx context.Context, event *Event) error
```

### Example

```go
package main

import (
    "context"
    "github.com/click2-run/dictamesh/pkg/events"
)

func main() {
    // Create producer
    producer, err := events.NewProducer(events.ProducerConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "dictamesh.entities",
    })
    if err != nil {
        panic(err)
    }
    defer producer.Close()

    // Publish event
    event := &events.Event{
        Type:       "entity.created",
        EntityID:   "prod-12345",
        EntityType: "product",
        Payload: map[string]interface{}{
            "name":  "Premium Headphones",
            "price": 299.99,
        },
        Timestamp: time.Now(),
    }

    if err := producer.Publish(context.Background(), event); err != nil {
        panic(err)
    }

    // Create consumer
    consumer, err := events.NewConsumer(events.ConsumerConfig{
        Brokers:  []string{"localhost:9092"},
        Topic:    "dictamesh.entities",
        GroupID:  "product-sync",
    })
    if err != nil {
        panic(err)
    }
    defer consumer.Close()

    // Subscribe to events
    consumer.Subscribe("entity.created", func(ctx context.Context, event *events.Event) error {
        log.Printf("Received event: %s for %s", event.Type, event.EntityID)
        return nil
    })

    // Start consuming
    if err := consumer.Start(context.Background()); err != nil {
        panic(err)
    }
}
```

## pkg/observability

Tracing, metrics, and logging.

### Tracing

```go
package observability

import (
    "context"
    "go.opentelemetry.io/otel/trace"
)

// Tracer provides distributed tracing
type Tracer struct {
    tracer trace.Tracer
}

// NewTracer creates a new tracer
func NewTracer(serviceName string) (*Tracer, error)

// StartSpan starts a new span
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)

// AddEvent adds an event to the current span
func (t *Tracer) AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue)

// RecordError records an error on the current span
func (t *Tracer) RecordError(ctx context.Context, err error)
```

### Metrics

```go
package observability

import "github.com/prometheus/client_golang/prometheus"

// Metrics provides Prometheus metrics
type Metrics struct {
    requestCounter   *prometheus.CounterVec
    requestDuration  *prometheus.HistogramVec
    entityGauge      *prometheus.GaugeVec
}

// NewMetrics creates metrics
func NewMetrics(namespace string) *Metrics

// RecordRequest records a request
func (m *Metrics) RecordRequest(method, path string, status int, duration time.Duration)

// RecordEntityCount records entity count
func (m *Metrics) RecordEntityCount(entityType string, count float64)

// Handler returns Prometheus HTTP handler
func (m *Metrics) Handler() http.Handler
```

### Logger

```go
package observability

import "go.uber.org/zap"

// Logger provides structured logging
type Logger struct {
    logger *zap.Logger
}

// NewLogger creates a new logger
func NewLogger(config LoggerConfig) (*Logger, error)

// Info logs info message
func (l *Logger) Info(msg string, fields ...zap.Field)

// Error logs error message
func (l *Logger) Error(msg string, err error, fields ...zap.Field)

// With adds fields to logger
func (l *Logger) With(fields ...zap.Field) *Logger

// WithContext adds context fields
func (l *Logger) WithContext(ctx context.Context) *Logger
```

## pkg/governance

Policy enforcement and audit.

### Policy

```go
package governance

// PolicyEngine enforces governance policies
type PolicyEngine struct {
    policies []Policy
}

// NewPolicyEngine creates a policy engine
func NewPolicyEngine(policies []Policy) *PolicyEngine

// Evaluate evaluates policies
func (e *PolicyEngine) Evaluate(ctx context.Context, action Action, resource Resource) (Decision, error)

// Policy represents a governance policy
type Policy interface {
    Evaluate(ctx context.Context, action Action, resource Resource) (Decision, error)
}

// Decision represents policy decision
type Decision struct {
    Allowed bool
    Reason  string
}
```

### Audit

```go
package governance

// AuditLogger logs audit events
type AuditLogger struct {
    db *gorm.DB
}

// NewAuditLogger creates audit logger
func NewAuditLogger(db *gorm.DB) *AuditLogger

// Log logs an audit event
func (a *AuditLogger) Log(ctx context.Context, event *AuditEvent) error

// AuditEvent represents audit event
type AuditEvent struct {
    Actor      string
    Action     string
    Resource   string
    Result     string
    Timestamp  time.Time
    TraceID    string
}
```

## Next Steps

- [REST API Reference](./rest-api.md) - REST API for metadata catalog
- [GraphQL API Reference](./graphql-api.md) - Query the unified graph
- [Event Schemas Reference](./event-schemas.md) - Event-driven integration
- [Building Adapters Guide](../guides/building-adapters.md) - Create your first adapter

---

**Previous**: [← GraphQL API](./graphql-api.md) | **Next**: [Event Schemas →](./event-schemas.md)
