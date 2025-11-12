# Connector Pattern Guide

## Overview

Connectors are **low-level drivers** that handle the technical details of connecting to specific types of data sources. They are protocol and technology-specific components that provide a standardized interface for accessing external systems.

**Key Distinction:**
- **Connectors** = Technical connectivity layer (HTTP, DB, File, etc.)
- **Adapters** = Business logic layer (Chatwoot, Salesforce, etc.)

## Connector vs Adapter

### Connector Responsibilities

✅ Protocol-level communication (HTTP, TCP, Database protocols)
✅ Authentication and session management
✅ Connection pooling and lifecycle management
✅ Low-level error handling and retries
✅ Data serialization/deserialization
✅ Query execution and result parsing

**Example Connectors:**
- `HTTPConnector` - For REST APIs, GraphQL, SOAP
- `PostgreSQLConnector` - For PostgreSQL databases
- `MongoDBConnector` - For MongoDB databases
- `FileConnector` - For file system access
- `KafkaConnector` - For Kafka message brokers

### Adapter Responsibilities

✅ Business logic and domain modeling
✅ Entity transformation and mapping
✅ Workflow orchestration
✅ Integration with DictaMesh framework
✅ Event publishing and metadata registration
✅ Relationship management

**Example Adapters:**
- `ChatwootAdapter` - Uses HTTPConnector to integrate Chatwoot
- `SalesforceAdapter` - Uses HTTPConnector for Salesforce API
- `PostgreSQLDataAdapter` - Uses PostgreSQLConnector for data integration

## When to Create a Connector

Create a new connector when:

1. **New Protocol/Technology**: You need to support a communication protocol or technology that doesn't have a connector yet
2. **Specialized Requirements**: The technology has specific requirements that generic connectors can't handle
3. **Performance Optimization**: You need optimized access patterns for a specific system type
4. **Reusability**: Multiple adapters will benefit from the same connectivity layer

## Connector Architecture

### Standard Connector Interface

```go
package connector

import (
    "context"
    "io"
)

// Connector is the base interface for all connectors
type Connector interface {
    // Name returns the connector name (e.g., "http", "postgresql")
    Name() string

    // Version returns the connector version
    Version() string

    // Connect establishes a connection to the external system
    Connect(ctx context.Context, config Config) error

    // Disconnect closes the connection
    Disconnect(ctx context.Context) error

    // Ping checks if the connection is alive
    Ping(ctx context.Context) error

    // IsConnected returns the connection status
    IsConnected() bool
}

// Config represents connector configuration
type Config interface {
    // Validate checks if the configuration is valid
    Validate() error

    // GetConnectionString returns the connection string
    GetConnectionString() string
}
```

### HTTP Connector Pattern

