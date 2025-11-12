<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 6
---

# Troubleshooting

This guide covers common issues, diagnostic procedures, and solutions for DictaMesh production deployments.

## Diagnostic Tools

### Essential Commands

```bash
# Check pod status
kubectl get pods -n dictamesh-system

# Describe pod for events and state
kubectl describe pod <pod-name> -n dictamesh-system

# View pod logs
kubectl logs <pod-name> -n dictamesh-system

# View previous pod logs (if crashed)
kubectl logs <pod-name> -n dictamesh-system --previous

# Follow logs in real-time
kubectl logs -f <pod-name> -n dictamesh-system

# Execute commands in pod
kubectl exec -it <pod-name> -n dictamesh-system -- /bin/bash

# Check resource usage
kubectl top pods -n dictamesh-system
kubectl top nodes

# View events
kubectl get events -n dictamesh-system --sort-by='.lastTimestamp'
```

### Quick Health Check Script

```bash
#!/bin/bash
# dictamesh-health-check.sh

echo "=== DictaMesh Health Check ==="

echo -e "\n1. Pod Status:"
kubectl get pods -n dictamesh-system

echo -e "\n2. Service Status:"
kubectl get svc -n dictamesh-system

echo -e "\n3. PVC Status:"
kubectl get pvc -n dictamesh-system

echo -e "\n4. Recent Events:"
kubectl get events -n dictamesh-system --sort-by='.lastTimestamp' | tail -10

echo -e "\n5. Health Endpoints:"
for service in metadata-catalog graphql-gateway; do
    echo "Testing ${service}..."
    kubectl exec -n dictamesh-system $(kubectl get pod -n dictamesh-system -l app=${service} -o jsonpath='{.items[0].metadata.name}') -- \
        curl -s http://localhost:8080/health || echo "Failed"
done

echo -e "\n6. Database Connectivity:"
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
    psql -U dictamesh -d dictamesh_catalog -c "SELECT 1;" > /dev/null && \
    echo "Database: OK" || echo "Database: FAILED"

echo -e "\n7. Kafka Brokers:"
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
    kafka-broker-api-versions.sh --bootstrap-server localhost:9092 > /dev/null && \
    echo "Kafka: OK" || echo "Kafka: FAILED"

echo -e "\n8. Redis Connectivity:"
kubectl exec -it dictamesh-redis-master-0 -n dictamesh-system -- \
    redis-cli ping > /dev/null && \
    echo "Redis: OK" || echo "Redis: FAILED"
```

## Common Issues

### Pod Issues

#### CrashLoopBackOff

**Symptoms:**
```
NAME                                 READY   STATUS             RESTARTS   AGE
dictamesh-metadata-catalog-0         0/1     CrashLoopBackOff   5          10m
```

**Diagnosis:**
```bash
# Check logs
kubectl logs dictamesh-metadata-catalog-0 -n dictamesh-system

# Check previous logs
kubectl logs dictamesh-metadata-catalog-0 -n dictamesh-system --previous

# Check pod events
kubectl describe pod dictamesh-metadata-catalog-0 -n dictamesh-system
```

**Common Causes and Solutions:**

1. **Database connection failure**
   ```bash
   # Verify database secret
   kubectl get secret dictamesh-postgresql -n dictamesh-system -o yaml

   # Test database connectivity
   kubectl run -it --rm postgres-test \
     --image=postgres:15 \
     --namespace dictamesh-system \
     --env="PGPASSWORD=<password>" \
     -- psql -h dictamesh-postgresql -U dictamesh -d dictamesh_catalog
   ```

2. **Missing environment variables**
   ```bash
   # Check environment variables
   kubectl describe pod dictamesh-metadata-catalog-0 -n dictamesh-system | grep -A 20 "Environment:"

   # Verify secrets exist
   kubectl get secrets -n dictamesh-system
   ```

3. **Insufficient resources**
   ```bash
   # Check resource limits
   kubectl describe pod dictamesh-metadata-catalog-0 -n dictamesh-system | grep -A 5 "Limits:"

   # Check node resources
   kubectl describe node <node-name> | grep -A 10 "Allocated resources:"
   ```

