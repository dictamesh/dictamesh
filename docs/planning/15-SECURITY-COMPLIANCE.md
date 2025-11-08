# Security & Compliance

[← Previous: Documentation Planning](14-DOCUMENTATION-PLANNING.md) | [Next: Data Governance →](16-DATA-GOVERNANCE.md)

---

## 🎯 Purpose

Security architecture, authentication, authorization, and compliance requirements.

---

## 🔒 Security Layers

### Authentication (OAuth 2.0 / OIDC)

```go
// Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractBearerToken(r)
        claims, err := validateJWT(token)
        if err != nil {
            http.Error(w, "Unauthorized", 401)
            return
        }
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Network Policies

```yaml
# Deny all by default
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
```

### Secrets Management (Sealed Secrets)

```bash
kubeseal -o yaml < secret.yaml > sealed-secret.yaml
kubectl apply -f sealed-secret.yaml
```

---

[← Previous: Documentation Planning](14-DOCUMENTATION-PLANNING.md) | [Next: Data Governance →](16-DATA-GOVERNANCE.md)
