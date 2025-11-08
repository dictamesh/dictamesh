// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/click2-run/dictamesh/pkg/observability"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// EventHandler processes consumed events
type EventHandler func(ctx context.Context, event *Event) error

// Consumer wraps Kafka consumer with observability
type Consumer struct {
	consumer *kafka.Consumer
	config   *Config
	logger   *observability.Logger
	handler  EventHandler
	running  bool
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *Config, logger *observability.Logger, handler EventHandler) (*Consumer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	kafkaConfig := cfg.GetConsumerConfig()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig["bootstrap.servers"],
		"group.id":          kafkaConfig["group.id"],
		"client.id":         kafkaConfig["client.id"],
		"auto.offset.reset": kafkaConfig["auto.offset.reset"],
		"enable.auto.commit": kafkaConfig["enable.auto.commit"],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	logger.Info("Kafka consumer created",
		"bootstrap_servers", cfg.BootstrapServersString(),
		"group_id", cfg.Consumer.GroupID,
	)

	return &Consumer{
		consumer: consumer,
		config:   cfg,
		logger:   logger,
		handler:  handler,
		running:  false,
	}, nil
}

// Subscribe subscribes to topics
func (c *Consumer) Subscribe(topics []string) error {
	fullTopics := make([]string, len(topics))
	for i, topic := range topics {
		fullTopics[i] = c.config.TopicName(topic)
	}

	if err := c.consumer.SubscribeTopics(fullTopics, nil); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	c.logger.Info("subscribed to topics", "topics", fullTopics)
	return nil
}

// Start starts consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	c.running = true

	c.logger.Info("starting consumer")

	for c.running {
		select {
		case <-ctx.Done():
			c.logger.Info("consumer stopped by context")
			return ctx.Err()
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				c.logger.Error("failed to read message", "error", err)
				continue
			}

			if err := c.processMessage(ctx, msg); err != nil {
				c.logger.Error("failed to process message",
					"topic", *msg.TopicPartition.Topic,
					"partition", msg.TopicPartition.Partition,
					"offset", msg.TopicPartition.Offset,
					"error", err,
				)
			} else {
				// Commit offset on success
				if _, err := c.consumer.CommitMessage(msg); err != nil {
					c.logger.Error("failed to commit offset", "error", err)
				}
			}
		}
	}

	return nil
}

// processMessage processes a single message
func (c *Consumer) processMessage(ctx context.Context, msg *kafka.Message) error {
	start := time.Now()

	// Deserialize event
	var event Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Extract trace context from headers
	for _, header := range msg.Headers {
		if header.Key == "trace_id" {
			// Add trace context to ctx if needed
		}
	}

	// Call handler
	if err := c.handler(ctx, &event); err != nil {
		return fmt.Errorf("handler error: %w", err)
	}

	duration := time.Since(start)
	c.logger.Debug("event processed",
		"topic", *msg.TopicPartition.Topic,
		"partition", msg.TopicPartition.Partition,
		"offset", msg.TopicPartition.Offset,
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}

// Stop stops the consumer
func (c *Consumer) Stop() {
	c.running = false
}

// Close closes the consumer
func (c *Consumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}
	c.logger.Info("Kafka consumer closed")
	return nil
}
