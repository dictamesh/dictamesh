<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 4
---

# Scaling and Performance

This guide covers horizontal and vertical scaling strategies, performance optimization, and capacity planning for DictaMesh.

## Scaling Strategy Overview

DictaMesh uses a multi-tier scaling approach:

1. **Horizontal Scaling**: Add more replicas (preferred for stateless services)
2. **Vertical Scaling**: Increase resource limits (for stateful services)
3. **Auto-scaling**: Dynamic scaling based on metrics
4. **Data Partitioning**: Shard data across multiple instances

## Horizontal Pod Autoscaling (HPA)

### GraphQL Gateway Autoscaling

The GraphQL Gateway is stateless and scales horizontally easily:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: dictamesh-graphql-gateway
  namespace: dictamesh-system
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: dictamesh-graphql-gateway
  minReplicas: 3
  maxReplicas: 20
  metrics:
    # CPU-based scaling
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70

    # Memory-based scaling
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80

    # Custom metric: requests per second
    - type: Pods
      pods:
        metric:
          name: http_requests_per_second
        target:
          type: AverageValue
          averageValue: "1000"

  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300  # Wait 5 min before scaling down
      policies:
        - type: Percent
          value: 50  # Scale down max 50% of current pods
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0  # Scale up immediately
      policies:
        - type: Percent
          value: 100  # Double pods if needed
          periodSeconds: 15
        - type: Pods
          value: 4  # Or add 4 pods
          periodSeconds: 15
      selectPolicy: Max  # Use the policy that scales fastest
```

Apply the HPA:

```bash
kubectl apply -f graphql-gateway-hpa.yaml

# Verify HPA status
kubectl get hpa -n dictamesh-system

# Watch autoscaling in action
kubectl get hpa dictamesh-graphql-gateway -n dictamesh-system -w
```

### Metadata Catalog Autoscaling

The Metadata Catalog can also scale horizontally:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: dictamesh-metadata-catalog
  namespace: dictamesh-system
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: StatefulSet
    name: dictamesh-metadata-catalog
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80

    # Custom metric: database connection pool utilization
    - type: Pods
      pods:
        metric:
          name: db_connection_pool_utilization
        target:
          type: AverageValue
          averageValue: "0.7"  # 70% pool usage
```

### Custom Metrics with Prometheus Adapter

Install Prometheus Adapter for custom metrics:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus-adapter prometheus-community/prometheus-adapter \
  --namespace dictamesh-system \
  --set prometheus.url=http://dictamesh-prometheus-server \
  --set prometheus.port=80 \
  --values - <<EOF
rules:
  default: false
  custom:
    - seriesQuery: 'http_requests_total{namespace="dictamesh-system"}'
      resources:
        overrides:
          namespace: {resource: "namespace"}
          pod: {resource: "pod"}
      name:
        matches: "^(.*)_total$"
        as: "http_requests_per_second"
      metricsQuery: 'sum(rate(<<.Series>>{<<.LabelMatchers>>}[2m])) by (<<.GroupBy>>)'

    - seriesQuery: 'pg_stat_activity_count{datname="dictamesh_catalog"}'
      resources:
        overrides:
          namespace: {resource: "namespace"}
          pod: {resource: "pod"}
      name:
        matches: "^(.*)$"
        as: "db_connection_pool_utilization"
      metricsQuery: '<<.Series>>{<<.LabelMatchers>>} / on() pg_settings_max_connections'
EOF

# Verify custom metrics are available
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1" | jq .
```

## Vertical Scaling

### When to Use Vertical Scaling

Use vertical scaling for:
- **Stateful services** (PostgreSQL, Kafka, Redis)
- **Memory-intensive workloads** (caching, aggregations)
- **CPU-intensive operations** (schema validation, transformations)

### Vertical Pod Autoscaler (VPA)

Install VPA:

```bash
git clone https://github.com/kubernetes/autoscaler.git
cd autoscaler/vertical-pod-autoscaler
./hack/vpa-up.sh
```

Create VPA for Metadata Catalog:

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: dictamesh-metadata-catalog-vpa
  namespace: dictamesh-system
spec:
  targetRef:
    apiVersion: apps/v1
    kind: StatefulSet
    name: dictamesh-metadata-catalog
  updatePolicy:
    updateMode: "Auto"  # Auto, Initial, Recreate, Off
  resourcePolicy:
    containerPolicies:
      - containerName: metadata-catalog
        minAllowed:
          cpu: 500m
          memory: 512Mi
        maxAllowed:
          cpu: 4000m
          memory: 8Gi
        controlledResources: ["cpu", "memory"]
        mode: Auto
```

