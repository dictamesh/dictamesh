# DictaMesh Framework Architecture Overview

[← Previous: Index](00-INDEX.md) | [Next: Implementation Phases →](02-IMPLEMENTATION-PHASES.md)

---

## 🎯 Purpose

This document provides framework developers with a comprehensive understanding of the DictaMesh architecture, design patterns, and component interactions.

**Reading Time:** 15 minutes
**Prerequisites:** Familiarity with PROJECT-SCOPE.md
**Outputs:** Understanding of framework layers, patterns, and extension points

---

## 🏛️ Framework Architecture Philosophy

### Core Design Principles

1. **Domain-Oriented Decentralization** (Data Mesh)
   - Framework enables source system ownership
   - Adapters provide domain-specific integration
   - No centralized data duplication

2. **Event-Driven Integration** (CQRS/Event Sourcing patterns)
   - Framework provides event bus integration
   - Immutable event log for audit and lineage
   - Asynchronous component communication

3. **Federated API Composition** (GraphQL Federation)
   - Framework provides unified API gateway
   - Adapters register their schemas
   - Intelligent cross-adapter query resolution

4. **Resilience and Observability**
   - Circuit breaker, retry, and timeout patterns
   - Built-in distributed tracing hooks
   - Comprehensive metrics and logging

### Pattern Validation

This framework architecture synthesizes proven enterprise patterns from:
- **Netflix:** Event-driven microservices, resilience patterns
- **Uber:** Real-time data mesh architecture, Kafka at scale
- **LinkedIn:** Federated data integration, schema evolution
- **Airbnb:** Multi-system integration, distributed tracing

---

## 📐 Framework Layers

### Framework Architecture Stack

