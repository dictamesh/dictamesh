// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/Click2-Run/dictamesh/pkg/billing/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// MetricsCollector handles usage metrics collection and aggregation
type MetricsCollector struct {
	db     *gorm.DB
	config *Config

	// Prometheus metrics
	apiCallsTotal      *prometheus.CounterVec
	storageBytes       *prometheus.GaugeVec
	transferBytesTotal *prometheus.CounterVec
	queryDuration      *prometheus.HistogramVec
	activeAdapters     *prometheus.GaugeVec
	kafkaEventsTotal   *prometheus.CounterVec
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(db *gorm.DB, config *Config) *MetricsCollector {
	return &MetricsCollector{
		db:     db,
		config: config,

		apiCallsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dictamesh_billing_api_calls_total",
				Help: "Total API calls by organization and endpoint",
			},
			[]string{"organization_id", "endpoint", "method"},
		),

		storageBytes: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "dictamesh_billing_storage_bytes",
				Help: "Current storage usage in bytes by organization",
			},
			[]string{"organization_id", "storage_type"},
		),

		transferBytesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dictamesh_billing_transfer_bytes_total",
				Help: "Total data transfer in bytes by organization",
			},
			[]string{"organization_id", "direction"},
		),

		queryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "dictamesh_billing_query_duration_seconds",
				Help:    "Query processing duration by organization",
				Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
			},
			[]string{"organization_id", "query_type"},
		),

		activeAdapters: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "dictamesh_billing_active_adapters",
				Help: "Number of active adapters by organization",
			},
			[]string{"organization_id"},
		),

		kafkaEventsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "dictamesh_billing_kafka_events_total",
				Help: "Total Kafka events by organization",
			},
			[]string{"organization_id", "topic"},
		),
	}
}

// RecordAPICall records an API call metric
func (mc *MetricsCollector) RecordAPICall(organizationID, endpoint, method string) {
	mc.apiCallsTotal.WithLabelValues(organizationID, endpoint, method).Inc()
}

// RecordStorage records current storage usage
func (mc *MetricsCollector) RecordStorage(organizationID, storageType string, bytes int64) {
	mc.storageBytes.WithLabelValues(organizationID, storageType).Set(float64(bytes))
}

// RecordTransfer records data transfer
func (mc *MetricsCollector) RecordTransfer(organizationID, direction string, bytes int64) {
	mc.transferBytesTotal.WithLabelValues(organizationID, direction).Add(float64(bytes))
}

// RecordQuery records a query execution
func (mc *MetricsCollector) RecordQuery(organizationID, queryType string, duration time.Duration) {
	mc.queryDuration.WithLabelValues(organizationID, queryType).Observe(duration.Seconds())
}

// RecordActiveAdapters records the number of active adapters
func (mc *MetricsCollector) RecordActiveAdapters(organizationID string, count int) {
	mc.activeAdapters.WithLabelValues(organizationID).Set(float64(count))
}

// RecordKafkaEvent records a Kafka event
func (mc *MetricsCollector) RecordKafkaEvent(organizationID, topic string) {
	mc.kafkaEventsTotal.WithLabelValues(organizationID, topic).Inc()
}

// AggregateUsageMetrics aggregates Prometheus metrics into database records
func (mc *MetricsCollector) AggregateUsageMetrics(ctx context.Context) error {
	now := time.Now()
	periodStart := now.Add(-mc.config.Usage.AggregationInterval)
	periodEnd := now

	// Get all organizations with subscriptions
	var subscriptions []models.Subscription
	if err := mc.db.WithContext(ctx).
		Preload("Organization").
		Where("status = ?", SubscriptionStatusActive).
		Find(&subscriptions).Error; err != nil {
		return fmt.Errorf("failed to fetch subscriptions: %w", err)
	}

	for _, sub := range subscriptions {
		orgID := sub.OrganizationID.String()

		// Aggregate API calls
		if err := mc.aggregateAPICallMetrics(ctx, orgID, sub.ID, periodStart, periodEnd); err != nil {
			return fmt.Errorf("failed to aggregate API calls for org %s: %w", orgID, err)
		}

		// Aggregate storage
		if err := mc.aggregateStorageMetrics(ctx, orgID, sub.ID, periodStart, periodEnd); err != nil {
			return fmt.Errorf("failed to aggregate storage for org %s: %w", orgID, err)
		}

		// Aggregate data transfer
		if err := mc.aggregateTransferMetrics(ctx, orgID, sub.ID, periodStart, periodEnd); err != nil {
			return fmt.Errorf("failed to aggregate transfer for org %s: %w", orgID, err)
		}

		// Aggregate query duration
		if err := mc.aggregateQueryMetrics(ctx, orgID, sub.ID, periodStart, periodEnd); err != nil {
			return fmt.Errorf("failed to aggregate queries for org %s: %w", orgID, err)
		}
	}

	return nil
}

// aggregateAPICallMetrics aggregates API call metrics
func (mc *MetricsCollector) aggregateAPICallMetrics(
	ctx context.Context,
	organizationID string,
	subscriptionID interface{},
	periodStart, periodEnd time.Time,
) error {
	// In a real implementation, you would query Prometheus for the metric values
	// For now, we'll simulate with a direct counter read
	// This is a simplified example - in production, you'd use the Prometheus API

	metric := &models.UsageMetric{
		OrganizationID: mustParseUUID(organizationID),
		SubscriptionID: subscriptionID.(interface{ String() string }).String(),
		MetricType:     string(MetricTypeAPICalls),
		MetricValue:    decimal.NewFromInt(0), // Would be fetched from Prometheus
		MetricUnit:     "count",
		RecordedAt:     time.Now(),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}

	return mc.db.WithContext(ctx).Create(metric).Error
}

