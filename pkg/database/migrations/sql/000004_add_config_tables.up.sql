-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

--
-- Migration: Add Configuration Management Tables
-- Description: Creates tables for centralized configuration management with versioning,
--              secrets encryption, and audit logging
-- IMPORTANT: All table names MUST use the dictamesh_ prefix for namespace isolation
--

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- Table: dictamesh_configurations
-- Purpose: Store all framework and service configurations
-- ============================================================================
CREATE TABLE IF NOT EXISTS dictamesh_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Configuration hierarchy
    environment VARCHAR(50) NOT NULL CHECK (environment IN ('dev', 'development', 'staging', 'production', 'test')),
    service VARCHAR(100) NOT NULL,        -- e.g., 'metadata-catalog', 'graphql-gateway', 'event-router'
    component VARCHAR(100),               -- e.g., 'database', 'cache', 'notifications', NULL for service-level
    key VARCHAR(255) NOT NULL,            -- configuration key (e.g., 'max_connections', 'timeout')

    -- Configuration value and metadata
    value JSONB NOT NULL,                 -- configuration value (can be any JSON type)
    value_type VARCHAR(50) NOT NULL CHECK (value_type IN ('string', 'number', 'boolean', 'object', 'array')),
    is_secret BOOLEAN DEFAULT false,      -- true if value is encrypted

    -- Validation and documentation
    schema JSONB,                         -- JSON schema for validation
    description TEXT,
    tags TEXT[],                          -- searchable tags for categorization

    -- Versioning
    version INTEGER NOT NULL DEFAULT 1,

    -- Audit fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),              -- user or service that created
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255),              -- user or service that last updated

    -- Constraints
    CONSTRAINT unique_config_key UNIQUE(environment, service, component, key)
);

-- Indexes for efficient querying
CREATE INDEX idx_dictamesh_configs_env_svc ON dictamesh_configurations(environment, service);
CREATE INDEX idx_dictamesh_configs_svc ON dictamesh_configurations(service);
CREATE INDEX idx_dictamesh_configs_env ON dictamesh_configurations(environment);
CREATE INDEX idx_dictamesh_configs_tags ON dictamesh_configurations USING GIN(tags);
CREATE INDEX idx_dictamesh_configs_is_secret ON dictamesh_configurations(is_secret) WHERE is_secret = true;
CREATE INDEX idx_dictamesh_configs_updated_at ON dictamesh_configurations(updated_at DESC);

-- Table comment
COMMENT ON TABLE dictamesh_configurations IS 'DictaMesh: Centralized configuration storage with versioning and encryption support';

-- Column comments
COMMENT ON COLUMN dictamesh_configurations.environment IS 'Deployment environment (dev, staging, production)';
COMMENT ON COLUMN dictamesh_configurations.service IS 'Service name this configuration belongs to';
COMMENT ON COLUMN dictamesh_configurations.component IS 'Component within service (e.g., database, cache) or NULL for service-level config';
COMMENT ON COLUMN dictamesh_configurations.key IS 'Configuration key identifier';
COMMENT ON COLUMN dictamesh_configurations.value IS 'Configuration value (encrypted if is_secret=true)';
COMMENT ON COLUMN dictamesh_configurations.is_secret IS 'Whether this value is encrypted (for sensitive data)';
COMMENT ON COLUMN dictamesh_configurations.schema IS 'JSON Schema for validating configuration values';
COMMENT ON COLUMN dictamesh_configurations.tags IS 'Tags for categorization and search';

-- ============================================================================
-- Table: dictamesh_config_versions
-- Purpose: Track complete history of configuration changes
-- ============================================================================
CREATE TABLE IF NOT EXISTS dictamesh_config_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_id UUID NOT NULL REFERENCES dictamesh_configurations(id) ON DELETE CASCADE,

    -- Version information
    version INTEGER NOT NULL,
    value JSONB NOT NULL,                 -- value at this version
    is_secret BOOLEAN DEFAULT false,      -- whether value was encrypted

    -- Change metadata
    change_description TEXT,              -- description of what changed and why

    -- Audit fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,

    -- Constraints
    CONSTRAINT unique_config_version UNIQUE(config_id, version)
);

