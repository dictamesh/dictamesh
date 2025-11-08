// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package events provides event-driven messaging infrastructure for the DictaMesh
// framework using Kafka/Redpanda. It includes producer/consumer wrappers with
// built-in observability, standardized event schemas, and topic management.
//
// Example usage:
//
//	// Create producer
//	config := events.DefaultConfig()
//	producer, _ := events.NewProducer(config, logger)
//
//	// Publish event
//	event := events.NewEvent("entity.created", "adapter", "entity:123", data)
//	producer.Publish(ctx, events.TopicEntityChanged, event)
//
//	// Create consumer
//	consumer, _ := events.NewConsumer(config, logger, handlerFunc)
//	consumer.Subscribe([]string{events.TopicEntityChanged})
//	consumer.Start(ctx)
package events

const (
	// Version is the package version
	Version = "0.1.0"
)