```go
package http

import (
    "context"
    "net/http"
    "time"

    "github.com/click2-run/dictamesh/pkg/connector"
)

// HTTPConnector implements HTTP/HTTPS communication
type HTTPConnector struct {
    client      *http.Client
    config      *Config
    connected   bool
    mu          sync.RWMutex
}

// Config represents HTTP connector configuration
type Config struct {
    // Connection settings
    BaseURL    string
    Timeout    time.Duration
    MaxRetries int

    // TLS settings
    TLSConfig  *tls.Config
    SkipVerify bool

    // Auth settings
    AuthType   AuthType
    Credentials Credentials

    // Connection pool
    MaxIdleConns        int
    MaxIdleConnsPerHost int
    MaxConnsPerHost     int
}

// AuthType represents authentication type
type AuthType string

const (
    AuthTypeNone   AuthType = "none"
    AuthTypeBasic  AuthType = "basic"
    AuthTypeBearer AuthType = "bearer"
    AuthTypeAPIKey AuthType = "apikey"
    AuthTypeOAuth2 AuthType = "oauth2"
)

// NewHTTPConnector creates a new HTTP connector
func NewHTTPConnector() *HTTPConnector {
    return &HTTPConnector{
        connected: false,
    }
}

// Connect establishes an HTTP client
func (c *HTTPConnector) Connect(ctx context.Context, cfg connector.Config) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    config, ok := cfg.(*Config)
    if !ok {
        return connector.ErrInvalidConfig
    }

    if err := config.Validate(); err != nil {
        return err
    }

    // Create transport
    transport := &http.Transport{
        MaxIdleConns:        config.MaxIdleConns,
        MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
        MaxConnsPerHost:     config.MaxConnsPerHost,
        TLSClientConfig:     config.TLSConfig,
        IdleConnTimeout:     90 * time.Second,
    }

    // Create client
    c.client = &http.Client{
        Transport: transport,
        Timeout:   config.Timeout,
    }

    c.config = config
    c.connected = true

    return nil
}

// Execute performs an HTTP request
func (c *HTTPConnector) Execute(ctx context.Context, req *Request) (*Response, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if !c.connected {
        return nil, connector.ErrNotConnected
    }

    // Build HTTP request
    httpReq, err := c.buildRequest(ctx, req)
    if err != nil {
        return nil, err
    }

    // Apply authentication
    if err := c.applyAuth(httpReq); err != nil {
        return nil, err
    }

    // Execute with retries
    resp, err := c.executeWithRetry(ctx, httpReq)
    if err != nil {
        return nil, err
    }

    return c.parseResponse(resp)
}

// Request represents an HTTP request
type Request struct {
    Method  string
    Path    string
    Headers map[string]string
    Body    []byte
    Query   map[string]string
}

// Response represents an HTTP response
type Response struct {
    StatusCode int
    Headers    map[string][]string
    Body       []byte
}
```

### Database Connector Pattern

```go
package database

import (
    "context"
    "database/sql"

    "github.com/click2-run/dictamesh/pkg/connector"
)

// DatabaseConnector implements database connectivity
type DatabaseConnector struct {
    db        *sql.DB
    config    *Config
    connected bool
    mu        sync.RWMutex
}

// Config represents database connector configuration
type Config struct {
    // Connection
    Host     string
    Port     int
    Database string
    Schema   string

    // Authentication
    Username string
    Password string

    // Connection pool
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration

    // SSL/TLS
    SSLMode     string
    SSLCert     string
    SSLKey      string
    SSLRootCert string

    // Driver-specific options
    Options map[string]string
}

// Connect establishes a database connection
func (c *DatabaseConnector) Connect(ctx context.Context, cfg connector.Config) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    config, ok := cfg.(*Config)
    if !ok {
        return connector.ErrInvalidConfig
    }

    if err := config.Validate(); err != nil {
        return err
    }

    // Build connection string
    connStr := c.buildConnectionString(config)

    // Open connection
    db, err := sql.Open(c.driverName(), connStr)
    if err != nil {
        return connector.NewConnectorError(
            connector.ErrorCodeConnectionFailed,
            "failed to open database connection",
            err,
        )
    }

    // Configure connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

    // Test connection
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return connector.NewConnectorError(
            connector.ErrorCodeConnectionFailed,
            "failed to ping database",
            err,
        )
    }

    c.db = db
    c.config = config
    c.connected = true

    return nil
}

// Query executes a query
func (c *DatabaseConnector) Query(ctx context.Context, query string, args ...interface{}) (*QueryResult, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if !c.connected {
        return nil, connector.ErrNotConnected
    }

    rows, err := c.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, connector.NewConnectorError(
            connector.ErrorCodeQueryFailed,
            "query execution failed",
            err,
        )
    }

    return c.parseRows(rows)
}

// Execute executes a statement (INSERT, UPDATE, DELETE)
func (c *DatabaseConnector) Execute(ctx context.Context, stmt string, args ...interface{}) (*ExecResult, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if !c.connected {
        return nil, connector.ErrNotConnected
    }

    result, err := c.db.ExecContext(ctx, stmt, args...)
    if err != nil {
        return nil, connector.NewConnectorError(
            connector.ErrorCodeExecFailed,
            "statement execution failed",
            err,
        )
    }

    return &ExecResult{
        LastInsertId: getLastInsertId(result),
        RowsAffected: getRowsAffected(result),
    }, nil
}

// Transaction creates a new transaction
func (c *DatabaseConnector) Transaction(ctx context.Context, opts *sql.TxOptions) (*Transaction, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if !c.connected {
        return nil, connector.ErrNotConnected
    }

    tx, err := c.db.BeginTx(ctx, opts)
    if err != nil {
        return nil, connector.NewConnectorError(
            connector.ErrorCodeTransactionFailed,
            "failed to begin transaction",
            err,
        )
    }

    return &Transaction{tx: tx}, nil
}
```

