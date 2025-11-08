<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 3
---

# Monitoring and Observability

This guide covers comprehensive monitoring, observability, and alerting for DictaMesh in production.

## Observability Stack

DictaMesh uses a complete observability stack:

- **Metrics**: Prometheus for metrics collection and storage
- **Visualization**: Grafana for dashboards and visualization
- **Tracing**: Jaeger for distributed tracing
- **Logging**: Structured JSON logs with correlation IDs
- **Alerting**: Alertmanager for alert routing and management

## Metrics Collection

### Prometheus Setup

The DictaMesh Helm chart includes Prometheus by default. For custom setup:

```yaml
# In dictamesh-values.yaml
monitoring:
  prometheus:
    enabled: true
    retention: 15d
    scrapeInterval: 30s
    evaluationInterval: 30s

    resources:
      requests:
        memory: 4Gi
        cpu: 1000m
      limits:
        memory: 8Gi
        cpu: 2000m

    persistence:
      enabled: true
      size: 100Gi
      storageClass: fast-ssd

    # Service monitors for auto-discovery
    serviceMonitor:
      enabled: true
      additionalLabels:
        monitoring: enabled

    # Alert manager configuration
    alertmanager:
      enabled: true
      config:
        global:
          resolve_timeout: 5m
          slack_api_url: ${SLACK_WEBHOOK_URL}

        route:
          group_by: ['alertname', 'cluster', 'service']
          group_wait: 10s
          group_interval: 10s
          repeat_interval: 12h
          receiver: 'default'
          routes:
            - match:
                severity: critical
              receiver: 'critical'
              continue: true
            - match:
                severity: warning
              receiver: 'warning'

        receivers:
          - name: 'default'
            slack_configs:
              - channel: '#dictamesh-alerts'
                title: 'DictaMesh Alert'
                text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'

          - name: 'critical'
            slack_configs:
              - channel: '#dictamesh-critical'
                title: 'CRITICAL: DictaMesh Alert'
                text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
            pagerduty_configs:
              - service_key: ${PAGERDUTY_SERVICE_KEY}

          - name: 'warning'
            slack_configs:
              - channel: '#dictamesh-alerts'
                title: 'Warning: DictaMesh Alert'
                text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```

### Key Metrics

#### System Metrics

**Pod Resource Usage**
```promql
# CPU usage by pod
sum(rate(container_cpu_usage_seconds_total{namespace="dictamesh-system"}[5m])) by (pod)

# Memory usage by pod
sum(container_memory_usage_bytes{namespace="dictamesh-system"}) by (pod)

# Network I/O
sum(rate(container_network_receive_bytes_total{namespace="dictamesh-system"}[5m])) by (pod)
sum(rate(container_network_transmit_bytes_total{namespace="dictamesh-system"}[5m])) by (pod)
```

#### Application Metrics

**Request Rate and Latency**
```promql
# Request rate
sum(rate(http_requests_total{service="metadata-catalog"}[5m])) by (method, path, status)

# Request latency (p95, p99)
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service))
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service))

# Error rate
sum(rate(http_requests_total{status=~"5.."}[5m])) by (service) /
sum(rate(http_requests_total[5m])) by (service)
```

**Database Metrics**
```promql
# Active connections
pg_stat_activity_count{datname="dictamesh_catalog"}

# Transaction rate
rate(pg_stat_database_xact_commit{datname="dictamesh_catalog"}[5m]) +
rate(pg_stat_database_xact_rollback{datname="dictamesh_catalog"}[5m])

# Query duration
histogram_quantile(0.95, sum(rate(pg_stat_statements_total_time_bucket[5m])) by (le))

# Cache hit ratio
pg_stat_database_blks_hit{datname="dictamesh_catalog"} /
(pg_stat_database_blks_hit{datname="dictamesh_catalog"} +
 pg_stat_database_blks_read{datname="dictamesh_catalog"})
```

**Kafka Metrics**
```promql
# Messages in per topic
sum(rate(kafka_server_brokertopicmetrics_messagesinpersec[5m])) by (topic)

# Bytes in/out
sum(rate(kafka_server_brokertopicmetrics_bytesinpersec[5m])) by (topic)
sum(rate(kafka_server_brokertopicmetrics_bytesoutpersec[5m])) by (topic)

# Consumer lag
kafka_consumergroup_lag{consumergroup="metadata-catalog-consumers"}

# Under-replicated partitions (should be 0)
kafka_server_replicamanager_underreplicatedpartitions
```

**Redis Metrics**
```promql
# Hit rate
rate(redis_keyspace_hits_total[5m]) /
(rate(redis_keyspace_hits_total[5m]) + rate(redis_keyspace_misses_total[5m]))

# Memory usage
redis_memory_used_bytes / redis_memory_max_bytes

# Connected clients
redis_connected_clients

# Commands per second
rate(redis_commands_processed_total[5m])
```

