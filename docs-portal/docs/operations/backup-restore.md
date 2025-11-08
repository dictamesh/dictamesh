<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 5
---

# Backup and Restore

This guide covers backup strategies, disaster recovery procedures, and data protection for DictaMesh production deployments.

## Backup Strategy Overview

DictaMesh requires backing up multiple data stores:

1. **PostgreSQL** - Metadata catalog (critical)
2. **Kafka** - Event stream (important for replay)
3. **Redis** - Cache layer (optional, can be rebuilt)
4. **Kubernetes State** - Deployments, ConfigMaps, Secrets
5. **Persistent Volumes** - StatefulSet data

## Backup Architecture

### Recovery Objectives

Define your recovery requirements:

**Recovery Point Objective (RPO)**: Maximum acceptable data loss
- **Production**: 15 minutes
- **Staging**: 1 hour
- **Development**: 24 hours

**Recovery Time Objective (RTO)**: Maximum acceptable downtime
- **Production**: 1 hour
- **Staging**: 4 hours
- **Development**: 24 hours

### Backup Types

1. **Continuous Backup**: WAL archiving, real-time replication (RPO ~0)
2. **Snapshot Backup**: Volume snapshots every hour (RPO ~1 hour)
3. **Logical Backup**: Database dumps daily (RPO ~24 hours)
4. **Off-site Backup**: Cross-region replication for disaster recovery

## PostgreSQL Backup

### Continuous WAL Archiving

Enable Write-Ahead Log (WAL) archiving for point-in-time recovery:

```yaml
# In dictamesh-values.yaml
postgresql:
  primary:
    extendedConfiguration: |-
      # WAL archiving
      wal_level = replica
      archive_mode = on
      archive_command = 'envdir /etc/wal-g.d/env wal-g wal-push %p'
      archive_timeout = 300  # Archive every 5 minutes

    # Mount S3 credentials
    extraVolumes:
      - name: wal-g-env
        secret:
          secretName: wal-g-credentials

    extraVolumeMounts:
      - name: wal-g-env
        mountPath: /etc/wal-g.d/env
```

Create WAL-G credentials:

```bash
# Create S3 credentials for WAL-G
kubectl create secret generic wal-g-credentials \
  --namespace dictamesh-system \
  --from-literal=AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  --from-literal=AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  --from-literal=AWS_REGION="us-east-1" \
  --from-literal=WALG_S3_PREFIX="s3://dictamesh-backups/wal" \
  --from-literal=PGHOST="localhost" \
  --from-literal=PGPORT="5432" \
  --from-literal=PGUSER="postgres" \
  --from-literal=PGDATABASE="dictamesh_catalog"
```

### Automated Daily Backups

Create a CronJob for daily base backups:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgresql-backup
  namespace: dictamesh-system
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  successfulJobsHistoryLimit: 7
  failedJobsHistoryLimit: 3
  concurrencyPolicy: Forbid
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: backup
              image: wal-g/wal-g:latest
              command:
                - /bin/bash
                - -c
                - |
                  set -e
                  echo "Starting PostgreSQL backup..."

                  # Create base backup
                  envdir /etc/wal-g.d/env wal-g backup-push /var/lib/postgresql/data

                  # Delete old backups (keep 30 days)
                  envdir /etc/wal-g.d/env wal-g delete retain FULL 30

                  echo "Backup completed successfully"
              env:
                - name: PGPASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: dictamesh-postgresql
                      key: postgres-password
              volumeMounts:
                - name: wal-g-env
                  mountPath: /etc/wal-g.d/env
                - name: data
                  mountPath: /var/lib/postgresql/data
          volumes:
            - name: wal-g-env
              secret:
                secretName: wal-g-credentials
            - name: data
              persistentVolumeClaim:
                claimName: data-dictamesh-postgresql-0
```

### Manual Database Backup

```bash
# Logical backup using pg_dump
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  pg_dump -U dictamesh -d dictamesh_catalog \
  --format=custom \
  --compress=9 \
  --file=/tmp/backup-$(date +%Y%m%d-%H%M%S).dump

