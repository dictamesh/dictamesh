# Testing Strategy

[← Previous: Layer 7 Saga Orchestration](12-LAYER7-SAGA-ORCHESTRATION.md) | [Next: Documentation Planning →](14-DOCUMENTATION-PLANNING.md)

---

## 🎯 Purpose

Comprehensive testing strategy covering unit, integration, end-to-end, and load testing.

---

## 🧪 Test Pyramid

```
        /\
       /  \      E2E Tests (10%)
      /    \
     /──────\    Integration Tests (30%)
    /        \
   /          \  Unit Tests (60%)
  /____________\
```

### Unit Testing

```go
// services/customer-adapter/internal/adapter/adapter_test.go
func TestGetEntity(t *testing.T) {
    mockClient := &MockDirectusClient{}
    adapter := NewCustomerAdapter(mockClient)

    entity, err := adapter.GetEntity(context.Background(), "test-id")
    assert.NoError(t, err)
    assert.Equal(t, "test-id", entity.ID)
}
```

### Integration Testing

```yaml
# docker-compose.test.yaml
services:
  postgres:
    image: postgres:15
  kafka:
    image: confluentinc/cp-kafka:7.5.0
  customer-adapter:
    build: ./services/customer-adapter
    depends_on: [postgres, kafka]
```

### Load Testing (k6)

```javascript
// tests/load/api-load-test.js
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 100 },
    { duration: '2m', target: 0 },
  ],
};

export default function () {
  const res = http.post('https://api.dictamesh.controle.digital/graphql', JSON.stringify({
    query: '{ customer(id: "123") { name email } }'
  }));
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'latency < 200ms': (r) => r.timings.duration < 200,
  });
}
```

---

[← Previous: Layer 7 Saga Orchestration](12-LAYER7-SAGA-ORCHESTRATION.md) | [Next: Documentation Planning →](14-DOCUMENTATION-PLANNING.md)
