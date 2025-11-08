# Architecture Overview

[← Previous: Index](00-INDEX.md) | [Next: Implementation Phases →](02-IMPLEMENTATION-PHASES.md)

---

## 🎯 Purpose

This document provides LLM agents with a comprehensive understanding of the DictaMesh architecture, enabling effective implementation planning and component development.

**Reading Time:** 15 minutes
**Prerequisites:** None
**Outputs:** Architectural comprehension, component relationship mapping

---

## 🏛️ Architecture Philosophy

### Core Principles

1. **Domain-Oriented Decentralization** (Data Mesh)
   - Each source system maintains data ownership
   - Domain teams manage their data products
   - No central data lake or warehouse bottleneck

2. **Event-Driven Coordination** (CQRS/Event Sourcing)
   - Immutable event log as source of truth
   - Asynchronous communication between components
   - Temporal query capabilities

3. **Federated Querying** (GraphQL Federation)
   - Unified API surface
   - Distributed schema ownership
   - Intelligent query resolution

4. **Distributed Transaction Management** (Saga Pattern)
   - Long-running business transactions
   - Compensating actions for rollback
   - Event-driven coordination

### Proven Patterns Source Validation

This architecture synthesizes battle-tested patterns from:
- **Netflix:** Event-driven microservices, chaos engineering
- **Uber:** Real-time data mesh, Kafka at scale
- **LinkedIn:** Federated data architecture, schema evolution
- **Airbnb:** Multi-system integration, distributed tracing

---

## 📐 System Layers

### Layer Stack Overview

```
┌─────────────────────────────────────────────────────────────┐
│ Layer 7: Saga Orchestration                                │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Distributed Transaction Coordination                    │ │
│ │ (Temporal Workflows, Compensation Logic)                │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 6: Multi-Tenancy & Isolation                         │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Tenant Management, Data Partitioning, Access Control   │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 5: Observability & Governance                        │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Tracing, Metrics, Logging, Audit, Compliance           │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 4: Federated API Gateway                             │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ GraphQL Federation, Unified Query Interface            │ │
│ │ ┌──────────┐  ┌──────────┐  ┌──────────┐              │ │
│ │ │ Customer │  │ Product  │  │ Invoice  │              │ │
│ │ │ Subgraph │  │ Subgraph │  │ Subgraph │              │ │
│ │ └──────────┘  └──────────┘  └──────────┘              │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 3: Metadata Catalog Service                          │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Entity Registry, Relationship Graph, Schema Registry   │ │
│ │ Data Lineage, Event Log, Cache Management              │ │
│ │ [PostgreSQL 15+ with TimescaleDB extensions]           │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 2: Event-Driven Integration Fabric                   │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Apache Kafka Event Bus + Schema Registry (Avro)        │ │
│ │                                                         │ │
│ │ Topics:                                                 │ │
│ │ • customers.directus.entity_changed                     │ │
│ │ • products.thirdparty.entity_changed                    │ │
│ │ • invoices.ecommerce.entity_changed                     │ │
│ │ • system.metadata.entity_registered                     │ │
│ │ • system.lineage.relationship_created                   │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Layer 1: Source System Adapters (Data Product Layer)       │
│ ┌────────────────┐  ┌────────────────┐  ┌───────────────┐  │
│ │ Directus       │  │ Third-Party    │  │ E-commerce    │  │
│ │ Customer       │  │ Product API    │  │ Invoice       │  │
│ │ Adapter        │  │ Adapter        │  │ Adapter       │  │
│ │                │  │                │  │               │  │
│ │ [Microservice] │  │ [Microservice] │  │ [Microservice]│  │
│ └────────┬───────┘  └────────┬───────┘  └───────┬───────┘  │
│          │                   │                   │          │
│          └───────────────────┴───────────────────┘          │
│                              │                              │
│                              ▼                              │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Source Systems (External Authority)                     │ │
│ │ • Directus CMS (Customers)                              │ │
│ │ • Third-Party APIs (Products)                           │ │
│ │ • E-commerce Platform (Invoices)                        │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔍 Layer 1: Source System Adapters

### Purpose
Transform heterogeneous external data sources into standardized Data Product interfaces.

### Key Components

#### 1. Data Product Adapter Interface (DPI)
**Standard contract all adapters must implement:**

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

#### 2. Adapter Implementations

| Adapter | Source | Domain | Technology |
|---------|--------|--------|------------|
| **DirectusCustomerAdapter** | Directus CMS | Customers | Go, Directus SDK |
| **ThirdPartyProductAdapter** | External APIs | Products | Go, HTTP client |
| **EcommerceInvoiceAdapter** | E-commerce Platform | Invoices | Go, Custom SDK |

#### 3. Adapter Responsibilities

**Data Transformation:**
```go
// Raw Directus entity → Canonical entity model
sourceData := directusClient.Get("customers", id)
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

### LLM Agent Implementation Checklist

- [ ] Scaffold adapter microservice project structure
- [ ] Implement DPI interface for specific source system
- [ ] Configure source system client/SDK
- [ ] Implement canonical entity transformation logic
- [ ] Set up Kafka producer for event emission
- [ ] Configure circuit breaker (Hystrix or similar)
- [ ] Implement multi-layer caching (in-memory + Redis)
- [ ] Add health check endpoint
- [ ] Configure Prometheus metrics
- [ ] Write integration tests with source system
- [ ] Create Kubernetes deployment manifests
- [ ] Document adapter-specific configuration

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

