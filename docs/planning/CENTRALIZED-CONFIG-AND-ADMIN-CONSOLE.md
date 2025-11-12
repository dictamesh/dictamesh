# Centralized Configuration and Admin Console Design

**Status:** Design Document
**Version:** 2.0 (Updated with Remix + Refine.dev)
**Date:** 2025-11-08
**Author:** Claude AI Assistant

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Problem Statement](#problem-statement)
3. [Current State Analysis](#current-state-analysis)
4. [Proposed Solution](#proposed-solution)
5. [Architecture Design](#architecture-design)
6. [Technology Stack](#technology-stack)
7. [Implementation Plan](#implementation-plan)
8. [Security Considerations](#security-considerations)
9. [Deployment Strategy](#deployment-strategy)
10. [Development Workflow](#development-workflow)

---

## Executive Summary

This document outlines the design and implementation of a **centralized configuration management system** and **web-based admin console** for the DictaMesh framework using **Remix** (full-stack framework) and **Refine.dev** (enterprise admin panel framework).

### Key Components

1. **Configuration Service** - Secure, versioned configuration storage with encryption
2. **Management API** - RESTful API for configuration CRUD operations (Go backend)
3. **Admin Console** - Modern web UI built with Remix + Refine.dev for operations, monitoring, and configuration management
4. **Secrets Management** - Encrypted credential storage with access control

### Goals

- ✅ Eliminate dependency on `.env` files for production
- ✅ Provide secure, encrypted storage for credentials and sensitive configuration
- ✅ Enable centralized management of multi-environment configurations
- ✅ Offer a full-featured admin console for health, monitoring, and deployment management
- ✅ Support development mode with Hot Module Reload (HMR)
- ✅ Integrate with existing observability stack (Prometheus, Grafana, Jaeger, Sentry)
- ✅ Leverage Remix for server-side rendering and progressive enhancement
- ✅ Use Refine.dev for rapid admin panel development with best practices

---

## Problem Statement

### Current Challenges

1. **Configuration Fragmentation**
   - Configuration spread across environment variables, Docker Compose files, and Go code defaults
   - No single source of truth for configuration
   - Difficult to manage multi-environment configurations (dev, staging, prod)

2. **Security Risks**
   - Credentials in environment variables or `.env` files
   - No encryption at rest for sensitive data
   - Limited audit trail for configuration changes
   - Risk of committing secrets to version control

3. **Operational Overhead**
   - Manual configuration management prone to human error
   - No centralized view of system health and configuration
   - Difficult to troubleshoot configuration-related issues
   - No versioning or rollback capability for configurations

4. **Lack of Admin Tooling**
   - No unified interface for operations team
   - Service health scattered across multiple UIs (Grafana, Jaeger, Redpanda Console, Sentry)
   - No deployment management capabilities
   - Limited visibility into framework internals

---

## Current State Analysis

### Existing Infrastructure (from develop branch)

**Strengths:**
- ✅ Comprehensive database infrastructure (PostgreSQL with migrations)
- ✅ Complete observability stack (Prometheus, Grafana, Jaeger, **Sentry**)
- ✅ Event bus infrastructure (Redpanda/Kafka)
- ✅ **Notifications service** (multi-channel: email, SMS, Slack, webhooks, etc.)
- ✅ Well-architected monorepo with Go Workspaces
- ✅ Docker Compose development environment

**New Features from develop:**

1. **Sentry Integration** (`docs/SENTRY-INTEGRATION.md`)
   - Self-hosted Sentry for error tracking
   - Available at http://localhost:9000
   - Integrated with ClickHouse for event storage
   - Support for Go, Node.js, Python SDKs
   - Performance monitoring and distributed tracing

2. **Notifications Service** (`pkg/notifications/`)
   - Multi-channel support (Email, SMS, Push, Slack, Webhooks, In-App, PagerDuty)
   - Template engine with Go templates
   - Event-driven with Kafka/Redpanda integration
   - User preferences and rate limiting
   - Database tables: `dictamesh_notifications`, `dictamesh_notification_templates`, etc.
   - Migration: `000003_add_notifications.up.sql`

**Configuration Management Today:**
```
1. Go Code Defaults
   └─> pkg/database/config.go (Config struct with defaults)
   └─> pkg/notifications/config.go (Notifications config)

2. Environment Variables
   └─> Docker Compose environment sections
   └─> Kubernetes ConfigMaps (planned, not implemented)

3. Infrastructure as Code
   └─> docker-compose.dev.yml (updated with Sentry services)
   └─> k8s manifests (Sentry production setup ready)
```

**Gaps Identified:**
- ❌ No centralized configuration service
- ❌ No configuration versioning or audit trail
- ❌ No secrets management solution
- ❌ No admin console for operations
- ❌ No unified health monitoring UI
- ❌ No deployment management tooling

---

## Proposed Solution

### Solution Overview

We will implement a **three-tier solution** using **Remix + Refine.dev**:

```
┌─────────────────────────────────────────────────────────────┐
│          Admin Console (Remix + Refine.dev Web UI)          │
│  - Configuration Management    - Health Monitoring           │
│  - Deployment Dashboard        - Metrics Integration         │
│  - Audit Logs                  - Service Control             │
│  - Notifications Management    - Sentry Error Dashboard      │
└────────────────────┬────────────────────────────────────────┘
                     │ HTTPS/REST API + Remix Loaders/Actions
┌────────────────────▼────────────────────────────────────────┐
│        Configuration Management API (Go Backend)             │
│  - Config CRUD Operations      - Secrets Management          │
│  - Version Control             - Access Control (RBAC)       │
│  - Health Endpoints            - Audit Logging               │
│  - Integration with Sentry     - Notifications Config        │
└────────────────────┬────────────────────────────────────────┘
                     │ Internal
┌────────────────────▼────────────────────────────────────────┐
│                Configuration Storage Layer                   │
│  PostgreSQL + Encryption                                     │
│  - dictamesh_configurations                                  │
│  - dictamesh_secrets (encrypted)                             │
│  - dictamesh_config_versions                                 │
│  - dictamesh_config_audit_logs                               │
│  - Integration with existing tables (notifications, etc.)    │
└─────────────────────────────────────────────────────────────┘
```

### Why Remix + Refine.dev?

**Remix Benefits:**
- 🚀 **Server-Side Rendering (SSR)** - Fast initial page loads, SEO-friendly
- 🔄 **Progressive Enhancement** - Works without JavaScript, then enhances
- 📦 **Nested Routing** - Perfect for admin panel sections
- 🎯 **Data Loading** - Loaders for server-side data fetching
- ✍️ **Mutations** - Actions for server-side form handling
- 🔒 **Security** - CSRF protection, secure cookies out of the box
- ⚡ **Performance** - Automatic code splitting, optimized bundles
- 🛠️ **Developer Experience** - TypeScript, Vite integration, HMR

**Refine.dev Benefits:**
- 🏗️ **Enterprise Admin Scaffolding** - Pre-built CRUD operations
- 📊 **Data Provider Pattern** - Clean separation of data layer
- 🎨 **UI Framework Agnostic** - Works with any React UI library
- 🔐 **Authentication Built-in** - Auth provider pattern
- 📝 **Form Management** - Integrated with React Hook Form
- 🔍 **Advanced Filtering** - Built-in filter, sort, pagination
- 📱 **Responsive** - Mobile-first design
- 🌍 **i18n Ready** - Multi-language support out of the box
- ♿ **Accessibility** - WCAG compliant components

### Core Features

#### 1. Configuration Management

- **Hierarchical Configuration**: Environment → Service → Component
- **Type-Safe Schema**: JSON Schema validation for configuration values
- **Versioning**: Full history with diff view and rollback capability
- **Multi-Environment**: Dev, staging, production with inheritance
- **Hot Reload**: Services can watch for configuration changes
- **Import/Export**: YAML/JSON bulk operations

#### 2. Secrets Management

- **Encryption at Rest**: AES-256-GCM encryption for sensitive values
- **Access Control**: RBAC with service-specific permissions
- **Secret Rotation**: Support for credential rotation workflows
- **Audit Trail**: Complete logging of secret access and modifications
- **Vault Integration Ready**: Future integration with HashiCorp Vault

#### 3. Admin Console Features

**Configuration Panel:**
- Browse and edit configurations by environment/service (Refine CRUD)
- Visual diff for version comparison
- Bulk import/export (YAML/JSON)
- Configuration validation before save
- Live preview of configuration changes

**Health Monitoring:**
- Real-time service health status dashboard
- Integration with Prometheus for metrics
- **Sentry error tracking dashboard** (embedded)
- Alerting rules configuration
- Dependency graph visualization

**Notifications Management:**
- Configure notification channels (Email, SMS, Slack, etc.)
- Template management with live preview
- User preference configuration
- Notification history and analytics
- Rate limit configuration

**Deployment Management:**
- Service deployment status
- Rolling update controls
- Rollback capabilities
- Deployment history and logs
- CI/CD integration

**Observability Integration:**
- Embedded Grafana dashboards
- Link to Jaeger traces
- Redpanda topic monitoring
- **Sentry error dashboard** (issues, releases, performance)
- Database connection pool status

**Audit & Security:**
- Configuration change history
- User activity logs
- Secret access logs
- Export audit reports (CSV, JSON)
- Compliance reporting

---

## Architecture Design

### Component Architecture

```
dictamesh/
├── services/
│   └── admin-console/          # NEW: Admin Console Service
│       ├── api/                # Backend API (Go)
│       │   ├── handlers/       # HTTP handlers
│       │   ├── middleware/     # Auth, CORS, logging
│       │   ├── models/         # Request/response models
│       │   ├── services/       # Business logic
│       │   └── main.go         # API server entry point
│       │
│       ├── app/                # Remix Application
│       │   ├── routes/         # Remix routes (file-based routing)
│       │   │   ├── _index.tsx                 # Dashboard
│       │   │   ├── _auth.tsx                  # Auth layout
│       │   │   ├── _auth.login.tsx            # Login page
│       │   │   ├── _dashboard.tsx             # Dashboard layout
│       │   │   ├── _dashboard.configs.tsx     # Config list
│       │   │   ├── _dashboard.configs.$id.tsx # Config detail
│       │   │   ├── _dashboard.health.tsx      # Health monitoring
│       │   │   ├── _dashboard.notifications.tsx  # Notifications
│       │   │   ├── _dashboard.sentry.tsx      # Sentry dashboard
│       │   │   ├── _dashboard.deployments.tsx # Deployments
│       │   │   └── _dashboard.audit.tsx       # Audit logs
│       │   │
│       │   ├── components/     # React components
│       │   │   ├── ui/         # Base UI components
│       │   │   ├── layout/     # Layout components
│       │   │   ├── config/     # Config-specific components
│       │   │   ├── charts/     # Data visualization
│       │   │   └── forms/      # Form components
│       │   │
│       │   ├── lib/            # Utilities and helpers
│       │   │   ├── api.ts      # API client
│       │   │   ├── auth.ts     # Auth utilities
│       │   │   ├── providers/  # Refine providers
│       │   │   │   ├── data-provider.ts     # Data provider
│       │   │   │   ├── auth-provider.ts     # Auth provider
│       │   │   │   └── access-control.ts    # Access control
│       │   │   └── utils.ts    # Misc utilities
│       │   │
│       │   ├── styles/         # CSS/Tailwind
│       │   │   └── app.css     # Global styles
│       │   │
│       │   ├── root.tsx        # Remix root component
│       │   └── entry.server.tsx # Server entry point
│       │
│       ├── public/             # Static assets
│       ├── package.json
│       ├── remix.config.js
│       ├── vite.config.ts
│       ├── tsconfig.json
│       ├── tailwind.config.ts
│       ├── go.mod              # Go backend dependencies
│       ├── Dockerfile          # Production build
│       └── Dockerfile.dev      # Development build
│
├── pkg/
│   └── config/                 # NEW: Configuration Package
│       ├── client.go           # Config client for services
│       ├── manager.go          # Config management logic
│       ├── store.go            # Storage abstraction
│       ├── encryption.go       # Secrets encryption
│       ├── validator.go        # Schema validation
│       ├── watcher.go          # Hot reload support
│       ├── models.go           # Core data models
│       ├── go.mod
│       └── migrations/         # Database migrations
│           └── sql/
│               └── 000004_config_tables.up.sql
│
└── infrastructure/
    └── docker-compose/
        └── docker-compose.dev.yml  # UPDATED: Add admin-console service
```

### Data Model

#### Configuration Tables

**dictamesh_configurations**
```sql
CREATE TABLE dictamesh_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    environment VARCHAR(50) NOT NULL,     -- dev, staging, prod
    service VARCHAR(100) NOT NULL,        -- metadata-catalog, graphql-gateway, etc.
    component VARCHAR(100),               -- database, cache, notifications, etc.
    key VARCHAR(255) NOT NULL,            -- configuration key
    value JSONB NOT NULL,                 -- configuration value
    value_type VARCHAR(50) NOT NULL,      -- string, number, boolean, object, array
    is_secret BOOLEAN DEFAULT false,      -- whether value is encrypted
    schema JSONB,                         -- JSON schema for validation
    description TEXT,
    tags TEXT[],                          -- searchable tags
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by VARCHAR(255),
    UNIQUE(environment, service, component, key)
);

CREATE INDEX idx_dictamesh_configs_env_svc ON dictamesh_configurations(environment, service);
CREATE INDEX idx_dictamesh_configs_tags ON dictamesh_configurations USING GIN(tags);
```

**dictamesh_config_versions**
```sql
CREATE TABLE dictamesh_config_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_id UUID REFERENCES dictamesh_configurations(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    value JSONB NOT NULL,
    is_secret BOOLEAN DEFAULT false,
    change_description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by VARCHAR(255),
    UNIQUE(config_id, version)
);

CREATE INDEX idx_dictamesh_config_versions_config ON dictamesh_config_versions(config_id);
```

**dictamesh_config_audit_logs**
```sql
CREATE TABLE dictamesh_config_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_id UUID REFERENCES dictamesh_configurations(id),
    action VARCHAR(50) NOT NULL,          -- CREATE, UPDATE, DELETE, ACCESS
    actor VARCHAR(255) NOT NULL,          -- user/service that performed action
    actor_type VARCHAR(50),               -- USER, SERVICE, API_KEY
    ip_address INET,
    user_agent TEXT,
    changes JSONB,                        -- before/after for updates
    metadata JSONB,                       -- additional context
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_dictamesh_config_audit_timestamp ON dictamesh_config_audit_logs(timestamp DESC);
CREATE INDEX idx_dictamesh_config_audit_actor ON dictamesh_config_audit_logs(actor);
CREATE INDEX idx_dictamesh_config_audit_config ON dictamesh_config_audit_logs(config_id);
```

**dictamesh_encryption_keys**
```sql
CREATE TABLE dictamesh_encryption_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_name VARCHAR(100) UNIQUE NOT NULL,
    encrypted_key BYTEA NOT NULL,         -- Master key encrypted with KEK
    algorithm VARCHAR(50) NOT NULL,       -- AES-256-GCM
    rotation_date TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    rotated_at TIMESTAMPTZ
);
```

### API Design

#### RESTful API Endpoints

**Configuration Management**
```
GET    /api/v1/configs                           # List configurations (filtered)
GET    /api/v1/configs/:id                       # Get configuration by ID
POST   /api/v1/configs                           # Create configuration
PUT    /api/v1/configs/:id                       # Update configuration
PATCH  /api/v1/configs/:id                       # Partial update
DELETE /api/v1/configs/:id                       # Delete configuration
GET    /api/v1/configs/:id/versions              # Get version history
POST   /api/v1/configs/:id/rollback/:version     # Rollback to version
GET    /api/v1/configs/export                    # Export configurations
POST   /api/v1/configs/import                    # Import configurations
POST   /api/v1/configs/validate                  # Validate configuration
```

**Environment/Service Queries**
```
GET    /api/v1/environments                      # List environments
GET    /api/v1/environments/:env/services        # Services in environment
GET    /api/v1/environments/:env/services/:svc/configs  # Service configs
GET    /api/v1/services                          # List all services
GET    /api/v1/services/:service/health          # Service health
```

**Health & Monitoring**
```
GET    /api/v1/health                            # API health check
GET    /api/v1/health/services                   # All services health
GET    /api/v1/health/services/:service          # Specific service health
GET    /api/v1/metrics                           # Prometheus metrics endpoint
GET    /api/v1/sentry/issues                     # Proxy to Sentry API
GET    /api/v1/sentry/stats                      # Sentry statistics
```

**Notifications Management**
```
GET    /api/v1/notifications/channels            # List notification channels
POST   /api/v1/notifications/channels            # Create channel
GET    /api/v1/notifications/templates           # List templates
POST   /api/v1/notifications/templates           # Create template
GET    /api/v1/notifications/history             # Notification history
```

**Deployment Management**
```
GET    /api/v1/deployments                       # List deployments
GET    /api/v1/deployments/:id                   # Deployment details
POST   /api/v1/deployments/:service              # Trigger deployment
POST   /api/v1/deployments/:id/rollback          # Rollback deployment
GET    /api/v1/deployments/:id/logs              # Deployment logs
```

**Audit Logs**
```
GET    /api/v1/audit                             # Query audit logs
GET    /api/v1/audit/export                      # Export audit logs
GET    /api/v1/audit/stats                       # Audit statistics
```

**Authentication & Authorization**
```
POST   /api/v1/auth/login                        # Login
POST   /api/v1/auth/logout                       # Logout
GET    /api/v1/auth/whoami                       # Current user info
POST   /api/v1/auth/refresh                      # Refresh token
GET    /api/v1/auth/permissions                  # User permissions
```

---

## Technology Stack

### Backend (Configuration Management API)

**Language & Framework:**
- **Go 1.21+** - Consistent with existing codebase
- **Chi Router** - Lightweight, idiomatic HTTP router
- **GORM** - Database ORM (already in use)
- **golang-migrate** - Database migrations (already in use)

**Dependencies:**
```go
// Core HTTP
github.com/go-chi/chi/v5
github.com/go-chi/cors
github.com/go-chi/httprate

// Authentication & Security
github.com/golang-jwt/jwt/v5
golang.org/x/crypto

// Configuration & Secrets
github.com/google/tink/go          // Encryption library
github.com/tidwall/gjson           // JSON query
github.com/xeipuuv/gojsonschema    // JSON schema validation

// Observability (already in use)
go.uber.org/zap
go.opentelemetry.io/otel
github.com/prometheus/client_golang
github.com/getsentry/sentry-go     // Sentry SDK
```

### Frontend (Admin Console - Remix + Refine.dev)

**Core Framework:**
- **Remix 2.x** - Full-stack React framework with SSR
- **React 18** - UI library
- **TypeScript 5** - Type safety
- **Vite** - Build tool with HMR

**Admin Panel Framework:**
- **Refine.dev 4.x** - Headless admin panel framework

**Dependencies:**
```json
{
  "dependencies": {
    // Core Remix
    "@remix-run/node": "^2.4.0",
    "@remix-run/react": "^2.4.0",
    "@remix-run/serve": "^2.4.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",

    // Refine Core
    "@refinedev/core": "^4.47.0",
    "@refinedev/remix-router": "^3.0.0",
    "@refinedev/react-hook-form": "^4.8.0",
    "@refinedev/react-table": "^5.6.0",

    // UI Components (using Mantine with Refine)
    "@refinedev/mantine": "^2.28.0",
    "@mantine/core": "^7.3.0",
    "@mantine/hooks": "^7.3.0",
    "@mantine/notifications": "^7.3.0",
    "@mantine/dates": "^7.3.0",
    "@tabler/icons-react": "^2.44.0",

    // Alternative: shadcn/ui (if preferred)
    // "@refinedev/cli": "^2.16.0",
    // "tailwindcss": "^3.4.0",
    // "@radix-ui/react-*": "latest",

    // Data Visualization
    "recharts": "^2.10.3",
    "@tremor/react": "^3.14.0",

    // Forms & Validation
    "react-hook-form": "^7.48.2",
    "zod": "^3.22.4",
    "@hookform/resolvers": "^3.3.2",

    // HTTP Client
    "axios": "^1.6.2",

    // Code/Diff Display
    "react-diff-viewer-continued": "^3.3.1",
    "react-syntax-highlighter": "^15.5.0",
    "js-yaml": "^4.1.0",

    // Utilities
    "dayjs": "^1.11.10",
    "lodash": "^4.17.21",
    "clsx": "^2.0.0"
  },
  "devDependencies": {
    "@remix-run/dev": "^2.4.0",
    "@types/react": "^18.2.42",
    "@types/react-dom": "^18.2.17",
    "typescript": "^5.3.3",
    "vite": "^5.0.7",
    "vite-tsconfig-paths": "^4.2.2",
    "eslint": "^8.55.0",
    "prettier": "^3.1.1"
  }
}
```

### Database

**PostgreSQL 16** - Already in use with extensions:
- `pgcrypto` - Encryption functions
- `uuid-ossp` - UUID generation (already enabled)
- Integration with existing tables: `dictamesh_notifications`, etc.

### Authentication

**Session-Based Authentication with JWT:**
- Server-side sessions stored in PostgreSQL or Redis
- JWT access tokens (short-lived, 15 minutes)
- Refresh tokens (long-lived, 7 days)
- Remix session management with secure HTTP-only cookies
- RBAC with roles: `admin`, `operator`, `viewer`, `readonly`

### Deployment

**Development:**
- Docker Compose with hot reload
- Remix dev server with Vite HMR on port 5173
- Go API server on port 8081
- Proxy configuration in Remix for API calls

**Production:**
- Multi-stage Docker build
- Remix production build served by Node.js or Go server
- Static assets on CDN (optional)
- Kubernetes deployment manifests

---

## Implementation Plan

### Phase 1: Foundation & Backend API (Week 1)

**Database Layer**
- [x] Review existing migrations (000001-000003 already exist)
- [ ] Create migration `000004_config_tables.up.sql`
- [ ] Implement database models in `pkg/config/models.go`
- [ ] Create GORM models with `TableName()` overrides
- [ ] Test migration with existing database

**Configuration Package**
- [ ] Initialize `pkg/config` Go module
- [ ] Implement `pkg/config/store.go` (CRUD operations)
- [ ] Implement `pkg/config/encryption.go` (secrets encryption)
- [ ] Implement `pkg/config/validator.go` (JSON schema validation)
- [ ] Implement `pkg/config/watcher.go` (hot reload)
- [ ] Write unit tests for core functionality

**API Service Setup**
- [ ] Initialize `services/admin-console/api/` module
- [ ] Set up Chi router with middleware
- [ ] Implement authentication middleware (JWT + sessions)
- [ ] Implement CORS and rate limiting
- [ ] Add Sentry error tracking integration

**API Handlers**
- [ ] Configuration CRUD endpoints
- [ ] Version management endpoints
- [ ] Health check endpoints
- [ ] Audit log query endpoints
- [ ] Sentry proxy endpoints
- [ ] Notifications configuration endpoints

### Phase 2: Remix + Refine.dev Setup (Week 2)

**Project Initialization**
- [ ] Initialize Remix app: `npx create-remix@latest`
- [ ] Install Refine dependencies
- [ ] Configure TypeScript, Vite, Tailwind
- [ ] Set up Refine providers (data, auth, access control)
- [ ] Configure Mantine UI (or shadcn/ui alternative)

**Authentication & Layout**
- [ ] Create auth routes (`_auth.login.tsx`, `_auth.logout.tsx`)
- [ ] Implement Remix session management
- [ ] Create dashboard layout (`_dashboard.tsx`)
- [ ] Implement navigation sidebar
- [ ] Add user profile menu

**Data Provider Implementation**
- [ ] Create Refine data provider for Go API
- [ ] Implement CRUD operations mapping
- [ ] Add error handling and validation
- [ ] Implement pagination, filtering, sorting
- [ ] Add optimistic updates

**Auth Provider Implementation**
- [ ] Create Refine auth provider
- [ ] Implement login/logout flows
- [ ] Add permission checks
- [ ] Implement role-based access control

### Phase 3: Core Admin Pages (Week 2-3)

**Dashboard Page**
- [ ] Service health overview cards
- [ ] Recent configuration changes
- [ ] Error summary (Sentry integration)
- [ ] Notification statistics
- [ ] System metrics charts

**Configuration Management**
- [ ] Config list page with Refine Table
  - [ ] Environment/service filters
  - [ ] Search functionality
  - [ ] Bulk actions (export, delete)
- [ ] Config detail/edit page
  - [ ] Form with validation (React Hook Form + Zod)
  - [ ] JSON/YAML editor with syntax highlighting
  - [ ] Secret masking for sensitive values
- [ ] Version history page
  - [ ] Timeline view of changes
  - [ ] Diff viewer for comparing versions
  - [ ] Rollback functionality

**Health Monitoring**
- [ ] Service health dashboard
- [ ] Integration with Prometheus metrics
- [ ] Real-time updates with polling or WebSockets
- [ ] Alerting configuration

**Notifications Management**
- [ ] Channel configuration pages
- [ ] Template editor with live preview
- [ ] Notification history table
- [ ] Test notification functionality

**Sentry Integration**
- [ ] Error dashboard (issues summary)
- [ ] Recent errors table
- [ ] Error detail view (redirect to Sentry)
- [ ] Release tracking

**Audit Logs**
- [ ] Audit log table with advanced filtering
- [ ] Export functionality (CSV, JSON)
- [ ] Log detail view
- [ ] Search and date range filtering

### Phase 4: Integration & Testing (Week 3-4)

**Docker Compose Integration**
- [ ] Add admin-console service to `docker-compose.dev.yml`
- [ ] Configure networking between services
- [ ] Add health checks
- [ ] Volume mounts for development

**Hot Reload Development**
- [ ] Configure Vite HMR for Remix frontend
- [ ] Set up Air for Go hot reload (optional)
- [ ] Test full-stack hot reload workflow

**Testing**
- [ ] Unit tests for Go API handlers
- [ ] Integration tests for database operations
- [ ] Remix loader/action tests
- [ ] E2E tests for critical flows (Playwright)
- [ ] Component tests (Vitest + React Testing Library)

**Documentation**
- [ ] API documentation (OpenAPI/Swagger)
- [ ] User guide for admin console
- [ ] Developer guide for configuration package
- [ ] Deployment guide

### Phase 5: Production & Polish (Week 4)

**Security Hardening**
- [ ] Security audit of authentication flow
- [ ] Implement rate limiting per user
- [ ] Add CSP headers
- [ ] Secrets rotation mechanism
- [ ] CSRF protection (Remix built-in)
- [ ] Input sanitization

**Production Deployment**
- [ ] Multi-stage Dockerfile optimization
- [ ] Kubernetes manifests (Deployment, Service, Ingress)
- [ ] Helm chart creation (optional)
- [ ] CI/CD pipeline setup
- [ ] Environment variable configuration

**Performance Optimization**
- [ ] React component lazy loading
- [ ] API response caching
- [ ] Database query optimization
- [ ] Bundle size analysis and optimization
- [ ] Lighthouse performance audit

**Monitoring & Alerting**
- [ ] Grafana dashboards for admin console
- [ ] Prometheus alerts for critical metrics
- [ ] Sentry performance monitoring
- [ ] Custom metrics for configuration changes

**UI/UX Polish**
- [ ] Responsive design testing
- [ ] Accessibility audit (WCAG)
- [ ] Loading states and skeletons
- [ ] Error boundaries
- [ ] Toast notifications
- [ ] Keyboard shortcuts

---

## Security Considerations

### Encryption Strategy

**Envelope Encryption Pattern:**
```
Master Key (KEK - Key Encryption Key)
    └─> Stored in environment variable or KMS (AWS KMS, GCP KMS, HashiCorp Vault)
    └─> Encrypts Data Encryption Keys (DEK)

Data Encryption Keys (DEK)
    └─> Stored in dictamesh_encryption_keys table (encrypted)
    └─> One per environment or rotation period
    └─> Used to encrypt configuration secrets

Configuration Secrets
    └─> Encrypted with active DEK using AES-256-GCM
    └─> Stored in dictamesh_configurations.value (when is_secret=true)
```

**Encryption Implementation:**
```go
// Use Google Tink for encryption
import "github.com/google/tink/go/aead"

type EncryptionService struct {
    masterKey []byte
    aead      tink.AEAD
}

// Encrypt secret with envelope encryption
func (s *EncryptionService) EncryptSecret(plaintext string) ([]byte, error)

// Decrypt secret with envelope encryption
func (s *EncryptionService) DecryptSecret(ciphertext []byte) (string, error)

// Rotate encryption keys (for key rotation workflow)
func (s *EncryptionService) RotateKeys() error
```

### Authentication & Authorization

**Remix Session-Based Authentication:**
```typescript
// app/lib/auth.ts
import { createCookieSessionStorage } from "@remix-run/node";

export const sessionStorage = createCookieSessionStorage({
  cookie: {
    name: "__session",
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    secrets: [process.env.SESSION_SECRET!],
    sameSite: "lax",
    maxAge: 60 * 60 * 24 * 7, // 7 days
  },
});

export async function requireUser(request: Request) {
  const session = await sessionStorage.getSession(
    request.headers.get("Cookie")
  );
  const userId = session.get("userId");
  if (!userId) {
    throw redirect("/login");
  }
  return userId;
}
```

**RBAC Roles:**
| Role | Permissions |
|------|-------------|
| `admin` | Full access to all operations |
| `operator` | Read/write configs, read-only secrets, view audit logs |
| `viewer` | Read-only access to configs and health |
| `readonly` | View-only access to dashboards |

**Access Control with Refine:**
```typescript
// app/lib/providers/access-control.ts
import { AccessControlProvider } from "@refinedev/core";

export const accessControlProvider: AccessControlProvider = {
  can: async ({ resource, action, params }) => {
    const user = await getCurrentUser();

    // Admin can do everything
    if (user.role === "admin") {
      return { can: true };
    }

    // Operators can't delete or modify secrets
    if (user.role === "operator") {
      if (action === "delete" || (resource === "secrets" && action !== "list")) {
        return { can: false, reason: "Insufficient permissions" };
      }
      return { can: true };
    }

    // Viewers can only read
    if (user.role === "viewer") {
      return { can: action === "list" || action === "show" };
    }

    return { can: false };
  },
};
```

### Audit Logging

**All sensitive operations are logged:**
- Configuration create/update/delete
- Secret access (read)
- Authentication events (login, logout, failed attempts)
- Authorization failures
- Key rotation events
- Deployment actions

**Audit log format:**
```json
{
  "id": "uuid",
  "action": "UPDATE_CONFIG",
  "actor": "admin@example.com",
  "actor_type": "USER",
  "resource": "/api/v1/configs/12345",
  "ip_address": "10.0.1.5",
  "user_agent": "Mozilla/5.0...",
  "changes": {
    "before": {"database.max_connections": 100},
    "after": {"database.max_connections": 200}
  },
  "timestamp": "2025-11-08T10:30:00Z"
}
```

### Network Security

**Development:**
- Services communicate within Docker network
- Admin console exposed on localhost only
- HTTPS via Caddy (optional for development)

**Production:**
- API behind Ingress with TLS termination (Let's Encrypt)
- mTLS between services (optional)
- Network policies to restrict traffic
- CORS configured for allowed origins only
- CSP headers configured

---

## Deployment Strategy

### Development Environment

**Docker Compose Service:**
```yaml
services:
  admin-console:
    build:
      context: ../../services/admin-console
      dockerfile: Dockerfile.dev
    container_name: dictamesh-admin-console
    ports:
      - "8081:8081"     # Go API
      - "5173:5173"     # Remix dev server
    environment:
      # Database
      - DATABASE_URL=postgres://dictamesh:dictamesh_dev_password@postgres:5432/metadata_catalog

      # External services
      - REDIS_URL=redis://redis:6379
      - PROMETHEUS_URL=http://prometheus:9090
      - JAEGER_URL=http://jaeger:14268
      - SENTRY_DSN=http://admin-console-key@sentry:9000/2
      - SENTRY_URL=http://sentry:9000

      # Authentication
      - JWT_SECRET=dev-secret-change-in-production
      - SESSION_SECRET=dev-session-secret-32-chars-min
      - MASTER_ENCRYPTION_KEY=dev-key-32-bytes-change-prod

      # App config
      - NODE_ENV=development
      - LOG_LEVEL=debug
      - API_BASE_URL=http://localhost:8081
      - PUBLIC_URL=http://localhost:5173

    volumes:
      - ../../services/admin-console:/app
      - /app/node_modules    # Anonymous volume for node_modules
      - /app/api/tmp         # Go build artifacts
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      sentry:
        condition: service_healthy
    networks:
      - dictamesh-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/api/v1/health"]
      interval: 10s
      timeout: 5s
      retries: 3
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '1'
```

**Development Workflow:**
```bash
# Start all services including admin console
cd infrastructure
make dev-up

# Initialize Sentry (first time only)
make sentry-init

# Access admin console (Remix dev server with HMR)
open http://localhost:5173

# Access API directly
curl http://localhost:8081/api/v1/health

# View logs
make logs service=admin-console

# Restart specific service
docker-compose restart admin-console

# Stop all services
make dev-down
```

### Production Kubernetes Deployment

**Deployment Manifest:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: admin-console
  namespace: dictamesh
  labels:
    app: admin-console
    component: admin
spec:
  replicas: 2
  selector:
    matchLabels:
      app: admin-console
  template:
    metadata:
      labels:
        app: admin-console
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8081"
        prometheus.io/path: "/api/v1/metrics"
    spec:
      containers:
      - name: admin-console
        image: dictamesh/admin-console:latest
        ports:
        - containerPort: 8081
          name: api
          protocol: TCP
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: dictamesh-db-credentials
              key: url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: admin-console-secrets
              key: jwt-secret
        - name: SESSION_SECRET
          valueFrom:
            secretKeyRef:
              name: admin-console-secrets
              key: session-secret
        - name: MASTER_ENCRYPTION_KEY
          valueFrom:
            secretKeyRef:
              name: admin-console-secrets
              key: master-key
        - name: SENTRY_DSN
          valueFrom:
            secretKeyRef:
              name: sentry-credentials
              key: admin-console-dsn
        - name: NODE_ENV
          value: production
        - name: LOG_LEVEL
          value: info
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8081
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            cpu: 100m
            memory: 256Mi
          limits:
            cpu: 1000m
            memory: 1Gi
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
```

**Service & Ingress:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: admin-console
  namespace: dictamesh
spec:
  selector:
    app: admin-console
  ports:
  - port: 8081
    targetPort: 8081
    name: api
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: admin-console
  namespace: dictamesh
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - admin.dictamesh.example.com
    secretName: admin-console-tls
  rules:
  - host: admin.dictamesh.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: admin-console
            port:
              number: 8081
```

---

## Development Workflow

### Local Development Setup

**Prerequisites:**
```bash
# Required
- Go 1.21+
- Node.js 20+
- Docker 24.0+
- Docker Compose 2.0+
- Make

# Optional (for native development)
- Air (Go hot reload): go install github.com/cosmtrek/air@latest
```

**Setup Steps:**
```bash
# 1. Clone and navigate to project
cd /home/user/dictamesh

# 2. Ensure on correct branch
git checkout claude/centralized-config-admin-console-011CUvx4PeRVkHGQ7qDgnTto

# 3. Start infrastructure services
cd infrastructure
make dev-up

# 4. Initialize Sentry (first time only)
make sentry-init

# 5. Run database migrations
cd ../pkg/config
go run ../database/migrations/migrate.go up

# 6. Start Admin Console via Docker Compose (recommended)
cd ../../infrastructure
docker-compose up admin-console

# OR start manually for development:

# 6a. Start API server with hot reload (terminal 1)
cd ../services/admin-console/api
air

# 6b. Start Remix dev server (terminal 2)
cd ../app
npm install
npm run dev

# 7. Access admin console
# Remix UI (with HMR): http://localhost:5173
# API direct: http://localhost:8081
# Sentry: http://localhost:9000
```

### Hot Module Reload (HMR) Configuration

**Remix Vite Configuration:**
```typescript
// services/admin-console/app/vite.config.ts
import { vitePlugin as remix } from "@remix-run/dev";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [
    remix({
      ignoredRouteFiles: ["**/.*"],
    }),
    tsconfigPaths(),
  ],
  server: {
    port: 5173,
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: "http://localhost:8081",
        changeOrigin: true,
      },
    },
    watch: {
      usePolling: true, // For Docker volumes
    },
  },
});
```

**Remix Configuration:**
```javascript
// services/admin-console/app/remix.config.js
/** @type {import('@remix-run/dev').AppConfig} */
export default {
  ignoredRouteFiles: ["**/.*"],
  serverModuleFormat: "esm",
  tailwind: true,
  future: {
    v3_fetcherPersist: true,
    v3_relativeSplatPath: true,
    v3_throwAbortReason: true,
  },
};
```

**Backend Air Configuration:**
```toml
# services/admin-console/api/.air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/admin-console-api"
  cmd = "go build -o ./tmp/admin-console-api ./main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

### Configuration Client Usage

**For other DictaMesh services to consume configurations:**

```go
// Example: Metadata Catalog Service
package main

import (
    "context"
    "log"

    "github.com/Click2-Run/dictamesh/pkg/config"
)

func main() {
    // Initialize config client
    client, err := config.NewClient(config.ClientOptions{
        APIEndpoint: "http://admin-console:8081",
        Environment: "production",
        Service:     "metadata-catalog",
        APIKey:      os.Getenv("CONFIG_API_KEY"),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Get configuration
    dbMaxConns, err := client.GetInt("database.max_connections")
    if err != nil {
        log.Fatal(err)
    }

    // Get secret (automatically decrypted)
    dbPassword, err := client.GetSecret("database.password")
    if err != nil {
        log.Fatal(err)
    }

    // Get complex configuration
    var notifConfig NotificationConfig
    err = client.GetStruct("notifications", &notifConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Watch for configuration changes (hot reload)
    client.Watch("database.*", func(key string, newValue interface{}) {
        log.Printf("Configuration changed: %s = %v", key, newValue)
        // Reload database connection pool, etc.
        reloadDatabasePool(newValue)
    })

    // Start service...
}
```

### Refine Data Provider Example

**Remix Loader with Refine:**
```typescript
// app/routes/_dashboard.configs.tsx
import { json, type LoaderFunctionArgs } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { List, useTable } from "@refinedev/core";
import { Table } from "@refinedev/mantine";

export async function loader({ request }: LoaderFunctionArgs) {
  const url = new URL(request.url);
  const page = parseInt(url.searchParams.get("page") || "1");
  const pageSize = parseInt(url.searchParams.get("pageSize") || "10");
  const environment = url.searchParams.get("environment") || "";

  // Fetch from Go API
  const response = await fetch(`http://localhost:8081/api/v1/configs?page=${page}&pageSize=${pageSize}&environment=${environment}`);
  const data = await response.json();

  return json(data);
}

export default function ConfigList() {
  const initialData = useLoaderData<typeof loader>();

  const { tableQueryResult } = useTable({
    resource: "configs",
    initialData,
  });

  return (
    <List>
      <Table {...tableQueryResult}>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Environment</Table.Th>
            <Table.Th>Service</Table.Th>
            <Table.Th>Key</Table.Th>
            <Table.Th>Value</Table.Th>
            <Table.Th>Actions</Table.Th>
          </Table.Tr>
        </Table.Thead>
        {/* Table body... */}
      </Table>
    </List>
  );
}
```

---

## Next Steps

### Immediate Actions

1. **Review & Approval**: Review this design document with stakeholders
2. **Technology Validation**: Create POC with Remix + Refine.dev integration
3. **Create Work Items**: Break down implementation into GitHub issues
4. **Set Up Project Board**: Track progress of implementation phases

### Implementation Order

**Priority 1 (Critical Path - Week 1):**
1. Database schema and migrations (`000004_config_tables.up.sql`)
2. Configuration package (`pkg/config`)
3. Basic Go API with health endpoints
4. Authentication middleware

**Priority 2 (Core Features - Week 2):**
1. Configuration CRUD operations (API)
2. Secrets encryption implementation
3. Remix + Refine.dev project setup
4. Authentication UI (login/logout)
5. Basic dashboard layout

**Priority 3 (Admin Features - Week 3):**
1. Configuration management pages (CRUD with Refine)
2. Version history and rollback
3. Health monitoring integration
4. Notifications management UI
5. Sentry dashboard integration

**Priority 4 (Polish & Production - Week 4):**
1. Audit log viewer
2. Advanced UI features (diff viewer, bulk operations)
3. Performance optimization
4. Comprehensive testing
5. Documentation
6. Production deployment setup

### Success Metrics

- ✅ All framework configurations centralized (0 `.env` files in production)
- ✅ 100% of secrets encrypted at rest
- ✅ Complete audit trail for configuration changes
- ✅ < 100ms API response time (p95)
- ✅ < 2s page load time for admin console (p95)
- ✅ Zero downtime configuration updates
- ✅ 95%+ test coverage for critical paths
- ✅ WCAG 2.1 Level AA accessibility compliance

---

## Appendix

### Why Remix + Refine.dev vs Alternatives?

**Compared to React SPA + Vite:**
- ✅ Remix: Better SEO, faster initial load, progressive enhancement
- ✅ Remix: Built-in form handling, no need for client-side form libraries
- ✅ Remix: Server-side data loading prevents waterfall requests
- ✅ Refine: Pre-built admin patterns, faster development

**Compared to Next.js:**
- ✅ Remix: Simpler mental model (no getServerSideProps/getStaticProps confusion)
- ✅ Remix: Better nested routing for admin panels
- ✅ Remix: More control over server-side logic
- ✅ Remix: Lighter bundle size

**Compared to Admin Panel Templates (React Admin, Django Admin):**
- ✅ Refine: Modern, type-safe, React-based
- ✅ Refine: Framework agnostic (works with Remix, Next.js, etc.)
- ✅ Refine: Better customization and extensibility
- ✅ Refine: Active development and community

### Alternative UI Component Libraries

**Option 1: Mantine (Recommended)**
- Modern, fully-featured React component library
- Excellent TypeScript support
- Built-in dark mode
- Official Refine integration: `@refinedev/mantine`
- Over 100 components

**Option 2: shadcn/ui + Tailwind**
- Copy-paste components (more control)
- Built on Radix UI (accessible)
- Tailwind CSS styling
- Requires more manual setup with Refine

**Option 3: Ant Design**
- Mature, enterprise-grade components
- Official Refine integration: `@refinedev/antd`
- Heavy bundle size (not recommended for this use case)

**Decision: Use Mantine** for better balance of features, performance, and DX.

### References

- [Remix Documentation](https://remix.run/docs)
- [Refine Documentation](https://refine.dev/docs)
- [Mantine UI](https://mantine.dev/)
- [12-Factor App: Config](https://12factor.net/config)
- [OWASP Secrets Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
- [Google Tink Cryptography](https://github.com/google/tink)
- [PostgreSQL Encryption Functions](https://www.postgresql.org/docs/current/pgcrypto.html)

---

**End of Design Document**
