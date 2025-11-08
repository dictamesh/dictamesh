<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 2
---

# Configuration Reference

This guide provides comprehensive configuration options for DictaMesh components in production environments.

## Configuration Methods

DictaMesh supports multiple configuration methods, applied in the following order (later sources override earlier ones):

1. **Default values** - Sensible defaults built into the application
2. **Configuration files** - YAML files mounted as ConfigMaps
3. **Environment variables** - Kubernetes environment variables
4. **Command-line flags** - Override any setting at runtime

## Metadata Catalog Configuration

### Basic Configuration

Create `metadata-catalog-config.yaml`:

```yaml
# Server configuration
server:
  port: 8080
  host: 0.0.0.0
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  shutdown_timeout: 30s
  max_header_bytes: 1048576  # 1MB

# Database configuration
database:
  url: ${DATABASE_URL}
  driver: postgres

  # Connection pool settings
  max_open_connections: 25
  max_idle_connections: 10
  connection_max_lifetime: 5m
  connection_max_idle_time: 10m

  # Performance tuning
  prepared_statement_cache_size: 100
  default_query_timeout: 30s
  slow_query_threshold: 1s

  # Migration settings
  auto_migrate: false  # Never use in production
  migration_lock_timeout: 10m

# Kafka configuration
kafka:
  brokers:
    - dictamesh-kafka-0.dictamesh-kafka-headless:9092
    - dictamesh-kafka-1.dictamesh-kafka-headless:9092
    - dictamesh-kafka-2.dictamesh-kafka-headless:9092

  # Client settings
  client_id: metadata-catalog

  # Producer settings
  producer:
    acks: all  # Wait for all in-sync replicas
    compression: snappy
    max_message_bytes: 1048576  # 1MB
    timeout: 30s
    retry_max: 3
    retry_backoff: 100ms
    idempotent: true

  # Consumer settings
  consumer:
    group_id: metadata-catalog-consumers
    session_timeout: 10s
    heartbeat_interval: 3s
    max_poll_interval: 300s
    max_poll_records: 500
    auto_offset_reset: earliest
    enable_auto_commit: false  # Manual commit for exactly-once

  # Topic settings
  topics:
    entity_events:
      name: dictamesh.entity.events
      partitions: 12
      replication_factor: 3
      retention_ms: 604800000  # 7 days
    schema_events:
      name: dictamesh.schema.events
      partitions: 6
      replication_factor: 3
      retention_ms: 2592000000  # 30 days

# Cache configuration
cache:
  # L1 cache (in-memory)
  l1:
    enabled: true
    max_size: 10000
    ttl: 5m
    cleanup_interval: 10m

  # L2 cache (Redis)
  redis:
    url: ${REDIS_URL}
    pool_size: 10
    min_idle_conns: 5
    max_conn_age: 30m
    pool_timeout: 4s
    idle_timeout: 5m

    # Sentinel configuration (if using Redis Sentinel)
    # sentinel:
    #   master_name: mymaster
    #   sentinel_addrs:
    #     - sentinel-0:26379
    #     - sentinel-1:26379
    #     - sentinel-2:26379

    # Cluster configuration (if using Redis Cluster)
    # cluster:
    #   addrs:
    #     - redis-0:6379
    #     - redis-1:6379
    #     - redis-2:6379

    # Default TTL per entity type
    ttl:
      entity: 5m
      schema: 1h
      lineage: 10m
      query_result: 2m

# Observability configuration
observability:
  # Tracing
  tracing:
    enabled: true
    provider: jaeger
    jaeger:
      endpoint: http://dictamesh-jaeger-collector:14268/api/traces
      sampler_type: probabilistic
      sampler_param: 0.1  # Sample 10% in production
      agent_host: localhost
      agent_port: 6831

    # Trace specific operations
    trace_sql: true
    trace_kafka: true
    trace_redis: true
    trace_http: true

  # Metrics
  metrics:
    enabled: true
    port: 9090
    path: /metrics

    # Custom metrics
    histograms:
      enabled: true
      buckets: [0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0, 10.0]

    # Metric cardinality limits
    max_labels: 1000

  # Logging
  logging:
    level: info  # debug, info, warn, error
    format: json  # json or text
    output: stdout

    # Log sampling (reduce volume in production)
    sampling:
      enabled: true
      initial: 100
      thereafter: 100

    # Structured fields
    add_caller: false
    add_stacktrace_on_error: true

    # Performance logging
    log_slow_queries: true
    slow_query_threshold: 1s

# Security configuration
security:
  # Authentication
  auth:
    enabled: true
    jwt:
      secret: ${JWT_SECRET}
      issuer: dictamesh
      audience: dictamesh-api
      expiration: 1h
      refresh_expiration: 24h

    # API key authentication
    api_key:
      enabled: true
      header: X-API-Key

  # Authorization (RBAC)
  rbac:
    enabled: true
    cache_ttl: 5m

  # Rate limiting
  rate_limit:
    enabled: true
    requests_per_second: 100
    burst: 200

    # Per-endpoint limits
    endpoints:
      - path: /api/v1/entities
        method: POST
        limit: 50
      - path: /api/v1/query
        method: POST
        limit: 20

  # TLS configuration
  tls:
    enabled: false  # Handled by Ingress
    cert_file: /etc/tls/tls.crt
    key_file: /etc/tls/tls.key
    min_version: "1.3"

# Feature flags
features:
  entity_versioning: true
  schema_validation: true
  lineage_tracking: true
  audit_logging: true
  advanced_search: true
  graphql_federation: true

# Resource limits
limits:
  max_query_complexity: 1000
  max_query_depth: 10
  max_batch_size: 100
  max_request_size: 10485760  # 10MB
  max_response_size: 52428800  # 50MB
```

