# Disaster Recovery

[← Previous: Monitoring & Alerting](17-MONITORING-ALERTING.md) | [Next: Migration Strategy →](19-MIGRATION-STRATEGY.md)

---

## 🎯 Purpose

Backup strategies, restore procedures, and disaster recovery testing.

---

## 💾 Backup Strategy

### PostgreSQL Backups (CloudNativePG)

```yaml
spec:
  backup:
    barmanObjectStore:
      destinationPath: s3://dictamesh-backups/postgres
      schedule: "0 0 * * *"  # Daily
      retentionPolicy: "30d"
```

### Kafka Backups

```bash
# Backup Kafka topic configuration
kubectl -n dictamesh-infra exec -it dictamesh-kafka-kafka-0 -- \
  kafka-configs.sh --describe --all > kafka-config-backup.txt
```

---

[← Previous: Monitoring & Alerting](17-MONITORING-ALERTING.md) | [Next: Migration Strategy →](19-MIGRATION-STRATEGY.md)
