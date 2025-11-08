-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

--
-- DictaMesh Metadata Catalog Schema
-- This initializes the PostgreSQL database for the framework's metadata catalog service
--

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search optimization

-- Entity Registry: Catalog of all entities across integrated systems
CREATE TABLE IF NOT EXISTS entity_catalog (
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

CREATE INDEX idx_entity_type ON entity_catalog(entity_type);
CREATE INDEX idx_domain ON entity_catalog(domain);
CREATE INDEX idx_source_system ON entity_catalog(source_system);
CREATE INDEX idx_status ON entity_catalog(status);

-- Entity Relationships: Cross-system relationship graph
CREATE TABLE IF NOT EXISTS entity_relationships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Subject (from)
    subject_catalog_id UUID REFERENCES entity_catalog(id) ON DELETE CASCADE,
    subject_entity_type VARCHAR(100) NOT NULL,
    subject_entity_id VARCHAR(255) NOT NULL,

    -- Predicate (relationship type)
    relationship_type VARCHAR(100) NOT NULL,
    relationship_cardinality VARCHAR(20),

    -- Object (to)
    object_catalog_id UUID REFERENCES entity_catalog(id) ON DELETE CASCADE,
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

CREATE INDEX idx_subject ON entity_relationships(subject_entity_type, subject_entity_id);
CREATE INDEX idx_object ON entity_relationships(object_entity_type, object_entity_id);
CREATE INDEX idx_relationship_type ON entity_relationships(relationship_type);
CREATE INDEX idx_temporal ON entity_relationships(valid_from, valid_to) WHERE valid_to IS NULL;

-- Schema Registry: Versioned schemas for all entities
CREATE TABLE IF NOT EXISTS schemas (
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

CREATE INDEX idx_schema_entity_type ON schemas(entity_type);
CREATE INDEX idx_schema_version ON schemas(version);

-- Event Log: Immutable audit trail
CREATE TABLE IF NOT EXISTS event_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id VARCHAR(255) UNIQUE NOT NULL,
    event_type VARCHAR(100) NOT NULL,

    catalog_id UUID REFERENCES entity_catalog(id) ON DELETE SET NULL,
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

CREATE INDEX idx_event_catalog ON event_log(catalog_id, event_timestamp DESC);
CREATE INDEX idx_event_type ON event_log(entity_type, entity_id, event_timestamp DESC);
CREATE INDEX idx_trace ON event_log(trace_id);
CREATE INDEX idx_event_timestamp ON event_log(event_timestamp DESC);

-- Data Lineage: Track data flow and transformations
CREATE TABLE IF NOT EXISTS data_lineage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Upstream (source)
    upstream_catalog_id UUID REFERENCES entity_catalog(id) ON DELETE CASCADE,
    upstream_system VARCHAR(100),

    -- Downstream (derived)
    downstream_catalog_id UUID REFERENCES entity_catalog(id) ON DELETE CASCADE,
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

CREATE INDEX idx_lineage_upstream ON data_lineage(upstream_catalog_id);
CREATE INDEX idx_lineage_downstream ON data_lineage(downstream_catalog_id);

-- Cache Status: Track cache freshness
CREATE TABLE IF NOT EXISTS cache_status (
    entity_catalog_id UUID REFERENCES entity_catalog(id) ON DELETE CASCADE,
    entity_id VARCHAR(255) NOT NULL,
    cache_layer VARCHAR(50) NOT NULL,

    cached_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    cache_key VARCHAR(500),
    hit_count INTEGER DEFAULT 0,

    PRIMARY KEY (entity_catalog_id, entity_id, cache_layer)
);

CREATE INDEX idx_cache_expiry ON cache_status(expires_at);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to entity_catalog
CREATE TRIGGER update_entity_catalog_updated_at BEFORE UPDATE ON entity_catalog
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO dictamesh;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO dictamesh;

-- Insert framework metadata
INSERT INTO entity_catalog (
    entity_type, domain, source_system, source_entity_id,
    api_base_url, api_path_template, status,
    data_classification
) VALUES (
    'framework_info', 'system', 'dictamesh', 'core',
    'http://localhost:8080', '/api/v1/info', 'active',
    'public'
) ON CONFLICT DO NOTHING;

COMMENT ON TABLE entity_catalog IS 'Registry of all entities across integrated data sources';
COMMENT ON TABLE entity_relationships IS 'Cross-system entity relationship graph';
COMMENT ON TABLE schemas IS 'Versioned entity schema registry';
COMMENT ON TABLE event_log IS 'Immutable audit trail of all entity events';
COMMENT ON TABLE data_lineage IS 'Data flow and transformation tracking';
COMMENT ON TABLE cache_status IS 'Cache freshness and hit rate tracking';
