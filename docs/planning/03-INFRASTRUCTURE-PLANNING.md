# Infrastructure Planning - K3S Cluster on Controle Digital

[← Previous: Implementation Phases](02-IMPLEMENTATION-PHASES.md) | [Next: Deployment Strategy →](04-DEPLOYMENT-STRATEGY.md)

---

## 🎯 Purpose

This document provides LLM agents with comprehensive infrastructure planning for deploying DictaMesh on the existing K3S cluster at Controle Digital Ltda.

**Reading Time:** 25 minutes
**Prerequisites:** [Architecture Overview](01-ARCHITECTURE-OVERVIEW.md), [Implementation Phases](02-IMPLEMENTATION-PHASES.md)
**Outputs:** Infrastructure manifests, resource specifications, storage planning

---

## 🏗️ K3S Cluster Overview

### Assumptions About Controle Digital K3S Cluster

```yaml
Cluster Type: K3S (Lightweight Kubernetes)
Container Runtime: containerd
Default Ingress: Traefik (K3S bundled)
Default Storage: local-path provisioner
Service LoadBalancer: Klipper (K3S bundled) or MetalLB
Kubernetes Version: v1.28+ (assumed)
```

### Cluster Discovery Commands for LLM Agents

```bash
# Verify cluster information
kubectl cluster-info
kubectl version --short
kubectl get nodes -o wide

# Check default storage class
kubectl get storageclass
kubectl describe storageclass local-path

# Check ingress controller
kubectl get pods -n kube-system | grep traefik
kubectl get svc -n kube-system | grep traefik

# Check available resources
kubectl top nodes
kubectl describe nodes
```

### Expected Output Analysis

```bash
# LLM Agent should analyze:
# 1. Available CPU/Memory per node
# 2. Storage class capabilities
# 3. Ingress controller type
# 4. Load balancer availability
# 5. Node labels and taints
```

---

## 📦 Storage Strategy

### Storage Classes Required

```yaml
# infrastructure/k8s/storage/storage-classes.yaml

# Option 1: If Longhorn is available (preferred for production)
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: longhorn-retain
provisioner: driver.longhorn.io
allowVolumeExpansion: true
reclaimPolicy: Retain
volumeBindingMode: Immediate
parameters:
  numberOfReplicas: "3"
  staleReplicaTimeout: "2880"
  fromBackup: ""
  fsType: "ext4"
---
# Option 2: If NFS is available
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: nfs-client
provisioner: k8s-sigs.io/nfs-subdir-external-provisioner
parameters:
  archiveOnDelete: "true"
reclaimPolicy: Retain
volumeBindingMode: Immediate
---
# Option 3: Local path (K3S default) - for dev/staging only
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-path-retain
provisioner: rancher.io/local-path
reclaimPolicy: Retain
volumeBindingMode: WaitForFirstConsumer
```

### Storage Assessment Script for LLM Agents

```bash
#!/bin/bash
# infrastructure/scripts/assess-storage.sh

echo "=== Storage Class Assessment ==="
kubectl get storageclass

echo -e "\n=== Longhorn Check ==="
if kubectl get ns longhorn-system &>/dev/null; then
    echo "✓ Longhorn detected"
    kubectl get pods -n longhorn-system
else
    echo "✗ Longhorn not installed"
fi

echo -e "\n=== NFS Provisioner Check ==="
if kubectl get pods --all-namespaces | grep nfs-subdir; then
    echo "✓ NFS provisioner detected"
else
    echo "✗ NFS provisioner not found"
fi

echo -e "\n=== Available PVs ==="
kubectl get pv

echo -e "\n=== Node Storage Capacity ==="
kubectl get nodes -o custom-columns=NAME:.metadata.name,STORAGE:.status.capacity.ephemeral-storage
```

### Storage Installation (If Not Available)

#### Option A: Install Longhorn (Recommended)

