# DictaMesh Billing System

## Overview

The DictaMesh Billing System is a comprehensive, enterprise-grade billing solution designed to handle subscription management, usage-based billing, invoicing, and payment processing for the DictaMesh platform.

## Features

✅ **Flexible Subscription Management**
- Multiple pricing tiers (Free, Starter, Professional, Enterprise)
- Monthly and annual billing cycles
- Seat-based pricing
- Custom pricing overrides for enterprise deals

✅ **Usage-Based Billing**
- Real-time metrics collection via Prometheus
- Multiple metric types (API calls, storage, data transfer, query processing)
- Fractional pricing with 6 decimal precision
- Hourly usage aggregation

✅ **Advanced Pricing**
- Tiered pricing for volume discounts
- Prorated billing for mid-cycle changes
- Account credits and promotions
- Tax calculation support

✅ **Automated Invoicing**
- Automatic invoice generation
- Detailed line items
- PDF generation
- Multi-currency support

✅ **Payment Processing**
- Stripe integration
- Multiple payment methods
- Automatic payment retry
- Webhook handling

✅ **Notifications**
- Integration with DictaMesh notification system
- Email notifications for all billing events
- Customizable templates
- Multi-channel delivery

✅ **Event-Driven Architecture**
- Kafka event publishing
- Subscription lifecycle events
- Payment events
- Usage threshold alerts

✅ **Observability**
- Prometheus metrics
- OpenTelemetry distributed tracing
- Comprehensive audit logging
- Real-time monitoring dashboards

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    BILLING SYSTEM ARCHITECTURE                   │
└─────────────────────────────────────────────────────────────────┘

┌──────────────────────┐      ┌──────────────────────┐
│   Usage Collector    │      │  Subscription Mgmt   │
│   • Prometheus       │      │  • Plans & tiers     │
│   • Metrics agg      │      │  • Lifecycle mgmt    │
└──────────┬───────────┘      └──────────┬───────────┘
           │                             │
           ▼                             ▼
┌────────────────────────────────────────────────────┐
│          PRICING CALCULATION ENGINE                │
│  • Rate cards         • Tiered pricing             │
│  • Fractional calc    • Proration                  │
└────────────────────┬───────────────────────────────┘
                     │
                     ▼
┌────────────────────────────────────────────────────┐
│          INVOICE GENERATION                        │
│  • Line items         • PDF generation             │
└────────────────────┬───────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│  Payment    │ │ Notification│ │   Event     │
│  (Stripe)   │ │   Service   │ │    Bus      │
└─────────────┘ └─────────────┘ └─────────────┘
```

## Package Structure

```
pkg/billing/
├── types.go              # Core type definitions
├── config.go             # Configuration management
├── models/
│   └── models.go         # GORM database models
├── pricing.go            # Pricing calculation engine
├── metrics.go            # Usage metrics collection
├── invoice.go            # Invoice generation
├── payment.go            # Payment processing (Stripe)
├── notifications.go      # Notification integration
├── events.go             # Kafka event publishing
├── observability.go      # Prometheus & OpenTelemetry
└── README.md            # This file
```

## Database Schema

### Core Tables

- `dictamesh_billing_organizations` - Billing accounts
- `dictamesh_billing_subscription_plans` - Product catalog
- `dictamesh_billing_subscriptions` - Active subscriptions
- `dictamesh_billing_usage_metrics` - Time-series usage data (partitioned)
- `dictamesh_billing_invoices` - Generated invoices
- `dictamesh_billing_invoice_line_items` - Invoice charges
- `dictamesh_billing_payments` - Payment transactions
- `dictamesh_billing_pricing_tiers` - Volume-based pricing
- `dictamesh_billing_credits` - Account credits
- `dictamesh_billing_audit_log` - Comprehensive audit trail

## Usage Examples

### Initialize the Billing System

```go
import (
    "github.com/dictamesh/dictamesh/pkg/billing"
)

// Load configuration
config, err := billing.LoadFromEnv()
if err != nil {
    log.Fatal(err)
}

