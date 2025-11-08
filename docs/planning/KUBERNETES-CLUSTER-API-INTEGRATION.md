# Kubernetes Cluster API Integration Module

**Planning Document**
**Version:** 1.0
**Date:** 2025-11-08
**Status:** Planning

## Executive Summary

This document outlines the complete design and implementation plan for a Kubernetes Cluster API integration module within the DictaMesh framework. This module will enable kubernetes management applications to use DictaMesh as an abstraction layer for managing Kubernetes clusters, providing:

- Unified API for multi-cluster management
- Event-driven cluster state synchronization
- Metadata catalog for cluster resources
- GraphQL API for cluster operations
- Advanced observability and governance for Kubernetes resources

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Components](#components)
4. [Resource Types](#resource-types)
5. [API Design](#api-design)
6. [Event Streaming](#event-streaming)
7. [Security & Governance](#security--governance)
8. [Implementation Plan](#implementation-plan)
9. [Deployment Strategy](#deployment-strategy)
10. [Testing Strategy](#testing-strategy)

---

## 1. Overview

### 1.1 Purpose

The Kubernetes Cluster API integration module provides a DictaMesh adapter that:

- **Connects** to one or more Kubernetes clusters via the Kubernetes API
- **Exposes** Kubernetes resources (Pods, Deployments, Services, etc.) as DictaMesh entities
- **Streams** cluster events in real-time to the DictaMesh event bus
- **Provides** a unified GraphQL API for kubernetes management applications
- **Enables** multi-cluster management through a single interface
- **Tracks** resource relationships and dependencies in the metadata catalog
- **Enforces** RBAC and governance policies on cluster operations

### 1.2 Use Cases

#### UC1: Multi-Cluster Management Platform
A platform team manages 50+ Kubernetes clusters across dev/staging/prod environments. They use the DictaMesh Kubernetes adapter to:
- Query all pods across all clusters via unified GraphQL API
- Track resource dependencies (e.g., which services depend on which deployments)
- Receive real-time events when resources change in any cluster
- Enforce governance policies (e.g., namespace quotas, security policies)

#### UC2: Kubernetes Monitoring Dashboard
A monitoring application uses the adapter to:
- Display real-time cluster status across multiple clusters
- Aggregate metrics and events from all clusters
- Correlate Kubernetes events with application logs and traces
- Alert on resource state changes

#### UC3: GitOps Automation
A GitOps workflow engine uses the adapter to:
- Query current cluster state
- Compare desired state (from Git) with actual state (from cluster)
- Apply changes via the adapter's mutation API
- Track deployment lineage and audit trail

#### UC4: Cost Optimization Service
A cost optimization service uses the adapter to:
- Query resource utilization across clusters
- Identify underutilized resources
- Recommend rightsizing based on actual usage
- Track resource costs and allocation

### 1.3 Key Features

- ✅ **Multi-Cluster Support** - Connect to multiple Kubernetes clusters simultaneously
- ✅ **Real-Time Event Streaming** - Stream cluster events via Kubernetes Watch API
- ✅ **Rich Resource Model** - Support for 20+ Kubernetes resource types
- ✅ **Relationship Tracking** - Automatic discovery of resource dependencies
- ✅ **RBAC Integration** - Respect Kubernetes RBAC policies
- ✅ **GraphQL API** - Unified query interface for cluster operations
- ✅ **Metadata Catalog** - Track all resources in DictaMesh metadata catalog
- ✅ **Observability** - Full OpenTelemetry tracing and metrics
- ✅ **Resilience** - Circuit breakers, retries, and graceful degradation

---

## 2. Architecture

### 2.1 Layered Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    MANAGEMENT APPLICATIONS                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   Dashboard  │  │  GitOps      │  │  Cost Mgmt   │          │
│  │   (GraphQL)  │  │  (GraphQL)   │  │  (GraphQL)   │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
          └──────────────────┴──────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                   DICTAMESH CORE FRAMEWORK                       │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  GraphQL Gateway  │  Event Bus  │  Metadata Catalog       │ │
│  │  Observability    │  Governance │  Resilience             │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│              KUBERNETES CLUSTER API ADAPTER                      │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Resource Manager │  Event Watcher │  Mutation Handler   │ │
│  │  Cache Layer      │  RBAC Enforcer │  Relationship Mapper│ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│              KUBERNETES API CONNECTOR (client-go)                │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Clientset        │  Dynamic Client  │  Discovery Client  │ │
│  │  Watch Interface  │  Informers       │  REST Client       │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┬──────────────────────────────────┘
                               │
          ┌────────────────────┴────────────────────┐
          │                                         │
┌─────────▼────────────┐                 ┌─────────▼────────────┐
│  Kubernetes Cluster  │                 │  Kubernetes Cluster  │
│  (Production)        │                 │  (Staging)           │
└──────────────────────┘                 └──────────────────────┘
```

### 2.2 Component Layers

#### Layer 1: Kubernetes API Connector
- **Technology**: Kubernetes client-go library
- **Purpose**: Low-level interaction with Kubernetes API server
- **Responsibilities**:
  - Authentication (kubeconfig, service account, OIDC)
  - Connection pooling and keep-alive
  - API version negotiation
  - Error handling and retries

#### Layer 2: Kubernetes Cluster API Adapter
- **Technology**: Go, implements DictaMesh `adapter.Adapter` interface
- **Purpose**: Business logic for Kubernetes resource management
- **Responsibilities**:
  - Resource type mapping (k8s → DictaMesh entities)
  - Event streaming from Kubernetes Watch API
  - Resource relationship discovery
  - RBAC policy enforcement
  - Cache management (informers)

#### Layer 3: DictaMesh Core Framework
- **Purpose**: Framework services (event bus, metadata catalog, GraphQL gateway)
- **Integration Points**:
  - Event publishing to Kafka topics
  - Metadata catalog registration
  - GraphQL schema federation
  - Observability (tracing, metrics, logs)

#### Layer 4: Management Applications
- **Purpose**: User-facing applications consuming the adapter
- **Access Methods**:
  - GraphQL queries and mutations
  - Event subscriptions (Kafka consumers)
  - REST API (via GraphQL gateway)

---

## 3. Components

### 3.1 Package Structure

```
pkg/
└── adapter/
    └── kubernetes/
        ├── adapter.go              # Main adapter implementation
        ├── config.go               # Configuration structures
        ├── types.go                # Type definitions
        ├── errors.go               # Error types
        │
        ├── connector/              # Kubernetes API connector
        │   ├── client.go           # Client-go wrapper
        │   ├── clientset.go        # Typed clientset
        │   ├── dynamic.go          # Dynamic client
        │   ├── discovery.go        # API discovery
        │   └── auth.go             # Authentication helpers
        │
        ├── resources/              # Resource managers
        │   ├── base.go             # Base resource manager
        │   ├── pod.go              # Pod resource manager
        │   ├── deployment.go       # Deployment resource manager
        │   ├── service.go          # Service resource manager
        │   ├── node.go             # Node resource manager
        │   ├── namespace.go        # Namespace resource manager
        │   ├── configmap.go        # ConfigMap resource manager
        │   ├── secret.go           # Secret resource manager
        │   ├── statefulset.go      # StatefulSet resource manager
        │   ├── daemonset.go        # DaemonSet resource manager
        │   ├── job.go              # Job resource manager
        │   ├── cronjob.go          # CronJob resource manager
        │   ├── ingress.go          # Ingress resource manager
        │   ├── persistentvolume.go # PersistentVolume manager
        │   └── custom.go           # Custom resource support
        │
        ├── watcher/                # Event watcher
        │   ├── watcher.go          # Watch coordinator
        │   ├── informer.go         # Informer-based watcher
        │   ├── events.go           # Event transformation
        │   └── filters.go          # Event filtering
        │
        ├── mutations/              # Mutation handlers
        │   ├── create.go           # Resource creation
        │   ├── update.go           # Resource updates
        │   ├── patch.go            # Strategic merge patch
        │   ├── delete.go           # Resource deletion
        │   └── scale.go            # Scaling operations
        │
        ├── relationships/          # Relationship discovery
        │   ├── mapper.go           # Relationship mapper
        │   ├── owner_refs.go       # OwnerReference tracking
        │   ├── selectors.go        # Label/field selector matching
        │   └── graph.go            # Resource graph builder
        │
        ├── cache/                  # Caching layer
        │   ├── cache.go            # Cache interface
        │   ├── informer_cache.go   # Informer-based cache
        │   ├── redis_cache.go      # Redis L2 cache
        │   └── sync.go             # Cache synchronization
        │
        ├── rbac/                   # RBAC integration
        │   ├── enforcer.go         # RBAC enforcement
        │   ├── impersonation.go    # User impersonation
        │   └── authorization.go    # Authorization checks
        │
        └── graphql/                # GraphQL schema
            ├── schema.graphql      # GraphQL type definitions
            ├── resolvers.go        # GraphQL resolvers
            ├── mutations.go        # Mutation resolvers
            └── subscriptions.go    # Real-time subscriptions
```

### 3.2 Core Components

#### 3.2.1 Adapter (adapter.go)

```go
// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package kubernetes

import (
    "context"
    "sync"

    "github.com/click2-run/dictamesh/pkg/adapter"
    "github.com/click2-run/dictamesh/pkg/events"
    "github.com/click2-run/dictamesh/pkg/observability"
)

// Adapter implements the DictaMesh adapter interface for Kubernetes
type Adapter struct {
    // Base adapter (provides common functionality)
    *adapter.BaseAdapter

    // Configuration
    config *Config

    // Kubernetes clients (per cluster)
    clusters map[string]*ClusterClient
    clustersMu sync.RWMutex

    // Resource managers
    resourceManagers map[string]ResourceManager

    // Event watcher
    watcher *EventWatcher

    // Relationship mapper
    relationshipMapper *RelationshipMapper

    // RBAC enforcer
    rbacEnforcer *RBACEnforcer

    // State
    initialized bool
    running bool
    mu sync.RWMutex
}

// NewAdapter creates a new Kubernetes adapter
func NewAdapter(obs *observability.Observability) (*Adapter, error)

// Initialize initializes the adapter
func (a *Adapter) Initialize(ctx context.Context, config adapter.Config) error

// Start starts the adapter (begins watching clusters)
func (a *Adapter) Start(ctx context.Context) error

// Stop stops the adapter gracefully
func (a *Adapter) Stop(ctx context.Context) error

// Health performs health check across all clusters
func (a *Adapter) Health(ctx context.Context) (*adapter.HealthStatus, error)

// GetCapabilities returns adapter capabilities
func (a *Adapter) GetCapabilities() []adapter.Capability

// AddCluster dynamically adds a new cluster
func (a *Adapter) AddCluster(ctx context.Context, cluster ClusterConfig) error

// RemoveCluster removes a cluster
func (a *Adapter) RemoveCluster(ctx context.Context, clusterID string) error

// ListClusters lists all connected clusters
func (a *Adapter) ListClusters(ctx context.Context) ([]*ClusterInfo, error)
```

#### 3.2.2 Configuration (config.go)

```go
// Config represents the Kubernetes adapter configuration
type Config struct {
    // Clusters to connect to
    Clusters []ClusterConfig `json:"clusters" yaml:"clusters"`

    // Default namespace for operations
    DefaultNamespace string `json:"default_namespace" yaml:"default_namespace"`

    // Resource types to watch (empty = all)
    WatchResources []string `json:"watch_resources" yaml:"watch_resources"`

    // Enable resource caching
    EnableCache bool `json:"enable_cache" yaml:"enable_cache"`

    // Cache TTL
    CacheTTL time.Duration `json:"cache_ttl" yaml:"cache_ttl"`

    // Enable RBAC enforcement
    EnableRBAC bool `json:"enable_rbac" yaml:"enable_rbac"`

    // Enable relationship discovery
    EnableRelationships bool `json:"enable_relationships" yaml:"enable_relationships"`

    // Resync period for informers
    ResyncPeriod time.Duration `json:"resync_period" yaml:"resync_period"`

    // Worker pool size for event processing
    WorkerPoolSize int `json:"worker_pool_size" yaml:"worker_pool_size"`

    // Enable mutations (create, update, delete)
    EnableMutations bool `json:"enable_mutations" yaml:"enable_mutations"`

    // Enable custom resources
    EnableCustomResources bool `json:"enable_custom_resources" yaml:"enable_custom_resources"`

    // Custom resource definitions to watch
    CustomResourceDefinitions []CRDConfig `json:"custom_resource_definitions" yaml:"custom_resource_definitions"`
}

// ClusterConfig represents configuration for a single cluster
type ClusterConfig struct {
    // Unique cluster identifier
    ID string `json:"id" yaml:"id"`

    // Human-readable cluster name
    Name string `json:"name" yaml:"name"`

    // Cluster environment (dev, staging, prod)
    Environment string `json:"environment" yaml:"environment"`

    // Cluster region/zone
    Region string `json:"region" yaml:"region"`

    // Authentication method
    AuthMethod AuthMethod `json:"auth_method" yaml:"auth_method"`

    // Kubeconfig path (for kubeconfig auth)
    KubeconfigPath string `json:"kubeconfig_path" yaml:"kubeconfig_path"`

    // Kubeconfig context (for kubeconfig auth)
    KubeconfigContext string `json:"kubeconfig_context" yaml:"kubeconfig_context"`

    // Service account token (for token auth)
    ServiceAccountToken string `json:"service_account_token" yaml:"service_account_token"`

    // API server URL
    APIServerURL string `json:"api_server_url" yaml:"api_server_url"`

    // TLS configuration
    TLSConfig TLSConfig `json:"tls_config" yaml:"tls_config"`

    // Rate limiting
    QPS float32 `json:"qps" yaml:"qps"`
    Burst int `json:"burst" yaml:"burst"`

    // Timeout
    Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

// AuthMethod represents authentication methods
type AuthMethod string

const (
    AuthMethodKubeconfig      AuthMethod = "kubeconfig"
    AuthMethodServiceAccount  AuthMethod = "service_account"
    AuthMethodToken           AuthMethod = "token"
    AuthMethodOIDC            AuthMethod = "oidc"
    AuthMethodInCluster       AuthMethod = "in_cluster"
)
```

#### 3.2.3 Resource Manager (resources/base.go)

```go
// ResourceManager defines the interface for managing Kubernetes resources
type ResourceManager interface {
    // GetResourceType returns the resource type (e.g., "pod", "deployment")
    GetResourceType() string

    // GetGroupVersionKind returns the GVK for this resource
    GetGroupVersionKind() schema.GroupVersionKind

    // Get retrieves a resource by namespace and name
    Get(ctx context.Context, cluster, namespace, name string) (*adapter.Resource, error)

    // List lists resources in a namespace (or cluster-wide if namespace is empty)
    List(ctx context.Context, cluster, namespace string, opts *adapter.ListOptions) (*adapter.ResourceList, error)

    // Watch streams resource changes
    Watch(ctx context.Context, cluster, namespace string) (<-chan *adapter.Event, error)

    // Create creates a new resource
    Create(ctx context.Context, cluster string, resource *adapter.Resource) (*adapter.Resource, error)

    // Update updates an existing resource
    Update(ctx context.Context, cluster string, resource *adapter.Resource) (*adapter.Resource, error)

    // Patch patches a resource (strategic merge patch)
    Patch(ctx context.Context, cluster, namespace, name string, patch []byte) (*adapter.Resource, error)

    // Delete deletes a resource
    Delete(ctx context.Context, cluster, namespace, name string) error

    // ToEntity converts Kubernetes resource to DictaMesh entity
    ToEntity(obj interface{}) (*adapter.Resource, error)

    // FromEntity converts DictaMesh entity to Kubernetes resource
    FromEntity(entity *adapter.Resource) (interface{}, error)

    // GetRelationships discovers relationships for this resource
    GetRelationships(ctx context.Context, cluster, namespace, name string) ([]*Relationship, error)
}
```

---

## 4. Resource Types

### 4.1 Supported Kubernetes Resources

The adapter will support the following Kubernetes resource types:

#### Core Resources (v1)
- **Pod** - Running containers
- **Service** - Service endpoints
- **Node** - Cluster nodes
- **Namespace** - Logical partitions
- **ConfigMap** - Configuration data
- **Secret** - Sensitive data
- **PersistentVolume** - Storage volumes
- **PersistentVolumeClaim** - Volume claims
- **ServiceAccount** - Service identities
- **Event** - Cluster events

#### Apps Resources (apps/v1)
- **Deployment** - Declarative pod management
- **StatefulSet** - Stateful applications
- **DaemonSet** - Node-level services
- **ReplicaSet** - Pod replication

#### Batch Resources (batch/v1)
- **Job** - One-time tasks
- **CronJob** - Scheduled tasks

#### Networking Resources (networking.k8s.io/v1)
- **Ingress** - HTTP routing
- **NetworkPolicy** - Network security

#### RBAC Resources (rbac.authorization.k8s.io/v1)
- **Role** - Namespace-scoped permissions
- **ClusterRole** - Cluster-scoped permissions
- **RoleBinding** - Role assignments
- **ClusterRoleBinding** - ClusterRole assignments

#### Custom Resources
- **CustomResourceDefinition** - CRD definitions
- **Dynamic Custom Resources** - Runtime CRD support

### 4.2 Resource Mapping

Each Kubernetes resource maps to a DictaMesh entity:

```go
// Entity mapping example for Pod
type PodEntity struct {
    // Standard entity fields
    ID        string                 `json:"id"`         // cluster:namespace:name
    Type      string                 `json:"type"`       // "kubernetes.pod"

    // Kubernetes-specific attributes
    Attributes map[string]interface{} `json:"attributes"`
    /*
    {
        "cluster": "prod-us-east-1",
        "namespace": "default",
        "name": "nginx-7c6f9d8b9-abc12",
        "uid": "12345678-1234-1234-1234-123456789abc",
        "labels": {"app": "nginx", "version": "1.0"},
        "annotations": {...},
        "phase": "Running",
        "node_name": "node-1",
        "pod_ip": "10.244.1.5",
        "host_ip": "192.168.1.10",
        "start_time": "2025-11-08T10:00:00Z",
        "containers": [
            {
                "name": "nginx",
                "image": "nginx:1.21",
                "state": "running",
                "ready": true
            }
        ],
        "conditions": [...],
        "resource_requests": {
            "cpu": "100m",
            "memory": "128Mi"
        },
        "resource_limits": {
            "cpu": "500m",
            "memory": "512Mi"
        }
    }
    */

    // Relationships
    Relationships map[string]interface{} `json:"relationships"`
    /*
    {
        "owned_by": {
            "type": "kubernetes.replicaset",
            "id": "cluster:namespace:replicaset-name"
        },
        "runs_on": {
            "type": "kubernetes.node",
            "id": "cluster::node-1"
        },
        "uses_services": [
            {
                "type": "kubernetes.service",
                "id": "cluster:namespace:service-name"
            }
        ],
        "uses_configmaps": [...],
        "uses_secrets": [...]
    }
    */

    // Metadata
    Metadata *adapter.ResourceMetadata `json:"metadata"`
}
```

---

## 5. API Design

### 5.1 GraphQL Schema

```graphql
# Cluster
type Cluster {
  id: ID!
  name: String!
  environment: String!
  region: String
  version: String!
  status: ClusterStatus!
  nodeCount: Int!
  namespaceCount: Int!
  health: HealthStatus!
  createdAt: DateTime!

  # Relationships
  nodes: [Node!]!
  namespaces: [Namespace!]!

  # Metrics
  metrics: ClusterMetrics
}

type ClusterStatus {
  ready: Boolean!
  message: String
  components: [ComponentStatus!]!
}

type ClusterMetrics {
  cpuUsage: Float!
  memoryUsage: Float!
  podCount: Int!
  storageUsage: Float!
}

# Node
type Node {
  id: ID!
  name: String!
  cluster: Cluster!
  labels: JSON!
  annotations: JSON!

  # Status
  status: NodeStatus!
  ready: Boolean!
  schedulable: Boolean!

  # Capacity
  capacity: ResourceQuantity!
  allocatable: ResourceQuantity!

  # Info
  osImage: String!
  kernelVersion: String!
  containerRuntimeVersion: String!
  kubeletVersion: String!

  # Relationships
  pods: [Pod!]!

  # Metrics
  metrics: NodeMetrics

  createdAt: DateTime!
}

type ResourceQuantity {
  cpu: String!
  memory: String!
  storage: String!
  pods: String!
}

# Namespace
type Namespace {
  id: ID!
  name: String!
  cluster: Cluster!
  labels: JSON!
  annotations: JSON!
  status: NamespaceStatus!

  # Quotas
  resourceQuota: ResourceQuota

  # Relationships
  pods: [Pod!]!
  deployments: [Deployment!]!
  services: [Service!]!
  configMaps: [ConfigMap!]!
  secrets: [Secret!]!

  createdAt: DateTime!
}

# Pod
type Pod {
  id: ID!
  name: String!
  namespace: Namespace!
  cluster: Cluster!
  uid: String!
  labels: JSON!
  annotations: JSON!

  # Status
  phase: PodPhase!
  reason: String
  message: String
  conditions: [PodCondition!]!

  # Spec
  nodeName: String
  serviceAccountName: String
  restartPolicy: RestartPolicy!

  # Networking
  podIP: String
  hostIP: String

  # Containers
  containers: [Container!]!
  initContainers: [Container!]!

  # Resources
  resources: PodResources!

  # Relationships
  node: Node
  ownedBy: Resource  # ReplicaSet, StatefulSet, DaemonSet, Job
  services: [Service!]!
  configMaps: [ConfigMap!]!
  secrets: [Secret!]!
  volumes: [Volume!]!

  # Metrics
  metrics: PodMetrics

  createdAt: DateTime!
  startedAt: DateTime
  deletedAt: DateTime
}

enum PodPhase {
  PENDING
  RUNNING
  SUCCEEDED
  FAILED
  UNKNOWN
}

type Container {
  name: String!
  image: String!
  state: ContainerState!
  ready: Boolean!
  restartCount: Int!
  resources: ContainerResources!
}

type ContainerState {
  state: String!  # running, waiting, terminated
  reason: String
  message: String
  startedAt: DateTime
  finishedAt: DateTime
  exitCode: Int
}

# Deployment
type Deployment {
  id: ID!
  name: String!
  namespace: Namespace!
  cluster: Cluster!
  labels: JSON!
  annotations: JSON!

  # Spec
  replicas: Int!
  selector: LabelSelector!
  template: PodTemplateSpec!
  strategy: DeploymentStrategy!

  # Status
  availableReplicas: Int!
  readyReplicas: Int!
  updatedReplicas: Int!
  conditions: [DeploymentCondition!]!

  # Relationships
  replicaSets: [ReplicaSet!]!
  pods: [Pod!]!

  createdAt: DateTime!
  updatedAt: DateTime
}

# Service
type Service {
  id: ID!
  name: String!
  namespace: Namespace!
  cluster: Cluster!
  labels: JSON!
  annotations: JSON!

  # Spec
  type: ServiceType!
  selector: JSON!
  ports: [ServicePort!]!
  clusterIP: String
  externalIPs: [String!]!
  loadBalancerIP: String

  # Relationships
  endpoints: Endpoints!
  pods: [Pod!]!
  ingresses: [Ingress!]!

  createdAt: DateTime!
}

enum ServiceType {
  CLUSTER_IP
  NODE_PORT
  LOAD_BALANCER
  EXTERNAL_NAME
}

# Query operations
type Query {
  # Cluster queries
  cluster(id: ID!): Cluster
  clusters(filter: ClusterFilter): [Cluster!]!

  # Resource queries
  node(cluster: String!, name: String!): Node
  nodes(cluster: String!, filter: NodeFilter): [Node!]!

  namespace(cluster: String!, name: String!): Namespace
  namespaces(cluster: String!, filter: NamespaceFilter): [Namespace!]!

  pod(cluster: String!, namespace: String!, name: String!): Pod
  pods(cluster: String!, namespace: String, filter: PodFilter): [Pod!]!

  deployment(cluster: String!, namespace: String!, name: String!): Deployment
  deployments(cluster: String!, namespace: String, filter: DeploymentFilter): [Deployment!]!

  service(cluster: String!, namespace: String!, name: String!): Service
  services(cluster: String!, namespace: String, filter: ServiceFilter): [Service!]!

  # Cross-cluster queries
  allPods(filter: PodFilter): [Pod!]!
  allDeployments(filter: DeploymentFilter): [Deployment!]!
  allServices(filter: ServiceFilter): [Service!]!

  # Relationship queries
  resourceGraph(
    cluster: String!
    namespace: String
    resourceType: String!
    name: String!
    depth: Int
  ): ResourceGraph!
}

# Mutation operations
type Mutation {
  # Pod operations
  createPod(cluster: String!, namespace: String!, pod: PodInput!): Pod!
  updatePod(cluster: String!, namespace: String!, name: String!, pod: PodInput!): Pod!
  deletePod(cluster: String!, namespace: String!, name: String!): Boolean!

  # Deployment operations
  createDeployment(cluster: String!, namespace: String!, deployment: DeploymentInput!): Deployment!
  updateDeployment(cluster: String!, namespace: String!, name: String!, deployment: DeploymentInput!): Deployment!
  scaleDeployment(cluster: String!, namespace: String!, name: String!, replicas: Int!): Deployment!
  rolloutRestart(cluster: String!, namespace: String!, name: String!): Deployment!
  deleteDeployment(cluster: String!, namespace: String!, name: String!): Boolean!

  # Service operations
  createService(cluster: String!, namespace: String!, service: ServiceInput!): Service!
  updateService(cluster: String!, namespace: String!, name: String!, service: ServiceInput!): Service!
  deleteService(cluster: String!, namespace: String!, name: String!): Boolean!

  # ConfigMap operations
  createConfigMap(cluster: String!, namespace: String!, configMap: ConfigMapInput!): ConfigMap!
  updateConfigMap(cluster: String!, namespace: String!, name: String!, configMap: ConfigMapInput!): ConfigMap!
  deleteConfigMap(cluster: String!, namespace: String!, name: String!): Boolean!

  # Namespace operations
  createNamespace(cluster: String!, namespace: NamespaceInput!): Namespace!
  deleteNamespace(cluster: String!, name: String!): Boolean!
}

# Subscription operations (real-time events)
type Subscription {
  # Pod events
  podEvents(cluster: String, namespace: String, labelSelector: String): PodEvent!

  # Deployment events
  deploymentEvents(cluster: String, namespace: String, labelSelector: String): DeploymentEvent!

  # Service events
  serviceEvents(cluster: String, namespace: String, labelSelector: String): ServiceEvent!

  # Generic resource events
  resourceEvents(
    cluster: String
    namespace: String
    resourceType: String
    labelSelector: String
  ): ResourceEvent!

  # Cluster events
  clusterEvents(cluster: String): ClusterEvent!
}

type PodEvent {
  type: EventType!
  pod: Pod!
  timestamp: DateTime!
  reason: String
  message: String
}

enum EventType {
  CREATED
  UPDATED
  DELETED
}
```

### 5.2 REST API Endpoints

For applications that prefer REST over GraphQL:

```
# Clusters
GET    /api/v1/clusters
GET    /api/v1/clusters/:cluster_id
POST   /api/v1/clusters
DELETE /api/v1/clusters/:cluster_id

# Nodes
GET    /api/v1/clusters/:cluster_id/nodes
GET    /api/v1/clusters/:cluster_id/nodes/:node_name

# Namespaces
GET    /api/v1/clusters/:cluster_id/namespaces
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace_name
POST   /api/v1/clusters/:cluster_id/namespaces
DELETE /api/v1/clusters/:cluster_id/namespaces/:namespace_name

# Pods
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace/pods
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace/pods/:pod_name
POST   /api/v1/clusters/:cluster_id/namespaces/:namespace/pods
PUT    /api/v1/clusters/:cluster_id/namespaces/:namespace/pods/:pod_name
DELETE /api/v1/clusters/:cluster_id/namespaces/:namespace/pods/:pod_name
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace/pods/:pod_name/logs

# Deployments
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace/deployments
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace/deployments/:deployment_name
POST   /api/v1/clusters/:cluster_id/namespaces/:namespace/deployments
PUT    /api/v1/clusters/:cluster_id/namespaces/:namespace/deployments/:deployment_name
PATCH  /api/v1/clusters/:cluster_id/namespaces/:namespace/deployments/:deployment_name/scale
POST   /api/v1/clusters/:cluster_id/namespaces/:namespace/deployments/:deployment_name/restart
DELETE /api/v1/clusters/:cluster_id/namespaces/:namespace/deployments/:deployment_name

# Services
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace/services
GET    /api/v1/clusters/:cluster_id/namespaces/:namespace/services/:service_name
POST   /api/v1/clusters/:cluster_id/namespaces/:namespace/services
PUT    /api/v1/clusters/:cluster_id/namespaces/:namespace/services/:service_name
DELETE /api/v1/clusters/:cluster_id/namespaces/:namespace/services/:service_name

# Events (SSE for streaming)
GET    /api/v1/events/stream?cluster=&namespace=&resource_type=
```

---

## 6. Event Streaming

### 6.1 Kubernetes Watch Integration

The adapter uses Kubernetes Watch API to stream real-time events:

```go
// EventWatcher manages Kubernetes event watching
type EventWatcher struct {
    clusters map[string]*ClusterWatcher
    publisher *events.Publisher
    filters []EventFilter
}

// ClusterWatcher watches resources in a single cluster
type ClusterWatcher struct {
    clusterID string
    informers map[string]cache.SharedIndexInformer
    stopCh chan struct{}
}

// Event types published to Kafka
const (
    TopicK8sResourceCreated = "kubernetes.resource.created"
    TopicK8sResourceUpdated = "kubernetes.resource.updated"
    TopicK8sResourceDeleted = "kubernetes.resource.deleted"
    TopicK8sClusterEvent    = "kubernetes.cluster.event"
)
```

### 6.2 Event Schema

```json
{
  "event_id": "uuid",
  "event_type": "kubernetes.resource.created",
  "timestamp": "2025-11-08T10:00:00Z",
  "cluster_id": "prod-us-east-1",
  "resource": {
    "type": "pod",
    "namespace": "default",
    "name": "nginx-7c6f9d8b9-abc12",
    "uid": "12345678-1234-1234-1234-123456789abc",
    "group_version_kind": {
      "group": "",
      "version": "v1",
      "kind": "Pod"
    }
  },
  "old_object": null,
  "new_object": {
    "spec": {...},
    "status": {...}
  },
  "diff": {
    "status.phase": {
      "old": "Pending",
      "new": "Running"
    }
  },
  "metadata": {
    "trace_id": "trace-id",
    "span_id": "span-id",
    "source": "kubernetes-adapter",
    "version": "1.0.0"
  }
}
```

### 6.3 Event Filtering

Support for flexible event filtering:

```go
type EventFilter struct {
    Clusters []string
    Namespaces []string
    ResourceTypes []string
    LabelSelector string
    FieldSelector string
    MinimumSeverity string
}
```

---

## 7. Security & Governance

### 7.1 RBAC Integration

The adapter respects Kubernetes RBAC policies:

```go
// RBACEnforcer enforces Kubernetes RBAC policies
type RBACEnforcer struct {
    authClient authorizationv1.AuthorizationV1Interface
}

// CheckAccess verifies if a user can perform an action
func (e *RBACEnforcer) CheckAccess(
    ctx context.Context,
    cluster string,
    user string,
    verb string,
    resource schema.GroupVersionResource,
    namespace string,
    name string,
) (bool, error) {
    // Use SubjectAccessReview API
    sar := &authorizationv1.SubjectAccessReview{
        Spec: authorizationv1.SubjectAccessReviewSpec{
            User: user,
            ResourceAttributes: &authorizationv1.ResourceAttributes{
                Namespace: namespace,
                Verb:      verb,
                Group:     resource.Group,
                Version:   resource.Version,
                Resource:  resource.Resource,
                Name:      name,
            },
        },
    }

    result, err := e.authClient.SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
    return result.Status.Allowed, err
}
```

### 7.2 Sensitive Data Handling

- **Secrets**: Never expose secret values in GraphQL/REST APIs
- **PII Marking**: Mark pods/services with PII labels for governance
- **Audit Logging**: Log all mutations to audit trail
- **Encryption**: Support encrypted etcd for sensitive data

### 7.3 Multi-Tenancy

- Namespace isolation
- User impersonation for access control
- Quota enforcement
- Network policies

---

## 8. Implementation Plan

### Phase 1: Foundation (Weeks 1-2)
- [ ] Package structure setup
- [ ] Connector layer (client-go integration)
- [ ] Basic adapter implementation
- [ ] Configuration management
- [ ] Unit tests for connector

### Phase 2: Core Resources (Weeks 3-4)
- [ ] Pod resource manager
- [ ] Deployment resource manager
- [ ] Service resource manager
- [ ] Node resource manager
- [ ] Namespace resource manager
- [ ] Resource caching (informers)

### Phase 3: Event Streaming (Week 5)
- [ ] Event watcher implementation
- [ ] Kafka integration
- [ ] Event filtering
- [ ] Event transformation

### Phase 4: Relationships (Week 6)
- [ ] Relationship mapper
- [ ] OwnerReference tracking
- [ ] Label selector matching
- [ ] Resource graph builder

### Phase 5: GraphQL API (Weeks 7-8)
- [ ] GraphQL schema definition
- [ ] Query resolvers
- [ ] Mutation resolvers
- [ ] Subscription resolvers (real-time events)
- [ ] DataLoader integration

### Phase 6: Advanced Resources (Weeks 9-10)
- [ ] StatefulSet, DaemonSet, Job, CronJob
- [ ] ConfigMap, Secret
- [ ] Ingress, NetworkPolicy
- [ ] RBAC resources
- [ ] PersistentVolume, PersistentVolumeClaim

### Phase 7: Mutations & RBAC (Week 11)
- [ ] Create operations
- [ ] Update operations
- [ ] Patch operations
- [ ] Delete operations
- [ ] RBAC enforcement
- [ ] User impersonation

### Phase 8: Custom Resources (Week 12)
- [ ] CRD support
- [ ] Dynamic client integration
- [ ] Generic resource handlers

### Phase 9: Multi-Cluster (Week 13)
- [ ] Dynamic cluster addition/removal
- [ ] Cross-cluster queries
- [ ] Cluster health monitoring

### Phase 10: Testing & Documentation (Week 14)
- [ ] Integration tests
- [ ] E2E tests
- [ ] Performance benchmarks
- [ ] API documentation
- [ ] Deployment guides

### Phase 11: Production Hardening (Week 15-16)
- [ ] Production deployment
- [ ] Performance optimization
- [ ] Security audit
- [ ] Load testing
- [ ] Chaos engineering tests

---

## 9. Deployment Strategy

### 9.1 Deployment Architecture

```
┌──────────────────────────────────────────────────────────┐
│              Kubernetes Cluster (Management)             │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │   DictaMesh Kubernetes Adapter                     │ │
│  │   - Deployment (3 replicas)                        │ │
│  │   - Service (ClusterIP)                            │ │
│  │   - ServiceAccount + RBAC                          │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │   DictaMesh Core Services                          │ │
│  │   - GraphQL Gateway                                │ │
│  │   - Metadata Catalog                               │ │
│  │   - Event Router                                   │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │   Infrastructure                                   │ │
│  │   - Kafka/Redpanda                                 │ │
│  │   - PostgreSQL                                     │ │
│  │   - Redis                                          │ │
│  └────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────┘
                          │
          ┌───────────────┴───────────────┐
          │                               │
┌─────────▼─────────┐           ┌─────────▼─────────┐
│  Cluster (Prod)   │           │  Cluster (Dev)    │
│  - ServiceAccount │           │  - ServiceAccount │
│  - ClusterRole    │           │  - ClusterRole    │
└───────────────────┘           └───────────────────┘
```

### 9.2 RBAC Configuration

```yaml
# ServiceAccount for the adapter
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dictamesh-k8s-adapter
  namespace: dictamesh

---
# ClusterRole with necessary permissions
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: dictamesh-k8s-adapter
rules:
# Read all resources
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["get", "list", "watch"]

# Mutations (optional, can be restricted)
- apiGroups: [""]
  resources: ["pods", "services", "configmaps", "secrets", "namespaces"]
  verbs: ["create", "update", "patch", "delete"]

- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets", "daemonsets", "replicasets"]
  verbs: ["create", "update", "patch", "delete"]

# RBAC checks
- apiGroups: ["authorization.k8s.io"]
  resources: ["subjectaccessreviews"]
  verbs: ["create"]

---
# ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: dictamesh-k8s-adapter
subjects:
- kind: ServiceAccount
  name: dictamesh-k8s-adapter
  namespace: dictamesh
roleRef:
  kind: ClusterRole
  name: dictamesh-k8s-adapter
  apiGroup: rbac.authorization.k8s.io
```

### 9.3 Deployment Manifest

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dictamesh-k8s-adapter
  namespace: dictamesh
  labels:
    app: dictamesh-k8s-adapter
spec:
  replicas: 3
  selector:
    matchLabels:
      app: dictamesh-k8s-adapter
  template:
    metadata:
      labels:
        app: dictamesh-k8s-adapter
    spec:
      serviceAccountName: dictamesh-k8s-adapter

      containers:
      - name: adapter
        image: dictamesh/kubernetes-adapter:v1.0.0
        imagePullPolicy: IfNotPresent

        env:
        # Cluster configurations
        - name: CLUSTERS_CONFIG
          valueFrom:
            configMapKeyRef:
              name: k8s-adapter-config
              key: clusters.yaml

        # Kafka configuration
        - name: KAFKA_BOOTSTRAP_SERVERS
          value: "kafka:9092"

        # PostgreSQL configuration
        - name: POSTGRES_DSN
          valueFrom:
            secretKeyRef:
              name: postgres-credentials
              key: dsn

        # Redis configuration
        - name: REDIS_URL
          value: "redis://redis:6379"

        # OpenTelemetry
        - name: OTEL_EXPORTER_OTLP_ENDPOINT
          value: "http://otel-collector:4317"

        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics

        resources:
          requests:
            cpu: "500m"
            memory: "512Mi"
          limits:
            cpu: "2000m"
            memory: "2Gi"

        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10

        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
```

---

## 10. Testing Strategy

### 10.1 Unit Tests
- Resource manager tests
- Relationship mapper tests
- RBAC enforcer tests
- Event transformation tests
- Cache layer tests

### 10.2 Integration Tests
- Full adapter lifecycle
- Multi-cluster scenarios
- Event streaming end-to-end
- GraphQL query execution
- Mutation operations

### 10.3 E2E Tests
- Deploy adapter to test cluster
- Create/update/delete resources via GraphQL
- Verify events published to Kafka
- Verify metadata catalog updates
- Test cross-cluster queries

### 10.4 Performance Tests
- Benchmark resource listing (10k+ pods)
- Measure event processing latency
- Test concurrent GraphQL queries
- Cache hit/miss ratios
- Memory usage under load

### 10.5 Chaos Tests
- Cluster connectivity loss
- API server unavailability
- Kafka outage
- PostgreSQL failure
- Network partitions

---

## 11. Success Metrics

### 11.1 Performance Targets
- **Query Latency**: P99 < 200ms for resource queries
- **Event Latency**: < 500ms from K8s event to Kafka publish
- **Throughput**: 1000+ resources/sec listing
- **Memory Usage**: < 1GB per 10k watched resources
- **CPU Usage**: < 50% under normal load

### 11.2 Reliability Targets
- **Availability**: 99.9% uptime
- **Error Rate**: < 0.1% for queries
- **Event Loss**: 0% (guaranteed delivery to Kafka)
- **Cache Hit Rate**: > 80%

### 11.3 Scalability Targets
- **Clusters**: Support 100+ clusters per adapter instance
- **Resources**: Watch 100k+ resources simultaneously
- **Replicas**: Horizontal scaling to 10+ replicas
- **Namespaces**: No limit on namespace count

---

## 12. Future Enhancements

### 12.1 Advanced Features
- **GitOps Integration**: Automated drift detection and reconciliation
- **Policy Engine**: OPA integration for policy enforcement
- **Cost Tracking**: Track resource costs and allocation
- **Capacity Planning**: Predict resource needs
- **Auto-Scaling**: Intelligent HPA recommendations

### 12.2 Additional Integrations
- **Prometheus**: Direct metrics scraping
- **Helm**: Chart deployment tracking
- **ArgoCD**: GitOps workflow integration
- **Vault**: Secrets management
- **Service Mesh**: Istio/Linkerd integration

---

## Appendix A: Technology Stack

- **Language**: Go 1.21+
- **Kubernetes Client**: client-go v0.28+
- **GraphQL**: gqlgen
- **Event Bus**: Kafka/Redpanda
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Observability**: OpenTelemetry, Prometheus, Jaeger
- **Testing**: testify, gomock, kind (for E2E)

## Appendix B: References

- Kubernetes API Conventions: https://kubernetes.io/docs/reference/using-api/api-concepts/
- client-go Documentation: https://pkg.go.dev/k8s.io/client-go
- Kubernetes Informers: https://pkg.go.dev/k8s.io/client-go/informers
- GraphQL Federation: https://www.apollographql.com/docs/federation/
- DictaMesh Architecture: /docs/PROJECT-SCOPE.md

---

**Document Status**: ✅ Ready for Implementation
**Next Steps**: Begin Phase 1 implementation

