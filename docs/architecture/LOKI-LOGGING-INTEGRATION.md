# Loki Logging Integration

**Document Version:** 1.0
**Last Updated:** 2025-11-08
**Status:** Implementation in Progress

## Executive Summary

This document describes the integration of Grafana Loki as the centralized log aggregation system for DictaMesh. Loki provides efficient log storage, indexing, and querying capabilities that complement our existing observability stack (Prometheus for metrics, Jaeger for traces, Sentry for errors).

## Problem Statement

**Current State:**
- Logs are output to stdout/stderr using Uber Zap (Go services)
- No centralized log aggregation - logs are ephemeral
- Logs are lost when containers are stopped or restarted
- No efficient way to search logs across services
- No correlation between logs, metrics, and traces in Grafana
- Manual docker-compose logs commands required for troubleshooting

**Gaps:**
- No long-term log retention
- No full-text search capabilities
- No log-based alerting
- Cannot correlate logs with metrics/traces in unified dashboard
- Difficult to debug distributed transactions across services

## Solution Design

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      DictaMesh Services                          │
│  (Go services with Zap logger outputting JSON to stdout)        │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ JSON logs to stdout
                         ↓
┌─────────────────────────────────────────────────────────────────┐
│                       Promtail Agent                             │
│  - Scrapes container logs via Docker socket                     │
│  - Adds labels (container_name, service, environment)           │
│  - Parses JSON logs and extracts fields                         │
│  - Enriches with metadata                                       │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ HTTP Push
                         ↓
┌─────────────────────────────────────────────────────────────────┐
│                         Loki                                     │
│  - Receives logs from Promtail                                  │
│  - Indexes by labels (not full-text)                            │
│  - Stores log chunks in filesystem                              │
│  - Provides LogQL query API                                     │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ LogQL queries
                         ↓