# Copy backup to local machine
kubectl cp dictamesh-system/dictamesh-postgresql-0:/tmp/backup-*.dump \
  ./backup-$(date +%Y%m%d-%H%M%S).dump

# Upload to S3
aws s3 cp ./backup-*.dump s3://dictamesh-backups/manual/
```

### Verify Backup Integrity

```bash
# List available backups
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  envdir /etc/wal-g.d/env wal-g backup-list

# Verify backup
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  envdir /etc/wal-g.d/env wal-g backup-fetch /tmp/verify LATEST

# Test restore (dry-run)
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  pg_restore --list /tmp/backup.dump
```

## Kafka Backup

### Mirror Maker for Cross-Cluster Replication

Set up Kafka MirrorMaker 2 for real-time replication:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-mirrormaker
  namespace: dictamesh-system
spec:
  replicas: 2
  selector:
    matchLabels:
      app: kafka-mirrormaker
  template:
    metadata:
      labels:
        app: kafka-mirrormaker
    spec:
      containers:
        - name: mirrormaker
          image: confluentinc/cp-kafka:7.5.0
          command:
            - /bin/bash
            - -c
            - |
              /usr/bin/connect-mirror-maker /etc/kafka/mm2.properties
          volumeMounts:
            - name: config
              mountPath: /etc/kafka
      volumes:
        - name: config
          configMap:
            name: mirrormaker-config
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mirrormaker-config
  namespace: dictamesh-system
data:
  mm2.properties: |
    # Source cluster
    clusters = source, backup
    source.bootstrap.servers = dictamesh-kafka-0:9092,dictamesh-kafka-1:9092,dictamesh-kafka-2:9092

    # Backup cluster (different region/cluster)
    backup.bootstrap.servers = backup-kafka-0:9092,backup-kafka-1:9092,backup-kafka-2:9092

    # Replication flows
    source->backup.enabled = true
    source->backup.topics = dictamesh\\..*

    # Replication settings
    replication.factor = 3
    refresh.topics.interval.seconds = 60
    sync.topic.acls.enabled = false
    emit.checkpoints.interval.seconds = 60
```

### Topic Snapshot Backup

For periodic backups of Kafka topics:

```bash
# Create snapshot backup script
cat > /tmp/kafka-backup.sh <<'EOF'
#!/bin/bash
set -e

BACKUP_DATE=$(date +%Y%m%d-%H%M%S)
BACKUP_DIR="/backups/kafka/${BACKUP_DATE}"
TOPICS=$(kafka-topics.sh --bootstrap-server localhost:9092 --list | grep "^dictamesh\.")

mkdir -p "${BACKUP_DIR}"

for topic in ${TOPICS}; do
    echo "Backing up topic: ${topic}"

    # Export topic to files
    kafka-console-consumer.sh \
        --bootstrap-server localhost:9092 \
        --topic "${topic}" \
        --from-beginning \
        --timeout-ms 60000 \
        --max-messages 1000000 > "${BACKUP_DIR}/${topic}.json"
done

# Compress backup
tar -czf "/backups/kafka-${BACKUP_DATE}.tar.gz" -C /backups/kafka "${BACKUP_DATE}"

# Upload to S3
aws s3 cp "/backups/kafka-${BACKUP_DATE}.tar.gz" s3://dictamesh-backups/kafka/

# Clean up local files
rm -rf "${BACKUP_DIR}"
rm -f "/backups/kafka-${BACKUP_DATE}.tar.gz"

echo "Kafka backup completed: kafka-${BACKUP_DATE}.tar.gz"
EOF

# Create CronJob
kubectl create configmap kafka-backup-script \
  --namespace dictamesh-system \
  --from-file=backup.sh=/tmp/kafka-backup.sh

kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: CronJob
metadata:
  name: kafka-backup
  namespace: dictamesh-system
spec:
  schedule: "0 3 * * *"  # Daily at 3 AM
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: backup
              image: confluentinc/cp-kafka:7.5.0
              command: ["/bin/bash", "/scripts/backup.sh"]
              volumeMounts:
                - name: scripts
                  mountPath: /scripts
                - name: backups
                  mountPath: /backups
          volumes:
            - name: scripts
              configMap:
                name: kafka-backup-script
                defaultMode: 0755
            - name: backups
              emptyDir: {}
EOF
```

