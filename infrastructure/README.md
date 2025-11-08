# DictaMesh Framework Infrastructure

Development and deployment infrastructure for the DictaMesh data mesh adapter framework.

## Overview

This directory contains infrastructure setup for developing and deploying the DictaMesh framework components:

- **Docker Compose**: Local development environment with all framework dependencies
- **Kubernetes/K3s**: Production-ready deployment manifests
- **Scripts**: Automation and utility scripts
- **Terraform**: Infrastructure as Code (optional, for cloud deployments)

## Quick Start - Development Environment

### Prerequisites

- Docker 24.0+ and Docker Compose 2.0+
- Make (optional, for convenience commands)
- 4GB RAM minimum, 8GB recommended

### Start Infrastructure

```bash
# Start all services
make dev-up

# Or without make:
docker-compose -f docker-compose/docker-compose.dev.yml up -d
```

### Access Services

Once running, access the following services:

| Service | URL | Credentials |
|---------|-----|-------------|
| **Redpanda Console** (Kafka UI) | http://localhost:8080 | - |
| **Grafana** | http://localhost:3000 | admin / admin |
| **Prometheus** | http://localhost:9090 | - |
| **Jaeger** (Tracing) | http://localhost:16686 | - |
| **Sentry** (Error Tracking) | http://localhost:9000 | admin@dictamesh.local / admin |
| **PostgreSQL** | localhost:5432 | dictamesh / dictamesh_dev_password |
| **Redis** | localhost:6379 | - |
| **Redpanda Kafka API** | localhost:19092 | - |
| **Schema Registry** | localhost:18081 | - |

### Useful Commands

```bash
# View logs
make dev-logs

# Check service status
make dev-ps

# Check health
make health

# Connect to PostgreSQL
make postgres-cli

# Connect to Redis
make redis-cli

# List Kafka topics
make kafka-topics

# Initialize Sentry (first-time only)
make sentry-init

# Open Sentry UI
make sentry

# View Sentry logs
make sentry-logs

# Stop services
make dev-down

# Reset everything (removes data!)
make dev-reset
```

## Architecture Components

### Event Bus - Redpanda

**Why Redpanda instead of Kafka?**
- Kafka-compatible API (drop-in replacement)
- **Much lighter**: ~500MB RAM vs Kafka's 2-4GB
- No JVM overhead
- Built-in Schema Registry
- Perfect for development and small deployments

**Configuration**:
- Kafka API: `localhost:19092`
- Schema Registry: `localhost:18081`
- REST API: `localhost:18082`
- Metrics: `localhost:19644`

### Metadata Catalog - PostgreSQL

The metadata catalog stores:
- Entity registry (all entities across integrated sources)
- Relationship graph (cross-system links)
- Schema registry (versioned entity schemas)
- Event log (immutable audit trail)
- Data lineage tracking

**Database**: `metadata_catalog`
**Schema**: Auto-initialized on first startup (see `init-scripts/postgres/`)

**Tables**:
- `entity_catalog` - Entity registry
- `entity_relationships` - Relationship graph
- `schemas` - Schema versions
- `event_log` - Audit trail
- `data_lineage` - Data flow tracking
- `cache_status` - Cache metadata

### Caching Layer - Redis

Multi-layer caching strategy:
- **L1**: In-memory (adapter-local)
- **L2**: Redis (shared across replicas) ← This service
- **L3**: PostgreSQL (metadata catalog)

**Configuration**:
- Max memory: 256MB
- Eviction policy: `allkeys-lru`

### Observability Stack

#### Prometheus
- Metrics storage and querying
- Retention: 30 days
- Scrapes: Redpanda, framework services

#### Grafana
- Metrics visualization
- Pre-configured datasources (Prometheus, Jaeger)
- Dashboards auto-provisioned (coming soon)

#### Jaeger
- Distributed tracing
- All-in-one deployment for development
- Storage: In-memory (for dev)

### Error Tracking and Monitoring - Sentry

**Self-Hosted Sentry** for comprehensive error tracking and application monitoring:
- Error and exception tracking
- Performance monitoring (APM)
- Release tracking
- User feedback collection
- Issue alerts and notifications

**Components**:
- **Sentry Web**: Main UI and API (port 9000)
- **Sentry Worker**: Background task processing
- **Sentry Cron**: Scheduled tasks
- **PostgreSQL**: Sentry database
- **Redis**: Cache and message broker
- **ClickHouse**: Event storage (optional but recommended)

**First-Time Setup**:
```bash
# After starting infrastructure
make sentry-init
```

**Configuration**:
See `docker-compose/sentry/README.md` for detailed configuration options.

**Integration**:
Framework components can integrate Sentry for error tracking. See the Sentry documentation for SDK integration guides.

## Resource Requirements

