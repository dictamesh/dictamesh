// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package notifications

import (
	"fmt"
	"time"
)

// Config represents the notifications service configuration
type Config struct {
	// Database configuration
	DatabaseDSN string

	// Redis configuration (for rate limiting and caching)
	RedisURL string

	// Kafka configuration
	KafkaBootstrapServers []string
	KafkaConsumerGroup    string

	// Channel configurations
	Channels ChannelConfig

	// Processing configuration
	Processing ProcessingConfig

	// Rate limiting
	RateLimits RateLimitConfig

	// Observability
	Observability ObservabilityConfig
}

// ChannelConfig contains configuration for all channels
type ChannelConfig struct {
	Email       EmailConfig
	SMS         SMSConfig
	Push        PushConfig
	Slack       SlackConfig
	Webhook     WebhookConfig
	InApp       InAppConfig
	BrowserPush BrowserPushConfig
	PagerDuty   PagerDutyConfig
}

// EmailConfig configures email delivery
type EmailConfig struct {
	Enabled  bool
	Provider string // smtp | ses | sendgrid | mailgun

	// SMTP configuration
	SMTP SMTPConfig

	// AWS SES configuration
	SES SESConfig

	// SendGrid configuration
	SendGrid SendGridConfig

	// Common settings
	From            string
	ReplyTo         string
	MaxAttachments  int
	MaxAttachmentMB int

	// Rate limiting
	RateLimit RateLimitDefinition
}

// SMTPConfig configures SMTP email delivery
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

// SESConfig configures AWS SES
type SESConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	ConfigurationSet string
}

// SendGridConfig configures SendGrid
type SendGridConfig struct {
	APIKey string
}

// SMSConfig configures SMS delivery
type SMSConfig struct {
	Enabled  bool
	Provider string // twilio | sns | messagebird

	// Twilio configuration
	Twilio TwilioConfig

	// AWS SNS configuration
	SNS SNSConfig

	// Common settings
	From      string
	MaxLength int

	// Rate limiting
	RateLimit RateLimitDefinition
}

// TwilioConfig configures Twilio SMS
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// SNSConfig configures AWS SNS
type SNSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

// PushConfig configures push notifications
type PushConfig struct {
	Enabled bool

	// Firebase Cloud Messaging (Android)
	FCM FCMConfig

	// Apple Push Notification Service (iOS)
	APNs APNsConfig

	// Web Push
	WebPush WebPushConfig

	// Rate limiting
	RateLimit RateLimitDefinition
}

// FCMConfig configures Firebase Cloud Messaging
type FCMConfig struct {
	Enabled         bool
	CredentialsFile string
	ProjectID       string
	Priority        string // high | normal
}

// APNsConfig configures Apple Push Notification Service
type APNsConfig struct {
	Enabled             bool
	CertificateFile     string
	CertificatePassword string
	KeyID               string
	TeamID              string
	BundleID            string
	Production          bool
}

// WebPushConfig configures Web Push API
type WebPushConfig struct {
	Enabled           bool
	VAPIDPublicKey    string
	VAPIDPrivateKey   string
	VAPIDSubscriber   string // Email address
}

// SlackConfig configures Slack notifications
type SlackConfig struct {
	Enabled bool

	// Webhook URL (simple integration)
	WebhookURL string

	// Bot token (advanced integration)
	BotToken string

	// Default settings
	DefaultChannel string
	Username       string
	IconEmoji      string

	// Rate limiting
	RateLimit RateLimitDefinition
}

// WebhookConfig configures webhook notifications
type WebhookConfig struct {
	Enabled bool

	// Default webhook settings
	Timeout time.Duration

	// Retry configuration
	Retry RetryConfig

	// Authentication
	Auth WebhookAuthConfig

	// Rate limiting
	RateLimit RateLimitDefinition
}

// WebhookAuthConfig configures webhook authentication
type WebhookAuthConfig struct {
	Type  string // none | bearer | apikey | oauth2
	Token string
}

// InAppConfig configures in-app notifications
type InAppConfig struct {
	Enabled bool

	// Transport mechanism
	Transport string // websocket | sse | longpoll

	// Persistence settings
	PersistenceDays int
	MaxUnread       int

	// WebSocket settings
	WebSocketPath     string
	WebSocketPingTime time.Duration
}

