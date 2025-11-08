# Documentation Planning

[← Previous: Testing Strategy](13-TESTING-STRATEGY.md) | [Next: Security & Compliance →](15-SECURITY-COMPLIANCE.md)

---

## 🎯 Purpose

This document provides LLM agents with a comprehensive plan for creating all necessary documentation for DictaMesh, including Usage, Administration, Troubleshooting, Development, and Contribution guides.

**Reading Time:** 20 minutes
**Prerequisites:** All implementation documents
**Outputs:** Complete documentation structure, templates, generation strategies

---

## 📚 Documentation Categories

### Overview of Required Documentation

```
docs/
├── usage/                      # End-user documentation
│   ├── getting-started.md
│   ├── graphql-api-guide.md
│   ├── querying-data.md
│   ├── authentication.md
│   └── examples/
│
├── administration/             # Operations and administration
│   ├── installation.md
│   ├── configuration.md
│   ├── scaling.md
│   ├── backup-restore.md
│   ├── upgrading.md
│   └── monitoring.md
│
├── troubleshooting/            # Problem resolution
│   ├── common-issues.md
│   ├── debugging-guide.md
│   ├── performance-tuning.md
│   ├── error-reference.md
│   └── runbooks/
│       ├── kafka-issues.md
│       ├── postgres-issues.md
│       └── adapter-failures.md
│
├── development/                # Developer documentation
│   ├── architecture.md
│   ├── local-setup.md
│   ├── adding-adapters.md
│   ├── api-reference/
│   ├── testing-guide.md
│   └── code-style.md
│
├── contribution/               # Contributing guidelines
│   ├── CONTRIBUTING.md
│   ├── code-of-conduct.md
│   ├── pull-request-template.md
│   ├── issue-templates/
│   └── development-workflow.md
│
└── reference/                  # Technical reference
    ├── api/
    ├── configuration-reference.md
    ├── metrics-reference.md
    └── events-reference.md
```

---

## 📖 1. Usage Documentation

### 1.1 Getting Started Guide

**File:** `docs/usage/getting-started.md`

**LLM Agent Template:**

```markdown
# Getting Started with DictaMesh

## What is DictaMesh?

DictaMesh is an enterprise-grade data mesh platform that provides unified access to distributed data sources through a federated GraphQL API.

## Quick Start

### Prerequisites
- Valid API credentials
- GraphQL client (Postman, Insomnia, or curl)

### Your First Query

#### Step 1: Authenticate

\`\`\`bash
curl -X POST https://api.dictamesh.controle.digital/auth/login \\
  -H "Content-Type: application/json" \\
  -d '{
    "username": "your-username",
    "password": "your-password"
  }'
\`\`\`

Response:
\`\`\`json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
\`\`\`

#### Step 2: Query Customer Data

\`\`\`graphql
query GetCustomer {
  customer(id: "cust-123") {
    id
    name
    email
    invoices {
      invoiceNumber
      total
      items {
        product {
          name
          price
        }
        quantity
      }
    }
  }
}
\`\`\`

#### Step 3: Execute Query

\`\`\`bash
curl -X POST https://api.dictamesh.controle.digital/graphql \\
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{
    "query": "query GetCustomer { customer(id: \"cust-123\") { id name email } }"
  }'
\`\`\`

## What's Next?

- [GraphQL API Guide](graphql-api-guide.md) - Learn advanced querying
- [Authentication Guide](authentication.md) - Configure OAuth/API keys
- [Examples](examples/) - See real-world use cases
```

### 1.2 GraphQL API Guide

**File:** `docs/usage/graphql-api-guide.md`

**Content Outline:**
- GraphQL basics for DictaMesh
- Schema exploration (using GraphQL Playground)
- Query patterns (filtering, pagination, sorting)
- Mutations (if write operations supported)
- Subscriptions (real-time data)
- Error handling
- Best practices
- Rate limiting

### 1.3 Code Examples

**File:** `docs/usage/examples/python-client.md`

