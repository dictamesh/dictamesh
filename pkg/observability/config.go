// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package observability

import "time"

// Config holds all observability configuration
type Config struct {
	// ServiceName is the name of the service for identification in traces and metrics
	ServiceName string

	// ServiceVersion is the version of the service
	ServiceVersion string

	// Environment specifies the deployment environment (dev, staging, prod)
	Environment string

	// Tracing configuration
	Tracing TracingConfig

	// Metrics configuration
	Metrics MetricsConfig

	// Logging configuration
	Logging LoggingConfig

	// Health check configuration
	Health HealthConfig
}

// TracingConfig configures distributed tracing
type TracingConfig struct {
	// Enabled determines if tracing is active
	Enabled bool

	// Endpoint is the OTLP collector endpoint (e.g., "localhost:4318" for HTTP or "localhost:4317" for gRPC)
	Endpoint string

	// SamplingRate is the fraction of traces to sample (0.0 to 1.0)
	// 1.0 means trace everything, 0.1 means trace 10% of requests
	SamplingRate float64

	// JaegerEndpoint is the legacy Jaeger endpoint (deprecated, use OTLP instead)
	// Example: "http://localhost:14268/api/traces"
	JaegerEndpoint string

	// UseJaeger determines if we should use legacy Jaeger exporter
	UseJaeger bool

	// Insecure determines if TLS should be disabled for OTLP
	Insecure bool

	// Headers are custom headers to send with OTLP requests
	Headers map[string]string

	// MaxExportBatchSize is the maximum number of spans to export in a batch
	MaxExportBatchSize int

	// MaxQueueSize is the maximum queue size for spans
	MaxQueueSize int

	// ExportTimeout is the timeout for exporting spans
	ExportTimeout time.Duration
}

// MetricsConfig configures Prometheus metrics
type MetricsConfig struct {
	// Enabled determines if metrics collection is active
	Enabled bool

	// Port is the HTTP port for the Prometheus metrics endpoint
	// Metrics will be available at http://localhost:{Port}/metrics
	Port int

	// Path is the HTTP path for metrics (default: "/metrics")
	Path string

	// Namespace is the Prometheus metrics namespace
	// All metrics will be prefixed with {Namespace}_
	Namespace string

	// DefaultHistogramBuckets defines custom histogram buckets for latency metrics
	// If nil, uses default buckets: [0.001, 0.01, 0.1, 0.5, 1, 2.5, 5, 10]
	DefaultHistogramBuckets []float64

	// EnableRuntimeMetrics enables Go runtime metrics (goroutines, memory, GC, etc.)
	EnableRuntimeMetrics bool
}

// LoggingConfig configures structured logging
type LoggingConfig struct {
	// Level is the minimum log level (debug, info, warn, error, fatal)
	Level string

	// Format is the log format (json or console)
	Format string

	// OutputPaths are the output destinations (stdout, stderr, or file paths)
	OutputPaths []string

	// ErrorOutputPaths are the error output destinations
	ErrorOutputPaths []string

	// EnableStackTrace enables stack traces for errors and above
	EnableStackTrace bool

	// EnableCaller enables logging the caller location
	EnableCaller bool

	// SamplingInitial and SamplingThereafter enable log sampling
	// After SamplingInitial messages are logged, only 1 in SamplingThereafter
	// messages will be logged (useful for high-volume logs)
	SamplingInitial    int
	SamplingThereafter int

	// Fields are custom fields to add to every log entry
	Fields map[string]interface{}
}

// HealthConfig configures health checks
type HealthConfig struct {
	// Enabled determines if health checks are active
	Enabled bool

	// Port is the HTTP port for health check endpoints
	Port int

	// LivenessPath is the HTTP path for liveness checks (default: "/health/live")
	LivenessPath string

	// ReadinessPath is the HTTP path for readiness checks (default: "/health/ready")
	ReadinessPath string

	// StartupPath is the HTTP path for startup checks (default: "/health/startup")
	StartupPath string

	// CheckInterval is how often to run health checks
	CheckInterval time.Duration

	// Timeout is the maximum time to wait for a health check
	Timeout time.Duration
}

// DefaultConfig returns a default configuration suitable for development
func DefaultConfig() *Config {
	return &Config{
		ServiceName:    "dictamesh-service",
		ServiceVersion: "0.1.0",
		Environment:    "development",
		Tracing: TracingConfig{
			Enabled:                true,
			Endpoint:               "localhost:4318", // OTLP HTTP endpoint
			SamplingRate:           1.0,              // Trace everything in dev
			Insecure:               true,
			MaxExportBatchSize:     512,
			MaxQueueSize:           2048,
			ExportTimeout:          30 * time.Second,
			UseJaeger:              false,
		},
		Metrics: MetricsConfig{
			Enabled:              true,
			Port:                 9090,
			Path:                 "/metrics",
			Namespace:            "dictamesh",
			EnableRuntimeMetrics: true,
			DefaultHistogramBuckets: []float64{
				0.001, // 1ms
				0.01,  // 10ms
				0.1,   // 100ms
				0.5,   // 500ms
				1.0,   // 1s
				2.5,   // 2.5s
				5.0,   // 5s
				10.0,  // 10s
			},
		},
		Logging: LoggingConfig{
			Level:              "info",
			Format:             "json",
			OutputPaths:        []string{"stdout"},
			ErrorOutputPaths:   []string{"stderr"},
			EnableStackTrace:   true,
			EnableCaller:       true,
			SamplingInitial:    100,
			SamplingThereafter: 100,
			Fields:             make(map[string]interface{}),
		},
		Health: HealthConfig{
			Enabled:       true,
			Port:          8081,
			LivenessPath:  "/health/live",
			ReadinessPath: "/health/ready",
			StartupPath:   "/health/startup",
			CheckInterval: 10 * time.Second,
			Timeout:       5 * time.Second,
		},
	}
}

// ProductionConfig returns a configuration optimized for production
func ProductionConfig() *Config {
	cfg := DefaultConfig()
	cfg.Environment = "production"
	cfg.Tracing.SamplingRate = 0.1 // Sample 10% in production
	cfg.Tracing.Insecure = false   // Use TLS in production
	cfg.Logging.Level = "info"
	cfg.Logging.EnableStackTrace = false // Disable stack traces for non-errors
	cfg.Logging.SamplingInitial = 100
	cfg.Logging.SamplingThereafter = 100
	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}