-- Indexes
CREATE INDEX idx_dictamesh_config_versions_config ON dictamesh_config_versions(config_id);
CREATE INDEX idx_dictamesh_config_versions_created_at ON dictamesh_config_versions(created_at DESC);

-- Table comment
COMMENT ON TABLE dictamesh_config_versions IS 'DictaMesh: Version history for all configuration changes';
COMMENT ON COLUMN dictamesh_config_versions.version IS 'Version number (increments with each change)';
COMMENT ON COLUMN dictamesh_config_versions.change_description IS 'Human-readable description of the change';

-- ============================================================================
-- Table: dictamesh_config_audit_logs
-- Purpose: Comprehensive audit trail for all configuration operations
-- ============================================================================
CREATE TABLE IF NOT EXISTS dictamesh_config_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_id UUID REFERENCES dictamesh_configurations(id) ON DELETE SET NULL,

    -- Action information
    action VARCHAR(50) NOT NULL CHECK (action IN ('CREATE', 'UPDATE', 'DELETE', 'ACCESS', 'ROLLBACK', 'EXPORT', 'IMPORT')),

    -- Actor information
    actor VARCHAR(255) NOT NULL,          -- user email or service account
    actor_type VARCHAR(50) CHECK (actor_type IN ('USER', 'SERVICE', 'API_KEY', 'SYSTEM')),

    -- Request metadata
    ip_address INET,
    user_agent TEXT,
    request_id VARCHAR(100),              -- correlation ID for distributed tracing

    -- Change details
    changes JSONB,                        -- before/after for updates, or full content for create/delete
    metadata JSONB,                       -- additional context (e.g., rollback target version)

    -- Timestamp
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for audit queries
CREATE INDEX idx_dictamesh_config_audit_timestamp ON dictamesh_config_audit_logs(timestamp DESC);
CREATE INDEX idx_dictamesh_config_audit_actor ON dictamesh_config_audit_logs(actor);
CREATE INDEX idx_dictamesh_config_audit_config ON dictamesh_config_audit_logs(config_id);
CREATE INDEX idx_dictamesh_config_audit_action ON dictamesh_config_audit_logs(action);
CREATE INDEX idx_dictamesh_config_audit_request ON dictamesh_config_audit_logs(request_id);

-- Table comment
COMMENT ON TABLE dictamesh_config_audit_logs IS 'DictaMesh: Audit trail for all configuration management operations';
COMMENT ON COLUMN dictamesh_config_audit_logs.action IS 'Type of operation performed';
COMMENT ON COLUMN dictamesh_config_audit_logs.actor IS 'User or service that performed the action';
COMMENT ON COLUMN dictamesh_config_audit_logs.changes IS 'Details of changes made (before/after values)';

-- ============================================================================
-- Table: dictamesh_encryption_keys
-- Purpose: Manage encryption keys for secrets (envelope encryption pattern)
-- ============================================================================
CREATE TABLE IF NOT EXISTS dictamesh_encryption_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Key identification
    key_name VARCHAR(100) UNIQUE NOT NULL,    -- e.g., 'master-key-2025-01', 'prod-dek-v1'
    key_type VARCHAR(50) NOT NULL CHECK (key_type IN ('MASTER', 'DATA_ENCRYPTION_KEY')),
    environment VARCHAR(50),                  -- NULL for master keys, environment for DEKs

    -- Key material (encrypted with KEK from environment)
    encrypted_key BYTEA NOT NULL,
    algorithm VARCHAR(50) NOT NULL DEFAULT 'AES-256-GCM',
    key_version INTEGER NOT NULL DEFAULT 1,

    -- Key lifecycle
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMPTZ,
    rotated_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,

    -- Metadata
    metadata JSONB,                           -- additional key metadata

    -- Constraints
    CHECK (key_type = 'MASTER' OR environment IS NOT NULL)
);