## Database Scaling

### PostgreSQL Read Replicas

Scale PostgreSQL reads by adding replicas:

```yaml
# In dictamesh-values.yaml
postgresql:
  readReplicas:
    replicaCount: 3  # Increase from 2 to 3

    resources:
      requests:
        memory: 4Gi
        cpu: 2000m
      limits:
        memory: 8Gi
        cpu: 4000m

    # Pod anti-affinity (spread across nodes)
    podAntiAffinityPreset: hard
```

Apply changes:

```bash
helm upgrade dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --values dictamesh-values.yaml

# Verify replicas
kubectl get pods -n dictamesh-system -l app.kubernetes.io/component=read
```

### Connection Pooling

Use PgBouncer for connection pooling:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pgbouncer
  namespace: dictamesh-system
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pgbouncer
  template:
    metadata:
      labels:
        app: pgbouncer
    spec:
      containers:
        - name: pgbouncer
          image: edoburu/pgbouncer:1.21.0
          ports:
            - containerPort: 5432
          env:
            - name: DATABASE_URL
              value: "postgres://dictamesh:password@dictamesh-postgresql:5432/dictamesh_catalog"
            - name: POOL_MODE
              value: "transaction"
            - name: MAX_CLIENT_CONN
              value: "1000"
            - name: DEFAULT_POOL_SIZE
              value: "25"
            - name: RESERVE_POOL_SIZE
              value: "5"
          resources:
            requests:
              memory: 256Mi
              cpu: 250m
            limits:
              memory: 512Mi
              cpu: 500m
---
apiVersion: v1
kind: Service
metadata:
  name: pgbouncer
  namespace: dictamesh-system
spec:
  selector:
    app: pgbouncer
  ports:
    - port: 5432
      targetPort: 5432
```

Update application to use PgBouncer:

```bash
# Update DATABASE_URL to point to PgBouncer
kubectl set env statefulset/dictamesh-metadata-catalog \
  -n dictamesh-system \
  DATABASE_URL="postgresql://dictamesh:password@pgbouncer:5432/dictamesh_catalog"
