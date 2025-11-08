# DictaMesh Core SDK - Planning Document

**Version:** 1.0.0
**Status:** Planning Phase
**Created:** 2025-11-08

## Executive Summary

The DictaMesh Core SDK is a TypeScript/JavaScript library that provides a framework-agnostic, protocol-driven client for interacting with DictaMesh data mesh infrastructure. It enables developers to build data-driven applications across React, Vue.js, Angular, Node.js, and other JavaScript environments with a unified, type-safe API.

## Vision

Create a **core SDK** that:
- Is completely **agnostic** to specific adapters, connectors, and backend implementations
- Accepts **modules/plugins** for different transport protocols and data sources
- Provides an **abstract protocol layer** for communication
- Works seamlessly in **browser and Node.js** environments
- Supports **real-time queries**, **data subscriptions**, **CRUD operations**, and **search**
- Implements a **self-written protocol pattern** that can work with any integrated solution

## Core Principles

### 1. Abstraction First
- Core SDK defines interfaces and contracts
- Implementations are provided via adapter modules
- No hard dependencies on specific backends

### 2. Protocol-Driven
- Communication via well-defined protocol messages
- Support multiple transport mechanisms (HTTP, WebSocket, SSE, gRPC-Web)
- Protocol can be implemented by any backend

### 3. Plugin Architecture
- Adapters: Backend-specific implementations (GraphQL, REST, custom protocols)
- Connectors: Transport-specific implementations (HTTP, WebSocket)
- Middleware: Request/response transformation, caching, auth

### 4. Real-Time First
- Built-in support for subscriptions and live queries
- Optimistic updates and offline support
- Conflict resolution strategies

### 5. Type Safety
- Full TypeScript support with generics
- Runtime type validation
- Schema introspection and code generation

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│              (React, Vue, Node.js, etc.)                    │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                   DictaMesh Core SDK                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Public API Surface                      │   │
│  │  - Client, Query, Mutation, Subscription             │   │
│  └──────────────────────────────────────────────────────┘   │
│                             │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           Core Abstractions Layer                    │   │
│  │  - IProtocol, IAdapter, IConnector, ICache           │   │
│  └──────────────────────────────────────────────────────┘   │
│                             │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Protocol Layer                          │   │
│  │  - Request/Response Types                            │   │
│  │  - Message Serialization                             │   │
│  │  - Protocol Versioning                               │   │
│  └──────────────────────────────────────────────────────┘   │
│                             │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Plugin System                           │   │
│  │  - Adapter Registry                                  │   │
│  │  - Middleware Pipeline                               │   │
│  │  - Hook System                                       │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│   Adapters   │    │  Connectors  │    │  Middleware  │
│              │    │              │    │              │
│ - GraphQL    │    │ - HTTP       │    │ - Auth       │
│ - REST       │    │ - WebSocket  │    │ - Cache      │
│ - gRPC-Web   │    │ - SSE        │    │ - Retry      │
│ - Custom     │    │ - Custom     │    │ - Logging    │
└──────────────┘    └──────────────┘    └──────────────┘
        │                    │                    │
        └────────────────────┼────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend Services                         │