```
┌─────────────────────────────────────────────────────────────┐
│ USER-BUILT APPLICATIONS & SERVICES (Out of framework scope)│
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Services consume data via GraphQL or Event Bus          │ │
│ │ (APIs, Workflows, ML, Analytics, etc.)                  │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│ DICTAMESH FRAMEWORK COMPONENTS (Provided by framework)     │
├─────────────────────────────────────────────────────────────┤
│ Layer 5: Observability & Governance Hooks                  │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ • OpenTelemetry tracing integration                     │ │
│ │ • Prometheus metrics exporters                          │ │
│ │ • PII tracking & audit hooks                            │ │
│ │ • Policy enforcement extension points                   │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 4: Federated GraphQL Gateway                         │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ • Apollo Federation engine                              │ │
│ │ • Schema composition & resolution                       │ │
│ │ • DataLoader pattern for N+1 prevention                 │ │
│ │ • Query routing to adapters                             │ │
│ │                                                         │ │
│ │   Example Subgraphs (user-defined):                    │ │
│ │   ┌──────────┐  ┌──────────┐  ┌──────────┐            │ │
│ │   │ Entity A │  │ Entity B │  │ Entity C │            │ │
│ │   │ Subgraph │  │ Subgraph │  │ Subgraph │            │ │
│ │   └──────────┘  └──────────┘  └──────────┘            │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 3: Metadata Catalog Service (Framework Component)    │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ • Entity registry (all entities across sources)         │ │
│ │ • Relationship graph (cross-system links)               │ │
│ │ • Schema registry (versioned entity schemas)            │ │
│ │ • Data lineage tracking                                 │ │
│ │ • Event log (immutable audit trail)                     │ │
│ │ [PostgreSQL-based catalog service]                      │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 2: Event Bus Integration (Framework Component)       │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ • Kafka producer/consumer abstractions                  │ │
│ │ • Schema Registry integration (Avro)                    │ │
│ │ • Topic naming conventions & patterns                   │ │
│ │ • Event schema definitions                              │ │
│ │                                                         │ │
│ │ Example Topics (user-created):                         │ │
│ │ • domain.source.entity_changed                          │ │
│ │ • system.metadata.entity_registered                     │ │
│ │ • system.lineage.relationship_created                   │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 1: Adapter Interface & Patterns (Framework Core)     │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ • DataProductAdapter interface definition               │ │
│ │ • Circuit breaker, retry, timeout patterns              │ │
│ │ • Multi-layer caching (L1 memory, L2 Redis, L3 DB)      │ │
│ │ • Event publishing abstractions                         │ │
│ │ • Health check & metrics interfaces                     │ │
│ │ • Reference implementations & examples                  │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│ USER-BUILT ADAPTERS (Out of framework scope)               │
│ ┌────────────────┐  ┌────────────────┐  ┌───────────────┐  │
│ │ Example:       │  │ Example:       │  │ Example:      │  │
│ │ CMS Adapter    │  │ API Adapter    │  │ DB Adapter    │  │
│ │ (implements    │  │ (implements    │  │ (implements   │  │
│ │ DPA interface) │  │ DPA interface) │  │ DPA interface)│  │
│ └────────┬───────┘  └────────┬───────┘  └───────┬───────┘  │
└──────────┼──────────────────┼──────────────────┼───────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│ USER'S SOURCE SYSTEMS (Out of framework scope)             │
│ • Any CMS, API, Database, File System, etc.                │
│ • Users integrate their own systems via adapters           │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔍 Layer 1: Adapter Interface & Base Patterns

### Purpose
Provide the foundational interface and implementation patterns that adapter developers use to integrate their data sources.

### What the Framework Provides

#### 1. Data Product Adapter Interface (DPI)
**The standard contract that all user-built adapters must implement:**

```go
type DataProductAdapter interface {
    // Core CRUD operations
    GetEntity(ctx context.Context, id string) (*Entity, error)
    QueryEntities(ctx context.Context, query Query) ([]Entity, error)

    // Metadata and discovery
    GetSchema() Schema
    GetSLA() ServiceLevelAgreement
    GetLineage() DataLineage

    // Event streaming
    StreamChanges(ctx context.Context) (<-chan ChangeEvent, error)

    // Health monitoring
    HealthCheck() HealthStatus
    GetMetrics() Metrics
}
```

#### 2. Reference Implementations (Examples)

The framework includes example adapters to demonstrate usage:

| Example Adapter | Purpose | Technology Stack |
|-----------------|---------|------------------|
| **CMS Example Adapter** | Shows CMS integration pattern | Go, REST client |
| **API Example Adapter** | Shows external API integration | Go, HTTP client |
| **Database Example Adapter** | Shows direct DB integration | Go, SQL driver |

**Note:** These are examples. Users build their own adapters for their systems.

#### 3. Base Implementation Patterns Provided

**Data Transformation Pattern:**
```go
// Framework provides transformation utilities
// Users implement their specific transformation logic
sourceData := sourceClient.Get("entity_type", id)
canonicalEntity := transformToCanonical(sourceData)
```

**Event Emission:**
```go
// Detect changes and emit events to Kafka
changeEvent := ChangeEvent{
    EventID:      uuid.New(),
    EventType:    "UPDATED",
    EntityType:   "customer",
    EntityID:     id,
    ChangedFields: []string{"email", "name"},
    Timestamp:    time.Now(),
}
kafkaProducer.Publish("customers.directus.entity_changed", changeEvent)
```

**Circuit Breaking:**
```go
if !circuitBreaker.Allow() {
    return nil, ErrServiceUnavailable
}
result, err := fetchFromSource(ctx, id)
if err != nil {
    circuitBreaker.RecordFailure()
} else {
    circuitBreaker.RecordSuccess()
}
```

**Multi-Layer Caching:**
```
L1: In-memory (local to adapter instance)
L2: Redis (shared across adapter replicas)
L3: PostgreSQL (metadata catalog cache)
```

### Framework Development Checklist for Layer 1

- [ ] Define DataProductAdapter interface in pkg/adapter
- [ ] Implement circuit breaker pattern (gobreaker integration)
- [ ] Implement multi-layer cache abstraction
- [ ] Create event publisher abstraction
- [ ] Implement health check interface
- [ ] Create metrics collection interface
- [ ] Write adapter contract test suite
- [ ] Create example CMS adapter implementation
- [ ] Create example API adapter implementation
- [ ] Document adapter development guide
- [ ] Provide adapter scaffolding CLI tool (optional)

---

## 🚀 Layer 2: Event-Driven Integration Fabric

### Purpose
Provide asynchronous, reliable communication between all system components using event streaming.

### Key Components

#### 1. Apache Kafka Cluster
- **Brokers:** 3+ for high availability
- **Replication Factor:** 3
- **Partitioning Strategy:** Hash by entity ID
- **Retention:** 7-90 days (topic-specific)

#### 2. Topic Taxonomy

**Pattern:** `<domain>.<source>.<event_type>`

```yaml
customers.directus.entity_changed:
  partitions: 12
  replication: 3
  retention_days: 30
  cleanup_policy: delete

products.thirdparty.entity_changed:
  partitions: 12
  replication: 3
  retention_days: 7
  cleanup_policy: delete

invoices.ecommerce.entity_changed:
  partitions: 12
  replication: 3
  retention_days: 90  # Compliance requirement
  cleanup_policy: delete

