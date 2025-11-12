# DictaMesh Observability Package

Comprehensive observability infrastructure for the DictaMesh framework, providing distributed tracing, metrics collection, structured logging, and health checks.

## Features

### 🔍 Distributed Tracing
- **OpenTelemetry Integration**: Industry-standard tracing with OTLP export
- **Jaeger Support**: Legacy Jaeger exporter for existing setups
- **Context Propagation**: W3C Trace Context standard
- **Automatic Instrumentation**: Helper functions for common patterns
- **Configurable Sampling**: Control trace volume in production

### 📊 Metrics Collection
- **Prometheus Integration**: Native Prometheus metrics export
- **Pre-built Metrics**: HTTP, database, event bus, adapter, and GraphQL metrics
- **Runtime Metrics**: Go runtime statistics (goroutines, memory, GC)
- **Custom Metrics**: Easy-to-use metric builders
- **Auto-discovery**: Automatic service registration

### 📝 Structured Logging
- **Zap-based Logging**: High-performance structured logging
- **Multiple Formats**: JSON for production, console for development
- **Trace Correlation**: Automatic trace ID/span ID injection
- **Context-aware**: Logger enrichment from context
- **Configurable Sampling**: Reduce log volume for high-throughput services

### 🏥 Health Checks
- **Kubernetes-ready**: Liveness, readiness, and startup probes
- **Pluggable Checks**: Easy to add custom health checks
- **Dependency Tracking**: Check external dependencies
- **HTTP Endpoints**: Standard REST endpoints
- **Timeout Protection**: Prevent hanging health checks

## Installation

```bash
go get github.com/click2-run/dictamesh/pkg/observability
```

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "log"

    "github.com/click2-run/dictamesh/pkg/observability"
)

func main() {
    // Create configuration
    config := observability.DefaultConfig()
    config.ServiceName = "my-service"
    config.ServiceVersion = "1.0.0"
    config.Environment = "production"

    // Initialize observability
    obs, err := observability.New(config)
    if err != nil {
        log.Fatal(err)
    }

    // Start all components (metrics server, health checks, tracing)
    if err := obs.Start(); err != nil {
        log.Fatal(err)
    }

    // Use in your application
    ctx := context.Background()
    obs.Logger().Info("service started")

    // Graceful shutdown
    defer obs.Shutdown(ctx)
}
```

### Production Configuration

```go
config := observability.ProductionConfig()
config.ServiceName = "dictamesh-adapter"
config.ServiceVersion = "1.0.0"

// Configure tracing for production (10% sampling)
config.Tracing.Endpoint = "otel-collector:4318"
config.Tracing.SamplingRate = 0.1
config.Tracing.Insecure = false

// Configure metrics
config.Metrics.Port = 9090
config.Metrics.Namespace = "dictamesh"

obs, err := observability.New(config)
```

## Usage Examples

### Distributed Tracing

#### Starting Spans

```go
ctx := context.Background()
obs, _ := observability.New(observability.DefaultConfig())

// Start a span
ctx, span := obs.StartSpan(ctx, "process-entity")
defer span.End()

// Do work...
result, err := processEntity(ctx, entityID)

// Record error if any
if err != nil {
    observability.RecordError(ctx, err)
    return err
}

// Set attributes
observability.SetAttributes(ctx,
    observability.AttrEntityType.String("customer"),
    observability.AttrEntityID.String(entityID),
)
```

#### Using WithSpan Helper

```go
err := observability.WithSpan(ctx, obs.Tracer(), "database-query", func(ctx context.Context) error {
    return db.Query(ctx, "SELECT * FROM customers")
})
```

#### Adding Custom Attributes

```go
observability.SetAttributes(ctx,
    observability.AttrHTTPMethod.String("POST"),
    observability.AttrHTTPURL.String("/api/entities"),
    observability.AttrUserID.String("user-123"),
)
```

### Structured Logging

#### Basic Logging

```go
logger := obs.Logger()

logger.Info("processing started",
    "entity_id", entityID,
    "entity_type", "customer",
)

