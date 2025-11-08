-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- DictaMesh Metadata Catalog Schema - Initial Migration
-- This migration sets up the core database schema for the metadata catalog
--
-- IMPORTANT: All DictaMesh tables use the 'dictamesh_' prefix to avoid naming conflicts
-- and clearly identify framework tables in shared database environments.

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search optimization
CREATE EXTENSION IF NOT EXISTS "btree_gin"; -- For composite indexes
CREATE EXTENSION IF NOT EXISTS "btree_gist"; -- For temporal queries

-- Entity Registry: Catalog of all entities across integrated systems
CREATE TABLE IF NOT EXISTS dictamesh_entity_catalog (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(100) NOT NULL,
    domain VARCHAR(100) NOT NULL,
    source_system VARCHAR(100) NOT NULL,
    source_entity_id VARCHAR(255) NOT NULL,

    -- API access information
    api_base_url TEXT NOT NULL,
    api_path_template TEXT NOT NULL,
    api_method VARCHAR(10) DEFAULT 'GET',
    api_auth_type VARCHAR(50),

    -- Schema reference
    schema_id UUID,
    schema_version VARCHAR(50),

    -- Lifecycle metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(50) DEFAULT 'active',

    -- SLA information
    availability_sla DECIMAL(5,4),
    latency_p99_ms INTEGER,
    freshness_sla_seconds INTEGER,

    -- Data classification
    contains_pii BOOLEAN DEFAULT FALSE,
    data_classification VARCHAR(50),
    retention_days INTEGER,

    UNIQUE(source_system, source_entity_id, entity_type)
);

-- Indexes for dictamesh_entity_catalog
CREATE INDEX idx_dictamesh_entity_type ON dictamesh_entity_catalog(entity_type);
CREATE INDEX idx_dictamesh_domain ON dictamesh_entity_catalog(domain);
CREATE INDEX idx_dictamesh_source_system ON dictamesh_entity_catalog(source_system);
CREATE INDEX idx_dictamesh_status ON dictamesh_entity_catalog(status);
CREATE INDEX idx_dictamesh_contains_pii ON dictamesh_entity_catalog(contains_pii) WHERE contains_pii = true;
CREATE INDEX idx_dictamesh_data_classification ON dictamesh_entity_catalog(data_classification);

-- Entity Relationships: Cross-system relationship graph
CREATE TABLE IF NOT EXISTS dictamesh_entity_relationships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Subject (from)
    subject_catalog_id UUID REFERENCES dictamesh_entity_catalog(id) ON DELETE CASCADE,
    subject_entity_type VARCHAR(100) NOT NULL,
    subject_entity_id VARCHAR(255) NOT NULL,

    -- Predicate (relationship type)
    relationship_type VARCHAR(100) NOT NULL,
    relationship_cardinality VARCHAR(20),

    -- Object (to)
    object_catalog_id UUID REFERENCES dictamesh_entity_catalog(id) ON DELETE CASCADE,
    object_entity_type VARCHAR(100) NOT NULL,
    object_entity_id VARCHAR(255) NOT NULL,

    -- Temporal validity
    valid_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to TIMESTAMPTZ,

    -- Denormalized display fields
    subject_display_name VARCHAR(255),
    object_display_name VARCHAR(255),

    -- Relationship metadata
    relationship_metadata JSONB,

    -- Lineage
    created_by_event_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT temporal_validity CHECK (valid_to IS NULL OR valid_to > valid_from)
);

-- Indexes for dictamesh_entity_relationships
CREATE INDEX idx_dictamesh_subject ON dictamesh_entity_relationships(subject_entity_type, subject_entity_id);
CREATE INDEX idx_dictamesh_object ON dictamesh_entity_relationships(object_entity_type, object_entity_id);
CREATE INDEX idx_dictamesh_relationship_type ON dictamesh_entity_relationships(relationship_type);
CREATE INDEX idx_dictamesh_temporal ON dictamesh_entity_relationships(valid_from, valid_to) WHERE valid_to IS NULL;
CREATE INDEX idx_dictamesh_relationship_metadata ON dictamesh_entity_relationships USING gin(relationship_metadata);

-- Schema Registry: Versioned schemas for all entities
CREATE TABLE IF NOT EXISTS dictamesh_schemas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(100) NOT NULL,
    version VARCHAR(50) NOT NULL,
    schema_format VARCHAR(50) NOT NULL,
    schema_definition JSONB NOT NULL,

    -- Compatibility
    backward_compatible BOOLEAN DEFAULT TRUE,
    forward_compatible BOOLEAN DEFAULT FALSE,

    -- Lifecycle
    published_at TIMESTAMPTZ DEFAULT NOW(),
    deprecated_at TIMESTAMPTZ,
    retired_at TIMESTAMPTZ,

    UNIQUE(entity_type, version)
);