#### ImagePullBackOff

**Symptoms:**
```
NAME                                 READY   STATUS             RESTARTS   AGE
dictamesh-graphql-gateway-abc123     0/1     ImagePullBackOff   0          5m
```

**Diagnosis:**
```bash
# Check image pull errors
kubectl describe pod dictamesh-graphql-gateway-abc123 -n dictamesh-system | grep -A 10 "Events:"
```

**Solutions:**

1. **Image does not exist**
   ```bash
   # Verify image name and tag
   kubectl get deployment dictamesh-graphql-gateway -n dictamesh-system -o jsonpath='{.spec.template.spec.containers[0].image}'

   # Check if image exists
   docker pull <image-name>
   ```

2. **Missing image pull secret**
   ```bash
   # Create image pull secret
   kubectl create secret docker-registry ghcr-secret \
     --namespace dictamesh-system \
     --docker-server=ghcr.io \
     --docker-username=<username> \
     --docker-password=<token>

   # Add to service account
   kubectl patch serviceaccount default \
     --namespace dictamesh-system \
     -p '{"imagePullSecrets": [{"name": "ghcr-secret"}]}'
   ```

#### Pending State

**Symptoms:**
```
NAME                               READY   STATUS    RESTARTS   AGE
dictamesh-kafka-0                  0/1     Pending   0          10m
```

**Diagnosis:**
```bash
# Check why pod is pending
kubectl describe pod dictamesh-kafka-0 -n dictamesh-system | grep -A 20 "Events:"
```

**Common Causes:**

1. **Insufficient node resources**
   ```bash
   # Check node capacity
   kubectl describe nodes | grep -A 5 "Allocatable:"

   # Solution: Add nodes or reduce resource requests
   kubectl scale deployment <deployment> --replicas=0
   ```

2. **PVC not bound**
   ```bash
   # Check PVC status
   kubectl get pvc -n dictamesh-system

   # Check storage class
   kubectl get storageclass

   # Solution: Fix storage provisioning or manually create PV
   ```

3. **Pod affinity/anti-affinity constraints**
   ```bash
   # Check affinity rules
   kubectl get pod dictamesh-kafka-0 -n dictamesh-system -o yaml | grep -A 10 "affinity:"

   # Solution: Relax constraints or add more nodes
   ```

### Database Issues

#### High Connection Count

**Symptoms:**
- Error: "sorry, too many clients already"
- Slow query performance
- Connection timeouts

**Diagnosis:**
```sql
-- Check current connections
SELECT count(*), state
FROM pg_stat_activity
WHERE datname = 'dictamesh_catalog'
GROUP BY state;

-- Check connections by application
SELECT application_name, count(*)
FROM pg_stat_activity
WHERE datname = 'dictamesh_catalog'
GROUP BY application_name;

-- Check max connections
SHOW max_connections;
```

**Solutions:**

1. **Increase max_connections**
   ```yaml
   # In dictamesh-values.yaml
   postgresql:
     primary:
       extendedConfiguration: |-
         max_connections = 300  # Increase from 200
   ```

2. **Use connection pooling (PgBouncer)**
   ```bash
   # Deploy PgBouncer (see Scaling guide)
   kubectl apply -f pgbouncer-deployment.yaml
   ```

3. **Kill idle connections**
   ```sql
   -- Kill idle connections older than 5 minutes
   SELECT pg_terminate_backend(pid)
   FROM pg_stat_activity
   WHERE datname = 'dictamesh_catalog'
     AND state = 'idle'
     AND state_change < NOW() - INTERVAL '5 minutes';
   ```

#### Slow Queries

**Symptoms:**
- High p95/p99 latency
- Request timeouts
- Database CPU at 100%

**Diagnosis:**
```sql
-- Find slow queries
SELECT
    query,
    calls,
    total_exec_time / 1000 as total_time_sec,
    mean_exec_time / 1000 as mean_time_sec,
    max_exec_time / 1000 as max_time_sec
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY mean_exec_time DESC
LIMIT 20;

-- Check for missing indexes
SELECT
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    seq_tup_read / seq_scan as avg_seq_tup_read
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC
LIMIT 20;

-- Check table bloat
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
    n_dead_tup,
    n_live_tup,
    round(n_dead_tup * 100.0 / NULLIF(n_live_tup + n_dead_tup, 0), 2) as dead_pct
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC
LIMIT 20;
```

