# CI/CD Pipeline with ArgoCD

[← Previous: Deployment Strategy](04-DEPLOYMENT-STRATEGY.md) | [Next: Layer 1 Adapters →](06-LAYER1-ADAPTERS.md)

---

## 🎯 Purpose

This document provides LLM agents with comprehensive CI/CD pipeline configuration using ArgoCD for GitOps-based continuous delivery on K3S.

**Reading Time:** 30 minutes
**Prerequisites:** [Infrastructure Planning](03-INFRASTRUCTURE-PLANNING.md), [Deployment Strategy](04-DEPLOYMENT-STRATEGY.md)
**Outputs:** ArgoCD configuration, GitHub Actions workflows, deployment automation

---

## 🏗️ GitOps Architecture

### GitOps Principles

```
┌─────────────────┐
│  Git Repository │  ← Single Source of Truth
│  (main/develop) │
└────────┬────────┘
         │
         │ ArgoCD watches
         │
         ▼
┌─────────────────┐
│     ArgoCD      │  ← Reconciliation Engine
│   (Controller)  │
└────────┬────────┘
         │
         │ kubectl apply
         │
         ▼
┌─────────────────┐
│  K3S Cluster    │  ← Desired State
│  (Actual State) │
└─────────────────┘
```

### Repository Structure

```
dictamesh/
├── services/                    # Application source code
│   ├── customer-adapter/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── product-adapter/
│   ├── invoice-adapter/
│   ├── metadata-catalog/
│   └── graphql-gateway/
│
├── infrastructure/
│   ├── k8s/                     # Kubernetes manifests (ArgoCD source)
│   │   ├── base/                # Kustomize base
│   │   │   ├── customer-adapter/
│   │   │   │   ├── deployment.yaml
│   │   │   │   ├── service.yaml
│   │   │   │   └── kustomization.yaml
│   │   │   ├── metadata-catalog/
│   │   │   └── graphql-gateway/
│   │   │
│   │   └── overlays/            # Environment-specific
│   │       ├── dev/
│   │       │   ├── kustomization.yaml
│   │       │   └── patches/
│   │       ├── staging/
│   │       └── prod/
│   │
│   ├── argocd/
│   │   ├── applications/        # ArgoCD Application manifests
│   │   │   ├── dev/
│   │   │   │   ├── customer-adapter.yaml
│   │   │   │   ├── metadata-catalog.yaml
│   │   │   │   └── app-of-apps.yaml
│   │   │   ├── staging/
│   │   │   └── prod/
│   │   │
│   │   └── projects/
│   │       └── dictamesh.yaml
│   │
│   └── helm/                    # Helm charts for complex components
│       ├── kafka/
│       ├── postgresql/
│       └── monitoring/
│
├── .github/
│   └── workflows/               # CI pipelines
│       ├── customer-adapter-ci.yaml
│       ├── metadata-catalog-ci.yaml
│       └── image-scan.yaml
│
└── docs/
    └── planning/
```

---

## 🔧 ArgoCD Installation & Configuration

### ArgoCD Installation

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Install ArgoCD CLI (Linux)
curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
chmod +x /usr/local/bin/argocd

# Wait for ArgoCD to be ready
kubectl wait --for=condition=available --timeout=600s \
  deployment/argocd-server -n argocd

# Get initial admin password
ARGOCD_PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d)
echo "ArgoCD admin password: $ARGOCD_PASSWORD"

# Expose ArgoCD server
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'

# Or use port-forward for testing
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

### ArgoCD Ingress with TLS

```yaml
# infrastructure/k8s/argocd/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: argocd-server
  namespace: argocd
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-passthrough: "true"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - argocd.dictamesh.com
      secretName: argocd-server-tls
  rules:
    - host: argocd.dictamesh.com
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

### ArgoCD Project Configuration

```yaml
# infrastructure/argocd/projects/dictamesh.yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: dictamesh
  namespace: argocd
