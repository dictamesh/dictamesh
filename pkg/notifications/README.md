# DictaMesh Notifications Service

Comprehensive multi-channel notification system for the DictaMesh framework.

## Overview

The Notifications Service provides enterprise-grade notification capabilities for both infrastructure monitoring and application-level notifications. It supports multiple delivery channels, template management, user preferences, rate limiting, and comprehensive audit tracking.

## Features

- **Multi-Channel Support**: Email, SMS, Push (FCM/APNs), Slack, Webhooks, In-App, Browser Push, PagerDuty
- **Template Engine**: Flexible Go templates with multi-channel and i18n support
- **Event-Driven**: Kafka/Redpanda integration for real-time notification triggering
- **User Preferences**: Fine-grained control over notification delivery
- **Rate Limiting**: Per-user, per-channel, and system-wide rate limits
- **Retry & Fallback**: Automatic retry with exponential backoff and channel fallback
- **Batching**: Intelligent notification grouping for efficiency
- **Audit Trail**: Complete tracking of all notification events
- **Observability**: Prometheus metrics, distributed tracing, structured logging

## Package Structure

```
pkg/notifications/
├── README.md                 # This file
├── types.go                  # Core types and interfaces
├── config.go                 # Configuration structures
├── service.go                # Main service implementation
├── repository.go             # Data access layer
├── processor.go              # Notification processing logic
├── delivery.go               # Delivery management
├── template/                 # Template engine
│   ├── engine.go
│   └── renderer.go
├── channels/                 # Channel providers
│   ├── channel.go           # Base interface
│   ├── email/               # Email provider
│   ├── sms/                 # SMS provider
│   ├── push/                # Push notifications
│   ├── slack/               # Slack integration
│   ├── webhook/             # Webhook delivery
│   ├── inapp/               # In-app notifications
│   └── pagerduty/           # PagerDuty integration
├── rules/                    # Rule engine
│   ├── engine.go
│   └── matcher.go
├── ratelimit/               # Rate limiting
│   └── limiter.go
├── models/                  # Database models
│   └── notification.go
└── events/                  # Event handling
    ├── consumer.go
    └── producer.go
```

## Quick Start

### Installation

```go
import "github.com/click2-run/dictamesh/pkg/notifications"
```

### Basic Usage

```go
package main

import (
    "context"
    "log"

    "github.com/click2-run/dictamesh/pkg/notifications"
    "go.uber.org/zap"
)

func main() {
    // Create configuration
    config := notifications.DefaultConfig()
    config.DatabaseDSN = "postgres://user:pass@localhost/dictamesh"
    config.KafkaBootstrapServers = []string{"localhost:19092"}

    // Enable email channel
    config.Channels.Email.Enabled = true
    config.Channels.Email.Provider = "smtp"
    config.Channels.Email.SMTP = notifications.SMTPConfig{
        Host:     "smtp.example.com",
        Port:     587,
        Username: "user",
        Password: "pass",
        UseTLS:   true,
    }

    // Create logger
    logger, _ := zap.NewProduction()

    // Create service
    svc, err := notifications.NewService(config, logger)
    if err != nil {
        log.Fatal(err)
    }
    defer svc.Close()

    // Start service
    ctx := context.Background()
    if err := svc.Start(ctx); err != nil {
        log.Fatal(err)
    }

    // Send a notification
    notification, err := svc.SendNotification(ctx, &notifications.SendNotificationRequest{
        RecipientType: notifications.RecipientTypeUser,
        RecipientID:   "user123",
        Priority:      notifications.PriorityNormal,
        Channels:      []notifications.Channel{notifications.ChannelEmail},
        Subject:       "Welcome!",
        Body:          "Welcome to DictaMesh!",
    })

    if err != nil {
        log.Printf("Failed to send notification: %v", err)
    } else {
        log.Printf("Notification sent: %s", notification.ID)
    }
}
```

### Sending with Templates

