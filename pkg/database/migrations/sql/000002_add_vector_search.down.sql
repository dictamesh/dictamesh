-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- Drop functions
DROP FUNCTION IF EXISTS hybrid_search;
DROP FUNCTION IF EXISTS find_relevant_chunks;
DROP FUNCTION IF EXISTS find_similar_entities;
DROP FUNCTION IF EXISTS update_embedding_search_vector;

-- Drop tables
DROP TABLE IF EXISTS document_chunks CASCADE;
DROP TABLE IF EXISTS entity_embeddings CASCADE;

-- Note: We don't drop the vector extension as it might be used elsewhere
