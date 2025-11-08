# Sentry Integration Implementation Plan

**Status**: ✅ Implemented
**Date**: 2025-01-08
**Component**: Infrastructure - Error Tracking & Monitoring

## Overview

This document describes the implementation of Sentry (self-hosted, open source) integration into the DictaMesh framework infrastructure. Sentry provides comprehensive error tracking, performance monitoring, and application observability.

## Objectives

1. ✅ Deploy self-hosted Sentry for error tracking
2. ✅ Integrate Sentry with existing infrastructure (Redpanda, PostgreSQL, Redis)
3. ✅ Support local development via Docker Compose
4. ✅ Support production deployment via Kubernetes
5. ✅ Provide developer documentation for SDK integration

## Architecture

### Components Added

#### Docker Compose (Local Development)

**New Services:**
- `sentry-postgres`: Dedicated PostgreSQL database for Sentry metadata
- `sentry-redis`: Dedicated Redis instance for Sentry cache and queues
- `clickhouse`: Event storage database (ClickHouse)
- `sentry-web`: Main Sentry web UI and API
- `sentry-worker`: Background task processor (async jobs)
- `sentry-cron`: Scheduled task executor
- `sentry-post-process-forwarder`: Event processing pipeline

**Resource Allocation:**
- Total additional RAM: ~4GB
- Total additional CPU: ~3-4 cores
- Total additional Disk: ~10GB

#### Kubernetes Deployment

**Manifests Created:**
- Base configurations (namespace, configmap, secrets, deployments, services)
- Development overlay (reduced resources, local ingress)
- Production overlay (HA setup, scaled replicas, production ingress)

**Key Features:**
- PersistentVolumeClaims for all stateful components
- Health checks and readiness probes
- Resource limits and requests
- Ingress with TLS support (production)
- Horizontal scaling support

### Integration Points

1. **Redpanda (Kafka)**: Used by Sentry for event streaming
2. **PostgreSQL**: Separate database for Sentry metadata
3. **Redis**: Separate instance for Sentry caching
4. **ClickHouse**: Dedicated event storage
5. **Existing observability**: Complements Prometheus, Grafana, and Jaeger

## Implementation Details

### File Structure

```
dictamesh/
├── infrastructure/
│   ├── docker-compose/
│   │   ├── docker-compose.dev.yml          # Updated with Sentry services
│   │   └── sentry/
│   │       ├── config/
│   │       │   ├── sentry.conf.py          # Sentry Python configuration
│   │       │   └── config.yml              # Sentry YAML configuration
│   │       ├── clickhouse/
│   │       │   └── 01-init-sentry.sql      # ClickHouse initialization
│   │       ├── init-sentry.sh              # Sentry initialization script
│   │       └── README.md                   # Configuration documentation
│   ├── k8s/
│   │   └── sentry/
│   │       ├── base/                       # Base Kubernetes manifests
│   │       │   ├── namespace.yaml
│   │       │   ├── configmap.yaml
│   │       │   ├── secret.yaml
│   │       │   ├── postgres.yaml
│   │       │   ├── redis.yaml
│   │       │   ├── clickhouse.yaml
│   │       │   ├── sentry-web.yaml
│   │       │   ├── sentry-worker.yaml
│   │       │   ├── ingress.yaml
│   │       │   └── kustomization.yaml
│   │       ├── dev/                        # Development overlay
│   │       │   └── kustomization.yaml
│   │       ├── prod/                       # Production overlay
│   │       │   └── kustomization.yaml
│   │       └── README.md                   # Kubernetes deployment guide
│   ├── Makefile                            # Updated with Sentry commands
│   └── README.md                           # Updated with Sentry info
└── docs/
    ├── SENTRY-INTEGRATION.md               # Developer integration guide
    └── planning/
        └── SENTRY-INTEGRATION-PLAN.md      # This document
```

### Configuration

#### Environment Variables

**Sentry Core:**
- `SENTRY_SECRET_KEY`: Cryptographic signing key (change in production!)
- `SENTRY_POSTGRES_HOST`: PostgreSQL hostname
- `SENTRY_DB_NAME`, `SENTRY_DB_USER`, `SENTRY_DB_PASSWORD`: Database credentials
- `SENTRY_REDIS_HOST`, `SENTRY_REDIS_PORT`: Redis connection
- `SENTRY_KAFKA_HOSTS`: Kafka/Redpanda brokers
- `SENTRY_CLICKHOUSE_HOST`: ClickHouse hostname