## Kubernetes State Backup with Velero

### Install Velero

```bash
# Add Velero Helm repository
helm repo add vmware-tanzu https://vmware-tanzu.github.io/helm-charts

# Create cloud credentials
cat > /tmp/credentials-velero <<EOF
[default]
aws_access_key_id=${AWS_ACCESS_KEY_ID}
aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}
EOF

# Install Velero
helm install velero vmware-tanzu/velero \
  --namespace velero \
  --create-namespace \
  --set-file credentials.secretContents.cloud=/tmp/credentials-velero \
  --set configuration.provider=aws \
  --set configuration.backupStorageLocation.name=default \
  --set configuration.backupStorageLocation.bucket=dictamesh-velero-backups \
  --set configuration.backupStorageLocation.config.region=us-east-1 \
  --set configuration.volumeSnapshotLocation.name=default \
  --set configuration.volumeSnapshotLocation.config.region=us-east-1 \
  --set snapshotsEnabled=true \
  --set deployNodeAgent=true

# Clean up credentials file
rm /tmp/credentials-velero
```

### Configure Backup Schedules

```yaml
# Full namespace backup - daily
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: dictamesh-daily-backup
  namespace: velero
spec:
  schedule: "0 1 * * *"  # Daily at 1 AM
  template:
    includedNamespaces:
      - dictamesh-system
    ttl: 720h0m0s  # 30 days retention
    snapshotVolumes: true
    includeClusterResources: true
---
# Configuration backup - hourly
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: dictamesh-hourly-config
  namespace: velero
spec:
  schedule: "0 * * * *"  # Every hour
  template:
    includedNamespaces:
      - dictamesh-system
    ttl: 168h0m0s  # 7 days retention
    snapshotVolumes: false
    includedResources:
      - configmaps
      - secrets
      - services
      - deployments
      - statefulsets
---
# Weekly full backup with off-site replication
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: dictamesh-weekly-offsite
  namespace: velero
spec:
  schedule: "0 0 * * 0"  # Weekly on Sunday at midnight
  template:
    includedNamespaces:
      - dictamesh-system
    ttl: 2160h0m0s  # 90 days retention
    snapshotVolumes: true
    storageLocation: offsite
```

### Manual Backup

```bash
# Create on-demand backup
velero backup create dictamesh-manual-$(date +%Y%m%d-%H%M%S) \
  --include-namespaces dictamesh-system \
  --snapshot-volumes \
  --wait

# Check backup status
velero backup describe dictamesh-manual-20251108-140000

# List all backups
velero backup get
```

## Restore Procedures

### PostgreSQL Point-in-Time Recovery

Restore to a specific point in time:

```bash
# 1. Stop the current PostgreSQL instance
kubectl scale statefulset dictamesh-postgresql \
  --namespace dictamesh-system \
  --replicas=0

# 2. Create restore job
kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: postgresql-restore
  namespace: dictamesh-system
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: restore
          image: wal-g/wal-g:latest
          command:
            - /bin/bash
            - -c
            - |
              set -e

              # Clean data directory
              rm -rf /var/lib/postgresql/data/*

              # Fetch base backup
              envdir /etc/wal-g.d/env wal-g backup-fetch \
                /var/lib/postgresql/data LATEST

              # Create recovery configuration
              cat > /var/lib/postgresql/data/recovery.signal <<RECOVERY
              restore_command = 'envdir /etc/wal-g.d/env wal-g wal-fetch %f %p'
              recovery_target_time = '2025-11-08 14:00:00'
              recovery_target_action = 'promote'
              RECOVERY

              echo "Restore prepared successfully"
          env:
            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: dictamesh-postgresql
                  key: postgres-password
          volumeMounts:
            - name: wal-g-env
              mountPath: /etc/wal-g.d/env
            - name: data
              mountPath: /var/lib/postgresql/data
      volumes:
        - name: wal-g-env
          secret:
            secretName: wal-g-credentials
        - name: data
          persistentVolumeClaim:
            claimName: data-dictamesh-postgresql-0
EOF

# 3. Wait for restore job to complete
kubectl wait --for=condition=complete --timeout=600s \
  job/postgresql-restore -n dictamesh-system

# 4. Start PostgreSQL
kubectl scale statefulset dictamesh-postgresql \
  --namespace dictamesh-system \
  --replicas=1

# 5. Verify recovery
kubectl logs -f dictamesh-postgresql-0 -n dictamesh-system
```

