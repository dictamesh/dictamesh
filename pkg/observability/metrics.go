// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package observability

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer manages Prometheus metrics collection and exposition
type MetricsServer struct {
	config   *MetricsConfig
	registry *prometheus.Registry
	server   *http.Server
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(cfg *MetricsConfig) (*MetricsServer, error) {
	if !cfg.Enabled {
		return &MetricsServer{config: cfg}, nil
	}

	registry := prometheus.NewRegistry()

	// Register default Go collectors if enabled
	if cfg.EnableRuntimeMetrics {
		registry.MustRegister(prometheus.NewGoCollector())
		registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	}

	mux := http.NewServeMux()
	mux.Handle(cfg.Path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &MetricsServer{
		config:   cfg,
		registry: registry,
		server:   server,
	}, nil
}

// Start starts the metrics HTTP server
func (m *MetricsServer) Start() error {
	if !m.config.Enabled || m.server == nil {
		return nil
	}

	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log error
			fmt.Printf("metrics server error: %v\n", err)
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the metrics server
func (m *MetricsServer) Shutdown(ctx context.Context) error {
	if m.server == nil {
		return nil
	}
	return m.server.Shutdown(ctx)
}

// Registry returns the Prometheus registry
func (m *MetricsServer) Registry() *prometheus.Registry {
	return m.registry
}

// Metrics holds all application metrics
type Metrics struct {
	registry *prometheus.Registry
	factory  promauto.Factory
	config   *MetricsConfig

	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestSize      *prometheus.HistogramVec
	HTTPResponseSize     *prometheus.HistogramVec
	HTTPRequestsInFlight *prometheus.GaugeVec

	// Database metrics
	DBQueriesTotal     *prometheus.CounterVec
	DBQueryDuration    *prometheus.HistogramVec
	DBConnectionsOpen  *prometheus.GaugeVec
	DBConnectionsIdle  *prometheus.GaugeVec
	DBConnectionsInUse *prometheus.GaugeVec

	// Event bus metrics
	EventsPublishedTotal  *prometheus.CounterVec
	EventsConsumedTotal   *prometheus.CounterVec
	EventPublishDuration  *prometheus.HistogramVec
	EventConsumeDuration  *prometheus.HistogramVec
	EventProcessingErrors *prometheus.CounterVec

	// Adapter metrics
	AdapterRequestsTotal   *prometheus.CounterVec
	AdapterRequestDuration *prometheus.HistogramVec
	AdapterErrors          *prometheus.CounterVec
	AdapterCacheHits       *prometheus.CounterVec
	AdapterCacheMisses     *prometheus.CounterVec

	// GraphQL metrics
	GraphQLQueriesTotal    *prometheus.CounterVec
	GraphQLQueryDuration   *prometheus.HistogramVec
	GraphQLResolverErrors  *prometheus.CounterVec
	GraphQLComplexity      *prometheus.HistogramVec

	// System metrics
	SystemInfo *prometheus.GaugeVec
}

// NewMetrics creates a new metrics collector
func NewMetrics(registry *prometheus.Registry, cfg *MetricsConfig) *Metrics {
	factory := promauto.With(registry)

	namespace := cfg.Namespace
	if namespace == "" {
		namespace = "dictamesh"
	}

	buckets := cfg.DefaultHistogramBuckets
	if buckets == nil {
		buckets = []float64{0.001, 0.01, 0.1, 0.5, 1, 2.5, 5, 10}
	}

	m := &Metrics{
		registry: registry,
		factory:  factory,
		config:   cfg,

		// HTTP metrics
		HTTPRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency in seconds",
				Buckets:   buckets,
			},
			[]string{"method", "path"},
		),
		HTTPRequestSize: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),
		HTTPResponseSize: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),
		HTTPRequestsInFlight: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "http_requests_in_flight",
				Help:      "Number of HTTP requests currently being processed",
			},
			[]string{"method", "path"},
		),

		// Database metrics
		DBQueriesTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table", "status"},
		),
		DBQueryDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_query_duration_seconds",
				Help:      "Database query latency in seconds",
				Buckets:   buckets,
			},
			[]string{"operation", "table"},
		),
		DBConnectionsOpen: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_open",
				Help:      "Number of open database connections",
			},
			[]string{"pool"},
		),
		DBConnectionsIdle: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_idle",
				Help:      "Number of idle database connections",
			},
			[]string{"pool"},
		),
		DBConnectionsInUse: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_in_use",
				Help:      "Number of database connections currently in use",
			},
			[]string{"pool"},
		),

		// Event bus metrics
		EventsPublishedTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "events_published_total",
				Help:      "Total number of events published",
			},
			[]string{"topic", "event_type", "status"},
		),
		EventsConsumedTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "events_consumed_total",
				Help:      "Total number of events consumed",
			},
			[]string{"topic", "consumer_group", "status"},
		),
		EventPublishDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "event_publish_duration_seconds",
				Help:      "Event publish latency in seconds",
				Buckets:   buckets,
			},
			[]string{"topic", "event_type"},
		),
		EventConsumeDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "event_consume_duration_seconds",
				Help:      "Event consumption latency in seconds",
				Buckets:   buckets,
			},
			[]string{"topic", "consumer_group"},
		),
		EventProcessingErrors: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "event_processing_errors_total",
				Help:      "Total number of event processing errors",
			},
			[]string{"topic", "consumer_group", "error_type"},
		),

		// Adapter metrics
		AdapterRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "adapter_requests_total",
				Help:      "Total number of adapter requests",
			},
			[]string{"adapter", "operation", "status"},
		),
		AdapterRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "adapter_request_duration_seconds",
				Help:      "Adapter request latency in seconds",
				Buckets:   buckets,
			},
			[]string{"adapter", "operation"},
		),
		AdapterErrors: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "adapter_errors_total",
				Help:      "Total number of adapter errors",
			},
			[]string{"adapter", "operation", "error_type"},
		),
		AdapterCacheHits: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "adapter_cache_hits_total",
				Help:      "Total number of adapter cache hits",
			},
			[]string{"adapter", "cache_layer"},
		),
		AdapterCacheMisses: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "adapter_cache_misses_total",
				Help:      "Total number of adapter cache misses",
			},
			[]string{"adapter", "cache_layer"},
		),

		// GraphQL metrics
		GraphQLQueriesTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "graphql_queries_total",
				Help:      "Total number of GraphQL queries",
			},
			[]string{"operation_type", "operation_name", "status"},
		),
		GraphQLQueryDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "graphql_query_duration_seconds",
				Help:      "GraphQL query latency in seconds",
				Buckets:   buckets,
			},
			[]string{"operation_type", "operation_name"},
		),
		GraphQLResolverErrors: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "graphql_resolver_errors_total",
				Help:      "Total number of GraphQL resolver errors",
			},
			[]string{"resolver", "error_type"},
		),
		GraphQLComplexity: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "graphql_query_complexity",
				Help:      "GraphQL query complexity score",
				Buckets:   prometheus.LinearBuckets(0, 50, 20),
			},
			[]string{"operation_type"},
		),

		// System metrics
		SystemInfo: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "system_info",
				Help:      "System information",
			},
			[]string{"version", "go_version", "environment"},
		),
	}

	// Set system info
	m.SystemInfo.WithLabelValues("0.1.0", runtime.Version(), "development").Set(1)

	return m
}

// RecordHTTPRequest records metrics for an HTTP request
func (m *Metrics) RecordHTTPRequest(method, path, status string, duration time.Duration, requestSize, responseSize int64) {
	m.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	m.HTTPRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	m.HTTPResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
}

// RecordDBQuery records metrics for a database query
func (m *Metrics) RecordDBQuery(operation, table, status string, duration time.Duration) {
	m.DBQueriesTotal.WithLabelValues(operation, table, status).Inc()
	m.DBQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordEventPublish records metrics for event publishing
func (m *Metrics) RecordEventPublish(topic, eventType, status string, duration time.Duration) {
	m.EventsPublishedTotal.WithLabelValues(topic, eventType, status).Inc()
	m.EventPublishDuration.WithLabelValues(topic, eventType).Observe(duration.Seconds())
}

// RecordEventConsume records metrics for event consumption
func (m *Metrics) RecordEventConsume(topic, consumerGroup, status string, duration time.Duration) {
	m.EventsConsumedTotal.WithLabelValues(topic, consumerGroup, status).Inc()
	m.EventConsumeDuration.WithLabelValues(topic, consumerGroup).Observe(duration.Seconds())
}
