// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package billing

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// Config represents the billing system configuration
type Config struct {
	// Database
	DatabaseDSN string

	// Stripe configuration
	Stripe StripeConfig

	// PayPal configuration
	PayPal PayPalConfig

	// Invoice settings
	Invoice InvoiceConfig

	// Usage metrics settings
	Usage UsageConfig

	// Notification settings
	Notifications NotificationConfig

	// Feature flags
	Features FeatureFlags

	// Rate limiting
	RateLimits RateLimitConfig
}

// StripeConfig contains Stripe payment provider settings
type StripeConfig struct {
	APIKey        string
	WebhookSecret string
	Enabled       bool
}

// PayPalConfig contains PayPal payment provider settings
type PayPalConfig struct {
	ClientID     string
	ClientSecret string
	Environment  string // sandbox or production
	Enabled      bool
}

// InvoiceConfig contains invoice generation settings
type InvoiceConfig struct {
	DueDays         int             // Number of days until invoice is due
	NumberPrefix    string          // Prefix for invoice numbers (e.g., "INV-")
	TaxRate         decimal.Decimal // Default tax rate (e.g., 0.10 for 10%)
	DefaultCurrency string          // Default currency code (ISO 4217)
	PDFStoragePath  string          // Path to store generated PDF files
}

// UsageConfig contains usage metrics collection settings
type UsageConfig struct {
	AggregationInterval time.Duration // How often to aggregate usage metrics
	RetentionDays       int           // How long to retain detailed usage data
	BatchSize           int           // Batch size for metric processing
	EnableRealTime      bool          // Enable real-time usage tracking
}

// NotificationConfig contains notification integration settings
type NotificationConfig struct {
	ServiceURL     string        // URL of the notification service
	RetryAttempts  int           // Number of retry attempts for failed notifications
	RetryDelay     time.Duration // Delay between retry attempts
	TimeoutSeconds int           // Timeout for notification requests
}

// FeatureFlags controls which features are enabled
type FeatureFlags struct {
	EnableAutoPayment   bool // Automatically charge payment methods
	EnableUsageMetrics  bool // Track and bill for usage
	EnableTieredPricing bool // Support volume-based pricing tiers
	EnableMultiCurrency bool // Support multiple currencies
	EnableCredits       bool // Support account credits
	EnableProration     bool // Prorate charges for mid-cycle changes
}

// RateLimitConfig contains API rate limiting settings
type RateLimitConfig struct {
	RequestsPerSecond int // Maximum requests per second
	BurstSize         int // Maximum burst size
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	config := &Config{
		DatabaseDSN: getEnv("BILLING_DATABASE_DSN", ""),

		Stripe: StripeConfig{
			APIKey:        getEnv("STRIPE_API_KEY", ""),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
			Enabled:       getEnvBool("STRIPE_ENABLED", true),
		},

		PayPal: PayPalConfig{
			ClientID:     getEnv("PAYPAL_CLIENT_ID", ""),
			ClientSecret: getEnv("PAYPAL_CLIENT_SECRET", ""),
			Environment:  getEnv("PAYPAL_ENVIRONMENT", "sandbox"),
			Enabled:      getEnvBool("PAYPAL_ENABLED", false),
		},

		Invoice: InvoiceConfig{
			DueDays:         getEnvInt("INVOICE_DUE_DAYS", 30),
			NumberPrefix:    getEnv("INVOICE_NUMBER_PREFIX", "INV-"),
			TaxRate:         getEnvDecimal("INVOICE_TAX_RATE", "0.00"),
			DefaultCurrency: getEnv("INVOICE_DEFAULT_CURRENCY", "USD"),
			PDFStoragePath:  getEnv("INVOICE_PDF_STORAGE_PATH", "/tmp/invoices"),
		},

		Usage: UsageConfig{
			AggregationInterval: getEnvDuration("USAGE_AGGREGATION_INTERVAL", "1h"),
			RetentionDays:       getEnvInt("USAGE_RETENTION_DAYS", 90),
			BatchSize:           getEnvInt("USAGE_BATCH_SIZE", 1000),
			EnableRealTime:      getEnvBool("USAGE_ENABLE_REALTIME", true),
		},

		Notifications: NotificationConfig{
			ServiceURL:     getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8080"),
			RetryAttempts:  getEnvInt("NOTIFICATION_RETRY_ATTEMPTS", 3),
			RetryDelay:     getEnvDuration("NOTIFICATION_RETRY_DELAY", "5s"),
			TimeoutSeconds: getEnvInt("NOTIFICATION_TIMEOUT_SECONDS", 30),
		},

		Features: FeatureFlags{
			EnableAutoPayment:   getEnvBool("FEATURE_AUTO_PAYMENT", true),
			EnableUsageMetrics:  getEnvBool("FEATURE_USAGE_METRICS", true),
			EnableTieredPricing: getEnvBool("FEATURE_TIERED_PRICING", true),
			EnableMultiCurrency: getEnvBool("FEATURE_MULTI_CURRENCY", false),
			EnableCredits:       getEnvBool("FEATURE_CREDITS", true),
			EnableProration:     getEnvBool("FEATURE_PRORATION", true),
		},

		RateLimits: RateLimitConfig{
			RequestsPerSecond: getEnvInt("RATE_LIMIT_RPS", 100),
			BurstSize:         getEnvInt("RATE_LIMIT_BURST", 200),
		},
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.DatabaseDSN == "" {
		return fmt.Errorf("database DSN is required")
	}

	if c.Stripe.Enabled && c.Stripe.APIKey == "" {
		return fmt.Errorf("Stripe API key is required when Stripe is enabled")
	}

	if c.PayPal.Enabled && (c.PayPal.ClientID == "" || c.PayPal.ClientSecret == "") {
		return fmt.Errorf("PayPal client ID and secret are required when PayPal is enabled")
	}

	if c.Invoice.DueDays <= 0 {
		return fmt.Errorf("invoice due days must be positive")
	}

	if c.Usage.AggregationInterval <= 0 {
		return fmt.Errorf("usage aggregation interval must be positive")
	}

	if c.Usage.RetentionDays <= 0 {
		return fmt.Errorf("usage retention days must be positive")
	}

	return nil
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getEnvDecimal(key string, defaultValue string) decimal.Decimal {
	if value := os.Getenv(key); value != "" {
		if dec, err := decimal.NewFromString(value); err == nil {
			return dec
		}
	}
	dec, _ := decimal.NewFromString(defaultValue)
	return dec
}
