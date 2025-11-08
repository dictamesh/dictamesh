// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package database provides comprehensive database infrastructure for DictaMesh.
// It includes connection pooling, migrations, ORM, vector search, caching, and audit logging.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents the main database connection manager
type Database struct {
	config *Config
	logger *zap.Logger

	// Connection pools
	pool     *pgxpool.Pool // pgx pool for high-performance queries
	gormDB   *gorm.DB      // GORM for ORM operations
	stdDB    *sql.DB       // Standard database/sql for compatibility

	// Cache layer
	cache *redis.Client

	// State management
	mu        sync.RWMutex
	connected bool
	metrics   *Metrics
}

// Metrics tracks database performance metrics
type Metrics struct {
	QueryCount       int64
	QueryErrors      int64
	CacheHits        int64
	CacheMisses      int64
	ConnectionsOpen  int32
	ConnectionsIdle  int32
	AvgQueryDuration time.Duration
}

// New creates a new database instance
func New(config *Config, logger *zap.Logger) (*Database, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	db := &Database{
		config:  config,
		logger:  logger,
		metrics: &Metrics{},
	}

	return db, nil
}

// Connect establishes database connections
func (db *Database) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.connected {
		return fmt.Errorf("database already connected")
	}

	// Create pgx pool for high-performance operations
	if err := db.createPgxPool(ctx); err != nil {
		return fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Create GORM instance for ORM operations
	if err := db.createGormDB(); err != nil {
		db.pool.Close()
		return fmt.Errorf("failed to create GORM instance: %w", err)
	}

	// Get standard database/sql instance
	var err error
	db.stdDB, err = db.gormDB.DB()
	if err != nil {
		db.pool.Close()
		return fmt.Errorf("failed to get standard DB: %w", err)
	}

	// Configure connection pool
	db.configureConnectionPool()

	db.connected = true
	db.logger.Info("database connected successfully",
		zap.String("host", db.config.Host),
		zap.Int("port", db.config.Port),
		zap.String("database", db.config.Database),
	)

	return nil
}

// createPgxPool creates a pgx connection pool
func (db *Database) createPgxPool(ctx context.Context) error {
	config, err := pgxpool.ParseConfig(db.config.DSN())
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Configure pool settings
	config.MaxConns = int32(db.config.MaxOpenConns)
	config.MinConns = int32(db.config.MaxIdleConns / 2)
	config.MaxConnLifetime = db.config.ConnMaxLifetime
	config.MaxConnIdleTime = db.config.ConnMaxIdleTime

	// Configure statement timeout
	config.ConnConfig.RuntimeParams["statement_timeout"] =
		fmt.Sprintf("%dms", db.config.StatementTimeout.Milliseconds())
	config.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] =
		fmt.Sprintf("%dms", db.config.IdleInTxTimeout.Milliseconds())

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.pool = pool
	return nil
}

// createGormDB creates a GORM instance
func (db *Database) createGormDB() error {
	// Configure GORM logger
	var gormLogger logger.Interface
	if db.config.LogLevel == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Create GORM instance
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  db.config.DSN(),
		PreferSimpleProtocol: true, // disables prepared statement cache
	}), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
		SkipDefaultTransaction:                   true, // Improve performance
		PrepareStmt:                              true, // Enable prepared statement cache
	})

	if err != nil {
		return fmt.Errorf("failed to connect with GORM: %w", err)
	}

	db.gormDB = gormDB
	return nil
}

// configureConnectionPool configures the connection pool settings
func (db *Database) configureConnectionPool() {
	db.stdDB.SetMaxOpenConns(db.config.MaxOpenConns)
	db.stdDB.SetMaxIdleConns(db.config.MaxIdleConns)
	db.stdDB.SetConnMaxLifetime(db.config.ConnMaxLifetime)
	db.stdDB.SetConnMaxIdleTime(db.config.ConnMaxIdleTime)
}

// Close closes all database connections
func (db *Database) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if !db.connected {
		return nil
	}

	var errs []error

	// Close pgx pool
	if db.pool != nil {
		db.pool.Close()
	}

	// Close standard DB (GORM uses it internally)
	if db.stdDB != nil {
		if err := db.stdDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close standard DB: %w", err))
		}
	}

	// Close cache
	if db.cache != nil {
		if err := db.cache.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close cache: %w", err))
		}
	}

	db.connected = false

	if len(errs) > 0 {
		return fmt.Errorf("errors closing database: %v", errs)
	}

	db.logger.Info("database connections closed")
	return nil
}

// Pool returns the pgx connection pool for high-performance queries
func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

// GORM returns the GORM instance for ORM operations
func (db *Database) GORM() *gorm.DB {
	return db.gormDB
}

// StdDB returns the standard database/sql instance
func (db *Database) StdDB() *sql.DB {
	return db.stdDB
}

// Ping checks if the database is reachable
func (db *Database) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Stats returns database statistics
func (db *Database) Stats() sql.DBStats {
	return db.stdDB.Stats()
}

// GetMetrics returns current performance metrics
func (db *Database) GetMetrics() *Metrics {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Update connection stats
	stats := db.stdDB.Stats()
	db.metrics.ConnectionsOpen = int32(stats.OpenConnections)
	db.metrics.ConnectionsIdle = int32(stats.Idle)

	return db.metrics
}

// WithTransaction executes a function within a database transaction
func (db *Database) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return db.gormDB.WithContext(ctx).Transaction(fn)
}

// WithPgxTransaction executes a function within a pgx transaction
func (db *Database) WithPgxTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
