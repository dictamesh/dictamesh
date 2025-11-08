// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/Click2-Run/dictamesh/pkg/billing/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// InvoiceService handles invoice generation and management
type InvoiceService struct {
	db             *gorm.DB
	config         *Config
	pricingEngine  *PricingEngine
	metricsCollector *MetricsCollector
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(
	db *gorm.DB,
	config *Config,
	pricingEngine *PricingEngine,
	metricsCollector *MetricsCollector,
) *InvoiceService {
	return &InvoiceService{
		db:             db,
		config:         config,
		pricingEngine:  pricingEngine,
		metricsCollector: metricsCollector,
	}
}

// GenerateInvoice generates an invoice for a subscription billing period
func (is *InvoiceService) GenerateInvoice(
	ctx context.Context,
	subscriptionID string,
) (*models.Invoice, error) {
	// 1. Fetch subscription with plan and organization
	var subscription models.Subscription
	if err := is.db.WithContext(ctx).
		Preload("Plan").
		Preload("Organization").
		First(&subscription, "id = ?", subscriptionID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	// 2. Fetch usage metrics for the billing period
	usage, err := is.metricsCollector.GetUsageForPeriod(
		ctx,
		subscription.OrganizationID.String(),
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage metrics: %w", err)
	}

	// 3. Fetch available credits
	var credits []models.Credit
	if err := is.db.WithContext(ctx).
		Where("organization_id = ?", subscription.OrganizationID).
		Where("status = ?", CreditStatusActive).
		Where("valid_from <= ?", time.Now()).
		Where("valid_until IS NULL OR valid_until >= ?", time.Now()).
		Where("remaining_amount > 0").
		Order("valid_from ASC").
		Find(&credits).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch credits: %w", err)
	}

	// 4. Calculate charges
	calc, err := is.pricingEngine.CalculateSubscriptionCharge(
		&subscription,
		&subscription.Plan,
		usage,
		credits,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate charges: %w", err)
	}

	// 5. Generate invoice number
	invoiceNumber, err := is.generateInvoiceNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	// 6. Create invoice record
	invoice := &models.Invoice{
		ID:             uuid.New(),
		OrganizationID: subscription.OrganizationID,
		SubscriptionID: subscription.ID,
		InvoiceNumber:  invoiceNumber,
		PeriodStart:    subscription.CurrentPeriodStart,
		PeriodEnd:      subscription.CurrentPeriodEnd,
		Subtotal:       calc.Subtotal,
		TaxAmount:      calc.TaxAmount,
		TotalAmount:    calc.Total,
		AmountDue:      calc.Total,
		AmountPaid:     decimal.Zero,
		Currency:       subscription.Plan.Currency,
		Status:         string(InvoiceStatusOpen),
		InvoiceDate:    time.Now(),
		DueDate:        time.Now().AddDate(0, 0, is.config.Invoice.DueDays),
	}

	// 7. Begin transaction
	tx := is.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 8. Save invoice
	if err := tx.Create(invoice).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// 9. Save line items
	for _, lineItem := range calc.LineItems {
		dbLineItem := &models.InvoiceLineItem{
			ID:          uuid.New(),
			InvoiceID:   invoice.ID,
			Description: lineItem.Description,
			Quantity:    lineItem.Quantity,
			UnitPrice:   lineItem.UnitPrice,
			Amount:      lineItem.Amount,
			ItemType:    string(lineItem.ItemType),
			MetricType:  string(lineItem.MetricType),
			PeriodStart: lineItem.PeriodStart,
			PeriodEnd:   lineItem.PeriodEnd,
			Metadata:    models.JSONB(lineItem.Metadata),
		}

		if err := tx.Create(dbLineItem).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create line item: %w", err)
		}
	}

	// 10. Update credits if applied
	if calc.Credits.GreaterThan(decimal.Zero) {
		if err := is.applyCreditsToInvoice(tx, credits, calc.Credits); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to apply credits: %w", err)
		}
	}

	// 11. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 12. Load invoice with line items
	if err := is.db.WithContext(ctx).
		Preload("LineItems").
		Preload("Organization").
		Preload("Subscription").
		First(invoice, "id = ?", invoice.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload invoice: %w", err)
	}

	return invoice, nil
}

