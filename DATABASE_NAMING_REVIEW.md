# Database Naming Convention Review

**Date:** 2025-11-08
**Branch:** `claude/review-database-connectors-011CUw4kLkqs8w4nzFTHbBvc`
**Reviewer:** Claude AI

## Executive Summary

A comprehensive review of all database schemas and connectors was conducted to ensure compliance with the DictaMesh naming convention requirement: **all database objects must use the `dictamesh_` prefix**.

### Status: ✅ RESOLVED

All database naming convention violations have been identified and corrected.

---

## Naming Convention Requirements

As specified in `AGENT.md`:

### Table Names
All database tables must follow the pattern:
```
dictamesh_{table_name}
```

### Other Database Objects
- **Indexes**: `idx_dictamesh_{meaningful_name}`
- **Functions**: `dictamesh_{function_name}()`
- **Triggers**: `{action}_dictamesh_{table_name}_{purpose}`

### GORM Models
All GORM models must override the `TableName()` method to return prefixed names.

---

## Review Findings

### ✅ Compliant Components

#### 1. Migration Files (pkg/database/migrations/sql/)
All migration files correctly use the `dictamesh_` prefix:

- **000001_initial_schema.up.sql** ✓
  - `dictamesh_entity_catalog`
  - `dictamesh_entity_relationships`
  - `dictamesh_schemas`
  - `dictamesh_event_log`
  - `dictamesh_data_lineage`
  - `dictamesh_cache_status`

- **000002_add_vector_search.up.sql** ✓
  - `dictamesh_entity_embeddings`
  - `dictamesh_document_chunks`
  - Functions: `dictamesh_find_similar_entities()`, `dictamesh_find_relevant_chunks()`, `dictamesh_hybrid_search()`

- **000003_add_notifications.up.sql** ✓
  - All 8 notification-related tables use `dictamesh_notification_*` prefix
  - All indexes use `idx_dictamesh_notification_*` prefix
  - Functions use `dictamesh_update_updated_at_column()`

- **000004_add_billing.up.sql** ✓
  - All 10 billing tables use `dictamesh_billing_*` prefix
  - All indexes use `idx_dictamesh_billing_*` prefix

#### 2. Go Models
All GORM models correctly override `TableName()`:

- **pkg/database/models/entity.go** ✓
  - `EntityCatalog` → `dictamesh_entity_catalog`
  - `EntityRelationship` → `dictamesh_entity_relationships`
  - `Schema` → `dictamesh_schemas`
  - `EventLog` → `dictamesh_event_log`
  - `DataLineage` → `dictamesh_data_lineage`
  - `CacheStatus` → `dictamesh_cache_status`

- **pkg/billing/models/models.go** ✓
  - All 10 billing models use `dictamesh_billing_*` tables

- **pkg/notifications/models/notification.go** ✓
  - All 8 notification models use `dictamesh_notification_*` tables

### ❌ Violation Found and Corrected

#### File: `infrastructure/docker-compose/init-scripts/postgres/01-init-metadata-catalog.sql`

**Original Issues:**
| Object Type | Incorrect Name | Correct Name |
|-------------|----------------|--------------|
| Table | `entity_catalog` | `dictamesh_entity_catalog` |
| Table | `entity_relationships` | `dictamesh_entity_relationships` |
| Table | `schemas` | `dictamesh_schemas` |
| Table | `event_log` | `dictamesh_event_log` |
| Table | `data_lineage` | `dictamesh_data_lineage` |
| Table | `cache_status` | `dictamesh_cache_status` |
| Index | `idx_entity_type` | `idx_dictamesh_entity_type` |
| Index | `idx_domain` | `idx_dictamesh_domain` |
| Index | `idx_source_system` | `idx_dictamesh_source_system` |
| Index | `idx_status` | `idx_dictamesh_status` |
| Index | `idx_subject` | `idx_dictamesh_subject` |
| Index | `idx_object` | `idx_dictamesh_object` |
| Index | `idx_relationship_type` | `idx_dictamesh_relationship_type` |
| Index | `idx_temporal` | `idx_dictamesh_temporal` |
| Index | `idx_schema_entity_type` | `idx_dictamesh_schema_entity_type` |
| Index | `idx_schema_version` | `idx_dictamesh_schema_version` |
| Index | `idx_event_catalog` | `idx_dictamesh_event_catalog` |
| Index | `idx_event_type` | `idx_dictamesh_event_type` |
| Index | `idx_trace` | `idx_dictamesh_trace` |
| Index | `idx_event_timestamp` | `idx_dictamesh_event_timestamp` |
| Index | `idx_lineage_upstream` | `idx_dictamesh_lineage_upstream` |
| Index | `idx_lineage_downstream` | `idx_dictamesh_lineage_downstream` |
| Index | `idx_cache_expiry` | `idx_dictamesh_cache_expiry` |
| Function | `update_updated_at_column()` | `dictamesh_update_updated_at_column()` |
| Trigger | `update_entity_catalog_updated_at` | `update_dictamesh_entity_catalog_updated_at` |

