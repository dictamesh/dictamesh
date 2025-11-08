// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"fmt"
	"time"

	"github.com/Click2-Run/dictamesh/pkg/billing/models"
	"github.com/shopspring/decimal"
)

// PricingEngine handles all pricing calculations
type PricingEngine struct {
	config *Config
}

// NewPricingEngine creates a new pricing engine
func NewPricingEngine(config *Config) *PricingEngine {
	return &PricingEngine{
		config: config,
	}
}

// CalculateSubscriptionCharge calculates the charge for a subscription period
func (pe *PricingEngine) CalculateSubscriptionCharge(
	subscription *models.Subscription,
	plan *models.SubscriptionPlan,
	usage *UsageAggregation,
	credits []models.Credit,
) (*ChargeCalculation, error) {
	calc := &ChargeCalculation{
		UsageCharges: make(map[MetricType]decimal.Decimal),
		LineItems:    []InvoiceLineItem{},
	}

	// 1. Base subscription charge
	baseCharge := plan.BasePrice.Mul(decimal.NewFromInt(int64(subscription.Quantity)))
	calc.BaseCharge = baseCharge
	calc.LineItems = append(calc.LineItems, InvoiceLineItem{
		Description: fmt.Sprintf("%s Plan (%s)", plan.Name, subscription.CurrentPeriodStart.Format("Jan 2006")),
		Quantity:    decimal.NewFromInt(int64(subscription.Quantity)),
		UnitPrice:   plan.BasePrice,
		Amount:      baseCharge,
		ItemType:    LineItemTypeSubscriptionBase,
		PeriodStart: &subscription.CurrentPeriodStart,
		PeriodEnd:   &subscription.CurrentPeriodEnd,
	})

	// 2. Usage-based charges
	if usage != nil && pe.config.Features.EnableUsageMetrics {
		// API Calls
		if apiCallsCharge, lineItem := pe.calculateUsageCharge(
			MetricTypeAPICalls,
			usage.Metrics[MetricTypeAPICalls],
			decimal.NewFromInt(int64(plan.IncludedAPICalls)),
			plan.PricePerAPICall,
			"API Call",
			subscription.CurrentPeriodStart,
			subscription.CurrentPeriodEnd,
		); apiCallsCharge.GreaterThan(decimal.Zero) {
			calc.UsageCharges[MetricTypeAPICalls] = apiCallsCharge
			calc.LineItems = append(calc.LineItems, lineItem)
		}

		// Storage
		if storageCharge, lineItem := pe.calculateUsageCharge(
			MetricTypeStorageGB,
			usage.Metrics[MetricTypeStorageGB],
			decimal.NewFromInt(int64(plan.IncludedStorageGB)),
			plan.PricePerGBStorage,
			"GB Storage",
			subscription.CurrentPeriodStart,
			subscription.CurrentPeriodEnd,
		); storageCharge.GreaterThan(decimal.Zero) {
			calc.UsageCharges[MetricTypeStorageGB] = storageCharge
			calc.LineItems = append(calc.LineItems, lineItem)
		}

		// Data Transfer
		totalTransfer := usage.Metrics[MetricTypeTransferGBIn].Add(usage.Metrics[MetricTypeTransferGBOut])
		if transferCharge, lineItem := pe.calculateUsageCharge(
			MetricTypeTransferGBOut,
			totalTransfer,
			decimal.NewFromInt(int64(plan.IncludedDataTransferGB)),
			plan.PricePerGBTransfer,
			"GB Data Transfer",
			subscription.CurrentPeriodStart,
			subscription.CurrentPeriodEnd,
		); transferCharge.GreaterThan(decimal.Zero) {
			calc.UsageCharges[MetricTypeTransferGBOut] = transferCharge
			calc.LineItems = append(calc.LineItems, lineItem)
		}
	}

	// 3. Add-on charges (additional seats)
	if subscription.Quantity > plan.IncludedSeats {
		additionalSeats := subscription.Quantity - plan.IncludedSeats
		addonCharge := plan.PricePerAdditionalSeat.Mul(decimal.NewFromInt(int64(additionalSeats)))
		calc.AddonCharges = calc.AddonCharges.Add(addonCharge)
		calc.LineItems = append(calc.LineItems, InvoiceLineItem{
			Description: fmt.Sprintf("Additional Seats (%d)", additionalSeats),
			Quantity:    decimal.NewFromInt(int64(additionalSeats)),
			UnitPrice:   plan.PricePerAdditionalSeat,
			Amount:      addonCharge,
			ItemType:    LineItemTypeAddonSeats,
			PeriodStart: &subscription.CurrentPeriodStart,
			PeriodEnd:   &subscription.CurrentPeriodEnd,
		})
	}

	// 4. Calculate subtotal
	calc.Subtotal = calc.BaseCharge.Add(calc.AddonCharges)
	for _, charge := range calc.UsageCharges {
		calc.Subtotal = calc.Subtotal.Add(charge)
	}

	// 5. Apply credits
	if pe.config.Features.EnableCredits {
		creditAmount := pe.applyCredits(credits, calc.Subtotal)
		if creditAmount.GreaterThan(decimal.Zero) {
			calc.Credits = creditAmount
			calc.LineItems = append(calc.LineItems, InvoiceLineItem{
				Description: "Account Credit Applied",
				Quantity:    decimal.NewFromInt(1),
				UnitPrice:   creditAmount.Neg(),
				Amount:      creditAmount.Neg(),
				ItemType:    LineItemTypeCredit,
			})
		}
	}

	// 6. Calculate tax
	taxableAmount := calc.Subtotal.Sub(calc.Credits)
	if taxableAmount.GreaterThan(decimal.Zero) {
		calc.TaxAmount = taxableAmount.Mul(pe.config.Invoice.TaxRate)
		if calc.TaxAmount.GreaterThan(decimal.Zero) {
			calc.LineItems = append(calc.LineItems, InvoiceLineItem{
				Description: fmt.Sprintf("Tax (%s%%)", pe.config.Invoice.TaxRate.Mul(decimal.NewFromInt(100)).String()),
				Quantity:    decimal.NewFromInt(1),
				UnitPrice:   calc.TaxAmount,
				Amount:      calc.TaxAmount,
				ItemType:    LineItemTypeTax,
			})
		}
	}

	// 7. Calculate total
	calc.Total = calc.Subtotal.Sub(calc.Credits).Add(calc.TaxAmount)

	return calc, nil
}

