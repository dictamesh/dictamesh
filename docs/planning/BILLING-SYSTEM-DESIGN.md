# DictaMesh Billing System - Comprehensive Design Document

## Table of Contents
1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Database Schema](#database-schema)
4. [Pricing Model](#pricing-model)
5. [Usage Metrics Collection](#usage-metrics-collection)
6. [Invoice Generation](#invoice-generation)
7. [Payment Processing](#payment-processing)
8. [Notification Integration](#notification-integration)
9. [API Design](#api-design)
10. [Event-Driven Architecture](#event-driven-architecture)
11. [Observability](#observability)
12. [Security Considerations](#security-considerations)
13. [Implementation Roadmap](#implementation-roadmap)

---

## Overview

### Objectives
The DictaMesh Billing System provides a comprehensive, enterprise-grade billing solution for:
- **Hosted Services**: Data mesh adapter hosting, query processing, storage
- **Support Services**: Premium support tiers, consulting, SLA guarantees
- **Usage-Based Billing**: Granular tracking of API calls, data transfer, storage, query complexity
- **Subscription Management**: Multiple tiers, add-ons, seat-based pricing
- **Automated Invoicing**: Generation, delivery, and payment tracking
- **Multi-Currency Support**: Global billing capabilities

### Key Features
- ✅ Real-time usage metrics collection
- ✅ Fractional pricing calculations (per-second, per-MB, per-API-call)
- ✅ Flexible pricing models (flat-rate, usage-based, hybrid, tiered)
- ✅ Automated invoice generation and delivery
- ✅ Multiple payment provider support (Stripe, PayPal)
- ✅ Integration with notification system for billing emails
- ✅ Comprehensive audit trail for compliance
- ✅ Multi-tenant support with organization isolation

---

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────────┐
│                    BILLING SYSTEM ARCHITECTURE                   │
└─────────────────────────────────────────────────────────────────┘

┌──────────────────────┐      ┌──────────────────────┐
│   Usage Collector    │      │  Subscription Mgmt   │
│   • API metrics      │      │  • Plans & tiers     │
│   • Storage metrics  │      │  • Add-ons           │
│   • Transfer metrics │      │  • Seat management   │
│   • Query metrics    │      │  • Upgrades/downgrd  │
└──────────┬───────────┘      └──────────┬───────────┘
           │                             │
           ▼                             ▼
┌────────────────────────────────────────────────────┐
│          PRICING CALCULATION ENGINE                │
│  • Rate cards (per-resource pricing)              │
│  • Tiered pricing (volume discounts)              │
│  • Fractional calculations (sub-unit precision)   │
│  • Currency conversion                            │
│  • Tax calculation                                │
└────────────────────┬───────────────────────────────┘
                     │
                     ▼
┌────────────────────────────────────────────────────┐
│          INVOICE GENERATION                        │
│  • Line item aggregation                          │
│  • PDF generation                                 │
│  • Invoice numbering                              │
│  • Multi-currency support                         │
└────────────────────┬───────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│  Payment    │ │ Notification│ │   Event     │
│  Processing │ │   Service   │ │    Bus      │
│  (Stripe)   │ │  (Email)    │ │  (Kafka)    │
└─────────────┘ └─────────────┘ └─────────────┘
```

### Data Flow

1. **Usage Collection**: Services emit usage metrics → Billing collector aggregates
2. **Metering**: Hourly/daily aggregation of usage data → Storage in time-series format
3. **Billing Cycle**: At cycle end, calculate charges based on rate cards
4. **Invoice Generation**: Create invoice with line items → Generate PDF
5. **Payment Processing**: Charge payment method → Handle success/failure
6. **Notification**: Send invoice email → Payment confirmations → Overdue notices

---

## Database Schema

### Entity Relationship Diagram

```
┌──────────────────────┐
│  Organization        │
│  (Billing Account)   │
└──────────┬───────────┘
           │
           │ 1:N
           │
┌──────────▼───────────┐       ┌──────────────────────┐
│  Subscription        │◄──────┤  SubscriptionPlan    │
│  (Active billing)    │  N:1  │  (Product catalog)   │
└──────────┬───────────┘       └──────────────────────┘
           │
           │ 1:N
           │
┌──────────▼───────────┐       ┌──────────────────────┐
│  Invoice             │       │  InvoiceLineItem     │
│  (Billing period)    │◄──────┤  (Charge details)    │
└──────────┬───────────┘  1:N  └──────────────────────┘
           │
           │ 1:N
           │
┌──────────▼───────────┐
│  Payment             │
│  (Transaction)       │
└──────────────────────┘

┌──────────────────────┐       ┌──────────────────────┐
│  UsageMetric         │       │  PricingTier         │
│  (Time-series data)  │       │  (Volume discounts)  │
└──────────────────────┘       └──────────────────────┘
```

### Core Tables

#### 1. `dictamesh_billing_organizations`
```sql
CREATE TABLE dictamesh_billing_organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    billing_email VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    tax_id VARCHAR(100),

    -- Address
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2), -- ISO 3166-1 alpha-2

    -- Billing settings
    currency VARCHAR(3) DEFAULT 'USD', -- ISO 4217
    billing_cycle VARCHAR(20) DEFAULT 'monthly', -- monthly, annual
    billing_day_of_month INT DEFAULT 1,
    timezone VARCHAR(50) DEFAULT 'UTC',

    -- Payment
    default_payment_method_id VARCHAR(255), -- Stripe payment method ID
    auto_pay BOOLEAN DEFAULT false,

    -- Status
    status VARCHAR(20) DEFAULT 'active', -- active, suspended, deleted

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT chk_billing_cycle CHECK (billing_cycle IN ('monthly', 'annual')),
    CONSTRAINT chk_status CHECK (status IN ('active', 'suspended', 'deleted'))
);

CREATE INDEX idx_dictamesh_billing_org_status ON dictamesh_billing_organizations(status);
CREATE INDEX idx_dictamesh_billing_org_email ON dictamesh_billing_organizations(billing_email);

COMMENT ON TABLE dictamesh_billing_organizations IS 'DictaMesh: Billing accounts and organization details';
```

#### 2. `dictamesh_billing_subscription_plans`
```sql
CREATE TABLE dictamesh_billing_subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,

    -- Pricing
    base_price DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    billing_interval VARCHAR(20) NOT NULL, -- monthly, annual

    -- Features (JSONB for flexibility)
    features JSONB DEFAULT '{}',

    -- Limits
    included_api_calls INT DEFAULT 0, -- 0 = unlimited
    included_storage_gb INT DEFAULT 0,
    included_data_transfer_gb INT DEFAULT 0,
    included_seats INT DEFAULT 1,
    max_adapters INT DEFAULT 0, -- 0 = unlimited

    -- Add-on pricing (per unit above included amount)
    price_per_api_call DECIMAL(12,6) DEFAULT 0,
    price_per_gb_storage DECIMAL(12,4) DEFAULT 0,
    price_per_gb_transfer DECIMAL(12,4) DEFAULT 0,
    price_per_additional_seat DECIMAL(12,2) DEFAULT 0,

    -- Status
    is_public BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_billing_interval CHECK (billing_interval IN ('monthly', 'annual'))
);

CREATE INDEX idx_dictamesh_billing_plan_slug ON dictamesh_billing_subscription_plans(slug);
CREATE INDEX idx_dictamesh_billing_plan_active ON dictamesh_billing_subscription_plans(is_active);

COMMENT ON TABLE dictamesh_billing_subscription_plans IS 'DictaMesh: Subscription plan catalog with pricing tiers';
```

#### 3. `dictamesh_billing_subscriptions`
```sql
CREATE TABLE dictamesh_billing_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),
    plan_id UUID NOT NULL REFERENCES dictamesh_billing_subscription_plans(id),

    -- Subscription details
    status VARCHAR(20) DEFAULT 'active', -- active, canceled, past_due, trialing
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,

    -- Trial
    trial_start TIMESTAMP,
    trial_end TIMESTAMP,

    -- Cancellation
    cancel_at_period_end BOOLEAN DEFAULT false,
    canceled_at TIMESTAMP,
    cancellation_reason TEXT,

    -- Pricing overrides (for custom deals)
    custom_pricing JSONB, -- Override plan pricing

    -- Seats
    quantity INT DEFAULT 1, -- Number of seats/licenses

    -- Payment
    stripe_subscription_id VARCHAR(255), -- External provider ID

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_subscription_status CHECK (status IN ('active', 'canceled', 'past_due', 'trialing', 'incomplete'))
);

CREATE INDEX idx_dictamesh_billing_sub_org ON dictamesh_billing_subscriptions(organization_id);
CREATE INDEX idx_dictamesh_billing_sub_status ON dictamesh_billing_subscriptions(status);
CREATE INDEX idx_dictamesh_billing_sub_period_end ON dictamesh_billing_subscriptions(current_period_end);

COMMENT ON TABLE dictamesh_billing_subscriptions IS 'DictaMesh: Active subscriptions linking organizations to plans';
```

#### 4. `dictamesh_billing_usage_metrics`
```sql
CREATE TABLE dictamesh_billing_usage_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),
    subscription_id UUID REFERENCES dictamesh_billing_subscriptions(id),

    -- Metric details
    metric_type VARCHAR(50) NOT NULL, -- api_calls, storage_gb, transfer_gb, query_seconds
    metric_value DECIMAL(20,6) NOT NULL,
    metric_unit VARCHAR(20) NOT NULL, -- count, GB, seconds, MB

    -- Time dimension
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,

    -- Metadata
    resource_id VARCHAR(255), -- Specific adapter, service, etc.
    metadata JSONB, -- Additional context

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_metric_type CHECK (metric_type IN (
        'api_calls', 'storage_gb', 'transfer_gb_in', 'transfer_gb_out',
        'query_seconds', 'graphql_operations', 'kafka_events', 'adapters_active'
    ))
) PARTITION BY RANGE (recorded_at);

-- Create monthly partitions
CREATE TABLE dictamesh_billing_usage_metrics_2025_01 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE INDEX idx_dictamesh_billing_usage_org ON dictamesh_billing_usage_metrics(organization_id, recorded_at);
CREATE INDEX idx_dictamesh_billing_usage_type ON dictamesh_billing_usage_metrics(metric_type, recorded_at);

COMMENT ON TABLE dictamesh_billing_usage_metrics IS 'DictaMesh: Time-series usage metrics for billing calculation (partitioned by month)';
```

#### 5. `dictamesh_billing_invoices`
```sql
CREATE TABLE dictamesh_billing_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),
    subscription_id UUID REFERENCES dictamesh_billing_subscriptions(id),

    -- Invoice identification
    invoice_number VARCHAR(50) NOT NULL UNIQUE, -- e.g., INV-2025-001234

    -- Billing period
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,

    -- Amounts
    subtotal DECIMAL(12,2) NOT NULL,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL,
    amount_due DECIMAL(12,2) NOT NULL,
    amount_paid DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Status
    status VARCHAR(20) DEFAULT 'draft', -- draft, open, paid, void, uncollectible

    -- Dates
    invoice_date TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL,
    paid_at TIMESTAMP,

    -- Payment
    stripe_invoice_id VARCHAR(255),

    -- PDF
    pdf_url TEXT,
    pdf_generated_at TIMESTAMP,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_invoice_status CHECK (status IN ('draft', 'open', 'paid', 'void', 'uncollectible'))
);

CREATE INDEX idx_dictamesh_billing_invoice_org ON dictamesh_billing_invoices(organization_id);
CREATE INDEX idx_dictamesh_billing_invoice_number ON dictamesh_billing_invoices(invoice_number);
CREATE INDEX idx_dictamesh_billing_invoice_status ON dictamesh_billing_invoices(status);
CREATE INDEX idx_dictamesh_billing_invoice_due_date ON dictamesh_billing_invoices(due_date);

COMMENT ON TABLE dictamesh_billing_invoices IS 'DictaMesh: Generated invoices for billing periods';
```

#### 6. `dictamesh_billing_invoice_line_items`
```sql
CREATE TABLE dictamesh_billing_invoice_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES dictamesh_billing_invoices(id) ON DELETE CASCADE,

    -- Line item details
    description TEXT NOT NULL,
    quantity DECIMAL(20,6) NOT NULL,
    unit_price DECIMAL(12,6) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,

    -- Categorization
    item_type VARCHAR(50) NOT NULL, -- subscription, usage, addon, credit, tax
    metric_type VARCHAR(50), -- Links to usage metric type

    -- Period (for usage items)
    period_start TIMESTAMP,
    period_end TIMESTAMP,

    -- Metadata
    metadata JSONB,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_item_type CHECK (item_type IN (
        'subscription_base', 'usage_api_calls', 'usage_storage', 'usage_transfer',
        'addon_seats', 'addon_support', 'credit', 'tax', 'discount'
    ))
);

CREATE INDEX idx_dictamesh_billing_line_invoice ON dictamesh_billing_invoice_line_items(invoice_id);
CREATE INDEX idx_dictamesh_billing_line_type ON dictamesh_billing_invoice_line_items(item_type);

COMMENT ON TABLE dictamesh_billing_invoice_line_items IS 'DictaMesh: Individual line items on invoices with detailed charges';
```

#### 7. `dictamesh_billing_payments`
```sql
CREATE TABLE dictamesh_billing_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),
    invoice_id UUID REFERENCES dictamesh_billing_invoices(id),

    -- Payment details
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Status
    status VARCHAR(20) DEFAULT 'pending', -- pending, succeeded, failed, refunded

    -- Payment method
    payment_method VARCHAR(50), -- card, ach, wire, paypal
    payment_method_id VARCHAR(255), -- Stripe payment method ID

    -- Provider details
    provider VARCHAR(20) DEFAULT 'stripe', -- stripe, paypal, manual
    provider_payment_id VARCHAR(255), -- External transaction ID
    provider_customer_id VARCHAR(255),

    -- Timestamps
    attempted_at TIMESTAMP,
    succeeded_at TIMESTAMP,
    failed_at TIMESTAMP,
    refunded_at TIMESTAMP,

    -- Error handling
    failure_code VARCHAR(50),
    failure_message TEXT,

    -- Metadata
    metadata JSONB,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_payment_status CHECK (status IN ('pending', 'succeeded', 'failed', 'refunded', 'canceled'))
);