// generateInvoiceNumber generates a unique invoice number
func (is *InvoiceService) generateInvoiceNumber(ctx context.Context) (string, error) {
	// Get the latest invoice for the current year
	var count int64
	year := time.Now().Year()
	prefix := fmt.Sprintf("%s%d-", is.config.Invoice.NumberPrefix, year)

	err := is.db.WithContext(ctx).
		Model(&models.Invoice{}).
		Where("invoice_number LIKE ?", prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	// Generate invoice number: INV-2025-001234
	invoiceNumber := fmt.Sprintf("%s%06d", prefix, count+1)
	return invoiceNumber, nil
}

// applyCreditsToInvoice deducts credits and updates their remaining amounts
func (is *InvoiceService) applyCreditsToInvoice(
	tx *gorm.DB,
	credits []models.Credit,
	totalCreditApplied decimal.Decimal,
) error {
	remainingToApply := totalCreditApplied

	for i := range credits {
		if remainingToApply.LessThanOrEqual(decimal.Zero) {
			break
		}

		credit := &credits[i]
		creditToApply := decimal.Min(credit.RemainingAmount, remainingToApply)

		// Update credit
		newRemaining := credit.RemainingAmount.Sub(creditToApply)
		updates := map[string]interface{}{
			"remaining_amount": newRemaining,
		}

		// If exhausted, update status
		if newRemaining.LessThanOrEqual(decimal.Zero) {
			updates["status"] = CreditStatusExhausted
		}

		if err := tx.Model(credit).Updates(updates).Error; err != nil {
			return err
		}

		remainingToApply = remainingToApply.Sub(creditToApply)
	}

	return nil
}

// FinalizeInvoice marks an invoice as finalized and ready for payment
func (is *InvoiceService) FinalizeInvoice(ctx context.Context, invoiceID string) error {
	return is.db.WithContext(ctx).
		Model(&models.Invoice{}).
		Where("id = ?", invoiceID).
		Where("status = ?", InvoiceStatusDraft).
		Update("status", InvoiceStatusOpen).Error
}

// MarkInvoiceAsPaid marks an invoice as paid
func (is *InvoiceService) MarkInvoiceAsPaid(
	ctx context.Context,
	invoiceID string,
	paymentID string,
	paidAmount decimal.Decimal,
) error {
	now := time.Now()

	updates := map[string]interface{}{
		"amount_paid": paidAmount,
		"status":      InvoiceStatusPaid,
		"paid_at":     now,
	}

	return is.db.WithContext(ctx).
		Model(&models.Invoice{}).
		Where("id = ?", invoiceID).
		Updates(updates).Error
}

// VoidInvoice voids an invoice
func (is *InvoiceService) VoidInvoice(ctx context.Context, invoiceID string) error {
	return is.db.WithContext(ctx).
		Model(&models.Invoice{}).
		Where("id = ?", invoiceID).
		Where("status != ?", InvoiceStatusPaid).
		Update("status", InvoiceStatusVoid).Error
}

// GetInvoice retrieves an invoice by ID
func (is *InvoiceService) GetInvoice(ctx context.Context, invoiceID string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := is.db.WithContext(ctx).
		Preload("LineItems").
		Preload("Organization").
		Preload("Subscription").
		First(&invoice, "id = ?", invoiceID).Error

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// ListInvoices retrieves invoices for an organization
func (is *InvoiceService) ListInvoices(
	ctx context.Context,
	organizationID string,
	limit, offset int,
) ([]models.Invoice, error) {
	var invoices []models.Invoice

	query := is.db.WithContext(ctx).
		Preload("LineItems").
		Where("organization_id = ?", organizationID).
		Order("invoice_date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&invoices).Error
	return invoices, err
}

// GetUpcomingInvoice calculates what the next invoice will look like
func (is *InvoiceService) GetUpcomingInvoice(
	ctx context.Context,
	subscriptionID string,
) (*models.Invoice, error) {
	// 1. Fetch subscription
	var subscription models.Subscription
	if err := is.db.WithContext(ctx).
		Preload("Plan").
		Preload("Organization").
		First(&subscription, "id = ?", subscriptionID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	// 2. Fetch current usage (estimated)
	usage, err := is.metricsCollector.GetUsageForPeriod(
		ctx,
		subscription.OrganizationID.String(),
		subscription.CurrentPeriodStart,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage: %w", err)
	}

	// 3. Fetch available credits
	var credits []models.Credit
	if err := is.db.WithContext(ctx).
		Where("organization_id = ?", subscription.OrganizationID).
		Where("status = ?", CreditStatusActive).
		Where("remaining_amount > 0").
		Find(&credits).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch credits: %w", err)
	}

	// 4. Calculate charges
	calc, err := is.pricingEngine.CalculateSubscriptionCharge(
		&subscription,
		&subscription.Plan,
		usage,
		credits,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate charges: %w", err)
	}

	// 5. Build preview invoice (not saved to database)
	invoice := &models.Invoice{
		OrganizationID: subscription.OrganizationID,
		SubscriptionID: subscription.ID,
		InvoiceNumber:  "UPCOMING",
		PeriodStart:    subscription.CurrentPeriodStart,
		PeriodEnd:      subscription.CurrentPeriodEnd,
		Subtotal:       calc.Subtotal,
		TaxAmount:      calc.TaxAmount,
		TotalAmount:    calc.Total,
		AmountDue:      calc.Total,
		Currency:       subscription.Plan.Currency,
		Status:         string(InvoiceStatusDraft),
		InvoiceDate:    subscription.CurrentPeriodEnd,
		DueDate:        subscription.CurrentPeriodEnd.AddDate(0, 0, is.config.Invoice.DueDays),
		Organization:   subscription.Organization,
		Subscription:   subscription,
	}

	// Convert line items
	for _, lineItem := range calc.LineItems {
		invoice.LineItems = append(invoice.LineItems, models.InvoiceLineItem{
			Description: lineItem.Description,
			Quantity:    lineItem.Quantity,
			UnitPrice:   lineItem.UnitPrice,
			Amount:      lineItem.Amount,
			ItemType:    string(lineItem.ItemType),
			MetricType:  string(lineItem.MetricType),
			PeriodStart: lineItem.PeriodStart,
			PeriodEnd:   lineItem.PeriodEnd,
		})
	}

	return invoice, nil
}

// ProcessOverdueInvoices marks overdue invoices and triggers notifications
func (is *InvoiceService) ProcessOverdueInvoices(ctx context.Context) error {
	now := time.Now()

	// Find invoices that are past due
	var overdueInvoices []models.Invoice
	err := is.db.WithContext(ctx).
		Where("status = ?", InvoiceStatusOpen).
		Where("due_date < ?", now).
		Find(&overdueInvoices).Error

	if err != nil {
		return fmt.Errorf("failed to fetch overdue invoices: %w", err)
	}

	for _, invoice := range overdueInvoices {
		// Update status (in a real implementation, you might have different overdue statuses)
		// For now, we'll just trigger an event for notification
		// The invoice remains "open" but we can track it's overdue by comparing due_date

		// TODO: Publish event for overdue invoice
		// eventBus.Publish(EventInvoiceOverdue, invoice)
	}

	return nil
}
