// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * In-memory cache implementation
 */

import { BaseCache, type CacheEntry } from '../abstractions/cache';

/**
 * Memory cache implementation
 */
export class MemoryCache extends BaseCache {
  private store: Map<string, CacheEntry> = new Map();
  private cleanupInterval?: NodeJS.Timeout;

  constructor(private defaultTTL: number = 300000) {
    super();
    this.startCleanup();
  }

  async get<T = any>(key: string): Promise<T | null> {
    const entry = this.store.get(key);

    if (!entry) {
      this.recordMiss();
      return null;
    }

    // Check if expired
    if (entry.ttl && Date.now() - entry.timestamp > entry.ttl) {
      this.store.delete(key);
      this.recordMiss();
      return null;
    }

    this.recordHit();
    return entry.value as T;
  }

  async set<T = any>(key: string, value: T, ttl?: number): Promise<void> {
    const entry: CacheEntry<T> = {
      key,
      value,
      timestamp: Date.now(),
      ttl: ttl || this.defaultTTL,
    };

    this.store.set(key, entry);
  }

  async has(key: string): Promise<boolean> {
    const entry = this.store.get(key);

    if (!entry) {
      return false;
    }

    // Check if expired
    if (entry.ttl && Date.now() - entry.timestamp > entry.ttl) {
      this.store.delete(key);
      return false;
    }

    return true;
  }

  async delete(key: string): Promise<boolean> {
    return this.store.delete(key);
  }

  async clear(pattern?: string): Promise<void> {
    if (!pattern) {
      this.store.clear();
      this.resetStats();
      return;
    }

    // Simple pattern matching (supports * wildcard)
    const regex = new RegExp('^' + pattern.replace(/\*/g, '.*') + '$');
    const keysToDelete: string[] = [];

    for (const key of this.store.keys()) {
      if (regex.test(key)) {
        keysToDelete.push(key);
      }
    }

    keysToDelete.forEach(key => this.store.delete(key));
  }

  async keys(pattern?: string): Promise<string[]> {
    const allKeys = Array.from(this.store.keys());

    if (!pattern) {
      return allKeys;
    }

    const regex = new RegExp('^' + pattern.replace(/\*/g, '.*') + '$');
    return allKeys.filter(key => regex.test(key));
  }

  getStats() {
    return {
      ...super.getStats(),
      size: this.store.size,
    };
  }

  /**
   * Clean up expired entries periodically
   */
  private startCleanup(): void {
    this.cleanupInterval = setInterval(() => {
      const now = Date.now();
      const keysToDelete: string[] = [];

      for (const [key, entry] of this.store.entries()) {
        if (entry.ttl && now - entry.timestamp > entry.ttl) {
          keysToDelete.push(key);
        }
      }

      keysToDelete.forEach(key => this.store.delete(key));
    }, 60000); // Clean up every minute

    // Allow cleanup to be stopped
    if (typeof this.cleanupInterval.unref === 'function') {
      this.cleanupInterval.unref();
    }
  }

  /**
   * Stop cleanup and dispose
   */
  dispose(): void {
    if (this.cleanupInterval) {
      clearInterval(this.cleanupInterval);
    }
    this.store.clear();
  }
}