**Sentry Options:**
- `SENTRY_SINGLE_ORGANIZATION`: Single vs multi-org mode
- `SENTRY_BEACON`: Telemetry beacon (disabled for self-hosted)
- `SENTRY_METRICS_SAMPLE_RATE`: Metrics sampling rate
- `SENTRY_PROFILES_SAMPLE_RATE`: Profiling sampling rate
- `SENTRY_EVENT_RETENTION_DAYS`: Event retention period (default: 90 days)

#### Ports

- `9000`: Sentry Web UI and API
- `5432`: Sentry PostgreSQL (internal)
- `6379`: Sentry Redis (internal)
- `8123`: ClickHouse HTTP (internal)

### Makefile Commands

New commands added to `infrastructure/Makefile`:

```bash
make sentry              # Open Sentry UI in browser
make sentry-init         # Initialize Sentry (first-time setup)
make sentry-logs         # Follow Sentry logs
make sentry-shell        # Open Sentry shell
make sentry-cli          # Access Sentry CLI
```

Updated commands:
```bash
make dev-up             # Now includes Sentry services
make health             # Now checks Sentry health
make clean-volumes      # Now removes Sentry volumes
```

## Usage Guide

### Local Development Setup

1. **Start infrastructure**:
   ```bash
   cd infrastructure
   make dev-up
   ```

2. **Initialize Sentry** (first-time only):
   ```bash
   make sentry-init
   ```

3. **Access Sentry**:
   - URL: http://localhost:9000
   - Email: `admin@dictamesh.local`
   - Password: `admin`

4. **Create a project** for your framework component

5. **Get DSN** from project settings

6. **Integrate SDK** (see docs/SENTRY-INTEGRATION.md)

### Kubernetes Deployment

#### Development

```bash
kubectl apply -k infrastructure/k8s/sentry/dev/
```

#### Production

```bash
# Update secrets first!
kubectl create secret generic sentry-secrets \
  --from-literal=SENTRY_SECRET_KEY="$(python3 -c 'import secrets; print(secrets.token_urlsafe(50))')" \
  --from-literal=POSTGRES_PASSWORD="your-strong-password" \
  --from-literal=SENTRY_DB_PASSWORD="your-strong-password" \
  --from-literal=CLICKHOUSE_PASSWORD="your-strong-password" \
  --namespace dictamesh-sentry

kubectl apply -k infrastructure/k8s/sentry/prod/
```

## SDK Integration

### Supported Languages

Comprehensive integration guides provided for:
- ✅ Go (with examples for net/http, echo, gin)
- ✅ Node.js/TypeScript (with Express middleware)
- ✅ Python (with Flask integration)

### Example: Go Integration

```go
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.ClientOptions{
    Dsn: "http://your-dsn@localhost:9000/1",
    Environment: "development",
    Release: "dictamesh-adapter@1.0.0",
    EnableTracing: true,
    TracesSampleRate: 1.0,
})
defer sentry.Flush(2 * time.Second)
```

See `docs/SENTRY-INTEGRATION.md` for complete integration examples.

## Resource Requirements

### Development Environment

With Sentry added:
- **RAM**: ~6-7 GB (increased from ~2.5 GB)
- **CPU**: ~4-6 cores (increased from ~2-3 cores)
- **Disk**: ~15 GB (increased from ~5 GB)

### Production Environment

Recommended minimum:
- **RAM**: 16 GB
- **CPU**: 8 cores
- **Disk**: 100+ GB (depends on retention and traffic)

## Security Considerations

### Secrets Management

**Development:**
- Default passwords provided for convenience
- Clearly marked as insecure

**Production:**
- ✅ Change `SENTRY_SECRET_KEY` (generate with `secrets.token_urlsafe(50)`)
- ✅ Use strong, unique passwords for all databases
- ✅ Store secrets in Kubernetes Secrets or external secret manager
- ✅ Enable TLS/SSL for all connections
- ✅ Configure proper network policies

### Data Privacy

- ✅ `send_default_pii: false` in production
- ✅ Filter sensitive data before sending to Sentry
- ✅ Use data scrubbing rules
- ✅ Configure appropriate retention policies

## Monitoring and Operations

### Health Checks

All Sentry services include:
- Liveness probes
- Readiness probes
- Healthcheck endpoints

**Check health:**
```bash
make health                           # Docker Compose
curl http://localhost:9000/_health/   # Direct check
```

### Logging

