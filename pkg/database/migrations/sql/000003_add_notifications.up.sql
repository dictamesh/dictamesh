-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- DictaMesh Notifications Service Database Schema
-- This migration creates all tables needed for the notifications service
-- All tables use the dictamesh_ prefix for namespace isolation

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- Notification Templates
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,

    -- Content (JSONB for flexibility)
    channels JSONB NOT NULL,
    translations JSONB,

    -- Template metadata
    variables JSONB,
    schema_version VARCHAR(50) DEFAULT '1.0',

    -- Lifecycle
    version VARCHAR(50) DEFAULT '1.0.0',
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by VARCHAR(255),

    -- Organization
    tags TEXT[],

    CONSTRAINT valid_channels CHECK (jsonb_typeof(channels) = 'object')
);

CREATE INDEX idx_dictamesh_notification_templates_name
    ON dictamesh_notification_templates(name);
CREATE INDEX idx_dictamesh_notification_templates_tags
    ON dictamesh_notification_templates USING GIN(tags);
CREATE INDEX idx_dictamesh_notification_templates_enabled
    ON dictamesh_notification_templates(enabled) WHERE enabled = TRUE;

COMMENT ON TABLE dictamesh_notification_templates IS
    'DictaMesh: Notification template definitions with multi-channel support';

-- ============================================================================
-- Notification Rules
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notification_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,

    -- Trigger conditions
    event_pattern TEXT NOT NULL,
    domains TEXT[],
    event_types TEXT[],

    -- Routing
    priority VARCHAR(20) NOT NULL,
    channels TEXT[] NOT NULL,
    fallback_channels TEXT[],

    -- Recipients
    recipient_selector JSONB NOT NULL,

    -- Timing
    schedule JSONB,
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
CREATE INDEX idx_dictamesh_notification_rules_template
    ON dictamesh_notification_rules(template_id);

COMMENT ON TABLE dictamesh_notification_rules IS
    'DictaMesh: Rules for triggering notifications based on events';

-- ============================================================================
-- Notifications (main tracking table)
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Source
    event_id VARCHAR(255),
    rule_id UUID REFERENCES dictamesh_notification_rules(id),
    template_id UUID REFERENCES dictamesh_notification_templates(id),

    -- Recipients
    recipient_type VARCHAR(50) NOT NULL,
    recipient_id VARCHAR(255) NOT NULL,

    -- Content
    subject TEXT,
    body TEXT,
    body_html TEXT,
    data JSONB,

    -- Routing
    priority VARCHAR(20) NOT NULL,
    channels TEXT[] NOT NULL,
    selected_channel VARCHAR(50),

    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    -- Timing
    scheduled_at TIMESTAMPTZ DEFAULT NOW(),
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    read_at TIMESTAMPTZ,

    -- Metadata
    metadata JSONB,
    trace_id VARCHAR(64),

    -- Error tracking
    error TEXT,
    retry_count INTEGER DEFAULT 0,
    next_retry_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_status CHECK (status IN (
        'PENDING', 'QUEUED', 'SENDING', 'SENT',
        'DELIVERED', 'FAILED', 'RETRYING', 'CANCELLED'
    )),
    CONSTRAINT valid_recipient_type CHECK (recipient_type IN ('USER', 'ROLE', 'GROUP', 'SYSTEM'))
) PARTITION BY RANGE (created_at);

-- Create initial partitions (12 months)
CREATE TABLE dictamesh_notifications_y2025m01 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE dictamesh_notifications_y2025m02 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
CREATE TABLE dictamesh_notifications_y2025m03 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');
CREATE TABLE dictamesh_notifications_y2025m04 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');
CREATE TABLE dictamesh_notifications_y2025m05 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');
CREATE TABLE dictamesh_notifications_y2025m06 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');
CREATE TABLE dictamesh_notifications_y2025m07 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');
CREATE TABLE dictamesh_notifications_y2025m08 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');
CREATE TABLE dictamesh_notifications_y2025m09 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');
CREATE TABLE dictamesh_notifications_y2025m10 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
CREATE TABLE dictamesh_notifications_y2025m11 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
CREATE TABLE dictamesh_notifications_y2025m12 PARTITION OF dictamesh_notifications
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

-- Indexes on partitioned table
CREATE INDEX idx_dictamesh_notifications_recipient
    ON dictamesh_notifications(recipient_type, recipient_id, created_at DESC);
CREATE INDEX idx_dictamesh_notifications_status
    ON dictamesh_notifications(status, scheduled_at);
CREATE INDEX idx_dictamesh_notifications_event
    ON dictamesh_notifications(event_id);
