# Loki Log Aggregation Setup

This directory contains the Loki configuration for DictaMesh development environment.

## Overview

Loki is a horizontally-scalable, highly-available log aggregation system inspired by Prometheus. It indexes metadata (labels) rather than full-text, making it extremely efficient for log storage and retrieval.

## Components

- **Loki**: Log aggregation server (port 3100)
- **Promtail**: Log collection agent (scrapes Docker container logs)
- **Grafana**: Log visualization and querying interface (port 3000)

## Quick Start

### Start the Stack

```bash
# From the infrastructure/docker-compose directory
docker-compose -f docker-compose.dev.yml up -d loki promtail grafana
```

### Verify Services

```bash
# Check Loki is running
curl http://localhost:3100/ready
# Expected: ready

# Check Promtail is running
curl http://localhost:9080/ready
# Expected: ready

# Check logs are being ingested
curl http://localhost:3100/loki/api/v1/label/service/values
# Expected: ["metadata-catalog", "postgres", "redis", ...]
```

### Access Logs

1. **Grafana UI**: http://localhost:3000
   - Username: `admin`
   - Password: `admin`
   - Navigate to Explore → Select "Loki" datasource

2. **Pre-built Dashboard**: http://localhost:3000/d/dictamesh-logs-overview
   - View all DictaMesh logs
   - Filter by service, level, or search text
   - See error/warning counts

## Configuration

### Loki Configuration (`loki-config.yml`)

Key settings:
- **Retention**: 30 days (720h)
- **Storage**: Local filesystem (`/loki/chunks`)
- **Indexing**: BoltDB with daily rotation
- **Ingestion Limit**: 10MB/s with 20MB bursts

### Promtail Configuration (`promtail-config.yml`)

Key features:
- **Auto-discovery**: Automatically discovers all `dictamesh-*` containers
- **JSON Parsing**: Parses Zap-formatted JSON logs
- **Label Extraction**: Adds `service`, `level`, `logger`, `trace_id` labels
- **Metric Generation**: Counts errors per service

## Querying Logs

### Basic Queries

```logql
# View all logs from metadata-catalog service
{service="metadata-catalog"}

# Show only errors
{service="metadata-catalog"} | json | level="error"

# Search for specific text
{job="dictamesh-logs"} |= "database connection"

# Filter by HTTP status
{service="metadata-catalog"} | json | status >= 500
```

### Advanced Queries

```logql
# Find slow requests (>1 second)
{job="dictamesh-logs"} | json | duration_ms > 1000

# Track a distributed trace
{job="dictamesh-logs"} | json | trace_id="abc123..."

# Error rate by service
sum by (service) (rate({job="dictamesh-logs"} | json | level="error" [5m]))
```

See [Loki Query Runbook](../../../docs/operations/LOKI-QUERY-RUNBOOK.md) for more examples.

## Troubleshooting

### No Logs Appearing

1. **Check Promtail is running**:
   ```bash
   docker logs dictamesh-promtail
   ```

2. **Check Promtail can access Docker socket**:
   ```bash
   docker exec dictamesh-promtail ls -la /var/run/docker.sock
   ```

3. **Check Loki is receiving logs**:
   ```bash
   curl http://localhost:3100/metrics | grep loki_distributor_bytes_received_total
   ```

### Query Timeout

- Reduce time range
- Add more specific label filters
- Avoid high-cardinality fields

### High Memory Usage

Loki is configured with:
- 512MB memory limit
- 100MB result cache
- 10MB ingestion rate

If you need more, adjust in `docker-compose.dev.yml`:
```yaml
loki:
  deploy:
    resources:
      limits:
        memory: 1G  # Increase if needed
```

## Data Retention

Logs are retained for **30 days** by default. To change:

1. Edit `loki-config.yml`:
   ```yaml
   limits_config:
     retention_period: 720h  # Change to desired hours
   ```

2. Restart Loki:
   ```bash
   docker-compose -f docker-compose.dev.yml restart loki
   ```

## Backup and Restore

### Backup Logs

Loki data is stored in the `dictamesh-loki-data` Docker volume:

```bash
# Backup to tar file
docker run --rm -v dictamesh-loki-data:/data -v $(pwd):/backup \
  alpine tar czf /backup/loki-backup-$(date +%Y%m%d).tar.gz -C /data .
```

### Restore Logs

```bash
# Stop Loki
docker-compose -f docker-compose.dev.yml stop loki

# Restore from tar file
docker run --rm -v dictamesh-loki-data:/data -v $(pwd):/backup \
  alpine sh -c "cd /data && tar xzf /backup/loki-backup-20250108.tar.gz"

# Start Loki
docker-compose -f docker-compose.dev.yml start loki
```

## Resource Usage

Typical resource consumption:
- **Loki**: 150-300MB RAM, 5-10% CPU
- **Promtail**: 50-100MB RAM, 2-5% CPU
- **Disk**: ~1GB per day (depends on log volume)

## Integration with Existing Stack

Loki integrates with:

- **Prometheus**: Log-based metrics exported to Prometheus
- **Jaeger**: Trace IDs in logs link to Jaeger UI
- **Grafana**: Unified observability dashboard
- **Sentry**: Complementary error tracking

## Production Considerations

For production deployment:

1. **Use S3/GCS for storage** instead of local filesystem
2. **Enable authentication** (`auth_enabled: true`)
3. **Deploy Loki cluster** (3+ replicas for HA)
4. **Increase retention** based on compliance requirements
5. **Set up alerts** for ingestion failures
6. **Monitor Loki metrics** in Prometheus

See [Infrastructure Planning](../../../docs/planning/03-INFRASTRUCTURE-PLANNING.md) for Kubernetes deployment.

## References

- [Loki Configuration Docs](https://grafana.com/docs/loki/latest/configuration/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)
- [Best Practices](https://grafana.com/docs/loki/latest/best-practices/)

## Support

Issues or questions:
1. Check container logs: `docker logs dictamesh-loki`
2. Review [Query Runbook](../../../docs/operations/LOKI-QUERY-RUNBOOK.md)
3. Check [Loki GitHub Discussions](https://github.com/grafana/loki/discussions)