│     (DictaMesh Gateway, Adapters, Event Bus)                │
└─────────────────────────────────────────────────────────────┘
```

## Package Structure

```
packages/sdk-core/
├── src/
│   ├── core/
│   │   ├── client.ts                 # Main SDK client
│   │   ├── context.ts                # Request context
│   │   ├── error.ts                  # Error types
│   │   └── types.ts                  # Core types
│   │
│   ├── protocol/
│   │   ├── index.ts                  # Protocol interface
│   │   ├── request.ts                # Request types
│   │   ├── response.ts               # Response types
│   │   ├── message.ts                # Message format
│   │   ├── serializer.ts             # Serialization
│   │   └── version.ts                # Protocol versioning
│   │
│   ├── abstractions/
│   │   ├── adapter.ts                # IAdapter interface
│   │   ├── connector.ts              # IConnector interface
│   │   ├── cache.ts                  # ICache interface
│   │   ├── middleware.ts             # IMiddleware interface
│   │   └── transport.ts              # ITransport interface
│   │
│   ├── query/
│   │   ├── builder.ts                # Query builder
│   │   ├── filter.ts                 # Filter DSL
│   │   ├── sort.ts                   # Sorting
│   │   ├── pagination.ts             # Pagination helpers
│   │   └── aggregation.ts            # Aggregation queries
│   │
│   ├── operations/
│   │   ├── query.ts                  # Query operations
│   │   ├── mutation.ts               # Mutation operations
│   │   ├── subscription.ts           # Subscription operations
│   │   ├── batch.ts                  # Batch operations
│   │   └── transaction.ts            # Transaction support
│   │
│   ├── realtime/
│   │   ├── subscription-manager.ts   # Subscription lifecycle
│   │   ├── live-query.ts             # Live query engine
│   │   ├── event-emitter.ts          # Event handling
│   │   └── reconnection.ts           # Reconnection logic
│   │
│   ├── cache/
│   │   ├── memory-cache.ts           # In-memory cache
│   │   ├── storage-cache.ts          # localStorage/sessionStorage
│   │   ├── indexed-db-cache.ts       # IndexedDB cache
│   │   ├── cache-policy.ts           # Cache policies
│   │   └── normalization.ts          # Data normalization
│   │
│   ├── plugins/
│   │   ├── plugin-system.ts          # Plugin registry
│   │   ├── middleware-pipeline.ts    # Middleware execution
│   │   └── hooks.ts                  # Lifecycle hooks
│   │
│   ├── adapters/
│   │   ├── graphql/
│   │   │   ├── adapter.ts            # GraphQL adapter
│   │   │   ├── query-builder.ts      # GraphQL query builder
│   │   │   └── subscription.ts       # GraphQL subscriptions
│   │   │
│   │   ├── rest/
│   │   │   ├── adapter.ts            # REST adapter
│   │   │   ├── url-builder.ts        # URL construction
│   │   │   └── resource.ts           # Resource mapping
│   │   │
│   │   └── dictamesh/
│   │       ├── adapter.ts            # DictaMesh protocol adapter
│   │       ├── protocol.ts           # DictaMesh-specific protocol
│   │       └── federation.ts         # Federation support
│   │
│   ├── connectors/
│   │   ├── http/
│   │   │   ├── connector.ts          # HTTP connector
│   │   │   ├── fetch.ts              # Fetch API wrapper
│   │   │   └── interceptors.ts       # HTTP interceptors
│   │   │
│   │   ├── websocket/
│   │   │   ├── connector.ts          # WebSocket connector
│   │   │   ├── protocol.ts           # WS protocol
│   │   │   └── reconnect.ts          # Reconnection logic
│   │   │
│   │   └── sse/
│   │       ├── connector.ts          # SSE connector
│   │       └── event-source.ts       # EventSource wrapper
│   │
│   ├── middleware/
│   │   ├── auth.ts                   # Authentication
│   │   ├── retry.ts                  # Retry logic
│   │   ├── logging.ts                # Request logging
│   │   ├── metrics.ts                # Performance metrics
│   │   └── error-handling.ts         # Error handling
│   │
│   ├── utils/
│   │   ├── validation.ts             # Validation utilities
│   │   ├── serialization.ts          # Data serialization
│   │   ├── observability.ts          # Tracing/metrics
│   │   └── platform.ts               # Platform detection
│   │
│   └── index.ts                      # Main entry point
│
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
├── examples/
│   ├── react/
│   ├── vue/
│   ├── node/
│   └── vanilla-js/
│
├── docs/
│   ├── API.md
│   ├── ADAPTERS.md
│   ├── PROTOCOL.md
│   ├── REAL-TIME.md
│   └── MIGRATION.md
│
├── package.json
├── tsconfig.json
├── rollup.config.js
├── jest.config.js
└── README.md
```

## Core Interfaces

### 1. Client Interface

```typescript
interface IDictaMeshClient {
  // Configuration
  configure(config: ClientConfig): void;

  // Adapter management
  registerAdapter(name: string, adapter: IAdapter): void;
  setDefaultAdapter(name: string): void;

  // Operations
  query<T = any>(operation: QueryOperation): Promise<QueryResult<T>>;
  mutate<T = any>(operation: MutationOperation): Promise<MutationResult<T>>;
  subscribe<T = any>(operation: SubscriptionOperation): Subscription<T>;

  // Batch operations
  batch(operations: Operation[]): Promise<BatchResult>;

  // Cache management
  getCache(): ICache;
  clearCache(options?: ClearCacheOptions): Promise<void>;

  // Lifecycle
  connect(): Promise<void>;
  disconnect(): Promise<void>;

