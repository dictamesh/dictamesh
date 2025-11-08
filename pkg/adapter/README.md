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
