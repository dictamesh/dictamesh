// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional functionality for distributed tracing
type Logger struct {
	*zap.Logger
	config *LoggingConfig
}

// NewLogger creates a new logger from configuration
func NewLogger(cfg *LoggingConfig) (*Logger, error) {
	zapCfg := zap.NewProductionConfig()

	// Set level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", cfg.Level, err)
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	// Set format
	if cfg.Format == "console" {
		zapCfg.Encoding = "console"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapCfg.Encoding = "json"
	}

	// Set output paths
	if len(cfg.OutputPaths) > 0 {
		zapCfg.OutputPaths = cfg.OutputPaths
	}
	if len(cfg.ErrorOutputPaths) > 0 {
		zapCfg.ErrorOutputPaths = cfg.ErrorOutputPaths
	}

	// Configure caller
	zapCfg.DisableCaller = !cfg.EnableCaller

	// Configure stack trace
	if cfg.EnableStackTrace {
		zapCfg.DisableStacktrace = false
	} else {
		// Only stack trace on errors and above
		zapCfg.DisableStacktrace = true
	}

	// Configure sampling
	if cfg.SamplingInitial > 0 && cfg.SamplingThereafter > 0 {
		zapCfg.Sampling = &zap.SamplingConfig{
			Initial:    cfg.SamplingInitial,
			Thereafter: cfg.SamplingThereafter,
		}
	}

	// Build logger
	zapLogger, err := zapCfg.Build(
		zap.AddCallerSkip(1), // Skip one level to show actual caller
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	// Add default fields
	if len(cfg.Fields) > 0 {
		fields := make([]zap.Field, 0, len(cfg.Fields))
		for k, v := range cfg.Fields {
			fields = append(fields, zap.Any(k, v))
		}
		zapLogger = zapLogger.With(fields...)
	}

	return &Logger{
		Logger: zapLogger,
		config: cfg,
	}, nil
}

// WithContext returns a logger with trace context information
// If the context contains a span, it adds trace_id and span_id fields
func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return l
	}

	spanCtx := span.SpanContext()
	return &Logger{
		Logger: l.Logger.With(
			zap.String("trace_id", spanCtx.TraceID().String()),
			zap.String("span_id", spanCtx.SpanID().String()),
			zap.Bool("trace_sampled", spanCtx.IsSampled()),
		),
		config: l.config,
	}
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{
		Logger: l.Logger.With(zapFields...),
		config: l.config,
	}
}

// WithField returns a logger with an additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(zap.Any(key, value)),
		config: l.config,
	}
}

// WithError returns a logger with an error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.With(zap.Error(err)),
		config: l.config,
	}
}

// Named creates a named logger (adds a "logger" field with the provided name)
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		Logger: l.Logger.Named(name),
		config: l.config,
	}
}

// DebugContext logs at debug level with trace context
func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Debug(msg, fields...)
}

// InfoContext logs at info level with trace context
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Info(msg, fields...)
}

// WarnContext logs at warn level with trace context
func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Warn(msg, fields...)
}

// ErrorContext logs at error level with trace context
func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Error(msg, fields...)
}

// FatalContext logs at fatal level with trace context and exits
func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Fatal(msg, fields...)
}

// LogWithContext is a helper that logs at the appropriate level based on error
func (l *Logger) LogWithContext(ctx context.Context, err error, msg string, fields ...zap.Field) {
	if err != nil {
		l.WithContext(ctx).Error(msg, append(fields, zap.Error(err))...)
	} else {
		l.WithContext(ctx).Info(msg, fields...)
	}
}
