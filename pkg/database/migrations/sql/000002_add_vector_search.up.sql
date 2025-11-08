-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- Add pgvector extension for vector similarity search and RAG support
-- IMPORTANT: All DictaMesh tables use the 'dictamesh_' prefix

CREATE EXTENSION IF NOT EXISTS vector;

-- Entity Embeddings: Store vector embeddings for semantic search
CREATE TABLE IF NOT EXISTS dictamesh_entity_embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    catalog_id UUID NOT NULL REFERENCES dictamesh_entity_catalog(id) ON DELETE CASCADE,

    -- Embedding metadata
    embedding_model VARCHAR(100) NOT NULL, -- e.g., 'text-embedding-ada-002', 'sentence-transformers'
    embedding_version VARCHAR(50) NOT NULL,
    embedding_dimensions INTEGER NOT NULL,

    -- Vector embedding (flexible dimensions, typically 384, 768, or 1536)
    embedding vector(1536),

    -- Source text used for embedding
    source_text TEXT NOT NULL,
    source_fields JSONB, -- Which fields were used to generate embedding

    -- Metadata for context
    metadata JSONB,

    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    -- Ensure one embedding per model version per entity
    UNIQUE(catalog_id, embedding_model, embedding_version)
);

-- Indexes for efficient vector search
-- HNSW index for fast approximate nearest neighbor search
CREATE INDEX idx_dictamesh_embedding_hnsw ON dictamesh_entity_embeddings
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

-- IVFFlat index as an alternative (faster build, slower query)
-- CREATE INDEX idx_dictamesh_embedding_ivfflat ON dictamesh_entity_embeddings
--     USING ivfflat (embedding vector_cosine_ops)
--     WITH (lists = 100);

-- Additional indexes
CREATE INDEX idx_dictamesh_embedding_catalog ON dictamesh_entity_embeddings(catalog_id);
CREATE INDEX idx_dictamesh_embedding_model ON dictamesh_entity_embeddings(embedding_model, embedding_version);
CREATE INDEX idx_dictamesh_embedding_metadata ON dictamesh_entity_embeddings USING gin(metadata);

-- Document Chunks: For RAG - store document chunks with embeddings
CREATE TABLE IF NOT EXISTS dictamesh_document_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    catalog_id UUID NOT NULL REFERENCES dictamesh_entity_catalog(id) ON DELETE CASCADE,

    -- Chunk metadata
    chunk_index INTEGER NOT NULL, -- Position in document
    chunk_text TEXT NOT NULL,
    chunk_tokens INTEGER, -- Token count for LLM context management

    -- Embedding
    embedding_model VARCHAR(100) NOT NULL,
    embedding vector(1536),

    -- Context for retrieval
    preceding_context TEXT, -- Text before chunk for context
    following_context TEXT, -- Text after chunk for context

    -- Metadata
    metadata JSONB, -- Page number, section, etc.

    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(catalog_id, chunk_index, embedding_model)
);

-- Indexes for document chunks
CREATE INDEX idx_dictamesh_chunk_embedding_hnsw ON dictamesh_document_chunks
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

CREATE INDEX idx_dictamesh_chunk_catalog ON dictamesh_document_chunks(catalog_id);
CREATE INDEX idx_dictamesh_chunk_metadata ON dictamesh_document_chunks USING gin(metadata);

-- Full-text search integration
ALTER TABLE dictamesh_entity_embeddings ADD COLUMN search_vector tsvector;

CREATE INDEX idx_dictamesh_embedding_search_vector ON dictamesh_entity_embeddings USING gin(search_vector);

-- Trigger to maintain search vector
CREATE OR REPLACE FUNCTION dictamesh_update_embedding_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', COALESCE(NEW.source_text, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_dictamesh_embedding_search_vector_trigger
    BEFORE INSERT OR UPDATE ON dictamesh_entity_embeddings
    FOR EACH ROW
    EXECUTE FUNCTION dictamesh_update_embedding_search_vector();

-- Semantic Search Functions

