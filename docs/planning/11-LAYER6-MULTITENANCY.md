# Layer 6: Multi-Tenancy & Isolation

[← Previous: Layer 5 Observability](10-LAYER5-OBSERVABILITY.md) | [Next: Layer 7 Saga Orchestration →](12-LAYER7-SAGA-ORCHESTRATION.md)

---

## 🎯 Purpose

Secure multi-tenant architecture with data isolation and tenant-specific configuration.

---

## 🏢 Tenant Isolation Strategies

### Database-Level Isolation

```sql
-- Row-level security
CREATE POLICY tenant_isolation ON entity_catalog
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE entity_catalog ENABLE ROW LEVEL SECURITY;
```

### Application-Level Context

```go
type TenantContext struct {
    TenantID   string
    TenantName string
    Features   []string
}

func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tenantID := r.Header.Get("X-Tenant-ID")
        ctx := context.WithValue(r.Context(), "tenant", &TenantContext{
            TenantID: tenantID,
        })
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

[← Previous: Layer 5 Observability](10-LAYER5-OBSERVABILITY.md) | [Next: Layer 7 Saga Orchestration →](12-LAYER7-SAGA-ORCHESTRATION.md)
