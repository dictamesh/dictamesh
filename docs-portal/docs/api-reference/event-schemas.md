<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 4
---

# Event Schemas Reference

Complete reference for DictaMesh event schemas using Apache Avro.

## Overview

DictaMesh uses Apache Avro for event serialization with the following benefits:

- **Schema Evolution** - Add fields without breaking consumers
- **Compact Binary Format** - Smaller message size than JSON
- **Strong Typing** - Compile-time validation
- **Schema Registry** - Centralized schema management
- **Code Generation** - Generate Go/Java/Python structs from schemas

## Schema Registry

### Configuration

```bash
# Schema Registry URL
SCHEMA_REGISTRY_URL=http://localhost:8081
```

### Register Schema

```bash
curl -X POST http://localhost:8081/subjects/dictamesh.entity.created-value/versions \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d @entity-created-schema.json
```

### Get Latest Schema

```bash
curl http://localhost:8081/subjects/dictamesh.entity.created-value/versions/latest
```

## Core Event Schemas

### EntityCreated

Event published when a new entity is created.

**Topic**: `dictamesh.entity.created`

**Schema**:
```json
{
  "type": "record",
  "name": "EntityCreated",
  "namespace": "com.dictamesh.events",
  "doc": "Event published when an entity is created in the metadata catalog",
  "fields": [
    {
      "name": "event_id",
      "type": "string",
      "doc": "Unique event identifier (UUID)"
    },
    {
      "name": "event_type",
      "type": "string",
      "default": "entity.created",
      "doc": "Event type identifier"
    },
    {
      "name": "event_version",
      "type": "string",
      "default": "1.0.0",
      "doc": "Event schema version"
    },
    {
      "name": "timestamp",
      "type": "long",
      "logicalType": "timestamp-millis",
      "doc": "Event timestamp in milliseconds since epoch"
    },
    {
      "name": "trace_id",
      "type": ["null", "string"],
      "default": null,
      "doc": "Distributed tracing trace ID"
    },
    {
      "name": "span_id",
      "type": ["null", "string"],
      "default": null,
      "doc": "Distributed tracing span ID"
    },
    {
      "name": "catalog_id",
      "type": "string",
      "doc": "Entity catalog ID (UUID)"
    },
    {
      "name": "entity_type",
      "type": "string",
      "doc": "Entity type (product, customer, order, etc.)"
    },
    {
      "name": "entity_id",
      "type": "string",
      "doc": "Entity identifier in source system"
    },
    {
      "name": "domain",
      "type": "string",
      "doc": "Domain (ecommerce, crm, etc.)"
    },
    {
      "name": "source_system",
      "type": "string",
      "doc": "Source system identifier"
    },
    {
      "name": "attributes",
      "type": {
        "type": "map",
        "values": ["null", "string", "long", "double", "boolean"]
      },
      "doc": "Entity attributes as key-value pairs"
    },
    {
      "name": "metadata",
      "type": {
        "type": "record",
        "name": "EventMetadata",
        "fields": [
          {
            "name": "created_by",
            "type": ["null", "string"],
            "default": null,
            "doc": "User or service that created the entity"
          },
          {
            "name": "source_event_id",
            "type": ["null", "string"],
            "default": null,
            "doc": "Source system event ID if applicable"
          },
          {
            "name": "correlation_id",
            "type": ["null", "string"],
            "default": null,
            "doc": "Correlation ID for related events"
          }
        ]
      }
    }
  ]
}
```

**Example Event**:
```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440001",
  "event_type": "entity.created",
  "event_version": "1.0.0",
  "timestamp": 1642248600000,
  "trace_id": "trace-abc123",
  "span_id": "span-def456",
  "catalog_id": "550e8400-e29b-41d4-a716-446655440002",
  "entity_type": "product",
  "entity_id": "prod-12345",
  "domain": "ecommerce",
  "source_system": "shopify",
  "attributes": {
    "name": "Premium Headphones",
    "price": 299.99,
    "sku": "WH-1000XM4",
    "inStock": true
  },
  "metadata": {
    "created_by": "shopify-adapter",
    "source_event_id": null,
    "correlation_id": "order-processing-123"
  }
}
```