```markdown
# Python Client Example

## Installation

\`\`\`bash
pip install gql[all]
\`\`\`

## Basic Usage

\`\`\`python
from gql import gql, Client
from gql.transport.requests import RequestsHTTPTransport

# Configure transport
transport = RequestsHTTPTransport(
    url="https://api.dictamesh.controle.digital/graphql",
    headers={"Authorization": "Bearer YOUR_TOKEN"},
)

# Create client
client = Client(transport=transport, fetch_schema_from_transport=True)

# Execute query
query = gql("""
    query GetCustomerWithInvoices($customerId: ID!) {
        customer(id: $customerId) {
            name
            email
            invoices {
                invoiceNumber
                total
            }
        }
    }
""")

result = client.execute(query, variable_values={"customerId": "cust-123"})
print(result)
\`\`\`
```

**Additional Examples:**
- JavaScript/TypeScript client
- Go client
- Java client
- cURL examples
- Pagination examples
- Batch queries

---

## 🛠️ 2. Administration Documentation

### 2.1 Installation Guide

**File:** `docs/administration/installation.md`

**Content Structure:**

```markdown
# DictaMesh Installation Guide

## Prerequisites

### Infrastructure Requirements
- Kubernetes cluster (K3S 1.28+)
- Storage: 200Gi+ (dev), 900Gi+ (prod)
- CPU: 10+ cores (dev), 60+ cores (prod)
- Memory: 20Gi+ (dev), 100Gi+ (prod)

### Required Tools
- kubectl 1.28+
- helm 3.10+
- argocd CLI 2.8+

## Installation Methods

### Method 1: GitOps (Recommended)

#### Step 1: Install ArgoCD
\`\`\`bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
\`\`\`

#### Step 2: Deploy DictaMesh
\`\`\`bash
kubectl apply -f https://raw.githubusercontent.com/controle-digital/dictamesh/main/infrastructure/argocd/applications/prod/app-of-apps.yaml
\`\`\`

#### Step 3: Verify Installation
\`\`\`bash
argocd app list
argocd app wait dictamesh-prod --health
\`\`\`

### Method 2: Helm Charts
[... detailed Helm installation ...]

### Method 3: Manual kubectl
[... manual installation steps ...]

## Post-Installation Configuration

### Configure External Access
[... ingress setup ...]

### Set Up Monitoring
[... prometheus/grafana setup ...]

### Configure Backups
[... backup configuration ...]

## Validation

### Health Checks
\`\`\`bash
kubectl get pods -n dictamesh-prod
kubectl get svc -n dictamesh-prod
\`\`\`

### Smoke Tests
\`\`\`bash
curl https://api.dictamesh.controle.digital/health
\`\`\`
```

### 2.2 Configuration Reference

**File:** `docs/administration/configuration.md`

**Content:**
- Environment variables
- ConfigMaps structure
- Secrets management
- Feature flags
- Performance tuning parameters
- Multi-tenancy configuration

### 2.3 Backup & Restore

**File:** `docs/administration/backup-restore.md`

```markdown
# Backup and Restore

## Backup Strategy

### PostgreSQL Backups

#### Automated Backups (CloudNativePG)
\`\`\`yaml
spec:
  backup:
    barmanObjectStore:
      destinationPath: s3://dictamesh-backups/postgres
      schedule: "0 0 * * *"  # Daily at midnight
      retentionPolicy: "30d"
\`\`\`

#### Manual Backup
\`\`\`bash
kubectl cnpg backup metadata-catalog-db -n dictamesh-infra
\`\`\`

### Kafka Backups

#### Topic Configuration Backup
\`\`\`bash
kubectl -n dictamesh-infra exec -it dictamesh-kafka-kafka-0 -- \\
  kafka-configs.sh --bootstrap-server localhost:9092 --describe --all > kafka-config-backup.txt
\`\`\`

### Configuration Backups

#### Export All ConfigMaps and Secrets
\`\`\`bash
kubectl get configmaps -n dictamesh-prod -o yaml > configmaps-backup.yaml
kubectl get secrets -n dictamesh-prod -o yaml > secrets-backup.yaml
\`\`\`

## Restore Procedures

### Restore PostgreSQL from Backup
[... detailed restore steps ...]

### Restore Kafka Topics
[... detailed restore steps ...]

## Disaster Recovery Testing

Schedule quarterly DR tests:
\`\`\`bash
./scripts/disaster-recovery-test.sh
\`\`\`
```

