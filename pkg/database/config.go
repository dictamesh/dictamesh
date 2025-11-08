// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package database

import (
	"fmt"
	"time"
)

// Config represents database configuration
type Config struct {
	// Connection settings
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string

	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// Performance settings
	StatementTimeout time.Duration
	IdleInTxTimeout  time.Duration

	// Feature flags
	EnableMigrations   bool
	EnableVectorSearch bool
	EnableAuditLog     bool

	// Observability
	EnableMetrics bool
	EnableTracing bool
	LogLevel      string
}

// DefaultConfig returns a production-ready default configuration
func DefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     5432,
		User:     "dictamesh",
		Password: "",
		Database: "metadata_catalog",
		SSLMode:  "prefer",

		// Conservative defaults for production
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,

		StatementTimeout: 30 * time.Second,
		IdleInTxTimeout:  60 * time.Second,

		EnableMigrations:   true,
		EnableVectorSearch: false,
		EnableAuditLog:     true,

		EnableMetrics: true,
		EnableTracing: true,
		LogLevel:      "info",
	}
}

// DSN returns the PostgreSQL connection string
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Port)
	}
	if c.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if c.MaxOpenConns < 1 {
		return fmt.Errorf("max open connections must be at least 1")
	}
	if c.MaxIdleConns < 1 {
		return fmt.Errorf("max idle connections must be at least 1")
	}
	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("max idle connections cannot exceed max open connections")
	}
	return nil
}
