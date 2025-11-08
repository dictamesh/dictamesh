// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed sql/*.sql
var migrationFS embed.FS

// Migrator handles database schema migrations
type Migrator struct {
	db      *sql.DB
	logger  *zap.Logger
	migrate *migrate.Migrate
}

// MigrationInfo represents information about a migration
type MigrationInfo struct {
	Version   uint
	Dirty     bool
	AppliedAt time.Time
	Name      string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, logger *zap.Logger) (*Migrator, error) {
	// Create source from embedded filesystem
	sourceDriver, err := iofs.New(migrationFS, "sql")
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	// Create database driver
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    "metadata_catalog",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &Migrator{
		db:      db,
		logger:  logger,
		migrate: m,
	}, nil
}

// Up runs all pending migrations
func (m *Migrator) Up(ctx context.Context) error {
	m.logger.Info("running database migrations...")

	if err := m.migrate.Up(); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Info("no pending migrations")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.migrate.Version()
	if err != nil {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.logger.Info("migrations completed successfully",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
	)

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down(ctx context.Context) error {
	m.logger.Warn("rolling back last migration...")

	if err := m.migrate.Down(); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Info("no migrations to roll back")
			return nil
		}
		return fmt.Errorf("failed to roll back migration: %w", err)
	}

	m.logger.Info("migration rolled back successfully")
	return nil
}

// MigrateTo migrates to a specific version
func (m *Migrator) MigrateTo(ctx context.Context, version uint) error {
	m.logger.Info("migrating to specific version",
		zap.Uint("target_version", version),
	)

	if err := m.migrate.Migrate(version); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Info("already at target version")
			return nil
		}
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	m.logger.Info("migrated to target version successfully",
		zap.Uint("version", version),
	)

	return nil
}

// Force sets the migration version without running migrations
// Use with caution - this is for fixing dirty states
func (m *Migrator) Force(version int) error {
	m.logger.Warn("forcing migration version",
		zap.Int("version", version),
	)

	if err := m.migrate.Force(version); err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	return nil
}

// Version returns the current migration version
func (m *Migrator) Version() (version uint, dirty bool, err error) {
	return m.migrate.Version()
}

// GetAppliedMigrations returns a list of applied migrations
func (m *Migrator) GetAppliedMigrations(ctx context.Context) ([]MigrationInfo, error) {
	query := `
		SELECT version, dirty
		FROM schema_migrations
		ORDER BY version DESC
	`

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	var migrations []MigrationInfo
	for rows.Next() {
		var info MigrationInfo
		if err := rows.Scan(&info.Version, &info.Dirty); err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, info)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migrations: %w", err)
	}

	return migrations, nil
}

// Close closes the migrator and releases resources
func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed to close source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close database: %w", dbErr)
	}
	return nil
}

// Validate checks if all migrations are valid
func (m *Migrator) Validate(ctx context.Context) error {
	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is in dirty state at version %d - run Force to fix", version)
	}

	m.logger.Info("database migrations are valid",
		zap.Uint("current_version", version),
	)

	return nil
}
