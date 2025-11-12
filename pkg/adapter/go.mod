// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

module github.com/click2-run/dictamesh/pkg/adapter

go 1.21

require (
	github.com/click2-run/dictamesh/pkg/observability v0.0.0
	github.com/click2-run/dictamesh/pkg/events v0.0.0
	go.uber.org/zap v1.26.0
)

replace (
	github.com/click2-run/dictamesh/pkg/observability => ../observability
	github.com/click2-run/dictamesh/pkg/events => ../events
)
