// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Click2-Run/dictamesh/pkg/billing/models"
)

// EventBus defines the interface for publishing events
type EventBus interface {
	Publish(ctx context.Context, topic string, key string, value interface{}) error
}

// BillingEventPublisher publishes billing events to Kafka
type BillingEventPublisher struct {
	eventBus EventBus
}

// NewBillingEventPublisher creates a new event publisher
func NewBillingEventPublisher(eventBus EventBus) *BillingEventPublisher {
	return &BillingEventPublisher{
		eventBus: eventBus,
	}
}

// Event structures for different billing events

// SubscriptionCreatedEvent represents a subscription creation event
type SubscriptionCreatedEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	SubscriptionID string    `json:"subscription_id"`
	OrganizationID string    `json:"organization_id"`
	PlanID         string    `json:"plan_id"`
	PlanName       string    `json:"plan_name"`
	Status         string    `json:"status"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
}

// SubscriptionUpdatedEvent represents a subscription update event
type SubscriptionUpdatedEvent struct {
	EventID        string                 `json:"event_id"`
	EventType      string                 `json:"event_type"`
	OccurredAt     time.Time              `json:"occurred_at"`
	SubscriptionID string                 `json:"subscription_id"`
	OrganizationID string                 `json:"organization_id"`
	Changes        map[string]interface{} `json:"changes"`
}

// SubscriptionCanceledEvent represents a subscription cancellation event
type SubscriptionCanceledEvent struct {
	EventID            string    `json:"event_id"`
	EventType          string    `json:"event_type"`
	OccurredAt         time.Time `json:"occurred_at"`
	SubscriptionID     string    `json:"subscription_id"`
	OrganizationID     string    `json:"organization_id"`
	CancellationReason string    `json:"cancellation_reason"`
	CanceledAt         time.Time `json:"canceled_at"`
	EndDate            time.Time `json:"end_date"`
}

// InvoiceCreatedEvent represents an invoice creation event
type InvoiceCreatedEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	InvoiceID      string    `json:"invoice_id"`
	InvoiceNumber  string    `json:"invoice_number"`
	OrganizationID string    `json:"organization_id"`
	SubscriptionID string    `json:"subscription_id"`
	TotalAmount    string    `json:"total_amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	DueDate        time.Time `json:"due_date"`
}

// InvoicePaidEvent represents an invoice payment event
type InvoicePaidEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	InvoiceID      string    `json:"invoice_id"`
	InvoiceNumber  string    `json:"invoice_number"`
	OrganizationID string    `json:"organization_id"`
	PaymentID      string    `json:"payment_id"`
	AmountPaid     string    `json:"amount_paid"`
	Currency       string    `json:"currency"`
	PaidAt         time.Time `json:"paid_at"`
}

// InvoiceOverdueEvent represents an overdue invoice event
type InvoiceOverdueEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	InvoiceID      string    `json:"invoice_id"`
	InvoiceNumber  string    `json:"invoice_number"`
	OrganizationID string    `json:"organization_id"`
	AmountDue      string    `json:"amount_due"`
	Currency       string    `json:"currency"`
	DueDate        time.Time `json:"due_date"`
	DaysOverdue    int       `json:"days_overdue"`
}

// PaymentSucceededEvent represents a successful payment event
type PaymentSucceededEvent struct {
	EventID          string    `json:"event_id"`
	EventType        string    `json:"event_type"`
	OccurredAt       time.Time `json:"occurred_at"`
	PaymentID        string    `json:"payment_id"`
	OrganizationID   string    `json:"organization_id"`
	InvoiceID        string    `json:"invoice_id"`
	Amount           string    `json:"amount"`
	Currency         string    `json:"currency"`
	PaymentMethod    string    `json:"payment_method"`
	ProviderPaymentID string   `json:"provider_payment_id"`
	SucceededAt      time.Time `json:"succeeded_at"`
}