### EntityUpdated

Event published when an entity is updated.

**Topic**: `dictamesh.entity.updated`

**Schema**:
```json
{
  "type": "record",
  "name": "EntityUpdated",
  "namespace": "com.dictamesh.events",
  "doc": "Event published when an entity is updated",
  "fields": [
    {
      "name": "event_id",
      "type": "string",
      "doc": "Unique event identifier (UUID)"
    },
    {
      "name": "event_type",
      "type": "string",
      "default": "entity.updated"
    },
    {
      "name": "event_version",
      "type": "string",
      "default": "1.0.0"
    },
    {
      "name": "timestamp",
      "type": "long",
      "logicalType": "timestamp-millis"
    },
    {
      "name": "trace_id",
      "type": ["null", "string"],
      "default": null
    },
    {
      "name": "span_id",
      "type": ["null", "string"],
      "default": null
    },
    {
      "name": "catalog_id",
      "type": "string"
    },
    {
      "name": "entity_type",
      "type": "string"
    },
    {
      "name": "entity_id",
      "type": "string"
    },
    {
      "name": "domain",
      "type": "string"
    },
    {
      "name": "source_system",
      "type": "string"
    },
    {
      "name": "changed_fields",
      "type": {
        "type": "array",
        "items": "string"
      },
      "doc": "List of fields that changed"
    },
    {
      "name": "previous_values",
      "type": {
        "type": "map",
        "values": ["null", "string", "long", "double", "boolean"]
      },
      "doc": "Previous values of changed fields"
    },
    {
      "name": "current_values",
      "type": {
        "type": "map",
        "values": ["null", "string", "long", "double", "boolean"]
      },
      "doc": "Current values of changed fields"
    },
    {
      "name": "metadata",
      "type": "EventMetadata"
    }
  ]
}
```

**Example Event**:
```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440003",
  "event_type": "entity.updated",
  "event_version": "1.0.0",
  "timestamp": 1642248700000,
  "catalog_id": "550e8400-e29b-41d4-a716-446655440002",
  "entity_type": "product",
  "entity_id": "prod-12345",
  "domain": "ecommerce",
  "source_system": "shopify",
  "changed_fields": ["price", "inStock"],
  "previous_values": {
    "price": 299.99,
    "inStock": true
  },
  "current_values": {
    "price": 249.99,
    "inStock": false
  },
  "metadata": {
    "created_by": "shopify-webhook-handler",
    "source_event_id": "webhook-789",
    "correlation_id": null
  }
}
```

### EntityDeleted

Event published when an entity is deleted.

**Topic**: `dictamesh.entity.deleted`

**Schema**:
```json
{
  "type": "record",
  "name": "EntityDeleted",
  "namespace": "com.dictamesh.events",
  "doc": "Event published when an entity is deleted",
  "fields": [
    {
      "name": "event_id",
      "type": "string"
    },
    {
      "name": "event_type",
      "type": "string",
      "default": "entity.deleted"
    },
    {
      "name": "event_version",
      "type": "string",
      "default": "1.0.0"
    },
    {
      "name": "timestamp",
      "type": "long",
      "logicalType": "timestamp-millis"
    },
    {
      "name": "trace_id",
      "type": ["null", "string"],
      "default": null
    },
    {
      "name": "span_id",
      "type": ["null", "string"],
      "default": null
    },
    {
      "name": "catalog_id",
      "type": "string"
    },
    {
      "name": "entity_type",
      "type": "string"
    },
    {
      "name": "entity_id",
      "type": "string"
    },
    {
      "name": "domain",
      "type": "string"
    },
    {
      "name": "source_system",
      "type": "string"
    },
    {
      "name": "deletion_reason",
      "type": ["null", "string"],
      "default": null,
      "doc": "Reason for deletion if available"
    },
    {
      "name": "metadata",
      "type": "EventMetadata"
    }
  ]
}
```

