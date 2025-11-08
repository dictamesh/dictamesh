# DictaMesh Database Package

Comprehensive database infrastructure for DictaMesh framework with advanced features for data source agnostic metadata management, vector search, caching, and compliance.

## Features

### Core Database Management
- **Multiple Connection Pools**: pgx for performance, GORM for ORM, standard database/sql for compatibility
- **Advanced Connection Pooling**: Configurable pool sizes, timeouts, and lifecycle management
- **Transaction Management**: Support for both GORM and pgx transactions

### Migration System
- **Versioned Migrations**: golang-migrate integration with embedded SQL files
- **Forward/Backward Support**: Up and down migrations for schema evolution
- **Migration Validation**: Check migration status and detect dirty states

### Vector Search & RAG
- **pgvector Integration**: Native PostgreSQL vector similarity search
- **Semantic Search**: Find similar entities using cosine similarity
- **Document Chunking**: Store and retrieve document chunks for RAG
- **Hybrid Search**: Combine full-text and vector search
- **HNSW Indexing**: Fast approximate nearest neighbor search

### Multi-Layer Caching
- **L1 In-Memory Cache**: Ultra-fast local caching with LRU eviction
- **L2 Redis Cache**: Distributed caching for scaling across replicas
- **L3 Database Cache**: Metadata tracking for cache status
- **Cache Metrics**: Hit rates, evictions, and performance tracking

### Health Monitoring
- **Comprehensive Health Checks**: Connection, query execution, and replication lag
- **Pool Statistics**: Track open connections, wait times, and resource usage
- **Table Statistics**: Monitor row counts, vacuum status, and index health
- **Extension Checking**: Verify required PostgreSQL extensions

### Audit & Compliance
- **Comprehensive Audit Logging**: Track all data access and modifications
- **PII Access Tracking**: Special logging for sensitive data access
- **Compliance Reports**: Query audit logs for compliance verification
- **Distributed Tracing Integration**: Correlate audit logs with traces

### Repository Pattern
- **Type-Safe Models**: GORM models with proper relationships
- **Repository Implementations**: Catalog, Relationship, Schema repositories
- **Query Builders**: Flexible filtering and pagination
- **Preloading Support**: Eager load related entities

## Installation

```bash
go get github.com/click2-run/dictamesh/pkg/database
```

## Quick Start

### Basic Setup

```go
import (
    "github.com/click2-run/dictamesh/pkg/database"
    "go.uber.org/zap"
)

// Create configuration
config := database.DefaultConfig()
config.Host = "localhost"
config.Port = 5432
config.User = "dictamesh"
config.Password = "your_password"
config.Database = "metadata_catalog"

// Create logger
logger, _ := zap.NewProduction()

// Create database instance
db, err := database.New(config, logger)
if err != nil {
    log.Fatal(err)
}

// Connect
ctx := context.Background()
if err := db.Connect(ctx); err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### Running Migrations

```go
import "github.com/click2-run/dictamesh/pkg/database/migrations"

// Create migrator
migrator, err := migrations.NewMigrator(db.StdDB(), logger)
if err != nil {
    log.Fatal(err)
}
defer migrator.Close()

// Run all pending migrations
if err := migrator.Up(ctx); err != nil {
    log.Fatal(err)
}

// Check migration status
version, dirty, err := migrator.Version()
fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)
```

### Vector Search

```go
import (
    "github.com/click2-run/dictamesh/pkg/database"
    "github.com/pgvector/pgvector-go"
)

// Create vector search instance
vs := database.NewVectorSearch(db)

// Store an embedding
embedding := &database.EntityEmbedding{
    CatalogID:          "entity-123",
    EmbeddingModel:     "text-embedding-ada-002",
    EmbeddingVersion:   "v1",
    EmbeddingDimensions: 1536,
    Embedding:          pgvector.NewVector([]float32{0.1, 0.2, ...}),
    SourceText:         "This is the text that was embedded",
}

if err := vs.StoreEmbedding(ctx, embedding); err != nil {
    log.Fatal(err)
}

// Find similar entities
queryVector := pgvector.NewVector([]float32{0.15, 0.25, ...})
similar, err := vs.FindSimilarEntities(
    ctx,
    queryVector,
    "text-embedding-ada-002",
    0.7,  // similarity threshold
    10,   // limit
)

for _, entity := range similar {
    fmt.Printf("Entity: %s, Similarity: %.4f\n", entity.CatalogID, entity.Similarity)
}

// RAG: Find relevant chunks
chunks, err := vs.FindRelevantChunks(
    ctx,
    queryVector,
    "text-embedding-ada-002",
    nil,  // no entity filter
    0.7,  // similarity threshold
    5,    // top 5 chunks
)

for _, chunk := range chunks {
    fmt.Printf("Chunk: %s\nText: %s\n\n", chunk.ChunkID, chunk.ChunkText)
}
```

### Caching

```go
import "github.com/click2-run/dictamesh/pkg/database/cache"

// Create cache
cacheConfig := cache.DefaultConfig()
cacheConfig.RedisURL = "redis://localhost:6379"

cache, err := cache.New(cacheConfig, logger)
if err != nil {
    log.Fatal(err)
}
defer cache.Close()

// Store value
data := []byte("cached data")
if err := cache.Set(ctx, "my-key", data, 5*time.Minute); err != nil {
    log.Error(err)
}

// Retrieve value
value, err := cache.Get(ctx, "my-key")
if err == nil {
    fmt.Printf("Retrieved: %s\n", string(value))
}

// Store JSON
type MyData struct {
    Name string
    Value int
}
if err := cache.SetJSON(ctx, "json-key", MyData{"test", 42}, 10*time.Minute); err != nil {
    log.Error(err)
}

