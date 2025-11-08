# DictaMesh Adapter Framework

The adapter framework provides the foundation for integrating third-party applications, APIs, databases, and other external systems into the DictaMesh data mesh platform.

## Overview

Adapters enable DictaMesh to integrate with external systems by providing:
- Standardized interfaces for data access
- Type-safe operations
- Error handling and retry logic
- Health monitoring
- Resource management

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    DictaMesh Core                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Catalog    │  │  Event Bus   │  │   Gateway    │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
└─────────┼──────────────────┼──────────────────┼─────────────┘
          │                  │                  │
          └──────────────────┴──────────────────┘
                             │
          ┌──────────────────┴──────────────────┐
          │                                     │
┌─────────▼────────────┐            ┌──────────▼──────────┐
│   Adapter Layer      │            │  Adapter Layer      │
│  ┌────────────────┐  │            │  ┌────────────────┐ │
│  │  Chatwoot      │  │            │  │  Salesforce    │ │
│  │  Adapter       │  │            │  │  Adapter       │ │
│  └───────┬────────┘  │            │  └───────┬────────┘ │
└──────────┼───────────┘            └──────────┼──────────┘
           │                                   │
┌──────────▼───────────┐            ┌──────────▼──────────┐
│   Connector Layer    │            │  Connector Layer    │
│  ┌────────────────┐  │            │  ┌────────────────┐ │
│  │  HTTP          │  │            │  │  REST API      │ │
│  │  Connector     │  │            │  │  Connector     │ │
│  └────────────────┘  │            │  └────────────────┘ │
└──────────────────────┘            └─────────────────────┘
```

## Package Structure

```
pkg/adapter/
├── adapter.go                # Core adapter interfaces
├── config.go                 # Configuration interfaces and utilities
├── errors.go                 # Error types and handling
├── http_client.go            # HTTP client with retry and rate limiting
├── README.md                 # This file
├── CONNECTOR-PATTERN.md      # Connector implementation guide
│
└── chatwoot/                 # Chatwoot adapter (reference implementation)
    ├── adapter.go
    ├── config.go
    ├── types.go
    ├── platform_client.go
    ├── application_client.go
    ├── application_client_extended.go
    ├── public_client.go
    └── README.md             # Complete adapter pattern guide
```

## Core Interfaces

### Adapter Interface

All adapters must implement the `Adapter` interface:

```go
type Adapter interface {
    Name() string
    Version() string
    Initialize(ctx context.Context, config Config) error
    Health(ctx context.Context) (*HealthStatus, error)
    Shutdown(ctx context.Context) error
    GetCapabilities() []Capability
}
```

### Resource Adapter Interface

For adapters that manage resources (CRUD operations):

```go
type ResourceAdapter interface {
    Adapter
    ListResources(ctx context.Context, opts *ListOptions) (*ResourceList, error)
    GetResource(ctx context.Context, resourceType, resourceID string) (*Resource, error)
    CreateResource(ctx context.Context, resource *Resource) (*Resource, error)
    UpdateResource(ctx context.Context, resource *Resource) (*Resource, error)
    DeleteResource(ctx context.Context, resourceType, resourceID string) error
}
```

### Streaming Adapter Interface

For adapters that support real-time data streaming:

```go
type StreamingAdapter interface {
    Adapter
    Subscribe(ctx context.Context, opts *SubscriptionOptions) (<-chan *Event, error)
    Unsubscribe(ctx context.Context, subscriptionID string) error
}
```

### Webhook Adapter Interface

For adapters that support webhooks:

```go
type WebhookAdapter interface {
    Adapter
    RegisterWebhook(ctx context.Context, webhook *WebhookConfig) (*Webhook, error)
    UnregisterWebhook(ctx context.Context, webhookID string) error
    ListWebhooks(ctx context.Context) ([]*Webhook, error)
    HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*Event, error)
}
```

## Available Adapters

### Chatwoot Adapter (Reference Implementation)

Full-featured adapter for Chatwoot customer engagement platform.

**Status**: ✅ Complete
**Documentation**: [pkg/adapter/chatwoot/README.md](./chatwoot/README.md)
**API Coverage**: Platform API, Application API, Public API

**Key Features:**
- Complete API coverage for all Chatwoot endpoints
- Support for all three API types (Platform, Application, Public)
- Comprehensive type definitions
- Retry logic and rate limiting
- Health checking
- Thread-safe operations

## Creating a New Adapter

### Quick Start

1. **Read the Documentation**
   - [Chatwoot Adapter README](./chatwoot/README.md) - Comprehensive adapter pattern guide
   - [Connector Pattern Guide](./CONNECTOR-PATTERN.md) - Connector implementation patterns

2. **Create Directory Structure**
   ```bash
   mkdir -p pkg/adapter/yourservice
   cd pkg/adapter/yourservice
   ```

3. **Create Required Files**
   ```bash
   touch adapter.go      # Main adapter implementation
   touch config.go       # Configuration
   touch types.go        # Domain types
   touch client.go       # API client
   touch README.md       # Documentation
   touch adapter_test.go # Tests
   ```

4. **Follow the Pattern**
   Use the Chatwoot adapter as your template:
   - Implement all required interfaces
   - Use `adapter.HTTPClient` for HTTP operations
   - Use `adapter.AdapterError` for errors
   - Follow naming conventions
   - Add comprehensive tests
   - Document everything

### Implementation Checklist

#### Core Files

- [ ] **adapter.go**
  - [ ] Implement `Adapter` interface
  - [ ] Add thread-safe state management (sync.RWMutex)
  - [ ] Implement health checking
  - [ ] Add proper cleanup in Shutdown()

- [ ] **config.go**
  - [ ] Implement `Config` interface
  - [ ] Add validation logic
  - [ ] Provide sensible defaults
  - [ ] Document all configuration options

- [ ] **types.go**
  - [ ] Define all domain types
  - [ ] Add JSON tags
  - [ ] Include request/response structures
  - [ ] Document type mappings

- [ ] **client.go**
  - [ ] Use `adapter.HTTPClient`
  - [ ] Implement CRUD operations
  - [ ] Add error handling
  - [ ] Include health check method

#### Testing

- [ ] Unit tests for all public methods
- [ ] Integration tests (if applicable)
- [ ] Mock server tests
- [ ] Error scenario tests
- [ ] Concurrent access tests

#### Documentation

- [ ] Comprehensive README
- [ ] Usage examples
- [ ] Configuration reference
- [ ] API coverage table
- [ ] Troubleshooting guide

## Utilities

### HTTP Client

The framework provides `HTTPClient` with built-in:
- Automatic retries with exponential backoff
- Rate limiting
- Request/response logging
- Timeout handling
- Connection pooling

```go
import "github.com/click2-run/dictamesh/pkg/adapter"

