# Sentry Integration Guide for DictaMesh Framework

This guide explains how to integrate Sentry error tracking and monitoring into your DictaMesh framework components.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Integration by Language](#integration-by-language)
  - [Go](#go-integration)
  - [Node.js/TypeScript](#nodejs--typescript-integration)
  - [Python](#python-integration)
- [Configuration](#configuration)
- [Best Practices](#best-practices)
- [Advanced Features](#advanced-features)
- [Troubleshooting](#troubleshooting)

## Overview

DictaMesh includes a self-hosted Sentry instance for comprehensive error tracking and application monitoring across all framework components.

### What Sentry Provides

- **Error Tracking**: Automatic capture and reporting of exceptions and errors
- **Performance Monitoring**: Application Performance Monitoring (APM) with transaction tracing
- **Release Tracking**: Track errors and performance across different releases
- **Breadcrumbs**: Detailed context leading up to errors
- **User Feedback**: Collect user feedback on errors
- **Alerts**: Configurable alerts for error thresholds

### Architecture

```
┌─────────────────┐
│ Your Framework  │
│   Component     │──▶ Sentry SDK ──▶ Sentry Web (localhost:9000)
└─────────────────┘                   │
                                      ├─▶ PostgreSQL (metadata)
                                      ├─▶ ClickHouse (events)
                                      └─▶ Redis (cache)
```

## Quick Start

### 1. Start Sentry

```bash
cd infrastructure
make dev-up
make sentry-init  # First-time only
```

### 2. Create a Project

1. Open http://localhost:9000
2. Log in with `admin@dictamesh.local` / `admin`
3. Create a new project for your component
4. Select your platform (Go, Node.js, Python, etc.)
5. Copy the DSN (Data Source Name)

### 3. Integrate SDK

See language-specific sections below for integration details.

## Integration by Language

### Go Integration

#### Installation

```bash
go get github.com/getsentry/sentry-go
```

#### Basic Setup

```go
// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package main

import (
    "log"
    "time"

    "github.com/getsentry/sentry-go"
)

func main() {
    // Initialize Sentry
    err := sentry.Init(sentry.ClientOptions{
        Dsn: "http://your-dsn@localhost:9000/1",
        Environment: "development",
        Release: "dictamesh-adapter@1.0.0",
        // Enable performance monitoring
        EnableTracing: true,
        TracesSampleRate: 1.0, // 100% of transactions in dev
    })
    if err != nil {
        log.Fatalf("sentry.Init: %s", err)
    }
    // Flush buffered events before the program terminates
    defer sentry.Flush(2 * time.Second)

    // Your application code here
}
```

#### Error Capture

```go
// Capture an error
if err := someOperation(); err != nil {
    sentry.CaptureException(err)
}

// Capture a message
sentry.CaptureMessage("Something went wrong")

// With additional context
sentry.WithScope(func(scope *sentry.Scope) {
    scope.SetTag("component", "metadata-catalog")
    scope.SetExtra("entity_id", entityID)
    scope.SetLevel(sentry.LevelError)
    sentry.CaptureException(err)
})
```

#### Performance Monitoring

```go
// Start a transaction
span := sentry.StartSpan(ctx, "query.entity")
defer span.Finish()

// Add data to span
span.SetTag("entity_type", "customer")
span.SetData("query_time_ms", 42)

// Child spans for nested operations
childSpan := span.StartChild("db.query")
// ... database query ...
childSpan.Finish()
```

#### HTTP Middleware

```go
import (
    sentryhttp "github.com/getsentry/sentry-go/http"
)

// For net/http
handler := sentryhttp.New(sentryhttp.Options{}).Handle(yourHandler)

// For echo
e.Use(sentryecho.New(sentryecho.Options{}))

// For gin
r.Use(sentrygin.New(sentrygin.Options{}))
```

### Node.js / TypeScript Integration

#### Installation

```bash
npm install @sentry/node
# For Profiling
npm install @sentry/profiling-node
```

#### Basic Setup

```typescript
// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import * as Sentry from '@sentry/node';
import { ProfilingIntegration } from '@sentry/profiling-node';

// Initialize Sentry
Sentry.init({
  dsn: 'http://your-dsn@localhost:9000/1',
  environment: 'development',
  release: 'dictamesh-gateway@1.0.0',

  // Performance Monitoring
  tracesSampleRate: 1.0, // 100% in dev, lower in prod

  // Profiling
  profilesSampleRate: 1.0,
  integrations: [
    new ProfilingIntegration(),
  ],
});
```

#### Error Capture

```typescript
try {
  await someAsyncOperation();
} catch (error) {
  Sentry.captureException(error);
}

// With context
Sentry.withScope((scope) => {
  scope.setTag('component', 'graphql-gateway');
  scope.setExtra('query', queryString);
  scope.setLevel('error');
  Sentry.captureException(error);
});
```

#### Performance Monitoring

```typescript
// Start a transaction
const transaction = Sentry.startTransaction({
  op: 'graphql.query',
  name: 'GetEntity',
});

try {
  // Add child span
  const span = transaction.startChild({
    op: 'db.query',
    description: 'Fetch entity from database',
  });

  const result = await fetchEntity(id);

  span.finish();
  transaction.finish();

  return result;
} catch (error) {
  transaction.setStatus('internal_error');
  transaction.finish();
  throw error;
}
```

#### Express Middleware

```typescript
import express from 'express';
import * as Sentry from '@sentry/node';

const app = express();

// Request handler must be the first middleware
app.use(Sentry.Handlers.requestHandler());

// TracingHandler creates a trace for every incoming request
app.use(Sentry.Handlers.tracingHandler());

// Your routes here
app.get('/', (req, res) => {
  res.send('Hello World!');
});

// Error handler must be registered before any other error middleware
app.use(Sentry.Handlers.errorHandler());

app.listen(3000);
```

### Python Integration

#### Installation

```bash
pip install sentry-sdk
```

#### Basic Setup

```python
# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2025 Controle Digital Ltda

import sentry_sdk

# Initialize Sentry
sentry_sdk.init(
    dsn="http://your-dsn@localhost:9000/1",
    environment="development",
    release="dictamesh-service@1.0.0",

    # Performance Monitoring
    traces_sample_rate=1.0,  # 100% in dev, lower in prod

    # Profiling
    profiles_sample_rate=1.0,
)
```

#### Error Capture

```python
try:
    result = some_operation()
except Exception as e:
    sentry_sdk.capture_exception(e)

# With context
with sentry_sdk.push_scope() as scope:
    scope.set_tag("component", "event-router")
    scope.set_extra("event_type", event_type)
    scope.level = "error"
    sentry_sdk.capture_exception(e)
```

#### Performance Monitoring

```python
with sentry_sdk.start_transaction(op="task", name="process_event"):
    # Your code here
    with sentry_sdk.start_span(op="db.query", description="Fetch data"):
        data = fetch_from_db()

    with sentry_sdk.start_span(op="kafka.publish", description="Publish event"):
        publish_to_kafka(data)
```

#### Flask Integration

```python
from flask import Flask
import sentry_sdk
from sentry_sdk.integrations.flask import FlaskIntegration

sentry_sdk.init(
    dsn="http://your-dsn@localhost:9000/1",
    integrations=[FlaskIntegration()],
    traces_sample_rate=1.0,
)

app = Flask(__name__)
```

## Configuration

### Environment Variables

Create a `.env` file for your component:

```bash
# Sentry Configuration
SENTRY_DSN=http://your-dsn@localhost:9000/1
SENTRY_ENVIRONMENT=development
SENTRY_RELEASE=dictamesh-component@1.0.0
SENTRY_TRACES_SAMPLE_RATE=1.0
SENTRY_PROFILES_SAMPLE_RATE=1.0
```

### Configuration by Environment

#### Development

```yaml
dsn: http://your-dsn@localhost:9000/1
environment: development
traces_sample_rate: 1.0  # Capture 100% of transactions
profiles_sample_rate: 1.0
debug: true
```

#### Staging

```yaml
dsn: http://your-dsn@sentry-staging.dictamesh.local/1
environment: staging
traces_sample_rate: 0.5  # Capture 50% of transactions
profiles_sample_rate: 0.5
debug: false
```

#### Production

```yaml
dsn: http://your-dsn@sentry.dictamesh.io/1
environment: production
traces_sample_rate: 0.1  # Capture 10% of transactions
profiles_sample_rate: 0.1
debug: false
send_default_pii: false  # Don't send PII
```

## Best Practices

### 1. Error Context

Always provide context with errors:

```go
sentry.WithScope(func(scope *sentry.Scope) {
    scope.SetTag("component", "adapter")
    scope.SetTag("adapter_type", "rest")
    scope.SetExtra("endpoint", endpoint)
    scope.SetExtra("method", method)
    scope.SetUser(sentry.User{
        ID: userID,
        Email: userEmail,
    })
    sentry.CaptureException(err)
})
```

### 2. Performance Monitoring

Use transactions for critical operations:

```go
span := sentry.StartSpan(ctx, "adapter.fetch_entity")
defer span.Finish()

span.SetTag("entity_type", "customer")
span.SetData("entity_id", id)
```

### 3. Release Tracking

Use semantic versioning for releases:

```bash
SENTRY_RELEASE=dictamesh-adapter@1.2.3
```

### 4. Environment Separation

Use different projects for different environments:
- `dictamesh-adapter-dev`
- `dictamesh-adapter-staging`
- `dictamesh-adapter-prod`

### 5. Sample Rates

Adjust sample rates based on traffic:
- Development: 100% (1.0)
- Staging: 50% (0.5)
- Production: 10-25% (0.1-0.25)

### 6. Error Filtering

Filter out expected errors:

```go
sentry.Init(sentry.ClientOptions{
    BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
        if err, ok := hint.OriginalException.(error); ok {
            // Don't send validation errors
            if errors.Is(err, ErrValidation) {
                return nil
            }
        }
        return event
    },
})
```

### 7. Breadcrumbs

Add breadcrumbs for context:

```go
sentry.AddBreadcrumb(&sentry.Breadcrumb{
    Category: "auth",
    Message:  "User logged in",
    Level:    sentry.LevelInfo,
    Data: map[string]interface{}{
        "user_id": userID,
    },
})
```

## Advanced Features

### Custom Tags

```go
// Framework-specific tags
sentry.ConfigureScope(func(scope *sentry.Scope) {
    scope.SetTag("framework", "dictamesh")
    scope.SetTag("component_type", "adapter")
    scope.SetTag("data_source", "rest_api")
})
```

### User Context

```go
sentry.ConfigureScope(func(scope *sentry.Scope) {
    scope.SetUser(sentry.User{
        ID:       userID,
        Email:    userEmail,
        Username: username,
        IPAddress: ipAddress,
    })
})
```

### Custom Fingerprinting

Group similar errors together:

```go
sentry.WithScope(func(scope *sentry.Scope) {
    scope.SetFingerprint([]string{
        "database-connection-error",
        dbHost,
    })
    sentry.CaptureException(err)
})
```

### Distributed Tracing

For microservices architecture:

```go
// Service A
span := sentry.StartSpan(ctx, "service_a.operation")
defer span.Finish()

// Pass trace information to Service B
traceID := span.TraceID

// Service B
ctx = sentry.SetHubOnContext(ctx, hub)
span := sentry.StartSpan(ctx, "service_b.operation")
defer span.Finish()
```

## Troubleshooting

### Events Not Appearing

1. **Check DSN**: Ensure the DSN is correct
2. **Check Sentry is running**: `curl http://localhost:9000/_health/`
3. **Enable debug mode**:
   ```go
   sentry.Init(sentry.ClientOptions{
       Debug: true,
   })
   ```
4. **Check network connectivity**: Ensure your service can reach Sentry

### Performance Issues

1. **Lower sample rates**: Reduce `traces_sample_rate`
2. **Filter transactions**: Don't trace health check endpoints
3. **Async sending**: Ensure SDK is sending events asynchronously

### Missing Context

1. **Add breadcrumbs**: Use breadcrumbs for important operations
2. **Set scope data**: Add relevant data to scope before capturing
3. **Use tags**: Tag events with component, environment, etc.

## Integration with DictaMesh Framework

### Adapter Integration

```go
// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package adapter

import (
    "context"
    "github.com/getsentry/sentry-go"
)

type BaseAdapter struct {
    sentryHub *sentry.Hub
}

func (a *BaseAdapter) GetEntity(ctx context.Context, id string) (*Entity, error) {
    span := sentry.StartSpan(ctx, "adapter.get_entity")
    defer span.Finish()

    span.SetTag("entity_id", id)

    entity, err := a.fetchEntity(ctx, id)
    if err != nil {
        sentry.CaptureException(err)
        return nil, err
    }

    return entity, nil
}
```

### Event Bus Integration

```go
func (e *EventBus) Publish(ctx context.Context, event Event) error {
    span := sentry.StartSpan(ctx, "eventbus.publish")
    defer span.Finish()

    span.SetTag("event_type", event.Type)

    if err := e.kafka.Publish(ctx, event); err != nil {
        sentry.WithScope(func(scope *sentry.Scope) {
            scope.SetExtra("event", event)
            sentry.CaptureException(err)
        })
        return err
    }

    return nil
}
```

### GraphQL Gateway Integration

```typescript
const server = new ApolloServer({
  typeDefs,
  resolvers,
  plugins: [
    {
      async requestDidStart() {
        return {
          async didEncounterErrors(requestContext) {
            for (const error of requestContext.errors) {
              Sentry.withScope((scope) => {
                scope.setTag('kind', 'graphql_error');
                scope.setExtra('query', requestContext.request.query);
                scope.setExtra('variables', requestContext.request.variables);
                Sentry.captureException(error);
              });
            }
          },
        };
      },
    },
  ],
});
```

## Resources

- [Sentry Documentation](https://docs.sentry.io/)
- [Sentry Go SDK](https://docs.sentry.io/platforms/go/)
- [Sentry Node.js SDK](https://docs.sentry.io/platforms/node/)
- [Sentry Python SDK](https://docs.sentry.io/platforms/python/)
- [DictaMesh Sentry Configuration](../infrastructure/docker-compose/sentry/README.md)

## Support

For issues with Sentry integration:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review Sentry logs: `make sentry-logs`
3. Check Sentry UI for error details
4. Consult the official Sentry documentation

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