// Get metrics
metrics := cache.GetMetrics()
fmt.Printf("L1 Hits: %d, L2 Hits: %d\n", metrics.L1Hits, metrics.L2Hits)
```

### Health Checks

```go
import "github.com/click2-run/dictamesh/pkg/database/health"

// Create health checker
checker := health.NewChecker(db.Pool(), db.StdDB(), logger)

// Perform health check
result := checker.Check(ctx)
fmt.Printf("Status: %s\nMessage: %s\nResponse Time: %v\n",
    result.Status, result.Message, result.ResponseTime)

// Check specific table
if err := checker.CheckTable(ctx, "entity_catalog"); err != nil {
    log.Printf("Table check failed: %v", err)
}

// Check extension
if err := checker.CheckExtension(ctx, "vector"); err != nil {
    log.Printf("Extension check failed: %v", err)
}

// Get table statistics
stats, err := checker.GetTableStats(ctx, "entity_catalog")
if err == nil {
    fmt.Printf("Live tuples: %d\n", stats["live_tuples"])
}
```

### Audit Logging

```go
import "github.com/click2-run/dictamesh/pkg/database/audit"

// Create audit logger
auditConfig := &audit.Config{Enabled: true}
auditLogger := audit.NewLogger(db.Pool(), logger, auditConfig)

// Create audit table
if err := auditLogger.CreateAuditTable(ctx); err != nil {
    log.Fatal(err)
}

// Log an operation
entry := &audit.AuditLog{
    UserID:       "user-123",
    UserEmail:    "user@example.com",
    Operation:    audit.OpUpdate,
    ResourceType: "customer",
    ResourceID:   "cust-456",
    Changes: map[string]interface{}{
        "email": "newemail@example.com",
    },
    Success:    true,
    TraceID:    "trace-789",
}

if err := auditLogger.Log(ctx, entry); err != nil {
    log.Error(err)
}

// Log PII access
if err := auditLogger.LogDataAccess(ctx, "user-123", "customer", "cust-456",
    []string{"ssn", "credit_card"}); err != nil {
    log.Error(err)
}

// Query audit logs
filters := &audit.QueryFilters{
    UserID:     "user-123",
    StartTime:  time.Now().Add(-24 * time.Hour),
    EndTime:    time.Now(),
    Limit:      100,
}

logs, err := auditLogger.Query(ctx, filters)
for _, log := range logs {
    fmt.Printf("%s: %s on %s\n", log.Timestamp, log.Operation, log.ResourceType)
}
```

### Repository Pattern

```go
import (
    "github.com/click2-run/dictamesh/pkg/database/models"
    "github.com/click2-run/dictamesh/pkg/database/repository"
)

// Create repositories
catalogRepo := repository.NewCatalogRepository(db.GORM())
relationshipRepo := repository.NewRelationshipRepository(db.GORM())

// Create entity
entity := &models.EntityCatalog{
    EntityType:     "customer",
    Domain:         "customers",
    SourceSystem:   "directus",
    SourceEntityID: "12345",
    APIBaseURL:     "https://api.example.com",
    APIPathTemplate: "/customers/{id}",
    Status:         "active",
}

if err := catalogRepo.Create(ctx, entity); err != nil {
    log.Fatal(err)
}

// Find entity
found, err := catalogRepo.FindBySource(ctx, "directus", "12345", "customer")
if err != nil {
    log.Fatal(err)
}

// List entities with filters
filters := &repository.CatalogFilters{
    EntityType: "customer",
    Domain:     "customers",
    Status:     "active",
    Limit:      50,
    Offset:     0,
}

entities, err := catalogRepo.List(ctx, filters)
fmt.Printf("Found %d entities\n", len(entities))
```

## Configuration

### Database Config

```go
config := &database.Config{
    Host:     "localhost",
    Port:     5432,
    User:     "dictamesh",
    Password: "password",
    Database: "metadata_catalog",
    SSLMode:  "prefer",

    // Connection pool
    MaxOpenConns:    25,
    MaxIdleConns:    10,
    ConnMaxLifetime: 30 * time.Minute,
    ConnMaxIdleTime: 10 * time.Minute,

    // Performance
    StatementTimeout: 30 * time.Second,
    IdleInTxTimeout:  60 * time.Second,

    // Features
    EnableMigrations:   true,
    EnableVectorSearch: true,
    EnableAuditLog:     true,

    // Observability
    EnableMetrics: true,
    EnableTracing: true,
    LogLevel:      "info",
}
```

## Schema

The database package includes comprehensive schema migrations:

- **000001_initial_schema.up.sql**: Core metadata catalog tables
- **000002_add_vector_search.up.sql**: Vector embeddings and RAG support

### Tables

- `entity_catalog`: Registry of all entities
- `entity_relationships`: Cross-system relationships with temporal validity
- `schemas`: Versioned entity schemas
- `event_log`: Immutable audit trail
- `data_lineage`: Data flow tracking
- `cache_status`: Cache freshness tracking
- `entity_embeddings`: Vector embeddings for semantic search
- `document_chunks`: Document chunks for RAG
- `audit_logs`: Comprehensive audit logging

## Performance Tips

1. **Use pgx for high-performance queries**: Direct pool access for bulk operations
2. **Enable caching**: Reduce database load with multi-layer caching
3. **Batch operations**: Use transactions for multiple operations
4. **Optimize indexes**: Vector search uses HNSW for fast retrieval
5. **Monitor connections**: Track pool stats and adjust configuration

## Security

- **Connection encryption**: SSL/TLS support
- **Prepared statements**: Protection against SQL injection
- **Audit logging**: Track all sensitive data access
- **PII tracking**: Special handling for personal information

## License

AGPL-3.0-or-later - Copyright (C) 2025 Controle Digital Ltda