### Restore from Logical Backup

```bash
# 1. Upload backup file to pod
kubectl cp ./backup-20251108.dump \
  dictamesh-system/dictamesh-postgresql-0:/tmp/restore.dump

# 2. Drop and recreate database (CAUTION!)
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  psql -U postgres -c "DROP DATABASE dictamesh_catalog;"

kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  psql -U postgres -c "CREATE DATABASE dictamesh_catalog OWNER dictamesh;"

# 3. Restore from dump
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  pg_restore -U dictamesh -d dictamesh_catalog \
  --clean \
  --if-exists \
  --no-owner \
  --no-acl \
  /tmp/restore.dump

# 4. Verify restoration
kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
  psql -U dictamesh -d dictamesh_catalog -c "\dt dictamesh_*"
```

### Kafka Topic Restore

```bash
# 1. Create topics if they don't exist
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- \
  kafka-topics.sh --bootstrap-server localhost:9092 \
  --create --topic dictamesh.entity.events \
  --partitions 12 --replication-factor 3

# 2. Download backup from S3
aws s3 cp s3://dictamesh-backups/kafka/kafka-20251108.tar.gz /tmp/

# 3. Extract backup
tar -xzf /tmp/kafka-20251108.tar.gz -C /tmp/

# 4. Restore messages to topics
kubectl exec -it dictamesh-kafka-0 -n dictamesh-system -- bash <<'EOF'
for file in /tmp/kafka-backup-20251108/*.json; do
    topic=$(basename "$file" .json)
    echo "Restoring topic: $topic"

    kafka-console-producer.sh \
        --bootstrap-server localhost:9092 \
        --topic "$topic" < "$file"
done
EOF
```

### Full Cluster Restore with Velero

```bash
# 1. List available backups
velero backup get

# 2. Restore entire namespace
velero restore create \
  --from-backup dictamesh-daily-backup-20251108 \
  --wait

# 3. Check restore status
velero restore describe dictamesh-daily-backup-20251108-20251108140000

# 4. Check logs if there are issues
velero restore logs dictamesh-daily-backup-20251108-20251108140000

# 5. Verify restored resources
kubectl get all -n dictamesh-system
kubectl get pvc -n dictamesh-system
kubectl get configmap,secret -n dictamesh-system
```

### Selective Restore

Restore specific resources only:

```bash
# Restore only ConfigMaps and Secrets
velero restore create \
  --from-backup dictamesh-daily-backup-20251108 \
  --include-resources configmaps,secrets \
  --namespace-mappings dictamesh-system:dictamesh-system-restored

# Restore specific PVCs
velero restore create \
  --from-backup dictamesh-daily-backup-20251108 \
  --include-resources pvc \
  --selector app=metadata-catalog
```

## Disaster Recovery Plan

### DR Architecture

```
┌─────────────────────────────────────┐
│   Primary Region (us-east-1)       │
│                                     │
│  ┌──────────────────────────────┐  │
│  │ DictaMesh Production         │  │
│  │ - Active workloads           │  │
│  │ - PostgreSQL primary         │  │
│  │ - Kafka cluster              │  │
│  └──────────────────────────────┘  │
│              │                      │
│              │ Continuous           │
│              │ Replication          │
│              ▼                      │
└─────────────────────────────────────┘
               │
               │
               ▼
┌─────────────────────────────────────┐
│   DR Region (us-west-2)             │
│                                     │
│  ┌──────────────────────────────┐  │
│  │ DictaMesh DR Standby         │  │
│  │ - Standby workloads (scaled) │  │
│  │ - PostgreSQL replica         │  │
│  │ - Kafka mirror               │  │
│  └──────────────────────────────┘  │
│                                     │
└─────────────────────────────────────┘
```