CREATE INDEX idx_dictamesh_billing_payment_org ON dictamesh_billing_payments(organization_id);
CREATE INDEX idx_dictamesh_billing_payment_invoice ON dictamesh_billing_payments(invoice_id);
CREATE INDEX idx_dictamesh_billing_payment_status ON dictamesh_billing_payments(status);
CREATE INDEX idx_dictamesh_billing_payment_provider ON dictamesh_billing_payments(provider, provider_payment_id);

COMMENT ON TABLE dictamesh_billing_payments IS 'DictaMesh: Payment transactions and processing records';
```

#### 8. `dictamesh_billing_pricing_tiers`
```sql
CREATE TABLE dictamesh_billing_pricing_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID REFERENCES dictamesh_billing_subscription_plans(id),

    -- Tier definition
    metric_type VARCHAR(50) NOT NULL,
    tier_start DECIMAL(20,2) NOT NULL, -- Inclusive lower bound
    tier_end DECIMAL(20,2), -- Exclusive upper bound (NULL = infinity)

    -- Pricing
    price_per_unit DECIMAL(12,6) NOT NULL,
    flat_fee DECIMAL(12,2) DEFAULT 0, -- Optional flat fee for entering tier

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dictamesh_billing_tier_plan ON dictamesh_billing_pricing_tiers(plan_id, metric_type);