```bash
# Add Longhorn Helm repository
helm repo add longhorn https://charts.longhorn.io
helm repo update

# Install Longhorn
helm install longhorn longhorn/longhorn \
  --namespace longhorn-system \
  --create-namespace \
  --set defaultSettings.defaultDataPath="/var/lib/longhorn" \
  --set defaultSettings.defaultReplicaCount=3

# Wait for deployment
kubectl -n longhorn-system wait --for=condition=available --timeout=600s deployment/longhorn-driver-deployer

# Access Longhorn UI
kubectl port-forward -n longhorn-system svc/longhorn-frontend 8000:80
```

#### Option B: Install NFS Client Provisioner

```bash
# Assuming NFS server is available at 192.168.1.100:/exports/k3s-storage
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
helm install nfs-subdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
    --set nfs.server=192.168.1.100 \
    --set nfs.path=/exports/k3s-storage \
    --set storageClass.name=nfs-client \
    --set storageClass.defaultClass=false
```

### Storage Requirements by Component

| Component | Storage Type | Size (Dev) | Size (Staging) | Size (Prod) | Retention |
|-----------|--------------|------------|----------------|-------------|-----------|
| **PostgreSQL (Metadata Catalog)** | Block (RWO) | 20Gi | 50Gi | 100Gi | Permanent |
| **Kafka Brokers (×3)** | Block (RWO) | 50Gi each | 100Gi each | 200Gi each | 30-90 days |
| **Zookeeper (×3)** | Block (RWO) | 10Gi each | 20Gi each | 30Gi each | Permanent |
| **Redis** | Block (RWO) | 5Gi | 10Gi | 20Gi | Volatile |
| **Prometheus** | Block (RWO) | 30Gi | 50Gi | 100Gi | 30 days |
| **Loki** | Block (RWO) | 30Gi | 50Gi | 100Gi | 14 days |
| **Grafana** | Block (RWO) | 5Gi | 5Gi | 10Gi | Permanent |
| **ArgoCD** | Block (RWO) | 5Gi | 5Gi | 10Gi | Permanent |
| **Backups** | Object/NFS | - | 100Gi | 500Gi | 90 days |

**Total Storage Required:**
- **Development:** ~200Gi
- **Staging:** ~420Gi
- **Production:** ~900Gi + backups

---

## 🌐 Networking & Ingress

### Ingress Strategy

#### Option 1: Use K3S Traefik (Default)

```yaml
# infrastructure/k8s/ingress/traefik-config.yaml
apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: traefik
  namespace: kube-system
spec:
  valuesContent: |-
    additionalArguments:
      - "--entrypoints.websecure.http.tls=true"
      - "--providers.kubernetescrd"
      - "--api.dashboard=true"
    ports:
      web:
        port: 80
        redirectTo: websecure
      websecure:
        port: 443
        tls:
          enabled: true
```

#### Option 2: Install NGINX Ingress (Alternative)

```bash
# Install NGINX Ingress Controller
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.type=LoadBalancer \
  --set controller.metrics.enabled=true \
  --set controller.podAnnotations."prometheus\.io/scrape"=true
```

### DNS & TLS Configuration

```yaml
# infrastructure/k8s/cert-manager/cert-manager-install.yaml

# Install cert-manager for automatic TLS certificates
apiVersion: v1
kind: Namespace
metadata:
  name: cert-manager
---
# Install via Helm
# helm install cert-manager jetstack/cert-manager \
#   --namespace cert-manager \
#   --set installCRDs=true
```

```yaml
# infrastructure/k8s/cert-manager/letsencrypt-issuer.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: devops@dictamesh.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - http01:
          ingress:
            class: traefik  # or nginx
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    email: devops@dictamesh.com
    privateKeySecretRef:
      name: letsencrypt-staging
    solvers:
      - http01:
          ingress:
            class: traefik
```

### Ingress Definitions