// aggregateStorageMetrics aggregates storage metrics
func (mc *MetricsCollector) aggregateStorageMetrics(
	ctx context.Context,
	organizationID string,
	subscriptionID interface{},
	periodStart, periodEnd time.Time,
) error {
	metric := &models.UsageMetric{
		OrganizationID: mustParseUUID(organizationID),
		SubscriptionID: subscriptionID.(interface{ String() string }).String(),
		MetricType:     string(MetricTypeStorageGB),
		MetricValue:    decimal.NewFromInt(0), // Would be fetched from Prometheus
		MetricUnit:     "GB",
		RecordedAt:     time.Now(),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}

	return mc.db.WithContext(ctx).Create(metric).Error
}

// aggregateTransferMetrics aggregates data transfer metrics
func (mc *MetricsCollector) aggregateTransferMetrics(
	ctx context.Context,
	organizationID string,
	subscriptionID interface{},
	periodStart, periodEnd time.Time,
) error {
	// Aggregate inbound transfer
	metricIn := &models.UsageMetric{
		OrganizationID: mustParseUUID(organizationID),
		SubscriptionID: subscriptionID.(interface{ String() string }).String(),
		MetricType:     string(MetricTypeTransferGBIn),
		MetricValue:    decimal.NewFromInt(0), // Would be fetched from Prometheus
		MetricUnit:     "GB",
		RecordedAt:     time.Now(),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}

	if err := mc.db.WithContext(ctx).Create(metricIn).Error; err != nil {
		return err
	}

	// Aggregate outbound transfer
	metricOut := &models.UsageMetric{
		OrganizationID: mustParseUUID(organizationID),
		SubscriptionID: subscriptionID.(interface{ String() string }).String(),
		MetricType:     string(MetricTypeTransferGBOut),
		MetricValue:    decimal.NewFromInt(0), // Would be fetched from Prometheus
		MetricUnit:     "GB",
		RecordedAt:     time.Now(),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}

	return mc.db.WithContext(ctx).Create(metricOut).Error
}

// aggregateQueryMetrics aggregates query processing metrics
func (mc *MetricsCollector) aggregateQueryMetrics(
	ctx context.Context,
	organizationID string,
	subscriptionID interface{},
	periodStart, periodEnd time.Time,
) error {
	metric := &models.UsageMetric{
		OrganizationID: mustParseUUID(organizationID),
		SubscriptionID: subscriptionID.(interface{ String() string }).String(),
		MetricType:     string(MetricTypeQuerySeconds),
		MetricValue:    decimal.NewFromInt(0), // Would be fetched from Prometheus
		MetricUnit:     "seconds",
		RecordedAt:     time.Now(),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}

	return mc.db.WithContext(ctx).Create(metric).Error
}

// GetUsageForPeriod retrieves aggregated usage for a billing period
func (mc *MetricsCollector) GetUsageForPeriod(
	ctx context.Context,
	organizationID string,
	periodStart, periodEnd time.Time,
) (*UsageAggregation, error) {
	var metrics []models.UsageMetric

	err := mc.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Where("period_start >= ?", periodStart).
		Where("period_end <= ?", periodEnd).
		Find(&metrics).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage metrics: %w", err)
	}

	agg := &UsageAggregation{
		OrganizationID: organizationID,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		Metrics:        make(map[MetricType]decimal.Decimal),
	}

	// Aggregate metrics by type
	for _, metric := range metrics {
		metricType := MetricType(metric.MetricType)
		if current, ok := agg.Metrics[metricType]; ok {
			agg.Metrics[metricType] = current.Add(metric.MetricValue)
		} else {
			agg.Metrics[metricType] = metric.MetricValue
		}
	}

	return agg, nil
}

// GetCurrentUsage retrieves current usage (real-time)
func (mc *MetricsCollector) GetCurrentUsage(
	ctx context.Context,
	organizationID string,
) (map[MetricType]decimal.Decimal, error) {
	// This would query Prometheus for real-time metrics
	// For now, we'll return the most recent aggregated values

	var metrics []models.UsageMetric

	err := mc.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Where("recorded_at >= ?", time.Now().Add(-1*time.Hour)).
		Find(&metrics).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch current usage: %w", err)
	}

	usage := make(map[MetricType]decimal.Decimal)
	for _, metric := range metrics {
		metricType := MetricType(metric.MetricType)
		if current, ok := usage[metricType]; ok {
			usage[metricType] = current.Add(metric.MetricValue)
		} else {
			usage[metricType] = metric.MetricValue
		}
	}

	return usage, nil
}

// StartAggregationWorker starts a background worker to aggregate metrics periodically
func (mc *MetricsCollector) StartAggregationWorker(ctx context.Context) {
	ticker := time.NewTicker(mc.config.Usage.AggregationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := mc.AggregateUsageMetrics(ctx); err != nil {
				// Log error (in production, use proper logging)
				fmt.Printf("Error aggregating metrics: %v\n", err)
			}
		}
	}
}

// Helper functions

func mustParseUUID(s string) interface{} {
	// In real implementation, properly parse UUID
	// This is simplified for the example
	return s
}