COMMENT ON TABLE dictamesh_billing_pricing_tiers IS 'DictaMesh: Tiered pricing for volume-based discounts';
```

#### 9. `dictamesh_billing_credits`
```sql
CREATE TABLE dictamesh_billing_credits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),

    -- Credit details
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    remaining_amount DECIMAL(12,2) NOT NULL,

    -- Reason
    reason VARCHAR(100) NOT NULL, -- promotional, refund, compensation, discount
    description TEXT,

    -- Validity
    valid_from TIMESTAMP NOT NULL DEFAULT NOW(),
    valid_until TIMESTAMP,

    -- Status
    status VARCHAR(20) DEFAULT 'active', -- active, exhausted, expired, voided

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_credit_status CHECK (status IN ('active', 'exhausted', 'expired', 'voided'))
);

CREATE INDEX idx_dictamesh_billing_credit_org ON dictamesh_billing_credits(organization_id);
CREATE INDEX idx_dictamesh_billing_credit_status ON dictamesh_billing_credits(status, valid_until);

COMMENT ON TABLE dictamesh_billing_credits IS 'DictaMesh: Account credits for discounts and promotions';
```

#### 10. `dictamesh_billing_audit_log`
```sql
CREATE TABLE dictamesh_billing_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Entity tracking
    entity_type VARCHAR(50) NOT NULL, -- subscription, invoice, payment
    entity_id UUID NOT NULL,

    -- Event details
    event_type VARCHAR(50) NOT NULL, -- created, updated, deleted, status_changed
    event_data JSONB NOT NULL,

    -- Actor
    actor_id VARCHAR(255), -- User ID or system
    actor_type VARCHAR(20) DEFAULT 'system', -- user, system, webhook

    -- Context
    ip_address INET,
    user_agent TEXT,

    -- Timestamp
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dictamesh_billing_audit_entity ON dictamesh_billing_audit_log(entity_type, entity_id);
CREATE INDEX idx_dictamesh_billing_audit_occurred ON dictamesh_billing_audit_log(occurred_at);

