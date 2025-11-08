<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 1
---

# REST API Reference

Complete reference for the DictaMesh REST API for metadata catalog operations.

## Base URL

```
http://localhost:8080/api/v1
```

In production, replace with your deployed service URL.

## Authentication

All API requests require authentication using Bearer tokens:

```bash
Authorization: Bearer <your-token>
```

### Obtaining a Token

```bash
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "your-client-id",
    "client_secret": "your-client-secret"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

## Entity Catalog API

### Create Entity

Register a new entity in the metadata catalog.

```http
POST /catalog/entities
```

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "entity_type": "product",
  "domain": "ecommerce",
  "source_system": "shopify",
  "source_entity_id": "prod-12345",
  "api_base_url": "https://api.shopify.com/v1",
  "api_path_template": "/products/{id}",
  "api_method": "GET",
  "api_auth_type": "bearer",
  "schema_id": "550e8400-e29b-41d4-a716-446655440000",
  "schema_version": "1.0.0",
  "status": "active",
  "availability_sla": 0.999,
  "latency_p99_ms": 200,
  "freshness_sla": 300,
  "contains_pii": true,
  "data_classification": "confidential",
  "retention_days": 365
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "entity_type": "product",
  "domain": "ecommerce",
  "source_system": "shopify",
  "source_entity_id": "prod-12345",
  "api_base_url": "https://api.shopify.com/v1",
  "api_path_template": "/products/{id}",
  "api_method": "GET",
  "api_auth_type": "bearer",
  "schema_id": "550e8400-e29b-41d4-a716-446655440000",
  "schema_version": "1.0.0",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z",
  "last_seen_at": "2025-01-15T10:30:00Z",
  "status": "active",
  "availability_sla": 0.999,
  "latency_p99_ms": 200,
  "freshness_sla": 300,
  "contains_pii": true,
  "data_classification": "confidential",
  "retention_days": 365
}
```

### Get Entity by ID

Retrieve a specific entity by its catalog ID.

```http
GET /catalog/entities/{id}
```

**Parameters:**
- `id` (path, required) - Entity catalog UUID

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "entity_type": "product",
  "domain": "ecommerce",
  "source_system": "shopify",
  "source_entity_id": "prod-12345",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z",
  "status": "active"
}
```

### List Entities

Retrieve a paginated list of entities with optional filtering.

```http
GET /catalog/entities
```

**Query Parameters:**
- `entity_type` (optional) - Filter by entity type (e.g., "product", "customer")
- `domain` (optional) - Filter by domain (e.g., "ecommerce", "crm")
- `source_system` (optional) - Filter by source system (e.g., "shopify", "salesforce")
- `status` (optional) - Filter by status ("active", "inactive", "deprecated")
- `contains_pii` (optional) - Filter by PII flag (true/false)
- `limit` (optional) - Results per page (default: 20, max: 100)
- `offset` (optional) - Offset for pagination (default: 0)

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/entities?entity_type=product&domain=ecommerce&limit=10" \
  -H "Authorization: Bearer <token>"
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "entity_type": "product",
      "domain": "ecommerce",
      "source_system": "shopify",
      "source_entity_id": "prod-12345",
      "status": "active"
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 10,
    "offset": 0,
    "has_next": true
  }
}
```

### Update Entity

Update an existing entity in the catalog.

```http
PUT /catalog/entities/{id}
```

**Request Body:**
```json
{
  "status": "deprecated",
  "availability_sla": 0.995,
  "latency_p99_ms": 250
}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "entity_type": "product",
  "status": "deprecated",
  "updated_at": "2025-01-15T11:00:00Z"
}
```

### Delete Entity

Remove an entity from the catalog.

```http
DELETE /catalog/entities/{id}
```

**Response (204 No Content)**

### Find Entity by Source

Locate an entity by its source system and source ID.

```http
GET /catalog/entities/by-source
```

**Query Parameters:**
- `source_system` (required) - Source system identifier
- `source_entity_id` (required) - Entity ID in source system
- `entity_type` (required) - Entity type

**Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/entities/by-source?source_system=shopify&source_entity_id=prod-12345&entity_type=product" \
  -H "Authorization: Bearer <token>"
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "entity_type": "product",
  "source_system": "shopify",
  "source_entity_id": "prod-12345"
}
```

## Entity Relationships API

### Create Relationship

Create a relationship between two entities.

```http
POST /catalog/relationships
```

**Request Body:**
```json
{
  "subject_catalog_id": "550e8400-e29b-41d4-a716-446655440001",
  "subject_entity_type": "product",
  "subject_entity_id": "prod-12345",
  "relationship_type": "belongs_to",
  "relationship_cardinality": "many_to_one",
  "object_catalog_id": "550e8400-e29b-41d4-a716-446655440002",
  "object_entity_type": "category",
  "object_entity_id": "cat-100",
  "subject_display_name": "Premium Headphones",
  "object_display_name": "Electronics",
  "relationship_metadata": {
    "confidence": 1.0,
    "source": "explicit"
  }
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "subject_catalog_id": "550e8400-e29b-41d4-a716-446655440001",
  "relationship_type": "belongs_to",
  "object_catalog_id": "550e8400-e29b-41d4-a716-446655440002",
  "valid_from": "2025-01-15T10:30:00Z",
  "valid_to": null,
  "created_at": "2025-01-15T10:30:00Z"
}
```

### Get Entity Relationships

Retrieve relationships for an entity.

```http
GET /catalog/entities/{id}/relationships
```

**Query Parameters:**
- `direction` (optional) - "outgoing", "incoming", or "all" (default: "all")
- `relationship_type` (optional) - Filter by relationship type

**Response (200 OK):**
```json
{
  "outgoing": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "relationship_type": "belongs_to",
      "object_entity_type": "category",
      "object_entity_id": "cat-100",
      "object_display_name": "Electronics"
    }
  ],
  "incoming": []
}
```

## Schema API

### Register Schema

Register a new schema version.

```http
POST /catalog/schemas
```

**Request Body:**
```json
{
  "entity_type": "product",
  "version": "1.0.0",
  "schema_format": "avro",
  "schema_definition": {
    "type": "record",
    "name": "Product",
    "namespace": "com.dictamesh.ecommerce",
    "fields": [
      {"name": "id", "type": "string"},
      {"name": "name", "type": "string"},
      {"name": "price", "type": "double"}
    ]
  },
  "backward_compatible": true,
  "forward_compatible": false
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "entity_type": "product",
  "version": "1.0.0",
  "schema_format": "avro",
  "published_at": "2025-01-15T10:00:00Z"
}
```

### Get Schema

Retrieve a specific schema version.

```http
GET /catalog/schemas/{entity_type}/{version}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "entity_type": "product",
  "version": "1.0.0",
  "schema_format": "avro",
  "schema_definition": {...},
  "backward_compatible": true,
  "forward_compatible": false,
  "published_at": "2025-01-15T10:00:00Z"
}
```

### List Schemas

List all schema versions for an entity type.

```http
GET /catalog/schemas/{entity_type}
```

**Response (200 OK):**
```json
{
  "entity_type": "product",
  "schemas": [
    {
      "version": "1.0.0",
      "published_at": "2025-01-15T10:00:00Z",
      "deprecated_at": null
    },
    {
      "version": "1.1.0",
      "published_at": "2025-02-01T10:00:00Z",
      "deprecated_at": null
    }
  ]
}
```

## Event Log API

### Query Events

Query the event log with filters.

```http
GET /catalog/events
```

**Query Parameters:**
- `event_type` (optional) - Filter by event type
- `entity_type` (optional) - Filter by entity type
- `entity_id` (optional) - Filter by entity ID
- `trace_id` (optional) - Filter by trace ID
- `start_time` (optional) - ISO 8601 timestamp
- `end_time` (optional) - ISO 8601 timestamp
- `limit` (optional) - Results per page (default: 50, max: 500)
- `offset` (optional) - Offset for pagination

**Response (200 OK):**
```json
{
  "events": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "event_id": "evt-12345",
      "event_type": "entity.created",
      "entity_type": "product",
      "entity_id": "prod-12345",
      "event_timestamp": "2025-01-15T10:30:00Z",
      "trace_id": "trace-abc123"
    }
  ],
  "pagination": {
    "total": 1000,
    "limit": 50,
    "offset": 0
  }
}
```

## Data Lineage API

### Get Lineage

Retrieve data lineage for an entity.

```http
GET /catalog/lineage/{catalog_id}
```

**Query Parameters:**
- `direction` (optional) - "upstream", "downstream", or "both" (default: "both")
- `depth` (optional) - Maximum traversal depth (default: 3, max: 10)

**Response (200 OK):**
```json
{
  "entity_id": "550e8400-e29b-41d4-a716-446655440001",
  "upstream": [
    {
      "catalog_id": "550e8400-e29b-41d4-a716-446655440020",
      "system": "shopify",
      "entity_type": "product",
      "transformation_type": "extract",
      "last_flow_at": "2025-01-15T10:25:00Z"
    }
  ],
  "downstream": [
    {
      "catalog_id": "550e8400-e29b-41d4-a716-446655440021",
      "system": "data_warehouse",
      "entity_type": "product_dimension",
      "transformation_type": "transform",
      "last_flow_at": "2025-01-15T10:35:00Z"
    }
  ]
}
```

## Health & Status API

### Health Check

Check service health status.

```http
GET /health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2025-01-15T10:30:00Z",
  "checks": {
    "database": {
      "status": "up",
      "latency_ms": 5
    },
    "kafka": {
      "status": "up",
      "latency_ms": 8
    },
    "redis": {
      "status": "up",
      "latency_ms": 2
    }
  }
}
```

### Readiness Check

Check if service is ready to accept requests.

```http
GET /ready
```

**Response (200 OK):**
```json
{
  "ready": true,
  "timestamp": "2025-01-15T10:30:00Z"
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "ENTITY_NOT_FOUND",
    "message": "Entity with ID 550e8400-e29b-41d4-a716-446655440001 not found",
    "details": {
      "entity_id": "550e8400-e29b-41d4-a716-446655440001"
    },
    "trace_id": "trace-abc123"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | Request validation failed |
| `UNAUTHORIZED` | 401 | Authentication required |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `ENTITY_NOT_FOUND` | 404 | Entity not found |
| `CONFLICT` | 409 | Entity already exists |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

## Rate Limiting

API requests are rate limited per client:

- **Default**: 1000 requests per minute
- **Burst**: 100 requests per second

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 950
X-RateLimit-Reset: 1642248600
```

## Pagination

All list endpoints support cursor-based pagination:

```bash
# First page
GET /catalog/entities?limit=20

# Next page
GET /catalog/entities?limit=20&offset=20
```

Response includes pagination metadata:

```json
{
  "data": [...],
  "pagination": {
    "total": 150,
    "limit": 20,
    "offset": 20,
    "has_next": true,
    "has_previous": true
  }
}
```

## Next Steps

- [GraphQL API Reference](./graphql-api.md) - Query the unified graph
- [Go Packages Reference](./go-packages.md) - Build adapters programmatically
- [Event Schemas Reference](./event-schemas.md) - Event-driven integration

---

**Previous**: [← Metadata Catalog](../architecture/metadata-catalog.md) | **Next**: [GraphQL API →](./graphql-api.md)