**GraphQL Metrics**
```promql
# Query execution time
histogram_quantile(0.95, sum(rate(graphql_query_duration_seconds_bucket[5m])) by (le, operation))

# Query complexity
histogram_quantile(0.95, sum(rate(graphql_query_complexity_bucket[5m])) by (le))

# Errors by type
sum(rate(graphql_errors_total[5m])) by (error_type)

# Dataloader efficiency
rate(graphql_dataloader_batch_size_sum[5m]) / rate(graphql_dataloader_batch_size_count[5m])
```

### ServiceMonitor Configuration

Auto-discovery of metrics endpoints:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: dictamesh-metadata-catalog
  namespace: dictamesh-system
  labels:
    app: metadata-catalog
spec:
  selector:
    matchLabels:
      app: metadata-catalog
  endpoints:
    - port: metrics
      interval: 30s
      path: /metrics
      scrapeTimeout: 10s
```

## Distributed Tracing

### Jaeger Configuration

```yaml
# In dictamesh-values.yaml
monitoring:
  jaeger:
    enabled: true

    collector:
      replicaCount: 2
      resources:
        requests:
          memory: 1Gi
          cpu: 500m
        limits:
          memory: 2Gi
          cpu: 1000m

    query:
      replicaCount: 2
      ingress:
        enabled: true
        hosts:
          - jaeger.dictamesh.example.com

    storage:
      type: elasticsearch
      elasticsearch:
        host: elasticsearch-master
        port: 9200
        indexPrefix: jaeger

    # Sampling strategy
    sampling:
      strategies: |
        {
          "service_strategies": [
            {
              "service": "metadata-catalog",
              "type": "probabilistic",
              "param": 0.1
            },
            {
              "service": "graphql-gateway",
              "type": "probabilistic",
              "param": 0.5
            }
          ],
          "default_strategy": {
            "type": "probabilistic",
            "param": 0.01
          }
        }
```

### Trace Analysis

**Common Trace Queries**

```bash
# Find slow queries (> 1s)
service=metadata-catalog AND duration > 1s

# Find errors
service=metadata-catalog AND error=true

# Find database queries
service=metadata-catalog AND operation=sql.query

# Find Kafka operations
service=metadata-catalog AND operation=kafka.publish
```

**Key Spans to Monitor**

1. **HTTP Request Span**: Total request duration
2. **GraphQL Resolver Span**: Field resolution time
3. **Database Query Span**: SQL execution time
4. **Kafka Publish Span**: Event publishing time
5. **Redis Operation Span**: Cache operation time

### Trace Correlation

All logs include `trace_id` and `span_id` for correlation:

```json
{
  "timestamp": "2025-11-08T10:30:45.123Z",
  "level": "info",
  "service": "metadata-catalog",
  "trace_id": "a1b2c3d4e5f6g7h8",
  "span_id": "i9j0k1l2m3n4",
  "message": "Entity created successfully",
  "entity_id": "ent_123",
  "duration_ms": 45
}
```

## Grafana Dashboards

### Install Dashboards

```bash
# Import DictaMesh dashboards
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: dictamesh-grafana-dashboards
  namespace: dictamesh-system
  labels:
    grafana_dashboard: "1"
