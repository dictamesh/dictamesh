// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

// EmbeddingModel represents an embedding model configuration
type EmbeddingModel struct {
	Name       string
	Version    string
	Dimensions int
}

// EntityEmbedding represents a vector embedding of an entity
type EntityEmbedding struct {
	ID                 string
	CatalogID          string
	EmbeddingModel     string
	EmbeddingVersion   string
	EmbeddingDimensions int
	Embedding          pgvector.Vector
	SourceText         string
	SourceFields       map[string]interface{}
	Metadata           map[string]interface{}
}

// DocumentChunk represents a chunked document for RAG
type DocumentChunk struct {
	ID               string
	CatalogID        string
	ChunkIndex       int
	ChunkText        string
	ChunkTokens      int
	EmbeddingModel   string
	Embedding        pgvector.Vector
	PrecedingContext string
	FollowingContext string
	Metadata         map[string]interface{}
}

// SimilarEntity represents a search result with similarity score
type SimilarEntity struct {
	CatalogID  string
	Similarity float64
	SourceText string
	Metadata   map[string]interface{}
}

// RelevantChunk represents a relevant document chunk for RAG
type RelevantChunk struct {
	ChunkID          string
	CatalogID        string
	ChunkText        string
	ChunkIndex       int
	PrecedingContext string
	FollowingContext string
	Similarity       float64
	Metadata         map[string]interface{}
}

// VectorSearch provides vector similarity search capabilities
type VectorSearch struct {
	db *Database
}

// NewVectorSearch creates a new vector search instance
func NewVectorSearch(db *Database) *VectorSearch {
	return &VectorSearch{db: db}
}

// StoreEmbedding stores an entity embedding
func (vs *VectorSearch) StoreEmbedding(ctx context.Context, embedding *EntityEmbedding) error {
	query := `
		INSERT INTO dictamesh_entity_embeddings (
			catalog_id, embedding_model, embedding_version, embedding_dimensions,
			embedding, source_text, source_fields, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (catalog_id, embedding_model, embedding_version)
		DO UPDATE SET
			embedding = EXCLUDED.embedding,
			source_text = EXCLUDED.source_text,
			source_fields = EXCLUDED.source_fields,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
		RETURNING id
	`

	err := vs.db.pool.QueryRow(ctx, query,
		embedding.CatalogID,
		embedding.EmbeddingModel,
		embedding.EmbeddingVersion,
		embedding.EmbeddingDimensions,
		embedding.Embedding,
		embedding.SourceText,
		embedding.SourceFields,
		embedding.Metadata,
	).Scan(&embedding.ID)

	if err != nil {
		return fmt.Errorf("failed to store embedding: %w", err)
	}

	return nil
}

// StoreDocumentChunk stores a document chunk with embedding
func (vs *VectorSearch) StoreDocumentChunk(ctx context.Context, chunk *DocumentChunk) error {
	query := `
		INSERT INTO dictamesh_document_chunks (
			catalog_id, chunk_index, chunk_text, chunk_tokens,
			embedding_model, embedding, preceding_context, following_context, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (catalog_id, chunk_index, embedding_model)
		DO UPDATE SET
			chunk_text = EXCLUDED.chunk_text,
			chunk_tokens = EXCLUDED.chunk_tokens,
			embedding = EXCLUDED.embedding,
			preceding_context = EXCLUDED.preceding_context,
			following_context = EXCLUDED.following_context,
			metadata = EXCLUDED.metadata
		RETURNING id
	`

	err := vs.db.pool.QueryRow(ctx, query,
		chunk.CatalogID,
		chunk.ChunkIndex,
		chunk.ChunkText,
		chunk.ChunkTokens,
		chunk.EmbeddingModel,
		chunk.Embedding,
		chunk.PrecedingContext,
		chunk.FollowingContext,
		chunk.Metadata,
	).Scan(&chunk.ID)

	if err != nil {
		return fmt.Errorf("failed to store document chunk: %w", err)
	}

	return nil
}