spec:
  description: DictaMesh Data Mesh Platform

  # Allowed source repositories
  sourceRepos:
    - https://github.com/controle-digital/dictamesh.git
    - https://charts.bitnami.com/bitnami
    - https://strimzi.io/charts/
    - https://prometheus-community.github.io/helm-charts

  # Allowed destination clusters and namespaces
  destinations:
    - namespace: 'dictamesh-*'
      server: https://kubernetes.default.svc
    - namespace: 'argocd'
      server: https://kubernetes.default.svc

  # Cluster resource whitelist
  clusterResourceWhitelist:
    - group: ''
      kind: Namespace
    - group: 'rbac.authorization.k8s.io'
      kind: ClusterRole
    - group: 'rbac.authorization.k8s.io'
      kind: ClusterRoleBinding

  # Namespace resource whitelist
  namespaceResourceWhitelist:
    - group: '*'
      kind: '*'

  # Deny certain resources
  namespaceResourceBlacklist:
    - group: ''
      kind: ResourceQuota
    - group: ''
      kind: LimitRange

  # Sync windows
  syncWindows:
    - kind: allow
      schedule: '0 9-17 * * 1-5'  # Mon-Fri, 9am-5pm
      duration: 8h
      applications:
        - '*-prod'
      manualSync: true
    - kind: allow
      schedule: '* * * * *'  # Always
      duration: 24h
      applications:
        - '*-dev'
        - '*-staging'

  # Orphaned resources monitoring
  orphanedResources:
    warn: true
```

---

## 📦 Application Definitions

### App of Apps Pattern

```yaml
# infrastructure/argocd/applications/dev/app-of-apps.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: dictamesh-dev
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: dictamesh

  source:
    repoURL: https://github.com/dictamesh/dictamesh.git
    targetRevision: develop
    path: infrastructure/argocd/applications/dev

  destination:
    server: https://kubernetes.default.svc
    namespace: argocd

  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
      - CreateNamespace=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

### Individual Application: Customer Adapter

```yaml
# infrastructure/argocd/applications/dev/customer-adapter.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: customer-adapter-dev
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: dictamesh

  source:
    repoURL: https://github.com/dictamesh/dictamesh.git
    targetRevision: develop
    path: infrastructure/k8s/overlays/dev/customer-adapter

  destination:
    server: https://kubernetes.default.svc
    namespace: dictamesh-dev

  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
      - ApplyOutOfSyncOnly=true
    retry:
      limit: 3
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 1m

  # Health assessment
  ignoreDifferences:
    - group: apps
      kind: Deployment
      jsonPointers:
        - /spec/replicas  # Ignore HPA changes

  # Hooks for pre/post sync
  syncPolicy:
    syncOptions:
      - PruneLast=true
```

### Helm-Based Application: Kafka

```yaml
# infrastructure/argocd/applications/dev/kafka.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: kafka-dev
  namespace: argocd
spec:
  project: dictamesh

  source:
    chart: kafka
    repoURL: https://charts.bitnami.com/bitnami
    targetRevision: 26.4.3
    helm:
      releaseName: kafka
      values: |
        replicaCount: 1
        persistence:
          enabled: true
          storageClass: longhorn-retain
          size: 50Gi
        metrics:
          kafka:
            enabled: true
          jmx:
            enabled: true

  destination:
    server: https://kubernetes.default.svc
    namespace: dictamesh-infra

  syncPolicy:
    automated:
      prune: false  # Don't auto-delete Kafka (data safety)
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

---

## 🔄 Kustomize Structure

### Base Deployment

```yaml
# infrastructure/k8s/base/customer-adapter/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: customer-adapter
  labels:
    app: customer-adapter
    component: adapter
    domain: customers