### RelationshipCreated

Event published when a relationship between entities is created.

**Topic**: `dictamesh.relationship.created`

**Schema**:
```json
{
  "type": "record",
  "name": "RelationshipCreated",
  "namespace": "com.dictamesh.events",
  "doc": "Event published when an entity relationship is created",
  "fields": [
    {
      "name": "event_id",
      "type": "string"
    },
    {
      "name": "event_type",
      "type": "string",
      "default": "relationship.created"
    },
    {
      "name": "event_version",
      "type": "string",
      "default": "1.0.0"
    },
    {
      "name": "timestamp",
      "type": "long",
      "logicalType": "timestamp-millis"
    },
    {
      "name": "trace_id",
      "type": ["null", "string"],
      "default": null
    },
    {
      "name": "relationship_id",
      "type": "string",
      "doc": "Relationship record ID"
    },
    {
      "name": "subject",
      "type": {
        "type": "record",
        "name": "EntityReference",
        "fields": [
          {"name": "catalog_id", "type": "string"},
          {"name": "entity_type", "type": "string"},
          {"name": "entity_id", "type": "string"}
        ]
      },
      "doc": "Subject entity"
    },
    {
      "name": "relationship_type",
      "type": "string",
      "doc": "Relationship type (belongs_to, has_many, etc.)"
    },
    {
      "name": "relationship_cardinality",
      "type": ["null", {
        "type": "enum",
        "name": "Cardinality",
        "symbols": ["ONE_TO_ONE", "ONE_TO_MANY", "MANY_TO_ONE", "MANY_TO_MANY"]
      }],
      "default": null
    },
    {
      "name": "object",
      "type": "EntityReference",
      "doc": "Object entity"
    },
    {
      "name": "relationship_metadata",
      "type": {
        "type": "map",
        "values": ["null", "string", "long", "double", "boolean"]
      },
      "default": {},
      "doc": "Additional relationship metadata"
    },
    {
      "name": "metadata",
      "type": "EventMetadata"
    }
  ]
}
```

### SchemaRegistered

Event published when a new schema version is registered.

**Topic**: `dictamesh.schema.registered`

**Schema**:
```json
{
  "type": "record",
  "name": "SchemaRegistered",
  "namespace": "com.dictamesh.events",
  "doc": "Event published when a schema version is registered",
  "fields": [
    {
      "name": "event_id",
      "type": "string"
    },
    {
      "name": "event_type",
      "type": "string",
      "default": "schema.registered"
    },
    {
      "name": "event_version",
      "type": "string",
      "default": "1.0.0"
    },
    {
      "name": "timestamp",
      "type": "long",
      "logicalType": "timestamp-millis"
    },
    {
      "name": "schema_id",
      "type": "string"
    },
    {
      "name": "entity_type",
      "type": "string"
    },
    {
      "name": "version",
      "type": "string",
      "doc": "Schema version (semver)"
    },
    {
      "name": "schema_format",
      "type": {
        "type": "enum",
        "name": "SchemaFormat",
        "symbols": ["AVRO", "JSON_SCHEMA", "PROTOBUF", "GRAPHQL"]
      }
    },
    {
      "name": "backward_compatible",
      "type": "boolean",
      "default": true
    },
    {
      "name": "forward_compatible",
      "type": "boolean",
      "default": false
    },
    {
      "name": "metadata",
      "type": "EventMetadata"
    }
  ]
}
```

## Domain-Specific Schemas

### Product Events

**ProductCreated** - Product entity created

**Topic**: `dictamesh.product.created`