-- Indexes
CREATE INDEX idx_dictamesh_encryption_keys_active ON dictamesh_encryption_keys(is_active, key_type);
CREATE INDEX idx_dictamesh_encryption_keys_env ON dictamesh_encryption_keys(environment) WHERE environment IS NOT NULL;

-- Table comment
COMMENT ON TABLE dictamesh_encryption_keys IS 'DictaMesh: Encryption key management for secrets (envelope encryption)';
COMMENT ON COLUMN dictamesh_encryption_keys.key_type IS 'Type of key: MASTER (encrypts DEKs) or DATA_ENCRYPTION_KEY (encrypts data)';
COMMENT ON COLUMN dictamesh_encryption_keys.encrypted_key IS 'Key material encrypted with KEK from environment variable';

-- ============================================================================
-- Table: dictamesh_config_watchers
-- Purpose: Track active configuration watchers for hot reload
-- ============================================================================
CREATE TABLE IF NOT EXISTS dictamesh_config_watchers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Watcher identification
    service_instance VARCHAR(255) NOT NULL,   -- service instance ID
    service_name VARCHAR(100) NOT NULL,

    -- Watch specification
    environment VARCHAR(50) NOT NULL,
    watch_pattern VARCHAR(255) NOT NULL,      -- e.g., 'database.*', 'cache.redis.*'

    -- Watcher status
    last_heartbeat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,

    -- Registration
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB
);

-- Indexes
CREATE INDEX idx_dictamesh_config_watchers_service ON dictamesh_config_watchers(service_name, environment);
CREATE INDEX idx_dictamesh_config_watchers_heartbeat ON dictamesh_config_watchers(last_heartbeat DESC) WHERE is_active = true;

-- Table comment
COMMENT ON TABLE dictamesh_config_watchers IS 'DictaMesh: Active configuration watchers for hot reload functionality';