client := adapter.NewHTTPClient(&adapter.HTTPClientConfig{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
    RetryConfig: adapter.DefaultRetryConfig(),
    Headers: map[string]string{
        "Authorization": "Bearer token",
    },
})

resp, err := client.Get(ctx, "/api/v1/resource", nil)
```

### Error Handling

Use `AdapterError` for consistent error handling:

```go
import "github.com/click2-run/dictamesh/pkg/adapter"

// Create an error
err := adapter.NewAdapterError(
    adapter.ErrorCodeNotFound,
    "resource not found",
    originalErr,
)

// Add details
err.WithDetail("resource_id", id).WithRetryable(false)

// Check error type
if adapter.IsNotFoundError(err) {
    // Handle not found
}

if adapter.IsRetryableError(err) {
    // Retry operation
}
```

### Configuration

Use `MapConfig` for simple configuration:

```go
import "github.com/click2-run/dictamesh/pkg/adapter"

config := adapter.NewMapConfig(map[string]interface{}{
    "base_url": "https://api.example.com",
    "api_key":  "secret",
    "timeout":  "30s",
})

baseURL, _ := config.GetString("base_url")
timeout, _ := config.GetDuration("timeout")
```

## Best Practices

### 1. Follow Established Patterns
✅ Use Chatwoot adapter as your reference
✅ Follow the same file structure
✅ Use the same naming conventions
✅ Implement the same interfaces

### 2. Type Safety
✅ Define types for all entities
✅ Use strong typing for IDs
✅ Use time.Time for timestamps
✅ Add JSON tags to all fields

### 3. Error Handling
✅ Use AdapterError for all errors
✅ Include context in error messages
✅ Set appropriate error codes
✅ Mark errors as retryable when appropriate

### 4. Thread Safety
✅ Protect shared state with mutexes
✅ Use RWMutex for read-heavy workloads
✅ Document thread-safety guarantees
✅ Test concurrent access

### 5. Resource Management
✅ Close HTTP responses
✅ Cleanup resources in Shutdown()
✅ Use defer for cleanup
✅ Handle context cancellation

### 6. Testing
✅ Test all public methods
✅ Test error scenarios
✅ Test concurrent access
✅ Use mock servers for integration tests

### 7. Documentation
✅ Document all public types and methods
✅ Include usage examples
✅ Document configuration options
✅ Provide troubleshooting guide

## Contributing

When contributing a new adapter:

1. Follow the pattern established by the Chatwoot adapter
2. Read both pattern guides (Adapter and Connector)
3. Add comprehensive tests (aim for >80% coverage)
4. Document everything thoroughly
5. Include usage examples
6. Update this README with your adapter

## Testing

```bash
# Test all adapters
go test ./pkg/adapter/...

# Test specific adapter
go test ./pkg/adapter/chatwoot/...

# Run with coverage
go test -cover ./pkg/adapter/...

# Run with verbose output
go test -v ./pkg/adapter/...
```

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

## Resources

- [Chatwoot Adapter README](./chatwoot/README.md) - Complete adapter pattern guide
- [Connector Pattern Guide](./CONNECTOR-PATTERN.md) - Connector implementation guide
- [DictaMesh Documentation](../../docs/) - Framework documentation
- [Project Scope](../../PROJECT-SCOPE.md) - Project overview

---

**Start Here**: To create a new adapter, read the [Chatwoot Adapter README](./chatwoot/README.md) which serves as the definitive pattern guide for all adapter implementations.