// Create services
pricingEngine := billing.NewPricingEngine(config)
metricsCollector := billing.NewMetricsCollector(db, config)
invoiceService := billing.NewInvoiceService(db, config, pricingEngine, metricsCollector)
paymentService := billing.NewPaymentService(db, config, invoiceService)
notificationService := billing.NewNotificationService(config)
```

### Create a Subscription

```go
subscription := &models.Subscription{
    OrganizationID:     orgID,
    PlanID:             planID,
    Status:             "active",
    CurrentPeriodStart: time.Now(),
    CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
    Quantity:           5, // 5 seats
}

if err := db.Create(subscription).Error; err != nil {
    return err
}

// Publish event
eventPublisher.PublishSubscriptionCreated(ctx, subscription)
```

### Record Usage Metrics

```go
// Record API call
metricsCollector.RecordAPICall(organizationID, "/graphql", "POST")

// Record storage
metricsCollector.RecordStorage(organizationID, "metadata", 50*1024*1024*1024) // 50GB

// Record data transfer
metricsCollector.RecordTransfer(organizationID, "out", 1024*1024*1024) // 1GB
```

### Generate an Invoice

```go
invoice, err := invoiceService.GenerateInvoice(ctx, subscriptionID)
if err != nil {
    return err
}

// Send notification
notificationService.SendInvoiceCreatedNotification(ctx, invoice)

// Publish event
eventPublisher.PublishInvoiceCreated(ctx, invoice)
```

### Process a Payment

```go
payment, err := paymentService.ChargeInvoice(ctx, invoiceID)
if err != nil {
    // Handle payment failure
    notificationService.SendPaymentFailedNotification(ctx, payment, invoice)
    eventPublisher.PublishPaymentFailed(ctx, payment)
    return err
}

// Payment succeeded
notificationService.SendPaymentSucceededNotification(ctx, payment, invoice)
eventPublisher.PublishPaymentSucceeded(ctx, payment)
```

### Calculate Pricing

```go
// Fetch usage for billing period
usage, err := metricsCollector.GetUsageForPeriod(ctx, orgID, periodStart, periodEnd)

// Fetch credits
var credits []models.Credit
db.Where("organization_id = ? AND status = ?", orgID, "active").Find(&credits)

// Calculate charges
calc, err := pricingEngine.CalculateSubscriptionCharge(
    subscription,
    plan,
    usage,
    credits,
)

// calc.Total contains the final amount
// calc.LineItems contains itemized charges
```

## Configuration

Configure the billing system via environment variables:

```bash
# Database
BILLING_DATABASE_DSN=postgres://user:pass@localhost/dictamesh

# Stripe
STRIPE_API_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_ENABLED=true

# Invoice Settings
INVOICE_DUE_DAYS=30
INVOICE_NUMBER_PREFIX=INV-
INVOICE_TAX_RATE=0.10
INVOICE_DEFAULT_CURRENCY=USD

# Usage Metrics
USAGE_AGGREGATION_INTERVAL=1h
USAGE_RETENTION_DAYS=90
USAGE_ENABLE_REALTIME=true

# Notifications
NOTIFICATION_SERVICE_URL=http://localhost:8080
NOTIFICATION_RETRY_ATTEMPTS=3

# Feature Flags
FEATURE_AUTO_PAYMENT=true
FEATURE_USAGE_METRICS=true
FEATURE_TIERED_PRICING=true
FEATURE_CREDITS=true
```

## Subscription Plans

### Default Plans

| Plan | Price | API Calls | Storage | Transfer | Adapters |
|------|-------|-----------|---------|----------|----------|
| **Free** | $0/mo | 10K | 1GB | 1GB | 1 |
| **Starter** | $99/mo | 1M | 50GB | 100GB | 5 |
| **Professional** | $499/mo | 10M | 500GB | 1TB | 25 |
| **Enterprise** | $2,499/mo | 100M | 5TB | 10TB | Unlimited |

### Overage Pricing

- **API Calls**: $0.01 - $0.000001 per 1K calls (volume discounts)
- **Storage**: $0.50 - $0.10 per GB (volume discounts)
- **Transfer**: $0.50 - $0.10 per GB (volume discounts)
- **Additional Seats**: $49 - $199 per seat (plan-dependent)

## Metrics

### Prometheus Metrics

```
# Subscriptions
dictamesh_billing_active_subscriptions{plan="professional"} 150

