# Layer 2: Event-Driven Integration Fabric

[← Previous: Layer 1 Adapters](06-LAYER1-ADAPTERS.md) | [Next: Layer 3 Metadata Catalog →](08-LAYER3-METADATA-CATALOG.md)

---

## 🎯 Purpose

Implementation guide for Apache Kafka event bus, Schema Registry, and event-driven communication patterns.

---

## 🔧 Kafka Cluster Setup

### Using Strimzi Operator on K3S

```bash
# Install Strimzi
kubectl create namespace kafka
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka

# Deploy Kafka cluster
kubectl apply -f infrastructure/k8s/kafka/kafka-cluster.yaml
```

```yaml
# infrastructure/k8s/kafka/kafka-cluster.yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: dictamesh-kafka
  namespace: dictamesh-infra
spec:
  kafka:
    version: 3.6.0
    replicas: 3
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
    config:
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      default.replication.factor: 3
    storage:
      type: persistent-claim
      size: 100Gi
      class: longhorn
  zookeeper:
    replicas: 3
    storage:
      type: persistent-claim
      size: 20Gi
```

---

## 📝 Topic Design

### Topic Naming Convention

Pattern: `<domain>.<source>.<event_type>`

```yaml
# infrastructure/k8s/kafka/topics.yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: customers.directus.entity-changed
  namespace: dictamesh-infra
  labels:
    strimzi.io/cluster: dictamesh-kafka
spec:
  partitions: 12
  replicas: 3
  config:
    retention.ms: 2592000000  # 30 days
    cleanup.policy: delete
    compression.type: lz4
```

---

## 🔐 Schema Registry

```yaml
# infrastructure/k8s/kafka/schema-registry.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: schema-registry
  namespace: dictamesh-infra
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: schema-registry
          image: confluentinc/cp-schema-registry:7.5.0
          env:
            - name: SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS
              value: dictamesh-kafka-kafka-bootstrap:9092
```

### Avro Schema Registration

```bash
# Register customer change event schema
curl -X POST http://schema-registry:8081/subjects/customers.directus.entity_changed-value/versions \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data @schemas/customer-change-event.avsc
```

---

[← Previous: Layer 1 Adapters](06-LAYER1-ADAPTERS.md) | [Next: Layer 3 Metadata Catalog →](08-LAYER3-METADATA-CATALOG.md)