-- ============================================================================
-- Trigger: Auto-increment version on configuration update
-- ============================================================================
CREATE OR REPLACE FUNCTION dictamesh_increment_config_version()
RETURNS TRIGGER AS $$
BEGIN
    -- Only increment version if value actually changed
    IF OLD.value IS DISTINCT FROM NEW.value THEN
        NEW.version := OLD.version + 1;
        NEW.updated_at := NOW();

        -- Create version history record
        INSERT INTO dictamesh_config_versions (
            config_id, version, value, is_secret, created_by
        ) VALUES (
            NEW.id, NEW.version, NEW.value, NEW.is_secret, NEW.updated_by
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_dictamesh_config_version
    BEFORE UPDATE ON dictamesh_configurations
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_increment_config_version();

COMMENT ON FUNCTION dictamesh_increment_config_version() IS 'DictaMesh: Auto-increment version and create history record on config update';

-- ============================================================================
-- Trigger: Create initial version on configuration insert
-- ============================================================================
CREATE OR REPLACE FUNCTION dictamesh_create_initial_config_version()
RETURNS TRIGGER AS $$
BEGIN
    -- Create initial version record
    INSERT INTO dictamesh_config_versions (
        config_id, version, value, is_secret, created_by
    ) VALUES (
        NEW.id, 1, NEW.value, NEW.is_secret, NEW.created_by
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_dictamesh_initial_config_version
    AFTER INSERT ON dictamesh_configurations
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_create_initial_config_version();

COMMENT ON FUNCTION dictamesh_create_initial_config_version() IS 'DictaMesh: Create initial version record when configuration is created';

-- ============================================================================
-- Function: Get configuration value with type coercion
-- ============================================================================
CREATE OR REPLACE FUNCTION dictamesh_get_config(
    p_environment VARCHAR,
    p_service VARCHAR,
    p_component VARCHAR,
    p_key VARCHAR
) RETURNS JSONB AS $$
DECLARE
    v_value JSONB;
BEGIN
    SELECT value INTO v_value
    FROM dictamesh_configurations
    WHERE environment = p_environment
      AND service = p_service
      AND (component = p_component OR (component IS NULL AND p_component IS NULL))
      AND key = p_key;

    RETURN v_value;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION dictamesh_get_config IS 'DictaMesh: Retrieve configuration value by environment, service, component, and key';

-- ============================================================================
-- Function: Cleanup old audit logs (for maintenance)
-- ============================================================================
CREATE OR REPLACE FUNCTION dictamesh_cleanup_old_audit_logs(
    p_retention_days INTEGER DEFAULT 90
) RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER;
BEGIN
    DELETE FROM dictamesh_config_audit_logs
    WHERE timestamp < NOW() - (p_retention_days || ' days')::INTERVAL;

    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;

    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION dictamesh_cleanup_old_audit_logs IS 'DictaMesh: Delete audit logs older than specified days (default 90)';

-- ============================================================================
-- Function: Cleanup stale watchers (heartbeat timeout)
-- ============================================================================
CREATE OR REPLACE FUNCTION dictamesh_cleanup_stale_watchers(
    p_timeout_minutes INTEGER DEFAULT 5
) RETURNS INTEGER AS $$
DECLARE
    v_updated_count INTEGER;
BEGIN
    UPDATE dictamesh_config_watchers
    SET is_active = false
    WHERE is_active = true
      AND last_heartbeat < NOW() - (p_timeout_minutes || ' minutes')::INTERVAL;

    GET DIAGNOSTICS v_updated_count = ROW_COUNT;

    RETURN v_updated_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION dictamesh_cleanup_stale_watchers IS 'DictaMesh: Mark watchers as inactive if no heartbeat within timeout';

-- ============================================================================
-- View: Active configurations by environment
-- ============================================================================
CREATE OR REPLACE VIEW dictamesh_active_configs AS
SELECT
    environment,
    service,
    component,
    COUNT(*) as config_count,
    SUM(CASE WHEN is_secret THEN 1 ELSE 0 END) as secret_count,
    MAX(updated_at) as last_updated
FROM dictamesh_configurations
GROUP BY environment, service, component
ORDER BY environment, service, component;

COMMENT ON VIEW dictamesh_active_configs IS 'DictaMesh: Summary of active configurations by environment and service';

-- ============================================================================
-- Sample Data for Development
-- ============================================================================
-- Insert sample configuration for development environment
INSERT INTO dictamesh_configurations (
    environment, service, component, key, value, value_type, is_secret, description, tags, created_by
) VALUES
(
    'development',
    'admin-console',
    'api',
    'port',
    '8081'::jsonb,
    'number',
    false,
    'API server port',
    ARRAY['api', 'network'],
    'system'
),
(
    'development',
    'admin-console',
    'ui',
    'dev_server_port',
    '5173'::jsonb,
    'number',
    false,
    'Remix dev server port',
    ARRAY['ui', 'development'],
    'system'
),
(
    'development',
    'metadata-catalog',
    'database',
    'max_connections',
    '100'::jsonb,
    'number',
    false,
    'Maximum number of database connections',
    ARRAY['database', 'performance'],
    'system'
),
(
    'development',
    'graphql-gateway',
    NULL,
    'introspection_enabled',
    'true'::jsonb,
    'boolean',
    false,
    'Enable GraphQL introspection in development',
    ARRAY['graphql', 'development'],
    'system'
);

-- Grant appropriate permissions
-- Note: Adjust these based on your actual database roles
-- GRANT SELECT, INSERT, UPDATE, DELETE ON dictamesh_configurations TO dictamesh_app;
-- GRANT SELECT ON dictamesh_config_versions TO dictamesh_app;
-- GRANT SELECT, INSERT ON dictamesh_config_audit_logs TO dictamesh_app;