  // Status
  isConnected(): boolean;
  getStatus(): ClientStatus;
}
```

### 2. Adapter Interface

```typescript
interface IAdapter {
  // Metadata
  name: string;
  version: string;
  capabilities: AdapterCapabilities;

  // Operations
  execute<T = any>(
    operation: Operation,
    context: OperationContext
  ): Promise<OperationResult<T>>;

  // Real-time
  subscribe<T = any>(
    operation: SubscriptionOperation,
    context: OperationContext
  ): Observable<T>;

  // Schema
  introspect(): Promise<Schema>;

  // Lifecycle
  initialize(config: AdapterConfig): Promise<void>;
  dispose(): Promise<void>;
}
```

### 3. Connector Interface

```typescript
interface IConnector {
  // Transport
  send(message: ProtocolMessage): Promise<ProtocolMessage>;

  // Real-time
  connect(): Promise<void>;
  disconnect(): Promise<void>;
  onMessage(handler: MessageHandler): Unsubscribe;

  // Status
  isConnected(): boolean;
  getConnectionState(): ConnectionState;
}
```

### 4. Protocol Interface

```typescript
interface IProtocol {
  // Message handling
  createRequest(operation: Operation): ProtocolMessage;
  parseResponse(message: ProtocolMessage): OperationResult;

  // Serialization
  serialize(message: ProtocolMessage): string | ArrayBuffer;
  deserialize(data: string | ArrayBuffer): ProtocolMessage;

  // Validation
  validate(message: ProtocolMessage): ValidationResult;

  // Versioning
  getVersion(): string;
  supportsVersion(version: string): boolean;
}
```

## DictaMesh Protocol Specification

### Message Format

```typescript
type ProtocolMessage = {
  version: string;               // Protocol version
  id: string;                    // Unique message ID
  type: MessageType;             // query | mutation | subscription | response | error
  timestamp: number;             // Unix timestamp

  // Request fields
  operation?: {
    type: OperationType;         // get | list | create | update | delete | search
    entity: string;              // Entity type
    params?: Record<string, any>; // Operation parameters
    options?: OperationOptions;  // Query options
  };

  // Response fields
  data?: any;                    // Response data
  meta?: {
    count?: number;
    total?: number;
    page?: number;
    hasMore?: boolean;
    took?: number;               // Execution time in ms
  };

  // Error fields
  error?: {
    code: string;
    message: string;
    details?: any;
    stack?: string;
  };

  // Tracing
  trace?: {
    traceId: string;
    spanId: string;
    parentSpanId?: string;
  };
};

type MessageType =
  | 'query'
  | 'mutation'
  | 'subscription'
  | 'subscription_data'
  | 'subscription_complete'
  | 'response'
  | 'error'
  | 'ping'
  | 'pong';

type OperationType =
  | 'get'           // Get single entity
  | 'list'          // List entities
  | 'create'        // Create entity
  | 'update'        // Update entity
  | 'delete'        // Delete entity
  | 'search'        // Full-text search
  | 'aggregate'     // Aggregation query
  | 'batch'         // Batch operation
  | 'transaction';  // Transactional operation
```

### Query DSL

```typescript
type QueryOptions = {
  // Filtering
  where?: FilterExpression;

  // Sorting
  orderBy?: SortExpression[];

  // Pagination
  limit?: number;
  offset?: number;
  cursor?: string;

  // Field selection
  select?: string[];
  include?: string[];
  exclude?: string[];

  // Relationships
  expand?: string[];

  // Caching
  cache?: CacheOptions;

  // Misc
  timeout?: number;
  signal?: AbortSignal;
};

type FilterExpression = {
  [field: string]:
    | any                          // Exact match
    | {
        $eq?: any;                 // Equal
        $ne?: any;                 // Not equal
        $gt?: any;                 // Greater than
        $gte?: any;                // Greater than or equal
        $lt?: any;                 // Less than
        $lte?: any;                // Less than or equal
        $in?: any[];               // In array
        $nin?: any[];              // Not in array
        $contains?: string;        // Contains substring
        $startsWith?: string;      // Starts with
        $endsWith?: string;        // Ends with
        $regex?: string;           // Regex match
      };
};

type SortExpression = {
  field: string;
  direction: 'asc' | 'desc';
};
```

## Usage Examples

### Basic Setup

```typescript
import { DictaMeshClient, GraphQLAdapter, HTTPConnector } from '@dictamesh/sdk-core';