spec:
  replicas: 1  # Overridden by overlay
  selector:
    matchLabels:
      app: customer-adapter
  template:
    metadata:
      labels:
        app: customer-adapter
        version: latest  # Overridden by overlay
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: dictamesh-app
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - name: customer-adapter
          image: ghcr.io/controle-digital/customer-adapter:latest
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
            - name: metrics
              containerPort: 8080
              protocol: TCP
          env:
            - name: ENVIRONMENT
              value: dev
            - name: LOG_LEVEL
              value: info
            - name: DIRECTUS_URL
              valueFrom:
                configMapKeyRef:
                  name: customer-adapter-config
                  key: directus_url
            - name: KAFKA_BROKERS
              valueFrom:
                configMapKeyRef:
                  name: kafka-config
                  key: brokers
            - name: REDIS_URL
              valueFrom:
                secretKeyRef:
                  name: redis-credentials
                  key: url
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          livenessProbe:
            httpGet:
              path: /health/live
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /health/ready
              port: http
            initialDelaySeconds: 10
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
          volumeMounts:
            - name: cache
              mountPath: /tmp/cache
      volumes:
        - name: cache
          emptyDir: {}
```

```yaml
# infrastructure/k8s/base/customer-adapter/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: customer-adapter
  labels:
    app: customer-adapter
spec:
  type: ClusterIP
  ports:
    - port: 8080
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app: customer-adapter
```

```yaml
# infrastructure/k8s/base/customer-adapter/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - deployment.yaml
  - service.yaml
  - configmap.yaml
  - serviceaccount.yaml

commonLabels:
  app.kubernetes.io/name: customer-adapter
  app.kubernetes.io/component: adapter
  app.kubernetes.io/part-of: dictamesh

images:
  - name: ghcr.io/dictamesh/customer-adapter
    newTag: latest
```

### Development Overlay

```yaml
# infrastructure/k8s/overlays/dev/customer-adapter/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: dictamesh-dev

bases:
  - ../../../base/customer-adapter

commonLabels:
  environment: dev

images:
  - name: ghcr.io/dictamesh/customer-adapter
    newTag: dev-latest

patches:
  - path: replica-count.yaml
    target:
      kind: Deployment
      name: customer-adapter

configMapGenerator:
  - name: customer-adapter-config
    behavior: merge
    literals:
      - LOG_LEVEL=debug
      - DIRECTUS_URL=https://directus-dev.dictamesh.com

replicas:
  - name: customer-adapter
    count: 1
```

```yaml
# infrastructure/k8s/overlays/dev/customer-adapter/replica-count.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: customer-adapter
spec:
  replicas: 1
```

### Production Overlay

```yaml
# infrastructure/k8s/overlays/prod/customer-adapter/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: dictamesh-prod

bases:
  - ../../../base/customer-adapter

commonLabels:
  environment: prod

images:
  - name: ghcr.io/dictamesh/customer-adapter
    newTag: v1.0.0  # Specific version tag

patches:
  - path: resources.yaml
  - path: hpa.yaml

configMapGenerator:
  - name: customer-adapter-config
    behavior: merge
    literals:
      - LOG_LEVEL=info
      - DIRECTUS_URL=https://directus.dictamesh.com

replicas:
  - name: customer-adapter
    count: 3
```

```yaml
# infrastructure/k8s/overlays/prod/customer-adapter/resources.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: customer-adapter
spec:
  template:
    spec:
      containers:
        - name: customer-adapter
          resources:
            requests:
              cpu: 500m
              memory: 512Mi
            limits:
              cpu: 2000m
              memory: 2Gi
```

```yaml
# infrastructure/k8s/overlays/prod/customer-adapter/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: customer-adapter-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: customer-adapter
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

---

## 🚀 CI Pipeline (GitHub Actions)

### Docker Build & Push

