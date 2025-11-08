# DictaMesh

**Enterprise-Grade Data Mesh Adapter Framework**

## What is DictaMesh?

DictaMesh is a comprehensive **framework and foundation** for building data mesh adapters. It provides the core infrastructure to integrate any type of data source (APIs, SDKs, databases, file systems) into a unified, event-driven data mesh architecture.

### This is NOT

❌ A specific implementation for particular systems (like Directus, Shopify, etc.)
❌ A pre-built integration tool for specific platforms
❌ A ready-to-run product for a specific use case

### This IS

✅ A **framework** that developers use to build their own data mesh integrations
✅ A set of **core components** (metadata catalog, event bus, GraphQL gateway)
✅ **Standard interfaces** and patterns for building adapters
✅ **Built-in observability, governance, and resilience** features
✅ **Example reference implementations** showing how to use the framework

## Quick Start

See [PROJECT-SCOPE.md](PROJECT-SCOPE.md) for the complete architecture and usage guide.

### What You Get

The framework provides:
- Data Product Adapter Interface (standard contract)
- Event-driven integration layer (Kafka)
- Metadata catalog service
- Federated GraphQL gateway
- Observability stack (tracing, metrics, logging)
- Governance engine
- Resilience patterns (circuit breakers, retries)
- Testing framework
- Deployment templates (Kubernetes, Helm)

### What You Build

Using DictaMesh, you build:
- Your custom adapters for your data sources
- GraphQL schemas for your domain entities
- Business logic specific to your systems
- Configuration to connect your sources

## Example Usage

```go
// 1. Implement the DataProductAdapter interface for your system
type MyCustomAdapter struct {
    // Your implementation using the framework's base components
}

// 2. The framework provides the rest (events, catalog, API, observability)
```

## License

This project is licensed under the GNU Affero General Public License v3.0.
See LICENSE file for details.

Commercial use is permitted. If you modify this software and provide it
as a network service, you must make your source code available under AGPL v3.

## Documentation

- [PROJECT-SCOPE.md](PROJECT-SCOPE.md) - Complete architecture and technical details
- [AGENT.md](AGENT.md) - Development guidelines for contributors
- [CLAUDE.md](CLAUDE.md) - AI assistant instructions