┌─────────────────────────────────────────────────────────────────┐
│                        Grafana                                   │
│  - Visualizes logs alongside metrics (Prometheus)               │
│  - Correlates logs with traces (Jaeger)                         │
│  - Provides unified observability dashboard                     │
│  - Enables log-based alerts                                     │
└─────────────────────────────────────────────────────────────────┘
```

### Components

#### 1. Loki
- **Version:** 2.9.3
- **Purpose:** Time-series log aggregation and storage
- **Storage:** Local filesystem (development) - S3/GCS for production
- **Retention:** 30 days (configurable)
- **Ports:** 3100 (HTTP API)
- **Resource Limits:** 512MB RAM, 0.5 CPU

**Key Features:**
- Label-based indexing (efficient storage)
- LogQL query language (similar to PromQL)
- Native Grafana integration
- Horizontal scalability
- Multi-tenancy support

#### 2. Promtail
- **Version:** 2.9.3
- **Purpose:** Log collection agent
- **Source:** Docker container logs via /var/lib/docker/containers
- **Transport:** HTTP push to Loki
- **Resource Limits:** 256MB RAM, 0.25 CPU

**Key Features:**
- Automatic service discovery (Docker labels)
- JSON log parsing
- Label extraction from log fields
- Pipeline processing (filtering, parsing, labeling)
- Batching and compression

### Label Strategy

Effective labeling is critical for Loki performance. We use the following label hierarchy:

**Static Labels (set by Promtail):**
```yaml
container_name: dictamesh-metadata-catalog
service: metadata-catalog
environment: development
job: dictamesh-logs
```

**Extracted Labels (from log content):**
```yaml
level: info|warn|error|fatal
logger: service.component
trace_id: <opentelemetry-trace-id>
span_id: <opentelemetry-span-id>
```

**Best Practices:**
- Keep cardinality low (avoid user IDs, request IDs as labels)
- Use labels for filtering, not for data
- Store high-cardinality data in log fields, not labels
- Maximum ~20 unique label combinations per service

### Log Format

**Current Zap JSON Output:**
```json
{
  "level": "info",
  "ts": "2025-11-08T10:30:45.123Z",
  "caller": "service/handler.go:42",
  "msg": "Request processed successfully",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "trace_sampled": true,
  "service": "metadata-catalog",
  "method": "GET",
  "path": "/api/v1/entities",
  "status": 200,
  "duration_ms": 45.3,
  "user_id": "user-123"
}
```

**Loki Storage:**
- Labels: `{container_name="dictamesh-metadata-catalog", level="info", logger="service.handler"}`
- Log line: Full JSON stored as log content
- Indexed fields: Only label values are indexed
- Searchable: Full JSON can be queried with LogQL filters

### LogQL Query Examples

**Filter by service and level:**
```logql
{container_name="dictamesh-metadata-catalog"} |= "error"
```

**Extract and filter by JSON field:**
```logql
{service="metadata-catalog"} | json | duration_ms > 1000
```

**Correlate with trace:**
```logql
{job="dictamesh-logs"} | json | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
```

**Count errors by service (metrics from logs):**
```logql
sum by (service) (count_over_time({job="dictamesh-logs"} | json | level="error" [5m]))
```

**Pattern extraction:**
```logql
{service="metadata-catalog"} | pattern `<_> method=<method> status=<status>` | status >= 500
```

## Implementation Plan

### Phase 1: Infrastructure Setup (Current Phase)

**Tasks:**
1. Add Loki service to docker-compose.dev.yml
2. Add Promtail service to docker-compose.dev.yml
3. Create Loki configuration file
4. Create Promtail configuration file
5. Add persistent volumes for Loki data
6. Configure health checks

**Deliverables:**
- `infrastructure/docker-compose/loki/loki-config.yml`
- `infrastructure/docker-compose/promtail/promtail-config.yml`
- Updated `docker-compose.dev.yml`

### Phase 2: Grafana Integration

**Tasks:**
1. Add Loki datasource to Grafana provisioning
2. Create log exploration dashboard
3. Create service-specific dashboards
4. Add log panels to existing metric dashboards
5. Configure exemplars (link logs ↔ traces)

**Deliverables:**
- Updated `grafana/provisioning/datasources/datasources.yml`
- `grafana/provisioning/dashboards/logs-overview.json`
- `grafana/provisioning/dashboards/service-logs.json`

### Phase 3: Application Configuration

**Tasks:**
1. Verify Go services output JSON logs (already done)
2. Add container labels for service discovery
3. Configure log levels via environment variables
4. Add trace context to all log statements (already done)

**Deliverables:**
- Environment variable configuration
- Updated service configurations

### Phase 4: Documentation & Testing

**Tasks:**
1. Document LogQL query patterns
2. Create runbook for common log queries
3. Add alerting rules for log patterns
4. Test end-to-end log flow
5. Validate log retention and rotation

**Deliverables:**
- Operations documentation
- Log query runbook
- Alert rules

## Configuration Details

### Loki Configuration

**Key Settings:**
```yaml
auth_enabled: false  # Disable multi-tenancy for dev

server:
  http_listen_port: 3100
  grpc_listen_port: 9096

ingester:
  lifecycler:
    ring:
      replication_factor: 1  # Single instance
  chunk_idle_period: 3m
  chunk_retain_period: 1m
  max_chunk_age: 1h

schema_config:
  configs:
    - from: 2024-01-01
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/boltdb-shipper-active
    cache_location: /loki/boltdb-shipper-cache
  filesystem:
    directory: /loki/chunks

limits_config:
  retention_period: 720h  # 30 days
  ingestion_rate_mb: 10
  ingestion_burst_size_mb: 20
```

### Promtail Configuration

**Key Settings:**
```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: dictamesh-containers
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s

    relabel_configs:
      # Only scrape dictamesh containers
      - source_labels: ['__meta_docker_container_name']
        regex: '/dictamesh-.*'
        action: keep

      # Add container name label
      - source_labels: ['__meta_docker_container_name']
        target_label: 'container_name'
        regex: '/(.*)'

      # Extract service name from container
      - source_labels: ['container_name']
        target_label: 'service'
        regex: 'dictamesh-(.*)'

    pipeline_stages:
      # Parse JSON logs
      - json:
          expressions:
            level: level
            logger: logger
            trace_id: trace_id
            span_id: span_id

      # Add labels from JSON
      - labels:
          level:
          logger:
          trace_id:
```

## Integration with Existing Observability Stack

### Prometheus Integration
- Use `{job="dictamesh-logs"} | json | __error__=""` to validate log parsing
- Create log-based metrics for alerting
- Correlate metric spikes with log patterns

### Jaeger Integration
- Link logs to traces via `trace_id` field
- Click trace ID in logs → jump to Jaeger
- View all logs for a specific trace
- Unified view in Grafana Explore

### Sentry Integration
- Sentry continues to handle error tracking and aggregation
- Loki provides raw log context around errors
- Use Sentry for error alerting, Loki for debugging context

### Grafana Dashboards

**Unified Observability Dashboard:**
```
┌─────────────────────────────────────────────────────────────┐
│  Service: metadata-catalog                                  │
├─────────────────────────────────────────────────────────────┤
│  [Prometheus] Request Rate │ [Prometheus] Error Rate       │
│  [Prometheus] Latency P95  │ [Prometheus] CPU Usage        │
├─────────────────────────────────────────────────────────────┤
│  [Loki] Recent Error Logs                                   │
│  [Shows last 50 error logs with trace_id links]            │
├─────────────────────────────────────────────────────────────┤
│  [Jaeger] Active Traces                                     │
│  [Click trace → see all logs for that trace in Loki]       │
└─────────────────────────────────────────────────────────────┘
```

## Resource Requirements

**Development Environment:**
- Loki: 512MB RAM, 0.5 CPU, 10GB disk
- Promtail: 256MB RAM, 0.25 CPU, negligible disk

**Production Environment (estimated):**
- Loki: 2-4GB RAM, 2 CPU, 500GB disk (for 30-day retention)
- Promtail: 512MB RAM, 0.5 CPU per node
- Consider Loki clustering (3+ nodes) for high availability

## Security Considerations

**Development:**
- No authentication (internal network only)
- Docker socket access for Promtail (read-only)

**Production:**
- Enable multi-tenancy with authentication
- Use TLS for Promtail → Loki communication
- Restrict Docker socket access
- Consider log data sensitivity (PII, secrets)
- Implement RBAC in Grafana for log access

## Monitoring Loki Itself

**Key Metrics:**
```promql
# Ingestion rate
sum(rate(loki_distributor_bytes_received_total[1m]))

# Query performance
histogram_quantile(0.99, rate(loki_request_duration_seconds_bucket[5m]))

# Storage usage
loki_ingester_chunk_stored_bytes_total
```

**Alerts:**
- Loki ingestion failures
- Disk space running low
- Query latency degradation

## Migration Path

**Phase 1 (Current):** Docker Compose development environment
**Phase 2:** Kubernetes deployment with Helm charts
**Phase 3:** Production with S3/GCS backend storage
**Phase 4:** Loki clustering for high availability

## Success Criteria

- [ ] Logs from all DictaMesh services visible in Grafana
- [ ] Trace ID correlation working (click log → view trace)
- [ ] Log search response time < 2 seconds for 24h queries
- [ ] 30-day log retention configured and working
- [ ] Zero log data loss during container restarts
- [ ] Developers can debug issues using log queries
- [ ] Log-based alerts configured for critical errors

## References

- [Grafana Loki Documentation](https://grafana.com/docs/loki/latest/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)
- [LogQL Language](https://grafana.com/docs/loki/latest/logql/)
- [Best Practices for Loki](https://grafana.com/docs/loki/latest/best-practices/)
- [OpenTelemetry Context Propagation](https://opentelemetry.io/docs/instrumentation/go/manual/#context-propagation)

## Appendix: Example Queries

### Debugging a Failed Request
```logql
{service="metadata-catalog"}
  | json
  | status >= 500
  | line_format "{{.ts}} [{{.level}}] {{.msg}} (trace: {{.trace_id}})"
```

### Finding Slow Queries
```logql
{service="metadata-catalog"}
  | json
  | duration_ms > 1000
  | line_format "Slow query: {{.path}} took {{.duration_ms}}ms"
```

### Error Rate by Service
```logql
sum by (service) (
  rate({job="dictamesh-logs"} | json | level="error" [5m])
)
```

### All Logs for a Trace
```logql
{job="dictamesh-logs"}
  | json
  | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
```

### Pattern Detection
```logql
{service="metadata-catalog"}
  | pattern `<_> Database connection <status>`
  | status =~ "failed|error"
```
