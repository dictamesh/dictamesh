-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- Migration: Add billing system tables
-- IMPORTANT: All billing tables use dictamesh_billing_ prefix for namespace isolation

-- =====================================================
-- 1. Organizations (Billing Accounts)
-- =====================================================

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
    billing_cycle VARCHAR(20) DEFAULT 'monthly',
    billing_day_of_month INT DEFAULT 1,
    timezone VARCHAR(50) DEFAULT 'UTC',

    -- Payment
    default_payment_method_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    auto_pay BOOLEAN DEFAULT false,

    -- Status
    status VARCHAR(20) DEFAULT 'active',

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT chk_billing_cycle CHECK (billing_cycle IN ('monthly', 'annual')),
    CONSTRAINT chk_org_status CHECK (status IN ('active', 'suspended', 'deleted'))
);

CREATE INDEX idx_dictamesh_billing_org_status ON dictamesh_billing_organizations(status);
CREATE INDEX idx_dictamesh_billing_org_email ON dictamesh_billing_organizations(billing_email);
CREATE INDEX idx_dictamesh_billing_org_stripe ON dictamesh_billing_organizations(stripe_customer_id);

COMMENT ON TABLE dictamesh_billing_organizations IS 'DictaMesh: Billing accounts and organization details';

-- =====================================================
-- 2. Subscription Plans (Product Catalog)
-- =====================================================

CREATE TABLE dictamesh_billing_subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,

    -- Pricing
    base_price DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    billing_interval VARCHAR(20) NOT NULL,

    -- Features (JSONB for flexibility)
    features JSONB DEFAULT '{}',

    -- Limits
    included_api_calls INT DEFAULT 0,
    included_storage_gb INT DEFAULT 0,
    included_data_transfer_gb INT DEFAULT 0,
    included_seats INT DEFAULT 1,
    max_adapters INT DEFAULT 0,

    -- Add-on pricing
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

    CONSTRAINT chk_plan_billing_interval CHECK (billing_interval IN ('monthly', 'annual'))
);

CREATE INDEX idx_dictamesh_billing_plan_slug ON dictamesh_billing_subscription_plans(slug);
CREATE INDEX idx_dictamesh_billing_plan_active ON dictamesh_billing_subscription_plans(is_active);

COMMENT ON TABLE dictamesh_billing_subscription_plans IS 'DictaMesh: Subscription plan catalog with pricing tiers';

-- =====================================================
-- 3. Subscriptions
-- =====================================================

CREATE TABLE dictamesh_billing_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),
    plan_id UUID NOT NULL REFERENCES dictamesh_billing_subscription_plans(id),

    -- Subscription details
    status VARCHAR(20) DEFAULT 'active',
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,

    -- Trial
    trial_start TIMESTAMP,
    trial_end TIMESTAMP,

    -- Cancellation
    cancel_at_period_end BOOLEAN DEFAULT false,
    canceled_at TIMESTAMP,
    cancellation_reason TEXT,

    -- Pricing overrides
    custom_pricing JSONB,

    -- Seats
    quantity INT DEFAULT 1,

    -- Payment provider
    stripe_subscription_id VARCHAR(255),

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_subscription_status CHECK (status IN ('active', 'canceled', 'past_due', 'trialing', 'incomplete'))
);

CREATE INDEX idx_dictamesh_billing_sub_org ON dictamesh_billing_subscriptions(organization_id);
CREATE INDEX idx_dictamesh_billing_sub_status ON dictamesh_billing_subscriptions(status);
CREATE INDEX idx_dictamesh_billing_sub_period_end ON dictamesh_billing_subscriptions(current_period_end);
CREATE INDEX idx_dictamesh_billing_sub_stripe ON dictamesh_billing_subscriptions(stripe_subscription_id);

COMMENT ON TABLE dictamesh_billing_subscriptions IS 'DictaMesh: Active subscriptions linking organizations to plans';

-- =====================================================
-- 4. Usage Metrics (Partitioned by Month)
-- =====================================================

CREATE TABLE dictamesh_billing_usage_metrics (
    id UUID DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    subscription_id UUID,

    -- Metric details
    metric_type VARCHAR(50) NOT NULL,
    metric_value DECIMAL(20,6) NOT NULL,
    metric_unit VARCHAR(20) NOT NULL,

    -- Time dimension
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,

    -- Metadata
    resource_id VARCHAR(255),
    metadata JSONB,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_metric_type CHECK (metric_type IN (
        'api_calls', 'storage_gb', 'transfer_gb_in', 'transfer_gb_out',
        'query_seconds', 'graphql_operations', 'kafka_events', 'adapters_active'
    )),
    PRIMARY KEY (id, recorded_at)
) PARTITION BY RANGE (recorded_at);

-- Create partitions for 2025
CREATE TABLE dictamesh_billing_usage_metrics_2025_01 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_02 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_03 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_04 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_05 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_06 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_07 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_08 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_09 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_10 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_11 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');

CREATE TABLE dictamesh_billing_usage_metrics_2025_12 PARTITION OF dictamesh_billing_usage_metrics
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE INDEX idx_dictamesh_billing_usage_org ON dictamesh_billing_usage_metrics(organization_id, recorded_at);
CREATE INDEX idx_dictamesh_billing_usage_type ON dictamesh_billing_usage_metrics(metric_type, recorded_at);
CREATE INDEX idx_dictamesh_billing_usage_period ON dictamesh_billing_usage_metrics(period_start, period_end);

COMMENT ON TABLE dictamesh_billing_usage_metrics IS 'DictaMesh: Time-series usage metrics for billing calculation (partitioned by month)';