## 🏢 Layer 6: Multi-Tenancy & Isolation

### Purpose
Enable secure, isolated operation for multiple tenants on shared infrastructure.

### Key Components

#### 1. Tenant Isolation Strategies

**Database Level:**
```sql
-- Separate schemas per tenant
CREATE SCHEMA tenant_acme;
CREATE SCHEMA tenant_globex;

-- Row-level security
CREATE POLICY tenant_isolation ON entity_catalog
    USING (tenant_id = current_setting('app.tenant_id')::uuid);
```

**Kafka Level:**
```yaml
# Topic naming: <tenant>.<domain>.<source>.<event>
acme.customers.directus.entity_changed
globex.customers.directus.entity_changed
```

**Kubernetes Level:**
```yaml
# Namespace per tenant
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-tenant-acme
  labels:
    tenant: acme
```

#### 2. Tenant Context Propagation

```go
type TenantContext struct {
    TenantID   string
    TenantName string
    Namespace  string
    Features   []string
}

func (tc *TenantContext) InjectIntoContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, "tenant", tc)
}
```

### LLM Agent Implementation Checklist

- [ ] Design tenant isolation strategy
- [ ] Implement tenant context middleware
- [ ] Configure row-level security in PostgreSQL
- [ ] Create tenant-specific Kafka topics
- [ ] Implement tenant routing in adapters
- [ ] Set up Kubernetes namespaces per tenant
- [ ] Configure network policies for isolation
- [ ] Implement tenant-aware caching
- [ ] Add tenant quotas and rate limits
- [ ] Create tenant onboarding automation
- [ ] Document multi-tenancy architecture

---

## 🔄 Layer 7: Saga Orchestration

### Purpose
Coordinate long-running, distributed transactions across multiple services.

### Key Components

#### 1. Saga Pattern Implementation (Temporal or custom)

**Saga Example: Create Invoice with Stock Validation**
```go
func CreateInvoiceSaga(ctx workflow.Context, order Order) error {
    // Step 1: Reserve inventory
    var reservationID string
    err := workflow.ExecuteActivity(ctx, ReserveInventory, order.Items).Get(ctx, &reservationID)
    if err != nil {
        return err  // No compensation needed
    }

    // Step 2: Create invoice
    var invoiceID string
    err = workflow.ExecuteActivity(ctx, CreateInvoice, order).Get(ctx, &invoiceID)
    if err != nil {
        // Compensation: release inventory
        workflow.ExecuteActivity(ctx, ReleaseInventory, reservationID)
        return err
    }

    // Step 3: Charge customer
    err = workflow.ExecuteActivity(ctx, ChargeCustomer, order.CustomerID, invoiceID)
    if err != nil {
        // Compensation: cancel invoice and release inventory
        workflow.ExecuteActivity(ctx, CancelInvoice, invoiceID)
        workflow.ExecuteActivity(ctx, ReleaseInventory, reservationID)
        return err
    }

    return nil
}
```

#### 2. Saga State Machine

```
[Start] → [Reserve Inventory] → [Create Invoice] → [Charge Customer] → [Complete]
              │                       │                    │
              ▼ (failure)             ▼ (failure)          ▼ (failure)
         [Rollback]              [Cancel Invoice]     [Refund + Cancel]
                                 [Release Inventory]   [Release Inventory]
```

### LLM Agent Implementation Checklist

- [ ] Choose saga framework (Temporal, Cadence, or custom)
- [ ] Define saga workflows for business processes
- [ ] Implement activity handlers
- [ ] Implement compensation logic
- [ ] Set up saga state persistence
- [ ] Configure saga timeouts and retries
- [ ] Add saga monitoring dashboard
- [ ] Write saga integration tests
- [ ] Document saga patterns

---

## 🔗 Component Integration Map

```
External Sources
       │
       ▼
[Layer 1: Adapters] ──┐
       │              │
       ▼              │
[Layer 2: Kafka] ◄────┤
       │              │
       ├──────────────┼────────> [Layer 3: Metadata Catalog]
       │              │                    │
       ▼              │                    ▼
[Layer 4: GraphQL] ◄──┘          [Layer 5: Observability]
       │                                   │
       ▼                                   │
[Layer 6: Multi-tenancy] ◄────────────────┘
       │
       ▼
[Layer 7: Saga Orchestration]
```

---

## 🎯 Success Criteria for Architecture Understanding

### LLM Agent Self-Check

- [ ] Can explain data flow from source system to client query
- [ ] Understands role of each layer and dependencies
- [ ] Can identify which components run as microservices
- [ ] Knows which components use Kafka for communication
- [ ] Understands caching strategy (L1/L2/L3)
- [ ] Can explain how federation resolves cross-domain queries
- [ ] Knows how distributed tracing works across components
- [ ] Understands tenant isolation mechanisms
- [ ] Can describe saga compensation flows

### Validation Questions

1. **Q:** If a customer email changes in Directus, what happens?
   **A:** Directus adapter detects change → emits event to Kafka → Metadata catalog consumes and updates → GraphQL cache invalidated

2. **Q:** How does GraphQL resolve `customer.invoices.items.product`?
   **A:** Customer subgraph → Invoice subgraph (via federation) → Product subgraph → DataLoader batches requests

3. **Q:** What happens if Product API is down?
   **A:** Circuit breaker opens → return cached data if available → fall back to degraded response → metrics/alerts triggered

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
- Version: 1.0.0
- Last Updated: 2025-11-08
- LLM Agent Checkpoint: Architecture comprehension complete
