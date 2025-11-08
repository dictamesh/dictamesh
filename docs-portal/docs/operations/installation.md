<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 1
---

# Production Installation

This guide covers production-grade installation of DictaMesh on Kubernetes. For development setup, see [Getting Started - Installation](../getting-started/installation.md).

## Prerequisites

### Infrastructure Requirements

**Kubernetes Cluster**
- Kubernetes 1.26 or later
- Minimum 3 worker nodes
- Each node: 8 CPU cores, 32GB RAM, 100GB SSD
- Support for PersistentVolumes (StorageClass with dynamic provisioning)
- LoadBalancer support (MetalLB, cloud provider, etc.)

**Storage**
- StorageClass with ReadWriteOnce support (block storage)
- StorageClass with ReadWriteMany support for shared volumes (optional)
- Minimum 500GB total storage capacity
- SSD or NVMe for PostgreSQL and Kafka

**Networking**
- CNI plugin installed (Calico, Cilium, Weave, etc.)
- NetworkPolicy support for security
- Ingress Controller (NGINX, Traefik, or cloud provider)
- DNS resolution working (CoreDNS)

### Required Tools

```bash
# Verify Kubernetes access
kubectl version --short

# Verify Helm installation
helm version

# Verify cluster access
kubectl cluster-info

# Verify storage classes
kubectl get storageclass
```

Expected output:
```
NAME                 PROVISIONER             RECLAIMPOLICY   VOLUMEBINDINGMODE
standard (default)   kubernetes.io/gce-pd    Delete          Immediate
fast-ssd             kubernetes.io/gce-pd    Delete          Immediate
```

## Installation Steps

### Step 1: Create Namespace and Labels

```bash
# Create namespace with labels
kubectl create namespace dictamesh-system

# Label namespace for monitoring
kubectl label namespace dictamesh-system \
  monitoring=enabled \
  app.kubernetes.io/name=dictamesh \
  app.kubernetes.io/managed-by=helm

# Label namespace for network policies
kubectl label namespace dictamesh-system \
  network-policy=enabled
```

### Step 2: Create Secrets

**Database Credentials**

```bash
# Generate strong passwords
export POSTGRES_PASSWORD=$(openssl rand -base64 32)
export REDIS_PASSWORD=$(openssl rand -base64 32)
export JWT_SECRET=$(openssl rand -base64 64)

# Create PostgreSQL secret
kubectl create secret generic dictamesh-postgresql \
  --namespace dictamesh-system \
  --from-literal=postgres-password="${POSTGRES_PASSWORD}" \
  --from-literal=password="${POSTGRES_PASSWORD}" \
  --from-literal=replication-password="${POSTGRES_PASSWORD}"

# Create Redis secret
kubectl create secret generic dictamesh-redis \
  --namespace dictamesh-system \
  --from-literal=redis-password="${REDIS_PASSWORD}"

# Create application secrets
kubectl create secret generic dictamesh-secrets \
  --namespace dictamesh-system \
  --from-literal=jwt-secret="${JWT_SECRET}" \
  --from-literal=database-url="postgresql://dictamesh:${POSTGRES_PASSWORD}@dictamesh-postgresql:5432/dictamesh_catalog" \
  --from-literal=redis-url="redis://:${REDIS_PASSWORD}@dictamesh-redis:6379/0"

# Save credentials securely (DO NOT commit to version control)
cat > /tmp/dictamesh-credentials.txt <<EOF
PostgreSQL Password: ${POSTGRES_PASSWORD}
Redis Password: ${REDIS_PASSWORD}
JWT Secret: ${JWT_SECRET}
EOF

echo "Credentials saved to /tmp/dictamesh-credentials.txt"
echo "Store these credentials in your password manager!"
```

**TLS Certificates**

```bash
# Option 1: Use cert-manager (recommended)
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: dictamesh-tls
  namespace: dictamesh-system
spec:
  secretName: dictamesh-tls-secret
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
    - api.dictamesh.example.com
    - grafana.dictamesh.example.com
    - jaeger.dictamesh.example.com
EOF

# Option 2: Use existing certificates
kubectl create secret tls dictamesh-tls-secret \
  --namespace dictamesh-system \
  --cert=/path/to/tls.crt \
  --key=/path/to/tls.key
```

