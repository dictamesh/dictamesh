// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps the OpenTelemetry tracer
type Tracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
	config   *TracingConfig
}

// NewTracer creates a new tracer from configuration
func NewTracer(cfg *TracingConfig, serviceName, serviceVersion, environment string) (*Tracer, error) {
	if !cfg.Enabled {
		// Return a no-op tracer
		return &Tracer{
			tracer: trace.NewNoopTracerProvider().Tracer(serviceName),
			config: cfg,
		}, nil
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			attribute.String("environment", environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create exporter based on configuration
	var exporter sdktrace.SpanExporter

	if cfg.UseJaeger && cfg.JaegerEndpoint != "" {
		// Legacy Jaeger exporter
		exporter, err = jaeger.New(
			jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
		}
	} else {
		// Use OTLP exporter (default)
		// Note: This would typically use the OTLP HTTP or gRPC exporter
		// For now, using Jaeger as a placeholder since OTLP requires additional setup
		exporter, err = jaeger.New(
			jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://" + cfg.Endpoint + "/api/traces")),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}
	}

	// Configure batch span processor
	batchTimeout := cfg.ExportTimeout
	if batchTimeout == 0 {
		batchTimeout = 5 * time.Second
	}

	maxExportBatchSize := cfg.MaxExportBatchSize
	if maxExportBatchSize == 0 {
		maxExportBatchSize = 512
	}

	maxQueueSize := cfg.MaxQueueSize
	if maxQueueSize == 0 {
		maxQueueSize = 2048
	}

	bsp := sdktrace.NewBatchSpanProcessor(
		exporter,
		sdktrace.WithBatchTimeout(batchTimeout),
		sdktrace.WithMaxExportBatchSize(maxExportBatchSize),
		sdktrace.WithMaxQueueSize(maxQueueSize),
	)

	// Create trace provider with sampling
	samplingRate := cfg.SamplingRate
	if samplingRate <= 0 || samplingRate > 1 {
		samplingRate = 1.0
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(samplingRate)),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator for context propagation (W3C Trace Context)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return &Tracer{
		provider: tp,
		tracer:   tp.Tracer(serviceName),
		config:   cfg,
	}, nil
}

// Shutdown gracefully shuts down the tracer provider
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.provider == nil {
		return nil
	}
	return t.provider.Shutdown(ctx)
}

// StartSpan starts a new span
func (t *Tracer) StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.RecordError(err, opts...)
		span.SetStatus(codes.Error, err.Error())
	}
}

// SetStatus sets the status of the current span
func SetStatus(ctx context.Context, code codes.Code, description string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetStatus(code, description)
	}
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, opts...)
	}
}

// SetAttributes sets attributes on the current span
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// SpanFromContext returns the current span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// WithSpan is a helper that creates a span, executes a function, and ends the span
// It automatically records errors and sets the span status
func WithSpan(ctx context.Context, tracer trace.Tracer, spanName string, fn func(context.Context) error, opts ...trace.SpanStartOption) error {
	ctx, span := tracer.Start(ctx, spanName, opts...)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// TraceID returns the trace ID from the current span context
func TraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// SpanID returns the span ID from the current span context
func SpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// Common attribute keys for consistency
var (
	AttrHTTPMethod     = attribute.Key("http.method")
	AttrHTTPURL        = attribute.Key("http.url")
	AttrHTTPStatusCode = attribute.Key("http.status_code")
	AttrDBSystem       = attribute.Key("db.system")
	AttrDBName         = attribute.Key("db.name")
	AttrDBOperation    = attribute.Key("db.operation")
	AttrDBStatement    = attribute.Key("db.statement")
	AttrMessageBus     = attribute.Key("messaging.system")
	AttrMessageTopic   = attribute.Key("messaging.destination")
	AttrMessageID      = attribute.Key("messaging.message_id")
	AttrUserID         = attribute.Key("user.id")
	AttrTenantID       = attribute.Key("tenant.id")
	AttrEntityType     = attribute.Key("entity.type")
	AttrEntityID       = attribute.Key("entity.id")
	AttrAdapterName    = attribute.Key("adapter.name")
	AttrSourceSystem   = attribute.Key("source.system")
)
