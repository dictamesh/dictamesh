---
sidebar_position: 3
---

# Installation

This guide covers different installation methods for DictaMesh, from local development to production deployments.

## Prerequisites

### Required

- **Docker Desktop** or **Podman** 4.0+
- **Go** 1.21 or later
- **Git** for cloning the repository
- **kubectl** for Kubernetes deployments

### Recommended

- **k3d** or **kind** for local Kubernetes testing
- **Helm** 3.10+ for Kubernetes deployments
- **Make** for build automation
- **Node.js** 20+ (for documentation portal development)

## Installation Methods

### Method 1: Docker Compose (Recommended for Development)

The fastest way to get DictaMesh running locally with all dependencies.

#### Step 1: Clone the Repository

```bash
git clone https://github.com/Click2-Run/dictamesh.git
cd dictamesh
```

#### Step 2: Start Infrastructure

```bash
cd infrastructure
make dev-up
```

This starts all required services:
- PostgreSQL (metadata catalog)
- Redpanda (Kafka-compatible message broker)
- Redis (caching layer)
- Prometheus (metrics)
- Grafana (dashboards)
- Jaeger (distributed tracing)

#### Step 3: Verify Services

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

#### Step 4: Run Migrations

```bash
cd ../services/metadata-catalog
make migrate-up
```

Your development environment is now ready!

---

### Method 2: Kubernetes (Production)

Deploy DictaMesh to a Kubernetes cluster using Helm.

#### Prerequisites

- Kubernetes cluster 1.26+
- Helm 3.10+
- kubectl configured for your cluster

#### Step 1: Add Helm Repository

```bash
helm repo add dictamesh https://charts.dictamesh.controle.digital
helm repo update
```

#### Step 2: Create Namespace

```bash
kubectl create namespace dictamesh-system
```

#### Step 3: Install with Helm

```bash
helm install dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --create-namespace \
  --values values.yaml
```

#### Step 4: Verify Installation

```bash
kubectl get pods -n dictamesh-system
```

Expected output:
```
NAME                                    READY   STATUS    RESTARTS   AGE
dictamesh-metadata-catalog-0            1/1     Running   0          2m
dictamesh-graphql-gateway-xyz           1/1     Running   0          2m
dictamesh-postgres-0                    1/1     Running   0          2m
dictamesh-kafka-0                       1/1     Running   0          2m
dictamesh-redis-0                       1/1     Running   0          2m
```

---

### Method 3: Manual Installation

For custom setups or when you need more control.

#### Step 1: Install Dependencies

**PostgreSQL 15+**
```bash
# Ubuntu/Debian
sudo apt-get install postgresql-15

# macOS
brew install postgresql@15

# Start service
sudo systemctl start postgresql  # Linux
brew services start postgresql@15  # macOS
```

**Redis 7+**
```bash
# Ubuntu/Debian
sudo apt-get install redis-server

# macOS
brew install redis

# Start service
sudo systemctl start redis  # Linux
brew services start redis  # macOS
```

**Apache Kafka** (or use Redpanda)
```bash
# Using Docker
docker run -d \
  --name kafka \
  -p 9092:9092 \
  docker.redpanda.com/redpandadata/redpanda:latest \
  redpanda start \
  --kafka-addr internal://0.0.0.0:9092
```

#### Step 2: Set Up Database

```sql
-- Create database
CREATE DATABASE dictamesh_catalog;

-- Create user
CREATE USER dictamesh WITH PASSWORD 'your-secure-password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE dictamesh_catalog TO dictamesh;
```

#### Step 3: Run Migrations

```bash
cd services/metadata-catalog
export DATABASE_URL="postgresql://dictamesh:your-secure-password@localhost:5432/dictamesh_catalog"
make migrate-up
```

#### Step 4: Build Services

```bash
# Build metadata catalog
cd services/metadata-catalog
go build -o bin/metadata-catalog cmd/server/main.go

# Build GraphQL gateway
cd ../graphql-gateway
go build -o bin/graphql-gateway cmd/server/main.go
```

#### Step 5: Run Services

```bash
# Terminal 1: Metadata Catalog
cd services/metadata-catalog
./bin/metadata-catalog

# Terminal 2: GraphQL Gateway
cd services/graphql-gateway
./bin/graphql-gateway
```

---

## Configuration

### Environment Variables

Create a `.env` file with the following variables:

```bash
# Database Configuration
DATABASE_URL=postgresql://dictamesh:password@localhost:5432/dictamesh_catalog
DATABASE_MAX_CONNECTIONS=25
DATABASE_MAX_IDLE_CONNECTIONS=10

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID=dictamesh
KAFKA_CONSUMER_GROUP=dictamesh-consumers

# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_MAX_CONNECTIONS=10

# Observability
JAEGER_ENDPOINT=http://localhost:14268/api/traces
PROMETHEUS_PORT=9090
LOG_LEVEL=info

# Security
JWT_SECRET=your-jwt-secret-key-here
API_KEY=your-api-key-here
```

### Configuration Files

DictaMesh uses YAML configuration files for more complex settings.

Create `config/metadata-catalog.yaml`:

```yaml
server:
  port: 8080
  host: 0.0.0.0
  read_timeout: 30s
  write_timeout: 30s

database:
  url: ${DATABASE_URL}
  max_connections: 25
  max_idle_connections: 10
  connection_max_lifetime: 5m

kafka:
  brokers:
    - localhost:9092
  client_id: dictamesh-metadata-catalog
  consumer_group: metadata-catalog-consumers

cache:
  redis:
    url: ${REDIS_URL}
    ttl: 5m

observability:
  tracing:
    enabled: true
    jaeger_endpoint: ${JAEGER_ENDPOINT}
  metrics:
    enabled: true
    port: 9090
  logging:
    level: info
    format: json
```

---

## Verification

### Check Service Health

```bash
# Metadata Catalog
curl http://localhost:8080/health

# GraphQL Gateway
curl http://localhost:8000/health
```

### Test GraphQL API

Open http://localhost:8000/playground and run:

```graphql
query {
  health {
    status
    version
    timestamp
  }
}
```

### View Metrics

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686

---

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Database Connection Failed

```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U dictamesh -d dictamesh_catalog -h localhost
```

### Kafka Connection Failed

```bash
# Check Kafka is running
docker ps | grep kafka

# View Kafka logs
docker logs kafka
```

### Services Won't Start

Check logs for specific errors:

```bash
# Docker Compose
docker-compose logs -f

# Kubernetes
kubectl logs -f deployment/metadata-catalog -n dictamesh-system

# Manual installation
tail -f logs/metadata-catalog.log
```

---

## Next Steps

- 📖 Learn [Core Concepts](./core-concepts.md)
- 🏗️ Build your first adapter: [Building Adapters Guide](../guides/building-adapters.md)
- 🚀 Deploy to production: [Deployment Guide](../guides/deployment.md)

---

**Previous**: [← Quick Start](./quickstart.md) | **Next**: [Core Concepts →](./core-concepts.md)
