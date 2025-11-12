// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package adapter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/click2-run/dictamesh/pkg/events"
	"github.com/click2-run/dictamesh/pkg/observability"
)

// BaseAdapter provides common functionality for all adapters
type BaseAdapter struct {
	config       *Config
	obs          *observability.Observability
	producer     *events.Producer
	status       AdapterStatus
	statusMu     sync.RWMutex
	metrics      AdapterMetrics
	metricsMu    sync.RWMutex
	capabilities AdapterCapabilities
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(config *Config, obs *observability.Observability) (*BaseAdapter, error) {
	// Create event producer if events are enabled
	var producer *events.Producer
	if config.EnableEvents {
		var err error
		producer, err = events.NewProducer(config.EventsConfig, obs.Logger())
		if err != nil {
			return nil, fmt.Errorf("failed to create event producer: %w", err)
		}
	}

	return &BaseAdapter{
		config:   config,
		obs:      obs,
		producer: producer,
		status:   StatusUninitialized,
		capabilities: AdapterCapabilities{
			SupportsRead:   true,
			SupportsWrite:  false,
			SupportsDelete: false,
			SupportsCache:  config.EnableCache,
			SupportsEvents: config.EnableEvents,
		},
	}, nil
}

// Initialize performs common initialization
func (b *BaseAdapter) Initialize(ctx context.Context) error {
	b.setStatus(StatusInitializing)

	b.obs.Logger().InfoContext(ctx, "initializing adapter",
		"adapter", b.config.Name,
		"source", b.config.SourceSystem,
	)

	// Register health check
	b.obs.Health().RegisterReadinessCheck(b.config.Name, func(ctx context.Context) error {
		return b.Health(ctx)
	})

	b.setStatus(StatusReady)
	return nil
}

// Start starts the adapter
func (b *BaseAdapter) Start(ctx context.Context) error {
	b.obs.Logger().InfoContext(ctx, "starting adapter",
		"adapter", b.config.Name,
	)

	// Publish adapter started event
	if b.producer != nil {
		event := events.NewEvent(
			events.EventTypeAdapterStarted,
			b.config.Name,
			"adapter:"+b.config.Name,
			map[string]interface{}{
				"adapter": b.config.Name,
				"version": b.config.Version,
			},
		)
		_ = b.producer.Publish(ctx, events.TopicSystemEvents, event)
	}

	return nil
}

// Stop stops the adapter
func (b *BaseAdapter) Stop(ctx context.Context) error {
	b.obs.Logger().InfoContext(ctx, "stopping adapter",
		"adapter", b.config.Name,
	)

	// Publish adapter stopped event
	if b.producer != nil {
		event := events.NewEvent(
			events.EventTypeAdapterStopped,
			b.config.Name,
			"adapter:"+b.config.Name,
			map[string]interface{}{
				"adapter": b.config.Name,
			},
		)
		_ = b.producer.Publish(ctx, events.TopicSystemEvents, event)

		// Close producer
		b.producer.Close()
	}

	b.setStatus(StatusStopped)
	return nil
}

// Health performs a health check
func (b *BaseAdapter) Health(ctx context.Context) error {
	status := b.GetStatus()
	if status != StatusReady && status != StatusDegraded {
		return fmt.Errorf("adapter not ready: status=%s", status)
	}
	return nil
}

// Metadata methods
func (b *BaseAdapter) Name() string           { return b.config.Name }
func (b *BaseAdapter) Version() string        { return b.config.Version }
func (b *BaseAdapter) Description() string    { return b.config.Description }
func (b *BaseAdapter) SourceSystem() string   { return b.config.SourceSystem }
func (b *BaseAdapter) Domain() string         { return b.config.Domain }

// GetStatus returns the current status
func (b *BaseAdapter) GetStatus() AdapterStatus {
	b.statusMu.RLock()
	defer b.statusMu.RUnlock()
	return b.status
}

// setStatus sets the adapter status
func (b *BaseAdapter) setStatus(status AdapterStatus) {
	b.statusMu.Lock()
	defer b.statusMu.Unlock()
	b.status = status
}

// GetMetrics returns current metrics
func (b *BaseAdapter) GetMetrics() AdapterMetrics {
	b.metricsMu.RLock()
	defer b.metricsMu.RUnlock()
	return b.metrics
}

// IncrementRequests increments request counters
func (b *BaseAdapter) IncrementRequests(success bool) {
	b.metricsMu.Lock()
	defer b.metricsMu.Unlock()

	b.metrics.RequestsTotal++
	if success {
		b.metrics.RequestsSucceeded++
	} else {
		b.metrics.RequestsFailed++
	}
}

// RecordLatency records request latency
func (b *BaseAdapter) RecordLatency(duration time.Duration) {
	b.metricsMu.Lock()
	defer b.metricsMu.Unlock()

	// Simple moving average
	if b.metrics.AvgLatencyMs == 0 {
		b.metrics.AvgLatencyMs = float64(duration.Milliseconds())
	} else {
		b.metrics.AvgLatencyMs = (b.metrics.AvgLatencyMs*0.9 + float64(duration.Milliseconds())*0.1)
	}
}

// RecordError records an error
func (b *BaseAdapter) RecordError(err error) {
	b.metricsMu.Lock()
	defer b.metricsMu.Unlock()

	b.metrics.LastError = err
	b.metrics.LastErrorTime = time.Now()
}

// IncrementCacheHit increments cache hit counter
func (b *BaseAdapter) IncrementCacheHit() {
	b.metricsMu.Lock()
	defer b.metricsMu.Unlock()
	b.metrics.CacheHits++
}

// IncrementCacheMiss increments cache miss counter
func (b *BaseAdapter) IncrementCacheMiss() {
	b.metricsMu.Lock()
	defer b.metricsMu.Unlock()
	b.metrics.CacheMisses++
}

// PublishEvent publishes an event
func (b *BaseAdapter) PublishEvent(ctx context.Context, event *events.Event) error {
	if b.producer == nil {
		return fmt.Errorf("events not enabled for adapter")
	}
	return b.producer.Publish(ctx, events.TopicEntityChanged, event)
}

// GetObservability returns the observability instance
func (b *BaseAdapter) GetObservability() *observability.Observability {
	return b.obs
}

// GetCapabilities returns adapter capabilities
func (b *BaseAdapter) GetCapabilities() AdapterCapabilities {
	return b.capabilities
}

// SetCapabilities sets adapter capabilities
func (b *BaseAdapter) SetCapabilities(caps AdapterCapabilities) {
	b.capabilities = caps
}

// WithSpan is a helper to wrap operations with tracing
func (b *BaseAdapter) WithSpan(ctx context.Context, operation string, fn func(context.Context) error) error {
	start := time.Now()
	ctx, span := b.obs.StartSpan(ctx, "adapter."+operation)
	defer span.End()

	// Add attributes
	observability.SetAttributes(ctx,
		observability.AttrAdapterName.String(b.config.Name),
		observability.AttrSourceSystem.String(b.config.SourceSystem),
	)

	err := fn(ctx)

	// Record metrics
	b.RecordLatency(time.Since(start))
	b.IncrementRequests(err == nil)

	if err != nil {
		observability.RecordError(ctx, err)
		b.RecordError(err)
		b.obs.LoggerWithContext(ctx).Error("adapter operation failed",
			"operation", operation,
			"error", err,
		)
	}

	return err
}