-- Function to find similar entities by vector similarity
CREATE OR REPLACE FUNCTION dictamesh_find_similar_entities(
    query_embedding vector(1536),
    model_name VARCHAR(100),
    similarity_threshold FLOAT DEFAULT 0.7,
    result_limit INTEGER DEFAULT 10
)
RETURNS TABLE (
    catalog_id UUID,
    similarity FLOAT,
    source_text TEXT,
    metadata JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        ee.catalog_id,
        1 - (ee.embedding <=> query_embedding) AS similarity,
        ee.source_text,
        ee.metadata
    FROM dictamesh_entity_embeddings ee
    WHERE ee.embedding_model = model_name
        AND (1 - (ee.embedding <=> query_embedding)) >= similarity_threshold
    ORDER BY ee.embedding <=> query_embedding
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql;

-- Function to find relevant document chunks for RAG
CREATE OR REPLACE FUNCTION dictamesh_find_relevant_chunks(
    query_embedding vector(1536),
    model_name VARCHAR(100),
    entity_filter UUID DEFAULT NULL,
    similarity_threshold FLOAT DEFAULT 0.7,
    result_limit INTEGER DEFAULT 5
)
RETURNS TABLE (
    chunk_id UUID,
    catalog_id UUID,
    chunk_text TEXT,
    chunk_index INTEGER,
    preceding_context TEXT,
    following_context TEXT,
    similarity FLOAT,
    metadata JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        dc.id,
        dc.catalog_id,
        dc.chunk_text,
        dc.chunk_index,
        dc.preceding_context,
        dc.following_context,
        1 - (dc.embedding <=> query_embedding) AS similarity,
        dc.metadata
    FROM dictamesh_document_chunks dc
    WHERE dc.embedding_model = model_name
        AND (entity_filter IS NULL OR dc.catalog_id = entity_filter)
        AND (1 - (dc.embedding <=> query_embedding)) >= similarity_threshold
    ORDER BY dc.embedding <=> query_embedding
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql;

-- Hybrid search: Combine vector similarity with full-text search
CREATE OR REPLACE FUNCTION dictamesh_hybrid_search(
    query_text TEXT,
    query_embedding vector(1536),
    model_name VARCHAR(100),
    text_weight FLOAT DEFAULT 0.5,
    vector_weight FLOAT DEFAULT 0.5,
    result_limit INTEGER DEFAULT 10
)
RETURNS TABLE (
    catalog_id UUID,
    combined_score FLOAT,
    text_rank FLOAT,
    vector_similarity FLOAT,
    source_text TEXT
) AS $$
BEGIN
    RETURN QUERY
    WITH text_scores AS (
        SELECT
            ee.catalog_id,
            ts_rank(ee.search_vector, plainto_tsquery('english', query_text)) AS rank
        FROM dictamesh_entity_embeddings ee
        WHERE ee.search_vector @@ plainto_tsquery('english', query_text)
    ),
    vector_scores AS (
        SELECT
            ee.catalog_id,
            1 - (ee.embedding <=> query_embedding) AS similarity,
            ee.source_text
        FROM dictamesh_entity_embeddings ee
        WHERE ee.embedding_model = model_name
    )
    SELECT
        COALESCE(ts.catalog_id, vs.catalog_id) AS catalog_id,
        (COALESCE(ts.rank, 0) * text_weight + COALESCE(vs.similarity, 0) * vector_weight) AS combined_score,
        COALESCE(ts.rank, 0) AS text_rank,
        COALESCE(vs.similarity, 0) AS vector_similarity,
        vs.source_text
    FROM text_scores ts
    FULL OUTER JOIN vector_scores vs ON ts.catalog_id = vs.catalog_id
    ORDER BY combined_score DESC
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql;

-- Comments
COMMENT ON TABLE dictamesh_entity_embeddings IS 'DictaMesh: Vector embeddings of entities for semantic search and similarity analysis';
COMMENT ON TABLE dictamesh_document_chunks IS 'DictaMesh: Chunked documents with embeddings for RAG (Retrieval-Augmented Generation)';
COMMENT ON FUNCTION dictamesh_find_similar_entities IS 'DictaMesh: Find entities similar to query embedding using cosine similarity';
COMMENT ON FUNCTION dictamesh_find_relevant_chunks IS 'DictaMesh: Find relevant document chunks for RAG based on vector similarity';
COMMENT ON FUNCTION dictamesh_hybrid_search IS 'DictaMesh: Combine full-text and vector search for improved relevance';
