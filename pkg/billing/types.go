// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"time"

	"github.com/shopspring/decimal"
)

// BillingCycle represents the billing frequency
type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleAnnual  BillingCycle = "annual"
)

// SubscriptionStatus represents the current state of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive     SubscriptionStatus = "active"
	SubscriptionStatusCanceled   SubscriptionStatus = "canceled"
	SubscriptionStatusPastDue    SubscriptionStatus = "past_due"
	SubscriptionStatusTrialing   SubscriptionStatus = "trialing"
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete"
)

// InvoiceStatus represents the current state of an invoice
type InvoiceStatus string

const (
	InvoiceStatusDraft         InvoiceStatus = "draft"
	InvoiceStatusOpen          InvoiceStatus = "open"
	InvoiceStatusPaid          InvoiceStatus = "paid"
	InvoiceStatusVoid          InvoiceStatus = "void"
	InvoiceStatusUncollectible InvoiceStatus = "uncollectible"
)

// PaymentStatus represents the current state of a payment
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
	PaymentStatusCanceled PaymentStatus = "canceled"
)

// OrganizationStatus represents the current state of a billing organization
type OrganizationStatus string

const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusSuspended OrganizationStatus = "suspended"
	OrganizationStatusDeleted   OrganizationStatus = "deleted"
)

// MetricType represents different billable metrics
type MetricType string

const (
	MetricTypeAPICalls          MetricType = "api_calls"
	MetricTypeStorageGB         MetricType = "storage_gb"
	MetricTypeTransferGBIn      MetricType = "transfer_gb_in"
	MetricTypeTransferGBOut     MetricType = "transfer_gb_out"
	MetricTypeQuerySeconds      MetricType = "query_seconds"
	MetricTypeGraphQLOperations MetricType = "graphql_operations"
	MetricTypeKafkaEvents       MetricType = "kafka_events"
	MetricTypeAdaptersActive    MetricType = "adapters_active"
)

// LineItemType represents different types of invoice line items
type LineItemType string

const (
	LineItemTypeSubscriptionBase LineItemType = "subscription_base"
	LineItemTypeUsageAPICalls    LineItemType = "usage_api_calls"
	LineItemTypeUsageStorage     LineItemType = "usage_storage"
	LineItemTypeUsageTransfer    LineItemType = "usage_transfer"
	LineItemTypeAddonSeats       LineItemType = "addon_seats"
	LineItemTypeAddonSupport     LineItemType = "addon_support"
	LineItemTypeCredit           LineItemType = "credit"
	LineItemTypeTax              LineItemType = "tax"
	LineItemTypeDiscount         LineItemType = "discount"
)

// PaymentProvider represents payment processing providers
type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "stripe"
	PaymentProviderPayPal PaymentProvider = "paypal"
	PaymentProviderManual PaymentProvider = "manual"
)

// CreditStatus represents the current state of a credit
type CreditStatus string

const (
	CreditStatusActive    CreditStatus = "active"
	CreditStatusExhausted CreditStatus = "exhausted"
	CreditStatusExpired   CreditStatus = "expired"
	CreditStatusVoided    CreditStatus = "voided"
)

// Money represents a monetary amount with currency
type Money struct {
	Amount   decimal.Decimal
	Currency string
}

// UsageRecord represents a single usage metric record
type UsageRecord struct {
	OrganizationID string
	SubscriptionID string
	MetricType     MetricType
	MetricValue    decimal.Decimal
	MetricUnit     string
	RecordedAt     time.Time
	PeriodStart    time.Time
	PeriodEnd      time.Time
	ResourceID     string
	Metadata       map[string]interface{}
}

// PricingTier represents a volume-based pricing tier
type PricingTier struct {
	TierStart    decimal.Decimal // Inclusive lower bound
	TierEnd      *decimal.Decimal // Exclusive upper bound (nil = infinity)
	PricePerUnit decimal.Decimal
	FlatFee      decimal.Decimal
}

// InvoiceLineItem represents a single charge on an invoice
type InvoiceLineItem struct {
	Description string
	Quantity    decimal.Decimal
	UnitPrice   decimal.Decimal
	Amount      decimal.Decimal
	ItemType    LineItemType
	MetricType  MetricType
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	Metadata    map[string]interface{}
}

// PaymentMethod represents a stored payment method
type PaymentMethod struct {
	ID             string
	Type           string // card, ach, paypal, etc.
	Last4          string
	ExpiryMonth    int
	ExpiryYear     int
	Brand          string
	IsDefault      bool
	ProviderID     string // Stripe payment method ID
	OrganizationID string
}

// BillingConfig represents the billing system configuration
type BillingConfig struct {
	// Database connection
	DatabaseDSN string

	// Payment providers
	StripeAPIKey       string
	StripeWebhookSecret string
	PayPalClientID     string
	PayPalClientSecret string

	// Invoice settings
	InvoiceDueDays       int
	InvoiceNumberPrefix  string
	TaxRate              decimal.Decimal
	DefaultCurrency      string

	// Usage aggregation
	UsageAggregationInterval time.Duration
	UsageRetentionDays       int

	// Notifications
	NotificationServiceURL string

	// Feature flags
	EnableAutoPayment     bool
	EnableUsageMetrics    bool
	EnableTieredPricing   bool
	EnableMultiCurrency   bool
}

// UsageAggregation represents aggregated usage for a billing period
type UsageAggregation struct {
	OrganizationID string
	SubscriptionID string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	Metrics        map[MetricType]decimal.Decimal
}

// ChargeCalculation represents the result of pricing calculation
type ChargeCalculation struct {
	BaseCharge      decimal.Decimal
	UsageCharges    map[MetricType]decimal.Decimal
	AddonCharges    decimal.Decimal
	Subtotal        decimal.Decimal
	Credits         decimal.Decimal
	TaxAmount       decimal.Decimal
	Total           decimal.Decimal
	LineItems       []InvoiceLineItem
}

// SubscriptionChange represents a change to a subscription (upgrade/downgrade)
type SubscriptionChange struct {
	SubscriptionID string
	FromPlanID     string
	ToPlanID       string
	ChangeType     string // upgrade, downgrade
	Proration      decimal.Decimal
	EffectiveDate  time.Time
}

// WebhookEvent represents a payment provider webhook event
type WebhookEvent struct {
	Provider  PaymentProvider
	EventType string
	EventID   string
	Payload   map[string]interface{}
	Signature string
	ReceivedAt time.Time
}

// BillingReport represents various billing reports
type BillingReport struct {
	ReportType  string // mrr, revenue, churn, usage
	PeriodStart time.Time
	PeriodEnd   time.Time
	Data        map[string]interface{}
	GeneratedAt time.Time
}

// EventType represents billing event types for Kafka
type EventType string

const (
	EventSubscriptionCreated      EventType = "billing.subscription.created"
	EventSubscriptionUpdated      EventType = "billing.subscription.updated"
	EventSubscriptionCanceled     EventType = "billing.subscription.canceled"
	EventInvoiceCreated           EventType = "billing.invoice.created"
	EventInvoicePaid              EventType = "billing.invoice.paid"
	EventInvoiceOverdue           EventType = "billing.invoice.overdue"
	EventPaymentSucceeded         EventType = "billing.payment.succeeded"
	EventPaymentFailed            EventType = "billing.payment.failed"
	EventUsageThresholdReached    EventType = "billing.usage.threshold_reached"
	EventCreditApplied            EventType = "billing.credit.applied"
)