### DR Runbook

**Scenario: Primary region failure**

1. **Assess the situation** (5 minutes)
   ```bash
   # Check primary region status
   kubectl cluster-info --context prod-us-east-1

   # Check application health
   curl https://api.dictamesh.example.com/health
   ```

2. **Activate DR site** (10 minutes)
   ```bash
   # Switch kubectl context to DR region
   kubectl config use-context prod-us-west-2

   # Scale up DR workloads
   kubectl scale deployment dictamesh-graphql-gateway \
     --namespace dictamesh-system \
     --replicas=3

   kubectl scale statefulset dictamesh-metadata-catalog \
     --namespace dictamesh-system \
     --replicas=3
   ```

3. **Promote PostgreSQL replica** (10 minutes)
   ```bash
   # Promote standby to primary
   kubectl exec -it dictamesh-postgresql-0 -n dictamesh-system -- \
     pg_ctl promote -D /var/lib/postgresql/data
   ```

4. **Update DNS** (5 minutes)
   ```bash
   # Update DNS to point to DR region
   aws route53 change-resource-record-sets \
     --hosted-zone-id Z1234567890ABC \
     --change-batch file:///tmp/dns-failover.json
   ```

5. **Verify services** (10 minutes)
   ```bash
   # Test health endpoints
   curl https://api.dictamesh.example.com/health

   # Check pod status
   kubectl get pods -n dictamesh-system

   # Verify database connectivity
   kubectl exec -it dictamesh-metadata-catalog-0 -n dictamesh-system -- \
     /app/health-check database
   ```

6. **Monitor and communicate** (Ongoing)
   - Update status page
   - Notify stakeholders
   - Monitor metrics and logs

## Backup Verification

### Automated Backup Testing

Create a CronJob to test backups weekly:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: backup-verification
  namespace: dictamesh-system
spec:
  schedule: "0 4 * * 6"  # Weekly on Saturday at 4 AM
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: verify
              image: postgres:15
              command:
                - /bin/bash
                - -c
                - |
                  set -e

                  # Download latest backup
                  aws s3 cp s3://dictamesh-backups/latest.dump /tmp/verify.dump

                  # Create test database
                  createdb -h localhost -U postgres test_restore

                  # Restore to test database
                  pg_restore -h localhost -U postgres -d test_restore /tmp/verify.dump

                  # Run verification queries
                  psql -h localhost -U postgres -d test_restore -c "SELECT COUNT(*) FROM dictamesh_entity_catalog;"

                  # Clean up
                  dropdb -h localhost -U postgres test_restore

                  echo "Backup verification successful!"
```

## Best Practices

### Backup Strategy

✅ **Do:**
- Follow 3-2-1 rule: 3 copies, 2 different media, 1 off-site
- Test restores regularly
- Encrypt backups at rest and in transit
- Monitor backup jobs for failures
- Document restore procedures
- Automate backup verification

❌ **Don't:**
- Keep backups in the same region only
- Skip testing restore procedures
- Ignore backup failures
- Store backups without encryption
- Rely on a single backup method

### Disaster Recovery

✅ **Do:**
- Maintain updated DR runbooks
- Practice DR drills quarterly
- Monitor replication lag
- Have clear escalation paths
- Document RTO/RPO requirements

❌ **Don't:**
- Assume backups work without testing
- Skip DR drills
- Ignore replication delays
- Have unclear responsibilities
- Forget to update documentation

## Next Steps

- **[Monitoring](./monitoring.md)** - Monitor backup jobs
- **[Troubleshooting](./troubleshooting.md)** - Debug backup failures
- **[Configuration](./configuration.md)** - Configure backup retention policies

---

**Previous**: [← Scaling](./scaling.md) | **Next**: [Troubleshooting →](./troubleshooting.md)
