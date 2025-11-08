# Layer 5: Observability & Governance

[← Previous: Layer 4 API Gateway](09-LAYER4-API-GATEWAY.md) | [Next: Layer 6 Multi-Tenancy →](11-LAYER6-MULTITENANCY.md)

---

## 🎯 Purpose

Comprehensive observability stack with distributed tracing, metrics, and logging.

---

## 📊 OpenTelemetry Implementation

```go
// pkg/telemetry/tracer.go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
)

func InitTracer(serviceName string) error {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger-collector:14268/api/traces"),
    ))
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.ServiceNameKey.String(serviceName),
        )),
    )
    
    otel.SetTracerProvider(tp)
    return nil
}
```

### Prometheus Metrics

```go
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Buckets: []float64{.001, .01, .1, .5, 1, 5},
        },
        []string{"service", "method", "status"},
    )
)
```

---

[← Previous: Layer 4 API Gateway](09-LAYER4-API-GATEWAY.md) | [Next: Layer 6 Multi-Tenancy →](11-LAYER6-MULTITENANCY.md)