### Step 3: Add Helm Repository

```bash
# Add DictaMesh Helm repository
helm repo add dictamesh https://charts.dictamesh.com

# Add dependency repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts

# Update repositories
helm repo update
```

### Step 4: Configure Values

Create `dictamesh-values.yaml`:

```yaml
# Global configuration
global:
  storageClass: "fast-ssd"
  domain: "dictamesh.example.com"

# PostgreSQL configuration
postgresql:
  enabled: true
  architecture: replication
  auth:
    existingSecret: dictamesh-postgresql
    database: dictamesh_catalog
    username: dictamesh
  primary:
    resources:
      requests:
        memory: 4Gi
        cpu: 2000m
      limits:
        memory: 8Gi
        cpu: 4000m
    persistence:
      size: 100Gi
      storageClass: "fast-ssd"
    podAntiAffinity:
      preset: hard
  readReplicas:
    replicaCount: 2
    resources:
      requests:
        memory: 4Gi
        cpu: 1000m
      limits:
        memory: 8Gi
        cpu: 2000m
    persistence:
      size: 100Gi
      storageClass: "fast-ssd"

# Kafka configuration (using Redpanda)
kafka:
  enabled: true
  replicaCount: 3
  resources:
    requests:
      memory: 8Gi
      cpu: 2000m
    limits:
      memory: 16Gi
      cpu: 4000m
  persistence:
    size: 200Gi
    storageClass: "fast-ssd"
  config:
    log.retention.hours: 168  # 7 days
    num.partitions: 12
    default.replication.factor: 3
    min.insync.replicas: 2
  zookeeper:
    replicaCount: 3
    resources:
      requests:
        memory: 1Gi
        cpu: 500m
      limits:
        memory: 2Gi
        cpu: 1000m
    persistence:
      size: 20Gi

# Redis configuration
redis:
  enabled: true
  architecture: replication
  auth:
    existingSecret: dictamesh-redis
  master:
    resources:
      requests:
        memory: 2Gi
        cpu: 1000m
      limits:
        memory: 4Gi
        cpu: 2000m
    persistence:
      size: 20Gi
      storageClass: "fast-ssd"
  replica:
    replicaCount: 2
    resources:
      requests:
        memory: 2Gi
        cpu: 500m
      limits:
        memory: 4Gi
        cpu: 1000m
    persistence:
      size: 20Gi

# Metadata Catalog service
metadataCatalog:
  replicaCount: 3
  image:
    repository: ghcr.io/dictamesh/dictamesh-metadata-catalog
    tag: "v0.1.0"
    pullPolicy: IfNotPresent
  resources:
    requests:
      memory: 1Gi
      cpu: 1000m
    limits:
      memory: 2Gi
      cpu: 2000m
  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70
    targetMemoryUtilizationPercentage: 80
  env:
    LOG_LEVEL: "info"
    DATABASE_MAX_CONNECTIONS: "25"
    DATABASE_MAX_IDLE_CONNECTIONS: "10"
  podAntiAffinity:
    preset: hard

# GraphQL Gateway service
graphqlGateway:
  replicaCount: 3
  image:
    repository: ghcr.io/dictamesh/dictamesh-graphql-gateway
    tag: "v0.1.0"
    pullPolicy: IfNotPresent
  resources:
    requests:
      memory: 512Mi
      cpu: 500m
    limits:
      memory: 1Gi
      cpu: 1000m
  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 20
    targetCPUUtilizationPercentage: 70
  podAntiAffinity:
    preset: hard

# Ingress configuration
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
  hosts:
    - host: api.dictamesh.example.com
      paths:
        - path: /
          pathType: Prefix
          service: graphql-gateway
  tls:
    - secretName: dictamesh-tls-secret
      hosts:
        - api.dictamesh.example.com

# Monitoring configuration
monitoring:
  prometheus:
    enabled: true
    serviceMonitor:
      enabled: true
  grafana:
    enabled: true
    adminPassword: "change-me-in-production"
    ingress:
      enabled: true
      hosts:
        - grafana.dictamesh.example.com
  jaeger:
    enabled: true
    collector:
      replicaCount: 2
    query:
      ingress:
        enabled: true
        hosts:
          - jaeger.dictamesh.example.com

# Security settings
security:
  podSecurityPolicy:
    enabled: true
  networkPolicy:
    enabled: true
  rbac:
    create: true
```

