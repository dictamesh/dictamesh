---
slug: welcome-to-dictamesh
title: Welcome to DictaMesh
authors: [dictamesh_team]
tags: [announcement, introduction, data-mesh]
---

<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

# Welcome to DictaMesh

We're excited to introduce **DictaMesh** - an enterprise-grade data mesh adapter framework that provides the foundational infrastructure for building federated data integrations.

<!--truncate-->

## What is DictaMesh?

DictaMesh is not a pre-built integration tool for specific platforms. Instead, it's a **comprehensive framework** that developers use to build their own data mesh integrations. Think of it as the foundation upon which you build custom adapters for any data source.

### Key Features

- **Event-Driven Architecture**: Built on Apache Kafka with structured event schemas
- **Federated GraphQL Gateway**: Apollo Federation for unified API across all data sources
- **Metadata Catalog**: Centralized registry with entity discovery and lineage tracking
- **Production-Ready Patterns**: Circuit breakers, retry logic, caching, and observability
- **Built-in Observability**: OpenTelemetry tracing, Prometheus metrics, structured logging

## Why Build DictaMesh?

In modern enterprises, data is scattered across dozens or hundreds of systems - CMS platforms, e-commerce systems, ERPs, CRMs, databases, and APIs. Integrating these systems traditionally requires:

1. **Custom point-to-point integrations** - Hard to maintain, brittle, doesn't scale
2. **Proprietary iPaaS solutions** - Expensive, vendor lock-in, limited customization
3. **Building everything from scratch** - Time-consuming, reinventing the wheel

DictaMesh provides a third way: a **framework with batteries included** that lets you build custom integrations while leveraging battle-tested patterns for resilience, observability, and governance.

## Who Should Use DictaMesh?

DictaMesh is ideal for:

- **Enterprises** building data mesh architectures
- **SaaS companies** needing to integrate with customer systems
- **System integrators** building integration solutions
- **Data engineering teams** centralizing data access
- **Development teams** needing unified APIs across microservices

## What Can You Build?

With DictaMesh, you can build adapters for:

- **Content Management Systems** (WordPress, Drupal, Contentful)
- **E-commerce Platforms** (Shopify, WooCommerce, Magento)
- **ERP Systems** (SAP, Oracle, Microsoft Dynamics)
- **CRM Systems** (Salesforce, HubSpot, Dynamics 365)
- **Databases** (PostgreSQL, MySQL, MongoDB, Oracle)
- **File Systems** (S3, Azure Blob, Google Cloud Storage)
- **Custom APIs** (REST, GraphQL, SOAP, gRPC)

Each adapter you build automatically gets:
- ✅ GraphQL API
- ✅ Event streaming
- ✅ Metadata registration
- ✅ Caching
- ✅ Distributed tracing
- ✅ Metrics collection
- ✅ Circuit breakers
- ✅ Retry logic

## Technology Stack

DictaMesh is built with modern, proven technologies:

- **Language**: Go 1.21+ (performance, concurrency, type safety)
- **Event Streaming**: Apache Kafka / Redpanda
- **Database**: PostgreSQL 15+ (JSONB, pgvector)
- **Cache**: Redis 7+
- **API**: GraphQL with Apollo Federation
- **Observability**: OpenTelemetry, Prometheus, Jaeger
- **Deployment**: Kubernetes with Helm

## Getting Started

Ready to build your first adapter?

1. **[Quick Start Guide](/docs/getting-started/quickstart)** - Get up and running in 10 minutes
2. **[Core Concepts](/docs/getting-started/core-concepts)** - Understand key concepts
3. **[Building Adapters](/docs/guides/building-adapters)** - Build your first adapter
4. **[Architecture Overview](/docs/architecture/overview)** - Deep dive into design

## Open Source & License

DictaMesh is licensed under **AGPL-3.0-or-later**, which means:

- ✅ **Free for commercial use**
- ✅ **Modify and distribute**
- ⚠️ **Network use requires source disclosure** (if you modify DictaMesh and run it as a service, you must make your source code available)

We chose AGPL to ensure that improvements to the framework benefit the entire community.

## Roadmap

We have ambitious plans for DictaMesh:

### Q4 2025
- ✅ Core framework (metadata catalog, event bus, GraphQL gateway)
- ✅ Observability stack integration
- ✅ Example adapters and documentation
- 🔄 Production deployments and battle-testing

### Q1 2026
- 🔲 Advanced governance features (data classification, lineage visualization)
- 🔲 Performance optimizations (query caching, connection pooling)
- 🔲 Additional connector types (gRPC, SOAP, legacy systems)
- 🔲 Visual schema designer

### Q2 2026
- 🔲 Data quality framework
- 🔲 Schema registry UI
- 🔲 Marketplace for community adapters
- 🔲 Advanced analytics and reporting

## Community

Join our growing community:

- 📖 **[Documentation](https://docs.dictamesh.com)**
- 💬 **[GitHub Discussions](https://github.com/dictamesh/dictamesh/discussions)**
- 🐛 **[Issue Tracker](https://github.com/dictamesh/dictamesh/issues)**
- 🤝 **[Contributing Guide](/docs/contributing)**

## What's Next?

Follow this blog for:
- Release announcements
- Technical deep dives
- Case studies
- Best practices
- Community highlights

We're excited to see what you'll build with DictaMesh!

---

**Questions?** Join our [GitHub Discussions](https://github.com/dictamesh/dictamesh/discussions) or open an [issue](https://github.com/dictamesh/dictamesh/issues).