### Development Environment

Total resource usage with all services (including Sentry):

| Resource | Usage |
|----------|-------|
| **RAM** | ~6-7 GB |
| **CPU** | ~4-6 cores |
| **Disk** | ~15 GB (with data volumes) |

Per-service limits (configured in docker-compose):

| Service | RAM Limit | CPU Limit |
|---------|-----------|-----------|
| Redpanda | 1 GB | 1 core |
| PostgreSQL | 512 MB | 1 core |
| Redis | 256 MB | 0.5 core |
| Prometheus | 512 MB | 0.5 core |
| Grafana | 256 MB | 0.5 core |
| Jaeger | 256 MB | 0.5 core |
| Redpanda Console | 256 MB | 0.5 core |
| Sentry PostgreSQL | 512 MB | 1 core |
| Sentry Redis | 256 MB | 0.5 core |
| ClickHouse | 1 GB | 1 core |
| Sentry Web | 1 GB | 1 core |
| Sentry Worker | 1 GB | 1 core |
| Sentry Cron | 256 MB | 0.5 core |
| Sentry Post-Process | 512 MB | 0.5 core |

## Production Deployment

### Kubernetes/K3s

For production deployments on Kubernetes:

```bash
# See k8s/ directory for manifests
cd k8s/

# Apply base infrastructure
kubectl apply -k base/

# Apply environment-specific config
kubectl apply -k dev/
# or
kubectl apply -k prod/

# Deploy Sentry to Kubernetes
kubectl apply -k sentry/dev/
# or
kubectl apply -k sentry/prod/
```

See `k8s/sentry/README.md` for detailed Kubernetes deployment instructions.

### Helm Charts (Coming Soon)

Helm charts for all framework components will be available:
- Metadata Catalog Service
- GraphQL Gateway
- Event Router
- Observability Stack

## Scaling Considerations

### Development → Production Changes

When moving to production, consider:

1. **Kafka/Redpanda**:
   - Dev: Single Redpanda instance
   - Prod: 3+ broker cluster with replication factor 3

2. **PostgreSQL**:
   - Dev: Single instance
   - Prod: HA setup with replicas (Patroni, CloudNativePG, or managed service)

3. **Redis**:
   - Dev: Single instance
   - Prod: Redis Cluster or Sentinel for HA

4. **Observability**:
   - Dev: All-in-one Jaeger, small Prometheus
   - Prod: Distributed tracing backend, larger Prometheus with federation

## Environment Variables

Configure services via environment variables in `docker-compose.dev.yml`:

```bash
# PostgreSQL
POSTGRES_USER=dictamesh
POSTGRES_PASSWORD=dictamesh_dev_password
POSTGRES_DB=metadata_catalog

# Grafana
GF_SECURITY_ADMIN_USER=admin
GF_SECURITY_ADMIN_PASSWORD=admin
```

For production, use Kubernetes Secrets or a secret management solution.

## Troubleshooting

### Services won't start

```bash
# Check Docker resources
docker info

# Check logs
make dev-logs

# Restart specific service
docker-compose -f docker-compose/docker-compose.dev.yml restart <service>
```

### Redpanda issues

```bash
# Check cluster health
docker exec dictamesh-redpanda rpk cluster health

# Check cluster info
docker exec dictamesh-redpanda rpk cluster info

# Create test topic
docker exec dictamesh-redpanda rpk topic create test-topic
```

### PostgreSQL issues

```bash
# Check if database is ready
docker exec dictamesh-postgres pg_isready -U dictamesh

# Connect to database
make postgres-cli

# Check tables
\dt

# Check entity_catalog
SELECT * FROM entity_catalog;
```

### Redis issues

```bash
# Connect to Redis
make redis-cli

# Check info
INFO

# Monitor commands
MONITOR
```

### Sentry issues

```bash
# Check Sentry web logs
docker logs dictamesh-sentry-web

# Check Sentry worker logs
docker logs dictamesh-sentry-worker

# Test Sentry health endpoint
curl http://localhost:9000/_health/

# Re-initialize Sentry
make sentry-init

# Access Sentry shell for debugging
make sentry-shell
```

## Development Workflow

1. **Start infrastructure**: `make dev-up`
2. **Develop framework components**: Code in `pkg/`, `services/`
3. **Run tests**: Tests connect to local infrastructure
4. **View metrics**: Check Grafana dashboards
5. **View traces**: Check Jaeger UI
6. **Debug**: Use Redpanda Console, Postgres CLI, Redis CLI
7. **Stop infrastructure**: `make dev-down`

## Next Steps

- Review framework architecture: `../docs/planning/01-ARCHITECTURE-OVERVIEW.md`
- Start building framework components: `../pkg/`
- Read development guide: `../CONTRIBUTING.md` (coming soon)

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
