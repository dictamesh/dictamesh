# DictaMesh Notifications Service

## Executive Summary

The DictaMesh Notifications Service is a comprehensive, multi-channel notification system designed to support both framework core operations (infrastructure monitoring, system alerts, technical operations) and extensible application-level notifications for solutions built with the framework.

**Version:** 1.0.0
**Status:** Planning
**Authors:** DictaMesh Core Team

## Table of Contents

1. [Overview](#overview)
2. [Use Cases](#use-cases)
3. [Architecture](#architecture)
4. [Channel Support](#channel-support)
5. [Core Components](#core-components)
6. [Database Schema](#database-schema)
7. [API Design](#api-design)
8. [Event Integration](#event-integration)
9. [Security & Compliance](#security--compliance)
10. [Deployment & Scaling](#deployment--scaling)

## Overview

### Purpose

The Notifications Service provides:

1. **Framework Core Survival**: Critical infrastructure monitoring and alerting
   - Database health alerts
   - Kafka/Redpanda connectivity issues
   - Service degradation warnings
   - Circuit breaker state changes
   - Performance threshold violations
   - Security incidents

2. **Extensible Application Layer**: Customizable notifications for solutions
   - Business event notifications
   - User-facing alerts
   - Workflow notifications
   - Custom integrations

### Key Features

- **Multi-Channel Support**: Email, SMS, Push (mobile/web), Slack, Webhooks, In-App
- **Event-Driven Architecture**: Kafka/Redpanda integration for real-time processing
- **Template Engine**: Flexible templating with i18n support
- **Priority & Routing**: Smart routing based on priority, preferences, and channels
- **Delivery Tracking**: Full audit trail and delivery status
- **Rate Limiting**: Per-channel and per-user rate limiting
- **Retry & Fallback**: Automatic retry with exponential backoff, channel fallback
- **Batching & Aggregation**: Intelligent notification grouping
- **User Preferences**: Fine-grained user notification preferences
- **A/B Testing**: Template and channel testing support

## Use Cases

### Framework Core Operations

#### 1. Infrastructure Health Monitoring

```
Scenario: PostgreSQL connection pool exhaustion
↓
Event: system.database.pool_exhausted
↓
Notification:
- Channel: Slack (#ops-critical)
- Channel: Email (ops@company.com)
- Channel: PagerDuty (on-call)
- Priority: CRITICAL
- Template: infrastructure-alert
```

#### 2. Circuit Breaker Alerts

```
Scenario: Adapter circuit breaker opens
↓
Event: system.circuit_breaker.opened
↓
Notification:
- Channel: Slack (#ops-alerts)
- Channel: Email (team@company.com)
- Priority: HIGH
- Template: circuit-breaker-alert
```

#### 3. Security Incidents

```
Scenario: Unauthorized access attempt
↓
Event: system.security.unauthorized_access
↓
Notification:
- Channel: Email (security@company.com)
- Channel: SMS (security team)
- Channel: PagerDuty
- Priority: CRITICAL
- Template: security-alert
```

### Application-Level Notifications

#### 4. Business Process Notifications

```
Scenario: Order shipped
↓
Event: commerce.order.shipped
↓
Notification:
- Channel: Email (customer)
- Channel: Push (mobile app)
- Channel: SMS (if opted in)
- Priority: NORMAL
- Template: order-shipped
```

#### 5. Workflow Updates

```
Scenario: Approval request pending
↓
Event: workflow.approval.requested
↓
Notification:
- Channel: Email (approver)
- Channel: In-App (dashboard)
- Channel: Slack (DM)
- Priority: HIGH
- Template: approval-request
```

#### 6. Data Pipeline Events

```
Scenario: ETL pipeline failed
↓
Event: pipeline.etl.failed
↓
Notification:
- Channel: Slack (#data-eng)
- Channel: Email (data-team@company.com)
- Priority: HIGH
- Template: pipeline-failure
```

## Architecture

### High-Level Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                     EVENT SOURCES                            │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Framework  │  │  Adapters    │  │ Applications │      │
│  │  Core       │  │              │  │              │      │
│  └──────┬──────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                │                  │              │
└─────────┼────────────────┼──────────────────┼──────────────┘
          │                │                  │
          └────────────────┴──────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                  KAFKA / REDPANDA                            │
│                                                              │
│  Topics:                                                     │
│  - system.notifications.requested                           │
│  - system.notifications.sent                                │
│  - system.notifications.failed                              │
│  - {domain}.events (business events)                        │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│           NOTIFICATIONS SERVICE CORE                         │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │         Event Consumer & Router                    │     │
│  │  - Event validation                                │     │
│  │  - Rule matching                                   │     │
│  │  - Priority assignment                             │     │
│  └────────────────┬───────────────────────────────────┘     │
│                   │                                          │
│  ┌────────────────▼───────────────────────────────────┐     │
│  │       Notification Processor                       │     │
│  │  - Template rendering                              │     │
│  │  - User preference checking                        │     │
│  │  - Rate limiting                                   │     │
│  │  - Batching & aggregation                          │     │
│  └────────────────┬───────────────────────────────────┘     │
│                   │                                          │
│  ┌────────────────▼───────────────────────────────────┐     │
│  │         Delivery Manager                           │     │
│  │  - Channel selection                               │     │
│  │  - Retry logic                                     │     │
│  │  - Fallback handling                               │     │
│  │  - Delivery tracking                               │     │
│  └────────────────┬───────────────────────────────────┘     │
└───────────────────┼──────────────────────────────────────────┘
                    │
          ┌─────────┴─────────┐
          │                   │
┌─────────▼────────┐ ┌───────▼──────────┐
│  Channel Router  │ │  Channel Router  │
│                  │ │                  │
│  Queues:         │ │  Queues:         │
│  - High Priority │ │  - Normal        │
│  - Time-based    │ │  - Batch         │
└─────────┬────────┘ └───────┬──────────┘
          │                   │
          └─────────┬─────────┘
                    │
┌───────────────────▼──────────────────────────────────────────┐
│               CHANNEL PROVIDERS                              │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │  Email   │ │   SMS    │ │   Push   │ │  Slack   │      │
│  │          │ │          │ │          │ │          │      │
│  │ - SMTP   │ │ - Twilio │ │ - FCM    │ │ - Webhook│      │
│  │ - AWS SES│ │ - AWS SNS│ │ - APNs   │ │ - Bot API│      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
│                                                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │ Webhooks │ │  In-App  │ │  Browser │ │PagerDuty │      │
│  │          │ │          │ │          │ │          │      │
│  │ - HTTP   │ │ - WebSock│ │ - Web    │ │ - API    │      │
│  │ - Retry  │ │ - SSE    │ │ - Push   │ │          │      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
└──────────────────────────────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────────┐
│              STORAGE & STATE                                 │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────┐  ┌────────────────────────────┐      │
│  │   PostgreSQL     │  │   Redis Cache              │      │
│  │                  │  │                            │      │
│  │ - Notifications  │  │ - Rate limiting counters   │      │
│  │ - Templates      │  │ - Delivery locks           │      │
│  │ - User prefs     │  │ - Active sessions          │      │
│  │ - Audit logs     │  │ - Batching queues          │      │
│  └──────────────────┘  └────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

#### 1. Event Consumer & Router

**Responsibilities:**
- Subscribe to Kafka topics for notification events
- Validate event schemas
- Match events to notification rules
- Assign priorities based on rules
- Forward to Notification Processor

**Technology:**
- Kafka consumer groups
- Event schema validation (Avro)
- Rule engine (CEL or similar)

#### 2. Notification Processor

**Responsibilities:**
- Render templates with event data
- Check user notification preferences
- Apply rate limiting (per-user, per-channel)
- Batch notifications when appropriate
- Aggregate similar notifications

**Technology:**
- Template engine (Go templates, Handlebars)
- Redis for rate limiting
- PostgreSQL for preferences

#### 3. Delivery Manager

**Responsibilities:**
- Select appropriate channels
- Queue notifications by priority
- Implement retry logic with exponential backoff
- Handle channel fallbacks
- Track delivery status
- Update audit logs

**Technology:**
- Channel-specific queues (Redis)
- Retry mechanisms
- Circuit breakers for channel providers

#### 4. Channel Providers

**Responsibilities:**
- Abstract channel-specific APIs
- Handle authentication and credentials
- Format messages for each channel
- Report delivery status
- Handle provider-specific errors

**Supported Channels:**
- Email (SMTP, AWS SES, SendGrid)
- SMS (Twilio, AWS SNS)
- Push (FCM for Android, APNs for iOS)
- Slack (Incoming Webhooks, Bot API)
- Webhooks (HTTP POST)
- In-App (WebSocket, SSE)
- Browser Push (Web Push API)
- PagerDuty (API)

## Channel Support

### Email Channel

**Providers:**
- SMTP (standard)
- AWS SES
- SendGrid
- Mailgun

**Features:**
- HTML and plain text
- Attachments
- CC/BCC support
- Reply-to handling
- Bounce tracking
- Open/click tracking (optional)

**Configuration:**
```yaml
email:
  provider: smtp  # smtp | ses | sendgrid
  smtp:
    host: smtp.example.com
    port: 587
    username: ${SMTP_USERNAME}
    password: ${SMTP_PASSWORD}
    from: "DictaMesh <noreply@example.com>"
  rate_limit:
    per_user: 100/hour
    per_system: 10000/hour
```

### SMS Channel

**Providers:**
- Twilio
- AWS SNS
- MessageBird

**Features:**
- International numbers
- Delivery receipts
- Fallback to voice

**Configuration:**
```yaml
sms:
  provider: twilio  # twilio | sns | messagebird
  twilio:
    account_sid: ${TWILIO_ACCOUNT_SID}
    auth_token: ${TWILIO_AUTH_TOKEN}
    from: "+1234567890"
  rate_limit:
    per_user: 10/hour
    per_system: 1000/hour
  max_length: 160
```

### Push Notification Channel

**Providers:**
- Firebase Cloud Messaging (FCM) - Android
- Apple Push Notification Service (APNs) - iOS
- Web Push API - Browser

**Features:**
- Badge updates
- Sound/vibration
- Custom data payloads
- Priority handling
- Time-to-live

**Configuration:**
```yaml
push:
  fcm:
    credentials_file: /secrets/fcm-credentials.json
    priority: high
  apns:
    certificate_file: /secrets/apns-cert.p12
    certificate_password: ${APNS_CERT_PASSWORD}
    production: true
  web_push:
    vapid_public_key: ${VAPID_PUBLIC_KEY}
    vapid_private_key: ${VAPID_PRIVATE_KEY}
```

### Slack Channel

**Integration Types:**
- Incoming Webhooks (simple)
- Bot API (interactive)

**Features:**
- Channel posting
- Direct messages
- Interactive buttons
- Thread replies
- Rich formatting (blocks)

**Configuration:**
```yaml
slack:
  webhook_url: ${SLACK_WEBHOOK_URL}
  bot_token: ${SLACK_BOT_TOKEN}
  default_channel: "#alerts"
  rate_limit:
    per_channel: 1/second
```

### Webhook Channel

**Features:**
- Custom HTTP endpoints
- Authentication (API key, OAuth)
- Retry with backoff
- Signature verification
- Custom headers

**Configuration:**
```yaml
webhooks:
  timeout: 10s
  retry:
    max_attempts: 3
    initial_interval: 1s
    multiplier: 2
  auth:
    type: bearer  # bearer | apikey | oauth2
    token: ${WEBHOOK_AUTH_TOKEN}
```

### In-App Channel

**Technologies:**
- WebSocket for real-time
- Server-Sent Events (SSE) for fallback
- Long polling for legacy browsers

**Features:**
- Real-time delivery
- Notification center UI
- Read/unread tracking
- Persistence

**Configuration:**
```yaml
in_app:
  transport: websocket  # websocket | sse | longpoll
  persistence: 30d
  max_unread: 1000
```

### Browser Push Channel

**Features:**
- Background notifications
- Click actions
- Icons and images
- Persistence

**Requirements:**
- HTTPS
- Service Worker
- User permission

## Core Components

### 1. Notification Rules Engine

Define when and how to send notifications:

```go
type NotificationRule struct {
    ID          string
    Name        string
    Description string

    // Trigger
    EventPattern string // CEL expression
    Domains      []string
    EventTypes   []string

    // Routing
    Priority     Priority
    Channels     []Channel
    Fallback     []Channel

    // Targeting
    Recipients   RecipientSelector

    // Timing
    Schedule     *Schedule // Optional: schedule for delivery
    Timezone     string

    // Batching
    BatchWindow  time.Duration
    BatchSize    int

    // Template
    TemplateID   string
    TemplateVars map[string]interface{}

    // Lifecycle
    Enabled      bool
    ValidFrom    time.Time
    ValidUntil   *time.Time
}

type RecipientSelector struct {
    Type string // role | user | group | dynamic

    // Static recipients
    UserIDs  []string
    Roles    []string
    Groups   []string

    // Dynamic recipients (evaluated at runtime)
    Expression string // CEL expression
}
```

**Example Rules:**

```yaml
# Framework core: Database health alert
- name: database-pool-exhaustion
  event_pattern: event.type == "system.database.pool_exhausted"
  priority: CRITICAL
  channels:
    - slack
    - email
    - pagerduty
  recipients:
    type: role
    roles: [ops-team, on-call]
  template_id: infrastructure-alert

# Application: Order shipped
- name: order-shipped-notification
  event_pattern: event.type == "commerce.order.shipped"
  domains: [commerce]
  priority: NORMAL
  channels:
    - email
    - push
  fallback:
    - sms
  recipients:
    type: dynamic
    expression: event.data.customer_id
  batch_window: 5m
  template_id: order-shipped
```

### 2. Template Engine

**Template Structure:**

```go
type NotificationTemplate struct {
    ID          string
    Name        string
    Description string

    // Content
    Channels    map[Channel]ChannelTemplate

    // Localization
    Locale      string
    Translations map[string]LocalizedTemplate

    // Metadata
    Version     string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type ChannelTemplate struct {
    Subject      string // For email, push title
    Body         string // Main content
    BodyHTML     string // HTML version (email)
    Data         map[string]interface{} // Channel-specific data
}

type LocalizedTemplate struct {
    Subject      string
    Body         string
    BodyHTML     string
}
```

**Example Template:**

```yaml
id: infrastructure-alert
name: Infrastructure Alert
channels:
  email:
    subject: "[{{.Priority}}] {{.Title}}"
    body_html: |
      <h1>{{.Title}}</h1>
      <p><strong>Priority:</strong> {{.Priority}}</p>
      <p><strong>Service:</strong> {{.Service}}</p>
      <p><strong>Message:</strong> {{.Message}}</p>
      <p><strong>Time:</strong> {{.Timestamp}}</p>
      {{if .Metrics}}
      <h2>Metrics</h2>
      <ul>
      {{range $key, $value := .Metrics}}
        <li><strong>{{$key}}:</strong> {{$value}}</li>
      {{end}}
      </ul>
      {{end}}
      <p><a href="{{.DashboardURL}}">View Dashboard</a></p>
    body: |
      {{.Title}}

      Priority: {{.Priority}}
      Service: {{.Service}}
      Message: {{.Message}}
      Time: {{.Timestamp}}

      Dashboard: {{.DashboardURL}}

  slack:
    body: |
      {
        "blocks": [
          {
            "type": "header",
            "text": {
              "type": "plain_text",
              "text": "{{.Title}}"
            }
          },
          {
            "type": "section",
            "fields": [
              {"type": "mrkdwn", "text": "*Priority:*\n{{.Priority}}"},
              {"type": "mrkdwn", "text": "*Service:*\n{{.Service}}"}
            ]
          },
          {
            "type": "section",
            "text": {
              "type": "mrkdwn",
              "text": "{{.Message}}"
            }
          },
          {
            "type": "actions",
            "elements": [
              {
                "type": "button",
                "text": {"type": "plain_text", "text": "View Dashboard"},
                "url": "{{.DashboardURL}}"
              }
            ]
          }
        ]
      }
```

### 3. User Preferences

Allow users to control notification delivery:

```go
type UserPreferences struct {
    UserID    string

    // Global settings
    Enabled   bool
    Timezone  string
    Locale    string

    // Channel preferences
    Channels  map[Channel]ChannelPreference

    // Do not disturb
    QuietHours *QuietHours

    // Category preferences
    Categories map[string]CategoryPreference
}

type ChannelPreference struct {
    Enabled bool
    Address string // Email, phone, device token
}

type QuietHours struct {
    Enabled   bool
    StartTime string // "22:00"
    EndTime   string // "08:00"
    Timezone  string
    AllowCritical bool // Allow CRITICAL notifications
}

type CategoryPreference struct {
    Enabled      bool
    Channels     []Channel
    MinPriority  Priority
}
```

### 4. Rate Limiting

Prevent notification spam:

```go
type RateLimiter struct {
    // Per-user limits
    UserLimits map[Channel]RateLimit

    // System-wide limits
    SystemLimits map[Channel]RateLimit

    // Category limits
    CategoryLimits map[string]RateLimit
}

type RateLimit struct {
    Count     int
    Duration  time.Duration
    Bucket    string // Redis key pattern
}
```

**Example Limits:**

```yaml
rate_limits:
  user:
    email: 100/hour
    sms: 10/hour
    push: 50/hour
  system:
    email: 10000/hour
    sms: 1000/hour
    slack: 100/minute
  category:
    marketing: 5/day
    transactional: 1000/hour
```

### 5. Delivery Tracking

Track notification lifecycle:

```go
type NotificationStatus string

const (
    StatusPending    NotificationStatus = "pending"
    StatusQueued     NotificationStatus = "queued"
    StatusSending    NotificationStatus = "sending"
    StatusSent       NotificationStatus = "sent"
    StatusDelivered  NotificationStatus = "delivered"
    StatusFailed     NotificationStatus = "failed"
    StatusRetrying   NotificationStatus = "retrying"
    StatusCancelled  NotificationStatus = "cancelled"
)

type DeliveryRecord struct {
    NotificationID string
    Channel        Channel
    Status         NotificationStatus
    Attempts       int
    LastAttempt    time.Time
    NextRetry      *time.Time
    Error          *string
    Metadata       map[string]interface{}
}
```

## Database Schema

### Core Tables

```sql
-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- Notification Templates
CREATE TABLE dictamesh_notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,

    -- Content (JSONB for flexibility)
    channels JSONB NOT NULL, -- {email: {...}, slack: {...}}
    translations JSONB, -- {en: {...}, pt: {...}}

    -- Template metadata
    variables JSONB, -- Expected variables
    schema_version VARCHAR(50) DEFAULT '1.0',

    -- Lifecycle
    version VARCHAR(50) DEFAULT '1.0.0',
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by VARCHAR(255),

    -- Indexing
    tags TEXT[],

    CONSTRAINT valid_channels CHECK (jsonb_typeof(channels) = 'object')
);

CREATE INDEX idx_dictamesh_notification_templates_name
    ON dictamesh_notification_templates(name);
CREATE INDEX idx_dictamesh_notification_templates_tags
    ON dictamesh_notification_templates USING GIN(tags);

-- Notification Rules
CREATE TABLE dictamesh_notification_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,

    -- Trigger conditions
    event_pattern TEXT NOT NULL, -- CEL expression
    domains TEXT[],
    event_types TEXT[],

    -- Routing
    priority VARCHAR(20) NOT NULL, -- CRITICAL, HIGH, NORMAL, LOW
    channels TEXT[] NOT NULL,
    fallback_channels TEXT[],

    -- Recipients
    recipient_selector JSONB NOT NULL,

    -- Timing
    schedule JSONB, -- Cron or interval
    timezone VARCHAR(50) DEFAULT 'UTC',

    -- Batching
    batch_window_seconds INTEGER,
    batch_size INTEGER,

    -- Template
    template_id UUID REFERENCES dictamesh_notification_templates(id),
    template_vars JSONB,

    -- Lifecycle
    enabled BOOLEAN DEFAULT TRUE,
    valid_from TIMESTAMPTZ DEFAULT NOW(),
    valid_until TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_priority CHECK (priority IN ('CRITICAL', 'HIGH', 'NORMAL', 'LOW'))
);

CREATE INDEX idx_dictamesh_notification_rules_enabled
    ON dictamesh_notification_rules(enabled) WHERE enabled = TRUE;
CREATE INDEX idx_dictamesh_notification_rules_domains
    ON dictamesh_notification_rules USING GIN(domains);
CREATE INDEX idx_dictamesh_notification_rules_event_types
    ON dictamesh_notification_rules USING GIN(event_types);

-- Notifications (main tracking table)
CREATE TABLE dictamesh_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Source
    event_id VARCHAR(255), -- Original event ID
    rule_id UUID REFERENCES dictamesh_notification_rules(id),
    template_id UUID REFERENCES dictamesh_notification_templates(id),

    -- Recipients
    recipient_type VARCHAR(50) NOT NULL, -- user, role, group, system
    recipient_id VARCHAR(255) NOT NULL,

    -- Content
    subject TEXT,
    body TEXT,
    body_html TEXT,
    data JSONB,

    -- Routing
    priority VARCHAR(20) NOT NULL,
    channels TEXT[] NOT NULL,
    selected_channel VARCHAR(50), -- Actually used channel

    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    -- Timing
    scheduled_at TIMESTAMPTZ DEFAULT NOW(),
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    read_at TIMESTAMPTZ,

    -- Metadata
    metadata JSONB,
    trace_id VARCHAR(64), -- OpenTelemetry trace ID

    -- Error tracking
    error TEXT,
    retry_count INTEGER DEFAULT 0,
    next_retry_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_status CHECK (status IN (
        'pending', 'queued', 'sending', 'sent',
        'delivered', 'failed', 'retrying', 'cancelled'
    ))
);

-- Partitioning by month for scalability
CREATE TABLE dictamesh_notifications_y2025m01 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE INDEX idx_dictamesh_notifications_recipient
    ON dictamesh_notifications(recipient_type, recipient_id, created_at DESC);
CREATE INDEX idx_dictamesh_notifications_status
    ON dictamesh_notifications(status, scheduled_at);
CREATE INDEX idx_dictamesh_notifications_event
    ON dictamesh_notifications(event_id);
CREATE INDEX idx_dictamesh_notifications_trace
    ON dictamesh_notifications(trace_id);

-- Delivery Attempts (detailed delivery tracking)
CREATE TABLE dictamesh_notification_delivery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL REFERENCES dictamesh_notifications(id),

    -- Delivery details
    channel VARCHAR(50) NOT NULL,
    provider VARCHAR(100), -- smtp, twilio, fcm, etc.

    -- Status
    status VARCHAR(20) NOT NULL,
    attempt_number INTEGER NOT NULL,

    -- Timing
    started_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    -- Result
    success BOOLEAN,
    error TEXT,
    provider_response JSONB,
    provider_message_id VARCHAR(255), -- For tracking with provider

    -- Metadata
    metadata JSONB,

    CONSTRAINT valid_delivery_status CHECK (status IN (
        'sending', 'sent', 'delivered', 'failed', 'bounced', 'rejected'
    ))
);

CREATE INDEX idx_dictamesh_notification_delivery_notification
    ON dictamesh_notification_delivery(notification_id, attempt_number DESC);
CREATE INDEX idx_dictamesh_notification_delivery_provider_id
    ON dictamesh_notification_delivery(provider_message_id);

-- User Preferences
CREATE TABLE dictamesh_notification_preferences (
    user_id VARCHAR(255) PRIMARY KEY,

    -- Global settings
    enabled BOOLEAN DEFAULT TRUE,
    timezone VARCHAR(50) DEFAULT 'UTC',
    locale VARCHAR(10) DEFAULT 'en',

    -- Channel addresses
    email VARCHAR(255),
    phone VARCHAR(20),
    push_tokens JSONB, -- Array of device tokens

    -- Channel preferences
    channel_prefs JSONB DEFAULT '{}', -- {email: {enabled: true}, ...}

    -- Quiet hours
    quiet_hours_enabled BOOLEAN DEFAULT FALSE,
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    quiet_hours_allow_critical BOOLEAN DEFAULT TRUE,

    -- Category preferences
    category_prefs JSONB DEFAULT '{}',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_dictamesh_notification_preferences_email
    ON dictamesh_notification_preferences(email);
CREATE INDEX idx_dictamesh_notification_preferences_phone
    ON dictamesh_notification_preferences(phone);

-- Notification Batches (for grouped notifications)
CREATE TABLE dictamesh_notification_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Batch config
    rule_id UUID REFERENCES dictamesh_notification_rules(id),
    batch_key VARCHAR(255) NOT NULL, -- Grouping key

    -- Timing
    window_start TIMESTAMPTZ NOT NULL,
    window_end TIMESTAMPTZ NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    sent_at TIMESTAMPTZ,

    -- Content
    notification_ids UUID[], -- Individual notifications in batch
    count INTEGER NOT NULL,

    -- Status
    status VARCHAR(20) DEFAULT 'pending',

    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_batch_status CHECK (status IN ('pending', 'sent', 'failed'))
);

CREATE INDEX idx_dictamesh_notification_batches_key_window
    ON dictamesh_notification_batches(batch_key, window_end);
CREATE INDEX idx_dictamesh_notification_batches_scheduled
    ON dictamesh_notification_batches(status, scheduled_at);

-- Rate Limiting Counters (Redis-backed, table for config)
CREATE TABLE dictamesh_notification_rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    scope VARCHAR(50) NOT NULL, -- user, system, category
    scope_id VARCHAR(255), -- user_id, category name, etc.
    channel VARCHAR(50) NOT NULL,

    -- Limit definition
    max_count INTEGER NOT NULL,
    window_seconds INTEGER NOT NULL,

    -- Metadata
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(scope, scope_id, channel)
);

-- Audit Log (comprehensive tracking)
CREATE TABLE dictamesh_notification_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    notification_id UUID REFERENCES dictamesh_notifications(id),

    -- Event details
    event_type VARCHAR(100) NOT NULL, -- created, sent, delivered, failed, etc.
    actor_type VARCHAR(50), -- system, user, api
    actor_id VARCHAR(255),

    -- Details
    details JSONB,

    -- Timing
    timestamp TIMESTAMPTZ DEFAULT NOW(),

    -- Tracing
    trace_id VARCHAR(64)
);

CREATE INDEX idx_dictamesh_notification_audit_notification
    ON dictamesh_notification_audit(notification_id, timestamp DESC);
CREATE INDEX idx_dictamesh_notification_audit_type
    ON dictamesh_notification_audit(event_type, timestamp DESC);

-- Comments on tables (for documentation)
COMMENT ON TABLE dictamesh_notification_templates IS
    'DictaMesh: Notification template definitions with multi-channel support';
COMMENT ON TABLE dictamesh_notification_rules IS
    'DictaMesh: Rules for triggering notifications based on events';
COMMENT ON TABLE dictamesh_notifications IS
    'DictaMesh: Notification instances and delivery tracking';
COMMENT ON TABLE dictamesh_notification_delivery IS
    'DictaMesh: Detailed delivery attempts and provider responses';
COMMENT ON TABLE dictamesh_notification_preferences IS
    'DictaMesh: User notification preferences and settings';
```

## API Design

### REST API

```go
// Notification Management API
type NotificationAPI interface {
    // Send notification directly
    SendNotification(ctx context.Context, req *SendNotificationRequest) (*Notification, error)

    // Bulk send
    SendBulkNotifications(ctx context.Context, req *BulkSendRequest) (*BulkSendResponse, error)

    // Query notifications
    ListNotifications(ctx context.Context, filters NotificationFilters) (*NotificationList, error)
    GetNotification(ctx context.Context, id string) (*Notification, error)

    // Cancel notification
    CancelNotification(ctx context.Context, id string) error

    // Templates
    CreateTemplate(ctx context.Context, template *NotificationTemplate) (*NotificationTemplate, error)
    UpdateTemplate(ctx context.Context, id string, template *NotificationTemplate) (*NotificationTemplate, error)
    GetTemplate(ctx context.Context, id string) (*NotificationTemplate, error)
    ListTemplates(ctx context.Context, filters TemplateFilters) (*TemplateList, error)
    DeleteTemplate(ctx context.Context, id string) error

    // Rules
    CreateRule(ctx context.Context, rule *NotificationRule) (*NotificationRule, error)
    UpdateRule(ctx context.Context, id string, rule *NotificationRule) (*NotificationRule, error)
    GetRule(ctx context.Context, id string) (*NotificationRule, error)
    ListRules(ctx context.Context, filters RuleFilters) (*RuleList, error)
    DeleteRule(ctx context.Context, id string) error

    // User Preferences
    GetPreferences(ctx context.Context, userID string) (*UserPreferences, error)
    UpdatePreferences(ctx context.Context, userID string, prefs *UserPreferences) (*UserPreferences, error)

    // Statistics
    GetStats(ctx context.Context, filters StatsFilters) (*NotificationStats, error)
}
```

### GraphQL Schema

```graphql
type Notification {
  id: ID!
  eventId: String
  recipientType: RecipientType!
  recipientId: String!

  subject: String
  body: String
  bodyHtml: String
  data: JSON

  priority: Priority!
  channels: [Channel!]!
  selectedChannel: Channel

  status: NotificationStatus!

  scheduledAt: DateTime!
  sentAt: DateTime
  deliveredAt: DateTime
  readAt: DateTime

  retryCount: Int!
  nextRetryAt: DateTime

  error: String
  metadata: JSON

  # Relationships
  rule: NotificationRule
  template: NotificationTemplate
  deliveryAttempts: [DeliveryAttempt!]!
}

type NotificationTemplate {
  id: ID!
  name: String!
  description: String

  channels: JSON!
  translations: JSON

  version: String!
  enabled: Boolean!

  createdAt: DateTime!
  updatedAt: DateTime!
}

type NotificationRule {
  id: ID!
  name: String!
  description: String

  eventPattern: String!
  domains: [String!]
  eventTypes: [String!]

  priority: Priority!
  channels: [Channel!]!
  fallbackChannels: [Channel!]

  recipientSelector: JSON!
  schedule: JSON

  batchWindow: Int
  batchSize: Int

  template: NotificationTemplate
  templateVars: JSON

  enabled: Boolean!
  validFrom: DateTime!
  validUntil: DateTime
}

enum Priority {
  CRITICAL
  HIGH
  NORMAL
  LOW
}

enum Channel {
  EMAIL
  SMS
  PUSH
  SLACK
  WEBHOOK
  IN_APP
  BROWSER_PUSH
}

enum NotificationStatus {
  PENDING
  QUEUED
  SENDING
  SENT
  DELIVERED
  FAILED
  RETRYING
  CANCELLED
}

type Query {
  notification(id: ID!): Notification
  notifications(
    filters: NotificationFilters
    pagination: PaginationInput
  ): NotificationConnection!

  notificationTemplate(id: ID!): NotificationTemplate
  notificationTemplates(
    filters: TemplateFilters
    pagination: PaginationInput
  ): TemplateConnection!

  notificationRule(id: ID!): NotificationRule
  notificationRules(
    filters: RuleFilters
    pagination: PaginationInput
  ): RuleConnection!

  userNotificationPreferences(userId: ID!): UserPreferences!
  notificationStats(filters: StatsFilters!): NotificationStats!
}

type Mutation {
  sendNotification(input: SendNotificationInput!): Notification!
  cancelNotification(id: ID!): Boolean!
  markNotificationRead(id: ID!): Notification!

  createNotificationTemplate(input: TemplateInput!): NotificationTemplate!
  updateNotificationTemplate(id: ID!, input: TemplateInput!): NotificationTemplate!
  deleteNotificationTemplate(id: ID!): Boolean!

  createNotificationRule(input: RuleInput!): NotificationRule!
  updateNotificationRule(id: ID!, input: RuleInput!): NotificationRule!
  deleteNotificationRule(id: ID!): Boolean!

  updateUserPreferences(userId: ID!, input: PreferencesInput!): UserPreferences!
}

type Subscription {
  notificationReceived(userId: ID!): Notification!
  notificationStatusChanged(notificationId: ID!): Notification!
}
```

## Event Integration

### Kafka Topics

```
# Notification requests (input)
system.notifications.requested
├─ Producer: Any service/adapter
├─ Consumer: Notifications Service
├─ Schema: NotificationRequestEvent (Avro)
└─ Partitioning: By recipient_id

# Notification status updates (output)
system.notifications.sent
├─ Producer: Notifications Service
├─ Consumer: Analytics, Audit services
├─ Schema: NotificationSentEvent (Avro)
└─ Partitioning: By notification_id

system.notifications.delivered
├─ Producer: Notifications Service
├─ Consumer: Analytics, Audit services
├─ Schema: NotificationDeliveredEvent (Avro)
└─ Partitioning: By notification_id

system.notifications.failed
├─ Producer: Notifications Service
├─ Consumer: Monitoring, Alert services
├─ Schema: NotificationFailedEvent (Avro)
└─ Partitioning: By notification_id
```

### Event Schemas (Avro)

```avro
{
  "namespace": "com.dictamesh.notifications.events",
  "type": "record",
  "name": "NotificationRequestEvent",
  "fields": [
    {"name": "event_id", "type": "string"},
    {"name": "event_type", "type": "string"},
    {"name": "timestamp", "type": "long", "logicalType": "timestamp-millis"},

    {"name": "recipient_type", "type": {
      "type": "enum",
      "name": "RecipientType",
      "symbols": ["USER", "ROLE", "GROUP", "SYSTEM"]
    }},
    {"name": "recipient_id", "type": "string"},

    {"name": "priority", "type": {
      "type": "enum",
      "name": "Priority",
      "symbols": ["CRITICAL", "HIGH", "NORMAL", "LOW"]
    }},

    {"name": "channels", "type": {"type": "array", "items": "string"}},

    {"name": "template_id", "type": ["null", "string"], "default": null},
    {"name": "template_vars", "type": ["null", {"type": "map", "values": "string"}], "default": null},

    {"name": "subject", "type": ["null", "string"], "default": null},
    {"name": "body", "type": ["null", "string"], "default": null},

    {"name": "metadata", "type": ["null", {"type": "map", "values": "string"}], "default": null},

    {"name": "trace_context", "type": {
      "type": "record",
      "name": "TraceContext",
      "fields": [
        {"name": "trace_id", "type": "string"},
        {"name": "span_id", "type": "string"}
      ]
    }}
  ]
}
```

## Security & Compliance

### Authentication & Authorization

- API authentication via JWT tokens
- Role-based access control (RBAC)
- Service-to-service auth via mTLS
- Channel provider credentials stored in secrets manager

### Data Privacy

- PII handling in compliance with GDPR/CCPA
- User consent tracking for marketing notifications
- Data retention policies
- Right to be forgotten support
- Encryption at rest and in transit

### Audit & Compliance

- Complete audit trail of all notifications
- Delivery receipts and tracking
- User consent logs
- Admin action logs
- Compliance reports (GDPR, HIPAA, SOC2)

### Rate Limiting & Abuse Prevention

- Per-user rate limits
- System-wide rate limits
- IP-based rate limiting for API
- Anomaly detection for abuse
- Circuit breakers for channel providers

## Deployment & Scaling

### Deployment Architecture

```
┌────────────────────────────────────────────────────────────┐
│                      KUBERNETES CLUSTER                     │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌──────────────────────────────────────────────────┐     │
│  │  Notifications Service (StatefulSet)             │     │
│  │  - Replicas: 3+                                  │     │
│  │  - Anti-affinity rules                           │     │
│  │  - HPA: CPU/Memory/Queue depth                   │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  ┌──────────────────────────────────────────────────┐     │
│  │  Channel Workers (Deployment per channel)        │     │
│  │  - Email worker: 5 replicas                      │     │
│  │  - SMS worker: 3 replicas                        │     │
│  │  - Push worker: 5 replicas                       │     │
│  │  - Slack worker: 2 replicas                      │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  ┌──────────────────────────────────────────────────┐     │
│  │  PostgreSQL (StatefulSet or External RDS)        │     │
│  │  - Primary + Replicas                            │     │
│  │  - Auto-failover                                 │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  ┌──────────────────────────────────────────────────┐     │
│  │  Redis (StatefulSet or ElastiCache)              │     │
│  │  - Cluster mode or Sentinel                      │     │
│  │  - Persistence enabled                           │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  ┌──────────────────────────────────────────────────┐     │
│  │  Kafka/Redpanda (StatefulSet)                    │     │
│  │  - 3+ brokers                                    │     │
│  │  - Replication factor: 3                         │     │
│  └──────────────────────────────────────────────────┘     │
└────────────────────────────────────────────────────────────┘
```

### Scaling Considerations

**Horizontal Scaling:**
- Notification processors: Scale based on queue depth
- Channel workers: Scale based on channel-specific queues
- API servers: Scale based on request rate

**Vertical Scaling:**
- Database: Increase resources for high throughput
- Redis: Increase memory for rate limiting data
- Kafka: Increase storage for event retention

### Performance Targets

```
Throughput:
- 10,000 notifications/second (sustained)
- 50,000 notifications/second (burst)

Latency:
- P50: < 100ms (queue to start processing)
- P95: < 500ms
- P99: < 1s

Delivery:
- Email: 99.5% delivered within 5 minutes
- SMS: 99.9% delivered within 1 minute
- Push: 99.9% delivered within 30 seconds
- In-App: 99.9% delivered immediately

Availability:
- Service SLA: 99.95%
- No single point of failure
- Graceful degradation
```

### Monitoring & Observability

**Metrics (Prometheus):**
- Notification rate (per channel, per priority)
- Delivery success rate
- Delivery latency (per channel)
- Queue depth (per channel)
- Rate limit hits
- Error rates
- Channel provider API health

**Traces (Jaeger):**
- End-to-end notification lifecycle
- Template rendering time
- Channel provider latency
- Database query performance

**Logs (Structured):**
- All notification events
- Delivery attempts
- Errors and retries
- Rule matching decisions
- Rate limit violations

**Dashboards:**
- Real-time notification throughput
- Channel health and performance
- Delivery success rates
- Error breakdown by channel
- User preference analytics
- Cost analysis (per channel)

### Disaster Recovery

**Backup Strategy:**
- Database: Continuous backup with PITR
- Templates: Version controlled in Git
- Configuration: Stored in ConfigMaps/Secrets with backups

**Failure Scenarios:**
- Channel provider outage: Automatic fallback to alternative channel
- Database failure: Automatic failover to replica
- Redis failure: Graceful degradation (skip rate limiting)
- Kafka failure: Buffer in memory, resume when recovered

**Recovery Time Objectives:**
- RTO: < 5 minutes for critical channels
- RPO: < 1 minute (event replay from Kafka)

## Implementation Phases

### Phase 1: Core Infrastructure (4 weeks)
- Database schema and migrations
- Basic notification models
- Template engine implementation
- Rule engine implementation
- REST API skeleton
- Unit tests

### Phase 2: Channel Providers (4 weeks)
- Email channel (SMTP, AWS SES)
- SMS channel (Twilio)
- Slack channel (Webhook)
- Push channel (FCM basic)
- Delivery tracking
- Integration tests

### Phase 3: Advanced Features (3 weeks)
- Rate limiting
- Batching and aggregation
- User preferences
- Retry logic with exponential backoff
- Circuit breakers
- Performance tests

### Phase 4: Framework Integration (2 weeks)
- Kafka event integration
- System alert rules (database, circuit breaker, etc.)
- GraphQL API
- WebSocket for in-app notifications
- End-to-end tests

### Phase 5: Production Readiness (2 weeks)
- Monitoring and observability
- Kubernetes manifests
- Documentation
- Performance optimization
- Load testing
- Security audit

### Phase 6: Extended Channels (3 weeks)
- Browser push notifications
- PagerDuty integration
- Webhook channel
- APNs for iOS
- Additional providers (SendGrid, etc.)

**Total Estimated Timeline: 18 weeks**

## Conclusion

The DictaMesh Notifications Service provides a comprehensive, production-ready solution for both framework infrastructure monitoring and extensible application-level notifications. The architecture is designed for:

- **Reliability**: Retry logic, fallbacks, circuit breakers
- **Scalability**: Horizontal scaling, queue-based processing
- **Flexibility**: Multiple channels, extensible provider system
- **Observability**: Complete tracking and monitoring
- **Developer Experience**: Simple API, template system, rule engine

The service integrates seamlessly with the DictaMesh framework's event-driven architecture while remaining extensible for application-specific needs.

## References

- DictaMesh Framework Architecture: `PROJECT-SCOPE.md`
- Database Infrastructure: `pkg/database/README.md`
- Event Bus Integration: Kafka/Redpanda documentation
- OpenTelemetry: https://opentelemetry.io/
- Prometheus: https://prometheus.io/
