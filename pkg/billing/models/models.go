// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// JSONB is a custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Organization represents a billing account
type Organization struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	BillingEmail string     `gorm:"type:varchar(255);not null" json:"billing_email"`
	CompanyName  string     `gorm:"type:varchar(255)" json:"company_name,omitempty"`
	TaxID        string     `gorm:"type:varchar(100)" json:"tax_id,omitempty"`

	// Address
	AddressLine1 string `gorm:"type:varchar(255)" json:"address_line1,omitempty"`
	AddressLine2 string `gorm:"type:varchar(255)" json:"address_line2,omitempty"`
	City         string `gorm:"type:varchar(100)" json:"city,omitempty"`
	State        string `gorm:"type:varchar(100)" json:"state,omitempty"`
	PostalCode   string `gorm:"type:varchar(20)" json:"postal_code,omitempty"`
	Country      string `gorm:"type:varchar(2)" json:"country,omitempty"` // ISO 3166-1 alpha-2

	// Billing settings
	Currency          string `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	BillingCycle      string `gorm:"type:varchar(20);default:'monthly'" json:"billing_cycle"`
	BillingDayOfMonth int    `gorm:"default:1" json:"billing_day_of_month"`
	Timezone          string `gorm:"type:varchar(50);default:'UTC'" json:"timezone"`

	// Payment
	DefaultPaymentMethodID string `gorm:"type:varchar(255)" json:"default_payment_method_id,omitempty"`
	StripeCustomerID       string `gorm:"type:varchar(255)" json:"stripe_customer_id,omitempty"`
	AutoPay                bool   `gorm:"default:false" json:"auto_pay"`

	// Status
	Status string `gorm:"type:varchar(20);default:'active'" json:"status"`

	// Audit
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the default table name
func (Organization) TableName() string {
	return "dictamesh_billing_organizations"
}

// SubscriptionPlan represents a product offering
type SubscriptionPlan struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description,omitempty"`

	// Pricing
	BasePrice       decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"base_price"`
	Currency        string          `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	BillingInterval string          `gorm:"type:varchar(20);not null" json:"billing_interval"`

	// Features
	Features JSONB `gorm:"type:jsonb;default:'{}'" json:"features,omitempty"`

	// Limits
	IncludedAPICalls      int `gorm:"default:0" json:"included_api_calls"`
	IncludedStorageGB     int `gorm:"default:0" json:"included_storage_gb"`
	IncludedDataTransferGB int `gorm:"default:0" json:"included_data_transfer_gb"`
	IncludedSeats         int `gorm:"default:1" json:"included_seats"`
	MaxAdapters           int `gorm:"default:0" json:"max_adapters"`

	// Add-on pricing
	PricePerAPICall       decimal.Decimal `gorm:"type:decimal(12,6);default:0" json:"price_per_api_call"`
	PricePerGBStorage     decimal.Decimal `gorm:"type:decimal(12,4);default:0" json:"price_per_gb_storage"`
	PricePerGBTransfer    decimal.Decimal `gorm:"type:decimal(12,4);default:0" json:"price_per_gb_transfer"`
	PricePerAdditionalSeat decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"price_per_additional_seat"`

	// Status
	IsPublic bool `gorm:"default:true" json:"is_public"`
	IsActive bool `gorm:"default:true" json:"is_active"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (SubscriptionPlan) TableName() string {
	return "dictamesh_billing_subscription_plans"
}

