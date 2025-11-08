# Layer 1: Source System Adapters Implementation

[← Previous: CI/CD Pipeline](05-CICD-PIPELINE.md) | [Next: Layer 2 Event Bus →](07-LAYER2-EVENT-BUS.md)

---

## 🎯 Purpose

Detailed implementation guide for creating Data Product Adapters that transform heterogeneous source systems into standardized interfaces.

**Reading Time:** 30 minutes
**Prerequisites:** [Architecture Overview](01-ARCHITECTURE-OVERVIEW.md), [Infrastructure Planning](03-INFRASTRUCTURE-PLANNING.md)

---

## 📐 Adapter Architecture

### Data Product Adapter Interface

All adapters must implement this standard interface:

```go
// pkg/adapter/interface.go
package adapter

type DataProductAdapter interface {
    // Core CRUD operations
    GetEntity(ctx context.Context, id string) (*Entity, error)
    QueryEntities(ctx context.Context, query Query) ([]Entity, error)
    
    // Metadata
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

---

## 🔨 Implementation: Customer Adapter (Directus)

### Project Structure

```
services/customer-adapter/
├── cmd/server/main.go
├── internal/
│   ├── adapter/
│   │   ├── adapter.go           # Main adapter implementation
│   │   ├── directus_client.go   # Directus API client
│   │   ├── transformer.go       # Data transformation
│   │   └── circuit_breaker.go   # Resilience
│   ├── events/
│   │   ├── publisher.go         # Kafka event publisher
│   │   └── schema.go            # Avro schema handling
│   ├── cache/
│   │   ├── multi_layer.go       # L1/L2/L3 cache
│   │   └── strategies.go        # Cache invalidation
│   └── server/
│       ├── handlers.go          # HTTP handlers
│       └── middleware.go        # Auth, logging, tracing
├── pkg/models/
│   └── entity.go                # Canonical entity models
├── deployments/k8s/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
├── Dockerfile
├── go.mod
└── README.md
```

### Step 1: Directus Client Implementation

```go
// internal/adapter/directus_client.go
package adapter

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

type DirectusClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