system.metadata.entity_registered:
  partitions: 3
  replication: 3
  retention_days: 365
  cleanup_policy: compact  # Keep latest entity state

system.lineage.relationship_created:
  partitions: 3
  replication: 3
  retention_days: 365
  cleanup_policy: compact
```

#### 3. Schema Registry (Confluent Schema Registry or Karapace)

**Avro Schema Example:**
```json
{
  "type": "record",
  "name": "EntityChangeEvent",
  "namespace": "com.dictamesh.events",
  "fields": [
    {"name": "event_id", "type": "string"},
    {"name": "event_type", "type": {"type": "enum", "symbols": ["CREATED", "UPDATED", "DELETED"]}},
    {"name": "timestamp", "type": "long", "logicalType": "timestamp-millis"},
    {"name": "entity", "type": {
      "type": "record",
      "name": "EntityReference",
      "fields": [
        {"name": "type", "type": "string"},
        {"name": "id", "type": "string"},
        {"name": "version", "type": "long"}
      ]
    }}
  ]
}
```

### Event Flow

```
┌─────────────┐     Produce      ┌──────────────┐     Consume      ┌──────────────┐
│  Adapter    │ ───────────────> │    Kafka     │ ───────────────> │   Metadata   │
│ (Publisher) │   Event Stream   │    Topic     │   Event Stream   │   Catalog    │
└─────────────┘                  └──────────────┘                  └──────────────┘
                                        │
                                        │ Consume
                                        ▼
                                 ┌──────────────┐
                                 │   GraphQL    │
                                 │   Gateway    │
                                 │  (Cache Inv.)│
                                 └──────────────┘
```

### LLM Agent Implementation Checklist

- [ ] Deploy Kafka cluster on Kubernetes (Strimzi operator or Helm)
- [ ] Configure topic creation with retention policies
- [ ] Deploy Schema Registry (Confluent or Karapace)
- [ ] Register Avro schemas for all event types
- [ ] Configure Kafka Connect for external integrations
- [ ] Set up Kafka UI for monitoring (kafka-ui or Conduktor)
- [ ] Configure authentication (SASL/SCRAM or mTLS)
- [ ] Set up authorization (Kafka ACLs)
- [ ] Configure monitoring (JMX → Prometheus)
- [ ] Test event production and consumption
- [ ] Document topic naming conventions
- [ ] Create runbook for topic management

---

## 🗄️ Layer 3: Metadata Catalog Service

### Purpose
Central intelligence layer that maintains entity registry, relationship graphs, schemas, and data lineage.

### Key Components

#### 1. PostgreSQL Database Schema

**Core Tables:**
```sql
entity_catalog          -- Registry of all entities
entity_relationships    -- Cross-system relationship graph
schemas                 -- Versioned entity schemas
event_log               -- Immutable event audit trail
data_lineage            -- Data flow and transformations
cache_status            -- Cache freshness tracking
```

**Indexing Strategy:**
```sql
-- High-cardinality lookups
CREATE INDEX idx_entity_type ON entity_catalog(entity_type);
CREATE INDEX idx_source_system ON entity_catalog(source_system);

-- Graph traversal optimization
CREATE INDEX idx_subject ON entity_relationships(subject_entity_type, subject_entity_id);
CREATE INDEX idx_object ON entity_relationships(object_entity_type, object_entity_id);

-- Time-series queries
CREATE INDEX idx_event_timeline ON event_log(event_timestamp DESC);
```

#### 2. Catalog Service API

**Key Operations:**
```go
RegisterEntity(EntityRegistration) error
RecordRelationship(Relationship) error
QueryRelationshipGraph(GraphQuery) (*RelationshipGraph, error)
GetEntityLocation(entityType, entityID) (*EntityLocation, error)
TrackLineage(upstream, downstream, transformation) error
```

#### 3. Relationship Graph Queries

**Recursive CTE for graph traversal:**
```sql
WITH RECURSIVE relationship_graph AS (
    -- Base: starting entity
    SELECT id, subject_entity_id, relationship_type, object_entity_id, 1 as depth
    FROM entity_relationships
    WHERE subject_entity_type = 'customer' AND subject_entity_id = '123'

    UNION ALL

    -- Recursive: traverse relationships
    SELECT er.id, er.subject_entity_id, er.relationship_type, er.object_entity_id,
           rg.depth + 1
    FROM entity_relationships er
    INNER JOIN relationship_graph rg ON er.subject_entity_id = rg.object_entity_id
    WHERE rg.depth < 5  -- Max depth
)
SELECT * FROM relationship_graph;
```

### Data Flow

```
Kafka Events ──> Catalog Consumer ──> PostgreSQL Tables
                                             │
                                             ├──> Entity Registry
                                             ├──> Relationship Graph
                                             ├──> Event Log
                                             └──> Data Lineage