CREATE INDEX idx_dictamesh_notifications_trace
    ON dictamesh_notifications(trace_id);
CREATE INDEX idx_dictamesh_notifications_rule
    ON dictamesh_notifications(rule_id);
CREATE INDEX idx_dictamesh_notifications_template
    ON dictamesh_notifications(template_id);

COMMENT ON TABLE dictamesh_notifications IS
    'DictaMesh: Notification instances and delivery tracking';

-- ============================================================================
-- Delivery Attempts (detailed delivery tracking)
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notification_delivery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL,

    -- Delivery details
    channel VARCHAR(50) NOT NULL,
    provider VARCHAR(100),

    -- Status
    status VARCHAR(20) NOT NULL,
    attempt_number INTEGER NOT NULL,

    -- Timing
    started_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    -- Result
    success BOOLEAN DEFAULT FALSE,
    error TEXT,
    provider_response JSONB,
    provider_message_id VARCHAR(255),

    -- Metadata
    metadata JSONB,

    CONSTRAINT valid_delivery_status CHECK (status IN (
        'SENDING', 'SENT', 'DELIVERED', 'FAILED', 'BOUNCED', 'REJECTED'
    ))
);

CREATE INDEX idx_dictamesh_notification_delivery_notification
    ON dictamesh_notification_delivery(notification_id, attempt_number DESC);
CREATE INDEX idx_dictamesh_notification_delivery_provider_id
    ON dictamesh_notification_delivery(provider_message_id);
CREATE INDEX idx_dictamesh_notification_delivery_status
    ON dictamesh_notification_delivery(status, started_at DESC);

COMMENT ON TABLE dictamesh_notification_delivery IS
    'DictaMesh: Detailed delivery attempts and provider responses';

-- ============================================================================
-- User Preferences
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notification_preferences (
    user_id VARCHAR(255) PRIMARY KEY,

    -- Global settings
    enabled BOOLEAN DEFAULT TRUE,
    timezone VARCHAR(50) DEFAULT 'UTC',
    locale VARCHAR(10) DEFAULT 'en',

    -- Channel addresses
    email VARCHAR(255),
    phone VARCHAR(20),
    push_tokens JSONB,

    -- Channel preferences
    channel_prefs JSONB DEFAULT '{}',

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

COMMENT ON TABLE dictamesh_notification_preferences IS
    'DictaMesh: User notification preferences and settings';

-- ============================================================================
-- Notification Batches (for grouped notifications)
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notification_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Batch config
    rule_id UUID REFERENCES dictamesh_notification_rules(id),
    batch_key VARCHAR(255) NOT NULL,

    -- Timing
    window_start TIMESTAMPTZ NOT NULL,
    window_end TIMESTAMPTZ NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    sent_at TIMESTAMPTZ,

    -- Content
    notification_ids UUID[],
    count INTEGER NOT NULL,

    -- Status
    status VARCHAR(20) DEFAULT 'PENDING',

    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT valid_batch_status CHECK (status IN ('PENDING', 'SENT', 'FAILED'))
);

CREATE INDEX idx_dictamesh_notification_batches_key_window
    ON dictamesh_notification_batches(batch_key, window_end);
CREATE INDEX idx_dictamesh_notification_batches_scheduled
    ON dictamesh_notification_batches(status, scheduled_at);
CREATE INDEX idx_dictamesh_notification_batches_rule
    ON dictamesh_notification_batches(rule_id);

COMMENT ON TABLE dictamesh_notification_batches IS
    'DictaMesh: Batched notifications for efficient delivery';

-- ============================================================================
-- Rate Limiting Configuration
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notification_rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    scope VARCHAR(50) NOT NULL,
    scope_id VARCHAR(255),
    channel VARCHAR(50) NOT NULL,

    -- Limit definition
    max_count INTEGER NOT NULL,
    window_seconds INTEGER NOT NULL,

    -- Metadata
    enabled BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(scope, COALESCE(scope_id, ''), channel)
);

CREATE INDEX idx_dictamesh_notification_rate_limits_scope
    ON dictamesh_notification_rate_limits(scope, scope_id, channel);

COMMENT ON TABLE dictamesh_notification_rate_limits IS
    'DictaMesh: Rate limiting configuration for notifications';

-- ============================================================================
-- Audit Log (comprehensive tracking)
-- ============================================================================

CREATE TABLE IF NOT EXISTS dictamesh_notification_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    notification_id UUID,

    -- Event details
    event_type VARCHAR(100) NOT NULL,
    actor_type VARCHAR(50),
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
CREATE INDEX idx_dictamesh_notification_audit_timestamp
    ON dictamesh_notification_audit(timestamp DESC);

