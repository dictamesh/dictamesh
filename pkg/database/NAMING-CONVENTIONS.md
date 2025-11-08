# DictaMesh Database Naming Conventions

**CRITICAL REQUIREMENT**: All DictaMesh database objects MUST use the `dictamesh_` prefix to avoid naming conflicts and clearly identify framework tables in shared database environments.

## Table Names

### ✅ Required Pattern
All tables must follow this pattern:
```
dictamesh_{table_name}
```

### Examples
- `dictamesh_entity_catalog`
- `dictamesh_entity_relationships`
- `dictamesh_schemas`
- `dictamesh_event_log`
- `dictamesh_data_lineage`
- `dictamesh_cache_status`
- `dictamesh_entity_embeddings`
- `dictamesh_document_chunks`
- `dictamesh_audit_logs`

### ❌ NEVER Use
- `entity_catalog` (missing prefix)
- `dm_entity_catalog` (wrong prefix)
- `EntityCatalog` (wrong case and missing prefix)

## Index Names

### Required Pattern
```
idx_dictamesh_{meaningful_name}
```

### Examples
- `idx_dictamesh_entity_type`
- `idx_dictamesh_source_system`
- `idx_dictamesh_embedding_hnsw`
- `idx_dictamesh_audit_user`

## Function Names

### Required Pattern
```
dictamesh_{function_name}
```

### Examples
- `dictamesh_update_updated_at_column()`
- `dictamesh_find_similar_entities()`
- `dictamesh_find_relevant_chunks()`
- `dictamesh_hybrid_search()`
- `dictamesh_update_embedding_search_vector()`

## Trigger Names

### Required Pattern
```
{action}_dictamesh_{table_name}_{trigger_purpose}
```

### Examples
- `update_dictamesh_entity_catalog_updated_at`
- `update_dictamesh_embedding_search_vector_trigger`

## Constraint Names

### Required Pattern
- Primary Keys: `dictamesh_{table_name}_pkey`
- Foreign Keys: `fk_dictamesh_{table_name}_{referenced_table}`
- Unique Constraints: `uq_dictamesh_{table_name}_{column_names}`
- Check Constraints: `chk_dictamesh_{table_name}_{purpose}`

### Examples
- `fk_dictamesh_relationships_subject_catalog`
- `uq_dictamesh_entity_catalog_source`
- `chk_dictamesh_relationships_temporal_validity`

## GORM Model Configuration

All GORM models must override the default table name:

```go
// EntityCatalog represents an entity in the catalog
type EntityCatalog struct {
    ID         string    `gorm:"type:uuid;primary_key"`
    // ... fields
}

// TableName returns the table name with dictamesh_ prefix
func (EntityCatalog) TableName() string {
    return "dictamesh_entity_catalog"
}
```

## Migration Files

### File Naming
```
{number}_{description}.{up|down}.sql
```

### Content Requirements
1. **Header Comment**: Include SPDX license and copyright
2. **Prefix Reminder**: Add comment about dictamesh_ prefix requirement
3. **Table Names**: All tables must have dictamesh_ prefix
4. **Documentation**: Add COMMENT ON TABLE with "DictaMesh:" prefix

### Example
```sql
-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- IMPORTANT: All DictaMesh tables use the 'dictamesh_' prefix

CREATE TABLE IF NOT EXISTS dictamesh_entity_catalog (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- ... columns
);

COMMENT ON TABLE dictamesh_entity_catalog IS 'DictaMesh: Registry of all entities across integrated data sources';
```

## Query Code References

When writing queries in Go code, always use the prefixed table names:

```go
// ✅ CORRECT
query := `
    SELECT id, entity_type FROM dictamesh_entity_catalog
    WHERE status = $1
`

// ❌ WRONG
query := `
    SELECT id, entity_type FROM entity_catalog
    WHERE status = $1
`
```

## Repository Pattern

Repository methods should reference GORM models, which automatically use the correct table names:

```go
func (r *CatalogRepository) FindByID(ctx context.Context, id string) (*models.EntityCatalog, error) {
    var entity models.EntityCatalog
    // GORM automatically uses dictamesh_entity_catalog
    result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
    return &entity, result.Error
}
```

