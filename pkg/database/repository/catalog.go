// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package repository provides repository pattern implementations for database access
package repository

import (
	"context"
	"fmt"

	"github.com/click2-run/dictamesh/pkg/database/models"
	"gorm.io/gorm"
)

// CatalogRepository provides access to entity catalog
type CatalogRepository struct {
	db *gorm.DB
}

// NewCatalogRepository creates a new catalog repository
func NewCatalogRepository(db *gorm.DB) *CatalogRepository {
	return &CatalogRepository{db: db}
}

// Create creates a new entity catalog entry
func (r *CatalogRepository) Create(ctx context.Context, entity *models.EntityCatalog) error {
	result := r.db.WithContext(ctx).Create(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to create entity: %w", result.Error)
	}
	return nil
}

// FindByID finds an entity by ID
func (r *CatalogRepository) FindByID(ctx context.Context, id string) (*models.EntityCatalog, error) {
	var entity models.EntityCatalog
	result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find entity: %w", result.Error)
	}
	return &entity, nil
}

// FindBySource finds an entity by source system and ID
func (r *CatalogRepository) FindBySource(ctx context.Context, sourceSystem, sourceEntityID, entityType string) (*models.EntityCatalog, error) {
	var entity models.EntityCatalog
	result := r.db.WithContext(ctx).Where(
		"source_system = ? AND source_entity_id = ? AND entity_type = ?",
		sourceSystem, sourceEntityID, entityType,
	).First(&entity)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find entity: %w", result.Error)
	}
	return &entity, nil
}

// List lists entities with pagination
func (r *CatalogRepository) List(ctx context.Context, filters *CatalogFilters) ([]models.EntityCatalog, error) {
	query := r.db.WithContext(ctx)

	if filters.EntityType != "" {
		query = query.Where("entity_type = ?", filters.EntityType)
	}
	if filters.Domain != "" {
		query = query.Where("domain = ?", filters.Domain)
	}
	if filters.SourceSystem != "" {
		query = query.Where("source_system = ?", filters.SourceSystem)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.ContainsPII != nil {
		query = query.Where("contains_pii = ?", *filters.ContainsPII)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var entities []models.EntityCatalog
	result := query.Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list entities: %w", result.Error)
	}

	return entities, nil
}

// Update updates an entity
func (r *CatalogRepository) Update(ctx context.Context, entity *models.EntityCatalog) error {
	result := r.db.WithContext(ctx).Save(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update entity: %w", result.Error)
	}
	return nil
}

// Delete deletes an entity
func (r *CatalogRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.EntityCatalog{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete entity: %w", result.Error)
	}
	return nil
}

// CatalogFilters represents filters for listing entities
type CatalogFilters struct {
	EntityType   string
	Domain       string
	SourceSystem string
	Status       string
	ContainsPII  *bool
	Limit        int
	Offset       int
}

// RelationshipRepository provides access to entity relationships
type RelationshipRepository struct {
	db *gorm.DB
}

// NewRelationshipRepository creates a new relationship repository
func NewRelationshipRepository(db *gorm.DB) *RelationshipRepository {
	return &RelationshipRepository{db: db}
}

// Create creates a new relationship
func (r *RelationshipRepository) Create(ctx context.Context, rel *models.EntityRelationship) error {
	result := r.db.WithContext(ctx).Create(rel)
	if result.Error != nil {
		return fmt.Errorf("failed to create relationship: %w", result.Error)
	}
	return nil
}

// FindBySubject finds relationships where entity is the subject
func (r *RelationshipRepository) FindBySubject(ctx context.Context, entityType, entityID string) ([]models.EntityRelationship, error) {
	var relationships []models.EntityRelationship
	result := r.db.WithContext(ctx).
		Where("subject_entity_type = ? AND subject_entity_id = ? AND valid_to IS NULL", entityType, entityID).
		Preload("ObjectCatalog").
		Find(&relationships)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find relationships: %w", result.Error)
	}

	return relationships, nil
}

// FindByObject finds relationships where entity is the object
func (r *RelationshipRepository) FindByObject(ctx context.Context, entityType, entityID string) ([]models.EntityRelationship, error) {
	var relationships []models.EntityRelationship
	result := r.db.WithContext(ctx).
		Where("object_entity_type = ? AND object_entity_id = ? AND valid_to IS NULL", entityType, entityID).
		Preload("SubjectCatalog").
		Find(&relationships)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find relationships: %w", result.Error)
	}

	return relationships, nil
}