**Solutions:**

1. **Add missing indexes**
   ```sql
   -- Example: Add index for common query
   CREATE INDEX CONCURRENTLY idx_entity_catalog_type_created
       ON dictamesh_entity_catalog(entity_type, created_at DESC);
   ```

2. **Run VACUUM ANALYZE**
   ```sql
   -- Vacuum specific table
   VACUUM ANALYZE dictamesh_entity_catalog;

   -- Or schedule autovacuum more aggressively
   ALTER TABLE dictamesh_entity_catalog SET (
       autovacuum_vacuum_scale_factor = 0.05,
       autovacuum_analyze_scale_factor = 0.02
   );
   ```

3. **Optimize query**
   ```sql
   -- Use EXPLAIN to analyze query plan
   EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM dictamesh_entity_catalog WHERE entity_type = 'Product';
   ```

#### Replication Lag

**Symptoms:**
- Read replicas serving stale data
- Replication lag alerts firing

**Diagnosis:**
```sql
-- Check replication lag (on primary)
SELECT
    client_addr,
    state,
    pg_wal_lsn_diff(pg_current_wal_lsn(), sent_lsn) as send_lag,
    pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) as replay_lag
FROM pg_stat_replication;
```

**Solutions:**

1. **Check network connectivity**
   ```bash
   # Test connectivity between primary and replica
   kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
       pg_isready -h dictamesh-postgresql-1
   ```

2. **Increase WAL sender resources**
   ```yaml
   postgresql:
     primary:
       extendedConfiguration: |-
         max_wal_senders = 10  # Increase from 5
         wal_keep_size = 2GB   # Increase from 1GB
   ```

3. **Check replica load**
   ```sql
   -- Check if replica is overloaded
   SELECT * FROM pg_stat_activity WHERE state != 'idle';
   ```

### Kafka Issues

#### High Consumer Lag

**Symptoms:**
- Events delayed
- Consumer lag alerts
- Messages piling up

**Diagnosis:**
```bash
# Check consumer lag
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
    kafka-consumer-groups.sh \
    --bootstrap-server localhost:9092 \
    --describe \
    --group metadata-catalog-consumers

# Output shows LAG column
# TOPIC           PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG
# entity.events   0          1000            5000            4000  <- High lag!
```

**Solutions:**

1. **Scale consumers**
   ```bash
   # Increase consumer replicas to match partitions
   kubectl scale deployment entity-event-consumer \
       --namespace dictamesh-system \
       --replicas=12  # Match partition count
   ```

2. **Increase partition count**
   ```bash
   # Add more partitions (cannot decrease!)
   kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
       kafka-topics.sh \
       --bootstrap-server localhost:9092 \
       --alter \
       --topic dictamesh.entity.events \
       --partitions 24
   ```

3. **Optimize consumer performance**
   ```go
   // Increase fetch size
   config.Consumer.Fetch.Min = 1024 * 1024  // 1MB
   config.Consumer.Fetch.Max = 10 * 1024 * 1024  // 10MB

   // Increase max poll records
   config.Consumer.MaxProcessingTime = 10 * time.Second
   ```

#### Under-Replicated Partitions

**Symptoms:**
- Data loss risk
- Under-replicated partition alerts
- Some brokers down

**Diagnosis:**
```bash
# Check under-replicated partitions
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
    kafka-topics.sh \
    --bootstrap-server localhost:9092 \
    --describe \
    --under-replicated-partitions

# Check broker status
kubectl get pods -n dictamesh-system -l app.kubernetes.io/component=kafka
```

**Solutions:**

1. **Restart failed brokers**
   ```bash
   # Delete failed pod (StatefulSet will recreate)
   kubectl delete pod dictamesh-kafka-2 -n dictamesh-system

   # Wait for pod to be ready
   kubectl wait --for=condition=ready pod/dictamesh-kafka-2 -n dictamesh-system
   ```

