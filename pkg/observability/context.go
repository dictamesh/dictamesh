// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package observability

import (
	"context"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// ContextKeyLogger is the context key for the logger
	ContextKeyLogger contextKey = "observability.logger"

	// ContextKeyRequestID is the context key for request ID
	ContextKeyRequestID contextKey = "observability.request_id"

	// ContextKeyUserID is the context key for user ID
	ContextKeyUserID contextKey = "observability.user_id"

	// ContextKeyTenantID is the context key for tenant ID
	ContextKeyTenantID contextKey = "observability.tenant_id"

	// ContextKeyCorrelationID is the context key for correlation ID
	ContextKeyCorrelationID contextKey = "observability.correlation_id"
)

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

// LoggerFromContext retrieves the logger from context
// If no logger is found, returns nil
func LoggerFromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ContextKeyLogger).(*Logger); ok {
		return logger
	}
	return nil
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

// RequestIDFromContext retrieves the request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		return requestID
	}
	return ""
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// UserIDFromContext retrieves the user ID from context
func UserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// WithTenantID adds a tenant ID to the context
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, ContextKeyTenantID, tenantID)
}

// TenantIDFromContext retrieves the tenant ID from context
func TenantIDFromContext(ctx context.Context) string {
	if tenantID, ok := ctx.Value(ContextKeyTenantID).(string); ok {
		return tenantID
	}
	return ""
}

// WithCorrelationID adds a correlation ID to the context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, ContextKeyCorrelationID, correlationID)
}

// CorrelationIDFromContext retrieves the correlation ID from context
func CorrelationIDFromContext(ctx context.Context) string {
	if correlationID, ok := ctx.Value(ContextKeyCorrelationID).(string); ok {
		return correlationID
	}
	return ""
}

// EnrichContext enriches a context with all available observability metadata
// This is useful when propagating context across service boundaries
func EnrichContext(ctx context.Context) context.Context {
	// Add trace context if available
	traceID := TraceID(ctx)
	spanID := SpanID(ctx)

	// Get logger from context and enrich it with trace context
	logger := LoggerFromContext(ctx)
	if logger != nil && traceID != "" {
		logger = logger.WithContext(ctx)
		ctx = WithLogger(ctx, logger)
	}

	return ctx
}