// Create client
const client = new DictaMeshClient({
  endpoint: 'https://api.example.com/graphql',
  transport: 'http',
  cache: {
    type: 'memory',
    ttl: 300000, // 5 minutes
  },
});

// Register adapter
const adapter = new GraphQLAdapter({
  connector: new HTTPConnector({
    baseURL: 'https://api.example.com',
  }),
});

client.registerAdapter('graphql', adapter);
client.setDefaultAdapter('graphql');

// Connect
await client.connect();
```

### Query Operations

```typescript
// Get single entity
const customer = await client.query({
  entity: 'customer',
  operation: 'get',
  params: { id: '123' },
  options: {
    include: ['email', 'name', 'createdAt'],
  },
});

// List entities with filtering
const customers = await client.query({
  entity: 'customer',
  operation: 'list',
  options: {
    where: {
      status: 'active',
      createdAt: { $gte: '2025-01-01' },
    },
    orderBy: [{ field: 'createdAt', direction: 'desc' }],
    limit: 20,
    offset: 0,
  },
});

// Search
const results = await client.query({
  entity: 'product',
  operation: 'search',
  params: {
    query: 'laptop',
    fields: ['name', 'description'],
  },
  options: {
    limit: 10,
  },
});
```

### Mutation Operations

```typescript
// Create
const newCustomer = await client.mutate({
  entity: 'customer',
  operation: 'create',
  data: {
    email: 'john@example.com',
    name: 'John Doe',
  },
});

// Update
const updated = await client.mutate({
  entity: 'customer',
  operation: 'update',
  params: { id: '123' },
  data: {
    name: 'John Updated',
  },
});

// Delete
await client.mutate({
  entity: 'customer',
  operation: 'delete',
  params: { id: '123' },
});
```

### Subscriptions

```typescript
// Subscribe to entity changes
const subscription = client.subscribe({
  entity: 'customer',
  operation: 'list',
  options: {
    where: { status: 'active' },
  },
});

subscription.subscribe({
  next: (data) => {
    console.log('Data updated:', data);
  },
  error: (error) => {
    console.error('Subscription error:', error);
  },
  complete: () => {
    console.log('Subscription complete');
  },
});

// Unsubscribe
subscription.unsubscribe();
```

### Query Builder API

```typescript
// Fluent query builder
const customers = await client
  .from('customer')
  .select(['id', 'name', 'email'])
  .where('status', 'active')
  .where('createdAt', '>=', '2025-01-01')
  .orderBy('createdAt', 'desc')
  .limit(20)
  .execute();

// With relationships
const invoices = await client
  .from('invoice')
  .select(['id', 'total', 'status'])
  .expand('customer', ['name', 'email'])
  .expand('items.product', ['name', 'price'])
  .where('status', 'pending')
  .execute();
```

### React Integration

```typescript
import { useDictaMesh, useQuery, useMutation, useSubscription } from '@dictamesh/react';

function CustomerList() {
  // Query
  const { data, loading, error, refetch } = useQuery({
    entity: 'customer',
    operation: 'list',
    options: {
      where: { status: 'active' },
      limit: 20,
    },
  });

  // Mutation
  const [createCustomer, { loading: creating }] = useMutation({
    entity: 'customer',
    operation: 'create',
  });

  // Subscription
  useSubscription({
    entity: 'customer',
    operation: 'list',
    onData: (data) => {
      // Auto-updates when data changes
      refetch();
    },
  });

  const handleCreate = async () => {
    await createCustomer({
      data: {
        name: 'New Customer',
        email: 'new@example.com',
      },
    });
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <button onClick={handleCreate} disabled={creating}>
        Create Customer
      </button>
      <ul>
        {data.map(customer => (
          <li key={customer.id}>{customer.name}</li>
        ))}
      </ul>
    </div>
  );
}
```

### Vue Integration

```typescript
import { useDictaMesh, useQuery, useMutation } from '@dictamesh/vue';

