<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 2
---

# GraphQL API Reference

Complete reference for the DictaMesh federated GraphQL API.

## Endpoint

```
http://localhost:4000/graphql
```

GraphQL Playground available at the same URL in development mode.

## Authentication

Include authentication token in the `Authorization` header:

```graphql
# HTTP Headers
{
  "Authorization": "Bearer <your-token>"
}
```

## Schema Overview

DictaMesh uses Apollo Federation to compose a unified schema from multiple adapter subgraphs.

### Core Types

#### Entity

Base entity type shared across all domains.

```graphql
interface Entity {
  """Unique entity identifier"""
  id: ID!

  """Entity type (product, customer, order, etc.)"""
  type: String!

  """Domain this entity belongs to"""
  domain: String!

  """Source system identifier"""
  sourceSystem: String!

  """Timestamps"""
  createdAt: DateTime!
  updatedAt: DateTime!
  lastSeenAt: DateTime!
}
```

#### Product

Product entity from e-commerce domain.

```graphql
type Product implements Entity @key(fields: "id") {
  """Unique product identifier"""
  id: ID!

  """Entity metadata"""
  type: String!
  domain: String!
  sourceSystem: String!

  """Product attributes"""
  name: String!
  description: String
  sku: String!
  price: Money!
  currency: String!
  inStock: Boolean!
  stockQuantity: Int

  """Relationships"""
  category: Category
  images: [ProductImage!]!
  variants: [ProductVariant!]!

  """Metadata"""
  metadata: ProductMetadata!

  """Timestamps"""
  createdAt: DateTime!
  updatedAt: DateTime!
  lastSeenAt: DateTime!
}
```

#### Category

Product category hierarchy.

```graphql
type Category @key(fields: "id") {
  id: ID!
  name: String!
  slug: String!
  description: String
  parentId: ID
  parent: Category
  children: [Category!]!
  productCount: Int!
}
```

#### Money

Monetary amount with currency.

```graphql
type Money @shareable {
  """Decimal amount"""
  amount: Decimal!

  """ISO 4217 currency code"""
  currency: String!

  """Formatted display value"""
  formatted: String!
}
```

#### ProductImage

Product image representation.

```graphql
type ProductImage {
  url: String!
  alt: String
  width: Int!
  height: Int!
  isPrimary: Boolean!
}
```

#### ProductVariant

Product variant (size, color, etc.).

```graphql
type ProductVariant {
  id: ID!
  sku: String!
  name: String!
  price: Money
  inStock: Boolean!
  attributes: [VariantAttribute!]!
}
```

#### ProductMetadata

Additional product metadata.

```graphql
type ProductMetadata {
  views: Int!
  favorites: Int!
  rating: Float
  reviewCount: Int!
  tags: [String!]!
}
```

### Pagination

Connection-based pagination following Relay specification.

```graphql
type ProductConnection {
  """List of product nodes"""
  nodes: [Product!]!

  """List of product edges"""
  edges: [ProductEdge!]!

  """Pagination information"""
  pageInfo: PageInfo!

  """Total count of products"""
  totalCount: Int!
}

type ProductEdge {
  """Product node"""
  node: Product!

  """Cursor for this edge"""
  cursor: String!
}

type PageInfo {
  """Whether more results exist"""
  hasNextPage: Boolean!

  """Whether previous results exist"""
  hasPreviousPage: Boolean!

  """Start cursor"""
  startCursor: String

  """End cursor"""
  endCursor: String
}
```

## Queries

### Query Root

