# Implementation Phases

[← Previous: Architecture Overview](01-ARCHITECTURE-OVERVIEW.md) | [Next: Infrastructure Planning →](03-INFRASTRUCTURE-PLANNING.md)

---

## 🎯 Purpose

This document provides LLM agents with a phased implementation roadmap, ensuring systematic delivery with minimal risk and maximum learning opportunities.

**Reading Time:** 20 minutes
**Prerequisites:** [Architecture Overview](01-ARCHITECTURE-OVERVIEW.md)
**Outputs:** Clear execution sequence, milestone definitions, dependency mapping

---

## 📅 Phase Overview

```
Phase 0: Foundation (Weeks 1-2)
   ↓
Phase 1: Infrastructure Bootstrap (Weeks 3-4)
   ↓
Phase 2: Core Integration Layer (Weeks 5-8)
   ↓
Phase 3: First Data Product (Weeks 9-10)
   ↓
Phase 4: Federation & API Gateway (Weeks 11-12)
   ↓
Phase 5: Remaining Data Products (Weeks 13-15)
   ↓
Phase 6: Advanced Features (Weeks 16-18)
   ↓
Phase 7: Production Hardening (Weeks 19-20)
   ↓
Phase 8: Go-Live & Monitoring (Week 21)
```

---

## 🔧 Phase 0: Foundation (Weeks 1-2)

### Objectives
- Establish development environment
- Set up version control and branching strategy
- Configure K3S cluster access
- Initialize project structure

### Tasks for LLM Agents

#### Task 0.1: Repository Setup
```bash
# Initialize monorepo structure
mkdir -p {services,infrastructure,docs,tools}
mkdir -p services/{customer-adapter,product-adapter,invoice-adapter,metadata-catalog,graphql-gateway}
mkdir -p infrastructure/{k8s,terraform,argocd}

# Initialize Go modules
cd services/customer-adapter && go mod init github.com/dictamesh/customer-adapter
cd services/product-adapter && go mod init github.com/dictamesh/product-adapter
cd services/invoice-adapter && go mod init github.com/dictamesh/invoice-adapter
cd services/metadata-catalog && go mod init github.com/dictamesh/metadata-catalog
cd services/graphql-gateway && go mod init github.com/dictamesh/graphql-gateway
```

#### Task 0.2: K3S Cluster Validation
```bash
# Verify cluster access
kubectl cluster-info
kubectl get nodes
kubectl get namespaces

# Test namespace creation
kubectl create namespace dictamesh-dev --dry-run=client -o yaml
```

#### Task 0.3: Development Tools Setup
```yaml
# .devcontainer/devcontainer.json
{
  "name": "DictaMesh Development",
  "image": "mcr.microsoft.com/devcontainers/go:1.21",
  "features": {
    "ghcr.io/devcontainers/features/kubectl-helm-minikube:1": {},
    "ghcr.io/devcontainers/features/docker-in-docker:2": {}
  },
  "customizations": {
    "vscode": {
      "extensions": [
        "golang.go",
        "ms-kubernetes-tools.vscode-kubernetes-tools",
        "redhat.vscode-yaml"
      ]
    }
  }
}
```

#### Task 0.4: CI/CD Repository Setup
```bash
# Create GitHub Actions workflows directory
mkdir -p .github/workflows

# Create ArgoCD application repository
mkdir -p infrastructure/argocd/applications
```