```yaml
# infrastructure/k8s/ingress/dictamesh-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: dictamesh-gateway
  namespace: dictamesh-prod
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    traefik.ingress.kubernetes.io/router.entrypoints: websecure
    traefik.ingress.kubernetes.io/router.tls: "true"
spec:
  ingressClassName: traefik
  tls:
    - hosts:
        - api.dictamesh.com
      secretName: dictamesh-api-tls
  rules:
    - host: api.dictamesh.com
      http:
        paths:
          - path: /graphql
            pathType: Prefix
            backend:
              service:
                name: graphql-gateway
                port:
                  number: 8080
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: dictamesh-monitoring
  namespace: dictamesh-monitoring
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    traefik.ingress.kubernetes.io/router.middlewares: dictamesh-monitoring-basic-auth@kubernetescrd
spec:
  ingressClassName: traefik
  tls:
    - hosts:
        - grafana.dictamesh.com
      secretName: grafana-tls
  rules:
    - host: grafana.dictamesh.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: prometheus-grafana
                port:
                  number: 80
```

### Network Policies

```yaml
# infrastructure/k8s/network-policies/default-deny.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: dictamesh-prod
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
---
# infrastructure/k8s/network-policies/allow-graphql-gateway.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-graphql-gateway
  namespace: dictamesh-prod
spec:
  podSelector:
    matchLabels:
      app: graphql-gateway
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx  # or kube-system for Traefik
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: customer-adapter
      ports:
        - protocol: TCP
          port: 8080
    - to:
        - namespaceSelector:
            matchLabels:
              name: dictamesh-infra
        - podSelector:
            matchLabels:
              app: metadata-catalog
      ports:
        - protocol: TCP
          port: 8080
    - to:  # DNS egress
        - namespaceSelector:
            matchLabels:
              name: kube-system
        - podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - protocol: UDP
          port: 53
```

---

## 💾 Resource Allocation

### Compute Resources by Component

#### Development Environment

```yaml
# infrastructure/k8s/dev/resource-quotas.yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: dictamesh-dev-quota
  namespace: dictamesh-dev
spec:
  hard:
    requests.cpu: "10"
    requests.memory: "20Gi"
    limits.cpu: "20"
    limits.memory: "40Gi"
    persistentvolumeclaims: "20"
    services.loadbalancers: "2"
```

#### Component Resource Specifications

```yaml
# Customer Adapter (Development)
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi

# Customer Adapter (Production)
resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 2000m
    memory: 2Gi
```

### Resource Matrix

| Component | Replicas (Dev) | CPU Request | Memory Request | CPU Limit | Memory Limit | Replicas (Prod) |
|-----------|----------------|-------------|----------------|-----------|--------------|-----------------|
| **Customer Adapter** | 1 | 100m | 128Mi | 500m | 512Mi | 3 |
| **Product Adapter** | 1 | 100m | 128Mi | 500m | 512Mi | 3 |
| **Invoice Adapter** | 1 | 100m | 128Mi | 500m | 512Mi | 3 |
| **Metadata Catalog** | 1 | 200m | 256Mi | 1000m | 1Gi | 3 |
| **GraphQL Gateway** | 1 | 200m | 256Mi | 1000m | 1Gi | 3 |
| **PostgreSQL** | 1 | 500m | 512Mi | 2000m | 4Gi | 3 |
| **Kafka Broker** | 1 | 500m | 1Gi | 2000m | 4Gi | 3 |
| **Zookeeper** | 1 | 100m | 256Mi | 500m | 1Gi | 3 |
| **Redis** | 1 | 100m | 128Mi | 500m | 512Mi | 3 |
| **Schema Registry** | 1 | 100m | 256Mi | 500m | 512Mi | 2 |
| **Prometheus** | 1 | 500m | 1Gi | 2000m | 4Gi | 2 |
| **Grafana** | 1 | 100m | 128Mi | 500m | 512Mi | 2 |
| **Loki** | 1 | 200m | 256Mi | 1000m | 1Gi | 2 |
| **Jaeger** | 1 | 200m | 256Mi | 1000m | 1Gi | 2 |

**Total Resources Required:**

**Development:**
- CPU Requests: ~3.5 cores
- CPU Limits: ~12 cores
- Memory Requests: ~5Gi
- Memory Limits: ~15Gi

**Production:**
- CPU Requests: ~20 cores
- CPU Limits: ~60 cores
- Memory Requests: ~40Gi
- Memory Limits: ~100Gi