2. **Verify ISR (In-Sync Replicas)**
   ```bash
   # Check topic ISR
   kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
       kafka-topics.sh \
       --bootstrap-server localhost:9092 \
       --describe \
       --topic dictamesh.entity.events
   ```

3. **Rebalance partitions**
   ```bash
   # Use kafka-reassign-partitions (see Scaling guide)
   ```

#### Disk Full

**Symptoms:**
- Kafka brokers crashing
- Cannot write new messages
- Disk usage at 100%

**Diagnosis:**
```bash
# Check disk usage
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- df -h

# Check Kafka data size
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
    du -sh /var/lib/kafka/data/*
```

**Solutions:**

1. **Increase retention**
   ```bash
   # Reduce retention time
   kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
       kafka-configs.sh \
       --bootstrap-server localhost:9092 \
       --entity-type topics \
       --entity-name dictamesh.entity.events \
       --alter \
       --add-config retention.ms=86400000  # 1 day instead of 7
   ```

2. **Expand PVC**
   ```bash
   # Increase PVC size (if storage class supports it)
   kubectl patch pvc data-dictamesh-kafka-0 \
       --namespace dictamesh-system \
       --patch '{"spec":{"resources":{"requests":{"storage":"300Gi"}}}}'
   ```

3. **Delete old topics**
   ```bash
   # Delete unused topics
   kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
       kafka-topics.sh \
       --bootstrap-server localhost:9092 \
       --delete \
       --topic old-unused-topic
   ```

### Redis Issues

#### Memory Exhaustion

**Symptoms:**
- OOM (Out of Memory) errors
- Pod killed by OOMKiller
- Cache misses increase

**Diagnosis:**
```bash
# Check Redis memory usage
kubectl exec -it dictamesh-redis-master-0 -n dictamesh-system -- \
    redis-cli INFO memory

# Check eviction stats
kubectl exec -it dictamesh-redis-master-0 -n dictamesh-system -- \
    redis-cli INFO stats | grep evicted
```

**Solutions:**

1. **Increase memory limit**
   ```yaml
   # In dictamesh-values.yaml
   redis:
     master:
       resources:
         limits:
           memory: 8Gi  # Increase from 4Gi
   ```

2. **Configure eviction policy**
   ```bash
   # Set LRU eviction
   kubectl exec -it dictamesh-redis-master-0 -n dictamesh-system -- \
       redis-cli CONFIG SET maxmemory-policy allkeys-lru
   ```

3. **Reduce TTL**
   ```yaml
   # In configuration
   cache:
     redis:
       ttl:
         entity: 2m  # Reduce from 5m
         query_result: 1m  # Reduce from 2m
   ```

#### Connection Timeout

**Symptoms:**
- "i/o timeout" errors
- "connection refused" errors
- Slow cache operations

**Diagnosis:**
```bash
# Check Redis connectivity
kubectl exec -it dictamesh-metadata-catalog-0 -n dictamesh-system -- \
    nc -zv dictamesh-redis-master 6379

# Check Redis logs
kubectl logs dictamesh-redis-master-0 -n dictamesh-system

# Check connection count
kubectl exec -it dictamesh-redis-master-0 -n dictamesh-system -- \
    redis-cli INFO clients
```

**Solutions:**

1. **Increase connection pool**
   ```yaml
   # In configuration
   cache:
     redis:
       pool_size: 20  # Increase from 10
       max_conn_age: 60m
       pool_timeout: 10s
   ```

2. **Check network policies**
   ```bash
   # Verify network policies allow connection
   kubectl get networkpolicies -n dictamesh-system
   ```

3. **Restart Redis**
   ```bash
   kubectl delete pod dictamesh-redis-master-0 -n dictamesh-system
   ```

### GraphQL Gateway Issues

#### High Query Complexity

**Symptoms:**
- Timeout errors
- "Query complexity exceeded" errors
- High CPU usage