COMMENT ON TABLE dictamesh_billing_audit_log IS 'DictaMesh: Comprehensive audit trail for billing operations';
```

---

## Pricing Model

### Subscription Tiers

| Tier | Base Price | Included | Overage Pricing |
|------|-----------|----------|-----------------|
| **Free** | $0/month | 10K API calls, 1GB storage, 1 adapter | $0.01/1K calls, $0.50/GB storage |
| **Starter** | $99/month | 1M API calls, 50GB storage, 5 adapters | $0.005/1K calls, $0.25/GB storage |
| **Professional** | $499/month | 10M API calls, 500GB storage, 25 adapters | $0.003/1K calls, $0.15/GB storage |
| **Enterprise** | $2,499/month | 100M API calls, 5TB storage, unlimited adapters | $0.001/1K calls, $0.10/GB storage |

### Usage-Based Metrics

#### API Calls
- **Granularity**: Per 1,000 calls
- **Calculation**: Round up to nearest 1,000
- **Example**: 1,500 calls = 2 billing units

#### Storage
- **Granularity**: Per GB-hour (fractional)
- **Calculation**: Average hourly storage × hours in period
- **Example**: 100GB for 15 days = (100 × 15 × 24) / (30 × 24) = 50 GB-month

#### Data Transfer
- **Granularity**: Per GB
- **Calculation**: Cumulative bytes / 1,073,741,824
- **Example**: 1.5GB in + 2.3GB out = 3.8 GB total

#### Query Processing
- **Granularity**: Per second of CPU time
- **Calculation**: Sum of query execution times
- **Example**: 1,000 queries averaging 50ms = 50 seconds

### Tiered Pricing Example

For API calls on **Professional** plan:
```
Tier 1: 0 - 10M calls     = Included (base price)
Tier 2: 10M - 50M calls   = $0.003 per 1K
Tier 3: 50M - 100M calls  = $0.002 per 1K
Tier 4: 100M+ calls       = $0.001 per 1K
```

**Calculation for 75M calls:**
```
Included: 10M calls      = $0
Tier 2:   40M calls      = 40,000 × $0.003 = $120
Tier 3:   25M calls      = 25,000 × $0.002 = $50
Total:                   = $170 overage charge
```

### Fractional Pricing Precision

All calculations use **6 decimal places** for unit pricing:
```go
type Money struct {
    Amount   decimal.Decimal // Use shopspring/decimal for precision
    Currency string
}

