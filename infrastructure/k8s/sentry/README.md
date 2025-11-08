# Sentry Kubernetes Deployment

This directory contains Kubernetes manifests for deploying Sentry to Kubernetes clusters.

## Directory Structure

```
sentry/
├── base/               # Base Kubernetes manifests
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── postgres.yaml
│   ├── redis.yaml
│   ├── clickhouse.yaml
│   ├── sentry-web.yaml
│   ├── sentry-worker.yaml
│   ├── ingress.yaml
│   └── kustomization.yaml
├── dev/                # Development overlay
│   └── kustomization.yaml
├── prod/               # Production overlay
│   └── kustomization.yaml
└── README.md          # This file
```

## Prerequisites

- Kubernetes cluster (v1.25+)
- kubectl configured
- kustomize (v4.0+) or kubectl with built-in kustomize
- Ingress controller (e.g., nginx-ingress)
- Optional: cert-manager for TLS certificates

## Quick Start

### Deploy to Development Environment

```bash
# From the infrastructure directory
kubectl apply -k k8s/sentry/dev/
```

### Deploy to Production Environment

```bash
# From the infrastructure directory
kubectl apply -k k8s/sentry/prod/
```

## Components

The deployment includes the following components:

### Core Services

1. **sentry-web**: Main Sentry web application and API
2. **sentry-worker**: Background task processors (2-5 replicas)
3. **sentry-cron**: Scheduled task executor
4. **sentry-post-process-forwarder**: Event processing pipeline

### Dependencies

1. **sentry-postgres**: PostgreSQL database (16-alpine)
2. **sentry-redis**: Redis cache and message broker
3. **clickhouse**: Event storage database

## Configuration

### Secrets

**Important**: Update secrets before deploying to production!

```bash
# Generate a new Sentry secret key
python3 -c "import secrets; print(secrets.token_urlsafe(50))"

# Create secrets manually
kubectl create secret generic sentry-secrets \
  --from-literal=SENTRY_SECRET_KEY="your-generated-secret-key" \
  --from-literal=POSTGRES_PASSWORD="your-postgres-password" \
  --from-literal=SENTRY_DB_PASSWORD="your-sentry-db-password" \
  --from-literal=CLICKHOUSE_PASSWORD="your-clickhouse-password" \
  --from-literal=SENTRY_ADMIN_EMAIL="admin@yourdomain.com" \
  --from-literal=SENTRY_ADMIN_PASSWORD="your-admin-password" \
  --namespace dictamesh-sentry \
  --dry-run=client -o yaml | kubectl apply -f -
```

### ConfigMap

Edit `base/configmap.yaml` to customize:
- Email settings
- System configuration
- Authentication settings
- Integration configurations

### Environment-Specific Configuration

**Development** (`dev/kustomization.yaml`):
- Reduced resource requirements
- Single worker replica
- Local ingress hostname

**Production** (`prod/kustomization.yaml`):
- High availability (multiple replicas)
- Increased resource limits
- Larger persistent volumes
- TLS-enabled ingress
- Production hostname

## First-Time Setup

1. **Create the namespace**:
   ```bash
   kubectl apply -k k8s/sentry/dev/
   ```

2. **Wait for all pods to be ready**:
   ```bash
   kubectl get pods -n dictamesh-sentry -w
   ```

3. **Initialize the database**:
   ```bash
   # Get the sentry-web pod name
   SENTRY_POD=$(kubectl get pods -n dictamesh-sentry -l app.kubernetes.io/name=sentry-web -o jsonpath='{.items[0].metadata.name}')

   # Run database migrations
   kubectl exec -n dictamesh-sentry $SENTRY_POD -- sentry upgrade --noinput

   # Create superuser
   kubectl exec -it -n dictamesh-sentry $SENTRY_POD -- sentry createuser \
     --email admin@dictamesh.local \
     --password admin \
     --superuser \
     --no-input
   ```

4. **Access Sentry**:
   - Development: http://sentry-dev.dictamesh.local
   - Production: https://sentry.dictamesh.io

   Add to `/etc/hosts` for local testing:
   ```
   127.0.0.1 sentry-dev.dictamesh.local
   ```

## Scaling

### Horizontal Scaling

Scale Sentry web or worker deployments:

```bash
# Scale web instances
kubectl scale deployment -n dictamesh-sentry dev-sentry-web --replicas=3

# Scale worker instances
kubectl scale deployment -n dictamesh-sentry dev-sentry-worker --replicas=5
```

### Vertical Scaling

Update resource requests/limits in the kustomization overlays:

```yaml
patches:
  - target:
      kind: Deployment
      name: sentry-web
    patch: |-
      - op: replace
        path: /spec/template/spec/containers/0/resources/requests/memory
        value: "2Gi"
```

## Persistence

All stateful components use PersistentVolumeClaims:

- **sentry-postgres-pvc**: PostgreSQL data (10Gi dev, 50Gi prod)
- **sentry-redis-pvc**: Redis data (5Gi)
- **clickhouse-pvc**: ClickHouse data (20Gi dev, 100Gi prod)
- **sentry-data-pvc**: Sentry file storage (10Gi dev, 50Gi prod)

### Backup Recommendations

1. **Database Backups**:
   ```bash
   # PostgreSQL backup
   kubectl exec -n dictamesh-sentry deployment/dev-sentry-postgres -- \
     pg_dump -U sentry sentry > sentry-backup.sql
   ```

2. **Volume Snapshots**:
   Use your cloud provider's volume snapshot feature or Velero for full backups.

## Monitoring

### Health Checks

Check the health of all services:

```bash
# Check all pods
kubectl get pods -n dictamesh-sentry

# Check Sentry web health endpoint
kubectl exec -n dictamesh-sentry deployment/dev-sentry-web -- \
  curl -f http://localhost:9000/_health/
```

### Logs

View logs for troubleshooting:

```bash
# Sentry web logs
kubectl logs -n dictamesh-sentry deployment/dev-sentry-web -f

# Sentry worker logs
kubectl logs -n dictamesh-sentry deployment/dev-sentry-worker -f

# PostgreSQL logs
kubectl logs -n dictamesh-sentry deployment/dev-sentry-postgres -f
```

### Metrics

Sentry exposes metrics that can be scraped by Prometheus:

```yaml
# Add to prometheus.yml
- job_name: 'sentry'
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names:
          - dictamesh-sentry
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_app_kubernetes_io_name]
      regex: sentry-web
      action: keep
```

## Troubleshooting

### Pods Not Starting

1. Check pod status:
   ```bash
   kubectl describe pod -n dictamesh-sentry <pod-name>
   ```

2. Check events:
   ```bash
   kubectl get events -n dictamesh-sentry --sort-by='.lastTimestamp'
   ```

### Database Connection Issues

1. Verify PostgreSQL is running:
   ```bash
   kubectl exec -n dictamesh-sentry deployment/dev-sentry-postgres -- pg_isready -U sentry
   ```

2. Check PostgreSQL logs:
   ```bash
   kubectl logs -n dictamesh-sentry deployment/dev-sentry-postgres
   ```

### ClickHouse Issues

1. Test ClickHouse connection:
   ```bash
   kubectl exec -n dictamesh-sentry deployment/dev-clickhouse -- \
     clickhouse-client --query "SELECT 1"
   ```

2. Check ClickHouse logs:
   ```bash
   kubectl logs -n dictamesh-sentry deployment/dev-clickhouse
   ```

### Sentry Web Not Accessible

1. Check ingress status:
   ```bash
   kubectl get ingress -n dictamesh-sentry
   ```

2. Verify service endpoints:
   ```bash
   kubectl get endpoints -n dictamesh-sentry dev-sentry-web
   ```

3. Test service directly (port-forward):
   ```bash
   kubectl port-forward -n dictamesh-sentry svc/dev-sentry-web 9000:9000
   # Visit http://localhost:9000
   ```

## Production Considerations

### High Availability

1. **Database HA**: Consider using a managed PostgreSQL service or Patroni
2. **Redis HA**: Use Redis Sentinel or Redis Cluster
3. **Multiple Replicas**: Run multiple web and worker instances
4. **Pod Disruption Budgets**: Ensure minimum availability during updates

### Security

1. **Network Policies**: Restrict network access between pods
2. **RBAC**: Use Kubernetes RBAC for access control
3. **Secrets Management**: Use external secret managers (Vault, AWS Secrets Manager)
4. **TLS/SSL**: Enable HTTPS with cert-manager
5. **Pod Security Policies**: Enforce security standards

### Performance

1. **Resource Limits**: Set appropriate CPU and memory limits
2. **Autoscaling**: Configure HPA for web and worker deployments
3. **Database Tuning**: Optimize PostgreSQL and ClickHouse
4. **Redis Optimization**: Configure Redis for production workload

### Backup and Disaster Recovery

1. **Regular Backups**: Schedule automated database backups
2. **Volume Snapshots**: Use volume snapshots for all PVCs
3. **Disaster Recovery Plan**: Document recovery procedures
4. **Test Restores**: Regularly test backup restoration

## Updating

### Update Sentry Version

1. Edit the image version in base manifests
2. Apply the changes:
   ```bash
   kubectl apply -k k8s/sentry/dev/
   ```

3. Monitor the rollout:
   ```bash
   kubectl rollout status -n dictamesh-sentry deployment/dev-sentry-web
   ```

### Rollback

If issues occur, rollback to the previous version:

```bash
kubectl rollout undo -n dictamesh-sentry deployment/dev-sentry-web
```

## Cleanup

To remove all Sentry resources:

```bash
# Development
kubectl delete -k k8s/sentry/dev/

# Production
kubectl delete -k k8s/sentry/prod/

# Delete the namespace (this will delete all resources)
kubectl delete namespace dictamesh-sentry
```

**Warning**: This will delete all data. Backup before cleanup!

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