```json
{
  "type": "record",
  "name": "ProductCreated",
  "namespace": "com.dictamesh.events.product",
  "fields": [
    {"name": "event_id", "type": "string"},
    {"name": "timestamp", "type": "long", "logicalType": "timestamp-millis"},
    {"name": "product_id", "type": "string"},
    {"name": "sku", "type": "string"},
    {"name": "name", "type": "string"},
    {"name": "description", "type": ["null", "string"], "default": null},
    {
      "name": "price",
      "type": {
        "type": "record",
        "name": "Money",
        "fields": [
          {"name": "amount", "type": "double"},
          {"name": "currency", "type": "string"}
        ]
      }
    },
    {"name": "category_id", "type": ["null", "string"], "default": null},
    {"name": "in_stock", "type": "boolean"},
    {"name": "stock_quantity", "type": ["null", "int"], "default": null},
    {
      "name": "images",
      "type": {
        "type": "array",
        "items": {
          "type": "record",
          "name": "ProductImage",
          "fields": [
            {"name": "url", "type": "string"},
            {"name": "alt", "type": ["null", "string"], "default": null},
            {"name": "is_primary", "type": "boolean", "default": false}
          ]
        }
      },
      "default": []
    },
    {"name": "tags", "type": {"type": "array", "items": "string"}, "default": []},
    {"name": "source_system", "type": "string"},
    {"name": "source_entity_id", "type": "string"}
  ]
}
```

**ProductPriceChanged** - Product price updated

**Topic**: `dictamesh.product.price_changed`

```json
{
  "type": "record",
  "name": "ProductPriceChanged",
  "namespace": "com.dictamesh.events.product",
  "fields": [
    {"name": "event_id", "type": "string"},
    {"name": "timestamp", "type": "long", "logicalType": "timestamp-millis"},
    {"name": "product_id", "type": "string"},
    {"name": "sku", "type": "string"},
    {"name": "previous_price", "type": "Money"},
    {"name": "new_price", "type": "Money"},
    {"name": "price_change_reason", "type": ["null", "string"], "default": null},
    {"name": "effective_date", "type": ["null", "long"], "logicalType": "timestamp-millis", "default": null},
    {"name": "source_system", "type": "string"}
  ]
}
```

**ProductStockChanged** - Product inventory updated

**Topic**: `dictamesh.product.stock_changed`

```json
{
  "type": "record",
  "name": "ProductStockChanged",
  "namespace": "com.dictamesh.events.product",
  "fields": [
    {"name": "event_id", "type": "string"},
    {"name": "timestamp", "type": "long", "logicalType": "timestamp-millis"},
    {"name": "product_id", "type": "string"},
    {"name": "sku", "type": "string"},
    {"name": "previous_quantity", "type": "int"},
    {"name": "new_quantity", "type": "int"},
    {"name": "quantity_delta", "type": "int"},
    {"name": "in_stock", "type": "boolean"},
    {"name": "warehouse_id", "type": ["null", "string"], "default": null},
    {"name": "source_system", "type": "string"}
  ]
}
```

## Event Envelope

All events are wrapped in a standard envelope for routing and metadata:

```json
{
  "type": "record",
  "name": "EventEnvelope",
  "namespace": "com.dictamesh.events",
  "fields": [
    {
      "name": "envelope_version",
      "type": "string",
      "default": "1.0.0"
    },
    {
      "name": "event_id",
      "type": "string",
      "doc": "Globally unique event ID"
    },
    {
      "name": "event_type",
      "type": "string",
      "doc": "Event type for routing"
    },
    {
      "name": "timestamp",
      "type": "long",
      "logicalType": "timestamp-millis"
    },
    {
      "name": "partition_key",
      "type": ["null", "string"],
      "default": null,
      "doc": "Key for partition assignment"
    },
    {
      "name": "headers",
      "type": {
        "type": "map",
        "values": "string"
      },
      "default": {},
      "doc": "Event headers for metadata"
    },
    {
      "name": "payload",
      "type": "bytes",
      "doc": "Serialized event payload"
    }
  ]
}
```

## Schema Evolution

### Backward Compatibility

Adding optional fields (with defaults) is backward compatible:

```json
// Version 1.0.0
{
  "type": "record",
  "name": "Product",
  "fields": [
    {"name": "id", "type": "string"},
    {"name": "name", "type": "string"}
  ]
}

// Version 1.1.0 - Backward compatible
{
  "type": "record",
  "name": "Product",
  "fields": [
    {"name": "id", "type": "string"},
    {"name": "name", "type": "string"},
    {"name": "description", "type": ["null", "string"], "default": null}
  ]
}
```

### Forward Compatibility

Consumers ignore unknown fields for forward compatibility.

### Breaking Changes

Require major version bump:
- Removing fields
- Changing field types
- Making optional fields required
- Renaming fields

## Code Generation

### Generate Go Structs

```bash
# Install avro-tools
go install github.com/hamba/avro/v2/cmd/avrogen@latest

# Generate Go code
avrogen -pkg events -o events.go entity-created.avsc
```

**Generated Code**:
```go
package events

import "time"

type EntityCreated struct {
    EventID      string                 `avro:"event_id"`
    EventType    string                 `avro:"event_type"`
    EventVersion string                 `avro:"event_version"`
    Timestamp    time.Time              `avro:"timestamp"`
    TraceID      *string                `avro:"trace_id"`
    SpanID       *string                `avro:"span_id"`
    CatalogID    string                 `avro:"catalog_id"`
    EntityType   string                 `avro:"entity_type"`
    EntityID     string                 `avro:"entity_id"`
    Domain       string                 `avro:"domain"`
    SourceSystem string                 `avro:"source_system"`
    Attributes   map[string]interface{} `avro:"attributes"`
    Metadata     EventMetadata          `avro:"metadata"`
}
```

## Serialization Example

### Publish Event

```go
package main

import (
    "github.com/hamba/avro/v2"
    "github.com/segmentio/kafka-go"
)

func publishEvent(event *EntityCreated) error {
    // Load schema
    schema, err := avro.Parse(entityCreatedSchema)
    if err != nil {
        return err
    }

    // Serialize
    data, err := avro.Marshal(schema, event)
    if err != nil {
        return err
    }

    // Publish to Kafka
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "dictamesh.entity.created",
    })
    defer writer.Close()

    return writer.WriteMessages(context.Background(), kafka.Message{
        Key:   []byte(event.EntityID),
        Value: data,
    })
}
```

### Consume Event

```go
func consumeEvents() error {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "dictamesh.entity.created",
        GroupID: "my-consumer-group",
    })
    defer reader.Close()

    schema, _ := avro.Parse(entityCreatedSchema)

    for {
        msg, err := reader.ReadMessage(context.Background())
        if err != nil {
            return err
        }

        var event EntityCreated
        if err := avro.Unmarshal(schema, msg.Value, &event); err != nil {
            log.Printf("Error deserializing: %v", err)
            continue
        }

        log.Printf("Received: %s for %s", event.EventType, event.EntityID)
    }
}
```

## Best Practices

### 1. Use Logical Types

```json
{
  "name": "timestamp",
  "type": "long",
  "logicalType": "timestamp-millis"
}
```

### 2. Provide Defaults

```json
{
  "name": "metadata",
  "type": ["null", "string"],
  "default": null
}
```

### 3. Document Fields

```json
{
  "name": "entity_id",
  "type": "string",
  "doc": "Unique entity identifier"
}
```

### 4. Version Events

Include version field in all events for schema evolution tracking.

### 5. Use Enums for Fixed Sets

```json
{
  "name": "status",
  "type": {
    "type": "enum",
    "name": "EntityStatus",
    "symbols": ["ACTIVE", "INACTIVE", "DEPRECATED"]
  }
}
```

## Next Steps

- [REST API Reference](./rest-api.md) - REST API for metadata catalog
- [GraphQL API Reference](./graphql-api.md) - Query the unified graph
- [Go Packages Reference](./go-packages.md) - Build adapters programmatically
- [Event Streaming Guide](../guides/event-streaming.md) - Master Kafka integration

---

**Previous**: [← Go Packages](./go-packages.md) | **Next**: [Event Streaming Guide →](../guides/event-streaming.md)
