// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package events

import (
	"fmt"
	"time"
)

// Config holds all event bus configuration
type Config struct {
	// BootstrapServers is the list of Kafka bootstrap servers
	BootstrapServers []string

	// SchemaRegistryURL is the URL of the Schema Registry
	SchemaRegistryURL string

	// Producer configuration
	Producer ProducerConfig

	// Consumer configuration
	Consumer ConsumerConfig

	// Topic configuration
	Topics TopicConfig

	// Security configuration
	Security SecurityConfig
}

// ProducerConfig configures the Kafka producer
type ProducerConfig struct {
	// ClientID is the identifier for this producer
	ClientID string

	// Acks defines the number of acknowledgments the producer requires
	// 0 = no acknowledgment, 1 = leader acknowledgment, -1/all = all replicas
	Acks string

	// Compression type (none, gzip, snappy, lz4, zstd)
	Compression string

	// MaxMessageBytes is the maximum size of a message
	MaxMessageBytes int

	// RequestTimeout is the timeout for producer requests
	RequestTimeout time.Duration

	// RetryAttempts is the number of retry attempts for failed sends
	RetryAttempts int

	// RetryBackoff is the backoff duration between retries
	RetryBackoff time.Duration

	// EnableIdempotence ensures messages are delivered exactly once
	EnableIdempotence bool

	// LingerMs is the time to wait before sending a batch
	LingerMs int

	// BatchSize is the maximum size of a message batch
	BatchSize int
}

// ConsumerConfig configures the Kafka consumer
type ConsumerConfig struct {
	// GroupID is the consumer group ID
	GroupID string

	// ClientID is the identifier for this consumer
	ClientID string

	// AutoOffsetReset determines where to start consuming when there's no offset
	// "earliest" = start from beginning, "latest" = start from end
	AutoOffsetReset string

	// EnableAutoCommit determines if offsets are committed automatically
	EnableAutoCommit bool

	// AutoCommitInterval is the frequency of auto commits
	AutoCommitInterval time.Duration

	// SessionTimeout is the timeout for consumer session
	SessionTimeout time.Duration

	// HeartbeatInterval is the interval between heartbeats
	HeartbeatInterval time.Duration

	// MaxPollInterval is the maximum time between polls
	MaxPollInterval time.Duration

	// MaxPollRecords is the maximum number of records per poll
	MaxPollRecords int

	// FetchMinBytes is the minimum amount of data to fetch
	FetchMinBytes int

	// FetchMaxWait is the maximum time to wait for fetch
	FetchMaxWait time.Duration

	// IsolationLevel determines the transaction isolation level
	// "read_uncommitted" or "read_committed"
	IsolationLevel string
}

// TopicConfig configures topic creation and management
type TopicConfig struct {
	// Prefix is prepended to all topic names (e.g., "dictamesh.")
	Prefix string

	// NumPartitions is the default number of partitions for new topics
	NumPartitions int

	// ReplicationFactor is the default replication factor
	ReplicationFactor int

	// RetentionMs is the default retention period in milliseconds
	RetentionMs int64

	// CleanupPolicy is the default cleanup policy (delete, compact, or both)
	CleanupPolicy string

	// AutoCreateTopics determines if topics should be created automatically
	AutoCreateTopics bool
}

// SecurityConfig configures Kafka security
type SecurityConfig struct {
	// Protocol is the security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL)
	Protocol string

	// SASLMechanism is the SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512)
	SASLMechanism string

	// SASLUsername is the SASL username
	SASLUsername string

	// SASLPassword is the SASL password
	SASLPassword string

	// SSLCALocation is the path to CA certificate
	SSLCALocation string

	// SSLCertLocation is the path to client certificate
	SSLCertLocation string

	// SSLKeyLocation is the path to client key
	SSLKeyLocation string

	// SSLKeyPassword is the password for client key
	SSLKeyPassword string
}

// DefaultConfig returns a default configuration suitable for development with Redpanda
func DefaultConfig() *Config {
	return &Config{
		BootstrapServers:  []string{"localhost:19092"},
		SchemaRegistryURL: "http://localhost:8081",
		Producer: ProducerConfig{
			ClientID:          "dictamesh-producer",
			Acks:              "all", // Wait for all replicas
			Compression:       "snappy",
			MaxMessageBytes:   1048576, // 1MB
			RequestTimeout:    30 * time.Second,
			RetryAttempts:     3,
			RetryBackoff:      100 * time.Millisecond,
			EnableIdempotence: true,
			LingerMs:          10,
			BatchSize:         16384, // 16KB
		},
		Consumer: ConsumerConfig{
			GroupID:            "dictamesh-consumer-group",
			ClientID:           "dictamesh-consumer",
			AutoOffsetReset:    "earliest",
			EnableAutoCommit:   false, // Manual commit for reliability
			AutoCommitInterval: 5 * time.Second,
			SessionTimeout:     10 * time.Second,
			HeartbeatInterval:  3 * time.Second,
			MaxPollInterval:    5 * time.Minute,
			MaxPollRecords:     500,
			FetchMinBytes:      1,
			FetchMaxWait:       500 * time.Millisecond,
			IsolationLevel:     "read_committed",
		},
		Topics: TopicConfig{
			Prefix:            "dictamesh.",
			NumPartitions:     12,
			ReplicationFactor: 1, // Single broker for dev
			RetentionMs:       7 * 24 * 60 * 60 * 1000, // 7 days
			CleanupPolicy:     "delete",
			AutoCreateTopics:  true,
		},
		Security: SecurityConfig{
			Protocol: "PLAINTEXT",
		},
	}
}