### Step 5: Install DictaMesh

```bash
# Dry-run first to validate
helm install dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --values dictamesh-values.yaml \
  --dry-run --debug

# Install for real
helm install dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --values dictamesh-values.yaml \
  --timeout 15m \
  --wait

# Watch the installation
kubectl get pods -n dictamesh-system -w
```

### Step 6: Verify Installation

**Check Pod Status**

```bash
# All pods should be Running
kubectl get pods -n dictamesh-system

# Expected output:
# NAME                                        READY   STATUS    RESTARTS   AGE
# dictamesh-metadata-catalog-0                1/1     Running   0          5m
# dictamesh-metadata-catalog-1                1/1     Running   0          4m
# dictamesh-metadata-catalog-2                1/1     Running   0          3m
# dictamesh-graphql-gateway-7d5f8b9c-abc123   1/1     Running   0          5m
# dictamesh-graphql-gateway-7d5f8b9c-def456   1/1     Running   0          5m
# dictamesh-graphql-gateway-7d5f8b9c-ghi789   1/1     Running   0          5m
# dictamesh-postgresql-0                      1/1     Running   0          5m
# dictamesh-postgresql-1                      1/1     Running   0          4m
# dictamesh-kafka-0                           1/1     Running   0          5m
# dictamesh-kafka-1                           1/1     Running   0          4m
# dictamesh-kafka-2                           1/1     Running   0          3m
# dictamesh-redis-master-0                    1/1     Running   0          5m
# dictamesh-redis-replica-0                   1/1     Running   0          5m
# dictamesh-redis-replica-1                   1/1     Running   0          4m
```

**Check Services**

```bash
kubectl get services -n dictamesh-system

# Expected output showing ClusterIP services
```

**Check PersistentVolumeClaims**

```bash
kubectl get pvc -n dictamesh-system

# All PVCs should be Bound
```

**Check Ingress**

```bash
kubectl get ingress -n dictamesh-system

# Verify ADDRESS is assigned
```

**Test Health Endpoints**

```bash
# Get LoadBalancer IP or use port-forward
kubectl port-forward -n dictamesh-system svc/dictamesh-graphql-gateway 8000:80

# In another terminal
curl http://localhost:8000/health

# Expected: {"status":"healthy","version":"v0.1.0","timestamp":"2025-11-08T..."}
```

### Step 7: Run Database Migrations

```bash
# Create migration job
kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: dictamesh-migrations
  namespace: dictamesh-system
spec:
  template:
    spec:
      restartPolicy: OnFailure
      containers:
      - name: migrate
        image: ghcr.io/dictamesh/dictamesh-metadata-catalog:v0.1.0
        command: ["/app/migrate", "up"]
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: dictamesh-secrets
              key: database-url
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
EOF

# Wait for migration to complete
kubectl wait --for=condition=complete --timeout=300s job/dictamesh-migrations -n dictamesh-system

# Check migration logs
kubectl logs job/dictamesh-migrations -n dictamesh-system

# Expected output:
# INFO: Running migrations...
# INFO: Applied migration 001_initial_schema.up.sql
# INFO: Applied migration 002_entity_catalog.up.sql
# INFO: Applied migration 003_event_log.up.sql
# INFO: All migrations completed successfully
```

### Step 8: Configure DNS

```bash
# Get LoadBalancer IP
kubectl get ingress dictamesh-graphql-gateway -n dictamesh-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}'

# Add DNS records (example for Cloudflare/Route53):
# A record: api.dictamesh.example.com -> <LoadBalancer-IP>
# A record: grafana.dictamesh.example.com -> <LoadBalancer-IP>
# A record: jaeger.dictamesh.example.com -> <LoadBalancer-IP>
```

### Step 9: Test External Access

```bash
# Test GraphQL endpoint
curl https://api.dictamesh.example.com/health

# Test GraphQL playground
open https://api.dictamesh.example.com/playground

# Test Grafana
open https://grafana.dictamesh.example.com

# Test Jaeger
open https://jaeger.dictamesh.example.com
```

## Post-Installation Configuration

### Configure Monitoring