**Diagnosis:**
```bash
# Check query complexity logs
kubectl logs -n dictamesh-system \
    $(kubectl get pod -n dictamesh-system -l app=graphql-gateway -o jsonpath='{.items[0].metadata.name}') \
    | grep "complexity"
```

**Solutions:**

1. **Optimize query**
   ```graphql
   # Bad: Fetches too much data
   query {
     entities {
       relationships {
         entity {
           relationships {
             entity {
               name
             }
           }
         }
       }
     }
   }

   # Good: Fetch only what's needed
   query {
     entities {
       name
       relationships(limit: 10) {
         entity {
           name
         }
       }
     }
   }
   ```

2. **Adjust complexity limits**
   ```yaml
   # In graphql-gateway-config.yaml
   execution:
     max_complexity: 2000  # Increase if needed
     max_depth: 15
   ```

3. **Enable query caching**
   ```yaml
   cache:
     query_cache:
       enabled: true
       ttl: 5m
   ```

#### DataLoader N+1 Issues

**Symptoms:**
- Slow field resolution
- High database query count
- Multiple identical queries

**Diagnosis:**
```bash
# Enable query logging
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
    psql -U dictamesh -d dictamesh_catalog -c \
    "ALTER DATABASE dictamesh_catalog SET log_statement = 'all';"

# Check for duplicate queries
kubectl logs dictamesh-metadata-catalog-0 -n dictamesh-system | grep "SELECT" | sort | uniq -c | sort -rn
```

**Solutions:**

1. **Verify DataLoader is enabled**
   ```yaml
   dataloader:
     enabled: true
     batch_wait: 10ms
     max_batch: 100
   ```

2. **Implement proper batching**
   ```typescript
   // Ensure DataLoader is used for relationships
   const relationshipLoader = new DataLoader(async (ids) => {
     return await repository.findByIds(ids);
   });
   ```

### Network Issues

#### DNS Resolution Failures

**Symptoms:**
- "no such host" errors
- Intermittent connection failures
- Service discovery issues

**Diagnosis:**
```bash
# Test DNS resolution
kubectl run -it --rm dns-test \
    --image=busybox \
    --namespace dictamesh-system \
    --restart=Never \
    -- nslookup dictamesh-postgresql

# Check CoreDNS logs
kubectl logs -n kube-system -l k8s-app=kube-dns
```

**Solutions:**

1. **Restart CoreDNS**
   ```bash
   kubectl rollout restart deployment/coredns -n kube-system
   ```

2. **Add DNS config to pod**
   ```yaml
   spec:
     dnsPolicy: ClusterFirst
     dnsConfig:
       options:
         - name: ndots
           value: "2"
         - name: timeout
           value: "2"
         - name: attempts
           value: "2"
   ```

#### Network Policy Blocking Traffic

**Symptoms:**
- Connection refused
- Timeout errors
- Services cannot communicate

**Diagnosis:**
```bash
# List network policies
kubectl get networkpolicies -n dictamesh-system

# Test connectivity
kubectl run -it --rm nettest \
    --image=nicolaka/netshoot \
    --namespace dictamesh-system \
    -- curl http://dictamesh-postgresql:5432
```

**Solutions:**

1. **Check network policy**
   ```bash
   kubectl describe networkpolicy dictamesh-allow-internal -n dictamesh-system
   ```

2. **Temporarily disable for testing**
   ```bash
   kubectl delete networkpolicy dictamesh-allow-internal -n dictamesh-system
   # Test connectivity
   # Re-create policy if needed
   ```

## Performance Debugging

### High CPU Usage

**Diagnosis:**
```bash
# Check CPU usage
kubectl top pods -n dictamesh-system

# Get CPU profile
kubectl exec -it dictamesh-metadata-catalog-0 -n dictamesh-system -- \
    curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof

# Analyze with pprof
go tool pprof cpu.prof
```

**Solutions:**
- Optimize hot code paths
- Add caching
- Scale horizontally

### High Memory Usage

**Diagnosis:**
```bash
# Check memory usage
kubectl top pods -n dictamesh-system

# Get memory profile
kubectl exec -it dictamesh-metadata-catalog-0 -n dictamesh-system -- \
    curl http://localhost:6060/debug/pprof/heap > heap.prof

# Analyze
go tool pprof heap.prof
```

