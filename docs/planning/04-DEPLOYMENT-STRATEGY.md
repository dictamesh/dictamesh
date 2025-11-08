# Deployment Strategy

[← Previous: Infrastructure Planning](03-INFRASTRUCTURE-PLANNING.md) | [Next: CI/CD Pipeline →](05-CICD-PIPELINE.md)

---

## 🎯 Purpose

Define deployment strategies, release patterns, and rollback procedures for DictaMesh on K3S.

**Reading Time:** 15 minutes
**Prerequisites:** [Infrastructure Planning](03-INFRASTRUCTURE-PLANNING.md)
**Outputs:** Deployment procedures, rollback plans, environment promotion strategies

---

## 🌍 Environment Strategy

### Environment Progression

```
Development → Staging → Production
(dictamesh-dev) → (dictamesh-staging) → (dictamesh-prod)
```

### Environment Characteristics

| Aspect | Development | Staging | Production |
|--------|-------------|---------|------------|
| **Purpose** | Feature development | Integration testing | Live traffic |
| **Data** | Synthetic/anonymized | Production snapshot | Real customer data |
| **Replicas** | 1 per service | 2 per service | 3+ per service |
| **Resources** | Minimal | 50% of prod | Full allocation |
| **Deployment** | Automatic (on push to develop) | Manual approval | Manual approval + change window |
| **Monitoring** | Basic | Full | Full + alerting |
| **Backup** | None | Daily | Continuous + PITR |

---

## 🚀 Deployment Patterns

### 1. Rolling Deployment (Default)

```yaml
# Kubernetes rolling update
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0  # Zero-downtime
```

**Use For:** Stateless services (adapters, gateway)
**Benefits:** Zero-downtime, gradual rollout
**Limitations:** Cannot handle breaking schema changes

### 2. Blue-Green Deployment

```yaml
# infrastructure/k8s/overlays/prod/blue-green-deployment.yaml
# Blue (current) environment
apiVersion: v1
kind: Service
metadata:
  name: graphql-gateway
spec:
  selector:
    app: graphql-gateway
    version: blue  # Initially points to blue

---
# Green (new) deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: graphql-gateway-green
spec:
  selector:
    matchLabels:
      version: green
```

**Cutover:**
```bash
# Test green deployment
kubectl port-forward deployment/graphql-gateway-green 8080:8080

# Switch traffic to green
kubectl patch service graphql-gateway -p '{"spec":{"selector":{"version":"green"}}}'

# Rollback if needed
kubectl patch service graphql-gateway -p '{"spec":{"selector":{"version":"blue"}}}'
```

**Use For:** Major releases, database migrations
**Benefits:** Instant rollback, full testing before cutover
**Limitations:** Requires 2x resources during deployment

### 3. Canary Deployment

```yaml
# Using Argo Rollouts
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: customer-adapter-canary
spec:
  replicas: 5
  strategy:
    canary:
      steps:
        - setWeight: 20  # 20% traffic to new version
        - pause: {duration: 5m}
        - setWeight: 50
        - pause: {duration: 5m}
        - setWeight: 80
        - pause: {duration: 5m}
      trafficRouting:
        istio:
          virtualService:
            name: customer-adapter-vsvc
```

**Use For:** High-risk changes, performance testing
**Benefits:** Gradual rollout, early issue detection
**Limitations:** Requires service mesh or ingress controller support

---

## 📋 Pre-Deployment Checklist

### For All Deployments

- [ ] Tests passing in CI
- [ ] Security scan passed (Trivy)
- [ ] Code review approved
- [ ] Documentation updated
- [ ] Database migrations prepared (if applicable)
- [ ] Rollback plan documented
- [ ] Stakeholders notified
- [ ] Monitoring dashboards ready

### For Production Deployments

- [ ] Staging deployment successful
- [ ] Load testing completed
- [ ] Backup verified within last 24h
- [ ] Change window scheduled
- [ ] On-call engineer available
- [ ] Runbook reviewed
- [ ] Feature flags configured (if applicable)
- [ ] Customer communication prepared (if user-facing changes)

---

## 🔄 Database Migration Strategy

### Schema Evolution Pattern

```go
// Use golang-migrate for versioned migrations
// migrations/000001_create_entity_catalog.up.sql
CREATE TABLE entity_catalog (
    id UUID PRIMARY KEY,
    entity_type VARCHAR(100) NOT NULL,
    -- ...
);

// migrations/000002_add_entity_catalog_index.up.sql
CREATE INDEX idx_entity_type ON entity_catalog(entity_type);

// migrations/000002_add_entity_catalog_index.down.sql
DROP INDEX idx_entity_type;
```

### Migration Execution

```bash
# Development: Auto-migrate
export DATABASE_URL="postgresql://user:pass@localhost:5432/metadata_catalog"
migrate -database $DATABASE_URL -path migrations up

# Production: Manual approval required
kubectl exec -it metadata-catalog-db-rw-0 -n dictamesh-infra -- psql
-- Review migration SQL
\i /migrations/000003_new_migration.up.sql
-- Verify
SELECT * FROM schema_migrations;
```