### 2.4 Monitoring & Alerting

**File:** `docs/administration/monitoring.md`

**Content:**
- Grafana dashboard access
- Key metrics to watch
- Alert rules explanation
- Log aggregation setup
- Distributed tracing access
- Performance benchmarks

### 2.5 Scaling Guide

**File:** `docs/administration/scaling.md`

**Content:**
- Horizontal scaling (HPA configuration)
- Vertical scaling (resource adjustments)
- Database scaling (read replicas)
- Kafka scaling (adding brokers)
- Cache scaling (Redis cluster expansion)
- Load testing procedures

---

## 🔧 3. Troubleshooting Documentation

### 3.1 Common Issues

**File:** `docs/troubleshooting/common-issues.md`

```markdown
# Common Issues and Solutions

## Issue: GraphQL Gateway Returns 503 Service Unavailable

### Symptoms
- API returns HTTP 503
- Logs show "connection refused" errors
- Unable to query data

### Diagnosis
\`\`\`bash
# Check gateway pods
kubectl get pods -n dictamesh-prod -l app=graphql-gateway

# Check gateway logs
kubectl logs -n dictamesh-prod deployment/graphql-gateway --tail=100

# Check adapter health
kubectl get pods -n dictamesh-prod -l component=adapter
\`\`\`

### Common Causes
1. **Adapter pods not ready**
   - Solution: Check adapter logs, verify source system connectivity
2. **Network policy blocking traffic**
   - Solution: Review network policies, ensure correct labels
3. **Circuit breaker open**
   - Solution: Check adapter metrics, verify source system health

### Resolution Steps
1. Verify all dependencies are healthy
2. Check network connectivity
3. Review recent deployments
4. Check resource limits (CPU/memory throttling)

---

## Issue: Kafka Consumer Lag Increasing

### Symptoms
- `kafka_consumer_lag` metric increasing
- Events not being processed
- Metadata catalog out of sync

### Diagnosis
\`\`\`bash
# Check consumer group lag
kubectl -n dictamesh-infra exec -it dictamesh-kafka-kafka-0 -- \\
  kafka-consumer-groups.sh --bootstrap-server localhost:9092 \\
  --describe --group metadata-catalog-consumer
\`\`\`

### Solutions
[... troubleshooting steps ...]
```

### 3.2 Debugging Guide

**File:** `docs/troubleshooting/debugging-guide.md`

**Content:**
- Enabling debug logging
- Accessing container logs
- Using port-forwarding for local debugging
- Trace ID propagation
- Database query debugging
- Network traffic inspection

### 3.3 Runbooks

**File:** `docs/troubleshooting/runbooks/kafka-issues.md`

```markdown
# Runbook: Kafka Issues

## Scenario: Broker Down

### Detection
- Alert: `kafka_broker_down`
- Symptom: 1 of 3 brokers not responding

### Impact Assessment
- [ ] Check replication factor: Should have 2 remaining replicas
- [ ] Check under-replicated partitions
- [ ] Verify producer/consumer functionality

### Immediate Actions
1. Check broker pod status
   \`\`\`bash
   kubectl get pods -n dictamesh-infra | grep kafka
   \`\`\`

2. Check broker logs
   \`\`\`bash
   kubectl logs -n dictamesh-infra dictamesh-kafka-kafka-0 --tail=200
   \`\`\`

3. If pod crashed, check resource limits
   \`\`\`bash
   kubectl describe pod -n dictamesh-infra dictamesh-kafka-kafka-0
   \`\`\`

### Recovery Steps
[... detailed recovery procedure ...]

### Post-Incident
- [ ] Document root cause
- [ ] Update monitoring/alerting if needed
- [ ] Schedule post-mortem
```