// calculateUsageCharge calculates the charge for a single usage metric
func (pe *PricingEngine) calculateUsageCharge(
	metricType MetricType,
	actualUsage, includedAmount, pricePerUnit decimal.Decimal,
	unitName string,
	periodStart, periodEnd time.Time,
) (decimal.Decimal, InvoiceLineItem) {
	charge := decimal.Zero
	lineItem := InvoiceLineItem{
		MetricType:  metricType,
		PeriodStart: &periodStart,
		PeriodEnd:   &periodEnd,
	}

	// Calculate overage
	overage := actualUsage.Sub(includedAmount)
	if overage.LessThanOrEqual(decimal.Zero) {
		return charge, lineItem
	}

	// Apply pricing
	charge = overage.Mul(pricePerUnit)

	// Round to 2 decimal places
	charge = charge.Round(2)

	// Build line item
	lineItem.Description = fmt.Sprintf("%s\n  Included: %s %s\n  Usage: %s %s\n  Overage: %s %s",
		unitName,
		includedAmount.StringFixed(0), unitName,
		actualUsage.StringFixed(2), unitName,
		overage.StringFixed(2), unitName,
	)
	lineItem.Quantity = overage
	lineItem.UnitPrice = pricePerUnit
	lineItem.Amount = charge

	// Set item type based on metric
	switch metricType {
	case MetricTypeAPICalls:
		lineItem.ItemType = LineItemTypeUsageAPICalls
	case MetricTypeStorageGB:
		lineItem.ItemType = LineItemTypeUsageStorage
	case MetricTypeTransferGBIn, MetricTypeTransferGBOut:
		lineItem.ItemType = LineItemTypeUsageTransfer
	}

	return charge, lineItem
}

