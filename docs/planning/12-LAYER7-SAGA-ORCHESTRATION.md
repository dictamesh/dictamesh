# Layer 7: Saga Orchestration

[← Previous: Layer 6 Multi-Tenancy](11-LAYER6-MULTITENANCY.md) | [Next: Testing Strategy →](13-TESTING-STRATEGY.md)

---

## 🎯 Purpose

Distributed transaction coordination using Saga pattern for long-running business processes.

---

## 🔄 Saga Pattern Implementation

### Using Temporal Workflows

```go
func CreateInvoiceSaga(ctx workflow.Context, order Order) error {
    // Step 1: Reserve inventory
    var reservationID string
    err := workflow.ExecuteActivity(ctx, ReserveInventory, order.Items).Get(ctx, &reservationID)
    if err != nil {
        return err
    }

    // Step 2: Create invoice
    var invoiceID string
    err = workflow.ExecuteActivity(ctx, CreateInvoice, order).Get(ctx, &invoiceID)
    if err != nil {
        // Compensation: release inventory
        workflow.ExecuteActivity(ctx, ReleaseInventory, reservationID)
        return err
    }

    return nil
}
```

---

[← Previous: Layer 6 Multi-Tenancy](11-LAYER6-MULTITENANCY.md) | [Next: Testing Strategy →](13-TESTING-STRATEGY.md)