```

### LLM Agent Implementation Checklist

- [ ] Deploy PostgreSQL 15+ on Kubernetes (StatefulSet)
- [ ] Apply database schema migrations (use golang-migrate or Flyway)
- [ ] Create database indexes for performance
- [ ] Set up connection pooling (PgBouncer)
- [ ] Implement Catalog Service microservice
- [ ] Configure Kafka consumer for metadata events
- [ ] Implement graph traversal queries
- [ ] Set up read replicas for query scaling
- [ ] Configure backups (pg_dump or continuous archiving)
- [ ] Add health check endpoints
- [ ] Implement metrics collection
- [ ] Write integration tests
- [ ] Create Kubernetes manifests
- [ ] Document API endpoints

---

## 🌐 Layer 4: Federated API Gateway

### Purpose
Provide unified GraphQL API that federates queries across distributed data sources.

### Key Components

#### 1. GraphQL Federation (Apollo Federation v2 or Gqlgen Federation)

**Subgraph Architecture:**
```
┌─────────────────────┐     ┌─────────────────────┐     ┌─────────────────────┐
│ Customer Subgraph   │     │ Product Subgraph    │     │ Invoice Subgraph    │
│ (Customer Adapter)  │     │ (Product Adapter)   │     │ (Invoice Adapter)   │
└─────────────────────┘     └─────────────────────┘     └─────────────────────┘
           │                           │                           │
           └───────────────────────────┴───────────────────────────┘
                                       │
                                       ▼
                            ┌──────────────────────┐
                            │  Federation Gateway  │
                            │   (Unified Schema)   │
                            └──────────────────────┘
                                       │
                                       ▼
                                  Client Query
```

#### 2. Schema Example

```graphql
# Customer subgraph
type Customer @key(fields: "id") {
  id: ID!
  email: String!
  name: String!
  invoices: [Invoice!]! @requires(fields: "id")
}

# Product subgraph
type Product @key(fields: "id") {
  id: ID!
  sku: String!
  name: String!
  price: Money!
}

# Invoice subgraph
type Invoice @key(fields: "id") {
  id: ID!
  invoiceNumber: String!
  customer: Customer! @provides(fields: "id")
  items: [InvoiceItem!]!
}

type InvoiceItem {
  product: Product!
  quantity: Int!
  total: Money!
}
```

#### 3. DataLoader Pattern (N+1 Query Prevention)

```go
type DataLoaders struct {
    CustomerLoader *dataloader.Loader
    ProductLoader  *dataloader.Loader
    InvoiceLoader  *dataloader.Loader
}

// Batches individual requests into single batch query
func batchGetCustomers(ctx context.Context, keys []string) []Result {
    customers := adapter.QueryEntities(ctx, Query{IDs: keys})
    // Map results back to keys
}
```

### Query Execution Flow

```
Client Query
    │
    ▼
┌────────────────────┐
│ Gateway analyzes   │
│ query plan         │
└────────┬───────────┘
         │
         ├──> Subgraph 1 (Customer)
         ├──> Subgraph 2 (Product)
         └──> Subgraph 3 (Invoice)
         │
         ▼
┌────────────────────┐
│ Merge results      │
│ Return to client   │
└────────────────────┘
```

### LLM Agent Implementation Checklist

- [ ] Choose federation framework (Apollo Federation or Gqlgen)
- [ ] Define federated GraphQL schemas
- [ ] Implement subgraph resolvers for each domain
- [ ] Set up DataLoader for batch loading
- [ ] Configure federation gateway
- [ ] Implement authentication middleware
- [ ] Add query complexity analysis
- [ ] Set up query caching (Redis)
- [ ] Configure rate limiting
- [ ] Add GraphQL Playground/Apollo Studio
- [ ] Implement subscription support (WebSocket)
- [ ] Write resolver integration tests
- [ ] Create Kubernetes deployment
- [ ] Document API schema and examples

---

## 📊 Layer 5: Observability & Governance

### Purpose
Provide comprehensive visibility into system behavior and enforce governance policies.

### Key Components

#### 1. Distributed Tracing (OpenTelemetry + Jaeger/Tempo)

**Trace Propagation:**
```go
import "go.opentelemetry.io/otel"

tracer := otel.Tracer("dictamesh")
ctx, span := tracer.Start(ctx, "adapter.get_entity")
defer span.End()

