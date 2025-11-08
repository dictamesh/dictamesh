# DictaMesh Adapter Package

Core adapter interface and base implementation for building data product adapters in the DictaMesh framework.

## Overview

The adapter package provides the `DataProductAdapter` interface - the contract that all data source adapters must implement. It includes a base implementation with common functionality like observability, event publishing, metrics collection, and lifecycle management.

## Features

- **Standard Interface**: `DataProductAdapter` contract for all adapters
- **Base Implementation**: Common functionality for all adapters
- **Lifecycle Management**: Initialize, start, stop, health checks
- **Event Publishing**: Automatic event publishing for entity changes
- **Observability**: Built-in tracing, metrics, and logging
- **Health Checks**: Automatic health check registration
- **Metrics Tracking**: Request counters, latency, cache statistics

## Installation

```bash
go get github.com/click2-run/dictamesh/pkg/adapter
```

## Creating an Adapter

### 1. Define Your Adapter

```go
package myadapter

import (
    "context"
    "github.com/click2-run/dictamesh/pkg/adapter"
    "github.com/click2-run/dictamesh/pkg/observability"
)

type MyAdapter struct {
    *adapter.BaseAdapter
    client *MyAPIClient
}

func NewMyAdapter(cfg *adapter.Config, obs *observability.Observability) (*MyAdapter, error) {
    base, err := adapter.NewBaseAdapter(cfg, obs)
    if err != nil {
        return nil, err
    }

    client := NewMyAPIClient(cfg)

    return &MyAdapter{
        BaseAdapter: base,
        client:      client,
    }, nil
}
```

### 2. Implement Required Methods

```go
// GetEntity retrieves a single entity
func (a *MyAdapter) GetEntity(ctx context.Context, entityType, id string) (*adapter.Entity, error) {
    return a.WithSpan(ctx, "get-entity", func(ctx context.Context) error {
        // Check cache first
        if a.config.EnableCache {
            if cached := a.getFromCache(entityType, id); cached != nil {
                a.IncrementCacheHit()
                return cached, nil
            }
            a.IncrementCacheMiss()
        }

        // Fetch from source
        data, err := a.client.Get(entityType, id)
        if err != nil {
            return nil, err
        }

        entity := &adapter.Entity{
            ID:     id,
            Type:   entityType,
            Domain: a.Domain(),
            Source: a.SourceSystem(),
            Data:   data,
        }

        // Publish event
        if a.config.EnableEvents {
            event := events.NewEvent(
                events.EventTypeEntityRead,
                a.Name(),
                entityType+":"+id,
                data,
            )
            a.PublishEvent(ctx, event)
        }

        return entity, nil
    })
}

// ListEntities lists entities with pagination
func (a *MyAdapter) ListEntities(ctx context.Context, entityType string, opts adapter.ListOptions) ([]*adapter.Entity, error) {
    // Implementation...
}

// CreateEntity creates a new entity
func (a *MyAdapter) CreateEntity(ctx context.Context, entityType string, data map[string]interface{}) (*adapter.Entity, error) {
    // Implementation...
}

// UpdateEntity updates an existing entity
func (a *MyAdapter) UpdateEntity(ctx context.Context, entityType, id string, data map[string]interface{}) (*adapter.Entity, error) {
    // Implementation...
}

// DeleteEntity deletes an entity
func (a *MyAdapter) DeleteEntity(ctx context.Context, entityType, id string) error {
    // Implementation...
}

// GetRelationships retrieves relationships for an entity
func (a *MyAdapter) GetRelationships(ctx context.Context, entityType, id string) ([]*adapter.Relationship, error) {
    // Implementation...
}

// GetSchema retrieves the schema for an entity type
func (a *MyAdapter) GetSchema(ctx context.Context, entityType string) (*adapter.Schema, error) {
    // Implementation...
}

// ListSchemas lists all available schemas
func (a *MyAdapter) ListSchemas(ctx context.Context) ([]*adapter.Schema, error) {
    // Implementation...
}

// InvalidateCache invalidates cache for an entity
func (a *MyAdapter) InvalidateCache(ctx context.Context, entityType, id string) error {
    // Implementation...
}
```

### 3. Use Your Adapter

```go
package main

import (
    "context"
    "github.com/click2-run/dictamesh/pkg/adapter"
    "github.com/click2-run/dictamesh/pkg/observability"
)

func main() {
    // Create observability
    obs, _ := observability.New(observability.DefaultConfig())
    obs.Start()
    defer obs.Shutdown(context.Background())

    // Create adapter configuration
    cfg := adapter.NewConfig(
        "my-adapter",
        "1.0.0",
        "my-api",
        "customers",
    ).WithCache(true).WithEvents(true)

    // Create adapter
    myAdapter, err := NewMyAdapter(cfg, obs)
    if err != nil {
        panic(err)
    }

    // Initialize and start
    ctx := context.Background()
    if err := myAdapter.Initialize(ctx); err != nil {
        panic(err)
    }
    if err := myAdapter.Start(ctx); err != nil {
        panic(err)
    }

    // Use adapter
    entity, err := myAdapter.GetEntity(ctx, "customer", "123")
    if err != nil {
        panic(err)
    }

    // Stop adapter
    defer myAdapter.Stop(ctx)
}
```

## DataProductAdapter Interface

```go
type DataProductAdapter interface {
    // Lifecycle
    Initialize(ctx context.Context) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Health(ctx context.Context) error

    // Metadata
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

    // Relationships
    GetRelationships(ctx context.Context, entityType, id string) ([]*Relationship, error)

    // Schemas
    GetSchema(ctx context.Context, entityType string) (*Schema, error)
    ListSchemas(ctx context.Context) ([]*Schema, error)

    // Cache
    InvalidateCache(ctx context.Context, entityType, id string) error
}
```

## Base Adapter Features

The `BaseAdapter` provides:

- **Lifecycle Management**: Automatic initialization and shutdown
- **Event Publishing**: Built-in event producer for entity changes
- **Observability**: Integrated tracing, metrics, and logging
- **Health Checks**: Automatic health check registration
- **Metrics Tracking**: Request counters, latency, errors
- **Status Management**: Track adapter status (ready, degraded, etc.)

### Helper Methods

```go
// WithSpan wraps operations with tracing
func (a *MyAdapter) MyOperation(ctx context.Context) error {
    return a.WithSpan(ctx, "my-operation", func(ctx context.Context) error {
        // Your logic here
        return nil
    })
}

// Increment metrics
a.IncrementRequests(true) // success
a.IncrementCacheHit()
a.RecordLatency(duration)

// Publish events
event := events.NewEvent(...)
a.PublishEvent(ctx, event)

// Get current status
status := a.GetStatus() // StatusReady, StatusDegraded, etc.

// Get metrics
metrics := a.GetMetrics()
```

## Configuration

```go
config := adapter.NewConfig("name", "1.0.0", "source", "domain")

// Enable/disable features
config.WithCache(true)
config.WithEvents(true)

// Custom settings
config.WithSetting("api_url", "https://api.example.com")
config.WithSetting("api_key", "secret")

// Retrieve settings
apiURL, _ := config.GetStringSetting("api_url")
```

## Best Practices

1. **Use BaseAdapter**: Always embed BaseAdapter for common functionality
2. **Wrap with WithSpan**: Use WithSpan for automatic tracing and metrics
3. **Publish Events**: Publish events for entity changes
4. **Implement Health Checks**: Ensure Health() checks dependencies
5. **Handle Errors**: Record errors with RecordError()
6. **Cache Appropriately**: Use caching for read-heavy workloads
7. **Validate Configuration**: Always validate config in NewAdapter()

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