// Subscription represents an active subscription
type Subscription struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	PlanID         uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`

	// Relationships
	Organization Organization     `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Plan         SubscriptionPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`

	// Subscription details
	Status             string    `gorm:"type:varchar(20);default:'active';index" json:"status"`
	CurrentPeriodStart time.Time `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd   time.Time `gorm:"not null;index" json:"current_period_end"`

	// Trial
	TrialStart *time.Time `json:"trial_start,omitempty"`
	TrialEnd   *time.Time `json:"trial_end,omitempty"`

	// Cancellation
	CancelAtPeriodEnd  bool       `gorm:"default:false" json:"cancel_at_period_end"`
	CanceledAt         *time.Time `json:"canceled_at,omitempty"`
	CancellationReason string     `gorm:"type:text" json:"cancellation_reason,omitempty"`

	// Pricing overrides
	CustomPricing JSONB `gorm:"type:jsonb" json:"custom_pricing,omitempty"`

	// Seats
	Quantity int `gorm:"default:1" json:"quantity"`

	// Payment provider
	StripeSubscriptionID string `gorm:"type:varchar(255)" json:"stripe_subscription_id,omitempty"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (Subscription) TableName() string {
	return "dictamesh_billing_subscriptions"
}

// UsageMetric represents a usage measurement
type UsageMetric struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_usage_org_time" json:"organization_id"`
	SubscriptionID uuid.UUID `gorm:"type:uuid;index" json:"subscription_id,omitempty"`

	// Metric details
	MetricType  string          `gorm:"type:varchar(50);not null;index:idx_usage_type_time" json:"metric_type"`
	MetricValue decimal.Decimal `gorm:"type:decimal(20,6);not null" json:"metric_value"`
	MetricUnit  string          `gorm:"type:varchar(20);not null" json:"metric_unit"`

	// Time dimension
	RecordedAt  time.Time `gorm:"not null;default:now();index:idx_usage_org_time,idx_usage_type_time" json:"recorded_at"`
	PeriodStart time.Time `gorm:"not null" json:"period_start"`
	PeriodEnd   time.Time `gorm:"not null" json:"period_end"`

	// Metadata
	ResourceID string `gorm:"type:varchar(255)" json:"resource_id,omitempty"`
	Metadata   JSONB  `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the default table name
func (UsageMetric) TableName() string {
	return "dictamesh_billing_usage_metrics"
}

// Invoice represents a billing invoice
type Invoice struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	SubscriptionID uuid.UUID `gorm:"type:uuid;index" json:"subscription_id,omitempty"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Subscription Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`

	// Invoice identification
	InvoiceNumber string `gorm:"type:varchar(50);not null;uniqueIndex" json:"invoice_number"`

	// Billing period
	PeriodStart time.Time `gorm:"not null" json:"period_start"`
	PeriodEnd   time.Time `gorm:"not null" json:"period_end"`

	// Amounts
	Subtotal    decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	TaxAmount   decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"tax_amount"`
	TotalAmount decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	AmountDue   decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount_due"`
	AmountPaid  decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"amount_paid"`
	Currency    string          `gorm:"type:varchar(3);default:'USD'" json:"currency"`

	// Status
	Status string `gorm:"type:varchar(20);default:'draft';index" json:"status"`

	// Dates
	InvoiceDate time.Time  `gorm:"not null;default:now()" json:"invoice_date"`
	DueDate     time.Time  `gorm:"not null;index" json:"due_date"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`

	// Payment provider
	StripeInvoiceID string `gorm:"type:varchar(255)" json:"stripe_invoice_id,omitempty"`

	// PDF
	PDFURL        string     `gorm:"type:text" json:"pdf_url,omitempty"`
	PDFGeneratedAt *time.Time `json:"pdf_generated_at,omitempty"`

	// Line items
	LineItems []InvoiceLineItem `gorm:"foreignKey:InvoiceID" json:"line_items,omitempty"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (Invoice) TableName() string {
	return "dictamesh_billing_invoices"
}