-- =====================================================
-- 5. Invoices
-- =====================================================

CREATE TABLE dictamesh_billing_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),
    subscription_id UUID REFERENCES dictamesh_billing_subscriptions(id),

    -- Invoice identification
    invoice_number VARCHAR(50) NOT NULL UNIQUE,

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
    status VARCHAR(20) DEFAULT 'draft',

    -- Dates
    invoice_date TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL,
    paid_at TIMESTAMP,

    -- Payment provider
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
CREATE INDEX idx_dictamesh_billing_invoice_stripe ON dictamesh_billing_invoices(stripe_invoice_id);

COMMENT ON TABLE dictamesh_billing_invoices IS 'DictaMesh: Generated invoices for billing periods';

-- =====================================================
-- 6. Invoice Line Items
-- =====================================================

CREATE TABLE dictamesh_billing_invoice_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES dictamesh_billing_invoices(id) ON DELETE CASCADE,

    -- Line item details
    description TEXT NOT NULL,
    quantity DECIMAL(20,6) NOT NULL,
    unit_price DECIMAL(12,6) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,

    -- Categorization
    item_type VARCHAR(50) NOT NULL,
    metric_type VARCHAR(50),

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

-- =====================================================
-- 7. Payments
-- =====================================================

CREATE TABLE dictamesh_billing_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),
    invoice_id UUID REFERENCES dictamesh_billing_invoices(id),

    -- Payment details
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Status
    status VARCHAR(20) DEFAULT 'pending',

    -- Payment method
    payment_method VARCHAR(50),
    payment_method_id VARCHAR(255),

    -- Provider details
    provider VARCHAR(20) DEFAULT 'stripe',
    provider_payment_id VARCHAR(255),
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

-- =====================================================
-- 8. Pricing Tiers (Volume Discounts)
-- =====================================================

CREATE TABLE dictamesh_billing_pricing_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID REFERENCES dictamesh_billing_subscription_plans(id),

    -- Tier definition
    metric_type VARCHAR(50) NOT NULL,
    tier_start DECIMAL(20,2) NOT NULL,
    tier_end DECIMAL(20,2),

    -- Pricing
    price_per_unit DECIMAL(12,6) NOT NULL,
    flat_fee DECIMAL(12,2) DEFAULT 0,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dictamesh_billing_tier_plan ON dictamesh_billing_pricing_tiers(plan_id, metric_type);

COMMENT ON TABLE dictamesh_billing_pricing_tiers IS 'DictaMesh: Tiered pricing for volume-based discounts';

-- =====================================================
-- 9. Credits
-- =====================================================

CREATE TABLE dictamesh_billing_credits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES dictamesh_billing_organizations(id),

    -- Credit details
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    remaining_amount DECIMAL(12,2) NOT NULL,

    -- Reason
    reason VARCHAR(100) NOT NULL,
    description TEXT,

    -- Validity
    valid_from TIMESTAMP NOT NULL DEFAULT NOW(),
    valid_until TIMESTAMP,

    -- Status
    status VARCHAR(20) DEFAULT 'active',

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_credit_status CHECK (status IN ('active', 'exhausted', 'expired', 'voided'))
);

CREATE INDEX idx_dictamesh_billing_credit_org ON dictamesh_billing_credits(organization_id);
CREATE INDEX idx_dictamesh_billing_credit_status ON dictamesh_billing_credits(status, valid_until);

COMMENT ON TABLE dictamesh_billing_credits IS 'DictaMesh: Account credits for discounts and promotions';

-- =====================================================
-- 10. Audit Log
-- =====================================================

CREATE TABLE dictamesh_billing_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Entity tracking
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,

    -- Event details
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB NOT NULL,

    -- Actor
    actor_id VARCHAR(255),
    actor_type VARCHAR(20) DEFAULT 'system',

    -- Context
    ip_address INET,
    user_agent TEXT,

    -- Timestamp
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dictamesh_billing_audit_entity ON dictamesh_billing_audit_log(entity_type, entity_id);
CREATE INDEX idx_dictamesh_billing_audit_occurred ON dictamesh_billing_audit_log(occurred_at);

COMMENT ON TABLE dictamesh_billing_audit_log IS 'DictaMesh: Comprehensive audit trail for billing operations';

-- =====================================================
-- Seed Data: Default Subscription Plans
-- =====================================================

INSERT INTO dictamesh_billing_subscription_plans (name, slug, description, base_price, currency, billing_interval,
    included_api_calls, included_storage_gb, included_data_transfer_gb, included_seats, max_adapters,
    price_per_api_call, price_per_gb_storage, price_per_gb_transfer, price_per_additional_seat, is_public, is_active)
VALUES
    ('Free', 'free', 'Perfect for getting started', 0.00, 'USD', 'monthly',
     10000, 1, 1, 1, 1,
     0.00001, 0.50, 0.50, 0, true, true),

    ('Starter', 'starter', 'For small teams and projects', 99.00, 'USD', 'monthly',
     1000000, 50, 100, 5, 5,
     0.000005, 0.25, 0.25, 49.00, true, true),

    ('Professional', 'professional', 'For growing businesses', 499.00, 'USD', 'monthly',
     10000000, 500, 1000, 25, 25,
     0.000003, 0.15, 0.15, 99.00, true, true),

    ('Enterprise', 'enterprise', 'For large-scale operations', 2499.00, 'USD', 'monthly',
     100000000, 5000, 10000, 100, 0,
     0.000001, 0.10, 0.10, 199.00, true, true);