// Example: $0.000123 per API call
unitPrice := decimal.NewFromFloat(0.000123)
calls := decimal.NewFromInt(1_500_000)
charge := unitPrice.Mul(calls) // = $184.50
```

---

## Usage Metrics Collection

### Collection Architecture

```
┌─────────────────────────────────────────────────────┐
│              SERVICES (Emit Metrics)                │
│  • GraphQL Gateway      • Metadata Catalog          │
│  • Adapters             • Event Router              │
└───────────────────┬─────────────────────────────────┘
                    │
                    ▼
          ┌─────────────────────┐
          │   Metrics Collector │
          │   (Prometheus)      │
          └─────────┬───────────┘
                    │
                    ▼
          ┌─────────────────────┐
          │  Metrics Aggregator │
          │  (Hourly rollup)    │
          └─────────┬───────────┘
                    │
                    ▼
          ┌─────────────────────┐
          │  Billing Database   │
          │  (usage_metrics)    │
          └─────────────────────┘
```

### Metric Types

#### 1. API Calls
```go
// Counter metric
apiCallsTotal := promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "dictamesh_api_calls_total",
        Help: "Total API calls by organization and endpoint",
    },
    []string{"organization_id", "endpoint", "method"},
)

// Record call
apiCallsTotal.WithLabelValues(orgID, "/graphql", "POST").Inc()
```

#### 2. Storage
```go
// Gauge metric (current value)
storageBytes := promauto.NewGaugeVec(
    prometheus.GaugeOpts{
        Name: "dictamesh_storage_bytes",
        Help: "Current storage usage in bytes",
    },
    []string{"organization_id", "storage_type"},
)

// Update storage
storageBytes.WithLabelValues(orgID, "metadata").Set(12345678)
```

#### 3. Data Transfer
```go
// Counter metric
transferBytes := promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "dictamesh_transfer_bytes_total",
        Help: "Total data transfer in bytes",
    },
    []string{"organization_id", "direction"}, // in/out
)

// Record transfer
transferBytes.WithLabelValues(orgID, "out").Add(1024000)
```

#### 4. Query Processing
```go
// Histogram metric (for percentiles)
queryDuration := promauto.NewHistogramVec(
    prometheus.HistogramOpts{
        Name: "dictamesh_query_duration_seconds",
        Help: "Query processing duration",
        Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
    },
    []string{"organization_id", "query_type"},
)

// Record query
timer := prometheus.NewTimer(queryDuration.WithLabelValues(orgID, "graphql"))
defer timer.ObserveDuration()
```

### Aggregation Strategy

**Hourly Rollup Job**:
```sql
-- Aggregate API calls per hour
INSERT INTO dictamesh_billing_usage_metrics (
    organization_id,
    subscription_id,
    metric_type,
    metric_value,
    metric_unit,
    period_start,
    period_end,
    recorded_at
)
SELECT
    organization_id,
    subscription_id,
    'api_calls' as metric_type,
    COUNT(*) as metric_value,
    'count' as metric_unit,
    date_trunc('hour', timestamp) as period_start,
    date_trunc('hour', timestamp) + interval '1 hour' as period_end,
    NOW() as recorded_at
FROM raw_api_logs
WHERE timestamp >= NOW() - INTERVAL '1 hour'
  AND timestamp < NOW()
GROUP BY organization_id, subscription_id, date_trunc('hour', timestamp);
```

---

## Invoice Generation

### Generation Flow

```
1. Billing Period Ends
   ↓
2. Aggregate Usage Metrics
   ↓
3. Calculate Charges
   • Base subscription fee
   • Usage overage charges
   • Add-on charges
   • Apply credits
   • Calculate tax
   ↓
4. Create Invoice Record
   ↓
5. Generate Line Items
   ↓
6. Generate PDF
   ↓
7. Send Notification
   ↓
8. Process Payment (if auto-pay)
```

### Invoice Line Item Examples

**Example Invoice for Professional Plan:**

```
Invoice #INV-2025-001234
Organization: Acme Corp
Period: Jan 1, 2025 - Jan 31, 2025

LINE ITEMS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Description                      Qty        Rate      Amount
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Professional Plan (Jan 2025)     1      $499.00     $499.00

API Calls
  Included: 10M calls
  Usage:    15M calls
  Overage:  5M calls           5,000    $0.003       $15.00

Storage (avg 75GB)
  Included: 50GB
  Overage:  25GB                  25    $0.25         $6.25

Data Transfer
  Outbound: 120GB                120    $0.10        $12.00

Additional Seat                    2   $99.00       $198.00
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                                           Subtotal: $730.25
                                          Tax (10%):  $73.03
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                                        TOTAL DUE:   $803.28
```

### PDF Generation

Using a Go library like `go-pdf` or `wkhtmltopdf`:

```go
func GenerateInvoicePDF(invoice *Invoice) ([]byte, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()

    // Header
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, fmt.Sprintf("Invoice #%s", invoice.InvoiceNumber))

    // Company details
    pdf.SetFont("Arial", "", 10)
    pdf.Cell(40, 10, invoice.Organization.Name)

    // Line items table
    for _, item := range invoice.LineItems {
        pdf.Cell(100, 10, item.Description)
        pdf.Cell(30, 10, fmt.Sprintf("%.2f", item.Amount))
    }

    // Total
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(40, 10, fmt.Sprintf("Total: $%.2f", invoice.TotalAmount))

    return pdf.OutputFileAndClose("invoice.pdf")
}
```

---

## Payment Processing

### Stripe Integration

```go
package payment