// BrowserPushConfig configures browser push notifications
type BrowserPushConfig struct {
	Enabled bool

	// VAPID keys (same as WebPush)
	VAPIDPublicKey  string
	VAPIDPrivateKey string
	VAPIDSubscriber string
}

// PagerDutyConfig configures PagerDuty integration
type PagerDutyConfig struct {
	Enabled bool

	// API configuration
	APIKey         string
	IntegrationKey string

	// Default settings
	DefaultSeverity string // critical | error | warning | info
}

// ProcessingConfig configures notification processing
type ProcessingConfig struct {
	// Worker pools
	WorkerCount int

	// Queue settings
	QueueBufferSize int
	QueueTimeout    time.Duration

	// Batch processing
	BatchEnabled      bool
	BatchMaxSize      int
	BatchMaxWait      time.Duration
	BatchFlushTicker  time.Duration

	// Retry configuration
	Retry RetryConfig

	// Template rendering
	TemplateTimeout time.Duration
	TemplateCaching bool
}

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          bool
}

// RateLimitConfig configures rate limiting
type RateLimitConfig struct {
	Enabled bool

	// User rate limits (per channel)
	UserLimits map[Channel]RateLimitDefinition

	// System-wide rate limits
	SystemLimits map[Channel]RateLimitDefinition

	// Category rate limits
	CategoryLimits map[string]RateLimitDefinition
}

// RateLimitDefinition defines a rate limit
type RateLimitDefinition struct {
	Count    int
	Duration time.Duration
}

// ObservabilityConfig configures observability
type ObservabilityConfig struct {
	// Metrics
	MetricsEnabled bool
	MetricsPort    int

	// Tracing
	TracingEnabled  bool
	TracingEndpoint string
	TracingSampler  float64

	// Logging
	LogLevel  string // debug | info | warn | error
	LogFormat string // json | text
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.DatabaseDSN == "" {
		return fmt.Errorf("database DSN is required")
	}

	if len(c.KafkaBootstrapServers) == 0 {
		return fmt.Errorf("kafka bootstrap servers are required")
	}

	if c.KafkaConsumerGroup == "" {
		return fmt.Errorf("kafka consumer group is required")
	}

	// Validate at least one channel is enabled
	hasEnabledChannel := c.Channels.Email.Enabled ||
		c.Channels.SMS.Enabled ||
		c.Channels.Push.Enabled ||
		c.Channels.Slack.Enabled ||
		c.Channels.Webhook.Enabled ||
		c.Channels.InApp.Enabled ||
		c.Channels.BrowserPush.Enabled ||
		c.Channels.PagerDuty.Enabled

	if !hasEnabledChannel {
		return fmt.Errorf("at least one notification channel must be enabled")
	}

	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		KafkaConsumerGroup: "dictamesh-notifications",
		Processing: ProcessingConfig{
			WorkerCount:       10,
			QueueBufferSize:   1000,
			QueueTimeout:      30 * time.Second,
			BatchEnabled:      true,
			BatchMaxSize:      100,
			BatchMaxWait:      5 * time.Minute,
			BatchFlushTicker:  1 * time.Minute,
			TemplateTimeout:   5 * time.Second,
			TemplateCaching:   true,
			Retry: RetryConfig{
				MaxAttempts:     3,
				InitialInterval: 1 * time.Second,
				MaxInterval:     30 * time.Second,
				Multiplier:      2.0,
				Jitter:          true,
			},
		},
		RateLimits: RateLimitConfig{
			Enabled: true,
			UserLimits: map[Channel]RateLimitDefinition{
				ChannelEmail: {Count: 100, Duration: 1 * time.Hour},
				ChannelSMS:   {Count: 10, Duration: 1 * time.Hour},
				ChannelPush:  {Count: 50, Duration: 1 * time.Hour},
			},
			SystemLimits: map[Channel]RateLimitDefinition{
				ChannelEmail: {Count: 10000, Duration: 1 * time.Hour},
				ChannelSMS:   {Count: 1000, Duration: 1 * time.Hour},
				ChannelPush:  {Count: 50000, Duration: 1 * time.Hour},
			},
		},
		Observability: ObservabilityConfig{
			MetricsEnabled:  true,
			MetricsPort:     9090,
			TracingEnabled:  true,
			TracingSampler:  0.1,
			LogLevel:        "info",
			LogFormat:       "json",
		},
	}
}