### Zero-Downtime Migration Guidelines

1. **Additive changes only** (add columns, tables)
2. **Deploy code that works with both old and new schema**
3. **Run migration**
4. **Deploy code that uses new schema**
5. **Remove old schema in next release**

**Example: Adding a column**
```sql
-- Step 1: Add column as nullable
ALTER TABLE entity_catalog ADD COLUMN new_field VARCHAR(255);

-- Step 2: Backfill data (batched)
UPDATE entity_catalog SET new_field = 'default' WHERE new_field IS NULL;

-- Step 3: Add NOT NULL constraint (next release)
ALTER TABLE entity_catalog ALTER COLUMN new_field SET NOT NULL;
```

---

## 🎯 Rollback Procedures

### Application Rollback

#### ArgoCD Rollback
```bash
# View history
argocd app history customer-adapter-prod

# Rollback to previous version
argocd app rollback customer-adapter-prod <REVISION_NUMBER>

# Verify
argocd app wait customer-adapter-prod --health
```

#### kubectl Rollback
```bash
# View deployment history
kubectl rollout history deployment/customer-adapter -n dictamesh-prod

# Rollback to previous revision
kubectl rollout undo deployment/customer-adapter -n dictamesh-prod

# Rollback to specific revision
kubectl rollout undo deployment/customer-adapter -n dictamesh-prod --to-revision=3
```

### Database Rollback

```bash
# Rollback one migration
migrate -database $DATABASE_URL -path migrations down 1

# Rollback to specific version
migrate -database $DATABASE_URL -path migrations force <VERSION>
```

### Kafka Topic Rollback

```bash
# Revert topic configuration
kubectl -n dictamesh-infra apply -f kafka-topics-backup.yaml

# Reset consumer group offset (if needed)
kubectl -n dictamesh-infra exec -it dictamesh-kafka-kafka-0 -- \
  kafka-consumer-groups.sh --reset-offsets --group metadata-catalog-consumer \
  --topic customers.entity_changed --to-datetime 2025-11-08T10:00:00.000
```

---

## 🔐 Production Deployment Procedure

### Phase 1: Preparation (T-24h)

```bash
# 1. Verify staging deployment
argocd app get dictamesh-staging --hard-refresh

# 2. Create production backup
kubectl cnpg backup metadata-catalog-db -n dictamesh-infra

# 3. Verify backup completion
kubectl get backups -n dictamesh-infra

# 4. Tag release
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### Phase 2: Deployment (T-0)

```bash
# 1. Enable maintenance mode (if applicable)
kubectl patch configmap feature-flags -n dictamesh-prod \
  -p '{"data":{"maintenance_mode":"true"}}'

# 2. Update image tags in Git
cd infrastructure/k8s/overlays/prod
kustomize edit set image customer-adapter=ghcr.io/controle-digital/customer-adapter:v1.2.0

# 3. Commit and push
git add .
git commit -m "deploy: customer-adapter v1.2.0 to production"
git push

# 4. Sync ArgoCD application
argocd app sync customer-adapter-prod

# 5. Wait for healthy status
argocd app wait customer-adapter-prod --health --timeout 600
```

### Phase 3: Verification (T+5m)

```bash
# 1. Check pod status
kubectl get pods -n dictamesh-prod -l app=customer-adapter

# 2. Check logs for errors
kubectl logs -n dictamesh-prod deployment/customer-adapter --tail=100 | grep -i error

# 3. Run smoke tests
kubectl apply -f tests/smoke/production-smoke-tests.yaml
kubectl wait --for=condition=complete job/smoke-tests -n dictamesh-prod

# 4. Verify metrics
curl https://grafana.dictamesh.com/api/health

# 5. Disable maintenance mode
kubectl patch configmap feature-flags -n dictamesh-prod \
  -p '{"data":{"maintenance_mode":"false"}}'
```

### Phase 4: Monitoring (T+30m)

- [ ] Monitor error rates (should not increase)
- [ ] Monitor latency (P95, P99)
- [ ] Monitor resource usage
- [ ] Check customer reports/support tickets
- [ ] Verify business metrics

---

## 📊 Health Checks

### Application Health Endpoints

```go
// /health/live - Liveness probe (am I running?)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// /health/ready - Readiness probe (can I serve traffic?)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
    // Check dependencies
    if err := h.db.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    if err := h.kafka.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
```

### Kubernetes Probe Configuration

```yaml
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
```

---

[← Previous: Infrastructure Planning](03-INFRASTRUCTURE-PLANNING.md) | [Next: CI/CD Pipeline →](05-CICD-PIPELINE.md)

---

**Document Metadata**
- Version: 1.0.0
- Last Updated: 2025-11-08
- Deployment Tool: ArgoCD (GitOps)