import (
    "github.com/stripe/stripe-go/v75"
    "github.com/stripe/stripe-go/v75/customer"
    "github.com/stripe/stripe-go/v75/paymentintent"
)

type StripeProvider struct {
    apiKey string
}

func (p *StripeProvider) CreateCustomer(org *Organization) (string, error) {
    stripe.Key = p.apiKey

    params := &stripe.CustomerParams{
        Email: stripe.String(org.BillingEmail),
        Name:  stripe.String(org.Name),
        Metadata: map[string]string{
            "organization_id": org.ID,
        },
    }

    cust, err := customer.New(params)
    if err != nil {
        return "", err
    }

    return cust.ID, nil
}

func (p *StripeProvider) ChargeInvoice(invoice *Invoice) (*Payment, error) {
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(int64(invoice.TotalAmount * 100)), // Cents
        Currency: stripe.String(strings.ToLower(invoice.Currency)),
        Customer: stripe.String(invoice.Organization.StripeCustomerID),
        Metadata: map[string]string{
            "invoice_id": invoice.ID,
            "organization_id": invoice.OrganizationID,
        },
    }

    pi, err := paymentintent.New(params)
    if err != nil {
        return nil, err
    }

    payment := &Payment{
        InvoiceID:          invoice.ID,
        OrganizationID:     invoice.OrganizationID,
        Amount:             invoice.TotalAmount,
        Currency:           invoice.Currency,
        Provider:           "stripe",
        ProviderPaymentID:  pi.ID,
        Status:             string(pi.Status),
    }

    return payment, nil
}
```

### Webhook Handling

```go
func HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
    payload, _ := ioutil.ReadAll(r.Body)
    event := stripe.Event{}

    if err := json.Unmarshal(payload, &event); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    switch event.Type {
    case "payment_intent.succeeded":
        var pi stripe.PaymentIntent
        json.Unmarshal(event.Data.Raw, &pi)
        handlePaymentSuccess(&pi)

    case "payment_intent.payment_failed":
        var pi stripe.PaymentIntent
        json.Unmarshal(event.Data.Raw, &pi)
        handlePaymentFailure(&pi)

    case "invoice.payment_succeeded":
        // Subscription renewal succeeded
        handleSubscriptionRenewal(event)
    }

    w.WriteHeader(http.StatusOK)
}
```

---

## Notification Integration

### Email Templates

#### 1. Invoice Generated
```yaml
template_code: billing_invoice_generated
channels: [email]
subject: "Your DictaMesh Invoice #{{.InvoiceNumber}}"
body_html: |
  <h1>Invoice #{{.InvoiceNumber}}</h1>
  <p>Dear {{.OrganizationName}},</p>

  <p>Your invoice for the period {{.PeriodStart}} - {{.PeriodEnd}} is ready.</p>

  <table>
    <tr><td>Subtotal:</td><td>${{.Subtotal}}</td></tr>
    <tr><td>Tax:</td><td>${{.Tax}}</td></tr>
    <tr><th>Total:</th><th>${{.Total}}</th></tr>
  </table>

  <p><a href="{{.InvoiceURL}}">View Invoice</a></p>

  {{if .AutoPay}}
  <p>Your payment method will be charged automatically on {{.DueDate}}.</p>
  {{else}}
  <p>Please pay by {{.DueDate}} to avoid service interruption.</p>
  {{end}}
```

#### 2. Payment Succeeded
```yaml
template_code: billing_payment_succeeded
channels: [email]
subject: "Payment Received - Invoice #{{.InvoiceNumber}}"
body_html: |
  <h1>Payment Confirmed</h1>
  <p>Thank you! We've received your payment of ${{.Amount}}.</p>

  <p>Invoice: #{{.InvoiceNumber}}<br>
  Payment Method: {{.PaymentMethod}}<br>
  Transaction ID: {{.TransactionID}}</p>

  <p><a href="{{.ReceiptURL}}">View Receipt</a></p>
```

#### 3. Payment Failed
```yaml
template_code: billing_payment_failed
channels: [email]
subject: "Action Required: Payment Failed for Invoice #{{.InvoiceNumber}}"
body_html: |
  <h1>Payment Failed</h1>
  <p>We were unable to process your payment for invoice #{{.InvoiceNumber}}.</p>

  <p>Reason: {{.FailureReason}}</p>

  <p>Please update your payment method or pay manually to avoid service interruption.</p>

  <p><a href="{{.PaymentURL}}">Update Payment Method</a></p>
