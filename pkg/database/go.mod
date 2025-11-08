// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

module github.com/click2-run/dictamesh/pkg/database

go 1.21

require (
	github.com/jackc/pgx/v5 v5.5.1
	github.com/jackc/pgconn v1.14.1
	github.com/jackc/pgio v1.0.0
	github.com/jackc/pgtype v1.14.0
	github.com/golang-migrate/migrate/v4 v4.17.0
	github.com/pgvector/pgvector-go v0.1.1
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
	github.com/redis/go-redis/v9 v9.3.0
	go.uber.org/zap v1.26.0
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	github.com/prometheus/client_golang v1.18.0
)
