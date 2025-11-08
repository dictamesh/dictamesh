# Layer 4: Federated API Gateway

[← Previous: Layer 3 Metadata Catalog](08-LAYER3-METADATA-CATALOG.md) | [Next: Layer 5 Observability →](10-LAYER5-OBSERVABILITY.md)

---

## 🎯 Purpose

GraphQL Federation implementation providing unified API access to distributed data sources.

---

## 🌐 GraphQL Federation Architecture

### Using gqlgen with Federation

```go
// services/graphql-gateway/graph/schema.graphqls
type Query {
  customer(id: ID!): Customer
  customers(limit: Int, offset: Int): [Customer!]!
}

type Customer @key(fields: "id") {
  id: ID!
  email: String!
  name: String!
  invoices: [Invoice!]!
}

extend type Invoice @key(fields: "id") {
  customer: Customer! @provides(fields: "id")
}
```

### DataLoader Implementation

```go
// internal/dataloaders/customer_loader.go
func NewCustomerLoader(adapter adapter.DataProductAdapter) *dataloader.Loader {
    return dataloader.NewBatchedLoader(
        func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
            // Batch fetch customers
            customers := adapter.QueryEntities(ctx, Query{IDs: keys.Keys()})
            // Map results back to original order
            return mapResults(customers, keys)
        },
        dataloader.WithWait(10*time.Millisecond),
        dataloader.WithBatchCapacity(100),
    )
}
```

---

[← Previous: Layer 3 Metadata Catalog](08-LAYER3-METADATA-CATALOG.md) | [Next: Layer 5 Observability →](10-LAYER5-OBSERVABILITY.md)
