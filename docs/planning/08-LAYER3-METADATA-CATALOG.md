# Layer 3: Metadata Catalog Service

[← Previous: Layer 2 Event Bus](07-LAYER2-EVENT-BUS.md) | [Next: Layer 4 API Gateway →](09-LAYER4-API-GATEWAY.md)

---

## 🎯 Purpose

Implementation guide for the Metadata Catalog Service - the central intelligence layer managing entity registry, relationship graphs, and data lineage.

**Reading Time:** 25 minutes
**Prerequisites:** [Layer 2 Event Bus](07-LAYER2-EVENT-BUS.md), PostgreSQL deployed
**Outputs:** Metadata catalog service, database schema, Kafka consumers

---

## 📐 Architecture

The Metadata Catalog is the system's brain, maintaining:
- **Entity Registry:** All entities across source systems
- **Relationship Graph:** Cross-system entity relationships
- **Schema Registry:** Versioned entity schemas
- **Event Log:** Immutable audit trail
- **Data Lineage:** Data flow and transformations

---

## 🗄️ PostgreSQL Deployment

### Using CloudNativePG on K3S

```yaml
# infrastructure/k8s/postgres/metadata-catalog-cluster.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: metadata-catalog-db
  namespace: dictamesh-infra
spec:
  instances: 3
  imageName: ghcr.io/cloudnative-pg/postgresql:15.5
  storage:
    size: 100Gi
    storageClass: longhorn-retain
  postgresql:
    parameters:
      max_connections: "200"
      shared_buffers: "2GB"
  monitoring:
    enabled: true
  backup:
    barmanObjectStore:
      destinationPath: s3://dictamesh-backups/postgres
```

---

## 📊 Database Schema

### Migration Tool Setup

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq create_entity_catalog
```

### Core Tables

```sql
-- migrations/000001_create_entity_catalog.up.sql
CREATE TABLE entity_catalog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(100) NOT NULL,
    domain VARCHAR(100) NOT NULL,
    source_system VARCHAR(100) NOT NULL,
    source_entity_id VARCHAR(255) NOT NULL,
    api_base_url TEXT NOT NULL,
    api_path_template TEXT NOT NULL,
    schema_version VARCHAR(50),
    contains_pii BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(source_system, source_entity_id, entity_type)
);

CREATE INDEX idx_entity_type ON entity_catalog(entity_type);
CREATE INDEX idx_source_system ON entity_catalog(source_system);

-- migrations/000002_create_relationships.up.sql
CREATE TABLE entity_relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_entity_type VARCHAR(100),
    subject_entity_id VARCHAR(255),
    relationship_type VARCHAR(100),
    object_entity_type VARCHAR(100),
    object_entity_id VARCHAR(255),
    valid_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_subject ON entity_relationships(subject_entity_type, subject_entity_id);
CREATE INDEX idx_object ON entity_relationships(object_entity_type, object_entity_id);
```

---

## 🚀 Catalog Service Implementation

### Project Structure

```
services/metadata-catalog/
├── cmd/server/main.go
├── internal/
│   ├── catalog/
│   │   ├── service.go          # Core catalog operations
│   │   └── handlers.go         # HTTP/gRPC handlers
│   ├── consumer/
│   │   └── kafka_consumer.go   # Event consumer
│   └── repository/
│       └── postgres.go         # Database operations
├── migrations/
└── go.mod
```

### Kafka Event Consumer

```go
// internal/consumer/kafka_consumer.go
package consumer

type KafkaConsumer struct {
    reader         *kafka.Reader
    catalogService *catalog.Service
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
    c.reader = kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"dictamesh-kafka-kafka-bootstrap:9092"},
        GroupID: "metadata-catalog-consumer",
        Topic:   "customers.directus.entity_changed",
    })

    for {
        msg, err := c.reader.ReadMessage(ctx)
        if err != nil {
            log.Printf("Error reading message: %v", err)
            continue
        }

        var event ChangeEvent
        if err := avro.Unmarshal(msg.Value, &event); err != nil {
            log.Printf("Error unmarshaling event: %v", err)
            continue
        }

        // Process event
        if err := c.catalogService.RegisterEntity(ctx, event); err != nil {
            log.Printf("Error processing event: %v", err)
        }
    }
}
```

---

## 🔍 Graph Queries

```go
// internal/catalog/service.go
func (s *Service) QueryRelationshipGraph(ctx context.Context, query GraphQuery) (*RelationshipGraph, error) {
    rows, err := s.db.QueryContext(ctx, `
        WITH RECURSIVE relationship_graph AS (
            SELECT id, subject_entity_id, relationship_type, object_entity_id, 1 as depth
            FROM entity_relationships
            WHERE subject_entity_type = $1 AND subject_entity_id = $2

            UNION ALL

            SELECT er.id, er.subject_entity_id, er.relationship_type, er.object_entity_id, rg.depth + 1
            FROM entity_relationships er
            INNER JOIN relationship_graph rg ON er.subject_entity_id = rg.object_entity_id
            WHERE rg.depth < $3
        )
        SELECT * FROM relationship_graph
    `, query.EntityType, query.EntityID, query.MaxDepth)

    // Process results...
}
```

---

[← Previous: Layer 2 Event Bus](07-LAYER2-EVENT-BUS.md) | [Next: Layer 4 API Gateway →](09-LAYER4-API-GATEWAY.md)
