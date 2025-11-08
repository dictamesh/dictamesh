# Chatwoot Adapter for DictaMesh

This is the reference implementation of a third-party adapter for the DictaMesh framework, providing comprehensive integration with the [Chatwoot](https://www.chatwoot.com) customer engagement platform.

**⚠️ IMPORTANT:** This adapter serves as the **definitive pattern and template** for all future third-party adapter implementations in the DictaMesh ecosystem. Please read this documentation thoroughly before implementing new adapters.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Usage Examples](#usage-examples)
6. [API Coverage](#api-coverage)
7. [Adapter Pattern Guide](#adapter-pattern-guide)
8. [Best Practices](#best-practices)
9. [Testing](#testing)
10. [Contributing](#contributing)

## Overview

The Chatwoot adapter integrates all three Chatwoot API types:

- **Platform API**: Multi-tenant account management (requires platform API key)
- **Application API**: Account-specific operations (requires user API key)
- **Public API**: Client-side integrations (requires inbox identifier)

### Key Features

✅ Complete API coverage for all Chatwoot endpoints
✅ Type-safe Go implementations with comprehensive type definitions
✅ Automatic retry logic with exponential backoff
✅ Rate limiting support
✅ Health checking and monitoring
✅ Error handling with detailed error types
✅ Context-aware operations
✅ Configurable timeouts and connection pools

## Architecture

### Package Structure

```
pkg/adapter/chatwoot/
├── adapter.go                      # Main adapter implementation
├── config.go                       # Configuration management
├── types.go                        # Comprehensive type definitions
├── platform_client.go              # Platform API client
├── application_client.go           # Application API client
├── application_client_extended.go  # Extended Application API features
├── public_client.go                # Public API client
└── README.md                       # This file
```

### Component Responsibilities

**adapter.go**
- Implements the core `adapter.Adapter` interface
- Manages lifecycle (Initialize, Health, Shutdown)
- Provides access to API clients
- Handles adapter state and thread safety

**config.go**
- Defines adapter-specific configuration
- Implements `adapter.Config` interface
- Provides configuration validation
- Supports multiple API modes

**types.go**
- Defines all domain-specific types
- Maps to Chatwoot API resources
- Includes request/response structures
- Documents field mappings

**Client Files**
- Implement API operations for specific API types
- Handle HTTP communication
- Provide type-safe methods
- Include error handling

## Installation

### As a Go Module

```bash
go get github.com/click2-run/dictamesh/pkg/adapter/chatwoot
```

### Direct Integration

```go
import (
    "github.com/click2-run/dictamesh/pkg/adapter"
    "github.com/click2-run/dictamesh/pkg/adapter/chatwoot"
)
```

## Configuration

### Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `base_url` | string | Yes | Chatwoot instance URL |
| `platform_api_key` | string | Conditional | Platform API key (for Platform API) |
| `user_api_key` | string | Conditional | User API key (for Application API) |
| `account_id` | int64 | Conditional | Account ID (for Application API) |
| `inbox_identifier` | string | Conditional | Inbox identifier (for Public API) |
| `timeout` | duration | No | Request timeout (default: 30s) |
| `max_retries` | int | No | Maximum retry attempts (default: 3) |
| `retry_backoff` | duration | No | Initial backoff duration (default: 1s) |
| `rate_limit_per_second` | int | No | Rate limit (default: 10) |
| `enable_request_logging` | bool | No | Enable HTTP logging (default: false) |
| `webhook_secret` | string | No | Webhook signature secret |

### Example Configurations

#### Platform API Only

```go
config := &chatwoot.Config{
    BaseURL:            "https://app.chatwoot.com",
    PlatformAPIKey:     "your-platform-api-key",
    EnablePlatformAPI:  true,
}
```

#### Application API Only

```go
config := &chatwoot.Config{
    BaseURL:               "https://app.chatwoot.com",
    UserAPIKey:            "your-user-api-key",
    AccountID:             12345,
    EnableApplicationAPI:  true,
}
```

#### All APIs

```go
config := &chatwoot.Config{
    BaseURL:               "https://app.chatwoot.com",
    PlatformAPIKey:        "your-platform-api-key",
    UserAPIKey:            "your-user-api-key",
    AccountID:             12345,
    InboxIdentifier:       "your-inbox-id",
    EnablePlatformAPI:     true,
    EnableApplicationAPI:  true,
    EnablePublicAPI:       true,
}
```

## Usage Examples

### Basic Setup

```go
package main

import (
    "context"
    "log"

    "github.com/click2-run/dictamesh/pkg/adapter/chatwoot"
)

func main() {
    // Create adapter
    adapter := chatwoot.NewAdapter()

    // Configure
    config := &chatwoot.Config{
        BaseURL:              "https://app.chatwoot.com",
        UserAPIKey:           "your-api-key",
        AccountID:            12345,
        EnableApplicationAPI: true,
    }

    // Initialize
    ctx := context.Background()
    if err := adapter.Initialize(ctx, config); err != nil {
        log.Fatal(err)
    }
    defer adapter.Shutdown(ctx)

    // Check health
    health, err := adapter.Health(ctx)
    if err != nil {
        log.Printf("Health check failed: %v", err)
    } else {
        log.Printf("Adapter status: %s", health.Status)
    }
}
```

### Working with Conversations

```go
// Get Application API client
client, err := adapter.GetApplicationClient()
if err != nil {
    log.Fatal(err)
}

// List conversations
opts := &chatwoot.ConversationListOptions{
    Status: "open",
    Page:   1,
}
conversations, err := client.ListConversations(ctx, opts)
if err != nil {
    log.Fatal(err)
}

// Get specific conversation
conversation, err := client.GetConversation(ctx, conversationID)
if err != nil {
    log.Fatal(err)
}

// Send a message
message := &chatwoot.Message{
    Content:     "Hello! How can I help you?",
    MessageType: 1, // outgoing
    Private:     false,
}
sent, err := client.CreateMessage(ctx, conversation.ID, message)
if err != nil {
    log.Fatal(err)
}
```

### Managing Contacts

```go
// Create a contact
contact := &chatwoot.Contact{
    Name:        "John Doe",
    Email:       "john@example.com",
    PhoneNumber: "+1234567890",
    CustomAttributes: map[string]interface{}{
        "plan": "premium",
    },
}

created, err := client.CreateContact(ctx, contact)
if err != nil {
    log.Fatal(err)
}

// Search contacts
results, err := client.SearchContacts(ctx, "john@example.com", 1)
if err != nil {
    log.Fatal(err)
}

// Update contact
contact.Name = "John Smith"
updated, err := client.UpdateContact(ctx, contact.ID, contact)
if err != nil {
    log.Fatal(err)
}
```

### Using Platform API

```go
// Get Platform API client
platformClient, err := adapter.GetPlatformClient()
if err != nil {
    log.Fatal(err)
}

// Create an account
account := &chatwoot.Account{
    Name:   "Acme Corporation",
    Locale: "en",
}
created, err := platformClient.CreateAccount(ctx, account)
if err != nil {
    log.Fatal(err)
}

// Create a user
user := &chatwoot.User{
    Name:  "Jane Admin",
    Email: "jane@acme.com",
    Role:  "administrator",
}
createdUser, err := platformClient.CreateUser(ctx, user)
if err != nil {
    log.Fatal(err)
}

// Add user to account
userReq := &chatwoot.AccountUserRequest{
    UserID: createdUser.ID,
    Role:   "administrator",
}
_, err = platformClient.AddAccountUser(ctx, created.ID, userReq)
if err != nil {
    log.Fatal(err)
}
```

### Using Public API

```go
// Get Public API client
publicClient, err := adapter.GetPublicClient()
if err != nil {
    log.Fatal(err)
}

// Create a contact
contact := &chatwoot.Contact{
    Name:       "Customer Name",
    Email:      "customer@example.com",
    Identifier: "user-12345",
}
created, err := publicClient.CreateContact(ctx, contact)
if err != nil {
    log.Fatal(err)
}

// Create a conversation
convReq := &chatwoot.ConversationRequest{
    Message: "I need help with my order",
}
conversation, err := publicClient.CreateConversation(ctx, contact.Identifier, convReq)
if err != nil {
    log.Fatal(err)
}

// Send a message
msgReq := &chatwoot.MessageRequest{
    Content:     "My order number is #12345",
    MessageType: "incoming",
}
message, err := publicClient.CreateMessage(ctx, contact.Identifier, conversation.ID, msgReq)
if err != nil {
    log.Fatal(err)
}
```

## API Coverage

### Platform API

| Resource | Operations | Status |
|----------|------------|--------|
| Accounts | Create, Get, Update, Delete | ✅ |
| Account Users | List, Add, Remove | ✅ |
| Agent Bots | List, Create, Get, Update, Delete | ✅ |
| Users | Create, Get, Update, Delete, SSO Login | ✅ |

### Application API

| Resource | Operations | Status |
|----------|------------|--------|
| Accounts | Get, Update | ✅ |
| Agents | List, Add, Update, Remove | ✅ |
| Agent Bots | List, Create, Get, Update, Delete | ✅ |
| Audit Logs | List | ✅ |
| Automation Rules | List, Create, Get, Update, Delete | ✅ |
| Canned Responses | List, Create, Update, Delete | ✅ |
| Contacts | List, Create, Get, Update, Delete, Search, Filter | ✅ |
| Conversations | List, Get, Update, Assign, Toggle Status | ✅ |
| Custom Attributes | List, Create, Get, Update, Delete | ✅ |
| Inboxes | List, Create, Get, Update, Delete | ✅ |
| Integrations | List Apps, List, Enable | ✅ |
| Labels | List, Create, Delete | ✅ |
| Messages | List, Create, Update, Delete | ✅ |
| Reports | Account Reports, Conversation Metrics | ✅ |
| Teams | List, Create, Update, Delete | ✅ |
| Webhooks | List, Create, Update, Delete | ✅ |

### Public API

| Resource | Operations | Status |
|----------|------------|--------|
| Inbox | Get | ✅ |
| Contacts | Create, Get, Update | ✅ |
| Conversations | List, Create, Get, Resolve, Toggle Status, Update Last Seen | ✅ |
| Messages | List, Create, Update | ✅ |

---

## Adapter Pattern Guide

**This section defines the canonical patterns for implementing DictaMesh adapters.**

### Core Principles

1. **Interface Compliance**: All adapters MUST implement the `adapter.Adapter` interface
2. **Type Safety**: Use strongly-typed structures for all domain entities
3. **Error Handling**: Use `adapter.AdapterError` for consistent error reporting
4. **Context Awareness**: All operations must accept `context.Context` as the first parameter
5. **Thread Safety**: Adapter state must be protected with appropriate synchronization
6. **Resource Management**: Implement proper cleanup in `Shutdown()`

### File Structure Pattern

Every adapter implementation must follow this structure:

```
pkg/adapter/{provider}/
├── adapter.go           # Main adapter (implements adapter.Adapter)
├── config.go            # Configuration (implements adapter.Config)
├── types.go             # Domain types and structures
├── {api}_client.go      # API client(s) for external system
├── README.md            # Documentation
└── *_test.go            # Tests
```

### Implementation Checklist

#### 1. adapter.go

```go
type Adapter struct {
    config      *Config
    client      *Client
    initialized bool
    mu          sync.RWMutex
}

// Must implement:
- Name() string
- Version() string
- Initialize(ctx context.Context, config adapter.Config) error
- Health(ctx context.Context) (*adapter.HealthStatus, error)
- Shutdown(ctx context.Context) error
- GetCapabilities() []adapter.Capability
```

**Key Points:**
- Use `sync.RWMutex` for state protection
- Store initialization state
- Validate configuration before accepting it
- Return detailed health information
- Clean up all resources in Shutdown()

#### 2. config.go

```go
type Config struct {
    // Connection settings
    BaseURL string
    APIKey  string

    // Timeouts and limits
    Timeout     time.Duration
    MaxRetries  int

    // Feature flags
    EnableFeatureX bool
}

// Must implement:
- GetString(key string) (string, error)
- GetInt(key string) (int, error)
- GetBool(key string) (bool, error)
- GetDuration(key string) (time.Duration, error)
- Validate() error
```

**Key Points:**
- Provide sensible defaults
- Validate all required fields
- Support multiple configuration sources
- Document all configuration options
- Include feature flags for optional capabilities

#### 3. types.go

```go
// Define domain types that map to external API resources

type Resource struct {
    ID         int64                  `json:"id"`
    Name       string                 `json:"name"`
    Attributes map[string]interface{} `json:"attributes"`
    CreatedAt  time.Time              `json:"created_at"`
    UpdatedAt  time.Time              `json:"updated_at"`
}

// Include response wrappers
type ListResponse struct {
    Payload interface{}     `json:"payload"`
    Meta    *PaginationMeta `json:"meta,omitempty"`
}

// Include request structures
type CreateResourceRequest struct {
    Name       string                 `json:"name"`
    Attributes map[string]interface{} `json:"attributes,omitempty"`
}
```

**Key Points:**
- Map all external API types
- Use consistent JSON tags
- Include omitempty for optional fields
- Document non-obvious mappings
- Use time.Time for timestamps
- Use appropriate numeric types (int64 for IDs)

#### 4. Client Implementation

```go
type Client struct {
    httpClient *adapter.HTTPClient
    apiKey     string
    baseURL    string
}

func NewClient(config *Config) *Client {
    httpClient := adapter.NewHTTPClient(&adapter.HTTPClientConfig{
        BaseURL: config.BaseURL,
        Timeout: config.Timeout,
        RetryConfig: adapter.DefaultRetryConfig(),
        Headers: map[string]string{
            "Authorization": "Bearer " + config.APIKey,
        },
    })

    return &Client{
        httpClient: httpClient,
        apiKey:     config.APIKey,
        baseURL:    config.BaseURL,
    }
}

// Implement CRUD operations
func (c *Client) ListResources(ctx context.Context, opts *ListOptions) ([]Resource, error)
func (c *Client) GetResource(ctx context.Context, id int64) (*Resource, error)
func (c *Client) CreateResource(ctx context.Context, resource *Resource) (*Resource, error)
func (c *Client) UpdateResource(ctx context.Context, id int64, resource *Resource) (*Resource, error)
func (c *Client) DeleteResource(ctx context.Context, id int64) error
```

**Key Points:**
- Use `adapter.HTTPClient` for HTTP operations
- Configure retry logic and timeouts
- Handle pagination consistently
- Parse errors using `adapter.HTTPErrorToAdapterError`
- Include health check method
- Implement Close() for cleanup

### Error Handling Pattern

```go
// Use adapter.AdapterError for all errors
func (c *Client) GetResource(ctx context.Context, id int64) (*Resource, error) {
    path := fmt.Sprintf("/api/v1/resources/%d", id)

    resp, err := c.httpClient.Get(ctx, path, nil)
    if err != nil {
        return nil, adapter.NewAdapterError(
            adapter.ErrorCodeConnectionFailed,
            "failed to fetch resource",
            err,
        )
    }

    var result Resource
    if err := adapter.ParseJSONResponse(resp, &result); err != nil {
        return nil, err // ParseJSONResponse returns AdapterError
    }

    return &result, nil
}
```

### Testing Pattern

```go
func TestAdapter_Initialize(t *testing.T) {
    adapter := NewAdapter()

    config := &Config{
        BaseURL: "https://api.example.com",
        APIKey:  "test-key",
    }

    err := adapter.Initialize(context.Background(), config)
    if err != nil {
        t.Fatalf("Initialize failed: %v", err)
    }

    if !adapter.IsInitialized() {
        t.Error("Adapter should be initialized")
    }
}

func TestClient_GetResource(t *testing.T) {
    // Create mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(&Resource{
            ID:   1,
            Name: "Test Resource",
        })
    }))
    defer server.Close()

    config := &Config{
        BaseURL: server.URL,
        APIKey:  "test",
    }

    client := NewClient(config)

    resource, err := client.GetResource(context.Background(), 1)
    if err != nil {
        t.Fatalf("GetResource failed: %v", err)
    }

    if resource.Name != "Test Resource" {
        t.Errorf("Expected name 'Test Resource', got '%s'", resource.Name)
    }
}
```

## Best Practices

### 1. Configuration Management

- ✅ Provide clear defaults for all optional parameters
- ✅ Validate configuration in `Validate()` method
- ✅ Support environment variables for sensitive data
- ✅ Document all configuration options with examples
- ❌ Don't hard-code timeouts or limits
- ❌ Don't store sensitive data in logs

### 2. Error Handling

- ✅ Use `adapter.AdapterError` for all errors
- ✅ Include context in error messages
- ✅ Set appropriate error codes
- ✅ Mark errors as retryable when appropriate
- ❌ Don't swallow errors
- ❌ Don't expose internal error details to users

### 3. HTTP Client Usage

- ✅ Use `adapter.HTTPClient` for all HTTP operations
- ✅ Configure retry logic with exponential backoff
- ✅ Set appropriate timeouts
- ✅ Use rate limiting when needed
- ✅ Log requests in debug mode only
- ❌ Don't create custom HTTP clients
- ❌ Don't ignore timeout configurations

### 4. Type Safety

- ✅ Define types for all domain entities
- ✅ Use strong typing for IDs (int64, string, etc.)
- ✅ Use time.Time for timestamps
- ✅ Include JSON tags on all fields
- ✅ Use pointers for optional fields
- ❌ Don't use `interface{}` unless necessary
- ❌ Don't mix pointer and value receivers

### 5. Thread Safety

- ✅ Protect shared state with mutexes
- ✅ Use RWMutex for read-heavy workloads
- ✅ Document thread-safety guarantees
- ❌ Don't assume single-threaded usage
- ❌ Don't share mutable state without protection

### 6. Resource Management

- ✅ Close HTTP responses
- ✅ Cleanup resources in Shutdown()
- ✅ Use defer for cleanup
- ✅ Handle context cancellation
- ❌ Don't leak goroutines
- ❌ Don't leave connections open

### 7. Documentation

- ✅ Document all public types and methods
- ✅ Include usage examples
- ✅ Document error conditions
- ✅ Provide configuration examples
- ✅ Explain non-obvious behavior
- ❌ Don't assume prior knowledge
- ❌ Don't leave TODOs in production code

## Testing

### Running Tests

```bash
# Run all tests
go test ./pkg/adapter/chatwoot/...

# Run with coverage
go test -cover ./pkg/adapter/chatwoot/...

# Run with verbose output
go test -v ./pkg/adapter/chatwoot/...
```

### Integration Testing

Integration tests require a running Chatwoot instance:

```bash
# Set environment variables
export CHATWOOT_BASE_URL="https://app.chatwoot.com"
export CHATWOOT_USER_API_KEY="your-api-key"
export CHATWOOT_ACCOUNT_ID="12345"

# Run integration tests
go test -tags=integration ./pkg/adapter/chatwoot/...
```

## Contributing

When adding new features or fixing bugs:

1. Follow the adapter pattern defined in this README
2. Add comprehensive tests
3. Update documentation
4. Ensure backward compatibility
5. Add usage examples

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda

---

**This adapter is the reference implementation for the DictaMesh adapter pattern. When creating new adapters, use this as your template and follow all established patterns.**
