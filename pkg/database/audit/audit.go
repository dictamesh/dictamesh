// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package audit provides comprehensive audit logging and compliance tracking
package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Operation represents an audit operation type
type Operation string

const (
	OpCreate Operation = "CREATE"
	OpRead   Operation = "READ"
	OpUpdate Operation = "UPDATE"
	OpDelete Operation = "DELETE"
	OpSearch Operation = "SEARCH"
	OpExport Operation = "EXPORT"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID            string
	UserID        string
	UserEmail     string
	Operation     Operation
	ResourceType  string
	ResourceID    string
	Changes       map[string]interface{}
	Metadata      map[string]interface{}
	IPAddress     string
	UserAgent     string
	Success       bool
	ErrorMessage  string
	TraceID       string
	Timestamp     time.Time
	DurationMs    int64
}

// Logger provides audit logging capabilities
type Logger struct {
	pool    *pgxpool.Pool
	logger  *zap.Logger
	enabled bool
}

// Config represents audit logger configuration
type Config struct {
	Enabled bool
}

// NewLogger creates a new audit logger
func NewLogger(pool *pgxpool.Pool, logger *zap.Logger, config *Config) *Logger {
	return &Logger{
		pool:    pool,
		logger:  logger,
		enabled: config.Enabled,
	}
}

