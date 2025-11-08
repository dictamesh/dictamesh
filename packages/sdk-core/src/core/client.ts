// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Main DictaMesh SDK Client
 */

import type {
  QueryOperation,
  MutationOperation,
  SubscriptionOperation,
  QueryResult,
  MutationResult,
  Subscription,
  Operation,
  BatchResult,
  ClientStatus,
  ConnectionState,
  OperationContext,
  ClearCacheOptions,
} from './types';
import type { IAdapter } from '../abstractions/adapter';
import type { ICache } from '../abstractions/cache';
import type { IMiddleware } from '../abstractions/middleware';
import { DictaMeshError, ConfigurationError } from './error';
import { generateMessageId } from '../protocol/message';

/**
 * Client configuration
 */
export interface ClientConfig {
  /**
   * Default endpoint URL
   */
  endpoint?: string;

  /**
   * Default adapter name
   */
  defaultAdapter?: string;

  /**
   * Cache configuration
   */
  cache?: {
    enabled?: boolean;
    type?: 'memory' | 'storage' | 'indexeddb';
    ttl?: number;
  };

  /**
   * Authentication configuration
   */
  auth?: {
    type: 'bearer' | 'apikey' | 'basic';
    token?: string;
    credentials?: any;
  };

  /**
   * Request timeout in milliseconds
   */
  timeout?: number;

  /**
   * Additional headers to include in all requests
   */
  headers?: Record<string, string>;

  /**
   * Enable tracing
   */
  tracing?: boolean;

  /**
   * Additional options
   */
  options?: Record<string, any>;
}

/**
 * DictaMesh SDK Client
 */
export class DictaMeshClient {
  private config: ClientConfig;
  private adapters: Map<string, IAdapter> = new Map();
  private defaultAdapter: string | null = null;
  private cache: ICache | null = null;
  private middleware: IMiddleware[] = [];
  private connectionState: ConnectionState = 'disconnected';
  private requestCount = 0;
  private errorCount = 0;

  constructor(config: ClientConfig = {}) {
    this.config = config;
  }

  /**
   * Configure the client
   */
  configure(config: Partial<ClientConfig>): void {
    this.config = { ...this.config, ...config };
  }

  /**
   * Register an adapter
   */
  registerAdapter(name: string, adapter: IAdapter): void {
    if (this.adapters.has(name)) {
      throw new ConfigurationError(
        `Adapter with name "${name}" is already registered`
      );
    }
    this.adapters.set(name, adapter);

    // Set as default if it's the first adapter or explicitly configured
    if (!this.defaultAdapter || this.config.defaultAdapter === name) {
      this.defaultAdapter = name;
    }
  }

  /**
   * Set the default adapter
   */
  setDefaultAdapter(name: string): void {
    if (!this.adapters.has(name)) {
      throw new ConfigurationError(`Adapter "${name}" is not registered`);
    }
    this.defaultAdapter = name;
  }

  /**
   * Get an adapter by name
   */
  getAdapter(name?: string): IAdapter {
    const adapterName = name || this.defaultAdapter;
    if (!adapterName) {
      throw new ConfigurationError('No default adapter configured');
    }

    const adapter = this.adapters.get(adapterName);
    if (!adapter) {
      throw new ConfigurationError(`Adapter "${adapterName}" not found`);
    }

    return adapter;
  }

  /**
   * Register middleware
   */
  use(middleware: IMiddleware): void {
    this.middleware.push(middleware);
  }

  /**
   * Set cache implementation
   */
  setCache(cache: ICache): void {
    this.cache = cache;
  }

  /**
   * Get cache instance
   */
  getCache(): ICache | null {
    return this.cache;
  }

  /**
   * Execute a query operation
   */
  async query<T = any>(operation: QueryOperation): Promise<QueryResult<T>> {
    return this.executeOperation<T>(operation);
  }