**Additional Runbooks:**
- PostgreSQL connection pool exhaustion
- Redis cache failures
- Adapter source system timeout
- High API latency investigation
- Certificate expiration

### 3.4 Error Reference

**File:** `docs/troubleshooting/error-reference.md`

**Format:**

| Error Code | Error Message | Cause | Solution |
|------------|---------------|-------|----------|
| `ADAPTER_001` | Source system unreachable | Network connectivity or source down | Check network policies, verify source system health |
| `CATALOG_002` | Entity not found in catalog | Entity never registered or deleted | Verify entity exists in source, check adapter logs |
| `GATEWAY_003` | Resolver timeout | Slow adapter response | Check adapter performance, increase timeout |

---

## 💻 4. Development Documentation

### 4.1 Architecture Documentation

**File:** `docs/development/architecture.md`

**Content:**
- System architecture diagrams
- Component interactions
- Data flow diagrams
- Technology stack
- Design decisions (ADRs)
- Scalability considerations

### 4.2 Local Development Setup

**File:** `docs/development/local-setup.md`

```markdown
# Local Development Setup

## Prerequisites
- Docker Desktop or Podman
- Go 1.21+
- Node.js 20+ (for GraphQL Playground)
- kubectl + k3d (for local Kubernetes)

## Quick Start with Docker Compose

### Step 1: Clone Repository
\`\`\`bash
git clone https://github.com/controle-digital/dictamesh.git
cd dictamesh
\`\`\`

### Step 2: Start Infrastructure
\`\`\`bash
docker-compose -f docker-compose.dev.yaml up -d
\`\`\`

This starts:
- PostgreSQL (localhost:5432)
- Kafka + Zookeeper (localhost:9092)
- Redis (localhost:6379)
- Schema Registry (localhost:8081)

### Step 3: Run Database Migrations
\`\`\`bash
cd services/metadata-catalog
make migrate-up
\`\`\`

### Step 4: Start Services
\`\`\`bash
# Terminal 1: Customer Adapter
cd services/customer-adapter
go run cmd/server/main.go

# Terminal 2: Metadata Catalog
cd services/metadata-catalog
go run cmd/server/main.go

# Terminal 3: GraphQL Gateway
cd services/graphql-gateway
go run cmd/server/main.go
\`\`\`

### Step 5: Access GraphQL Playground
Open http://localhost:8080/playground

## Local Kubernetes Setup (k3d)

\`\`\`bash
# Create local cluster
k3d cluster create dictamesh-local --agents 2

# Deploy local version
kubectl apply -k infrastructure/k8s/overlays/local
\`\`\`

## Hot Reload Development

Use Air for Go hot reload:
\`\`\`bash
go install github.com/cosmtrek/air@latest
cd services/customer-adapter
air
\`\`\`
```

### 4.3 Adding New Adapters

**File:** `docs/development/adding-adapters.md`