```graphql
type Query {
  """Get entity by ID"""
  entity(id: ID!): Entity

  """Get product by ID"""
  product(id: ID!): Product

  """List products with pagination"""
  products(
    first: Int = 20
    after: String
    filters: ProductFilters
    sort: ProductSort
  ): ProductConnection!

  """Search products"""
  searchProducts(
    query: String!
    filters: ProductFilters
    first: Int = 20
  ): ProductConnection!

  """Get category by ID"""
  category(id: ID!): Category

  """List categories"""
  categories(
    parentId: ID
    includeEmpty: Boolean = false
  ): [Category!]!

  """Get metadata catalog entry"""
  catalogEntry(id: ID!): CatalogEntry

  """Search catalog entries"""
  searchCatalog(
    query: String!
    filters: CatalogFilters
    first: Int = 20
  ): CatalogConnection!
}
```

### Input Types

```graphql
"""Product filtering options"""
input ProductFilters {
  categoryId: ID
  minPrice: Decimal
  maxPrice: Decimal
  inStock: Boolean
  tags: [String!]
  sourceSystem: String
}

"""Product sorting options"""
input ProductSort {
  field: ProductSortField!
  direction: SortDirection!
}

enum ProductSortField {
  NAME
  PRICE
  CREATED_AT
  UPDATED_AT
}

enum SortDirection {
  ASC
  DESC
}

"""Catalog filtering options"""
input CatalogFilters {
  entityType: String
  domain: String
  sourceSystem: String
  containsPII: Boolean
  status: EntityStatus
}

enum EntityStatus {
  ACTIVE
  INACTIVE
  DEPRECATED
  ARCHIVED
}
```

### Scalar Types

```graphql
"""ISO 8601 DateTime"""
scalar DateTime

"""Arbitrary precision decimal"""
scalar Decimal

"""JSON object"""
scalar JSON
```

## Example Queries

### Get Product Details

```graphql
query GetProduct($id: ID!) {
  product(id: $id) {
    id
    name
    description
    sku
    price {
      amount
      currency
      formatted
    }
    inStock
    stockQuantity
    category {
      id
      name
      slug
    }
    images {
      url
      alt
      isPrimary
    }
    metadata {
      views
      rating
      reviewCount
      tags
    }
  }
}
```

**Variables:**
```json
{
  "id": "prod-12345"
}
```

**Response:**
```json
{
  "data": {
    "product": {
      "id": "prod-12345",
      "name": "Premium Wireless Headphones",
      "description": "High-quality noise-canceling headphones",
      "sku": "WH-1000XM4",
      "price": {
        "amount": 299.99,
        "currency": "USD",
        "formatted": "$299.99"
      },
      "inStock": true,
      "stockQuantity": 42,
      "category": {
        "id": "cat-electronics",
        "name": "Electronics",
        "slug": "electronics"
      },
      "images": [
        {
          "url": "https://cdn.example.com/headphones.jpg",
          "alt": "Premium Wireless Headphones",
          "isPrimary": true
        }
      ],
      "metadata": {
        "views": 1523,
        "rating": 4.7,
        "reviewCount": 89,
        "tags": ["wireless", "noise-canceling", "premium"]
      }
    }
  }
}
```

### List Products with Filters

```graphql
query ListProducts($first: Int!, $filters: ProductFilters, $sort: ProductSort) {
  products(first: $first, filters: $filters, sort: $sort) {
    nodes {
      id
      name
      price {
        amount
        currency
      }
      inStock
      category {
        name
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
  }
}
```

**Variables:**
```json
{
  "first": 10,
  "filters": {
    "categoryId": "cat-electronics",
    "inStock": true,
    "minPrice": 100,
    "maxPrice": 500
  },
  "sort": {
    "field": "PRICE",
    "direction": "ASC"
  }
}
```

### Search Products

```graphql
query SearchProducts($query: String!, $first: Int!) {
  searchProducts(query: $query, first: $first) {
    nodes {
      id
      name
      description
      price {
        formatted
      }
      images {
        url
        isPrimary
      }
    }
    totalCount
  }
}
```

**Variables:**
```json
{
  "query": "wireless headphones",
  "first": 20
}
```

### Federated Query Across Subgraphs

This query fetches product data and extends it with reviews from another subgraph:

```graphql
query ProductWithReviews($id: ID!) {
  product(id: $id) {
    # From products subgraph
    id
    name
    price {
      formatted
    }

    # Extended from reviews subgraph
    reviews(first: 5) {
      nodes {
        id
        rating
        title
        comment
        author {
          name
        }
        createdAt
      }
      averageRating
      totalCount
    }

    # Extended from inventory subgraph
    inventory {
      warehouseId
      quantity
      reservedQuantity
      availableQuantity
    }
  }
}
```

### Get Category Hierarchy

```graphql
query GetCategoryTree {
  categories(parentId: null) {
    id
    name
    slug
    productCount
    children {
      id
      name
      slug
      productCount
      children {
        id
        name
        slug
        productCount
      }
    }
  }
}
```

### Search Metadata Catalog

```graphql
query SearchCatalog($query: String!, $filters: CatalogFilters) {
  searchCatalog(query: $query, filters: $filters, first: 20) {
    nodes {
      id
      entityType
      domain
      sourceSystem
      status
      containsPII
      availabilitySLA
      latencyP99Ms
      createdAt
    }
    totalCount
  }
}
```

**Variables:**
```json
{
  "query": "product",
  "filters": {
    "domain": "ecommerce",
    "status": "ACTIVE",
    "containsPII": false
  }
}
```

## Mutations

### Mutation Root

```graphql
type Mutation {
  """Create a new catalog entry"""
  createCatalogEntry(input: CreateCatalogEntryInput!): CatalogEntry!

  """Update a catalog entry"""
  updateCatalogEntry(id: ID!, input: UpdateCatalogEntryInput!): CatalogEntry!

  """Delete a catalog entry"""
  deleteCatalogEntry(id: ID!): Boolean!

  """Create entity relationship"""
  createRelationship(input: CreateRelationshipInput!): EntityRelationship!

  """Register schema version"""
  registerSchema(input: RegisterSchemaInput!): SchemaVersion!
}
```

### Create Catalog Entry

```graphql
mutation CreateCatalogEntry($input: CreateCatalogEntryInput!) {
  createCatalogEntry(input: $input) {
    id
    entityType
    domain
    sourceSystem
    status
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "entityType": "product",
    "domain": "ecommerce",
    "sourceSystem": "shopify",
    "sourceEntityId": "prod-12345",
    "apiBaseUrl": "https://api.shopify.com/v1",
    "apiPathTemplate": "/products/{id}",
    "apiMethod": "GET",
    "containsPII": false,
    "availabilitySLA": 0.999
  }
}
```

### Update Catalog Entry

```graphql
mutation UpdateCatalogEntry($id: ID!, $input: UpdateCatalogEntryInput!) {
  updateCatalogEntry(id: $id, input: $input) {
    id
    status
    availabilitySLA
    updatedAt
  }
}
```

**Variables:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "input": {
    "status": "DEPRECATED",
    "availabilitySLA": 0.95
  }
}
```

## Subscriptions

Real-time updates via GraphQL subscriptions.

```graphql
type Subscription {
  """Subscribe to entity changes"""
  entityChanged(entityType: String, entityId: ID): EntityChangeEvent!

  """Subscribe to catalog changes"""
  catalogChanged(domain: String): CatalogChangeEvent!
}

type EntityChangeEvent {
  eventType: ChangeEventType!
  entity: Entity!
  changedFields: [String!]!
  timestamp: DateTime!
}

type CatalogChangeEvent {
  eventType: ChangeEventType!
  catalogEntry: CatalogEntry!
  timestamp: DateTime!
}