export default {
  setup() {
    const { data: customers, loading, error, refetch } = useQuery({
      entity: 'customer',
      operation: 'list',
      options: {
        where: { status: 'active' },
      },
    });

    const { mutate: createCustomer, loading: creating } = useMutation({
      entity: 'customer',
      operation: 'create',
    });

    const handleCreate = async () => {
      await createCustomer({
        data: {
          name: 'New Customer',
          email: 'new@example.com',
        },
      });
      refetch();
    };

    return {
      customers,
      loading,
      error,
      creating,
      handleCreate,
    };
  },
};
```

## Implementation Phases

### Phase 1: Foundation (Week 1-2)
- [ ] Set up TypeScript package with build configuration
- [ ] Implement core abstractions (IAdapter, IConnector, IProtocol)
- [ ] Create protocol message format and serialization
- [ ] Implement basic client with lifecycle management
- [ ] Add unit tests for core functionality

### Phase 2: Operations (Week 3-4)
- [ ] Implement query operations (get, list, search)
- [ ] Implement mutation operations (create, update, delete)
- [ ] Create query builder with fluent API
- [ ] Add filter DSL and validation
- [ ] Implement batch operations

### Phase 3: Real-time (Week 5-6)
- [ ] Implement subscription manager
- [ ] Add WebSocket connector
- [ ] Create live query engine
- [ ] Implement reconnection logic
- [ ] Add subscription lifecycle management

### Phase 4: Caching (Week 7-8)
- [ ] Implement memory cache
- [ ] Add storage-based cache (localStorage)
- [ ] Create IndexedDB cache for large datasets
- [ ] Implement cache normalization
- [ ] Add cache invalidation strategies

### Phase 5: Adapters (Week 9-10)
- [ ] Create GraphQL adapter
- [ ] Implement REST adapter
- [ ] Build DictaMesh protocol adapter
- [ ] Add adapter tests
- [ ] Create adapter documentation

### Phase 6: Middleware & Plugins (Week 11-12)
- [ ] Implement plugin system
- [ ] Create middleware pipeline
- [ ] Add auth middleware
- [ ] Implement retry middleware
- [ ] Create logging and metrics middleware

### Phase 7: Framework Integrations (Week 13-14)
- [ ] Create React hooks and components
- [ ] Build Vue composables
- [ ] Add Angular services (optional)
- [ ] Create example applications
- [ ] Add integration tests

### Phase 8: Documentation & Release (Week 15-16)
- [ ] Write comprehensive API documentation
- [ ] Create adapter development guide
- [ ] Write migration guides
- [ ] Set up documentation website
- [ ] Prepare for npm release

## Success Criteria

1. **Functionality**
   - [ ] All CRUD operations working
   - [ ] Real-time subscriptions functional
   - [ ] Query builder API complete
   - [ ] Multiple adapters implemented
   - [ ] Framework integrations ready

2. **Performance**
   - [ ] Initial load < 50KB gzipped
   - [ ] Query execution < 100ms (cached)
   - [ ] Real-time latency < 500ms
   - [ ] Memory usage < 10MB for typical app

3. **Developer Experience**
   - [ ] Full TypeScript support
   - [ ] Comprehensive documentation
   - [ ] Example applications
   - [ ] Clear error messages
   - [ ] Easy adapter development

4. **Quality**
   - [ ] 90%+ test coverage
   - [ ] Zero critical bugs
   - [ ] Performance benchmarks
   - [ ] Security audit passed

## Technical Decisions

### TypeScript Configuration
- Target: ES2020
- Module: ESNext with dual ESM/CJS output
- Strict mode enabled
- Declaration maps for debugging

### Build System
- Rollup for bundling
- Tree-shaking enabled
- Multiple output formats (ESM, CJS, UMD)
- Source maps included

### Testing
- Jest for unit tests
- Testing Library for integration tests
- Playwright for E2E tests
- Coverage threshold: 90%

### Code Quality
- ESLint with strict rules
- Prettier for formatting
- Husky for pre-commit hooks
- Conventional commits

### Documentation
- TSDoc for API documentation
- Markdown for guides
- Live code examples
- Interactive playground

## Dependencies

### Core Dependencies
- None (zero dependencies for core)

### Optional Dependencies
- `graphql` - For GraphQL adapter
- `ws` - For WebSocket connector (Node.js)
- `eventsource` - For SSE connector

### Dev Dependencies
- TypeScript
- Rollup
- Jest
- ESLint
- Prettier
- TSDoc

## Next Steps

1. Create package structure and configuration
2. Implement core abstractions
3. Create protocol specification
4. Build basic client
5. Add first adapter (GraphQL)
6. Create example application
7. Write documentation

## Questions to Resolve

1. Should we support offline-first mode from the start?
2. What level of GraphQL compatibility do we need?
3. Should we include React/Vue integrations in core or separate packages?
4. Do we need schema code generation?
5. What's the migration path from existing solutions?

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