```

### Database Partitioning

For large tables, implement partitioning:

```sql
-- Partition entity catalog by creation date
CREATE TABLE dictamesh_entity_catalog (
    id UUID PRIMARY KEY,
    entity_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    -- other columns
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE dictamesh_entity_catalog_2025_11
    PARTITION OF dictamesh_entity_catalog
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');

CREATE TABLE dictamesh_entity_catalog_2025_12
    PARTITION OF dictamesh_entity_catalog
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

-- Create index on each partition
CREATE INDEX idx_entity_catalog_2025_11_type
    ON dictamesh_entity_catalog_2025_11(entity_type);
```

## Kafka Scaling

### Add Kafka Brokers

Scale Kafka horizontally by adding brokers:

```yaml
# In dictamesh-values.yaml
kafka:
  replicaCount: 5  # Increase from 3 to 5
```

Apply changes:

```bash
helm upgrade dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --values dictamesh-values.yaml

# Verify brokers
kubectl get pods -n dictamesh-system -l app.kubernetes.io/component=kafka
```

### Rebalance Partitions

After adding brokers, rebalance partitions:

```bash
# Generate reassignment plan
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-reassign-partitions.sh \
  --bootstrap-server localhost:9092 \
  --topics-to-move-json-file /tmp/topics.json \
  --broker-list "0,1,2,3,4" \
  --generate

# Execute reassignment
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-reassign-partitions.sh \
  --bootstrap-server localhost:9092 \
  --reassignment-json-file /tmp/reassignment.json \
  --execute

# Verify reassignment
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-reassign-partitions.sh \
  --bootstrap-server localhost:9092 \
  --reassignment-json-file /tmp/reassignment.json \
  --verify
```

### Increase Topic Partitions

Scale topic throughput by adding partitions:

```bash
# Increase partitions for high-throughput topics
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --alter \
  --topic dictamesh.entity.events \
  --partitions 24  # Increase from 12 to 24

# Verify
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --describe \
  --topic dictamesh.entity.events
```

### Consumer Group Scaling

Scale consumers to match partition count:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: entity-event-consumer
  namespace: dictamesh-system
spec:
  replicas: 24  # Match number of partitions
  selector:
    matchLabels:
      app: entity-event-consumer
  template:
    metadata:
      labels:
        app: entity-event-consumer
    spec:
      containers:
        - name: consumer
          image: ghcr.io/click2-run/dictamesh-event-consumer:v0.1.0
          env:
            - name: KAFKA_TOPIC
              value: dictamesh.entity.events
            - name: KAFKA_GROUP_ID
              value: entity-consumers
```

## Redis Scaling

### Redis Cluster Mode

For high-throughput caching, use Redis Cluster:

```yaml
# In dictamesh-values.yaml
redis:
  architecture: cluster
  cluster:
    nodes: 6  # 3 masters + 3 replicas
    replicas: 1

  master:
    resources:
      requests:
        memory: 4Gi
        cpu: 1000m
      limits:
        memory: 8Gi
        cpu: 2000m

  replica:
    resources:
      requests:
        memory: 4Gi
        cpu: 500m
      limits:
        memory: 8Gi
        cpu: 1000m
```

### Redis Sentinel (High Availability)

For HA without clustering:

```yaml
redis:
  architecture: replication
  sentinel:
    enabled: true
    masterSet: mymaster
    quorum: 2

  master:
    count: 1

  replica:
    replicaCount: 2
```

## Performance Optimization

### Database Query Optimization

**Add Indexes for Common Queries**

```sql
-- Analyze slow queries
SELECT
    query,
    calls,
    total_time / 1000 as total_time_sec,
    mean_time / 1000 as mean_time_sec,
    max_time / 1000 as max_time_sec
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY total_time DESC
LIMIT 20;

-- Add indexes based on analysis
CREATE INDEX CONCURRENTLY idx_entity_catalog_type_created
    ON dictamesh_entity_catalog(entity_type, created_at DESC);

CREATE INDEX CONCURRENTLY idx_entity_catalog_search
    ON dictamesh_entity_catalog USING GIN(to_tsvector('english', name || ' ' || description));

-- Update statistics
ANALYZE dictamesh_entity_catalog;
```

**Enable Query Result Caching**

```go
// In application code
type QueryCache struct {
    cache *redis.Client
    ttl   time.Duration
}

func (qc *QueryCache) GetOrFetch(key string, fetcher func() (interface{}, error)) (interface{}, error) {
    // Check cache
    val, err := qc.cache.Get(context.Background(), key).Result()
    if err == nil {
        return val, nil
    }

    // Cache miss - fetch from database
    data, err := fetcher()
    if err != nil {
        return nil, err
    }

    // Store in cache
    qc.cache.Set(context.Background(), key, data, qc.ttl)
    return data, nil
}
```

### GraphQL Performance Optimization

**Enable DataLoader Batching**

```typescript
// Example DataLoader configuration
const entityLoader = new DataLoader(
  async (ids: string[]) => {
    // Batch load entities
    const entities = await entityRepository.findByIds(ids);

    // Return in same order as requested
    return ids.map(id => entities.find(e => e.id === id));
  },
  {
    cache: true,
    maxBatchSize: 100,
    batchScheduleFn: (callback) => setTimeout(callback, 10),
  }
);
```

**Implement Query Complexity Limits**

```yaml
# In graphql-gateway-config.yaml
execution:
  max_complexity: 1000
  max_depth: 10

  # Assign complexity to fields
  field_complexity:
    Query.entities: 10
    Entity.relationships: 5
    Entity.schema: 2
```

### Caching Strategy

**Multi-Level Cache Architecture**

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ CDN Cache   │ ← HTTP Cache-Control headers
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ API Gateway │ ← Query result cache (Redis)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Service    │ ← L1 in-memory cache
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Database   │ ← PostgreSQL query cache
└─────────────┘
```

**Cache TTL Strategy**

```yaml
cache_ttl:
  # Static/rarely changing data
  schemas: 1h
  metadata: 1h

  # Semi-static data
  entities: 5m
  relationships: 5m

  # Dynamic data
  query_results: 2m
  aggregations: 1m

  # Real-time data
  events: 10s
  metrics: 30s
```

## Load Testing

### Prepare Load Tests

Install k6 for load testing:

```bash
# Install k6
brew install k6  # macOS
# or
apt-get install k6  # Ubuntu
```

Create load test script `load-test.js`:

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },   // Ramp up to 100 users
    { duration: '5m', target: 100 },   // Stay at 100 users
    { duration: '2m', target: 200 },   // Ramp up to 200 users
    { duration: '5m', target: 200 },   // Stay at 200 users
    { duration: '2m', target: 0 },     // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],  // 95% < 500ms, 99% < 1s
    http_req_failed: ['rate<0.01'],  // Error rate < 1%
  },
};

export default function() {
  const payload = JSON.stringify({
    query: `
      query GetEntities {
        entities(limit: 20) {
          id
          type
          name
          createdAt
        }
      }
    `,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ${__ENV.API_TOKEN}',
    },
  };

  let res = http.post('https://api.dictamesh.example.com/graphql', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
    'no errors': (r) => !r.json().errors,
  });

  sleep(1);
}
```

Run load test:

```bash
k6 run --vus 100 --duration 10m load-test.js
```

### Analyze Results

Monitor during load test:

```bash
# Watch HPA scaling
kubectl get hpa -n dictamesh-system -w

# Watch pod metrics
kubectl top pods -n dictamesh-system

# Check Prometheus metrics
curl 'http://prometheus:9090/api/v1/query?query=rate(http_requests_total[5m])'
```

## Capacity Planning

### Calculate Resource Needs

**Requests per Second (RPS) Capacity**

```
Single Pod Capacity: 1000 RPS (based on testing)
Peak Load: 10,000 RPS
Required Pods: 10,000 / 1000 = 10 pods
With 30% headroom: 10 * 1.3 = 13 pods
```

**Database Connections**

```
Pods: 10
Connections per pod: 25
Total connections: 10 * 25 = 250
Database max_connections: 250 * 1.5 = 375 (with headroom)
```

**Storage Growth**

```
Current data: 100GB
Monthly growth: 10GB
6-month capacity: 100 + (10 * 6) = 160GB
With headroom: 160 * 1.5 = 240GB
```

### Resource Quotas

Set namespace resource quotas:

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: dictamesh-quota
  namespace: dictamesh-system
spec:
  hard:
    requests.cpu: "50"
    requests.memory: 100Gi
    limits.cpu: "100"
    limits.memory: 200Gi
    persistentvolumeclaims: "50"
    requests.storage: 1Ti
```

## Best Practices

### Scaling

✅ **Do:**
- Scale horizontally for stateless services
- Use HPA with multiple metrics
- Test scaling behavior under load
- Monitor scaling events
- Set appropriate min/max replicas

❌ **Don't:**
- Scale based on CPU alone
- Set aggressive scale-down policies
- Forget PodDisruptionBudgets
- Exceed node capacity
- Scale without testing

### Performance

✅ **Do:**
- Cache aggressively at all layers
- Use connection pooling
- Batch database queries
- Index frequently queried columns
- Monitor query performance

❌ **Don't:**
- Cache everything indefinitely
- Ignore slow query logs
- Create too many indexes
- Fetch data unnecessarily
- Skip load testing

## Next Steps

- **[Monitoring](./monitoring.md)** - Monitor scaling metrics
- **[Backup & Restore](./backup-restore.md)** - Backup before scaling operations
- **[Troubleshooting](./troubleshooting.md)** - Debug scaling issues

---

**Previous**: [← Monitoring](./monitoring.md) | **Next**: [Backup & Restore →](./backup-restore.md)