### Deliverables
- [ ] Monorepo structure initialized
- [ ] Git repository configured with branching strategy (main, develop, feature/*)
- [ ] K3S cluster accessible via kubectl
- [ ] Development environment documented
- [ ] Team onboarding guide created

### Success Criteria
```bash
# All checks pass
kubectl get nodes | grep Ready
go version | grep "go1.21"
docker --version
helm version
```

---

## 🏗️ Phase 1: Infrastructure Bootstrap (Weeks 3-4)

### Objectives
- Deploy Kubernetes infrastructure components
- Set up observability stack
- Configure networking and ingress
- Establish GitOps with ArgoCD

### Tasks for LLM Agents

#### Task 1.1: Namespace Creation
```yaml
# infrastructure/k8s/namespaces/namespaces.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-dev
  labels:
    environment: dev
---
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-staging
  labels:
    environment: staging
---
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-prod
  labels:
    environment: prod
---
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-infra
  labels:
    environment: shared
---
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-monitoring
  labels:
    environment: shared
---
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-cicd
  labels:
    environment: shared
```

```bash
kubectl apply -f infrastructure/k8s/namespaces/namespaces.yaml
```

#### Task 1.2: ArgoCD Installation
```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for rollout
kubectl wait --for=condition=available --timeout=600s deployment/argocd-server -n argocd

# Get initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# Expose ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

#### Task 1.3: Observability Stack Deployment
```bash
# Add Helm repos
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Install Prometheus
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace dictamesh-monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.retention=30d \
  --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage=50Gi

# Install Loki
helm install loki grafana/loki-stack \
  --namespace dictamesh-monitoring \
  --set promtail.enabled=true \
  --set grafana.enabled=false

# Install Jaeger
kubectl create namespace observability
kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.51.0/jaeger-operator.yaml -n observability
```

#### Task 1.4: Ingress Configuration
```yaml
# infrastructure/k8s/ingress/argocd-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: argocd-ingress
  namespace: argocd
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - argocd.dictamesh.controle.digital
      secretName: argocd-tls
  rules:
    - host: argocd.dictamesh.controle.digital
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: argocd-server
                port:
                  number: 443
```

### Deliverables
- [ ] All namespaces created and labeled
- [ ] ArgoCD operational and accessible
- [ ] Prometheus + Grafana deployed
- [ ] Loki logging stack operational
- [ ] Jaeger tracing ready
- [ ] Ingress controller configured
- [ ] SSL certificates configured (cert-manager)

### Success Criteria
```bash
# Verify all systems operational
kubectl get pods -n dictamesh-monitoring | grep Running
kubectl get pods -n argocd | grep Running
kubectl get pods -n observability | grep Running

# Access UIs
# - ArgoCD: https://argocd.dictamesh.controle.digital
# - Grafana: kubectl port-forward -n dictamesh-monitoring svc/prometheus-grafana 3000:80
# - Jaeger: kubectl port-forward -n observability svc/jaeger-query 16686:16686
```

---

## 🔗 Phase 2: Core Integration Layer (Weeks 5-8)

### Objectives
- Deploy PostgreSQL for Metadata Catalog
- Deploy Kafka cluster for event streaming
- Deploy Redis for caching
- Set up Schema Registry

### Tasks for LLM Agents

#### Task 2.1: PostgreSQL Deployment
```bash
# Using CloudNativePG operator
kubectl apply -f https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.21/releases/cnpg-1.21.0.yaml

# Create PostgreSQL cluster
cat <<EOF | kubectl apply -f -
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: metadata-catalog-db
  namespace: dictamesh-infra
spec:
  instances: 3
  imageName: ghcr.io/cloudnative-pg/postgresql:15.5
  storage:
    size: 100Gi
    storageClass: longhorn
  postgresql:
    parameters:
      max_connections: "200"
      shared_buffers: "2GB"
      effective_cache_size: "6GB"
      work_mem: "16MB"
  monitoring:
    enabled: true
  backup:
    barmanObjectStore:
      destinationPath: s3://dictamesh-backups/postgres
      s3Credentials:
        accessKeyId:
          name: s3-credentials
          key: ACCESS_KEY_ID
        secretAccessKey:
          name: s3-credentials
          key: ACCESS_SECRET_KEY
EOF
```

#### Task 2.2: Kafka Deployment (Strimzi)
```bash
# Install Strimzi operator
kubectl create namespace kafka
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka

# Deploy Kafka cluster
cat <<EOF | kubectl apply -f -
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
      transaction.state.log.min.isr: 2
      default.replication.factor: 3
      min.insync.replicas: 2
      inter.broker.protocol.version: "3.6"
    storage:
      type: persistent-claim
      size: 100Gi
      class: longhorn
    resources:
      requests:
        memory: 4Gi
        cpu: "1"
      limits:
        memory: 8Gi
        cpu: "2"
    metricsConfig:
      type: jmxPrometheusExporter
      valueFrom:
        configMapKeyRef:
          name: kafka-metrics
          key: kafka-metrics-config.yml
  zookeeper:
    replicas: 3
    storage:
      type: persistent-claim
      size: 20Gi
      class: longhorn
    resources:
      requests:
        memory: 1Gi
        cpu: "0.5"
      limits:
        memory: 2Gi
        cpu: "1"
  entityOperator:
    topicOperator: {}
    userOperator: {}
EOF
```

#### Task 2.3: Schema Registry Deployment
```yaml
# infrastructure/k8s/kafka/schema-registry.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: schema-registry
  namespace: dictamesh-infra
spec:
  replicas: 2
  selector:
    matchLabels:
      app: schema-registry
  template:
    metadata:
      labels:
        app: schema-registry
    spec:
      containers:
        - name: schema-registry
          image: confluentinc/cp-schema-registry:7.5.0
          ports:
            - containerPort: 8081
          env:
            - name: SCHEMA_REGISTRY_HOST_NAME
              value: schema-registry
            - name: SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS
              value: dictamesh-kafka-kafka-bootstrap.dictamesh-infra.svc:9092
            - name: SCHEMA_REGISTRY_LISTENERS
              value: http://0.0.0.0:8081
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: schema-registry
  namespace: dictamesh-infra
spec:
  ports:
    - port: 8081
      targetPort: 8081
  selector:
    app: schema-registry
```

#### Task 2.4: Redis Deployment
```bash
# Install Redis operator
helm repo add ot-helm https://ot-container-kit.github.io/helm-charts/
helm install redis-operator ot-helm/redis-operator --namespace dictamesh-infra

# Deploy Redis cluster
cat <<EOF | kubectl apply -f -
apiVersion: redis.redis.opstreelabs.in/v1beta1
kind: RedisCluster
metadata:
  name: dictamesh-cache
  namespace: dictamesh-infra
spec:
  clusterSize: 3
  kubernetesConfig:
    image: redis:7.2
    imagePullPolicy: IfNotPresent
  redisExporter:
    enabled: true
    image: quay.io/opstree/redis-exporter:v1.44.0
  storage:
    volumeClaimTemplate:
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
        storageClassName: longhorn
EOF
```

### Deliverables
- [ ] PostgreSQL cluster operational (3 replicas)
- [ ] Kafka cluster operational (3 brokers)
- [ ] Schema Registry deployed
- [ ] Redis cluster operational
- [ ] All services monitored in Prometheus
- [ ] Connection secrets created and documented

### Success Criteria
```bash
# PostgreSQL
kubectl get cluster -n dictamesh-infra metadata-catalog-db
kubectl get pods -n dictamesh-infra | grep metadata-catalog-db | grep Running

# Kafka
kubectl get kafka -n dictamesh-infra dictamesh-kafka
kubectl get pods -n dictamesh-infra | grep dictamesh-kafka | grep Running

# Schema Registry
kubectl get pods -n dictamesh-infra | grep schema-registry | grep Running
curl http://schema-registry.dictamesh-infra.svc:8081/subjects

# Redis
kubectl get rediscluster -n dictamesh-infra dictamesh-cache
kubectl exec -it dictamesh-cache-0 -n dictamesh-infra -- redis-cli ping
```

---

## 📦 Phase 3: First Data Product (Weeks 9-10)

### Objectives
- Implement Customer Adapter (Directus integration)
- Deploy Metadata Catalog Service
- Establish end-to-end data flow
- Validate event-driven architecture

### Tasks for LLM Agents

#### Task 3.1: Metadata Catalog Service Implementation
```bash
cd services/metadata-catalog

# Initialize project structure
mkdir -p {cmd/server,internal/{catalog,consumer,repository},pkg/models,migrations}

# Create main.go
cat > cmd/server/main.go << 'EOF'
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/dictamesh/metadata-catalog/internal/catalog"
    "github.com/dictamesh/metadata-catalog/internal/consumer"
    "github.com/dictamesh/metadata-catalog/internal/repository"
)

func main() {
    // Initialize database connection
    db, err := repository.NewPostgresDB(os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Run migrations
    if err := repository.RunMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Initialize catalog service
    catalogService := catalog.NewService(db)

    // Start Kafka consumer
    kafkaConsumer := consumer.NewKafkaConsumer(
        os.Getenv("KAFKA_BROKERS"),
        catalogService,
    )
    go kafkaConsumer.Start(context.Background())

    // Start HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: catalog.NewHTTPHandler(catalogService),
    }

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    go func() {
        log.Println("Server starting on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    <-stop
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown error: %v", err)
    }

    log.Println("Server stopped")
}
EOF
```

#### Task 3.2: Customer Adapter Implementation
```bash
cd services/customer-adapter

# See detailed implementation in 06-LAYER1-ADAPTERS.md

# Key files to create:
# - cmd/server/main.go
# - internal/adapter/directus_client.go
# - internal/adapter/event_publisher.go
# - internal/adapter/cache_layer.go
# - deployments/k8s/deployment.yaml
```

#### Task 3.3: Kafka Topic Creation
```bash
# Create topic for customer events
kubectl -n dictamesh-infra exec -it dictamesh-kafka-kafka-0 -- bin/kafka-topics.sh \
  --create \
  --topic customers.directus.entity_changed \
  --bootstrap-server localhost:9092 \
  --partitions 12 \
  --replication-factor 3 \
  --config retention.ms=2592000000 \
  --config cleanup.policy=delete
```

#### Task 3.4: Schema Registration
```bash
# Register Avro schema for customer events
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data @schemas/customer-change-event.avsc \
  http://schema-registry.dictamesh-infra.svc:8081/subjects/customers.directus.entity_changed-value/versions
```

### Deliverables
- [ ] Metadata Catalog Service deployed and operational
- [ ] Customer Adapter deployed and connected to Directus
- [ ] Kafka topic created for customer events
- [ ] Avro schema registered
- [ ] End-to-end event flow validated (Directus → Adapter → Kafka → Catalog)
- [ ] Monitoring dashboards configured
- [ ] Integration tests passing

### Success Criteria
```bash
# Trigger customer change in Directus
# Verify event published to Kafka
kubectl -n dictamesh-infra exec -it dictamesh-kafka-kafka-0 -- bin/kafka-console-consumer.sh \
  --topic customers.directus.entity_changed \
  --bootstrap-server localhost:9092 \
  --from-beginning \
  --max-messages 1

# Verify metadata catalog received event
kubectl logs -n dictamesh-dev deployment/metadata-catalog | grep "Entity registered"

# Query catalog API
curl http://metadata-catalog.dictamesh-dev.svc:8080/api/v1/entities/customer/123
```

---

## 🌐 Phase 4: Federation & API Gateway (Weeks 11-12)

### Objectives
- Deploy GraphQL Gateway
- Implement Customer subgraph
- Set up federation
- Validate federated queries

### Tasks for LLM Agents

Detailed in [09-LAYER4-API-GATEWAY.md](09-LAYER4-API-GATEWAY.md)

### Deliverables
- [ ] GraphQL Gateway deployed
- [ ] Customer subgraph operational
- [ ] Federation configuration complete
- [ ] DataLoader implemented for batching
- [ ] GraphQL Playground accessible
- [ ] Query performance benchmarked

---

## 📊 Phase 5: Remaining Data Products (Weeks 13-15)

### Objectives
- Implement Product Adapter
- Implement Invoice Adapter
- Add Product and Invoice subgraphs
- Validate cross-domain queries

### Tasks for LLM Agents

Parallel implementation of:
- Product Adapter (see [06-LAYER1-ADAPTERS.md](06-LAYER1-ADAPTERS.md))
- Invoice Adapter (see [06-LAYER1-ADAPTERS.md](06-LAYER1-ADAPTERS.md))
- Product Subgraph (see [09-LAYER4-API-GATEWAY.md](09-LAYER4-API-GATEWAY.md))
- Invoice Subgraph (see [09-LAYER4-API-GATEWAY.md](09-LAYER4-API-GATEWAY.md))

### Deliverables
- [ ] Product Adapter operational
- [ ] Invoice Adapter operational
- [ ] All Kafka topics created
- [ ] All schemas registered
- [ ] Product and Invoice subgraphs federated
- [ ] Cross-domain queries working (e.g., customer.invoices.items.product)

---

## 🚀 Phase 6: Advanced Features (Weeks 16-18)

### Objectives
- Implement multi-tenancy
- Add Saga orchestration
- Configure advanced caching
- Set up distributed tracing

### Tasks for LLM Agents

See detailed guides:
- [11-LAYER6-MULTITENANCY.md](11-LAYER6-MULTITENANCY.md)
- [12-LAYER7-SAGA-ORCHESTRATION.md](12-LAYER7-SAGA-ORCHESTRATION.md)
- [10-LAYER5-OBSERVABILITY.md](10-LAYER5-OBSERVABILITY.md)

---

## 🛡️ Phase 7: Production Hardening (Weeks 19-20)

### Objectives
- Security hardening
- Performance optimization
- Chaos testing
- Documentation completion

### Tasks for LLM Agents

#### Task 7.1: Security Audit
- [ ] Enable network policies
- [ ] Configure Pod Security Standards
- [ ] Set up Vault for secrets management
- [ ] Enable audit logging
- [ ] Scan container images (Trivy)
- [ ] Penetration testing

#### Task 7.2: Performance Optimization
- [ ] Load testing with k6
- [ ] Database query optimization
- [ ] Cache hit rate optimization
- [ ] Kafka tuning
- [ ] Resource limit calibration

#### Task 7.3: Chaos Engineering
```bash
# Install Chaos Mesh
kubectl create namespace chaos-testing
helm install chaos-mesh chaos-mesh/chaos-mesh -n chaos-testing

# Run pod kill experiment
cat <<EOF | kubectl apply -f -
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: pod-failure-example
  namespace: dictamesh-dev
spec:
  action: pod-failure
  mode: one
  duration: '30s'
  selector:
    namespaces:
      - dictamesh-dev
    labelSelectors:
      'app': 'customer-adapter'
EOF
```

### Deliverables
- [ ] Security audit passed
- [ ] Performance benchmarks documented
- [ ] Chaos tests passing
- [ ] Disaster recovery tested
- [ ] Documentation complete

---

## ✅ Phase 8: Go-Live & Monitoring (Week 21)

### Objectives
- Production deployment
- Monitoring validation
- On-call setup
- Knowledge transfer

### Tasks for LLM Agents

#### Task 8.1: Production Deployment
```bash
# Deploy to production namespace via ArgoCD
kubectl apply -f infrastructure/argocd/applications/production.yaml

# Verify deployment
argocd app sync dictamesh-production
argocd app wait dictamesh-production --health
```

#### Task 8.2: Smoke Tests
```bash
# Run production smoke tests
kubectl apply -f tests/smoke/production-smoke-tests.yaml
kubectl wait --for=condition=complete job/production-smoke-tests -n dictamesh-prod
```

#### Task 8.3: Documentation Handoff
- [ ] Usage documentation complete (see [14-DOCUMENTATION-PLANNING.md](14-DOCUMENTATION-PLANNING.md))
- [ ] Administration runbooks complete
- [ ] Troubleshooting guides complete
- [ ] Development onboarding complete
- [ ] Contribution guidelines complete

### Success Criteria
- [ ] All services healthy in production
- [ ] Monitoring alerts configured
- [ ] On-call runbooks tested
- [ ] SLOs defined and measured
- [ ] Stakeholders trained

---

## 📊 Progress Tracking Template

Create `IMPLEMENTATION-STATUS.md`:

```markdown
# DictaMesh Implementation Status

Last Updated: 2025-11-08

## Phase 0: Foundation ✅
- [x] Repository setup
- [x] K3S cluster access
- [x] Development environment

## Phase 1: Infrastructure Bootstrap 🟡
- [x] Namespaces created
- [x] ArgoCD installed
- [ ] Observability stack (80% complete)

## Phase 2: Core Integration Layer 🔴
- [ ] PostgreSQL deployment (not started)
- [ ] Kafka cluster (not started)

... continue for all phases ...
```

---

## 🔄 Rollback Procedures

### Phase Rollback Protocol

```bash
# Rollback to previous phase
argocd app rollback <app-name> <revision-id>

# Rollback database migrations
migrate -database $DATABASE_URL -path migrations down 1

# Rollback Helm releases
helm rollback <release-name> <revision>
```

---

[← Previous: Architecture Overview](01-ARCHITECTURE-OVERVIEW.md) | [Next: Infrastructure Planning →](03-INFRASTRUCTURE-PLANNING.md)

---

**Document Metadata**
- Version: 1.0.0
- Last Updated: 2025-11-08
- Next Review: After Phase 2 completion
