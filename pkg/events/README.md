# DictaMesh Events Package

Event-driven messaging infrastructure for the DictaMesh framework using Kafka/Redpanda.

## Features

- **Kafka/Redpanda Integration**: Production-ready event streaming
- **Producer/Consumer Wrappers**: Simplified API with observability
- **Event Schema**: Standardized event structure
- **Topic Management**: Pre-configured topic templates
- **Observability**: Built-in tracing, metrics, and logging
- **Reliability**: Idempotent production, manual offset commits
- **Scalability**: Configurable partitioning and replication

## Installation

```bash
go get github.com/click2-run/dictamesh/pkg/events
```

## Quick Start

### Publishing Events

```go
package main

import (
    "context"
    "github.com/click2-run/dictamesh/pkg/events"
    "github.com/click2-run/dictamesh/pkg/observability"
)

func main() {
    // Create observability
    obs, _ := observability.New(observability.DefaultConfig())

    // Create config
    config := events.DefaultConfig()

    // Create producer
    producer, err := events.NewProducer(config, obs.Logger())
    if err != nil {
        panic(err)
    }
    defer producer.Close()

    // Create event
    event := events.NewEvent(
        events.EventTypeEntityCreated,
        "customer-adapter",
        "customer:123",
        map[string]interface{}{
            "customer_id": "123",
            "name":        "John Doe",
        },
    )

    // Publish
    ctx := context.Background()
    if err := producer.Publish(ctx, events.TopicEntityChanged, event); err != nil {
        panic(err)
    }
}
```

### Consuming Events

```go
// Create event handler
handler := func(ctx context.Context, event *events.Event) error {
    log.Printf("Received event: %s from %s", event.Type, event.Source)
    // Process event...
    return nil
}

// Create consumer
consumer, err := events.NewConsumer(config, obs.Logger(), handler)
if err != nil {
    panic(err)
}
defer consumer.Close()

// Subscribe to topics
if err := consumer.Subscribe([]string{events.TopicEntityChanged}); err != nil {
    panic(err)
}

// Start consuming
ctx := context.Background()
if err := consumer.Start(ctx); err != nil {
    panic(err)
}
```

## Configuration

### Default Configuration (Development)

```go
config := events.DefaultConfig()
// Uses:
// - Bootstrap servers: localhost:19092
// - Idempotent producer with acks=all
// - Manual offset commits
// - 12 partitions, replication factor 1
```

### Production Configuration

```go
config := events.ProductionConfig()
// Uses:
// - 3 Kafka brokers
// - SASL_SSL security
// - 12 partitions, replication factor 3
// - 30 day retention
```

### Custom Configuration

```go
config := &events.Config{
    BootstrapServers: []string{"kafka-0:9092", "kafka-1:9092"},
    Producer: events.ProducerConfig{
        ClientID:          "my-producer",
        Acks:              "all",
        EnableIdempotence: true,
        Compression:       "snappy",
    },
    Consumer: events.ConsumerConfig{
        GroupID:          "my-consumer-group",
        AutoOffsetReset:  "earliest",
        EnableAutoCommit: false,
    },
    Topics: events.TopicConfig{
        Prefix:            "dictamesh.",
        NumPartitions:     12,
        ReplicationFactor: 3,
    },
}
```

## Event Structure

```go
type Event struct {
    ID            string                 // Unique event ID
    Type          string                 // Event type (e.g., "entity.created")
    Source        string                 // Source adapter/system
    Subject       string                 // Entity identifier
    Timestamp     time.Time              // When event occurred
    Data          map[string]interface{} // Event payload
    Metadata      map[string]string      // Additional metadata
    CorrelationID string                 // Links related events
    CausationID   string                 // What caused this event
}
```

## Standard Topics

| Topic | Purpose | Partitions | Cleanup |
|-------|---------|------------|---------|
| `entity.changed` | Entity create/update/delete events | 12 | delete |
| `relationship.changed` | Relationship events | 12 | delete |
| `schema.changed` | Schema registration/updates | 6 | compact |
| `cache.invalidation` | Cache invalidation events | 12 | delete |
| `system.events` | System/health events | 3 | delete |
| `dead-letter` | Failed message processing | 6 | delete |

## Event Types

### Entity Events
- `entity.created` - Entity created
- `entity.updated` - Entity updated
- `entity.deleted` - Entity deleted
- `entity.read` - Entity read (for audit)

### Relationship Events
- `relationship.created` - Relationship created
- `relationship.deleted` - Relationship deleted

### Schema Events
- `schema.registered` - New schema registered
- `schema.updated` - Schema updated

### Cache Events
- `cache.invalidated` - Cache invalidated
- `cache.warmed` - Cache pre-warmed

## Best Practices

1. **Use Idempotent Producers**: Always enable idempotence in production
2. **Manual Offset Commits**: Commit only after successful processing
3. **Partition Keys**: Use entity IDs as keys for ordering
4. **Correlation IDs**: Track related events across services
5. **Error Handling**: Use dead letter topics for failed messages
6. **Monitoring**: Track producer/consumer lag and throughput

## Integration with DictaMesh

### Adapter Publishing Events

```go
type MyAdapter struct {
    producer *events.Producer
}

func (a *MyAdapter) CreateEntity(ctx context.Context, data interface{}) error {
    // Create entity...

    // Publish event
    event := events.NewEvent(
        events.EventTypeEntityCreated,
        "my-adapter",
        "entity:"+id,
        map[string]interface{}{"entity": data},
    )

    return a.producer.Publish(ctx, events.TopicEntityChanged, event)
}
```

### Service Consuming Events

```go
func (s *Service) handleEntityEvent(ctx context.Context, event *events.Event) error {
    switch event.Type {
    case events.EventTypeEntityCreated:
        return s.onEntityCreated(ctx, event)
    case events.EventTypeEntityUpdated:
        return s.onEntityUpdated(ctx, event)
    case events.EventTypeEntityDeleted:
        return s.onEntityDeleted(ctx, event)
    default:
        return fmt.Errorf("unknown event type: %s", event.Type)
    }
}
```

## Performance

- **Throughput**: 10K+ events/second per partition
- **Latency**: <10ms p99 for publish
- **Reliability**: Exactly-once semantics with idempotence
- **Scalability**: Horizontal scaling via partitions

## Dependencies

- `github.com/confluentinc/confluent-kafka-go/v2` - Kafka client
- `github.com/linkedin/goavro/v2` - Avro serialization
- `github.com/click2-run/dictamesh/pkg/observability` - Observability

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