// PaymentFailedEvent represents a failed payment event
type PaymentFailedEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	PaymentID      string    `json:"payment_id"`
	OrganizationID string    `json:"organization_id"`
	InvoiceID      string    `json:"invoice_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	FailureCode    string    `json:"failure_code"`
	FailureMessage string    `json:"failure_message"`
	FailedAt       time.Time `json:"failed_at"`
}

// UsageThresholdReachedEvent represents a usage threshold event
type UsageThresholdReachedEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	OrganizationID string    `json:"organization_id"`
	MetricType     string    `json:"metric_type"`
	CurrentUsage   string    `json:"current_usage"`
	Threshold      string    `json:"threshold"`
	PercentUsed    int       `json:"percent_used"`
}

// CreditAppliedEvent represents a credit application event
type CreditAppliedEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	CreditID       string    `json:"credit_id"`
	OrganizationID string    `json:"organization_id"`
	InvoiceID      string    `json:"invoice_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	Reason         string    `json:"reason"`
}

// PublishSubscriptionCreated publishes a subscription created event
func (p *BillingEventPublisher) PublishSubscriptionCreated(
	ctx context.Context,
	subscription *models.Subscription,
) error {
	event := SubscriptionCreatedEvent{
		EventID:        generateEventID(),
		EventType:      string(EventSubscriptionCreated),
		OccurredAt:     time.Now(),
		SubscriptionID: subscription.ID.String(),
		OrganizationID: subscription.OrganizationID.String(),
		PlanID:         subscription.PlanID.String(),
		PlanName:       subscription.Plan.Name,
		Status:         subscription.Status,
		PeriodStart:    subscription.CurrentPeriodStart,
		PeriodEnd:      subscription.CurrentPeriodEnd,
		Amount:         subscription.Plan.BasePrice.String(),
		Currency:       subscription.Plan.Currency,
	}

	return p.publish(ctx, string(EventSubscriptionCreated), subscription.OrganizationID.String(), event)
}

// PublishSubscriptionCanceled publishes a subscription canceled event
func (p *BillingEventPublisher) PublishSubscriptionCanceled(
	ctx context.Context,
	subscription *models.Subscription,
) error {
	event := SubscriptionCanceledEvent{
		EventID:            generateEventID(),
		EventType:          string(EventSubscriptionCanceled),
		OccurredAt:         time.Now(),
		SubscriptionID:     subscription.ID.String(),
		OrganizationID:     subscription.OrganizationID.String(),
		CancellationReason: subscription.CancellationReason,
		CanceledAt:         *subscription.CanceledAt,
		EndDate:            subscription.CurrentPeriodEnd,
	}

	return p.publish(ctx, string(EventSubscriptionCanceled), subscription.OrganizationID.String(), event)
}

// PublishInvoiceCreated publishes an invoice created event
func (p *BillingEventPublisher) PublishInvoiceCreated(
	ctx context.Context,
	invoice *models.Invoice,
) error {
	event := InvoiceCreatedEvent{
		EventID:        generateEventID(),
		EventType:      string(EventInvoiceCreated),
		OccurredAt:     time.Now(),
		InvoiceID:      invoice.ID.String(),
		InvoiceNumber:  invoice.InvoiceNumber,
		OrganizationID: invoice.OrganizationID.String(),
		SubscriptionID: invoice.SubscriptionID.String(),
		TotalAmount:    invoice.TotalAmount.String(),
		Currency:       invoice.Currency,
		Status:         invoice.Status,
		DueDate:        invoice.DueDate,
	}

	return p.publish(ctx, string(EventInvoiceCreated), invoice.OrganizationID.String(), event)
}

// PublishInvoicePaid publishes an invoice paid event
func (p *BillingEventPublisher) PublishInvoicePaid(
	ctx context.Context,
	invoice *models.Invoice,
	paymentID string,
) error {
	event := InvoicePaidEvent{
		EventID:        generateEventID(),
		EventType:      string(EventInvoicePaid),
		OccurredAt:     time.Now(),
		InvoiceID:      invoice.ID.String(),
		InvoiceNumber:  invoice.InvoiceNumber,
		OrganizationID: invoice.OrganizationID.String(),
		PaymentID:      paymentID,
		AmountPaid:     invoice.AmountPaid.String(),
		Currency:       invoice.Currency,
		PaidAt:         *invoice.PaidAt,
	}

	return p.publish(ctx, string(EventInvoicePaid), invoice.OrganizationID.String(), event)
}