```yaml
# .github/workflows/customer-adapter-ci.yaml
name: Customer Adapter CI

on:
  push:
    branches: [develop, staging, main]
    paths:
      - 'services/customer-adapter/**'
      - '.github/workflows/customer-adapter-ci.yaml'
  pull_request:
    branches: [develop, staging, main]
    paths:
      - 'services/customer-adapter/**'

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: dictamesh/customer-adapter

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          cache-dependency-path: services/customer-adapter/go.sum

      - name: Run tests
        working-directory: services/customer-adapter
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./services/customer-adapter/coverage.out
          flags: customer-adapter

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          working-directory: services/customer-adapter

  build-and-push:
    needs: [test, lint]
    runs-on: ubuntu-latest
    if: github.event_name != 'pull_request'
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=sha,prefix={{branch}}-
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: services/customer-adapter
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
          format: 'sarif'
          output: 'trivy-results.sarif'

      - name: Upload Trivy results to GitHub Security
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'

  update-manifest:
    needs: build-and-push
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    steps:
      - uses: actions/checkout@v4

      - name: Update Kustomize image tag
        working-directory: infrastructure/k8s/overlays/dev/customer-adapter
        run: |
          kustomize edit set image \
            ghcr.io/dictamesh/customer-adapter=ghcr.io/dictamesh/customer-adapter:develop-${{ github.sha }}

      - name: Commit and push changes
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add infrastructure/k8s/overlays/dev/customer-adapter/kustomization.yaml
          git commit -m "Update customer-adapter dev image to develop-${{ github.sha }}"
          git push
```

### ArgoCD Sync Trigger

```yaml
# .github/workflows/argocd-sync.yaml
name: ArgoCD Sync

on:
  workflow_run:
    workflows: ["Customer Adapter CI", "Metadata Catalog CI", "GraphQL Gateway CI"]
    types:
      - completed
    branches: [develop]

jobs:
  argocd-sync:
    runs-on: ubuntu-latest
    if: ${{ github.event.workflow_run.conclusion == 'success' }}
    steps:
      - name: Install ArgoCD CLI
        run: |
          curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
          chmod +x /usr/local/bin/argocd

      - name: ArgoCD Login
        run: |
          argocd login argocd.dictamesh.com \
            --username admin \
            --password ${{ secrets.ARGOCD_PASSWORD }} \
            --grpc-web

      - name: Sync Application
        run: |
          argocd app sync dictamesh-dev --prune
          argocd app wait dictamesh-dev --health --timeout 600
```

---

## 📊 Monitoring ArgoCD

### ServiceMonitor for Prometheus

```yaml
# infrastructure/k8s/monitoring/argocd-servicemonitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: argocd-metrics
  namespace: argocd
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: argocd-metrics
  endpoints:
    - port: metrics
      interval: 30s
```

### Grafana Dashboard

```yaml
# infrastructure/k8s/monitoring/argocd-dashboard-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-dashboard
  namespace: dictamesh-monitoring
  labels:
    grafana_dashboard: "1"
data:
  argocd-dashboard.json: |
    {
      "dashboard": {
        "title": "ArgoCD",
        "panels": [
          {
            "title": "Application Health",
            "targets": [
              {
                "expr": "argocd_app_info{health_status=\"Healthy\"}"
              }
            ]
          }
        ]
      }
    }
```

---

## 🎯 LLM Agent CI/CD Checklist

### Setup Phase

- [ ] ArgoCD installed and accessible
- [ ] ArgoCD project created
- [ ] Repository connected to ArgoCD
- [ ] Kustomize base manifests created
- [ ] Environment overlays configured (dev/staging/prod)
- [ ] App-of-Apps configured
- [ ] GitHub Actions workflows created
- [ ] Container registry authentication configured
- [ ] Image scanning enabled (Trivy)

### Deployment Phase

- [ ] Base application deployed via ArgoCD
- [ ] Health checks passing
- [ ] Metrics being collected
- [ ] Logs aggregated
- [ ] Alerts configured
- [ ] Rollback tested

### Validation

```bash
# Check ArgoCD application status
argocd app list
argocd app get customer-adapter-dev

# Verify deployment
kubectl get deployments -n dictamesh-dev
kubectl get pods -n dictamesh-dev

# Check sync status
argocd app sync customer-adapter-dev --dry-run
```

---

[← Previous: Deployment Strategy](04-DEPLOYMENT-STRATEGY.md) | [Next: Layer 1 Adapters →](06-LAYER1-ADAPTERS.md)

---

**Document Metadata**
- Version: 1.0.0
- Last Updated: 2025-11-08
- GitOps Tool: ArgoCD
- CI Platform: GitHub Actions