**Solutions:**
- Fix memory leaks
- Reduce cache sizes
- Increase memory limits

### Trace Analysis

**Use Jaeger for distributed tracing:**

```bash
# Find slow traces
# 1. Open Jaeger UI
kubectl port-forward -n dictamesh-system svc/dictamesh-jaeger-query 16686:16686

# 2. Navigate to http://localhost:16686

# 3. Search for:
# - Service: metadata-catalog
# - Min Duration: 1s
# - Limit Results: 20

# 4. Analyze spans to find bottlenecks
```

## Emergency Procedures

### Complete Service Outage

```bash
# 1. Check overall cluster health
kubectl get nodes
kubectl get pods --all-namespaces

# 2. Restart all DictaMesh components
kubectl rollout restart statefulset/dictamesh-metadata-catalog -n dictamesh-system
kubectl rollout restart deployment/dictamesh-graphql-gateway -n dictamesh-system

# 3. Check data layer
kubectl get pods -n dictamesh-system -l app.kubernetes.io/component=postgresql
kubectl get pods -n dictamesh-system -l app.kubernetes.io/component=kafka
kubectl get pods -n dictamesh-system -l app.kubernetes.io/component=redis

# 4. Verify recovery
./dictamesh-health-check.sh
```

### Data Corruption

```bash
# 1. Stop writes
kubectl scale deployment dictamesh-graphql-gateway -n dictamesh-system --replicas=0

# 2. Assess damage
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
    psql -U dictamesh -d dictamesh_catalog -c \
    "SELECT COUNT(*) FROM dictamesh_entity_catalog;"

# 3. Restore from backup (see Backup & Restore guide)

# 4. Verify data integrity

# 5. Resume operations
kubectl scale deployment dictamesh-graphql-gateway -n dictamesh-system --replicas=3
```

## Best Practices

### Troubleshooting

✅ **Do:**
- Start with the basics (pod status, logs, events)
- Use structured logging with correlation IDs
- Check metrics before diving into logs
- Document issues and solutions
- Create runbooks for common issues

❌ **Don't:**
- Make multiple changes at once
- Delete pods without understanding the issue
- Skip checking recent changes
- Ignore warning signs (high CPU, memory, etc.)

### Prevention

✅ **Do:**
- Monitor proactively
- Test in staging first
- Have rollback plans
- Use canary deployments
- Maintain up-to-date documentation

❌ **Don't:**
- Deploy without testing
- Ignore alerts
- Skip backups
- Disable health checks

## Getting Help

### Escalation Path

1. **Check Documentation**: Review relevant operations guides
2. **Search Issues**: Check GitHub issues for similar problems
3. **Community Support**: Ask in DictaMesh Slack/Discord
4. **Commercial Support**: Contact support team (if applicable)

### Providing Information

When reporting issues, include:

```bash
# 1. DictaMesh version
kubectl get deployment dictamesh-metadata-catalog -n dictamesh-system -o jsonpath='{.spec.template.spec.containers[0].image}'

# 2. Kubernetes version
kubectl version --short

# 3. Pod status
kubectl get pods -n dictamesh-system -o wide

# 4. Recent events
kubectl get events -n dictamesh-system --sort-by='.lastTimestamp' | tail -20

# 5. Relevant logs
kubectl logs dictamesh-metadata-catalog-0 -n dictamesh-system --tail=100

# 6. Resource usage
kubectl top pods -n dictamesh-system

# 7. Configuration (sanitize secrets!)
kubectl get configmap dictamesh-config -n dictamesh-system -o yaml
```

## Next Steps

- **[Installation](./installation.md)** - Reinstall if needed
- **[Configuration](./configuration.md)** - Adjust configuration
- **[Monitoring](./monitoring.md)** - Set up better monitoring
- **[Backup & Restore](./backup-restore.md)** - Restore from backup

---

**Previous**: [← Backup & Restore](./backup-restore.md) | **Up**: [Operations Overview](./installation.md)