### Deploy as ConfigMap

```bash
kubectl create configmap metadata-catalog-config \
  --namespace dictamesh-system \
  --from-file=config.yaml=metadata-catalog-config.yaml

# Mount in deployment
kubectl patch deployment dictamesh-metadata-catalog \
  --namespace dictamesh-system \
  --type strategic \
  --patch '
spec:
  template:
    spec:
      containers:
      - name: metadata-catalog
        volumeMounts:
        - name: config
          mountPath: /etc/dictamesh
          readOnly: true
      volumes:
      - name: config
        configMap:
          name: metadata-catalog-config
'
```

## GraphQL Gateway Configuration

### Gateway Configuration

Create `graphql-gateway-config.yaml`:

```yaml
# Server configuration
server:
  port: 8000
  host: 0.0.0.0
  read_timeout: 60s
  write_timeout: 60s
  playground_enabled: false  # Disable in production
  introspection_enabled: true

# Apollo Federation
federation:
  # Service list (auto-discovered from Kubernetes services)
  services:
    - name: metadata-catalog
      url: http://dictamesh-metadata-catalog:8080/graphql
    - name: products-adapter
      url: http://products-adapter:8080/graphql
    - name: orders-adapter
      url: http://orders-adapter:8080/graphql

  # Polling interval for schema updates
  schema_poll_interval: 30s

  # Service health checks
  health_check:
    enabled: true
    interval: 10s
    timeout: 5s
    retries: 3

# Query execution
execution:
  # Query complexity limits
  max_complexity: 1000
  max_depth: 10

  # Timeout settings
  timeout: 30s

  # Batch settings
  batch_enabled: true
  batch_wait: 10ms
  batch_max_size: 100

# DataLoader configuration
dataloader:
  enabled: true
  batch_capacity: 1000
  wait: 10ms
  max_batch: 100

  # Per-entity loader configuration
  loaders:
    entity:
      cache_enabled: true
      cache_ttl: 5m
    schema:
      cache_enabled: true
      cache_ttl: 1h

# Caching
cache:
  # Query result caching
  query_cache:
    enabled: true
    ttl: 5m
    max_size: 1000

  # Persisted queries
  persisted_queries:
    enabled: true
    cache:
      enabled: true
      ttl: 24h

# CORS configuration
cors:
  enabled: true
  allowed_origins:
    - https://app.example.com
    - https://admin.example.com
  allowed_methods:
    - GET
    - POST
    - OPTIONS
  allowed_headers:
    - Content-Type
    - Authorization
    - X-Request-ID
  expose_headers:
    - X-Request-ID
  allow_credentials: true
  max_age: 3600

# Security
security:
  # Authentication
  auth:
    enabled: true
    required_for_introspection: true
    required_for_playground: true

  # Field-level authorization
  field_auth:
    enabled: true

  # Query depth limiting
  depth_limit:
    enabled: true
    max_depth: 10

  # Disable introspection in production
  disable_introspection: false

# Observability
observability:
  tracing:
    enabled: true
    trace_resolvers: true
    trace_validation: true
    trace_parsing: true

  metrics:
    enabled: true
    port: 9091

    # Track query performance
    track_query_performance: true
    track_field_performance: true

  logging:
    level: info
    format: json
    log_queries: false  # Set to true for debugging
    log_errors: true

# Error handling
errors:
  # Include stack traces in development
  include_stacktrace: false

  # Include extensions in errors
  include_extensions: true

  # Mask internal errors
  mask_internal_errors: true
  error_message: "Internal server error"

# Performance
performance:
  # Connection pooling
  max_idle_conns: 10
  max_open_conns: 100

  # Compression
  compression:
    enabled: true
    level: 6  # 1-9, 6 is default
    min_size: 1024  # Minimum bytes to compress
```