### Horizontal Pod Autoscaling (HPA)

```yaml
# infrastructure/k8s/hpa/customer-adapter-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: customer-adapter-hpa
  namespace: dictamesh-prod
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
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
        - type: Percent
          value: 100
          periodSeconds: 30
```

---

## 🔒 Security Infrastructure

### Pod Security Standards

```yaml
# infrastructure/k8s/security/pod-security-standards.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: dictamesh-prod
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### Service Accounts & RBAC

```yaml
# infrastructure/k8s/security/rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dictamesh-app
  namespace: dictamesh-prod
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: dictamesh-app-role
  namespace: dictamesh-prod
rules:
  - apiGroups: [""]
    resources: ["configmaps", "secrets"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dictamesh-app-rolebinding
  namespace: dictamesh-prod
subjects:
  - kind: ServiceAccount
    name: dictamesh-app
    namespace: dictamesh-prod
roleRef:
  kind: Role
  name: dictamesh-app-role
  apiGroup: rbac.authorization.k8s.io
```

### Secrets Management

#### Option 1: Kubernetes Secrets (Baseline)

```yaml
# infrastructure/k8s/secrets/postgres-credentials.yaml
apiVersion: v1
kind: Secret
metadata:
  name: postgres-credentials
  namespace: dictamesh-prod
type: Opaque
stringData:
  username: dictamesh_user
  password: <GENERATE_STRONG_PASSWORD>
  database: metadata_catalog
  connection_string: postgresql://dictamesh_user:<PASSWORD>@metadata-catalog-db-rw.dictamesh-infra.svc:5432/metadata_catalog
```

#### Option 2: Sealed Secrets (Recommended)

```bash
# Install Sealed Secrets controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# Create sealed secret
echo -n "supersecretpassword" | kubectl create secret generic postgres-password \
  --dry-run=client --from-file=password=/dev/stdin -o yaml | \
  kubeseal -o yaml > postgres-password-sealed.yaml
```

#### Option 3: External Secrets Operator (Advanced)

```yaml
# If using Vault, AWS Secrets Manager, or similar
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: vault-backend
  namespace: dictamesh-prod
spec:
  provider:
    vault:
      server: "https://vault.dictamesh.com"
      path: "secret"
      version: "v2"
      auth:
        kubernetes:
          mountPath: "kubernetes"
          role: "dictamesh-app"
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: postgres-credentials
  namespace: dictamesh-prod
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: postgres-credentials
    creationPolicy: Owner
  data:
    - secretKey: password
      remoteRef:
        key: dictamesh/postgres
        property: password
```

---

## 📊 Monitoring Infrastructure

### Metrics Server

```bash
# K3S includes metrics-server by default, verify:
kubectl get deployment metrics-server -n kube-system

# If not present, install:
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

### Prometheus Stack

```bash
# Install kube-prometheus-stack via Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace dictamesh-monitoring \
  --create-namespace \
  --values infrastructure/k8s/monitoring/prometheus-values.yaml
```

```yaml
# infrastructure/k8s/monitoring/prometheus-values.yaml
prometheus:
  prometheusSpec:
    retention: 30d
    storageSpec:
      volumeClaimTemplate:
        spec:
          storageClassName: longhorn-retain
          accessModes: ["ReadWriteOnce"]
          resources:
            requests:
              storage: 50Gi
    resources:
      requests:
        cpu: 500m
        memory: 2Gi
      limits:
        cpu: 2000m
        memory: 8Gi
    serviceMonitorSelectorNilUsesHelmValues: false
    podMonitorSelectorNilUsesHelmValues: false

grafana:
  adminPassword: <GENERATE_STRONG_PASSWORD>
  persistence:
    enabled: true
    storageClassName: longhorn-retain
    size: 10Gi
  datasources:
    datasources.yaml:
      apiVersion: 1
      datasources:
        - name: Prometheus
          type: prometheus
          url: http://prometheus-kube-prometheus-prometheus:9090
          isDefault: true
        - name: Loki
          type: loki
          url: http://loki:3100

alertmanager:
  alertmanagerSpec:
    storage:
      volumeClaimTemplate:
        spec:
          storageClassName: longhorn-retain
          accessModes: ["ReadWriteOnce"]
          resources:
            requests:
              storage: 5Gi
```

---

## 🎯 LLM Agent Infrastructure Checklist

### Pre-Deployment Validation

```bash
#!/bin/bash
# infrastructure/scripts/pre-deployment-check.sh

echo "=== DictaMesh Infrastructure Pre-Deployment Check ==="

# Check 1: Cluster access
echo -e "\n[1/10] Checking cluster access..."
if kubectl cluster-info &>/dev/null; then
    echo "✓ Cluster accessible"
else
    echo "✗ Cannot access cluster"
    exit 1
fi

# Check 2: Required nodes
echo -e "\n[2/10] Checking node readiness..."
READY_NODES=$(kubectl get nodes --no-headers | grep -c " Ready")
echo "✓ $READY_NODES nodes ready"

# Check 3: Storage classes
echo -e "\n[3/10] Checking storage classes..."
kubectl get storageclass

# Check 4: Available resources
echo -e "\n[4/10] Checking available resources..."
kubectl top nodes

# Check 5: Namespaces
echo -e "\n[5/10] Checking namespaces..."
for ns in dictamesh-dev dictamesh-staging dictamesh-prod dictamesh-infra dictamesh-monitoring; do
    if kubectl get namespace $ns &>/dev/null; then
        echo "✓ $ns exists"
    else
        echo "✗ $ns missing - will create"
    fi
done

# Check 6: Ingress controller
echo -e "\n[6/10] Checking ingress controller..."
if kubectl get pods -n kube-system | grep -q traefik; then
    echo "✓ Traefik ingress detected"
elif kubectl get pods -n ingress-nginx | grep -q nginx; then
    echo "✓ NGINX ingress detected"
else
    echo "⚠ No ingress controller detected"
fi

# Check 7: Metrics server
echo -e "\n[7/10] Checking metrics server..."
if kubectl get deployment metrics-server -n kube-system &>/dev/null; then
    echo "✓ Metrics server available"
else
    echo "⚠ Metrics server not found"
fi

# Check 8: cert-manager
echo -e "\n[8/10] Checking cert-manager..."
if kubectl get pods -n cert-manager &>/dev/null; then
    echo "✓ cert-manager installed"
else
    echo "✗ cert-manager not installed - required for TLS"
fi

# Check 9: ArgoCD
echo -e "\n[9/10] Checking ArgoCD..."
if kubectl get pods -n argocd &>/dev/null; then
    echo "✓ ArgoCD installed"
else
    echo "✗ ArgoCD not installed - required for GitOps"
fi

# Check 10: Load balancer
echo -e "\n[10/10] Checking load balancer..."
LB_SERVICE=$(kubectl get svc -A | grep LoadBalancer | head -n 1)
if [ -n "$LB_SERVICE" ]; then
    echo "✓ Load balancer available"
    echo "$LB_SERVICE"
else
    echo "⚠ No load balancer detected - may need MetalLB"
fi

echo -e "\n=== Pre-Deployment Check Complete ==="
```

### Deployment Steps

- [ ] Run pre-deployment check script
- [ ] Install missing infrastructure (Longhorn, cert-manager, ArgoCD)
- [ ] Create namespaces
- [ ] Configure storage classes
- [ ] Set up ingress and TLS
- [ ] Deploy monitoring stack
- [ ] Configure network policies
- [ ] Create service accounts and RBAC
- [ ] Set up secrets management
- [ ] Validate all infrastructure components healthy

---

[← Previous: Implementation Phases](02-IMPLEMENTATION-PHASES.md) | [Next: Deployment Strategy →](04-DEPLOYMENT-STRATEGY.md)

---

**Document Metadata**
- Version: 1.0.0
- Last Updated: 2025-11-08
- Target Environment: K3S Cluster @ Controle Digital Ltda
