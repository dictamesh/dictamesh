---
sidebar_position: 2
---

# Quick Start

Get up and running with DictaMesh in 10 minutes.

## Prerequisites

- Docker Desktop or Podman
- Go 1.21+
- kubectl + k3d (optional, for Kubernetes)
- Git

## Step 1: Start Local Infrastructure

DictaMesh includes a complete development environment with all dependencies.

```bash
# Clone the repository
git clone https://github.com/dictamesh/dictamesh.git
cd dictamesh

# Start infrastructure with Docker Compose
cd infrastructure
make dev-up
```

This starts:
- **PostgreSQL** - Metadata catalog database (localhost:5432)
- **Redpanda** - Kafka-compatible message broker (localhost:9092)
- **Redis** - L2 caching layer (localhost:6379)
- **Prometheus** - Metrics collection (localhost:9090)
- **Grafana** - Dashboards (localhost:3000)
- **Jaeger** - Distributed tracing (localhost:16686)

Wait for all services to be healthy:

```bash
make health
```

Expected output:
```
✓ PostgreSQL is healthy
✓ Redpanda is healthy
✓ Redis is healthy
✓ Prometheus is healthy
✓ Grafana is healthy
✓ Jaeger is healthy
```

## Step 2: Run Database Migrations

Initialize the metadata catalog schema:

```bash
cd ../services/metadata-catalog
make migrate-up
```

This creates the complete database schema for the metadata catalog.

## Step 3: Build Your First Adapter

Let's create a simple adapter for a hypothetical "Products API":

```bash
# Scaffold a new adapter
cd ../../
./scripts/scaffold-adapter.sh products-adapter products
```

This generates:
```
services/products-adapter/
├── cmd/server/main.go
├── internal/
│   ├── adapter/adapter.go
│   ├── client/client.go
│   └── transformer/transformer.go
├── Dockerfile
├── go.mod
└── README.md
```

## Step 4: Implement the Adapter

Edit `services/products-adapter/internal/adapter/adapter.go`:

```go
package adapter

import (
    "context"
    "github.com/dictamesh/dictamesh/pkg/adapter"
    "github.com/dictamesh/dictamesh/pkg/models"
)

type ProductsAdapter struct {
    client         *Client
    eventPublisher *events.Publisher
    cache          cache.Cache
}

func (a *ProductsAdapter) GetEntity(ctx context.Context, id string) (*models.Entity, error) {
    // 1. Check cache
    if cached, err := a.cache.Get(ctx, "product:"+id); err == nil {
        return cached.(*models.Entity), nil
    }

    // 2. Fetch from source
    product, err := a.client.GetProduct(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. Transform to canonical model
    entity := &models.Entity{
        ID:   product.ID,
        Type: "product",
        Data: map[string]interface{}{
            "name":        product.Name,
            "price":       product.Price,
            "description": product.Description,
        },
    }

    // 4. Cache result
    a.cache.Set(ctx, "product:"+id, entity, 5*time.Minute)

    return entity, nil
}

func (a *ProductsAdapter) GetSchema() adapter.Schema {
    return adapter.Schema{
        Entity:  "product",
        Version: "1.0.0",
        Fields: []adapter.Field{
            {Name: "id", Type: "uuid", Required: true},
            {Name: "name", Type: "string", Required: true},
            {Name: "price", Type: "decimal", Required: true},
            {Name: "description", Type: "string"},
        },
    }
}
```

## Step 5: Run the Adapter

```bash
cd services/products-adapter
go run cmd/server/main.go
```

The adapter will:
1. Connect to PostgreSQL (metadata catalog)
2. Connect to Redpanda (event bus)
3. Connect to Redis (cache)
4. Register its schema
5. Start serving requests

## Step 6: Query via GraphQL

Open the GraphQL Playground at http://localhost:8080/playground

Try this query:

```graphql
query GetProduct {
  product(id: "prod-123") {
    id
    name
    price
    description
  }
}
```

## Step 7: Explore Observability

### Metrics (Prometheus)
http://localhost:9090

Example queries:
```promql
# Adapter request rate
rate(adapter_requests_total[5m])

# Adapter latency (p99)
histogram_quantile(0.99, adapter_request_duration_seconds_bucket)

# Cache hit rate
rate(cache_hits_total[5m]) / rate(cache_requests_total[5m])
```

### Dashboards (Grafana)
http://localhost:3000

Default credentials: `admin / admin`

Explore pre-built dashboards:
- DictaMesh Overview
- Adapter Performance
- Event Bus Metrics
- Cache Performance

### Distributed Tracing (Jaeger)
http://localhost:16686

Search for traces from your adapter to see:
- Request flow through components
- Latency breakdown
- External API calls
- Database queries

## What's Next?

### Learn Core Concepts
Understand the key concepts: [Core Concepts](./core-concepts.md)

### Build a Real Adapter
Follow the complete guide: [Building Adapters](../guides/building-adapters.md)

### Deploy to Production
Set up Kubernetes: [Deployment Guide](../guides/deployment.md)

### Explore Architecture
Deep dive into design: [Architecture Overview](../architecture/overview.md)

## Common Issues

### Port Already in Use

If ports 5432, 9092, 6379, etc. are already in use:

```bash
# Check what's using the port
lsof -i :5432

# Stop conflicting services
docker stop $(docker ps -q)
```

### Services Not Starting

Check Docker resource limits:
- Memory: At least 4GB
- CPU: At least 2 cores
- Disk: At least 10GB free

### Connection Refused Errors

Ensure all services are healthy:

```bash
cd infrastructure
make health
```

If any service is unhealthy:

```bash
# View logs
docker-compose logs <service-name>

# Restart specific service
docker-compose restart <service-name>
```

## Next Steps

- 📖 Read [Core Concepts](./core-concepts.md)
- 🏗️ Follow [Building Adapters Guide](../guides/building-adapters.md)
- 🚀 Set up [Production Deployment](../guides/deployment.md)
- 💬 Join [GitHub Discussions](https://github.com/dictamesh/dictamesh/discussions)

---

**Previous**: [← Introduction](./introduction.md) | **Next**: [Core Concepts →](./core-concepts.md)