```markdown
# Adding a New Adapter

## Step 1: Scaffold Adapter Project

\`\`\`bash
./scripts/scaffold-adapter.sh <adapter-name> <domain>
# Example: ./scripts/scaffold-adapter.sh supplier-adapter suppliers
\`\`\`

This generates:
\`\`\`
services/<adapter-name>/
├── cmd/server/main.go
├── internal/
│   ├── adapter/
│   │   ├── adapter.go          # Implements DataProductAdapter interface
│   │   ├── client.go            # Source system client
│   │   └── transformer.go       # Data transformation logic
│   ├── events/
│   │   └── publisher.go         # Kafka event publisher
│   └── cache/
│       └── cache.go             # Caching layer
├── Dockerfile
├── go.mod
└── README.md
\`\`\`

## Step 2: Implement DataProductAdapter Interface

\`\`\`go
// internal/adapter/adapter.go
package adapter

import (
    "context"
    "github.com/dictamesh/common/pkg/models"
)

type SupplierAdapter struct {
    client         *SupplierAPIClient
    eventPublisher *kafka.Producer
    cache          *redis.Client
}

func (a *SupplierAdapter) GetEntity(ctx context.Context, id string) (*models.Entity, error) {
    // 1. Check cache
    // 2. Fetch from source
    // 3. Transform to canonical model
    // 4. Cache result
    // 5. Return entity
}

func (a *SupplierAdapter) StreamChanges(ctx context.Context) (<-chan models.ChangeEvent, error) {
    // Implement change detection (webhook or polling)
}
\`\`\`

## Step 3: Register Schema

\`\`\`bash
# Create Avro schema
cat > schemas/supplier-change-event.avsc << EOF
{
  "type": "record",
  "name": "SupplierChangeEvent",
  "namespace": "com.dictamesh.events",
  "fields": [
    {"name": "event_id", "type": "string"},
    {"name": "supplier_id", "type": "string"},
    {"name": "event_type", "type": {"type": "enum", "symbols": ["CREATED", "UPDATED", "DELETED"]}}
  ]
}
EOF

# Register with Schema Registry
curl -X POST http://schema-registry:8081/subjects/suppliers.api.entity_changed-value/versions \\
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \\
  --data @schemas/supplier-change-event.avsc
\`\`\`

## Step 4: Create Kubernetes Manifests
[... k8s setup ...]

## Step 5: Add to GraphQL Federation
[... federation setup ...]

## Step 6: Write Tests
[... testing guide ...]
```

### 4.4 API Reference

**File:** `docs/development/api-reference/metadata-catalog-api.md`

**Auto-generated using tools:**
- OpenAPI/Swagger for REST APIs
- GraphQL schema introspection for GraphQL APIs

### 4.5 Testing Guide

**File:** `docs/development/testing-guide.md`

**Content:**
- Unit testing guidelines
- Integration testing setup
- End-to-end testing
- Load testing procedures
- Mocking strategies
- Test data management

---

## 🤝 5. Contribution Documentation

### 5.1 Contributing Guide

**File:** `docs/contribution/CONTRIBUTING.md`

```markdown
# Contributing to DictaMesh

Thank you for your interest in contributing to DictaMesh!

## Code of Conduct

Please read our [Code of Conduct](code-of-conduct.md) before contributing.

## How to Contribute

### Reporting Bugs

Use our [bug report template](.github/ISSUE_TEMPLATE/bug_report.md):

1. Clear description of the issue
2. Steps to reproduce
3. Expected vs actual behavior
4. Environment details
5. Logs/screenshots

### Suggesting Features

Use our [feature request template](.github/ISSUE_TEMPLATE/feature_request.md):

1. Problem statement
2. Proposed solution
3. Alternatives considered
4. Impact assessment

### Pull Request Process

1. **Fork and clone**
   \`\`\`bash
   git clone https://github.com/YOUR_USERNAME/dictamesh.git
   cd dictamesh
   git remote add upstream https://github.com/controle-digital/dictamesh.git
   \`\`\`

2. **Create feature branch**
   \`\`\`bash
   git checkout -b feature/your-feature-name
   \`\`\`

3. **Make changes**
   - Follow [code style guide](code-style.md)
   - Add tests for new functionality
   - Update documentation

4. **Commit with conventional commits**
   \`\`\`bash
   git commit -m "feat(adapter): add supplier adapter support"
   \`\`\`

   Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

5. **Push and create PR**
   \`\`\`bash
   git push origin feature/your-feature-name
   \`\`\`

6. **PR Review Process**
   - Automated checks must pass (tests, linting, security scans)
   - At least 2 approvals required
   - Address review feedback
   - Squash commits before merge

## Development Workflow

### Before Starting
- Check existing issues/PRs
- Discuss major changes in an issue first
- Ensure your local environment is set up (see [local-setup.md](../development/local-setup.md))

### Code Quality Standards
- Test coverage: minimum 80%
- All tests must pass
- No linting errors
- Documentation updated

### Commit Message Format
\`\`\`
<type>(<scope>): <subject>

<body>

<footer>
\`\`\`

Example:
\`\`\`
feat(customer-adapter): add webhook support for Directus events

Implements real-time event detection via Directus webhooks instead of polling.
This improves event freshness from ~5s to <1s.

Closes #123
\`\`\`

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
```

