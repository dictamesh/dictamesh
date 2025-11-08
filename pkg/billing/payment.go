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
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/paymentintent"
	"github.com/stripe/stripe-go/v75/paymentmethod"
	"gorm.io/gorm"
)

// PaymentService handles payment processing
type PaymentService struct {
	db             *gorm.DB
	config         *Config
	invoiceService *InvoiceService
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	db *gorm.DB,
	config *Config,
	invoiceService *InvoiceService,
) *PaymentService {
	// Initialize Stripe
	if config.Stripe.Enabled {
		stripe.Key = config.Stripe.APIKey
	}

	return &PaymentService{
		db:             db,
		config:         config,
		invoiceService: invoiceService,
	}
}

// CreateStripeCustomer creates a Stripe customer for an organization
func (ps *PaymentService) CreateStripeCustomer(
	ctx context.Context,
	org *models.Organization,
) (string, error) {
	if !ps.config.Stripe.Enabled {
		return "", fmt.Errorf("Stripe is not enabled")
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(org.BillingEmail),
		Name:  stripe.String(org.Name),
		Metadata: map[string]string{
			"organization_id": org.ID.String(),
		},
	}

	if org.AddressLine1 != "" {
		params.Address = &stripe.AddressParams{
			Line1:      stripe.String(org.AddressLine1),
			Line2:      stripe.String(org.AddressLine2),
			City:       stripe.String(org.City),
			State:      stripe.String(org.State),
			PostalCode: stripe.String(org.PostalCode),
			Country:    stripe.String(org.Country),
		}
	}

	cust, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Update organization with Stripe customer ID
	if err := ps.db.WithContext(ctx).
		Model(org).
		Update("stripe_customer_id", cust.ID).Error; err != nil {
		return "", fmt.Errorf("failed to update organization: %w", err)
	}

	return cust.ID, nil
}

// AttachPaymentMethod attaches a payment method to a customer
func (ps *PaymentService) AttachPaymentMethod(
	ctx context.Context,
	organizationID, paymentMethodID string,
	setAsDefault bool,
) error {
	if !ps.config.Stripe.Enabled {
		return fmt.Errorf("Stripe is not enabled")
	}

	// Fetch organization
	var org models.Organization
	if err := ps.db.WithContext(ctx).First(&org, "id = ?", organizationID).Error; err != nil {
		return fmt.Errorf("failed to fetch organization: %w", err)
	}

	// Ensure organization has a Stripe customer
	if org.StripeCustomerID == "" {
		customerID, err := ps.CreateStripeCustomer(ctx, &org)
		if err != nil {
			return err
		}
		org.StripeCustomerID = customerID
	}

	// Attach payment method to customer
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(org.StripeCustomerID),
	}

	if _, err := paymentmethod.Attach(paymentMethodID, params); err != nil {
		return fmt.Errorf("failed to attach payment method: %w", err)
	}

	// Set as default if requested
	if setAsDefault {
		updateParams := &stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(paymentMethodID),
			},
		}

		if _, err := customer.Update(org.StripeCustomerID, updateParams); err != nil {
			return fmt.Errorf("failed to set default payment method: %w", err)
		}

		// Update organization
		if err := ps.db.WithContext(ctx).
			Model(&org).
			Update("default_payment_method_id", paymentMethodID).Error; err != nil {
			return fmt.Errorf("failed to update organization: %w", err)
		}
	}

	return nil
}