// ProductionConfig returns a configuration optimized for production
func ProductionConfig() *Config {
	cfg := DefaultConfig()

	// Production Kafka settings
	cfg.BootstrapServers = []string{
		"kafka-0.kafka-headless:9092",
		"kafka-1.kafka-headless:9092",
		"kafka-2.kafka-headless:9092",
	}

	// Producer settings for production
	cfg.Producer.Acks = "all"
	cfg.Producer.EnableIdempotence = true
	cfg.Producer.RetryAttempts = 5
	cfg.Producer.LingerMs = 100 // Batch more for efficiency

	// Consumer settings for production
	cfg.Consumer.EnableAutoCommit = false // Always manual commit in production
	cfg.Consumer.SessionTimeout = 30 * time.Second
	cfg.Consumer.MaxPollInterval = 10 * time.Minute

	// Topic settings for production
	cfg.Topics.ReplicationFactor = 3 // High availability
	cfg.Topics.NumPartitions = 12
	cfg.Topics.RetentionMs = 30 * 24 * 60 * 60 * 1000 // 30 days

	// Security for production
	cfg.Security.Protocol = "SASL_SSL"
	cfg.Security.SASLMechanism = "SCRAM-SHA-512"

	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.BootstrapServers) == 0 {
		return fmt.Errorf("bootstrap servers cannot be empty")
	}

	if c.Producer.Acks != "0" && c.Producer.Acks != "1" && c.Producer.Acks != "all" && c.Producer.Acks != "-1" {
		return fmt.Errorf("invalid producer acks value: %s", c.Producer.Acks)
	}

	if c.Consumer.AutoOffsetReset != "earliest" && c.Consumer.AutoOffsetReset != "latest" {
		return fmt.Errorf("invalid auto offset reset value: %s", c.Consumer.AutoOffsetReset)
	}

	if c.Topics.NumPartitions < 1 {
		return fmt.Errorf("number of partitions must be at least 1")
	}

	if c.Topics.ReplicationFactor < 1 {
		return fmt.Errorf("replication factor must be at least 1")
	}

	return nil
}

// GetProducerConfig returns Kafka producer configuration map
func (c *Config) GetProducerConfig() map[string]interface{} {
	config := map[string]interface{}{
		"bootstrap.servers":  c.BootstrapServersString(),
		"client.id":          c.Producer.ClientID,
		"acks":               c.Producer.Acks,
		"compression.type":   c.Producer.Compression,
		"message.max.bytes":  c.Producer.MaxMessageBytes,
		"request.timeout.ms": int(c.Producer.RequestTimeout.Milliseconds()),
		"retries":            c.Producer.RetryAttempts,
		"retry.backoff.ms":   int(c.Producer.RetryBackoff.Milliseconds()),
		"enable.idempotence": c.Producer.EnableIdempotence,
		"linger.ms":          c.Producer.LingerMs,
		"batch.size":         c.Producer.BatchSize,
	}

	// Add security configuration
	c.addSecurityConfig(config)

	return config
}

// GetConsumerConfig returns Kafka consumer configuration map
func (c *Config) GetConsumerConfig() map[string]interface{} {
	config := map[string]interface{}{
		"bootstrap.servers":        c.BootstrapServersString(),
		"group.id":                 c.Consumer.GroupID,
		"client.id":                c.Consumer.ClientID,
		"auto.offset.reset":        c.Consumer.AutoOffsetReset,
		"enable.auto.commit":       c.Consumer.EnableAutoCommit,
		"auto.commit.interval.ms":  int(c.Consumer.AutoCommitInterval.Milliseconds()),
		"session.timeout.ms":       int(c.Consumer.SessionTimeout.Milliseconds()),
		"heartbeat.interval.ms":    int(c.Consumer.HeartbeatInterval.Milliseconds()),
		"max.poll.interval.ms":     int(c.Consumer.MaxPollInterval.Milliseconds()),
		"max.poll.records":         c.Consumer.MaxPollRecords,
		"fetch.min.bytes":          c.Consumer.FetchMinBytes,
		"fetch.max.wait.ms":        int(c.Consumer.FetchMaxWait.Milliseconds()),
		"isolation.level":          c.Consumer.IsolationLevel,
		"go.application.rebalance.enable": true,
	}

	// Add security configuration
	c.addSecurityConfig(config)

	return config
}

// BootstrapServersString returns bootstrap servers as a comma-separated string
func (c *Config) BootstrapServersString() string {
	result := ""
	for i, server := range c.BootstrapServers {
		if i > 0 {
			result += ","
		}
		result += server
	}
	return result
}

// addSecurityConfig adds security configuration to the config map
func (c *Config) addSecurityConfig(config map[string]interface{}) {
	if c.Security.Protocol != "PLAINTEXT" {
		config["security.protocol"] = c.Security.Protocol
	}

	if c.Security.SASLMechanism != "" {
		config["sasl.mechanism"] = c.Security.SASLMechanism
		config["sasl.username"] = c.Security.SASLUsername
		config["sasl.password"] = c.Security.SASLPassword
	}

	if c.Security.SSLCALocation != "" {
		config["ssl.ca.location"] = c.Security.SSLCALocation
	}

	if c.Security.SSLCertLocation != "" {
		config["ssl.certificate.location"] = c.Security.SSLCertLocation
	}

	if c.Security.SSLKeyLocation != "" {
		config["ssl.key.location"] = c.Security.SSLKeyLocation
	}

	if c.Security.SSLKeyPassword != "" {
		config["ssl.key.password"] = c.Security.SSLKeyPassword
	}
}

// TopicName returns the full topic name with prefix
func (c *Config) TopicName(name string) string {
	if c.Topics.Prefix == "" {
		return name
	}
	return c.Topics.Prefix + name
}