span.SetAttributes(
    attribute.String("entity.id", id),
    attribute.String("entity.type", "customer"),
)
```

**Trace Context in Events:**
```json
{
  "trace_context": {
    "trace_id": "5e8c05a8",
    "span_id": "7ab8c3f2",
    "parent_span_id": "2cd4e1b9"
  }
}
```

#### 2. Metrics (Prometheus + Grafana)

**Key Metrics:**
```
# Request metrics
http_requests_total{service="customer_adapter", endpoint="/api/v1/customers"}
http_request_duration_seconds{service="customer_adapter"}

# Kafka metrics
kafka_consumer_lag{topic="customers.entity_changed", consumer_group="metadata_catalog"}
kafka_producer_batch_size{topic="customers.entity_changed"}

# Database metrics
postgres_connections_active{database="metadata_catalog"}
postgres_query_duration_seconds{query="get_entity_location"}

# Cache metrics
redis_cache_hit_rate{layer="l2", entity_type="customer"}
```

#### 3. Logging (Structured JSON Logs → Loki)

```json
{
  "timestamp": "2025-11-08T10:30:45Z",
  "level": "info",
  "service": "customer_adapter",
  "trace_id": "5e8c05a8",
  "message": "Entity fetched successfully",
  "entity_id": "cust-123",
  "duration_ms": 45
}
```

### LLM Agent Implementation Checklist

- [ ] Deploy OpenTelemetry Collector
- [ ] Deploy Jaeger or Tempo for tracing
- [ ] Deploy Prometheus for metrics
- [ ] Deploy Grafana for visualization
- [ ] Deploy Loki + Promtail for logging
- [ ] Instrument all services with OpenTelemetry SDK
- [ ] Create Grafana dashboards (Golden Signals)
- [ ] Configure alert rules in Prometheus
- [ ] Set up log aggregation
- [ ] Implement audit logging
- [ ] Configure data retention policies
- [ ] Write runbooks for common alerts

---

## 🔗 Framework Component Integration Map

```
User's Data Sources (Out of scope)
       │
       ▼
User's Adapters (Implement DPA interface)
       │
       ├─────────────┐
       │             │
       ▼             │ Events
[Layer 2: Event Bus] ◄┘
       │ (Framework Component)
       │
       ├──────────────┬────────> [Layer 3: Metadata Catalog]
       │              │          (Framework Component)
       │              │                    │
       ▼              │                    ▼
[Layer 4: GraphQL    │          [Layer 5: Observability]
  Gateway]           │          (Framework Hooks & Middleware)
(Framework Component)│                    │
       │              │                    │
       │◄─────────────┴────────────────────┘
       │
       ▼
User Applications & Services
(Consume via GraphQL or Events)
```

---

## 🎯 Architecture Comprehension Validation

### Framework Developer Understanding Checklist

- [ ] Can explain data flow from user adapter to GraphQL consumer
- [ ] Understands the role of each framework layer
- [ ] Can identify which components are framework-provided vs user-built
- [ ] Knows how adapters communicate via event bus
- [ ] Understands the caching abstraction (L1/L2/L3)
- [ ] Can explain how GraphQL federation works
- [ ] Knows how observability hooks integrate
- [ ] Can describe the adapter interface contract

### Knowledge Validation Questions

1. **Q:** What does a user need to build to integrate a new data source?
   **A:** A DataProductAdapter implementation that transforms their source data to canonical entities and emits change events.

2. **Q:** How does the framework enable cross-domain queries like `entityA.relationshipB.fieldC`?
   **A:** GraphQL federation gateway composes schemas from multiple adapters, resolves relationships via metadata catalog, and uses DataLoader for efficient batching.

3. **Q:** What happens when a source system is unavailable?
   **A:** Circuit breaker pattern (provided by framework) protects the adapter → falls back to cached data → returns degraded response → emits metrics/alerts.

---

## 📚 Related Documents

- **Next Step:** [Implementation Phases](02-IMPLEMENTATION-PHASES.md)
- **Infrastructure:** [Infrastructure Planning](03-INFRASTRUCTURE-PLANNING.md)
- **Deployment:** [Deployment Strategy](04-DEPLOYMENT-STRATEGY.md)
- **Reference:** `../PROJECT-SCOPE.md` (original specification)

---

[← Previous: Index](00-INDEX.md) | [Next: Implementation Phases →](02-IMPLEMENTATION-PHASES.md)

---

**Document Metadata**
- Version: 2.0.0
- Last Updated: 2025-11-08
- Status: Updated for framework scope compliance
- Maintained By: DictaMesh Framework Contributors
