// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Cache abstractions for DictaMesh SDK
 *
 * Caching layer for optimizing data access and offline support
 */

/**
 * Cache entry metadata
 */
export interface CacheEntry<T = any> {
  key: string;
  value: T;
  timestamp: number;
  ttl?: number;
  metadata?: Record<string, any>;
}

/**
 * Cache statistics
 */
export interface CacheStats {
  hits: number;
  misses: number;
  size: number;
  hitRate: number;
}

/**
 * Cache interface
 *
 * All cache implementations must implement this interface
 */
export interface ICache {
  /**
   * Get a value from cache
   */
  get<T = any>(key: string): Promise<T | null>;

  /**
   * Set a value in cache
   */
  set<T = any>(key: string, value: T, ttl?: number): Promise<void>;

  /**
   * Check if key exists in cache
   */
  has(key: string): Promise<boolean>;

  /**
   * Delete a value from cache
   */
  delete(key: string): Promise<boolean>;

  /**
   * Clear cache (optionally by pattern)
   */
  clear(pattern?: string): Promise<void>;

  /**
   * Get cache statistics
   */
  getStats(): CacheStats;

  /**
   * Get all keys in cache
   */
  keys(pattern?: string): Promise<string[]>;

  /**
   * Get multiple values at once
   */
  getMany<T = any>(keys: string[]): Promise<Array<T | null>>;

  /**
   * Set multiple values at once
   */
  setMany<T = any>(entries: Array<{ key: string; value: T; ttl?: number }>): Promise<void>;

  /**
   * Delete multiple keys at once
   */
  deleteMany(keys: string[]): Promise<number>;
}

/**
 * Base cache class with common functionality
 */
export abstract class BaseCache implements ICache {
  protected hits = 0;
  protected misses = 0;

  abstract get<T = any>(key: string): Promise<T | null>;
  abstract set<T = any>(key: string, value: T, ttl?: number): Promise<void>;
  abstract has(key: string): Promise<boolean>;
  abstract delete(key: string): Promise<boolean>;
  abstract clear(pattern?: string): Promise<void>;
  abstract keys(pattern?: string): Promise<string[]>;

  async getMany<T = any>(keys: string[]): Promise<Array<T | null>> {
    return Promise.all(keys.map(key => this.get<T>(key)));
  }

  async setMany<T = any>(entries: Array<{ key: string; value: T; ttl?: number }>): Promise<void> {
    await Promise.all(entries.map(entry => this.set(entry.key, entry.value, entry.ttl)));
  }

  async deleteMany(keys: string[]): Promise<number> {
    const results = await Promise.all(keys.map(key => this.delete(key)));
    return results.filter(Boolean).length;
  }

  getStats(): CacheStats {
    const total = this.hits + this.misses;
    return {
      hits: this.hits,
      misses: this.misses,
      size: 0, // Override in implementations
      hitRate: total > 0 ? this.hits / total : 0,
    };
  }

  protected recordHit(): void {
    this.hits++;
  }

  protected recordMiss(): void {
    this.misses++;
  }

  protected resetStats(): void {
    this.hits = 0;
    this.misses = 0;
  }
}