# Revenue
dictamesh_billing_mrr 74850.00
dictamesh_billing_arr 898200.00

# Invoices
dictamesh_billing_invoices_generated_total{status="paid"} 1234

# Payments
dictamesh_billing_payments_processed_total{status="succeeded",provider="stripe"} 987
dictamesh_billing_payment_failures_total{failure_code="card_declined"} 23
```

### Events

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

## Notification Templates

### Available Templates

1. **billing_invoice_generated** - New invoice created
2. **billing_payment_succeeded** - Payment confirmation
3. **billing_payment_failed** - Payment failure alert
4. **billing_invoice_overdue** - Overdue invoice notice
5. **billing_subscription_created** - Welcome email
6. **billing_subscription_canceled** - Cancellation confirmation
7. **billing_usage_threshold_reached** - Usage alert
8. **billing_upcoming_renewal** - Renewal reminder

## API Integration

### REST API Endpoints

```
# Organizations
POST   /api/v1/billing/organizations
GET    /api/v1/billing/organizations/:id
PUT    /api/v1/billing/organizations/:id

# Subscriptions
POST   /api/v1/billing/subscriptions
GET    /api/v1/billing/subscriptions/:id
POST   /api/v1/billing/subscriptions/:id/cancel

# Invoices
GET    /api/v1/billing/invoices
GET    /api/v1/billing/invoices/:id
GET    /api/v1/billing/invoices/:id/pdf

# Payments
POST   /api/v1/billing/payments
GET    /api/v1/billing/payments/:id
POST   /api/v1/billing/payment-methods

# Usage
GET    /api/v1/billing/usage/current
GET    /api/v1/billing/usage/history
```

## Security

### Best Practices

- ✅ Never store raw card numbers (use Stripe tokens)
- ✅ All payment data encrypted at rest and in transit
- ✅ PCI DSS compliant payment processing
- ✅ Webhook signature verification
- ✅ Role-based access control
- ✅ Comprehensive audit logging
- ✅ Rate limiting on billing APIs

### Webhook Verification

```go
// Verify Stripe webhook signature
stripe.Key = config.Stripe.APIKey

event, err := webhook.ConstructEvent(
    payload,
    signature,
    config.Stripe.WebhookSecret,
)
```

## Testing

### Unit Tests

```bash
go test ./pkg/billing/...
```

### Integration Tests

```bash
go test -tags=integration ./pkg/billing/...
```

### Test Coverage

```bash
go test -cover ./pkg/billing/...
```

## Deployment

### Database Migration

```bash
# Run migrations
migrate -path pkg/database/migrations/sql \
        -database postgres://localhost/dictamesh up

# Rollback
migrate -path pkg/database/migrations/sql \
        -database postgres://localhost/dictamesh down 1
```

### Docker Deployment

```bash
# Build
docker build -t dictamesh-billing .

# Run
docker run -p 8080:8080 \
  -e BILLING_DATABASE_DSN=... \
  -e STRIPE_API_KEY=... \
  dictamesh-billing
```

## Monitoring

### Grafana Dashboards

Import the provided Grafana dashboard to monitor:
- Monthly Recurring Revenue (MRR)
- Active Subscriptions
- Payment Success Rate
- Invoice Aging
- Usage Trends

### Alerts

Configure alerts for:
- Payment failure rate > 5%
- Invoice collection rate < 95%
- Unusual usage spikes
- Credit balance alerts

## Troubleshooting

### Common Issues

**Q: Invoices not generating**
- Check subscription period_end dates
- Verify usage metrics are being collected
- Check invoice generation cron job

**Q: Payments failing**
- Verify Stripe API keys
- Check webhook endpoint configuration
- Review payment method status

**Q: Usage metrics missing**
- Verify Prometheus scraping configuration
- Check metrics aggregation worker
- Review database partitions

## Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for guidelines.

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda

## Support

For issues and questions:
- GitHub Issues: https://github.com/dictamesh/dictamesh/issues
- Documentation: https://docs.dictamesh.com/billing