// ChargeInvoice charges a payment method for an invoice
func (ps *PaymentService) ChargeInvoice(
	ctx context.Context,
	invoiceID string,
) (*models.Payment, error) {
	// 1. Fetch invoice
	invoice, err := ps.invoiceService.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoice: %w", err)
	}

	// 2. Check if already paid
	if invoice.Status == string(InvoiceStatusPaid) {
		return nil, fmt.Errorf("invoice already paid")
	}

	// 3. Fetch organization
	var org models.Organization
	if err := ps.db.WithContext(ctx).First(&org, "id = ?", invoice.OrganizationID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch organization: %w", err)
	}

	// 4. Create payment record
	payment := &models.Payment{
		ID:             uuid.New(),
		OrganizationID: invoice.OrganizationID,
		InvoiceID:      invoice.ID,
		Amount:         invoice.AmountDue,
		Currency:       invoice.Currency,
		Status:         string(PaymentStatusPending),
		Provider:       string(PaymentProviderStripe),
		PaymentMethodID: org.DefaultPaymentMethodID,
		ProviderCustomerID: org.StripeCustomerID,
	}

	// Save payment record
	if err := ps.db.WithContext(ctx).Create(payment).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// 5. Process payment with Stripe
	if ps.config.Stripe.Enabled {
		if err := ps.processStripePayment(ctx, payment, invoice, &org); err != nil {
			// Update payment as failed
			now := time.Now()
			ps.db.WithContext(ctx).Model(payment).Updates(map[string]interface{}{
				"status":         PaymentStatusFailed,
				"failed_at":      now,
				"failure_message": err.Error(),
			})
			return payment, err
		}
	}

	// 6. Reload payment with updates
	if err := ps.db.WithContext(ctx).First(payment, "id = ?", payment.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload payment: %w", err)
	}

	return payment, nil
}