logger.Error("processing failed",
    "entity_id", entityID,
    "error", err.Error(),
)
```

#### Context-Aware Logging

```go
// Automatically includes trace_id and span_id
logger.InfoContext(ctx, "entity processed successfully",
    "entity_id", entityID,
    "duration_ms", duration.Milliseconds(),
)
```

#### Logger Enrichment

```go
// Add fields to logger
enrichedLogger := logger.WithFields(map[string]interface{}{
    "tenant_id": "tenant-123",
    "user_id": "user-456",
})

// Add single field
enrichedLogger = enrichedLogger.WithField("request_id", requestID)

// Add error
enrichedLogger = enrichedLogger.WithError(err)

enrichedLogger.Info("processing entity")
```

#### Named Loggers

```go
adapterLogger := logger.Named("adapter")
adapterLogger.Info("adapter initialized") // logger: "service.adapter"

dbLogger := adapterLogger.Named("database")
dbLogger.Info("connection established") // logger: "service.adapter.database"
```

### Metrics Collection

#### HTTP Metrics

```go
start := time.Now()
status, requestSize, responseSize := handleRequest(w, r)
duration := time.Since(start)

obs.Metrics().RecordHTTPRequest(
    r.Method,
    r.URL.Path,
    strconv.Itoa(status),
    duration,
    requestSize,
    responseSize,
)
```

#### Database Metrics

```go
start := time.Now()
err := db.Exec("INSERT INTO customers ...")
duration := time.Since(start)

status := "success"
if err != nil {
    status = "error"
}

obs.Metrics().RecordDBQuery("INSERT", "customers", status, duration)
```

#### Event Bus Metrics

```go
// Publishing events
start := time.Now()
err := publishEvent(topic, event)
duration := time.Since(start)

status := "success"
if err != nil {
    status = "error"
}

obs.Metrics().RecordEventPublish(topic, eventType, status, duration)

// Consuming events
obs.Metrics().RecordEventConsume(topic, consumerGroup, status, duration)
```

#### Custom Metrics

```go
// Access Prometheus metrics directly
metrics := obs.Metrics()

// Increment counter
metrics.AdapterRequestsTotal.WithLabelValues("customer-adapter", "get", "success").Inc()

// Record histogram
metrics.AdapterRequestDuration.WithLabelValues("customer-adapter", "get").Observe(duration.Seconds())

// Set gauge
metrics.DBConnectionsOpen.WithLabelValues("primary").Set(float64(openConnections))
```

### Health Checks

#### Registering Health Checks

```go
health := obs.Health()

// Liveness check (is the process alive?)
health.RegisterLivenessCheck("ping", observability.PingCheck())

// Readiness check (can we serve traffic?)
health.RegisterReadinessCheck("database", func(ctx context.Context) error {
    return db.Ping(ctx)
})

health.RegisterReadinessCheck("cache", func(ctx context.Context) error {
    return redis.Ping(ctx).Err()
})

// Startup check (has the service initialized?)
health.RegisterStartupCheck("migrations", func(ctx context.Context) error {
    return migrator.Verify(ctx)
})
```

#### Custom Health Checks with Timeout

```go
health.RegisterReadinessCheck("external-api",
    observability.TimeoutCheck(2*time.Second, func(ctx context.Context) error {
        resp, err := http.Get("https://api.example.com/health")
        if err != nil {
            return err
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            return fmt.Errorf("unhealthy: status %d", resp.StatusCode)
        }
        return nil
    }),
)
```

#### Health Check Endpoints

Once started, health checks are available at:
- **Liveness**: `http://localhost:8081/health/live`
- **Readiness**: `http://localhost:8081/health/ready`
- **Startup**: `http://localhost:8081/health/startup`

Response format:
```json
{
  "status": "healthy",
  "timestamp": "2025-11-08T10:30:00Z",
  "checks": {
    "database": {
      "status": "healthy",
      "duration": "5ms"
    },
    "cache": {
      "status": "healthy",
      "duration": "2ms"
    }
  }
}
```

### Context Propagation

#### Adding Context Values

```go
ctx = observability.WithRequestID(ctx, "req-123")
ctx = observability.WithUserID(ctx, "user-456")
ctx = observability.WithTenantID(ctx, "tenant-789")

// Retrieve values
requestID := observability.RequestIDFromContext(ctx)
userID := observability.UserIDFromContext(ctx)
tenantID := observability.TenantIDFromContext(ctx)
```