## File Structure for Connectors

```
pkg/connector/
├── connector.go              # Base interfaces
├── config.go                 # Common configuration
├── errors.go                 # Error types
├── http/                     # HTTP connector
│   ├── connector.go
│   ├── config.go
│   ├── auth.go              # Authentication handlers
│   ├── retry.go             # Retry logic
│   └── README.md
├── postgresql/               # PostgreSQL connector
│   ├── connector.go
│   ├── config.go
│   ├── query_builder.go
│   └── README.md
├── mongodb/                  # MongoDB connector
│   ├── connector.go
│   ├── config.go
│   └── README.md
└── README.md
```

## Connector Implementation Checklist

### 1. Core Interface
- [ ] Implement `Connector` interface
- [ ] Thread-safe connection management
- [ ] Proper resource cleanup in `Disconnect()`
- [ ] Health checking in `Ping()`

### 2. Configuration
- [ ] Define `Config` struct with all options
- [ ] Implement `Validate()` method
- [ ] Provide sensible defaults
- [ ] Support connection strings
- [ ] Document all configuration options

### 3. Connection Management
- [ ] Connection pooling (where applicable)
- [ ] Automatic reconnection on failure
- [ ] Connection timeout handling
- [ ] Idle connection cleanup
- [ ] Resource limits enforcement

### 4. Error Handling
- [ ] Use `connector.ConnectorError`
- [ ] Classify errors (transient vs permanent)
- [ ] Include retry hints
- [ ] Preserve error context
- [ ] Log errors appropriately

### 5. Authentication
- [ ] Support multiple auth methods
- [ ] Secure credential storage
- [ ] Token refresh handling
- [ ] Certificate management (for TLS)

### 6. Testing
- [ ] Unit tests for all methods
- [ ] Integration tests with real systems
- [ ] Mock implementations for testing
- [ ] Connection failure scenarios
- [ ] Timeout handling tests

### 7. Documentation
- [ ] README with usage examples
- [ ] Configuration reference
- [ ] Authentication guide
- [ ] Troubleshooting section
- [ ] Performance tuning tips

## Best Practices

### Connection Management

✅ **DO:**
- Use connection pooling
- Implement automatic reconnection
- Set appropriate timeouts
- Clean up resources properly
- Monitor connection health

❌ **DON'T:**
- Create connections per request
- Leave connections open indefinitely
- Ignore connection failures
- Share connections unsafely

### Configuration

✅ **DO:**
- Validate all configuration
- Provide defaults
- Support environment variables
- Document all options
- Version configuration schemas

❌ **DON'T:**
- Hard-code values
- Expose sensitive data
- Use complex config structures
- Ignore validation errors

### Error Handling

✅ **DO:**
- Use typed errors
- Include context
- Log errors
- Classify error types
- Provide recovery hints

❌ **DON'T:**
- Swallow errors
- Expose internal details
- Use generic error messages
- Panic on recoverable errors

### Thread Safety

✅ **DO:**
- Protect shared state
- Use appropriate locks
- Document thread-safety
- Test concurrent access