**View logs:**
```bash
make sentry-logs                     # All Sentry services
docker logs dictamesh-sentry-web     # Specific service
```

### Backup and Recovery

**Important data to backup:**
1. PostgreSQL database (Sentry metadata)
2. ClickHouse data (event storage)
3. Sentry file storage (attachments, etc.)

**Backup commands:**
```bash
# PostgreSQL backup
docker exec dictamesh-sentry-postgres pg_dump -U sentry sentry > sentry-backup.sql

# Volume snapshots (recommended)
# Use Docker volume backup or cloud provider snapshots
```

## Scaling Considerations

### Horizontal Scaling

**Scalable components:**
- Sentry Web (multiple replicas behind load balancer)
- Sentry Worker (scale based on queue depth)

**Single-instance components:**
- Sentry Cron (single instance to avoid duplicate jobs)
- PostgreSQL (requires HA setup for scaling)
- Redis (can use Sentinel/Cluster for HA)

### Performance Tuning

**Sample Rates:**
- Development: 100% (1.0)
- Staging: 50% (0.5)
- Production: 10-25% (0.1-0.25)

**ClickHouse:**
- Configure appropriate compression
- Set up proper retention policies
- Monitor disk usage

## Testing

### Smoke Tests

```bash
# 1. Start services
make dev-up

# 2. Wait for services to be ready
make health

# 3. Initialize Sentry
make sentry-init

# 4. Check Sentry UI is accessible
curl -f http://localhost:9000/_health/

# 5. Test SDK integration (see docs/SENTRY-INTEGRATION.md)
```

### Validation Checklist

- ✅ All services start successfully
- ✅ Health checks pass
- ✅ Sentry UI is accessible
- ✅ Can create projects
- ✅ Can send test events
- ✅ Events appear in UI
- ✅ Performance monitoring works
- ✅ Kubernetes manifests apply successfully

## Future Enhancements

### Potential Improvements

1. **Helm Chart**: Package Sentry deployment as Helm chart
2. **Auto-scaling**: Configure HPA for Sentry web/worker
3. **High Availability**:
   - PostgreSQL with Patroni or CloudNativePG
   - Redis Cluster/Sentinel
   - ClickHouse replication
4. **Advanced Features**:
   - Session replay
   - Source maps upload automation
   - Release automation with CI/CD
5. **Integration**:
   - GitHub integration for issue creation
   - Slack notifications
   - PagerDuty alerts

### Monitoring Improvements

1. Grafana dashboards for Sentry metrics
2. Prometheus alerts for Sentry health
3. Log aggregation for Sentry services

## Documentation

### Created Documentation

1. **`infrastructure/docker-compose/sentry/README.md`**
   - Configuration guide
   - First-time setup
   - Environment variables
   - Production considerations

2. **`infrastructure/k8s/sentry/README.md`**
   - Kubernetes deployment guide
   - Scaling instructions
   - Troubleshooting
   - Production checklist

3. **`docs/SENTRY-INTEGRATION.md`**
   - Comprehensive SDK integration guide
   - Language-specific examples (Go, Node.js, Python)
   - Best practices
   - Framework integration patterns

4. **`infrastructure/README.md`** (updated)
   - Added Sentry to services table
   - Added Sentry commands
   - Updated resource requirements

5. **`docs/planning/SENTRY-INTEGRATION-PLAN.md`** (this document)
   - Implementation overview
   - Architecture decisions
   - Usage guide

## Troubleshooting

Common issues and solutions documented in:
- `infrastructure/README.md` (Operations)
- `infrastructure/k8s/sentry/README.md` (Kubernetes)
- `docs/SENTRY-INTEGRATION.md` (SDK integration)

## Success Criteria

All objectives achieved:
- ✅ Self-hosted Sentry deployed and operational
- ✅ Integrated with existing infrastructure (Redpanda, PostgreSQL, Redis)
- ✅ Docker Compose support for local development
- ✅ Kubernetes manifests for dev and production
- ✅ Comprehensive developer documentation
- ✅ Makefile automation
- ✅ Security considerations addressed
- ✅ Resource requirements documented
- ✅ Troubleshooting guides provided

## References

- [Sentry Self-Hosted Documentation](https://develop.sentry.dev/self-hosted/)
- [Sentry Docker Repository](https://github.com/getsentry/self-hosted)
- [Sentry SDK Documentation](https://docs.sentry.io/platforms/)
- [ClickHouse Documentation](https://clickhouse.com/docs)

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