#### Logger in Context

```go
// Add logger to context
ctx = observability.WithLogger(ctx, logger)

// Retrieve logger from context
logger := observability.LoggerFromContext(ctx)
if logger != nil {
    logger.Info("processing request")
}
```

## Configuration

### Default Configuration

```go
config := observability.DefaultConfig()
// Returns configuration suitable for development:
// - Tracing: 100% sampling, local Jaeger
// - Metrics: enabled on port 9090
// - Logging: info level, JSON format
// - Health: enabled on port 8081
```

### Production Configuration

```go
config := observability.ProductionConfig()
// Returns configuration optimized for production:
// - Tracing: 10% sampling, TLS enabled
// - Metrics: enabled with runtime metrics
// - Logging: info level, no stack traces
// - Health: enabled with checks
```

### Custom Configuration

```go
config := &observability.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Environment:    "staging",

    Tracing: observability.TracingConfig{
        Enabled:      true,
        Endpoint:     "otel-collector:4318",
        SamplingRate: 0.5, // 50% sampling
        Insecure:     false,
    },

    Metrics: observability.MetricsConfig{
        Enabled:              true,
        Port:                 9090,
        Namespace:            "dictamesh",
        EnableRuntimeMetrics: true,
        DefaultHistogramBuckets: []float64{0.001, 0.01, 0.1, 1, 10},
    },

    Logging: observability.LoggingConfig{
        Level:            "debug",
        Format:           "console",
        OutputPaths:      []string{"stdout", "/var/log/app.log"},
        EnableStackTrace: true,
        EnableCaller:     true,
    },

    Health: observability.HealthConfig{
        Enabled:       true,
        Port:          8081,
        CheckInterval: 10 * time.Second,
        Timeout:       5 * time.Second,
    },
}
```

## Integration with DictaMesh

### Adapter Integration

```go
type MyAdapter struct {
    obs *observability.Observability
}

func (a *MyAdapter) GetEntity(ctx context.Context, id string) (*Entity, error) {
    // Start span
    ctx, span := a.obs.StartSpan(ctx, "adapter.get-entity")
    defer span.End()

    // Log with context
    a.obs.LoggerWithContext(ctx).Info("fetching entity", "id", id)

    // Record metrics
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        a.obs.Metrics().AdapterRequestDuration.
            WithLabelValues("my-adapter", "get").
            Observe(duration.Seconds())
    }()

    // Do work...
    return entity, nil
}
```

### Service Integration

```go
func main() {
    // Initialize observability
    obs, _ := observability.New(observability.DefaultConfig())
    obs.Start()
    defer obs.Shutdown(context.Background())

    // Register health checks
    obs.Health().RegisterReadinessCheck("database", dbHealthCheck)

    // Use throughout service
    logger := obs.Logger()
    tracer := obs.Tracer()
    metrics := obs.Metrics()

    // Start service...
}
```

## Best Practices

1. **Always propagate context**: Pass context through all function calls to maintain trace context
2. **Use structured logging**: Add context as fields, not in the message string
3. **Record metrics consistently**: Use the same label values across the application
4. **Set meaningful span names**: Use descriptive names like "adapter.get-entity" not "fetch"
5. **Add relevant attributes**: Include entity type, ID, user ID, etc. in spans
6. **Handle errors properly**: Record errors on spans and log them
7. **Use sampling in production**: Set sampling rate to 0.1 (10%) or lower for high-volume services
8. **Register health checks early**: Add them during service initialization
9. **Close spans**: Always defer span.End() immediately after starting
10. **Test observability**: Verify traces, metrics, and logs in development

## Performance

- **Logging**: Zap provides excellent performance with minimal allocations
- **Tracing**: OpenTelemetry uses batching to minimize overhead (<1% typically)
- **Metrics**: Prometheus client is highly optimized for concurrent access
- **Health Checks**: Run in separate goroutines with timeout protection

## Dependencies

- `go.opentelemetry.io/otel` - OpenTelemetry SDK
- `go.uber.org/zap` - Structured logging
- `github.com/prometheus/client_golang` - Prometheus metrics
- Compatible with Jaeger, Grafana, and other observability tools

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