data:
  dictamesh-overview.json: |
    $(curl -s https://raw.githubusercontent.com/Click2-Run/dictamesh/main/deployments/monitoring/dashboards/overview.json)
  dictamesh-performance.json: |
    $(curl -s https://raw.githubusercontent.com/Click2-Run/dictamesh/main/deployments/monitoring/dashboards/performance.json)
  dictamesh-infrastructure.json: |
    $(curl -s https://raw.githubusercontent.com/Click2-Run/dictamesh/main/deployments/monitoring/dashboards/infrastructure.json)
EOF
```

### Key Dashboards

#### 1. Overview Dashboard

**Panels:**
- System Health: Overall status of all components
- Request Rate: Requests per second across services
- Error Rate: 5xx errors percentage
- Latency: p50, p95, p99 latencies
- Active Users: Concurrent active sessions
- Resource Usage: CPU and memory utilization

#### 2. Performance Dashboard

**Panels:**
- GraphQL Query Performance: Query execution times
- Database Performance: Query duration, connection pool
- Cache Performance: Hit rates, memory usage
- Kafka Performance: Message throughput, consumer lag
- Slow Query Log: Queries exceeding threshold

#### 3. Infrastructure Dashboard

**Panels:**
- Pod Status: Running/Failed pods
- Node Resources: CPU, memory, disk per node
- Network Traffic: Ingress/egress bandwidth
- Storage Usage: PVC utilization
- Pod Restarts: Restart count per pod

#### 4. Business Metrics Dashboard

**Panels:**
- Entity Operations: Creates, updates, deletes per hour
- Schema Registrations: New schemas over time
- Adapter Activity: Active adapters and data products
- Query Volume: GraphQL queries by type
- Data Lineage Events: Lineage tracking activity

### Example Dashboard JSON

```json
{
  "dashboard": {
    "title": "DictaMesh Overview",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{namespace=\"dictamesh-system\"}[5m])) by (service)",
            "legendFormat": "{{service}}"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{namespace=\"dictamesh-system\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{namespace=\"dictamesh-system\"}[5m]))",
            "legendFormat": "Error Rate"
          }
        ],
        "type": "singlestat",
        "format": "percentunit"
      }
    ]
  }
}
```

## Alerting Rules

### PrometheusRule Configuration

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: dictamesh-alerts
  namespace: dictamesh-system
spec:
  groups:
    - name: dictamesh.availability
      interval: 30s
      rules:
        # High error rate
        - alert: HighErrorRate
          expr: |
            sum(rate(http_requests_total{namespace="dictamesh-system",status=~"5.."}[5m])) by (service) /
            sum(rate(http_requests_total{namespace="dictamesh-system"}[5m])) by (service) > 0.05
          for: 5m
          labels:
            severity: critical
          annotations:
            summary: "High error rate on {{ $labels.service }}"
            description: "{{ $labels.service }} has error rate of {{ $value | humanizePercentage }} (> 5%)"

        # Service down
        - alert: ServiceDown
          expr: up{namespace="dictamesh-system"} == 0
          for: 1m
          labels:
            severity: critical
          annotations:
            summary: "Service {{ $labels.job }} is down"
            description: "{{ $labels.job }} has been down for more than 1 minute"

    - name: dictamesh.performance
      interval: 30s
      rules:
        # High latency
        - alert: HighLatency
          expr: |
            histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{namespace="dictamesh-system"}[5m])) by (le, service)) > 1
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "High latency on {{ $labels.service }}"
            description: "{{ $labels.service }} p95 latency is {{ $value }}s (> 1s)"

        # Database slow queries
        - alert: DatabaseSlowQueries
          expr: |
            histogram_quantile(0.95, sum(rate(pg_stat_statements_total_time_bucket[5m])) by (le)) > 1000
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "Database has slow queries"
            description: "Database p95 query time is {{ $value }}ms (> 1000ms)"

    - name: dictamesh.resources
      interval: 30s
      rules:
        # High CPU usage
        - alert: HighCPUUsage
          expr: |
            sum(rate(container_cpu_usage_seconds_total{namespace="dictamesh-system"}[5m])) by (pod) /
            sum(container_spec_cpu_quota{namespace="dictamesh-system"}/container_spec_cpu_period{namespace="dictamesh-system"}) by (pod) > 0.8
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "High CPU usage on {{ $labels.pod }}"
            description: "{{ $labels.pod }} CPU usage is {{ $value | humanizePercentage }} (> 80%)"

        # High memory usage
        - alert: HighMemoryUsage
          expr: |
            sum(container_memory_usage_bytes{namespace="dictamesh-system"}) by (pod) /
            sum(container_spec_memory_limit_bytes{namespace="dictamesh-system"}) by (pod) > 0.9
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "High memory usage on {{ $labels.pod }}"
            description: "{{ $labels.pod }} memory usage is {{ $value | humanizePercentage }} (> 90%)"

        # Disk space low
        - alert: DiskSpaceLow
          expr: |
            (1 - kubelet_volume_stats_available_bytes{namespace="dictamesh-system"} /
            kubelet_volume_stats_capacity_bytes{namespace="dictamesh-system"}) > 0.85
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "Disk space low on {{ $labels.persistentvolumeclaim }}"
            description: "{{ $labels.persistentvolumeclaim }} is {{ $value | humanizePercentage }} full (> 85%)"

    - name: dictamesh.kafka
      interval: 30s
      rules:
        # High consumer lag
        - alert: HighConsumerLag
          expr: |
            kafka_consumergroup_lag{consumergroup="metadata-catalog-consumers"} > 1000
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "High consumer lag on {{ $labels.topic }}"
            description: "Consumer lag is {{ $value }} messages (> 1000)"

        # Under-replicated partitions
        - alert: UnderReplicatedPartitions
          expr: kafka_server_replicamanager_underreplicatedpartitions > 0
          for: 5m
          labels:
            severity: critical
          annotations:
            summary: "Kafka has under-replicated partitions"
            description: "{{ $value }} partitions are under-replicated"

    - name: dictamesh.database
      interval: 30s
      rules:
        # Connection pool exhaustion
        - alert: ConnectionPoolExhaustion
          expr: |
            pg_stat_activity_count{datname="dictamesh_catalog"} /
            pg_settings_max_connections > 0.9
          for: 5m
          labels:
            severity: critical
          annotations:
            summary: "Database connection pool near exhaustion"
            description: "Database using {{ $value | humanizePercentage }} of max connections (> 90%)"

        # Low cache hit ratio
        - alert: LowCacheHitRatio
          expr: |
            pg_stat_database_blks_hit{datname="dictamesh_catalog"} /
            (pg_stat_database_blks_hit{datname="dictamesh_catalog"} +
             pg_stat_database_blks_read{datname="dictamesh_catalog"}) < 0.9
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "Low database cache hit ratio"
            description: "Cache hit ratio is {{ $value | humanizePercentage }} (< 90%)"
```

## Log Aggregation

### Structured Logging

All DictaMesh services emit structured JSON logs:

```json
{
  "timestamp": "2025-11-08T10:30:45.123Z",
  "level": "info",
  "service": "metadata-catalog",
  "version": "v0.1.0",
  "trace_id": "a1b2c3d4e5f6g7h8",
  "span_id": "i9j0k1l2m3n4",
  "request_id": "req_abc123",
  "user_id": "user_xyz789",
  "operation": "CreateEntity",
  "entity_type": "Product",
  "entity_id": "ent_123",
  "duration_ms": 45,
  "status": "success",
  "message": "Entity created successfully"
}
```

### Log Collection with Loki

```yaml
# Install Loki
helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki-stack \
  --namespace dictamesh-system \
  --set loki.persistence.enabled=true \
  --set loki.persistence.size=100Gi \
  --set promtail.enabled=true

# Query logs in Grafana
# LogQL query examples:
{namespace="dictamesh-system"} |= "error"
{namespace="dictamesh-system", service="metadata-catalog"} | json | level="error"
{namespace="dictamesh-system"} | json | duration_ms > 1000
```

## Health Checks

### Liveness and Readiness Probes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dictamesh-metadata-catalog
spec:
  template:
    spec:
      containers:
        - name: metadata-catalog
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3

          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3

          startupProbe:
            httpGet:
              path: /health/startup
              port: 8080
            initialDelaySeconds: 0
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 30
```

### Health Check Endpoints

**Liveness**: `/health/live`
- Returns 200 if application is running
- Does not check dependencies

**Readiness**: `/health/ready`
- Returns 200 if ready to serve traffic
- Checks database, Kafka, Redis connectivity

**Startup**: `/health/startup`
- Returns 200 when fully initialized
- Used during pod startup

## SLIs and SLOs

### Service Level Indicators (SLIs)

1. **Availability**: Percentage of successful requests
2. **Latency**: p95 and p99 response times
3. **Error Rate**: Percentage of 5xx errors
4. **Throughput**: Requests per second

### Service Level Objectives (SLOs)

```yaml
# Example SLO definitions
slos:
  - name: availability
    target: 99.9%  # 99.9% uptime
    indicator: |
      sum(rate(http_requests_total{status!~"5.."}[30d])) /
      sum(rate(http_requests_total[30d]))

  - name: latency_p95
    target: 200ms
    indicator: |
      histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

  - name: error_rate
    target: 0.1%  # < 0.1% errors
    indicator: |
      sum(rate(http_requests_total{status=~"5.."}[5m])) /
      sum(rate(http_requests_total[5m]))
```

### Error Budget

```promql
# Calculate error budget burn rate
# 99.9% SLO = 0.1% error budget
# 30 days = 43.2 minutes allowed downtime

# Current burn rate
(1 - (sum(rate(http_requests_total{status!~"5.."}[1h])) /
      sum(rate(http_requests_total[1h])))) / 0.001

# Budget remaining this month
1 - ((1 - (sum(rate(http_requests_total{status!~"5.."}[30d])) /
           sum(rate(http_requests_total[30d])))) / 0.001)
```

## Best Practices

### Monitoring

✅ **Do:**
- Monitor golden signals: latency, traffic, errors, saturation
- Set up alerting on SLO violations
- Use distributed tracing for debugging
- Retain metrics for at least 15 days
- Use structured logging with correlation IDs

❌ **Don't:**
- Alert on everything (alert fatigue)
- Ignore high-cardinality metrics
- Keep dashboards without owners
- Log sensitive information (PII, secrets)

### Alerting

✅ **Do:**
- Alert on symptoms, not causes
- Include runbooks in alert annotations
- Use severity levels appropriately
- Test alerts regularly
- Route critical alerts to on-call

❌ **Don't:**
- Create alerts without actions
- Set thresholds too sensitive (false positives)
- Forget to update alerts after changes
- Alert without context

## Next Steps

- **[Scaling](./scaling.md)** - Scale based on metrics
- **[Troubleshooting](./troubleshooting.md)** - Use observability for debugging
- **[Backup & Restore](./backup-restore.md)** - Backup monitoring data

---

**Previous**: [← Configuration](./configuration.md) | **Next**: [Scaling →](./scaling.md)