-- Indexes for dictamesh_schemas
CREATE INDEX idx_dictamesh_schema_entity_type ON dictamesh_schemas(entity_type);
CREATE INDEX idx_dictamesh_schema_version ON dictamesh_schemas(version);
CREATE INDEX idx_dictamesh_schema_definition ON dictamesh_schemas USING gin(schema_definition);

-- Event Log: Immutable audit trail
CREATE TABLE IF NOT EXISTS dictamesh_event_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id VARCHAR(255) UNIQUE NOT NULL,
    event_type VARCHAR(100) NOT NULL,

    catalog_id UUID REFERENCES dictamesh_entity_catalog(id) ON DELETE SET NULL,
    entity_type VARCHAR(100),
    entity_id VARCHAR(255),

    changed_fields TEXT[],
    event_payload JSONB,

    -- Tracing
    trace_id VARCHAR(64),
    span_id VARCHAR(16),

    -- Time
    event_timestamp TIMESTAMPTZ NOT NULL,
    ingested_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for dictamesh_event_log
CREATE INDEX idx_dictamesh_event_catalog ON dictamesh_event_log(catalog_id, event_timestamp DESC);
CREATE INDEX idx_dictamesh_event_type ON dictamesh_event_log(entity_type, entity_id, event_timestamp DESC);
CREATE INDEX idx_dictamesh_trace ON dictamesh_event_log(trace_id);
CREATE INDEX idx_dictamesh_event_timestamp ON dictamesh_event_log(event_timestamp DESC);
CREATE INDEX idx_dictamesh_event_payload ON dictamesh_event_log USING gin(event_payload);

-- Data Lineage: Track data flow and transformations
CREATE TABLE IF NOT EXISTS dictamesh_data_lineage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Upstream (source)
    upstream_catalog_id UUID REFERENCES dictamesh_entity_catalog(id) ON DELETE CASCADE,
    upstream_system VARCHAR(100),

    -- Downstream (derived)
    downstream_catalog_id UUID REFERENCES dictamesh_entity_catalog(id) ON DELETE CASCADE,
    downstream_system VARCHAR(100),

    -- Transformation metadata
    transformation_type VARCHAR(50),
    transformation_logic TEXT,

    -- Observability
    data_flow_active BOOLEAN DEFAULT TRUE,
    last_flow_at TIMESTAMPTZ,
    average_latency_ms INTEGER,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for dictamesh_data_lineage
CREATE INDEX idx_dictamesh_lineage_upstream ON dictamesh_data_lineage(upstream_catalog_id);
CREATE INDEX idx_dictamesh_lineage_downstream ON dictamesh_data_lineage(downstream_catalog_id);
CREATE INDEX idx_dictamesh_lineage_active ON dictamesh_data_lineage(data_flow_active) WHERE data_flow_active = true;

-- Cache Status: Track cache freshness
CREATE TABLE IF NOT EXISTS dictamesh_cache_status (
    entity_catalog_id UUID REFERENCES dictamesh_entity_catalog(id) ON DELETE CASCADE,
    entity_id VARCHAR(255) NOT NULL,
    cache_layer VARCHAR(50) NOT NULL,

    cached_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    cache_key VARCHAR(500),
    hit_count INTEGER DEFAULT 0,

    PRIMARY KEY (entity_catalog_id, entity_id, cache_layer)
);

-- Indexes for dictamesh_cache_status
CREATE INDEX idx_dictamesh_cache_expiry ON dictamesh_cache_status(expires_at);
CREATE INDEX idx_dictamesh_cache_layer ON dictamesh_cache_status(cache_layer);

-- Triggers and Functions
CREATE OR REPLACE FUNCTION dictamesh_update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to dictamesh_entity_catalog
CREATE TRIGGER update_dictamesh_entity_catalog_updated_at
    BEFORE UPDATE ON dictamesh_entity_catalog
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE dictamesh_entity_catalog IS 'DictaMesh: Registry of all entities across integrated data sources';
COMMENT ON TABLE dictamesh_entity_relationships IS 'DictaMesh: Cross-system entity relationship graph with temporal validity';
COMMENT ON TABLE dictamesh_schemas IS 'DictaMesh: Versioned entity schema registry with compatibility tracking';
COMMENT ON TABLE dictamesh_event_log IS 'DictaMesh: Immutable audit trail of all entity events with distributed tracing';
COMMENT ON TABLE dictamesh_data_lineage IS 'DictaMesh: Data flow and transformation tracking for lineage analysis';
COMMENT ON TABLE dictamesh_cache_status IS 'DictaMesh: Cache freshness and hit rate tracking for performance optimization';