// processStripePayment processes a payment via Stripe
func (ps *PaymentService) processStripePayment(
	ctx context.Context,
	payment *models.Payment,
	invoice *models.Invoice,
	org *models.Organization,
) error {
	// Convert amount to cents
	amountCents := payment.Amount.Mul(decimal.NewFromInt(100)).IntPart()

	// Create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(invoice.Currency),
		Customer: stripe.String(org.StripeCustomerID),
		PaymentMethod: stripe.String(payment.PaymentMethodID),
		Confirm: stripe.Bool(true), // Automatically confirm
		OffSession: stripe.Bool(true), // For subscription billing
		Metadata: map[string]string{
			"invoice_id":      invoice.ID.String(),
			"organization_id": org.ID.String(),
			"payment_id":      payment.ID.String(),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Update payment record
	now := time.Now()
	updates := map[string]interface{}{
		"provider_payment_id": pi.ID,
		"attempted_at":        now,
	}

	if pi.Status == stripe.PaymentIntentStatusSucceeded {
		updates["status"] = PaymentStatusSucceeded
		updates["succeeded_at"] = now

		// Mark invoice as paid
		if err := ps.invoiceService.MarkInvoiceAsPaid(ctx, invoice.ID.String(), payment.ID.String(), payment.Amount); err != nil {
			return fmt.Errorf("failed to mark invoice as paid: %w", err)
		}
	} else if pi.Status == stripe.PaymentIntentStatusRequiresAction ||
		pi.Status == stripe.PaymentIntentStatusRequiresPaymentMethod {
		updates["status"] = PaymentStatusPending
	} else {
		updates["status"] = PaymentStatusFailed
		updates["failed_at"] = now
		if pi.LastPaymentError != nil {
			updates["failure_code"] = pi.LastPaymentError.Code
			updates["failure_message"] = pi.LastPaymentError.Message
		}
	}

	return ps.db.WithContext(ctx).Model(payment).Updates(updates).Error
}

// HandleWebhook processes payment provider webhooks
func (ps *PaymentService) HandleWebhook(
	ctx context.Context,
	provider PaymentProvider,
	eventType string,
	payload map[string]interface{},
) error {
	switch provider {
	case PaymentProviderStripe:
		return ps.handleStripeWebhook(ctx, eventType, payload)
	default:
		return fmt.Errorf("unsupported payment provider: %s", provider)
	}
}

// handleStripeWebhook handles Stripe webhook events
func (ps *PaymentService) handleStripeWebhook(
	ctx context.Context,
	eventType string,
	payload map[string]interface{},
) error {
	switch eventType {
	case "payment_intent.succeeded":
		return ps.handlePaymentIntentSucceeded(ctx, payload)
	case "payment_intent.payment_failed":
		return ps.handlePaymentIntentFailed(ctx, payload)
	case "customer.subscription.updated":
		// Handle subscription updates
		return nil
	case "invoice.payment_succeeded":
		// Handle invoice payment success
		return nil
	default:
		// Unknown event type, ignore
		return nil
	}
}

// handlePaymentIntentSucceeded handles successful payment intents
func (ps *PaymentService) handlePaymentIntentSucceeded(
	ctx context.Context,
	payload map[string]interface{},
) error {
	// Extract payment intent ID
	paymentIntentID, ok := payload["id"].(string)
	if !ok {
		return fmt.Errorf("invalid payment intent ID")
	}

	// Find payment by provider payment ID
	var payment models.Payment
	if err := ps.db.WithContext(ctx).
		Where("provider_payment_id = ?", paymentIntentID).
		First(&payment).Error; err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Update payment status
	now := time.Now()
	updates := map[string]interface{}{
		"status":       PaymentStatusSucceeded,
		"succeeded_at": now,
	}

	if err := ps.db.WithContext(ctx).Model(&payment).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Mark invoice as paid
	if payment.InvoiceID.String() != "" {
		if err := ps.invoiceService.MarkInvoiceAsPaid(
			ctx,
			payment.InvoiceID.String(),
			payment.ID.String(),
			payment.Amount,
		); err != nil {
			return fmt.Errorf("failed to mark invoice as paid: %w", err)
		}
	}

	// TODO: Publish event
	// eventBus.Publish(EventPaymentSucceeded, payment)

	return nil
}

// handlePaymentIntentFailed handles failed payment intents
func (ps *PaymentService) handlePaymentIntentFailed(
	ctx context.Context,
	payload map[string]interface{},
) error {
	// Extract payment intent ID
	paymentIntentID, ok := payload["id"].(string)
	if !ok {
		return fmt.Errorf("invalid payment intent ID")
	}

	// Find payment by provider payment ID
	var payment models.Payment
	if err := ps.db.WithContext(ctx).
		Where("provider_payment_id = ?", paymentIntentID).
		First(&payment).Error; err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Extract failure reason
	var failureMessage string
	if lastError, ok := payload["last_payment_error"].(map[string]interface{}); ok {
		if msg, ok := lastError["message"].(string); ok {
			failureMessage = msg
		}
	}

	// Update payment status
	now := time.Now()
	updates := map[string]interface{}{
		"status":          PaymentStatusFailed,
		"failed_at":       now,
		"failure_message": failureMessage,
	}

	if err := ps.db.WithContext(ctx).Model(&payment).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// TODO: Publish event
	// eventBus.Publish(EventPaymentFailed, payment)

	return nil
}

// ListPayments retrieves payments for an organization
func (ps *PaymentService) ListPayments(
	ctx context.Context,
	organizationID string,
	limit, offset int,
) ([]models.Payment, error) {
	var payments []models.Payment

	query := ps.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&payments).Error
	return payments, err
}

// RefundPayment refunds a payment
func (ps *PaymentService) RefundPayment(
	ctx context.Context,
	paymentID string,
	amount *decimal.Decimal,
) error {
	// Fetch payment
	var payment models.Payment
	if err := ps.db.WithContext(ctx).First(&payment, "id = ?", paymentID).Error; err != nil {
		return fmt.Errorf("failed to fetch payment: %w", err)
	}

	if payment.Status != string(PaymentStatusSucceeded) {
		return fmt.Errorf("can only refund succeeded payments")
	}

	// Determine refund amount
	refundAmount := payment.Amount
	if amount != nil {
		refundAmount = *amount
	}

	if refundAmount.GreaterThan(payment.Amount) {
		return fmt.Errorf("refund amount cannot exceed payment amount")
	}

	// TODO: Process refund with Stripe
	// For now, just update the status

	now := time.Now()
	updates := map[string]interface{}{
		"status":      PaymentStatusRefunded,
		"refunded_at": now,
	}

	return ps.db.WithContext(ctx).Model(&payment).Updates(updates).Error
}
