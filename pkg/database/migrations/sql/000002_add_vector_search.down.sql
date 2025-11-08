-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- Drop functions
DROP FUNCTION IF EXISTS dictamesh_hybrid_search;
DROP FUNCTION IF EXISTS dictamesh_find_relevant_chunks;
DROP FUNCTION IF EXISTS dictamesh_find_similar_entities;
DROP FUNCTION IF EXISTS dictamesh_update_embedding_search_vector;

-- Drop tables
DROP TABLE IF EXISTS dictamesh_document_chunks CASCADE;
DROP TABLE IF EXISTS dictamesh_entity_embeddings CASCADE;

-- Note: We don't drop the vector extension as it might be used elsewhere
