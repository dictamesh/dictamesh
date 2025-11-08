// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package observability provides comprehensive observability infrastructure
// for the DictaMesh framework, including distributed tracing, metrics collection,
// structured logging, and health checks.
package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

// Observability manages all observability components
type Observability struct {
	config        *Config
	logger        *Logger
	tracer        *Tracer
	metrics       *Metrics
	metricsServer *MetricsServer
	health        *HealthChecker
}

// New creates a new Observability instance with all components
func New(config *Config) (*Observability, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create logger
	logger, err := NewLogger(&config.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Add service metadata to logger
	logger = logger.WithFields(map[string]interface{}{
		"service":     config.ServiceName,
		"version":     config.ServiceVersion,
		"environment": config.Environment,
	})

	// Create tracer
	tracer, err := NewTracer(
		&config.Tracing,
		config.ServiceName,
		config.ServiceVersion,
		config.Environment,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	// Create metrics server
	metricsServer, err := NewMetricsServer(&config.Metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics server: %w", err)
	}

	// Create metrics collector
	var metrics *Metrics
	if metricsServer.Registry() != nil {
		metrics = NewMetrics(metricsServer.Registry(), &config.Metrics)
	}

	// Create health checker
	health := NewHealthChecker(&config.Health)

	return &Observability{
		config:        config,
		logger:        logger,
		tracer:        tracer,
		metrics:       metrics,
		metricsServer: metricsServer,
		health:        health,
	}, nil
}

// Start starts all observability components
func (o *Observability) Start() error {
	// Start metrics server
	if err := o.metricsServer.Start(); err != nil {
		return fmt.Errorf("failed to start metrics server: %w", err)
	}

	o.logger.Info("metrics server started",
		"port", o.config.Metrics.Port,
		"path", o.config.Metrics.Path,
	)

	// Start health check server
	if err := o.health.Start(); err != nil {
		return fmt.Errorf("failed to start health check server: %w", err)
	}

	o.logger.Info("health check server started",
		"port", o.config.Health.Port,
	)

	o.logger.Info("observability initialized",
		"tracing_enabled", o.config.Tracing.Enabled,
		"metrics_enabled", o.config.Metrics.Enabled,
	)

	return nil
}

// Shutdown gracefully shuts down all observability components
func (o *Observability) Shutdown(ctx context.Context) error {
	o.logger.Info("shutting down observability")

	// Shutdown metrics server
	if err := o.metricsServer.Shutdown(ctx); err != nil {
		o.logger.Error("failed to shutdown metrics server", "error", err)
	}

	// Shutdown health check server
	if err := o.health.Shutdown(ctx); err != nil {
		o.logger.Error("failed to shutdown health check server", "error", err)
	}

	// Shutdown tracer
	if err := o.tracer.Shutdown(ctx); err != nil {
		o.logger.Error("failed to shutdown tracer", "error", err)
	}

	// Sync logger
	if err := o.logger.Sync(); err != nil {
		// Ignore sync errors on stdout/stderr
		// This is a known issue with zap
	}

	return nil
}

// Logger returns the logger instance
func (o *Observability) Logger() *Logger {
	return o.logger
}

// Tracer returns the OpenTelemetry tracer
func (o *Observability) Tracer() trace.Tracer {
	return o.tracer.tracer
}

// Metrics returns the metrics collector
func (o *Observability) Metrics() *Metrics {
	return o.metrics
}

// Health returns the health checker
func (o *Observability) Health() *HealthChecker {
	return o.health
}

// Config returns the configuration
func (o *Observability) Config() *Config {
	return o.config
}

// StartSpan is a convenience method to start a new span
func (o *Observability) StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return o.tracer.StartSpan(ctx, spanName, opts...)
}

// LoggerWithContext returns a logger enriched with trace context
func (o *Observability) LoggerWithContext(ctx context.Context) *Logger {
	return o.logger.WithContext(ctx)
}