```

### Notification Rules

```go
// Create notification rule for invoice generation
rule := &NotificationRule{
    Name:        "Invoice Generated",
    EventType:   "billing.invoice.created",
    Condition:   "event.invoice.status == 'open'",
    TemplateID:  "billing_invoice_generated",
    Channels:    []string{"email"},
    Recipients: []string{"${event.invoice.organization.billing_email}"},
}
```

### Integration Code

```go
func SendInvoiceNotification(invoice *Invoice) error {
    notification := &Notification{
        RecipientID:   invoice.OrganizationID,
        RecipientType: "organization",
        TemplateCode:  "billing_invoice_generated",
        Channels:      []string{"email"},
        Priority:      "high",
        Data: map[string]interface{}{
            "InvoiceNumber":    invoice.InvoiceNumber,
            "OrganizationName": invoice.Organization.Name,
            "PeriodStart":      invoice.PeriodStart.Format("Jan 2, 2006"),
            "PeriodEnd":        invoice.PeriodEnd.Format("Jan 2, 2006"),
            "Subtotal":         invoice.Subtotal,
            "Tax":              invoice.TaxAmount,
            "Total":            invoice.TotalAmount,
            "DueDate":          invoice.DueDate.Format("Jan 2, 2006"),
            "InvoiceURL":       fmt.Sprintf("https://app.dictamesh.io/invoices/%s", invoice.ID),
            "AutoPay":          invoice.Organization.AutoPay,
        },
    }

    return notificationService.Send(notification)
}
```

---

## API Design

### REST API Endpoints

#### Organizations

```
POST   /api/v1/billing/organizations
GET    /api/v1/billing/organizations/:id
PUT    /api/v1/billing/organizations/:id
DELETE /api/v1/billing/organizations/:id
```

#### Subscriptions

```
POST   /api/v1/billing/subscriptions
GET    /api/v1/billing/subscriptions/:id
PUT    /api/v1/billing/subscriptions/:id
POST   /api/v1/billing/subscriptions/:id/cancel
POST   /api/v1/billing/subscriptions/:id/upgrade
POST   /api/v1/billing/subscriptions/:id/downgrade
GET    /api/v1/billing/subscriptions/:id/usage
```

#### Invoices

```
GET    /api/v1/billing/invoices
GET    /api/v1/billing/invoices/:id
GET    /api/v1/billing/invoices/:id/pdf
POST   /api/v1/billing/invoices/:id/pay
GET    /api/v1/billing/invoices/upcoming
```

#### Payments

```
GET    /api/v1/billing/payments
GET    /api/v1/billing/payments/:id
POST   /api/v1/billing/payments
POST   /api/v1/billing/payment-methods
GET    /api/v1/billing/payment-methods
DELETE /api/v1/billing/payment-methods/:id
```

#### Usage Metrics

```
GET    /api/v1/billing/usage/current
GET    /api/v1/billing/usage/history?start=2025-01-01&end=2025-01-31
GET    /api/v1/billing/usage/breakdown?metric=api_calls
```

### Example Request/Response

**Create Subscription:**
```http
POST /api/v1/billing/subscriptions
Content-Type: application/json