**Additional Issues:**
- Foreign key references pointed to non-prefixed table names
- Table comments lacked "DictaMesh:" prefix
- Missing header comment about prefix requirement

---

## Changes Made

### Modified File: `infrastructure/docker-compose/init-scripts/postgres/01-init-metadata-catalog.sql`

1. **Added header comment** about naming convention requirement
2. **Renamed all tables** to use `dictamesh_` prefix
3. **Updated all foreign key references** to point to prefixed tables
4. **Renamed all indexes** to use `idx_dictamesh_` prefix
5. **Renamed function** to `dictamesh_update_updated_at_column()`
6. **Renamed trigger** to `update_dictamesh_entity_catalog_updated_at`
7. **Updated INSERT statement** to use `dictamesh_entity_catalog`
8. **Updated table comments** to include "DictaMesh:" prefix

---

## Impact Analysis

### Database Compatibility
- ✅ Go models already expect prefixed table names
- ✅ Migration files already use prefixed names
- ✅ Changes ensure Docker initialization matches application expectations

### Potential Issues Prevented
1. **Runtime errors** - Application code would fail to find unprefixed tables
2. **Data isolation** - Ensures DictaMesh tables don't conflict with user tables
3. **Multi-tenancy** - Enables multiple frameworks to coexist in same database
4. **Clear ownership** - Unambiguous identification of framework vs. user objects

---

## Database Name Preference

As requested, the recommended database name is:
```
dictamesh
```

This can be configured via environment variables or database configuration files.

---

## Verification Checklist

- [x] All tables have `dictamesh_` prefix
- [x] All indexes have `idx_dictamesh_` prefix
- [x] All functions have `dictamesh_` prefix
- [x] All triggers reference prefixed tables/functions
- [x] GORM models override TableName()
- [x] Table comments start with "DictaMesh:"
- [x] Migration includes prefix reminder comment
- [x] All foreign key references use prefixed names
- [x] INSERT/UPDATE statements use prefixed names
- [x] All queries use prefixed names

---

## Recommendations

1. **Database Creation**: Use `dictamesh` as the database name
2. **Testing**: Verify Docker initialization script creates correct schema
3. **Documentation**: Update deployment docs to specify database naming
4. **CI/CD**: Add validation to check for naming convention compliance
5. **Code Review**: Include naming convention check in PR review process

---

## Files Modified

1. `infrastructure/docker-compose/init-scripts/postgres/01-init-metadata-catalog.sql` - Complete rewrite to comply with naming conventions

## Files Reviewed (No Changes Needed)

1. `pkg/database/migrations/sql/000001_initial_schema.up.sql`
2. `pkg/database/migrations/sql/000002_add_vector_search.up.sql`
3. `pkg/database/migrations/sql/000003_add_notifications.up.sql`
4. `pkg/database/migrations/sql/000004_add_billing.up.sql`
5. `pkg/database/models/entity.go`
6. `pkg/billing/models/models.go`
7. `pkg/notifications/models/notification.go`

---

## Conclusion

All database connectors and schemas now fully comply with the DictaMesh naming convention requirement. The Docker initialization script has been corrected to match the expectations of the Go application code and migration files, preventing runtime errors and ensuring proper namespace isolation.

**Review Status: ✅ COMPLETE**
**Naming Compliance: ✅ 100%**