// FindSimilarEntities finds entities similar to the query embedding
func (vs *VectorSearch) FindSimilarEntities(
	ctx context.Context,
	queryEmbedding pgvector.Vector,
	modelName string,
	similarityThreshold float64,
	limit int,
) ([]SimilarEntity, error) {
	query := `
		SELECT catalog_id, similarity, source_text, metadata
		FROM dictamesh_find_similar_entities($1, $2, $3, $4)
	`

	rows, err := vs.db.pool.Query(ctx, query,
		queryEmbedding,
		modelName,
		similarityThreshold,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar entities: %w", err)
	}
	defer rows.Close()

	var results []SimilarEntity
	for rows.Next() {
		var entity SimilarEntity
		if err := rows.Scan(
			&entity.CatalogID,
			&entity.Similarity,
			&entity.SourceText,
			&entity.Metadata,
		); err != nil {
			return nil, fmt.Errorf("failed to scan similar entity: %w", err)
		}
		results = append(results, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating similar entities: %w", err)
	}

	return results, nil
}

// FindRelevantChunks finds relevant document chunks for RAG
func (vs *VectorSearch) FindRelevantChunks(
	ctx context.Context,
	queryEmbedding pgvector.Vector,
	modelName string,
	catalogID *string,
	similarityThreshold float64,
	limit int,
) ([]RelevantChunk, error) {
	query := `
		SELECT chunk_id, catalog_id, chunk_text, chunk_index,
		       preceding_context, following_context, similarity, metadata
		FROM dictamesh_find_relevant_chunks($1, $2, $3, $4, $5)
	`

	rows, err := vs.db.pool.Query(ctx, query,
		queryEmbedding,
		modelName,
		catalogID,
		similarityThreshold,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find relevant chunks: %w", err)
	}
	defer rows.Close()

	var results []RelevantChunk
	for rows.Next() {
		var chunk RelevantChunk
		if err := rows.Scan(
			&chunk.ChunkID,
			&chunk.CatalogID,
			&chunk.ChunkText,
			&chunk.ChunkIndex,
			&chunk.PrecedingContext,
			&chunk.FollowingContext,
			&chunk.Similarity,
			&chunk.Metadata,
		); err != nil {
			return nil, fmt.Errorf("failed to scan relevant chunk: %w", err)
		}
		results = append(results, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating relevant chunks: %w", err)
	}

	return results, nil
}

// HybridSearchResult represents a result from hybrid search
type HybridSearchResult struct {
	CatalogID        string
	CombinedScore    float64
	TextRank         float64
	VectorSimilarity float64
	SourceText       string
}

// HybridSearch performs combined full-text and vector search
func (vs *VectorSearch) HybridSearch(
	ctx context.Context,
	queryText string,
	queryEmbedding pgvector.Vector,
	modelName string,
	textWeight float64,
	vectorWeight float64,
	limit int,
) ([]HybridSearchResult, error) {
	query := `
		SELECT catalog_id, combined_score, text_rank, vector_similarity, source_text
		FROM dictamesh_hybrid_search($1, $2, $3, $4, $5, $6)
	`

	rows, err := vs.db.pool.Query(ctx, query,
		queryText,
		queryEmbedding,
		modelName,
		textWeight,
		vectorWeight,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform hybrid search: %w", err)
	}
	defer rows.Close()

	var results []HybridSearchResult
	for rows.Next() {
		var result HybridSearchResult
		if err := rows.Scan(
			&result.CatalogID,
			&result.CombinedScore,
			&result.TextRank,
			&result.VectorSimilarity,
			&result.SourceText,
		); err != nil {
			return nil, fmt.Errorf("failed to scan hybrid search result: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating hybrid search results: %w", err)
	}

	return results, nil
}

// DeleteEmbeddings deletes all embeddings for a catalog entry
func (vs *VectorSearch) DeleteEmbeddings(ctx context.Context, catalogID string) error {
	query := `DELETE FROM dictamesh_entity_embeddings WHERE catalog_id = $1`
	_, err := vs.db.pool.Exec(ctx, query, catalogID)
	if err != nil {
		return fmt.Errorf("failed to delete embeddings: %w", err)
	}
	return nil
}

// DeleteDocumentChunks deletes all chunks for a catalog entry
func (vs *VectorSearch) DeleteDocumentChunks(ctx context.Context, catalogID string) error {
	query := `DELETE FROM dictamesh_document_chunks WHERE catalog_id = $1`
	_, err := vs.db.pool.Exec(ctx, query, catalogID)
	if err != nil {
		return fmt.Errorf("failed to delete document chunks: %w", err)
	}
	return nil
}

// BatchStoreChunks stores multiple document chunks in a transaction
func (vs *VectorSearch) BatchStoreChunks(ctx context.Context, chunks []DocumentChunk) error {
	return vs.db.WithPgxTransaction(ctx, func(tx pgx.Tx) error {
		for i := range chunks {
			query := `
				INSERT INTO dictamesh_document_chunks (
					catalog_id, chunk_index, chunk_text, chunk_tokens,
					embedding_model, embedding, preceding_context, following_context, metadata
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT (catalog_id, chunk_index, embedding_model)
				DO UPDATE SET
					chunk_text = EXCLUDED.chunk_text,
					chunk_tokens = EXCLUDED.chunk_tokens,
					embedding = EXCLUDED.embedding,
					preceding_context = EXCLUDED.preceding_context,
					following_context = EXCLUDED.following_context,
					metadata = EXCLUDED.metadata
			`

			_, err := tx.Exec(ctx, query,
				chunks[i].CatalogID,
				chunks[i].ChunkIndex,
				chunks[i].ChunkText,
				chunks[i].ChunkTokens,
				chunks[i].EmbeddingModel,
				chunks[i].Embedding,
				chunks[i].PrecedingContext,
				chunks[i].FollowingContext,
				chunks[i].Metadata,
			)

			if err != nil {
				return fmt.Errorf("failed to store chunk %d: %w", i, err)
			}
		}
		return nil
	})
}