enum ChangeEventType {
  CREATED
  UPDATED
  DELETED
}
```

### Subscribe to Product Changes

```graphql
subscription WatchProduct($productId: ID!) {
  entityChanged(entityType: "product", entityId: $productId) {
    eventType
    entity {
      ... on Product {
        id
        name
        price {
          formatted
        }
        inStock
      }
    }
    changedFields
    timestamp
  }
}
```

## Introspection

### Get Schema

```graphql
query GetSchema {
  __schema {
    types {
      name
      kind
      description
    }
    queryType {
      name
    }
    mutationType {
      name
    }
    subscriptionType {
      name
    }
  }
}
```

### Get Type Information

```graphql
query GetTypeInfo($typeName: String!) {
  __type(name: $typeName) {
    name
    kind
    description
    fields {
      name
      type {
        name
        kind
      }
      description
    }
  }
}
```

## Error Handling

GraphQL errors follow the standard format:

```json
{
  "errors": [
    {
      "message": "Entity not found",
      "locations": [{"line": 2, "column": 3}],
      "path": ["product"],
      "extensions": {
        "code": "ENTITY_NOT_FOUND",
        "entityId": "prod-12345",
        "traceId": "trace-abc123"
      }
    }
  ],
  "data": {
    "product": null
  }
}
```

### Error Codes

| Code | Description |
|------|-------------|
| `UNAUTHENTICATED` | Authentication required |
| `FORBIDDEN` | Insufficient permissions |
| `ENTITY_NOT_FOUND` | Requested entity not found |
| `VALIDATION_ERROR` | Input validation failed |
| `INTERNAL_ERROR` | Internal server error |
| `RATE_LIMIT_EXCEEDED` | Too many requests |

## Performance Optimization

### DataLoader Batching

DictaMesh automatically batches requests using DataLoader within a 10ms window:

```graphql
# This query will batch-load all products in one database query
query GetMultipleProducts {
  product1: product(id: "prod-1") { name }
  product2: product(id: "prod-2") { name }
  product3: product(id: "prod-3") { name }
}
```

### Query Complexity

Queries are limited by complexity score to prevent expensive operations:

- Default limit: 1000
- Each field adds complexity based on expected cost
- Pagination multiplies complexity by item count

### Automatic Persisted Queries (APQ)

Enable APQ to reduce bandwidth:

```javascript
// Client sends query hash
{
  "operationName": "GetProduct",
  "variables": {"id": "prod-12345"},
  "extensions": {
    "persistedQuery": {
      "version": 1,
      "sha256Hash": "abc123..."
    }
  }
}
```

## Best Practices

### 1. Request Only Needed Fields

```graphql
# Good - specific fields
query {
  product(id: "prod-12345") {
    id
    name
    price { amount currency }
  }
}

# Avoid - requesting everything
query {
  product(id: "prod-12345") {
    id
    name
    description
    sku
    price { amount currency formatted }
    images { url alt width height isPrimary }
    variants { ... }
    metadata { ... }
  }
}
```

### 2. Use Pagination

```graphql
# Good - paginated
query {
  products(first: 20) {
    nodes { id name }
    pageInfo { hasNextPage endCursor }
  }
}

# Avoid - no pagination
query {
  products(first: 1000) {
    nodes { id name }
  }
}
```

### 3. Leverage Fragments

```graphql
fragment ProductBasic on Product {
  id
  name
  price { formatted }
  inStock
}

query GetProducts {
  featured: products(first: 5, filters: {tags: ["featured"]}) {
    nodes { ...ProductBasic }
  }

  sale: products(first: 5, filters: {tags: ["sale"]}) {
    nodes { ...ProductBasic }
  }
}
```

### 4. Use Variables

```graphql
# Good - with variables
query GetProduct($id: ID!) {
  product(id: $id) { name }
}

# Avoid - hardcoded values
query {
  product(id: "prod-12345") { name }
}
```

## Next Steps

- [REST API Reference](./rest-api.md) - REST API for metadata catalog
- [Go Packages Reference](./go-packages.md) - Build adapters programmatically
- [Event Schemas Reference](./event-schemas.md) - Event-driven integration
- [GraphQL Federation Guide](../guides/graphql-federation.md) - Build federated schemas

---

**Previous**: [← REST API](./rest-api.md) | **Next**: [Go Packages →](./go-packages.md)
