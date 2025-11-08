// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// Package cache provides multi-layer caching with Redis integration
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// CacheLayer represents different cache levels
type CacheLayer string

const (
	L1Memory   CacheLayer = "l1_memory"
	L2Redis    CacheLayer = "l2_redis"
	L3Database CacheLayer = "l3_postgres"
)

// Cache provides multi-layer caching capabilities
type Cache struct {
	logger *zap.Logger
	redis  *redis.Client

	// L1 in-memory cache
	l1     map[string]*cacheEntry
	l1Mu   sync.RWMutex
	l1TTL  time.Duration
	l1Size int

	// L2 Redis cache
	l2TTL time.Duration

	// Metrics
	metrics *CacheMetrics
}

// cacheEntry represents an entry in the L1 cache
type cacheEntry struct {
	Value     []byte
	ExpiresAt time.Time
}

// CacheMetrics tracks cache performance
type CacheMetrics struct {
	L1Hits      int64
	L1Misses    int64
	L2Hits      int64
	L2Misses    int64
	L1Size      int
	L1Evictions int64
	mu          sync.RWMutex
}

// Config represents cache configuration
type Config struct {
	RedisURL  string
	L1TTL     time.Duration
	L2TTL     time.Duration
	L1MaxSize int
}

// DefaultConfig returns default cache configuration
func DefaultConfig() *Config {
	return &Config{
		RedisURL:  "redis://localhost:6379",
		L1TTL:     5 * time.Minute,
		L2TTL:     30 * time.Minute,
		L1MaxSize: 1000,
	}
}

// New creates a new multi-layer cache
func New(config *Config, logger *zap.Logger) (*Cache, error) {
	// Parse Redis URL
	opts, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %w", err)
	}

	// Create Redis client
	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	cache := &Cache{
		logger:  logger,
		redis:   client,
		l1:      make(map[string]*cacheEntry, config.L1MaxSize),
		l1TTL:   config.L1TTL,
		l2TTL:   config.L2TTL,
		l1Size:  config.L1MaxSize,
		metrics: &CacheMetrics{},
	}

	logger.Info("cache initialized",
		zap.String("redis_url", opts.Addr),
		zap.Duration("l1_ttl", config.L1TTL),
		zap.Duration("l2_ttl", config.L2TTL),
		zap.Int("l1_max_size", config.L1MaxSize),
	)

	return cache, nil
}

// Get retrieves a value from cache (tries L1, then L2)
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	// Try L1 first
	if value, found := c.getL1(key); found {
		c.recordHit(L1Memory)
		return value, nil
	}
	c.recordMiss(L1Memory)

	// Try L2 (Redis)
	value, err := c.getL2(ctx, key)
	if err == nil {
		c.recordHit(L2Redis)
		// Backfill L1
		c.setL1(key, value)
		return value, nil
	}

	if err != redis.Nil {
		c.logger.Error("Redis get error", zap.String("key", key), zap.Error(err))
	}
	c.recordMiss(L2Redis)

	return nil, fmt.Errorf("cache miss")
}

// Set stores a value in all cache layers
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	// Set in L1
	c.setL1(key, value)

	// Set in L2 with specified TTL or default
	if ttl == 0 {
		ttl = c.l2TTL
	}

	if err := c.setL2(ctx, key, value, ttl); err != nil {
		c.logger.Error("Redis set error",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set in Redis: %w", err)
	}

	return nil
}

// SetJSON stores a JSON-serializable value
func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.Set(ctx, key, data, ttl)
}

// GetJSON retrieves and unmarshals a JSON value
func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a value from all cache layers
func (c *Cache) Delete(ctx context.Context, key string) error {
	// Delete from L1
	c.deleteL1(key)

	// Delete from L2
	if err := c.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from Redis: %w", err)
	}

	return nil
}

// DeletePattern deletes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	// L1: Delete matching keys
	c.l1Mu.Lock()
	for key := range c.l1 {
		// Simple pattern matching (can be enhanced)
		if matchPattern(key, pattern) {
			delete(c.l1, key)
		}
	}
	c.l1Mu.Unlock()

	// L2: Use Redis SCAN to delete matching keys
	iter := c.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}

// Clear clears all cache layers
func (c *Cache) Clear(ctx context.Context) error {
	// Clear L1
	c.l1Mu.Lock()
	c.l1 = make(map[string]*cacheEntry, c.l1Size)
	c.l1Mu.Unlock()

	// Clear L2 (flush all databases)
	if err := c.redis.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush Redis: %w", err)
	}

	return nil
}

// GetMetrics returns current cache metrics
func (c *Cache) GetMetrics() *CacheMetrics {
	c.metrics.mu.RLock()
	defer c.metrics.mu.RUnlock()

	c.l1Mu.RLock()
	c.metrics.L1Size = len(c.l1)
	c.l1Mu.RUnlock()

	return c.metrics
}

// Close closes the cache connections
func (c *Cache) Close() error {
	if err := c.redis.Close(); err != nil {
		return fmt.Errorf("failed to close Redis: %w", err)
	}
	return nil
}

// L1 cache operations

func (c *Cache) getL1(key string) ([]byte, bool) {
	c.l1Mu.RLock()
	defer c.l1Mu.RUnlock()

	entry, found := c.l1[key]
	if !found {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

func (c *Cache) setL1(key string, value []byte) {
	c.l1Mu.Lock()
	defer c.l1Mu.Unlock()

	// Evict old entries if cache is full
	if len(c.l1) >= c.l1Size {
		c.evictL1()
	}

	c.l1[key] = &cacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.l1TTL),
	}
}

func (c *Cache) deleteL1(key string) {
	c.l1Mu.Lock()
	defer c.l1Mu.Unlock()
	delete(c.l1, key)
}

func (c *Cache) evictL1() {
	// Simple LRU: Remove expired entries first
	now := time.Now()
	evicted := 0

	for key, entry := range c.l1 {
		if now.After(entry.ExpiresAt) {
			delete(c.l1, key)
			evicted++
		}
	}

	c.metrics.mu.Lock()
	c.metrics.L1Evictions += int64(evicted)
	c.metrics.mu.Unlock()

	// If still full, remove oldest 10%
	if len(c.l1) >= c.l1Size {
		toRemove := c.l1Size / 10
		for key := range c.l1 {
			if toRemove <= 0 {
				break
			}
			delete(c.l1, key)
			toRemove--
			evicted++
		}
	}
}

// L2 cache operations

func (c *Cache) getL2(ctx context.Context, key string) ([]byte, error) {
	return c.redis.Get(ctx, key).Bytes()
}

func (c *Cache) setL2(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.redis.Set(ctx, key, value, ttl).Err()
}

// Metrics recording

func (c *Cache) recordHit(layer CacheLayer) {
	c.metrics.mu.Lock()
	defer c.metrics.mu.Unlock()

	switch layer {
	case L1Memory:
		c.metrics.L1Hits++
	case L2Redis:
		c.metrics.L2Hits++
	}
}

func (c *Cache) recordMiss(layer CacheLayer) {
	c.metrics.mu.Lock()
	defer c.metrics.mu.Unlock()

	switch layer {
	case L1Memory:
		c.metrics.L1Misses++
	case L2Redis:
		c.metrics.L2Misses++
	}
}

// Helper functions

func matchPattern(key, pattern string) bool {
	// Simple wildcard matching (* at end)
	if len(pattern) == 0 {
		return false
	}

	if pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}

	return key == pattern
}