// Log records an audit log entry
func (al *Logger) Log(ctx context.Context, entry *AuditLog) error {
	if !al.enabled {
		return nil
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	query := `
		INSERT INTO audit_logs (
			user_id, user_email, operation, resource_type, resource_id,
			changes, metadata, ip_address, user_agent, success, error_message,
			trace_id, timestamp, duration_ms
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id
	`

	changesJSON, err := json.Marshal(entry.Changes)
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %w", err)
	}

	metadataJSON, err := json.Marshal(entry.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = al.pool.QueryRow(ctx, query,
		entry.UserID,
		entry.UserEmail,
		string(entry.Operation),
		entry.ResourceType,
		entry.ResourceID,
		changesJSON,
		metadataJSON,
		entry.IPAddress,
		entry.UserAgent,
		entry.Success,
		entry.ErrorMessage,
		entry.TraceID,
		entry.Timestamp,
		entry.DurationMs,
	).Scan(&entry.ID)

	if err != nil {
		al.logger.Error("failed to write audit log",
			zap.String("user_id", entry.UserID),
			zap.String("operation", string(entry.Operation)),
			zap.String("resource_type", entry.ResourceType),
			zap.Error(err),
		)
		return fmt.Errorf("failed to write audit log: %w", err)
	}

	return nil
}

// LogDataAccess logs PII/sensitive data access
func (al *Logger) LogDataAccess(ctx context.Context, userID, resourceType, resourceID string, fields []string) error {
	metadata := map[string]interface{}{
		"accessed_fields": fields,
		"pii_access":      true,
	}

	return al.Log(ctx, &AuditLog{
		UserID:       userID,
		Operation:    OpRead,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Metadata:     metadata,
		Success:      true,
	})
}

// Query searches audit logs
func (al *Logger) Query(ctx context.Context, filters *QueryFilters) ([]AuditLog, error) {
	query := `
		SELECT
			id, user_id, user_email, operation, resource_type, resource_id,
			changes, metadata, ip_address, user_agent, success, error_message,
			trace_id, timestamp, duration_ms
		FROM audit_logs
		WHERE 1=1
	`

	args := []interface{}{}
	argNum := 1

	// Build WHERE clause
	if filters.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, filters.UserID)
		argNum++
	}

	if filters.Operation != "" {
		query += fmt.Sprintf(" AND operation = $%d", argNum)
		args = append(args, filters.Operation)
		argNum++
	}

	if filters.ResourceType != "" {
		query += fmt.Sprintf(" AND resource_type = $%d", argNum)
		args = append(args, filters.ResourceType)
		argNum++
	}

	if filters.ResourceID != "" {
		query += fmt.Sprintf(" AND resource_id = $%d", argNum)
		args = append(args, filters.ResourceID)
		argNum++
	}

	if !filters.StartTime.IsZero() {
		query += fmt.Sprintf(" AND timestamp >= $%d", argNum)
		args = append(args, filters.StartTime)
		argNum++
	}

	if !filters.EndTime.IsZero() {
		query += fmt.Sprintf(" AND timestamp <= $%d", argNum)
		args = append(args, filters.EndTime)
		argNum++
	}

	// Add ordering and limit
	query += " ORDER BY timestamp DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, filters.Limit)
		argNum++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argNum)
		args = append(args, filters.Offset)
	}

	rows, err := al.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		var changesJSON, metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.UserEmail,
			&log.Operation,
			&log.ResourceType,
			&log.ResourceID,
			&changesJSON,
			&metadataJSON,
			&log.IPAddress,
			&log.UserAgent,
			&log.Success,
			&log.ErrorMessage,
			&log.TraceID,
			&log.Timestamp,
			&log.DurationMs,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if err := json.Unmarshal(changesJSON, &log.Changes); err != nil {
			al.logger.Warn("failed to unmarshal changes", zap.Error(err))
		}

		if err := json.Unmarshal(metadataJSON, &log.Metadata); err != nil {
			al.logger.Warn("failed to unmarshal metadata", zap.Error(err))
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// QueryFilters represents audit log query filters
type QueryFilters struct {
	UserID       string
	Operation    string
	ResourceType string
	ResourceID   string
	StartTime    time.Time
	EndTime      time.Time
	Limit        int
	Offset       int
}

// GetStatistics returns audit statistics
func (al *Logger) GetStatistics(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) AS total_operations,
			COUNT(DISTINCT user_id) AS unique_users,
			COUNT(CASE WHEN success = false THEN 1 END) AS failed_operations,
			COUNT(CASE WHEN metadata->>'pii_access' = 'true' THEN 1 END) AS pii_accesses,
			AVG(duration_ms) AS avg_duration_ms,
			MAX(duration_ms) AS max_duration_ms
		FROM audit_logs
		WHERE timestamp >= $1 AND timestamp <= $2
	`

	var stats struct {
		TotalOperations  int64
		UniqueUsers      int64
		FailedOperations int64
		PIIAccesses      int64
		AvgDurationMs    sql.NullFloat64
		MaxDurationMs    sql.NullInt64
	}

	err := al.pool.QueryRow(ctx, query, startTime, endTime).Scan(
		&stats.TotalOperations,
		&stats.UniqueUsers,
		&stats.FailedOperations,
		&stats.PIIAccesses,
		&stats.AvgDurationMs,
		&stats.MaxDurationMs,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return map[string]interface{}{
		"total_operations":  stats.TotalOperations,
		"unique_users":      stats.UniqueUsers,
		"failed_operations": stats.FailedOperations,
		"pii_accesses":      stats.PIIAccesses,
		"avg_duration_ms":   stats.AvgDurationMs.Float64,
		"max_duration_ms":   stats.MaxDurationMs.Int64,
	}, nil
}

// CreateAuditTable creates the audit_logs table if it doesn't exist
func (al *Logger) CreateAuditTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS audit_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR(255) NOT NULL,
			user_email VARCHAR(255),
			operation VARCHAR(50) NOT NULL,
			resource_type VARCHAR(100) NOT NULL,
			resource_id VARCHAR(255),
			changes JSONB,
			metadata JSONB,
			ip_address VARCHAR(45),
			user_agent TEXT,
			success BOOLEAN DEFAULT true,
			error_message TEXT,
			trace_id VARCHAR(64),
			timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			duration_ms BIGINT
		);

		CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_logs(user_id, timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_resource ON audit_logs(resource_type, resource_id, timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_operation ON audit_logs(operation, timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_logs(timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_trace ON audit_logs(trace_id);
		CREATE INDEX IF NOT EXISTS idx_audit_metadata ON audit_logs USING gin(metadata);
		CREATE INDEX IF NOT EXISTS idx_audit_pii ON audit_logs(timestamp DESC)
			WHERE (metadata->>'pii_access')::boolean = true;

		COMMENT ON TABLE audit_logs IS 'Comprehensive audit trail for compliance and security monitoring';
	`

	_, err := al.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create audit table: %w", err)
	}

	return nil
}