// CalculateTieredPrice calculates price using volume-based tiers
func (pe *PricingEngine) CalculateTieredPrice(
	usage decimal.Decimal,
	tiers []models.PricingTier,
) decimal.Decimal {
	if !pe.config.Features.EnableTieredPricing || len(tiers) == 0 {
		return decimal.Zero
	}

	totalCharge := decimal.Zero
	remainingUsage := usage

	for _, tier := range tiers {
		if remainingUsage.LessThanOrEqual(decimal.Zero) {
			break
		}

		// Determine tier capacity
		var tierCapacity decimal.Decimal
		if tier.TierEnd == nil {
			// Infinite tier
			tierCapacity = remainingUsage
		} else {
			tierCapacity = tier.TierEnd.Sub(tier.TierStart)
		}

		// Calculate usage in this tier
		tierUsage := decimal.Min(remainingUsage, tierCapacity)

		// Calculate charge for this tier
		tierCharge := tierUsage.Mul(tier.PricePerUnit).Add(tier.FlatFee)
		totalCharge = totalCharge.Add(tierCharge)

		// Reduce remaining usage
		remainingUsage = remainingUsage.Sub(tierUsage)
	}

	return totalCharge.Round(2)
}

// applyCredits applies available credits to the charge
func (pe *PricingEngine) applyCredits(credits []models.Credit, amount decimal.Decimal) decimal.Decimal {
	appliedCredit := decimal.Zero
	remainingAmount := amount

	for _, credit := range credits {
		if remainingAmount.LessThanOrEqual(decimal.Zero) {
			break
		}

		if credit.Status != string(CreditStatusActive) {
			continue
		}

		// Check if credit is still valid
		now := time.Now()
		if credit.ValidUntil != nil && credit.ValidUntil.Before(now) {
			continue
		}

		// Apply credit
		creditToApply := decimal.Min(credit.RemainingAmount, remainingAmount)
		appliedCredit = appliedCredit.Add(creditToApply)
		remainingAmount = remainingAmount.Sub(creditToApply)
	}

	return appliedCredit
}

// CalculateProration calculates prorated charges for mid-cycle changes
func (pe *PricingEngine) CalculateProration(
	oldPrice, newPrice decimal.Decimal,
	periodStart, periodEnd, changeDate time.Time,
) decimal.Decimal {
	if !pe.config.Features.EnableProration {
		return decimal.Zero
	}

	totalSeconds := periodEnd.Sub(periodStart).Seconds()
	remainingSeconds := periodEnd.Sub(changeDate).Seconds()

	if totalSeconds <= 0 {
		return decimal.Zero
	}

	// Calculate prorated portion
	priceDiff := newPrice.Sub(oldPrice)
	proratedRatio := decimal.NewFromFloat(remainingSeconds / totalSeconds)
	proration := priceDiff.Mul(proratedRatio)

	return proration.Round(2)
}

// EstimateMonthlyCharge estimates the monthly charge for a subscription
func (pe *PricingEngine) EstimateMonthlyCharge(
	plan *models.SubscriptionPlan,
	quantity int,
	estimatedUsage map[MetricType]decimal.Decimal,
) decimal.Decimal {
	// Base charge
	estimate := plan.BasePrice.Mul(decimal.NewFromInt(int64(quantity)))

	// Additional seats
	if quantity > plan.IncludedSeats {
		additionalSeats := quantity - plan.IncludedSeats
		estimate = estimate.Add(plan.PricePerAdditionalSeat.Mul(decimal.NewFromInt(int64(additionalSeats))))
	}

	// Usage estimates
	if pe.config.Features.EnableUsageMetrics {
		// API calls
		if apiCalls, ok := estimatedUsage[MetricTypeAPICalls]; ok {
			overage := apiCalls.Sub(decimal.NewFromInt(int64(plan.IncludedAPICalls)))
			if overage.GreaterThan(decimal.Zero) {
				estimate = estimate.Add(overage.Mul(plan.PricePerAPICall))
			}
		}

		// Storage
		if storage, ok := estimatedUsage[MetricTypeStorageGB]; ok {
			overage := storage.Sub(decimal.NewFromInt(int64(plan.IncludedStorageGB)))
			if overage.GreaterThan(decimal.Zero) {
				estimate = estimate.Add(overage.Mul(plan.PricePerGBStorage))
			}
		}

		// Transfer
		if transfer, ok := estimatedUsage[MetricTypeTransferGBOut]; ok {
			overage := transfer.Sub(decimal.NewFromInt(int64(plan.IncludedDataTransferGB)))
			if overage.GreaterThan(decimal.Zero) {
				estimate = estimate.Add(overage.Mul(plan.PricePerGBTransfer))
			}
		}
	}

	return estimate.Round(2)
}