COMMENT ON TABLE dictamesh_notification_audit IS
    'DictaMesh: Comprehensive audit trail for all notification events';

-- ============================================================================
-- Functions and Triggers
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION dictamesh_update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at trigger to relevant tables
CREATE TRIGGER update_dictamesh_notification_templates_updated_at
    BEFORE UPDATE ON dictamesh_notification_templates
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_update_updated_at_column();

CREATE TRIGGER update_dictamesh_notification_rules_updated_at
    BEFORE UPDATE ON dictamesh_notification_rules
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_update_updated_at_column();

CREATE TRIGGER update_dictamesh_notifications_updated_at
    BEFORE UPDATE ON dictamesh_notifications
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_update_updated_at_column();

CREATE TRIGGER update_dictamesh_notification_preferences_updated_at
    BEFORE UPDATE ON dictamesh_notification_preferences
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_update_updated_at_column();

CREATE TRIGGER update_dictamesh_notification_rate_limits_updated_at
    BEFORE UPDATE ON dictamesh_notification_rate_limits
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_update_updated_at_column();

-- ============================================================================
-- Sample Data (for development/testing)
-- ============================================================================

-- Insert default rate limits
INSERT INTO dictamesh_notification_rate_limits (scope, scope_id, channel, max_count, window_seconds)
VALUES
    ('USER', NULL, 'EMAIL', 100, 3600),
    ('USER', NULL, 'SMS', 10, 3600),
    ('USER', NULL, 'PUSH', 50, 3600),
    ('USER', NULL, 'SLACK', 30, 3600),
    ('SYSTEM', NULL, 'EMAIL', 10000, 3600),
    ('SYSTEM', NULL, 'SMS', 1000, 3600),
    ('SYSTEM', NULL, 'PUSH', 50000, 3600),
    ('SYSTEM', NULL, 'SLACK', 100, 60)
ON CONFLICT DO NOTHING;

-- Insert sample infrastructure alert template
INSERT INTO dictamesh_notification_templates (name, description, channels, variables, tags)
VALUES (
    'infrastructure-alert',
    'Template for infrastructure monitoring alerts',
    '{
        "EMAIL": {
            "subject": "[{{.Priority}}] {{.Title}}",
            "body": "{{.Title}}\n\nPriority: {{.Priority}}\nService: {{.Service}}\nMessage: {{.Message}}\nTime: {{.Timestamp}}\n\nDashboard: {{.DashboardURL}}",
            "body_html": "<h1>{{.Title}}</h1><p><strong>Priority:</strong> {{.Priority}}</p><p><strong>Service:</strong> {{.Service}}</p><p><strong>Message:</strong> {{.Message}}</p><p><strong>Time:</strong> {{.Timestamp}}</p><p><a href=\"{{.DashboardURL}}\">View Dashboard</a></p>"
        },
        "SLACK": {
            "body": "{\"text\": \"[{{.Priority}}] {{.Title}}\", \"blocks\": [{\"type\": \"header\", \"text\": {\"type\": \"plain_text\", \"text\": \"{{.Title}}\"}}, {\"type\": \"section\", \"fields\": [{\"type\": \"mrkdwn\", \"text\": \"*Priority:*\\n{{.Priority}}\"}, {\"type\": \"mrkdwn\", \"text\": \"*Service:*\\n{{.Service}}\"}]}, {\"type\": \"section\", \"text\": {\"type\": \"mrkdwn\", \"text\": \"{{.Message}}\"}}]}"
        }
    }',
    '["Title", "Priority", "Service", "Message", "Timestamp", "DashboardURL"]',
    ARRAY['infrastructure', 'alert', 'ops']
)
ON CONFLICT (name) DO NOTHING;

-- Insert sample rule for database alerts
INSERT INTO dictamesh_notification_rules (
    name,
    description,
    event_pattern,
    event_types,
    priority,
    channels,
    recipient_selector,
    template_id,
    enabled
)
SELECT
    'database-health-alert',
    'Alert for database health issues',
    'event.type == "system.database.health" && event.data.status == "unhealthy"',
    ARRAY['system.database.health'],
    'CRITICAL',
    ARRAY['EMAIL', 'SLACK'],
    '{"type": "role", "roles": ["ops-team", "on-call"]}',
    id,
    TRUE
FROM dictamesh_notification_templates
WHERE name = 'infrastructure-alert'
ON CONFLICT (name) DO NOTHING;