  /**
   * Execute a mutation operation
   */
  async mutate<T = any>(operation: MutationOperation): Promise<MutationResult<T>> {
    return this.executeOperation<T>(operation);
  }

  /**
   * Subscribe to real-time updates
   */
  subscribe<T = any>(operation: SubscriptionOperation): Subscription<T> {
    const adapter = this.getAdapter();
    const context = this.createContext(operation);

    if (!adapter.capabilities.subscription) {
      throw new DictaMeshError(
        `Adapter "${adapter.name}" does not support subscriptions`,
        'UNSUPPORTED_OPERATION'
      );
    }

    return adapter.subscribe<T>(operation, context);
  }

  /**
   * Execute batch operations
   */
  async batch(operations: Operation[]): Promise<BatchResult> {
    const startTime = Date.now();
    const results: any[] = [];
    const errors: Error[] = [];

    await Promise.all(
      operations.map(async operation => {
        try {
          const result = await this.executeOperation(operation);
          results.push(result);
        } catch (error) {
          errors.push(error as Error);
          results.push({ error });
        }
      })
    );

    return {
      results,
      errors,
      meta: {
        total: operations.length,
        successful: results.filter(r => !r.error).length,
        failed: errors.length,
        took: Date.now() - startTime,
      },
    };
  }

  /**
   * Clear cache
   */
  async clearCache(options?: ClearCacheOptions): Promise<void> {
    if (!this.cache) {
      return;
    }

    if (options?.entity) {
      const pattern = options.id
        ? `${options.entity}:${options.id}`
        : `${options.entity}:*`;
      await this.cache.clear(pattern);
    } else if (options?.pattern) {
      await this.cache.clear(options.pattern);
    } else {
      await this.cache.clear();
    }
  }

  /**
   * Connect to the server
   */
  async connect(): Promise<void> {
    this.connectionState = 'connecting';
    // Initialize adapters if needed
    this.connectionState = 'connected';
  }

  /**
   * Disconnect from the server
   */
  async disconnect(): Promise<void> {
    this.connectionState = 'disconnected';
    // Cleanup adapters
  }

  /**
   * Check if connected
   */
  isConnected(): boolean {
    return this.connectionState === 'connected';
  }

  /**
   * Get client status
   */
  getStatus(): ClientStatus {
    const total = this.requestCount;
    const cacheHitRate = this.cache ? this.cache.getStats().hitRate : 0;

    return {
      connected: this.isConnected(),
      adapter: this.defaultAdapter,
      connectionState: this.connectionState,
      metrics: {
        requestCount: this.requestCount,
        errorCount: this.errorCount,
        cacheHitRate,
        avgResponseTime: 0, // TODO: Implement
        activeSubscriptions: 0, // TODO: Implement
      },
    };
  }

  /**
   * Execute an operation with middleware pipeline
   */
  private async executeOperation<T = any>(
    operation: Operation
  ): Promise<QueryResult<T>> {
    this.requestCount++;
    const adapter = this.getAdapter();
    const context = this.createContext(operation);

    try {
      // Execute through middleware pipeline
      let index = 0;
      const executeNext = async (): Promise<any> => {
        if (index < this.middleware.length) {
          const middleware = this.middleware[index++];
          return middleware.execute(
            { operation, context },
            executeNext
          );
        } else {
          // Final execution through adapter
          return await adapter.execute<T>(operation, context);
        }
      };

      return await executeNext();
    } catch (error) {
      this.errorCount++;
      throw error;
    }
  }

  /**
   * Create operation context
   */
  private createContext(operation: Operation): OperationContext {
    const context: OperationContext = {
      requestId: generateMessageId(),
      timestamp: Date.now(),
      headers: { ...this.config.headers },
    };

    if (this.config.auth) {
      context.auth = this.config.auth;
    }

    if (this.config.tracing) {
      context.trace = {
        traceId: generateMessageId(),
        spanId: generateMessageId(),
      };
    }

    return context;
  }
}