{
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "plan_slug": "professional",
  "quantity": 3,
  "billing_interval": "monthly",
  "payment_method_id": "pm_1234567890"
}
```

Response:
```json
{
  "id": "sub_abc123",
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "plan": {
    "id": "plan_xyz",
    "name": "Professional",
    "base_price": 499.00,
    "currency": "USD"
  },
  "status": "active",
  "current_period_start": "2025-01-01T00:00:00Z",
  "current_period_end": "2025-02-01T00:00:00Z",
  "quantity": 3,
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

## Event-Driven Architecture

### Kafka Topics

```
billing.subscription.created
billing.subscription.updated
billing.subscription.canceled
billing.invoice.created
billing.invoice.paid
billing.invoice.overdue
billing.payment.succeeded
billing.payment.failed
billing.usage.threshold_reached
billing.credit.applied
```

### Event Schema (Avro)

```json
{
  "type": "record",
  "name": "InvoiceCreated",
  "namespace": "io.dictamesh.billing.events",
  "fields": [
    {"name": "invoice_id", "type": "string"},
    {"name": "organization_id", "type": "string"},
    {"name": "invoice_number", "type": "string"},
    {"name": "total_amount", "type": "double"},
    {"name": "currency", "type": "string"},
    {"name": "status", "type": "string"},
    {"name": "due_date", "type": "long", "logicalType": "timestamp-millis"},
    {"name": "created_at", "type": "long", "logicalType": "timestamp-millis"}
  ]
}
```

### Event Producer

```go
func PublishInvoiceCreatedEvent(invoice *Invoice) error {
    event := map[string]interface{}{
        "invoice_id":      invoice.ID,
        "organization_id": invoice.OrganizationID,
        "invoice_number":  invoice.InvoiceNumber,
        "total_amount":    invoice.TotalAmount,
        "currency":        invoice.Currency,
        "status":          invoice.Status,
        "due_date":        invoice.DueDate.UnixMilli(),
        "created_at":      invoice.CreatedAt.UnixMilli(),
    }

    return kafkaProducer.Publish("billing.invoice.created", event)
}
```

---

## Observability

### Prometheus Metrics

```go
var (
    // Subscriptions
    activeSubscriptions = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "dictamesh_billing_active_subscriptions",
            Help: "Number of active subscriptions by plan",
        },
        []string{"plan"},
    )

    // Revenue
    monthlyRecurringRevenue = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "dictamesh_billing_mrr",
            Help: "Monthly recurring revenue in USD",
        },
    )

    // Invoices
    invoicesGenerated = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "dictamesh_billing_invoices_generated_total",
            Help: "Total invoices generated",
        },
        []string{"status"},
    )

    // Payments
    paymentAmount = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "dictamesh_billing_payment_amount",
            Help:    "Payment amounts",
            Buckets: prometheus.ExponentialBuckets(1, 2, 15),
        },
        []string{"currency", "status"},
    )

    // Payment failures
    paymentFailures = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "dictamesh_billing_payment_failures_total",
            Help: "Total payment failures",
        },
        []string{"failure_code"},
    )
)
```

### OpenTelemetry Traces

```go
func ProcessInvoice(ctx context.Context, invoice *Invoice) error {
    ctx, span := tracer.Start(ctx, "billing.process_invoice")
    defer span.End()

    span.SetAttributes(
        attribute.String("invoice.id", invoice.ID),
        attribute.String("organization.id", invoice.OrganizationID),
        attribute.Float64("invoice.amount", invoice.TotalAmount),
    )

    // Calculate charges
    _, span2 := tracer.Start(ctx, "billing.calculate_charges")
    charges, err := calculateCharges(ctx, invoice)
    span2.End()
    if err != nil {
        span.RecordError(err)
        return err
    }

    // Generate PDF
    _, span3 := tracer.Start(ctx, "billing.generate_pdf")
    pdf, err := generatePDF(ctx, invoice)
    span3.End()

    // Send notification
    _, span4 := tracer.Start(ctx, "billing.send_notification")
    err = sendNotification(ctx, invoice)
    span4.End()

    return nil
}
```

---

## Security Considerations

### 1. Payment Data Security

- **PCI DSS Compliance**: Never store card numbers; use Stripe tokens
- **Encryption**: All payment data encrypted at rest and in transit
- **Access Control**: Role-based access to billing data
- **Audit Logging**: Complete audit trail for all financial transactions

### 2. API Security

```go
// Middleware for billing API authentication
func BillingAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify JWT token
        token := extractToken(r)
        claims, err := validateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Check organization access
        orgID := r.URL.Query().Get("organization_id")
        if !hasOrganizationAccess(claims.UserID, orgID) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### 3. Webhook Signature Verification

```go
func VerifyStripeWebhook(r *http.Request) error {
    payload, _ := ioutil.ReadAll(r.Body)
    signature := r.Header.Get("Stripe-Signature")

    _, err := webhook.ConstructEvent(
        payload,
        signature,
        webhookSecret,
    )

    return err
}
```

### 4. Rate Limiting

```go
// Rate limit billing API to prevent abuse
var limiter = rate.NewLimiter(rate.Limit(100), 200) // 100 req/s, burst 200

func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1-2)
- [ ] Database schema implementation
- [ ] Core models (Organization, Subscription, Plan)
- [ ] Database migrations
- [ ] Basic CRUD operations

### Phase 2: Usage Metrics (Week 3)
- [ ] Prometheus metrics collection
- [ ] Metrics aggregation service
- [ ] Hourly rollup jobs
- [ ] Usage API endpoints

### Phase 3: Pricing Engine (Week 4)
- [ ] Pricing calculation logic
- [ ] Tiered pricing implementation
- [ ] Fractional pricing support
- [ ] Credits system

### Phase 4: Invoicing (Week 5-6)
- [ ] Invoice generation service
- [ ] Line item calculation
- [ ] PDF generation
- [ ] Invoice API endpoints

### Phase 5: Payment Processing (Week 7)
- [ ] Stripe integration
- [ ] Payment method management
- [ ] Webhook handlers
- [ ] Payment retry logic

### Phase 6: Notifications (Week 8)
- [ ] Email template creation
- [ ] Notification rules setup
- [ ] Integration with notification service
- [ ] Testing email delivery

### Phase 7: API & Service (Week 9-10)
- [ ] REST API implementation
- [ ] GraphQL schema (optional)
- [ ] Authentication & authorization
- [ ] Rate limiting

### Phase 8: Observability (Week 11)
- [ ] Prometheus metrics
- [ ] OpenTelemetry traces
- [ ] Grafana dashboards
- [ ] Alerting rules

### Phase 9: Testing & Documentation (Week 12)
- [ ] Unit tests
- [ ] Integration tests
- [ ] API documentation
- [ ] User guides

### Phase 10: Deployment (Week 13)
- [ ] Docker containerization
- [ ] Kubernetes manifests
- [ ] CI/CD pipeline
- [ ] Production deployment

---

## Conclusion

This design provides a comprehensive, production-ready billing system with:
- ✅ Flexible pricing models
- ✅ Accurate usage tracking
- ✅ Automated invoicing
- ✅ Multiple payment providers
- ✅ Comprehensive notifications
- ✅ Full observability
- ✅ Enterprise-grade security

The system is designed to scale with DictaMesh's growth while maintaining precision in billing calculations and providing excellent customer experience through clear invoicing and timely notifications.