## PostgreSQL Configuration

### PostgreSQL Tuning

Create `postgresql-extended.conf`:

```ini
# Connection settings
max_connections = 200
superuser_reserved_connections = 3

# Memory settings (for 32GB RAM server)
shared_buffers = 8GB
effective_cache_size = 24GB
work_mem = 64MB
maintenance_work_mem = 2GB
huge_pages = try

# Checkpoint settings
checkpoint_timeout = 15min
checkpoint_completion_target = 0.9
max_wal_size = 4GB
min_wal_size = 1GB

# Write-Ahead Log
wal_buffers = 16MB
wal_compression = on
wal_level = replica
max_wal_senders = 5
wal_keep_size = 1GB

# Query planning
random_page_cost = 1.1  # For SSD
effective_io_concurrency = 200  # For SSD
default_statistics_target = 100

# Logging
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB
log_min_duration_statement = 1000  # Log queries > 1s
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_checkpoints = on
log_connections = on
log_disconnections = on
log_lock_waits = on
log_temp_files = 0

# Performance
shared_preload_libraries = 'pg_stat_statements'
pg_stat_statements.track = all
pg_stat_statements.max = 10000

# Autovacuum tuning
autovacuum_max_workers = 4
autovacuum_naptime = 10s
autovacuum_vacuum_scale_factor = 0.05
autovacuum_analyze_scale_factor = 0.02

# Replication
hot_standby = on
hot_standby_feedback = on
max_standby_streaming_delay = 30s
```

Apply via ConfigMap:

```bash
kubectl create configmap postgresql-extended-config \
  --namespace dictamesh-system \
  --from-file=extended.conf=postgresql-extended.conf

# Update Helm values
cat <<EOF >> dictamesh-values.yaml
postgresql:
  primary:
    extendedConfiguration: |-
      $(cat postgresql-extended.conf)
  readReplicas:
    extendedConfiguration: |-
      $(cat postgresql-extended.conf)
EOF

helm upgrade dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --values dictamesh-values.yaml
```

## Kafka Configuration

### Broker Configuration

Key Kafka settings for production:

```yaml
# In dictamesh-values.yaml
kafka:
  config:
    # Replication and durability
    default.replication.factor: 3
    min.insync.replicas: 2
    unclean.leader.election.enable: false

    # Performance
    num.network.threads: 8
    num.io.threads: 16
    socket.send.buffer.bytes: 102400
    socket.receive.buffer.bytes: 102400
    socket.request.max.bytes: 104857600  # 100MB

    # Log retention
    log.retention.hours: 168  # 7 days
    log.retention.bytes: 107374182400  # 100GB per partition
    log.segment.bytes: 1073741824  # 1GB
    log.cleanup.policy: delete

    # Compression
    compression.type: snappy

    # Partition management
    num.partitions: 12
    auto.create.topics.enable: false  # Explicit topic creation only

    # Consumer group settings
    group.initial.rebalance.delay.ms: 3000

    # Transaction settings (for exactly-once)
    transaction.state.log.replication.factor: 3
    transaction.state.log.min.isr: 2

    # Zookeeper settings
    zookeeper.connection.timeout.ms: 18000
    zookeeper.session.timeout.ms: 18000
```

### Topic Configuration

Create topics with specific configurations:

```bash
# Create topic with specific config
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-topics.sh --create \
  --bootstrap-server localhost:9092 \
  --topic dictamesh.entity.events \
  --partitions 12 \
  --replication-factor 3 \
  --config retention.ms=604800000 \
  --config min.insync.replicas=2 \
  --config compression.type=snappy \
  --config segment.ms=86400000 \
  --config cleanup.policy=delete

# Verify topic configuration
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-topics.sh --describe \
  --bootstrap-server localhost:9092 \
  --topic dictamesh.entity.events
```

## Redis Configuration

### Redis Tuning

```yaml
# In dictamesh-values.yaml
redis:
  master:
    extraFlags:
      - "--maxmemory-policy allkeys-lru"
      - "--maxmemory 4gb"
      - "--save ''"  # Disable RDB snapshots
      - "--appendonly yes"  # Enable AOF
      - "--appendfsync everysec"
      - "--auto-aof-rewrite-percentage 100"
      - "--auto-aof-rewrite-min-size 64mb"
      - "--timeout 300"
      - "--tcp-keepalive 60"
      - "--maxclients 10000"

  replica:
    extraFlags:
      - "--maxmemory-policy allkeys-lru"
      - "--maxmemory 4gb"
      - "--save ''"
      - "--appendonly yes"
      - "--appendfsync everysec"
      - "--replica-read-only yes"
```

## Environment Variables Reference

### Metadata Catalog

