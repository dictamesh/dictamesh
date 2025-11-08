// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer = otel.Tracer("billing")
)

// Prometheus Metrics

var (
	// Subscription metrics
	activeSubscriptionsGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dictamesh_billing_active_subscriptions",
			Help: "Number of active subscriptions by plan",
		},
		[]string{"plan"},
	)

	subscriptionCreatedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_subscriptions_created_total",
			Help: "Total subscriptions created",
		},
		[]string{"plan"},
	)

	subscriptionCanceledCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_subscriptions_canceled_total",
			Help: "Total subscriptions canceled",
		},
		[]string{"plan", "reason"},
	)

	// Revenue metrics
	monthlyRecurringRevenueGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "dictamesh_billing_mrr",
			Help: "Monthly recurring revenue in USD",
		},
	)

	annualRecurringRevenueGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "dictamesh_billing_arr",
			Help: "Annual recurring revenue in USD",
		},
	)

	// Invoice metrics
	invoicesGeneratedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_invoices_generated_total",
			Help: "Total invoices generated",
		},
		[]string{"status"},
	)

	invoiceAmountHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dictamesh_billing_invoice_amount",
			Help:    "Invoice amounts in USD",
			Buckets: prometheus.ExponentialBuckets(1, 2, 15),
		},
		[]string{"currency"},
	)

	// Payment metrics
	paymentsProcessedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_payments_processed_total",
			Help: "Total payments processed",
		},
		[]string{"status", "provider"},
	)

	paymentAmountHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dictamesh_billing_payment_amount",
			Help:    "Payment amounts",
			Buckets: prometheus.ExponentialBuckets(1, 2, 15),
		},
		[]string{"currency", "status"},
	)

	paymentFailuresCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_payment_failures_total",
			Help: "Total payment failures",
		},
		[]string{"failure_code", "provider"},
	)

	paymentProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dictamesh_billing_payment_processing_duration_seconds",
			Help:    "Payment processing duration",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
		},
		[]string{"provider"},
	)

	// Usage metrics
	usageMetricsCollectedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_usage_metrics_collected_total",
			Help: "Total usage metrics collected",
		},
		[]string{"metric_type"},
	)

	// Credit metrics
	creditsIssuedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_credits_issued_total",
			Help: "Total credits issued",
		},
		[]string{"reason"},
	)

	creditsAppliedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dictamesh_billing_credits_applied_total",
			Help: "Total credits applied to invoices",
		},
		[]string{"reason"},
	)
)

// ObservabilityService provides observability instrumentation
type ObservabilityService struct{}

// NewObservabilityService creates a new observability service
func NewObservabilityService() *ObservabilityService {
	return &ObservabilityService{}
}

// Subscription Metrics

// RecordSubscriptionCreated records a subscription creation
func (o *ObservabilityService) RecordSubscriptionCreated(plan string) {
	subscriptionCreatedCounter.WithLabelValues(plan).Inc()
	// Update active subscriptions gauge would require querying DB
}

// RecordSubscriptionCanceled records a subscription cancellation
func (o *ObservabilityService) RecordSubscriptionCanceled(plan, reason string) {
	subscriptionCanceledCounter.WithLabelValues(plan, reason).Inc()
}

// UpdateMRR updates the monthly recurring revenue metric
func (o *ObservabilityService) UpdateMRR(mrr float64) {
	monthlyRecurringRevenueGauge.Set(mrr)
}

// UpdateARR updates the annual recurring revenue metric
func (o *ObservabilityService) UpdateARR(arr float64) {
	annualRecurringRevenueGauge.Set(arr)
}

// Invoice Metrics

// RecordInvoiceGenerated records an invoice generation
func (o *ObservabilityService) RecordInvoiceGenerated(status string, amount float64, currency string) {
	invoicesGeneratedCounter.WithLabelValues(status).Inc()
	invoiceAmountHistogram.WithLabelValues(currency).Observe(amount)
}

// Payment Metrics

// RecordPaymentProcessed records a payment attempt
func (o *ObservabilityService) RecordPaymentProcessed(status, provider string, amount float64, currency string) {
	paymentsProcessedCounter.WithLabelValues(status, provider).Inc()
	paymentAmountHistogram.WithLabelValues(currency, status).Observe(amount)
}

// RecordPaymentFailure records a payment failure
func (o *ObservabilityService) RecordPaymentFailure(failureCode, provider string) {
	paymentFailuresCounter.WithLabelValues(failureCode, provider).Inc()
}

// RecordPaymentDuration records payment processing duration
func (o *ObservabilityService) RecordPaymentDuration(provider string, seconds float64) {
	paymentProcessingDuration.WithLabelValues(provider).Observe(seconds)
}

// Usage Metrics

// RecordUsageMetricCollected records a usage metric collection
func (o *ObservabilityService) RecordUsageMetricCollected(metricType string) {
	usageMetricsCollectedCounter.WithLabelValues(metricType).Inc()
}

// Credit Metrics

// RecordCreditIssued records a credit issuance
func (o *ObservabilityService) RecordCreditIssued(reason string) {
	creditsIssuedCounter.WithLabelValues(reason).Inc()
}

// RecordCreditApplied records a credit application
func (o *ObservabilityService) RecordCreditApplied(reason string) {
	creditsAppliedCounter.WithLabelValues(reason).Inc()
}

// OpenTelemetry Tracing Helpers

// TraceInvoiceGeneration creates a span for invoice generation
func TraceInvoiceGeneration(ctx context.Context, subscriptionID string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "billing.generate_invoice")
	span.SetAttributes(
		attribute.String("subscription.id", subscriptionID),
	)
	return ctx, span
}

// TracePaymentProcessing creates a span for payment processing
func TracePaymentProcessing(ctx context.Context, invoiceID, provider string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "billing.process_payment")
	span.SetAttributes(
		attribute.String("invoice.id", invoiceID),
		attribute.String("payment.provider", provider),
	)
	return ctx, span
}

// TraceUsageCollection creates a span for usage collection
func TraceUsageCollection(ctx context.Context, organizationID string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "billing.collect_usage")
	span.SetAttributes(
		attribute.String("organization.id", organizationID),
	)
	return ctx, span
}

// TracePricingCalculation creates a span for pricing calculation
func TracePricingCalculation(ctx context.Context, subscriptionID string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "billing.calculate_pricing")
	span.SetAttributes(
		attribute.String("subscription.id", subscriptionID),
	)
	return ctx, span
}

// TraceNotificationSend creates a span for notification sending
func TraceNotificationSend(ctx context.Context, notificationType, recipientID string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "billing.send_notification")
	span.SetAttributes(
		attribute.String("notification.type", notificationType),
		attribute.String("recipient.id", recipientID),
	)
	return ctx, span
}

// RecordSpanError records an error in the current span
func RecordSpanError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
	}
}

// RecordSpanSuccess records success in the current span
func RecordSpanSuccess(span trace.Span) {
	span.SetAttributes(attribute.Bool("success", true))
}
