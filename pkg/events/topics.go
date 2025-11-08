// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package events

// Common topic names used in DictaMesh
const (
	// Entity topics
	TopicEntityChanged      = "entity.changed"
	TopicEntityRead         = "entity.read"

	// Relationship topics
	TopicRelationshipChanged = "relationship.changed"

	// Schema topics
	TopicSchemaChanged = "schema.changed"

	// Cache topics
	TopicCacheInvalidation = "cache.invalidation"

	// System topics
	TopicSystemEvents = "system.events"
	TopicHealthEvents = "health.events"

	// Dead letter topics
	TopicDeadLetter = "dead-letter"
)

// TopicConfig represents configuration for a specific topic
type TopicConfiguration struct {
	Name              string
	Partitions        int
	ReplicationFactor int
	RetentionMs       int64
	CleanupPolicy     string
}

// GetStandardTopics returns standard topic configurations
func GetStandardTopics(config *Config) []TopicConfiguration {
	return []TopicConfiguration{
		{
			Name:              config.TopicName(TopicEntityChanged),
			Partitions:        config.Topics.NumPartitions,
			ReplicationFactor: config.Topics.ReplicationFactor,
			RetentionMs:       config.Topics.RetentionMs,
			CleanupPolicy:     "delete",
		},
		{
			Name:              config.TopicName(TopicRelationshipChanged),
			Partitions:        config.Topics.NumPartitions,
			ReplicationFactor: config.Topics.ReplicationFactor,
			RetentionMs:       config.Topics.RetentionMs,
			CleanupPolicy:     "delete",
		},
		{
			Name:              config.TopicName(TopicSchemaChanged),
			Partitions:        config.Topics.NumPartitions / 2,
			ReplicationFactor: config.Topics.ReplicationFactor,
			RetentionMs:       config.Topics.RetentionMs * 2, // Keep longer
			CleanupPolicy:     "compact",                      // Schema history
		},
		{
			Name:              config.TopicName(TopicCacheInvalidation),
			Partitions:        config.Topics.NumPartitions,
			ReplicationFactor: config.Topics.ReplicationFactor,
			RetentionMs:       60000, // 1 minute, short-lived
			CleanupPolicy:     "delete",
		},
		{
			Name:              config.TopicName(TopicSystemEvents),
			Partitions:        config.Topics.NumPartitions / 4,
			ReplicationFactor: config.Topics.ReplicationFactor,
			RetentionMs:       config.Topics.RetentionMs,
			CleanupPolicy:     "delete",
		},
		{
			Name:              config.TopicName(TopicDeadLetter),
			Partitions:        config.Topics.NumPartitions / 2,
			ReplicationFactor: config.Topics.ReplicationFactor,
			RetentionMs:       config.Topics.RetentionMs * 4, // Keep longer for investigation
			CleanupPolicy:     "delete",
		},
	}
}
