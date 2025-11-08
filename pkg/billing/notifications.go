// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Click2-Run/dictamesh/pkg/billing/models"
)

// NotificationService handles sending billing-related notifications
type NotificationService struct {
	config *Config
	client *http.Client
}

// NewNotificationService creates a new notification service
func NewNotificationService(config *Config) *NotificationService {
	return &NotificationService{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Notifications.TimeoutSeconds) * time.Second,
		},
	}
}

// NotificationRequest represents a request to the notification service
type NotificationRequest struct {
	RecipientID   string                 `json:"recipient_id"`
	RecipientType string                 `json:"recipient_type"`
	TemplateCode  string                 `json:"template_code"`
	Channels      []string               `json:"channels"`
	Priority      string                 `json:"priority"`
	Data          map[string]interface{} `json:"data"`
}

// SendInvoiceCreatedNotification sends notification when invoice is created
func (ns *NotificationService) SendInvoiceCreatedNotification(
	ctx context.Context,
	invoice *models.Invoice,
) error {
	data := map[string]interface{}{
		"InvoiceNumber":    invoice.InvoiceNumber,
		"OrganizationName": invoice.Organization.Name,
		"PeriodStart":      invoice.PeriodStart.Format("Jan 2, 2006"),
		"PeriodEnd":        invoice.PeriodEnd.Format("Jan 2, 2006"),
		"Subtotal":         invoice.Subtotal.StringFixed(2),
		"Tax":              invoice.TaxAmount.StringFixed(2),
		"Total":            invoice.TotalAmount.StringFixed(2),
		"DueDate":          invoice.DueDate.Format("Jan 2, 2006"),
		"InvoiceURL":       fmt.Sprintf("https://app.dictamesh.io/invoices/%s", invoice.ID),
		"AutoPay":          invoice.Organization.AutoPay,
		"Currency":         invoice.Currency,
	}

	notification := &NotificationRequest{
		RecipientID:   invoice.OrganizationID.String(),
		RecipientType: "organization",
		TemplateCode:  "billing_invoice_generated",
		Channels:      []string{"email"},
		Priority:      "high",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// SendPaymentSucceededNotification sends notification when payment succeeds
func (ns *NotificationService) SendPaymentSucceededNotification(
	ctx context.Context,
	payment *models.Payment,
	invoice *models.Invoice,
) error {
	data := map[string]interface{}{
		"InvoiceNumber":  invoice.InvoiceNumber,
		"Amount":         payment.Amount.StringFixed(2),
		"Currency":       payment.Currency,
		"PaymentMethod":  payment.PaymentMethod,
		"TransactionID":  payment.ProviderPaymentID,
		"PaymentDate":    payment.SucceededAt.Format("Jan 2, 2006"),
		"ReceiptURL":     fmt.Sprintf("https://app.dictamesh.io/payments/%s/receipt", payment.ID),
	}

	notification := &NotificationRequest{
		RecipientID:   payment.OrganizationID.String(),
		RecipientType: "organization",
		TemplateCode:  "billing_payment_succeeded",
		Channels:      []string{"email"},
		Priority:      "normal",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// SendPaymentFailedNotification sends notification when payment fails
func (ns *NotificationService) SendPaymentFailedNotification(
	ctx context.Context,
	payment *models.Payment,
	invoice *models.Invoice,
) error {
	data := map[string]interface{}{
		"InvoiceNumber": invoice.InvoiceNumber,
		"Amount":        payment.Amount.StringFixed(2),
		"Currency":      payment.Currency,
		"FailureReason": payment.FailureMessage,
		"FailureCode":   payment.FailureCode,
		"PaymentURL":    fmt.Sprintf("https://app.dictamesh.io/invoices/%s/pay", invoice.ID),
		"DueDate":       invoice.DueDate.Format("Jan 2, 2006"),
	}

	notification := &NotificationRequest{
		RecipientID:   payment.OrganizationID.String(),
		RecipientType: "organization",
		TemplateCode:  "billing_payment_failed",
		Channels:      []string{"email"},
		Priority:      "urgent",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// SendInvoiceOverdueNotification sends notification when invoice is overdue
func (ns *NotificationService) SendInvoiceOverdueNotification(
	ctx context.Context,
	invoice *models.Invoice,
) error {
	daysOverdue := int(time.Since(invoice.DueDate).Hours() / 24)

	data := map[string]interface{}{
		"InvoiceNumber": invoice.InvoiceNumber,
		"Amount":        invoice.AmountDue.StringFixed(2),
		"Currency":      invoice.Currency,
		"DueDate":       invoice.DueDate.Format("Jan 2, 2006"),
		"DaysOverdue":   daysOverdue,
		"PaymentURL":    fmt.Sprintf("https://app.dictamesh.io/invoices/%s/pay", invoice.ID),
	}

	notification := &NotificationRequest{
		RecipientID:   invoice.OrganizationID.String(),
		RecipientType: "organization",
		TemplateCode:  "billing_invoice_overdue",
		Channels:      []string{"email"},
		Priority:      "urgent",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// SendSubscriptionCreatedNotification sends notification when subscription is created
func (ns *NotificationService) SendSubscriptionCreatedNotification(
	ctx context.Context,
	subscription *models.Subscription,
) error {
	data := map[string]interface{}{
		"PlanName":            subscription.Plan.Name,
		"BillingCycle":        subscription.Plan.BillingInterval,
		"Amount":              subscription.Plan.BasePrice.StringFixed(2),
		"Currency":            subscription.Plan.Currency,
		"CurrentPeriodStart":  subscription.CurrentPeriodStart.Format("Jan 2, 2006"),
		"CurrentPeriodEnd":    subscription.CurrentPeriodEnd.Format("Jan 2, 2006"),
		"SubscriptionURL":     fmt.Sprintf("https://app.dictamesh.io/subscriptions/%s", subscription.ID),
	}

	notification := &NotificationRequest{
		RecipientID:   subscription.OrganizationID.String(),
		RecipientType: "organization",
		TemplateCode:  "billing_subscription_created",
		Channels:      []string{"email"},
		Priority:      "normal",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// SendSubscriptionCanceledNotification sends notification when subscription is canceled
func (ns *NotificationService) SendSubscriptionCanceledNotification(
	ctx context.Context,
	subscription *models.Subscription,
) error {
	data := map[string]interface{}{
		"PlanName":         subscription.Plan.Name,
		"CancellationDate": subscription.CanceledAt.Format("Jan 2, 2006"),
		"EndDate":          subscription.CurrentPeriodEnd.Format("Jan 2, 2006"),
		"Reason":           subscription.CancellationReason,
	}

	notification := &NotificationRequest{
		RecipientID:   subscription.OrganizationID.String(),
		RecipientType: "organization",
		TemplateCode:  "billing_subscription_canceled",
		Channels:      []string{"email"},
		Priority:      "normal",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// SendUsageThresholdNotification sends notification when usage threshold is reached
func (ns *NotificationService) SendUsageThresholdNotification(
	ctx context.Context,
	organizationID string,
	metricType MetricType,
	currentUsage, threshold string,
	percentUsed int,
) error {
	data := map[string]interface{}{
		"MetricType":   metricType,
		"CurrentUsage": currentUsage,
		"Threshold":    threshold,
		"PercentUsed":  percentUsed,
		"UsageURL":     "https://app.dictamesh.io/usage",
	}

	notification := &NotificationRequest{
		RecipientID:   organizationID,
		RecipientType: "organization",
		TemplateCode:  "billing_usage_threshold_reached",
		Channels:      []string{"email"},
		Priority:      "normal",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// SendUpcomingRenewalNotification sends notification before subscription renewal
func (ns *NotificationService) SendUpcomingRenewalNotification(
	ctx context.Context,
	subscription *models.Subscription,
	upcomingInvoice *models.Invoice,
) error {
	daysUntilRenewal := int(time.Until(subscription.CurrentPeriodEnd).Hours() / 24)

	data := map[string]interface{}{
		"PlanName":        subscription.Plan.Name,
		"RenewalDate":     subscription.CurrentPeriodEnd.Format("Jan 2, 2006"),
		"DaysUntilRenewal": daysUntilRenewal,
		"Amount":          upcomingInvoice.TotalAmount.StringFixed(2),
		"Currency":        upcomingInvoice.Currency,
		"InvoiceURL":      fmt.Sprintf("https://app.dictamesh.io/invoices/upcoming"),
	}

	notification := &NotificationRequest{
		RecipientID:   subscription.OrganizationID.String(),
		RecipientType: "organization",
		TemplateCode:  "billing_upcoming_renewal",
		Channels:      []string{"email"},
		Priority:      "normal",
		Data:          data,
	}

	return ns.sendNotification(ctx, notification)
}

// sendNotification sends a notification request to the notification service
func (ns *NotificationService) sendNotification(
	ctx context.Context,
	notification *NotificationRequest,
) error {
	// Marshal notification to JSON
	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Build request URL
	url := fmt.Sprintf("%s/api/v1/notifications", ns.config.Notifications.ServiceURL)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request with retries
	var lastErr error
	for i := 0; i < ns.config.Notifications.RetryAttempts; i++ {
		resp, err := ns.client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(ns.config.Notifications.RetryDelay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		lastErr = fmt.Errorf("notification service returned status %d", resp.StatusCode)
		time.Sleep(ns.config.Notifications.RetryDelay)
	}

	return fmt.Errorf("failed to send notification after %d attempts: %w",
		ns.config.Notifications.RetryAttempts, lastErr)
}

// CreateBillingTemplates creates email templates for billing notifications
// This function should be called during system initialization
func (ns *NotificationService) CreateBillingTemplates(ctx context.Context) error {
	templates := []map[string]interface{}{
		{
			"template_code": "billing_invoice_generated",
			"name":          "Invoice Generated",
			"description":   "Sent when a new invoice is generated",
			"channels":      []string{"email"},
			"subject":       "Your DictaMesh Invoice #{{.InvoiceNumber}}",
			"body_html":     getInvoiceGeneratedTemplate(),
		},
		{
			"template_code": "billing_payment_succeeded",
			"name":          "Payment Received",
			"description":   "Sent when a payment is successfully processed",
			"channels":      []string{"email"},
			"subject":       "Payment Received - Invoice #{{.InvoiceNumber}}",
			"body_html":     getPaymentSucceededTemplate(),
		},
		{
			"template_code": "billing_payment_failed",
			"name":          "Payment Failed",
			"description":   "Sent when a payment fails",
			"channels":      []string{"email"},
			"subject":       "Action Required: Payment Failed for Invoice #{{.InvoiceNumber}}",
			"body_html":     getPaymentFailedTemplate(),
		},
		{
			"template_code": "billing_invoice_overdue",
			"name":          "Invoice Overdue",
			"description":   "Sent when an invoice becomes overdue",
			"channels":      []string{"email"},
			"subject":       "Overdue Invoice #{{.InvoiceNumber}} - Payment Required",
			"body_html":     getInvoiceOverdueTemplate(),
		},
		{
			"template_code": "billing_subscription_created",
			"name":          "Subscription Created",
			"description":   "Sent when a new subscription is created",
			"channels":      []string{"email"},
			"subject":       "Welcome to {{.PlanName}}!",
			"body_html":     getSubscriptionCreatedTemplate(),
		},
		{
			"template_code": "billing_subscription_canceled",
			"name":          "Subscription Canceled",
			"description":   "Sent when a subscription is canceled",
			"channels":      []string{"email"},
			"subject":       "Your {{.PlanName}} subscription has been canceled",
			"body_html":     getSubscriptionCanceledTemplate(),
		},
		{
			"template_code": "billing_usage_threshold_reached",
			"name":          "Usage Threshold Reached",
			"description":   "Sent when usage reaches a threshold",
			"channels":      []string{"email"},
			"subject":       "Usage Alert: {{.MetricType}} at {{.PercentUsed}}%",
			"body_html":     getUsageThresholdTemplate(),
		},
		{
			"template_code": "billing_upcoming_renewal",
			"name":          "Upcoming Renewal",
			"description":   "Sent before subscription renewal",
			"channels":      []string{"email"},
			"subject":       "Your {{.PlanName}} subscription renews in {{.DaysUntilRenewal}} days",
			"body_html":     getUpcomingRenewalTemplate(),
		},
	}

	// Send each template to the notification service
	url := fmt.Sprintf("%s/api/v1/templates", ns.config.Notifications.ServiceURL)

	for _, template := range templates {
		payload, err := json.Marshal(template)
		if err != nil {
			return fmt.Errorf("failed to marshal template: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := ns.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to create template: %w", err)
		}
		resp.Body.Close()
	}

	return nil
}

// Email template HTML content

func getInvoiceGeneratedTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h1>Invoice #{{.InvoiceNumber}}</h1>
<p>Dear Customer,</p>
<p>Your invoice for the period {{.PeriodStart}} - {{.PeriodEnd}} is ready.</p>
<table border="1" cellpadding="10">
<tr><td>Subtotal:</td><td>{{.Currency}} {{.Subtotal}}</td></tr>
<tr><td>Tax:</td><td>{{.Currency}} {{.Tax}}</td></tr>
<tr><th>Total:</th><th>{{.Currency}} {{.Total}}</th></tr>
</table>
<p><a href="{{.InvoiceURL}}">View Invoice</a></p>
{{if .AutoPay}}
<p>Your payment method will be charged automatically on {{.DueDate}}.</p>
{{else}}
<p>Please pay by {{.DueDate}} to avoid service interruption.</p>
{{end}}
<p>Thank you for using DictaMesh!</p>
</body>
</html>
`
}

func getPaymentSucceededTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h1>Payment Confirmed</h1>
<p>Thank you! We've received your payment of {{.Currency}} {{.Amount}}.</p>
<p>Invoice: #{{.InvoiceNumber}}<br>
Payment Method: {{.PaymentMethod}}<br>
Transaction ID: {{.TransactionID}}<br>
Payment Date: {{.PaymentDate}}</p>
<p><a href="{{.ReceiptURL}}">View Receipt</a></p>
<p>Thank you for your business!</p>
</body>
</html>
`
}

func getPaymentFailedTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;color:#c00;}</style></head>
<body>
<h1>Payment Failed</h1>
<p>We were unable to process your payment for invoice #{{.InvoiceNumber}}.</p>
<p><strong>Reason:</strong> {{.FailureReason}}</p>
<p>Please update your payment method or pay manually to avoid service interruption.</p>
<p><a href="{{.PaymentURL}}">Update Payment Method</a></p>
<p>Due Date: {{.DueDate}}</p>
</body>
</html>
`
}

func getInvoiceOverdueTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h1>Overdue Invoice #{{.InvoiceNumber}}</h1>
<p><strong>Your invoice is {{.DaysOverdue}} days overdue.</strong></p>
<p>Amount Due: {{.Currency}} {{.Amount}}<br>
Due Date: {{.DueDate}}</p>
<p>Please make payment immediately to avoid service suspension.</p>
<p><a href="{{.PaymentURL}}">Pay Now</a></p>
</body>
</html>
`
}

func getSubscriptionCreatedTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h1>Welcome to {{.PlanName}}!</h1>
<p>Your subscription has been successfully created.</p>
<p>Plan: {{.PlanName}}<br>
Billing Cycle: {{.BillingCycle}}<br>
Amount: {{.Currency}} {{.Amount}}<br>
Current Period: {{.CurrentPeriodStart}} - {{.CurrentPeriodEnd}}</p>
<p><a href="{{.SubscriptionURL}}">View Subscription</a></p>
<p>Thank you for choosing DictaMesh!</p>
</body>
</html>
`
}

func getSubscriptionCanceledTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h1>Subscription Canceled</h1>
<p>Your {{.PlanName}} subscription has been canceled.</p>
<p>Cancellation Date: {{.CancellationDate}}<br>
Service End Date: {{.EndDate}}</p>
{{if .Reason}}
<p>Reason: {{.Reason}}</p>
{{end}}
<p>We're sorry to see you go. You can reactivate your subscription at any time.</p>
</body>
</html>
`
}

func getUsageThresholdTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h1>Usage Alert</h1>
<p>You've reached {{.PercentUsed}}% of your {{.MetricType}} limit.</p>
<p>Current Usage: {{.CurrentUsage}}<br>
Limit: {{.Threshold}}</p>
<p><a href="{{.UsageURL}}">View Usage Details</a></p>
<p>Consider upgrading your plan if you need more capacity.</p>
</body>
</html>
`
}

func getUpcomingRenewalTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head><style>body{font-family:Arial,sans-serif;}</style></head>
<body>
<h1>Upcoming Renewal</h1>
<p>Your {{.PlanName}} subscription will renew in {{.DaysUntilRenewal}} days.</p>
<p>Renewal Date: {{.RenewalDate}}<br>
Amount: {{.Currency}} {{.Amount}}</p>
<p><a href="{{.InvoiceURL}}">Preview Upcoming Invoice</a></p>
<p>Your payment method will be charged automatically on the renewal date.</p>
</body>
</html>
`
}