// InvoiceLineItem represents a line item on an invoice
type InvoiceLineItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	InvoiceID uuid.UUID `gorm:"type:uuid;not null;index" json:"invoice_id"`

	// Line item details
	Description string          `gorm:"type:text;not null" json:"description"`
	Quantity    decimal.Decimal `gorm:"type:decimal(20,6);not null" json:"quantity"`
	UnitPrice   decimal.Decimal `gorm:"type:decimal(12,6);not null" json:"unit_price"`
	Amount      decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`

	// Categorization
	ItemType   string `gorm:"type:varchar(50);not null;index" json:"item_type"`
	MetricType string `gorm:"type:varchar(50)" json:"metric_type,omitempty"`

	// Period (for usage items)
	PeriodStart *time.Time `json:"period_start,omitempty"`
	PeriodEnd   *time.Time `json:"period_end,omitempty"`

	// Metadata
	Metadata JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the default table name
func (InvoiceLineItem) TableName() string {
	return "dictamesh_billing_invoice_line_items"
}

// Payment represents a payment transaction
type Payment struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	InvoiceID      uuid.UUID `gorm:"type:uuid;index" json:"invoice_id,omitempty"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Invoice      Invoice      `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`

	// Payment details
	Amount   decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency string          `gorm:"type:varchar(3);default:'USD'" json:"currency"`

	// Status
	Status string `gorm:"type:varchar(20);default:'pending';index" json:"status"`

	// Payment method
	PaymentMethod   string `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	PaymentMethodID string `gorm:"type:varchar(255)" json:"payment_method_id,omitempty"`

	// Provider details
	Provider          string `gorm:"type:varchar(20);default:'stripe';index:idx_payment_provider" json:"provider"`
	ProviderPaymentID string `gorm:"type:varchar(255);index:idx_payment_provider" json:"provider_payment_id,omitempty"`
	ProviderCustomerID string `gorm:"type:varchar(255)" json:"provider_customer_id,omitempty"`

	// Timestamps
	AttemptedAt *time.Time `json:"attempted_at,omitempty"`
	SucceededAt *time.Time `json:"succeeded_at,omitempty"`
	FailedAt    *time.Time `json:"failed_at,omitempty"`
	RefundedAt  *time.Time `json:"refunded_at,omitempty"`

	// Error handling
	FailureCode    string `gorm:"type:varchar(50)" json:"failure_code,omitempty"`
	FailureMessage string `gorm:"type:text" json:"failure_message,omitempty"`

	// Metadata
	Metadata JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (Payment) TableName() string {
	return "dictamesh_billing_payments"
}

// PricingTier represents volume-based pricing
type PricingTier struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID uuid.UUID `gorm:"type:uuid;index:idx_tier_plan_metric" json:"plan_id,omitempty"`

	// Tier definition
	MetricType string           `gorm:"type:varchar(50);not null;index:idx_tier_plan_metric" json:"metric_type"`
	TierStart  decimal.Decimal  `gorm:"type:decimal(20,2);not null" json:"tier_start"`
	TierEnd    *decimal.Decimal `gorm:"type:decimal(20,2)" json:"tier_end,omitempty"` // NULL = infinity

	// Pricing
	PricePerUnit decimal.Decimal `gorm:"type:decimal(12,6);not null" json:"price_per_unit"`
	FlatFee      decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"flat_fee"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (PricingTier) TableName() string {
	return "dictamesh_billing_pricing_tiers"
}

// Credit represents account credits
type Credit struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`

	// Credit details
	Amount          decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string          `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	RemainingAmount decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"remaining_amount"`

	// Reason
	Reason      string `gorm:"type:varchar(100);not null" json:"reason"`
	Description string `gorm:"type:text" json:"description,omitempty"`

	// Validity
	ValidFrom  time.Time  `gorm:"not null;default:now()" json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`

	// Status
	Status string `gorm:"type:varchar(20);default:'active';index:idx_credit_status_validity" json:"status"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (Credit) TableName() string {
	return "dictamesh_billing_credits"
}

// AuditLog represents billing audit trail
type AuditLog struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Entity tracking
	EntityType string    `gorm:"type:varchar(50);not null;index:idx_audit_entity" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null;index:idx_audit_entity" json:"entity_id"`

	// Event details
	EventType string `gorm:"type:varchar(50);not null" json:"event_type"`
	EventData JSONB  `gorm:"type:jsonb;not null" json:"event_data"`

	// Actor
	ActorID   string `gorm:"type:varchar(255)" json:"actor_id,omitempty"`
	ActorType string `gorm:"type:varchar(20);default:'system'" json:"actor_type"`

	// Context
	IPAddress string `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent string `gorm:"type:text" json:"user_agent,omitempty"`

	// Timestamp
	OccurredAt time.Time `gorm:"not null;default:now();index" json:"occurred_at"`
}

// TableName overrides the default table name
func (AuditLog) TableName() string {
	return "dictamesh_billing_audit_log"
}