```bash
# Import Grafana dashboards
kubectl apply -f https://raw.githubusercontent.com/dictamesh/dictamesh/main/deployments/monitoring/grafana-dashboards.yaml

# Configure Prometheus alerts
kubectl apply -f https://raw.githubusercontent.com/dictamesh/dictamesh/main/deployments/monitoring/prometheus-rules.yaml
```

### Configure Backup

```bash
# Install Velero for cluster backups
helm repo add vmware-tanzu https://vmware-tanzu.github.io/helm-charts
helm install velero vmware-tanzu/velero \
  --namespace velero \
  --create-namespace \
  --set-file credentials.secretContents.cloud=/path/to/cloud-credentials \
  --set configuration.provider=aws \
  --set configuration.backupStorageLocation.bucket=dictamesh-backups \
  --set configuration.backupStorageLocation.config.region=us-east-1 \
  --set snapshotsEnabled=true

# Create backup schedule
kubectl apply -f - <<EOF
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: dictamesh-daily-backup
  namespace: velero
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  template:
    includedNamespaces:
      - dictamesh-system
    ttl: 720h0m0s  # 30 days retention
EOF
```

### Enable Network Policies

```bash
# Apply network policies for security
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: dictamesh-default-deny
  namespace: dictamesh-system
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: dictamesh-allow-internal
  namespace: dictamesh-system
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: dictamesh
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: dictamesh-system
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: dictamesh-system
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: UDP
      port: 53  # DNS
EOF
```

## Upgrading

### Pre-Upgrade Checklist

- [ ] Review release notes for breaking changes
- [ ] Backup current installation (use Velero)
- [ ] Test upgrade in staging environment
- [ ] Notify users of maintenance window
- [ ] Verify resource quotas are sufficient

### Upgrade Process

```bash
# Update Helm repository
helm repo update

# Check what will change
helm diff upgrade dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --values dictamesh-values.yaml

# Perform upgrade
helm upgrade dictamesh dictamesh/dictamesh \
  --namespace dictamesh-system \
  --values dictamesh-values.yaml \
  --timeout 15m \
  --wait

# Verify upgrade
kubectl rollout status statefulset/dictamesh-metadata-catalog -n dictamesh-system
kubectl rollout status deployment/dictamesh-graphql-gateway -n dictamesh-system

# Run any new migrations
kubectl delete job dictamesh-migrations -n dictamesh-system
kubectl apply -f /path/to/migration-job.yaml
```

## Uninstallation

```bash
# Delete Helm release
helm uninstall dictamesh --namespace dictamesh-system

# Delete PVCs (WARNING: This deletes all data!)
kubectl delete pvc -n dictamesh-system --all

# Delete namespace
kubectl delete namespace dictamesh-system

# Delete CRDs (if any)
kubectl delete crd -l app.kubernetes.io/name=dictamesh
```

## Troubleshooting Installation

### Pods Not Starting

```bash
# Check pod status
kubectl describe pod <pod-name> -n dictamesh-system

# Check logs
kubectl logs <pod-name> -n dictamesh-system --previous

# Common issues:
# - ImagePullBackOff: Check image registry credentials
# - CrashLoopBackOff: Check environment variables and secrets
# - Pending: Check PVC binding and resource availability
```

### Database Connection Issues

```bash
# Test database connectivity
kubectl run -it --rm --restart=Never postgres-test \
  --image=postgres:15 \
  --namespace dictamesh-system \
  --env="PGPASSWORD=${POSTGRES_PASSWORD}" \
  -- psql -h dictamesh-postgresql -U dictamesh -d dictamesh_catalog -c "SELECT 1;"
```

### Helm Installation Failures

```bash
# Check Helm release status
helm status dictamesh -n dictamesh-system

# Get detailed release information
helm get all dictamesh -n dictamesh-system

# Rollback if needed
helm rollback dictamesh -n dictamesh-system
```

## Next Steps

- **[Configuration](./configuration.md)** - Detailed configuration options and tuning
- **[Monitoring](./monitoring.md)** - Set up comprehensive monitoring and alerting
- **[Scaling](./scaling.md)** - Scale your deployment for production workloads
- **[Backup & Restore](./backup-restore.md)** - Implement backup and disaster recovery

---

**Previous**: [← Deployment Guide](../guides/deployment.md) | **Next**: [Configuration →](./configuration.md)