```go
// Create a template
template := &notifications.NotificationTemplate{
    Name: "welcome-email",
    Channels: map[notifications.Channel]notifications.ChannelTemplate{
        notifications.ChannelEmail: {
            Subject:  "Welcome {{.UserName}}!",
            Body:     "Hi {{.UserName}},\n\nWelcome to our platform!",
            BodyHTML: "<h1>Welcome {{.UserName}}!</h1><p>Welcome to our platform!</p>",
        },
    },
}

err := svc.CreateTemplate(ctx, template)

// Send using template
notification, err := svc.SendNotification(ctx, &notifications.SendNotificationRequest{
    RecipientType: notifications.RecipientTypeUser,
    RecipientID:   "user123",
    Priority:      notifications.PriorityNormal,
    Channels:      []notifications.Channel{notifications.ChannelEmail},
    TemplateID:    template.ID,
    TemplateVars:  map[string]interface{}{
        "UserName": "John Doe",
    },
})
```

### Event-Driven Notifications

```go
// Create a notification rule
rule := &notifications.NotificationRule{
    Name:         "order-shipped",
    Description:  "Notify user when order is shipped",
    EventPattern: `event.type == "commerce.order.shipped"`,
    EventTypes:   []string{"commerce.order.shipped"},
    Priority:     notifications.PriorityNormal,
    Channels:     []notifications.Channel{notifications.ChannelEmail, notifications.ChannelPush},
    RecipientSelector: notifications.RecipientSelector{
        Type:       "dynamic",
        Expression: "event.data.customer_id",
    },
    TemplateID: "order-shipped",
    Enabled:    true,
}

err := svc.CreateRule(ctx, rule)

// Now when an event is published to Kafka topic "commerce.order.shipped",
// a notification will automatically be sent to the customer
```

## Configuration

See `config.go` for full configuration options. Key areas:

### Database

```go
config.DatabaseDSN = "postgres://user:pass@localhost:5432/dictamesh"
```

### Kafka

```go
config.KafkaBootstrapServers = []string{"localhost:19092"}
config.KafkaConsumerGroup = "dictamesh-notifications"
```

### Email Channel

```go
config.Channels.Email.Enabled = true
config.Channels.Email.Provider = "smtp" // or "ses" or "sendgrid"
config.Channels.Email.SMTP = notifications.SMTPConfig{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-app-password",
    UseTLS:   true,
}
config.Channels.Email.From = "DictaMesh <noreply@example.com>"
```

### SMS Channel

```go
config.Channels.SMS.Enabled = true
config.Channels.SMS.Provider = "twilio"
config.Channels.SMS.Twilio = notifications.TwilioConfig{
    AccountSID: "your-account-sid",
    AuthToken:  "your-auth-token",
    FromNumber: "+1234567890",
}
```

### Push Notifications

```go
config.Channels.Push.Enabled = true
config.Channels.Push.FCM = notifications.FCMConfig{
    Enabled:         true,
    CredentialsFile: "/path/to/firebase-credentials.json",
    Priority:        "high",
}
```

### Rate Limiting

```go
config.RateLimits.Enabled = true
config.RateLimits.UserLimits = map[notifications.Channel]notifications.RateLimitDefinition{
    notifications.ChannelEmail: {Count: 100, Duration: time.Hour},
    notifications.ChannelSMS:   {Count: 10, Duration: time.Hour},
}
```

## Architecture

### Event Flow

```
Event Source → Kafka → Event Consumer → Rule Matcher → Notification Processor
                                                              ↓
                                             Template Renderer ← Templates DB
                                                              ↓
                                             User Preferences ← Preferences DB
                                                              ↓
                                                Rate Limiter ← Redis
                                                              ↓
                                            Delivery Manager → Channel Queues
                                                              ↓
                                             Channel Providers → External APIs
                                                              ↓
                                            Delivery Tracking → Database
```

### Database Schema

The service uses PostgreSQL with the following tables:

- `dictamesh_notification_templates` - Notification templates
- `dictamesh_notification_rules` - Event-to-notification rules
- `dictamesh_notifications` - Notification instances (partitioned by month)
- `dictamesh_notification_delivery` - Delivery attempts
- `dictamesh_notification_preferences` - User preferences
- `dictamesh_notification_batches` - Batched notifications
- `dictamesh_notification_rate_limits` - Rate limit configuration
- `dictamesh_notification_audit` - Audit trail

### Channel Providers

Each channel provider implements the `ChannelProvider` interface:

```go
type ChannelProvider interface {
    Send(ctx context.Context, notification *Notification) (*DeliveryResult, error)
    GetChannel() Channel
    HealthCheck(ctx context.Context) error
}
```

Available providers:

- **Email**: SMTP, AWS SES, SendGrid, Mailgun
- **SMS**: Twilio, AWS SNS, MessageBird
- **Push**: Firebase (FCM), Apple (APNs), Web Push
- **Slack**: Webhooks, Bot API
- **Webhook**: Generic HTTP POST
- **In-App**: WebSocket, SSE
- **PagerDuty**: API integration

## Observability

### Metrics (Prometheus)

The service exposes metrics on port 9090 (configurable):

```
# Notification metrics
dictamesh_notifications_total{channel, priority, status}
dictamesh_notifications_delivery_duration_seconds{channel}
dictamesh_notifications_delivery_success_rate{channel}
dictamesh_notifications_queue_depth{channel, priority}

# Rate limit metrics
dictamesh_notifications_rate_limit_hits{scope, channel}

# Channel health
dictamesh_notifications_channel_health{channel}
```

### Distributed Tracing

All notifications are traced end-to-end using OpenTelemetry:

- Event consumption from Kafka
- Rule matching
- Template rendering
- User preference lookup
- Rate limiting check
- Channel delivery
- Database updates

### Logging

Structured logging with configurable levels:

```
{
  "level": "info",
  "ts": "2025-01-08T10:30:00Z",
  "msg": "notification sent",
  "notification_id": "uuid",
  "recipient_id": "user123",
  "channel": "EMAIL",
  "priority": "NORMAL",
  "delivery_time_ms": 245,
  "trace_id": "abc123"
}
```

## Use Cases

### Framework Infrastructure Alerts

Monitor framework health and alert operations team:

```go
// Database health alert rule
rule := &notifications.NotificationRule{
    Name:         "database-pool-exhausted",
    EventPattern: `event.type == "system.database.pool_exhausted"`,
    Priority:     notifications.PriorityCritical,
    Channels:     []notifications.Channel{
        notifications.ChannelSlack,
        notifications.ChannelEmail,
        notifications.ChannelPagerDuty,
    },
    RecipientSelector: notifications.RecipientSelector{
        Type:  "role",
        Roles: []string{"ops-team", "on-call"},
    },
    TemplateID: "infrastructure-alert",
    Enabled:    true,
}
```

### Application Notifications

User-facing notifications:

```go
// Order shipped notification
rule := &notifications.NotificationRule{
    Name:         "order-shipped",
    EventPattern: `event.type == "commerce.order.shipped"`,
    Priority:     notifications.PriorityNormal,
    Channels:     []notifications.Channel{
        notifications.ChannelEmail,
        notifications.ChannelPush,
    },
    FallbackChannels: []notifications.Channel{
        notifications.ChannelSMS,
    },
    RecipientSelector: notifications.RecipientSelector{
        Type:       "dynamic",
        Expression: "event.data.customer_id",
    },
    TemplateID: "order-shipped",
    Enabled:    true,
}
```

## Performance

Expected performance characteristics:

- **Throughput**: 10,000+ notifications/second (sustained)
- **Latency**: P99 < 1s (queue to delivery start)
- **Availability**: 99.95% SLA
- **Delivery Success**: 99.5%+ (per channel)

## Testing

Run tests:

```bash
cd pkg/notifications
go test ./...
```

Integration tests (requires infrastructure):

```bash
# Start infrastructure
cd ../../infrastructure
make dev-up

# Run integration tests
cd ../pkg/notifications
go test -tags=integration ./...
```

## Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for development guidelines.

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