func NewDirectusClient(baseURL, apiKey string) *DirectusClient {
    return &DirectusClient{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *DirectusClient) GetCustomer(ctx context.Context, id string) (*DirectusCustomer, error) {
    url := fmt.Sprintf("%s/items/customers/%s", c.baseURL, id)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("directus request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("directus returned status %d", resp.StatusCode)
    }
    
    var result DirectusCustomer
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

### Step 2: Main Adapter Implementation

```go
// internal/adapter/adapter.go
package adapter

import (
    "context"
    "github.com/dictamesh/customer-adapter/pkg/models"
    "github.com/sony/gobreaker"
)

type CustomerAdapter struct {
    directusClient *DirectusClient
    eventPublisher *EventPublisher
    cache          *MultiLayerCache
    circuitBreaker *gobreaker.CircuitBreaker
    metrics        *Metrics
}

func NewCustomerAdapter(cfg Config) *CustomerAdapter {
    cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "directus",
        MaxRequests: 3,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            return counts.ConsecutiveFailures > 5
        },
    })
    
    return &CustomerAdapter{
        directusClient: NewDirectusClient(cfg.DirectusURL, cfg.APIKey),
        eventPublisher: NewEventPublisher(cfg.KafkaBrokers),
        cache:          NewMultiLayerCache(cfg.RedisURL),
        circuitBreaker: cb,
        metrics:        NewMetrics(),
    }
}

func (a *CustomerAdapter) GetEntity(ctx context.Context, id string) (*models.Entity, error) {
    span, ctx := tracer.Start(ctx, "CustomerAdapter.GetEntity")
    defer span.End()
    
    // L1 Cache: In-memory
    if cached, ok := a.cache.GetFromL1(id); ok {
        a.metrics.CacheHits.WithLabelValues("l1").Inc()
        return cached, nil
    }
    
    // L2 Cache: Redis
    if cached, err := a.cache.GetFromL2(ctx, id); err == nil {
        a.metrics.CacheHits.WithLabelValues("l2").Inc()
        a.cache.SetL1(id, cached)
        return cached, nil
    }
    
    // Circuit breaker protection
    result, err := a.circuitBreaker.Execute(func() (interface{}, error) {
        return a.fetchFromSource(ctx, id)
    })
    
    if err != nil {
        a.metrics.SourceErrors.Inc()
        return nil, err
    }
    
    entity := result.(*models.Entity)
    
    // Cache the result
    a.cache.SetL1(id, entity)
    a.cache.SetL2(ctx, id, entity, 5*time.Minute)
    
    return entity, nil
}

func (a *CustomerAdapter) fetchFromSource(ctx context.Context, id string) (*models.Entity, error) {
    directusCustomer, err := a.directusClient.GetCustomer(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Transform to canonical model
    entity := a.transformToEntity(directusCustomer)
    
    return entity, nil
}

func (a *CustomerAdapter) StreamChanges(ctx context.Context) (<-chan models.ChangeEvent, error) {
    changeChan := make(chan models.ChangeEvent, 100)
    
    // Webhook listener (if Directus supports webhooks)
    go a.listenWebhooks(ctx, changeChan)
    
    // Polling fallback
    go a.pollChanges(ctx, changeChan)
    
    return changeChan, nil
}

func (a *CustomerAdapter) GetSchema() Schema {
    return Schema{
        Entity:  "customer",
        Version: "1.0.0",
        Fields: []Field{
            {Name: "id", Type: "uuid", Required: true},
            {Name: "email", Type: "string", Required: true, PII: true},
            {Name: "name", Type: "string", Required: true, PII: true},
        },
        SLA: ServiceLevelAgreement{
            Availability: 0.999,
            LatencyP99:   500 * time.Millisecond,
            Freshness:    5 * time.Second,
        },
    }
}
```

### Step 3: Event Publishing

```go
// internal/events/publisher.go
package events

import (
    "context"
    "github.com/segmentio/kafka-go"
)

type EventPublisher struct {
    writer *kafka.Writer
}

func NewEventPublisher(brokers []string) *EventPublisher {
    return &EventPublisher{
        writer: &kafka.Writer{
            Addr:                   kafka.TCP(brokers...),
            Topic:                  "customers.directus.entity_changed",
            Balancer:               &kafka.Hash{},
            RequiredAcks:           kafka.RequireAll,
            AllowAutoTopicCreation: false,
        },
    }
}

func (p *EventPublisher) PublishChange(ctx context.Context, event ChangeEvent) error {
    value, err := avro.Marshal(event)
    if err != nil {
        return err
    }
    
    return p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(event.EntityID),
        Value: value,
    })
}
```

### Step 4: Multi-Layer Caching

```go
// internal/cache/multi_layer.go
package cache

import (
    "context"
    "github.com/go-redis/redis/v8"
    "sync"
    "time"
)

type MultiLayerCache struct {
    l1     sync.Map // In-memory
    l2     *redis.Client
    l3     *sql.DB // Metadata catalog
}

func (c *MultiLayerCache) GetFromL1(key string) (*Entity, bool) {
    if val, ok := c.l1.Load(key); ok {
        return val.(*Entity), true
    }
    return nil, false
}

func (c *MultiLayerCache) GetFromL2(ctx context.Context, key string) (*Entity, error) {
    val, err := c.l2.Get(ctx, "customer:"+key).Result()
    if err == redis.Nil {
        return nil, ErrCacheMiss
    }
    if err != nil {
        return nil, err
    }
    
    var entity Entity
    json.Unmarshal([]byte(val), &entity)
    return &entity, nil
}

func (c *MultiLayerCache) SetL1(key string, entity *Entity) {
    c.l1.Store(key, entity)
}

func (c *MultiLayerCache) SetL2(ctx context.Context, key string, entity *Entity, ttl time.Duration) error {
    data, _ := json.Marshal(entity)
    return c.l2.Set(ctx, "customer:"+key, data, ttl).Err()
}
```

---

## 🚀 Deployment

See [03-INFRASTRUCTURE-PLANNING.md](03-INFRASTRUCTURE-PLANNING.md) and [04-DEPLOYMENT-STRATEGY.md](04-DEPLOYMENT-STRATEGY.md) for deployment procedures.

---

[← Previous: CI/CD Pipeline](05-CICD-PIPELINE.md) | [Next: Layer 2 Event Bus →](07-LAYER2-EVENT-BUS.md)