```bash
# Database
DATABASE_URL=postgresql://user:pass@host:5432/dbname
DATABASE_MAX_CONNECTIONS=25
DATABASE_MAX_IDLE_CONNECTIONS=10
DATABASE_CONNECTION_MAX_LIFETIME=5m

# Kafka
KAFKA_BROKERS=kafka-0:9092,kafka-1:9092,kafka-2:9092
KAFKA_CLIENT_ID=metadata-catalog
KAFKA_CONSUMER_GROUP=metadata-catalog-consumers

# Redis
REDIS_URL=redis://:password@host:6379/0
REDIS_POOL_SIZE=10

# Observability
JAEGER_ENDPOINT=http://jaeger:14268/api/traces
PROMETHEUS_PORT=9090
LOG_LEVEL=info
LOG_FORMAT=json

# Security
JWT_SECRET=your-secret-here
ENABLE_AUTH=true
ENABLE_RBAC=true

# Features
ENABLE_ENTITY_VERSIONING=true
ENABLE_SCHEMA_VALIDATION=true
ENABLE_LINEAGE_TRACKING=true
```

### GraphQL Gateway

```bash
# Server
SERVER_PORT=8000
PLAYGROUND_ENABLED=false
INTROSPECTION_ENABLED=true

# Federation
METADATA_CATALOG_URL=http://metadata-catalog:8080/graphql
SCHEMA_POLL_INTERVAL=30s

# Cache
ENABLE_QUERY_CACHE=true
QUERY_CACHE_TTL=5m
ENABLE_DATALOADER=true
DATALOADER_WAIT=10ms

# Security
ENABLE_AUTH=true
JWT_SECRET=your-secret-here
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=https://app.example.com

# Performance
MAX_QUERY_COMPLEXITY=1000
MAX_QUERY_DEPTH=10
REQUEST_TIMEOUT=30s
```

## Configuration Validation

### Validate Configuration

```bash
# Create validation script
cat > /tmp/validate-config.sh <<'EOF'
#!/bin/bash
set -e

echo "Validating DictaMesh configuration..."

# Check metadata catalog config
kubectl exec -n dictamesh-system dictamesh-metadata-catalog-0 -- \
  /app/metadata-catalog validate-config --config /etc/dictamesh/config.yaml

# Check GraphQL gateway config
kubectl exec -n dictamesh-system \
  $(kubectl get pod -n dictamesh-system -l app=graphql-gateway -o jsonpath='{.items[0].metadata.name}') -- \
  /app/graphql-gateway validate-config --config /etc/dictamesh/config.yaml

echo "Configuration validation successful!"
EOF

chmod +x /tmp/validate-config.sh
/tmp/validate-config.sh
```

## Dynamic Configuration Updates

### Hot Reload Configuration

```bash
# Update ConfigMap
kubectl create configmap metadata-catalog-config \
  --namespace dictamesh-system \
  --from-file=config.yaml=metadata-catalog-config.yaml \
  --dry-run=client -o yaml | kubectl apply -f -

# Trigger rolling restart to pick up new config
kubectl rollout restart statefulset/dictamesh-metadata-catalog -n dictamesh-system
kubectl rollout restart deployment/dictamesh-graphql-gateway -n dictamesh-system

# Watch rollout
kubectl rollout status statefulset/dictamesh-metadata-catalog -n dictamesh-system
```

## Configuration Best Practices

### Security

✅ **Do:**
- Store secrets in Kubernetes Secrets, never in ConfigMaps
- Use environment variable substitution in config files
- Rotate secrets regularly (every 90 days)
- Enable audit logging
- Use TLS for all external connections

❌ **Don't:**
- Commit secrets to version control
- Use default passwords
- Disable authentication in production
- Share secrets across environments

### Performance

✅ **Do:**
- Tune connection pools based on workload
- Enable caching at all layers
- Use prepared statements
- Monitor slow queries
- Set appropriate timeouts

❌ **Don't:**
- Use unlimited connection pools
- Disable caching in production
- Set timeouts too low (causes failures) or too high (wastes resources)
- Ignore performance metrics

### Reliability

✅ **Do:**
- Set resource limits and requests
- Configure health checks
- Enable auto-scaling
- Use circuit breakers
- Implement retry logic with backoff

❌ **Don't:**
- Run without resource limits
- Disable health checks
- Ignore pod disruption budgets
- Use infinite retries

## Next Steps

- **[Monitoring](./monitoring.md)** - Set up monitoring and alerting
- **[Scaling](./scaling.md)** - Scale your deployment
- **[Troubleshooting](./troubleshooting.md)** - Debug configuration issues

---

**Previous**: [← Installation](./installation.md) | **Next**: [Monitoring →](./monitoring.md)