// PublishInvoiceOverdue publishes an invoice overdue event
func (p *BillingEventPublisher) PublishInvoiceOverdue(
	ctx context.Context,
	invoice *models.Invoice,
) error {
	daysOverdue := int(time.Since(invoice.DueDate).Hours() / 24)

	event := InvoiceOverdueEvent{
		EventID:        generateEventID(),
		EventType:      string(EventInvoiceOverdue),
		OccurredAt:     time.Now(),
		InvoiceID:      invoice.ID.String(),
		InvoiceNumber:  invoice.InvoiceNumber,
		OrganizationID: invoice.OrganizationID.String(),
		AmountDue:      invoice.AmountDue.String(),
		Currency:       invoice.Currency,
		DueDate:        invoice.DueDate,
		DaysOverdue:    daysOverdue,
	}

	return p.publish(ctx, string(EventInvoiceOverdue), invoice.OrganizationID.String(), event)
}

// PublishPaymentSucceeded publishes a payment succeeded event
func (p *BillingEventPublisher) PublishPaymentSucceeded(
	ctx context.Context,
	payment *models.Payment,
) error {
	event := PaymentSucceededEvent{
		EventID:           generateEventID(),
		EventType:         string(EventPaymentSucceeded),
		OccurredAt:        time.Now(),
		PaymentID:         payment.ID.String(),
		OrganizationID:    payment.OrganizationID.String(),
		InvoiceID:         payment.InvoiceID.String(),
		Amount:            payment.Amount.String(),
		Currency:          payment.Currency,
		PaymentMethod:     payment.PaymentMethod,
		ProviderPaymentID: payment.ProviderPaymentID,
		SucceededAt:       *payment.SucceededAt,
	}

	return p.publish(ctx, string(EventPaymentSucceeded), payment.OrganizationID.String(), event)
}

// PublishPaymentFailed publishes a payment failed event
func (p *BillingEventPublisher) PublishPaymentFailed(
	ctx context.Context,
	payment *models.Payment,
) error {
	event := PaymentFailedEvent{
		EventID:        generateEventID(),
		EventType:      string(EventPaymentFailed),
		OccurredAt:     time.Now(),
		PaymentID:      payment.ID.String(),
		OrganizationID: payment.OrganizationID.String(),
		InvoiceID:      payment.InvoiceID.String(),
		Amount:         payment.Amount.String(),
		Currency:       payment.Currency,
		FailureCode:    payment.FailureCode,
		FailureMessage: payment.FailureMessage,
		FailedAt:       *payment.FailedAt,
	}

	return p.publish(ctx, string(EventPaymentFailed), payment.OrganizationID.String(), event)
}

// PublishUsageThresholdReached publishes a usage threshold reached event
func (p *BillingEventPublisher) PublishUsageThresholdReached(
	ctx context.Context,
	organizationID string,
	metricType MetricType,
	currentUsage, threshold string,
	percentUsed int,
) error {
	event := UsageThresholdReachedEvent{
		EventID:        generateEventID(),
		EventType:      string(EventUsageThresholdReached),
		OccurredAt:     time.Now(),
		OrganizationID: organizationID,
		MetricType:     string(metricType),
		CurrentUsage:   currentUsage,
		Threshold:      threshold,
		PercentUsed:    percentUsed,
	}

	return p.publish(ctx, string(EventUsageThresholdReached), organizationID, event)
}

// PublishCreditApplied publishes a credit applied event
func (p *BillingEventPublisher) PublishCreditApplied(
	ctx context.Context,
	credit *models.Credit,
	invoiceID string,
	amountApplied string,
) error {
	event := CreditAppliedEvent{
		EventID:        generateEventID(),
		EventType:      string(EventCreditApplied),
		OccurredAt:     time.Now(),
		CreditID:       credit.ID.String(),
		OrganizationID: credit.OrganizationID.String(),
		InvoiceID:      invoiceID,
		Amount:         amountApplied,
		Currency:       credit.Currency,
		Reason:         credit.Reason,
	}

	return p.publish(ctx, string(EventCreditApplied), credit.OrganizationID.String(), event)
}

// publish publishes an event to Kafka
func (p *BillingEventPublisher) publish(ctx context.Context, topic string, key string, event interface{}) error {
	if p.eventBus == nil {
		return fmt.Errorf("event bus not configured")
	}

	// Validate event
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Validate JSON
	if !json.Valid(eventBytes) {
		return fmt.Errorf("invalid JSON event")
	}

	// Publish to event bus
	return p.eventBus.Publish(ctx, topic, key, event)
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