### 5.2 Code of Conduct

**File:** `docs/contribution/code-of-conduct.md`

(Standard Contributor Covenant)

### 5.3 Pull Request Template

**File:** `.github/pull_request_template.md`

```markdown
## Description
<!-- Describe your changes in detail -->

## Related Issue
<!-- Link to related issue: Closes #123 -->

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] All tests passing

## Checklist
- [ ] My code follows the code style of this project
- [ ] I have updated the documentation accordingly
- [ ] I have added tests to cover my changes
- [ ] All new and existing tests passed
- [ ] My commits follow conventional commit format
- [ ] I have run linting and fixed all issues

## Screenshots (if applicable)

## Additional Notes
```

---

## 🤖 LLM Agent Documentation Generation Strategy

### Automated Documentation

```yaml
# .github/workflows/generate-docs.yaml
name: Generate Documentation

on:
  push:
    branches: [main, develop]
    paths:
      - 'services/**/*.go'
      - 'infrastructure/**/*.yaml'

jobs:
  generate-api-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Generate OpenAPI specs
        run: |
          go install github.com/swaggo/swag/cmd/swag@latest
          cd services/metadata-catalog
          swag init -g cmd/server/main.go -o docs/api

      - name: Generate GraphQL schema docs
        run: |
          npm install -g graphql-markdown
          graphql-markdown http://localhost:8080/graphql > docs/api/graphql-schema.md

      - name: Commit generated docs
        run: |
          git config user.name "Documentation Bot"
          git config user.email "docs@dictamesh.io"
          git add docs/api/
          git commit -m "docs: update API documentation [skip ci]" || true
          git push
```

### Documentation Validation

```bash
#!/bin/bash
# scripts/validate-docs.sh

echo "Validating documentation..."

# Check for broken links
docker run --rm -v $(pwd):/docs ghcr.io/tcort/markdown-link-check:stable /docs/**/*.md

# Check code examples
for file in docs/**/*.md; do
    # Extract code blocks
    awk '/```bash/,/```/' "$file" > /tmp/code-block.sh
    # Validate bash syntax
    bash -n /tmp/code-block.sh || echo "Syntax error in $file"
done

echo "Documentation validation complete"
```

---

## 📅 Documentation Maintenance Schedule

### Weekly
- [ ] Update CHANGELOG.md with merged PRs
- [ ] Review and close stale issues
- [ ] Update metrics/dashboards documentation

### Monthly
- [ ] Review and update troubleshooting guides based on support tickets
- [ ] Validate all code examples still work
- [ ] Update performance benchmarks

### Quarterly
- [ ] Full documentation audit
- [ ] Update architecture diagrams
- [ ] Review and update all runbooks
- [ ] User feedback review and documentation improvements

---

## 🎯 Documentation Completeness Checklist

### For Each Component

- [ ] README.md with overview
- [ ] Architecture documentation
- [ ] API reference (auto-generated)
- [ ] Configuration reference
- [ ] Deployment guide
- [ ] Troubleshooting section
- [ ] Examples
- [ ] Testing guide

### Cross-Cutting Documentation

- [ ] Overall system architecture
- [ ] Getting started guide
- [ ] API quickstart
- [ ] Administration guide
- [ ] Security documentation
- [ ] Performance tuning guide
- [ ] Contribution guidelines
- [ ] Runbooks for common issues

---

[← Previous: Testing Strategy](13-TESTING-STRATEGY.md) | [Next: Security & Compliance →](15-SECURITY-COMPLIANCE.md)

---

**Document Metadata**
- Version: 1.0.0
- Last Updated: 2025-11-08
- Documentation Standard: Markdown with embedded code examples
- Auto-generation: OpenAPI, GraphQL introspection