❌ **DON'T:**
- Assume single-threaded
- Share mutable state
- Use locks unnecessarily
- Create deadlocks

## Example: Creating a GraphQL Connector

```go
package graphql

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"

    "github.com/click2-run/dictamesh/pkg/connector"
    "github.com/click2-run/dictamesh/pkg/connector/http"
)

// GraphQLConnector wraps HTTP connector for GraphQL
type GraphQLConnector struct {
    httpConn *http.HTTPConnector
    config   *Config
}

// Config extends HTTP config with GraphQL specifics
type Config struct {
    http.Config
    Endpoint string
    Schema   string
}

// NewGraphQLConnector creates a GraphQL connector
func NewGraphQLConnector() *GraphQLConnector {
    return &GraphQLConnector{
        httpConn: http.NewHTTPConnector(),
    }
}

// Connect establishes connection
func (c *GraphQLConnector) Connect(ctx context.Context, cfg connector.Config) error {
    config, ok := cfg.(*Config)
    if !ok {
        return connector.ErrInvalidConfig
    }

    c.config = config

    // Initialize HTTP connector
    return c.httpConn.Connect(ctx, &config.Config)
}

// Query executes a GraphQL query
func (c *GraphQLConnector) Query(ctx context.Context, query string, variables map[string]interface{}) (*QueryResult, error) {
    payload := map[string]interface{}{
        "query": query,
    }
    if variables != nil {
        payload["variables"] = variables
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    req := &http.Request{
        Method: "POST",
        Path:   c.config.Endpoint,
        Body:   body,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }

    resp, err := c.httpConn.Execute(ctx, req)
    if err != nil {
        return nil, err
    }

    return c.parseResponse(resp)
}

// Mutation executes a GraphQL mutation
func (c *GraphQLConnector) Mutation(ctx context.Context, mutation string, variables map[string]interface{}) (*MutationResult, error) {
    // Similar to Query but for mutations
    return c.executeMutation(ctx, mutation, variables)
}
```

## Integration with Adapters

Adapters use connectors to access external systems:

```go
// In adapter implementation
type ChatwootAdapter struct {
    httpConn *http.HTTPConnector
    config   *Config
}

func (a *ChatwootAdapter) Initialize(ctx context.Context, cfg adapter.Config) error {
    // Create HTTP connector
    a.httpConn = http.NewHTTPConnector()

    // Configure connector
    connConfig := &http.Config{
        BaseURL:    a.config.BaseURL,
        Timeout:    a.config.Timeout,
        AuthType:   http.AuthTypeAPIKey,
        Credentials: http.Credentials{
            APIKey: a.config.APIKey,
        },
    }

    // Connect
    if err := a.httpConn.Connect(ctx, connConfig); err != nil {
        return err
    }

    return nil
}

// Use connector for operations
func (a *ChatwootAdapter) GetContact(ctx context.Context, id string) (*Contact, error) {
    req := &http.Request{
        Method: "GET",
        Path:   fmt.Sprintf("/api/v1/contacts/%s", id),
    }

    resp, err := a.httpConn.Execute(ctx, req)
    if err != nil {
        return nil, err
    }

    var contact Contact
    if err := json.Unmarshal(resp.Body, &contact); err != nil {
        return nil, err
    }

    return &contact, nil
}
```

## Summary

**Connectors** provide the plumbing for adapters to communicate with external systems. They handle:
- Protocol-level details
- Connection management
- Authentication
- Error handling
- Resource pooling

**Adapters** build on connectors to provide:
- Business logic
- Domain modeling
- Framework integration
- Event publishing
- Relationship management

When implementing a new integration:
1. Check if a suitable connector exists
2. If not, create a connector for the protocol/technology
3. Create an adapter that uses the connector
4. Follow the patterns established in this guide

---

**Next Steps**: See `/pkg/adapter/chatwoot/README.md` for the complete adapter implementation pattern.
