// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package health provides database health checking and monitoring
package health

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Status represents the health status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

// HealthCheck represents a health check result
type HealthCheck struct {
	Status        Status
	Message       string
	CheckedAt     time.Time
	ResponseTime  time.Duration
	Details       map[string]interface{}
}

// Checker provides database health checking
type Checker struct {
	pool   *pgxpool.Pool
	stdDB  *sql.DB
	logger *zap.Logger
}

// NewChecker creates a new health checker
func NewChecker(pool *pgxpool.Pool, stdDB *sql.DB, logger *zap.Logger) *Checker {
	return &Checker{
		pool:   pool,
		stdDB:  stdDB,
		logger: logger,
	}
}

// Check performs a comprehensive health check
func (c *Checker) Check(ctx context.Context) *HealthCheck {
	start := time.Now()

	result := &HealthCheck{
		CheckedAt: start,
		Details:   make(map[string]interface{}),
	}

	// Check 1: Basic connectivity
	if err := c.checkConnectivity(ctx); err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("Connectivity check failed: %v", err)
		result.ResponseTime = time.Since(start)
		return result
	}

	// Check 2: Query execution
	if err := c.checkQueryExecution(ctx); err != nil {
		result.Status = StatusDegraded
		result.Message = fmt.Sprintf("Query execution check failed: %v", err)
		result.ResponseTime = time.Since(start)
		return result
	}

	// Check 3: Connection pool stats
	stats := c.getPoolStats()
	result.Details["pool_stats"] = stats

	// Check if pool is degraded
	if stats.IdleConnections == 0 {
		result.Status = StatusDegraded
		result.Message = "No idle connections available"
	}

	// Check 4: Replication lag (if applicable)
	if lag, err := c.checkReplicationLag(ctx); err == nil {
		result.Details["replication_lag_ms"] = lag
		if lag > 5000 { // More than 5 seconds
			result.Status = StatusDegraded
			result.Message = fmt.Sprintf("High replication lag: %dms", lag)
		}
	}

	// Default to healthy if no issues found
	if result.Status == "" {
		result.Status = StatusHealthy
		result.Message = "All checks passed"
	}

	result.ResponseTime = time.Since(start)
	return result
}

// checkConnectivity checks if the database is reachable
func (c *Checker) checkConnectivity(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return c.pool.Ping(ctx)
}

// checkQueryExecution checks if queries can be executed
func (c *Checker) checkQueryExecution(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var result int
	err := c.pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected query result: %d", result)
	}

	return nil
}

// getPoolStats returns connection pool statistics
func (c *Checker) getPoolStats() map[string]interface{} {
	stats := c.stdDB.Stats()

	return map[string]interface{}{
		"open_connections":      stats.OpenConnections,
		"idle_connections":      stats.Idle,
		"in_use_connections":    stats.InUse,
		"wait_count":            stats.WaitCount,
		"wait_duration_ms":      stats.WaitDuration.Milliseconds(),
		"max_idle_closed":       stats.MaxIdleClosed,
		"max_idle_time_closed":  stats.MaxIdleTimeClosed,
		"max_lifetime_closed":   stats.MaxLifetimeClosed,
	}
}

// checkReplicationLag checks replication lag in milliseconds
func (c *Checker) checkReplicationLag(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Query replication lag
	// This works for PostgreSQL replicas
	query := `
		SELECT CASE
			WHEN pg_is_in_recovery() = false THEN 0
			ELSE EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())) * 1000
		END AS lag_ms
	`

	var lagMs sql.NullFloat64
	err := c.pool.QueryRow(ctx, query).Scan(&lagMs)
	if err != nil {
		return 0, fmt.Errorf("failed to check replication lag: %w", err)
	}

	if !lagMs.Valid {
		return 0, nil // Primary database, no lag
	}

	return int64(lagMs.Float64), nil
}

// CheckTable checks if a specific table exists and is accessible
func (c *Checker) CheckTable(ctx context.Context, tableName string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`

	var exists bool
	if err := c.pool.QueryRow(ctx, query, tableName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check table: %w", err)
	}

	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	return nil
}

// CheckExtension checks if a PostgreSQL extension is installed
func (c *Checker) CheckExtension(ctx context.Context, extName string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT FROM pg_extension
			WHERE extname = $1
		)
	`

	var exists bool
	if err := c.pool.QueryRow(ctx, query, extName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check extension: %w", err)
	}

	if !exists {
		return fmt.Errorf("extension %s is not installed", extName)
	}

	return nil
}

// GetDatabaseSize returns the database size in bytes
func (c *Checker) GetDatabaseSize(ctx context.Context, dbName string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `SELECT pg_database_size($1)`

	var size int64
	if err := c.pool.QueryRow(ctx, query, dbName).Scan(&size); err != nil {
		return 0, fmt.Errorf("failed to get database size: %w", err)
	}

	return size, nil
}

// GetTableStats returns statistics for a table
// Note: tableName should include the dictamesh_ prefix (e.g., "dictamesh_entity_catalog")
func (c *Checker) GetTableStats(ctx context.Context, tableName string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		SELECT
			n_live_tup AS live_tuples,
			n_dead_tup AS dead_tuples,
			n_tup_ins AS inserts,
			n_tup_upd AS updates,
			n_tup_del AS deletes,
			last_vacuum,
			last_autovacuum,
			last_analyze,
			last_autoanalyze
		FROM pg_stat_user_tables
		WHERE relname = $1
	`

	var stats struct {
		LiveTuples      int64
		DeadTuples      int64
		Inserts         int64
		Updates         int64
		Deletes         int64
		LastVacuum      sql.NullTime
		LastAutovacuum  sql.NullTime
		LastAnalyze     sql.NullTime
		LastAutoanalyze sql.NullTime
	}

	err := c.pool.QueryRow(ctx, query, tableName).Scan(
		&stats.LiveTuples,
		&stats.DeadTuples,
		&stats.Inserts,
		&stats.Updates,
		&stats.Deletes,
		&stats.LastVacuum,
		&stats.LastAutovacuum,
		&stats.LastAnalyze,
		&stats.LastAutoanalyze,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get table stats: %w", err)
	}

	return map[string]interface{}{
		"live_tuples":      stats.LiveTuples,
		"dead_tuples":      stats.DeadTuples,
		"inserts":          stats.Inserts,
		"updates":          stats.Updates,
		"deletes":          stats.Deletes,
		"last_vacuum":      stats.LastVacuum.Time,
		"last_autovacuum":  stats.LastAutovacuum.Time,
		"last_analyze":     stats.LastAnalyze.Time,
		"last_autoanalyze": stats.LastAutoanalyze.Time,
	}, nil
}