## Schema Comments

All table and function comments should start with "DictaMesh:" to clearly identify framework objects:

```sql
COMMENT ON TABLE dictamesh_entity_catalog IS 'DictaMesh: Registry of all entities across integrated data sources';
COMMENT ON FUNCTION dictamesh_find_similar_entities IS 'DictaMesh: Find entities similar to query embedding using cosine similarity';
```

## Rationale

### Why Use the dictamesh_ Prefix?

1. **Namespace Isolation**: Prevents conflicts with user tables in shared databases
2. **Clear Ownership**: Immediately identifies framework vs. user tables
3. **Multi-Tenancy Support**: Enables multiple frameworks in same database
4. **Migration Safety**: Reduces risk of accidentally modifying user tables
5. **Documentation**: Self-documenting code - prefix indicates framework responsibility

### Shared Database Scenarios

DictaMesh may be deployed in environments where:
- Multiple applications share the same PostgreSQL database
- Users have their own application tables alongside framework tables
- Database administrators need to quickly identify framework vs. application tables
- Backup and restore procedures need to target framework tables specifically

## Checklist for New Database Objects

When creating new database objects, verify:

- [ ] Table name has `dictamesh_` prefix
- [ ] All indexes have `idx_dictamesh_` prefix
- [ ] All functions have `dictamesh_` prefix
- [ ] All triggers follow naming convention
- [ ] GORM models override TableName()
- [ ] Table comments start with "DictaMesh:"
- [ ] Migration file includes prefix reminder comment
- [ ] All queries in Go code use prefixed names
- [ ] Documentation updated with new table names

## Common Mistakes to Avoid

### 1. Forgetting Prefix in Migrations
```sql
-- ❌ WRONG
CREATE TABLE entity_catalog (...)

-- ✅ CORRECT
CREATE TABLE dictamesh_entity_catalog (...)
```

### 2. Missing TableName() Override
```go
// ❌ WRONG - GORM will use default "entity_catalogs"
type EntityCatalog struct { ... }

// ✅ CORRECT
type EntityCatalog struct { ... }
func (EntityCatalog) TableName() string {
    return "dictamesh_entity_catalog"
}
```

### 3. Inconsistent Query References
```go
// ❌ WRONG - mixing prefixed and non-prefixed
query := `
    SELECT * FROM dictamesh_entity_catalog ec
    JOIN entity_relationships er ON ec.id = er.catalog_id
`

// ✅ CORRECT
query := `
    SELECT * FROM dictamesh_entity_catalog ec
    JOIN dictamesh_entity_relationships er ON ec.id = er.catalog_id
`
```

## Validation Tools

Use these queries to verify naming conventions:

### Check All Tables
```sql
SELECT tablename
FROM pg_tables
WHERE schemaname = 'public'
    AND tablename NOT LIKE 'dictamesh_%'
    AND tablename NOT IN ('schema_migrations');
```

### Check All Indexes
```sql
SELECT indexname
FROM pg_indexes
WHERE schemaname = 'public'
    AND indexname NOT LIKE 'idx_dictamesh_%'
    AND indexname NOT LIKE 'dictamesh_%_pkey';
```

### Check All Functions
```sql
SELECT routine_name
FROM information_schema.routines
WHERE routine_schema = 'public'
    AND routine_name NOT LIKE 'dictamesh_%';
```

## For LLM Agents

When implementing database features:

1. **ALWAYS** use `dictamesh_` prefix for all tables, indexes, and functions
2. **VERIFY** GORM models have TableName() method returning prefixed name
3. **UPDATE** all SQL queries to use prefixed table names
4. **ADD** comment in migration files reminding about prefix requirement
5. **DOCUMENT** new tables with "DictaMesh:" prefix in comments
6. **TEST** queries against actual database to verify table names

## References

- **Migrations**: `pkg/database/migrations/sql/`
- **Models**: `pkg/database/models/entity.go`
- **Repositories**: `pkg/database/repository/`
- **Query Examples**: `pkg/database/README.md`
