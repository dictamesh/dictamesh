# DictaMesh

**Enterprise-Grade Data Mesh Adapter Framework**

## What is DictaMesh?

DictaMesh is a comprehensive **framework and foundation** for building data mesh adapters. It provides the core infrastructure to integrate any type of data source (APIs, SDKs, databases, file systems) into a unified, event-driven data mesh architecture.

### This is NOT

❌ A specific implementation for particular systems (like specific CMS, e-commerce platforms, etc.)
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

## Project Structure

DictaMesh is organized into four distinct layers:

### 1. **Core Framework** (This Repository)
The foundational infrastructure that everything builds upon:
- Adapter Interface (DataProductAdapter contract)
- Event Bus (Kafka integration, event schemas)
- Metadata Catalog (entity registry, relationships, lineage)
- GraphQL Gateway (Apollo Federation)
- Observability (OpenTelemetry, Prometheus, logging)
- Governance (access control, PII tracking, compliance)
- Resilience patterns (circuit breakers, retries, rate limiting)
- Testing framework and deployment templates

### 2. **Connectors** (Data Source Drivers)
Low-level drivers for different protocols and technologies:
- Database connectors (PostgreSQL, MySQL, MongoDB, Oracle)
- API connectors (REST, GraphQL, gRPC, SOAP, **OpenAPI 3.0**)
- File system connectors (CSV, JSON, XML, Parquet)
- Message queue connectors (RabbitMQ, Redis, SQS)
- Legacy system connectors (ODBC, JDBC, FTP)

Repository: `github.com/click2-run/dictamesh-connectors/*`

### 3. **Adapters** (Data Integration Layer)
Domain-specific implementations using connectors:
- Implement DataProductAdapter interface
- Use connectors to access data sources
- Transform data to canonical models
- Handle business logic and validation
- Publish events and manage metadata

Repository: `github.com/click2-run/dictamesh-adapters/*`

### 4. **Services** (Business Logic & Applications)
Higher-level applications consuming the data mesh:
- API services (REST/GraphQL endpoints)
- Data pipelines (ETL, synchronization)
- Workflow engines (business processes)
- AI/ML services (recommendations, analytics)
- Agents and automation

Repository: Your own or `github.com/click2-run/dictamesh-services/*`

### What You Build

Using DictaMesh, you typically build:
- **Connectors** (if you need a new protocol/technology)
- **Adapters** (for your specific data sources)
- **Services** (your business applications)
- GraphQL schemas for your domain entities

## Example Usage

```go
// 1. Use a connector to access your data source
import "github.com/click2-run/dictamesh-connectors/rest"

restConn := rest.NewConnector(rest.Config{
    BaseURL: "https://api.example.com",
    Auth:    rest.BearerToken("your-token"),
})

// 2. Build an adapter using the connector
type MyCustomAdapter struct {
    connector *rest.Connector
    // Framework injects: cache, events, metrics, etc.
}

func (a *MyCustomAdapter) GetEntity(ctx context.Context, id string) (*Entity, error) {
    // Use connector to fetch data
    data, err := a.connector.Get(ctx, "/customers/"+id)
    // Transform and return
    return transformToEntity(data), err
}

// 3. Register with framework - it handles the rest
app.RegisterAdapter("my_domain", myAdapter)

// Framework automatically provides:
// - GraphQL API, Event streaming, Metadata catalog
// - Observability, Governance, Resilience patterns
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
