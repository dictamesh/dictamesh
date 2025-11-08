// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Adapter abstractions for DictaMesh SDK
 *
 * Adapters are responsible for translating SDK operations into
 * backend-specific protocols (GraphQL, REST, custom protocols, etc.)
 */

import type {
  Operation,
  OperationContext,
  OperationResult,
  SubscriptionOperation,
  Schema,
  AdapterCapabilities,
  Observable,
} from '../core/types';

/**
 * Adapter configuration
 */
export interface AdapterConfig {
  /**
   * Adapter name
   */
  name: string;

  /**
   * Adapter version
   */
  version: string;

  /**
   * Base URL or endpoint
   */
  endpoint?: string;

  /**
   * Additional configuration options
   */
  options?: Record<string, any>;
}

/**
 * Adapter interface
 *
 * All adapters must implement this interface to integrate with the SDK
 */
export interface IAdapter {
  /**
   * Adapter metadata
   */
  readonly name: string;
  readonly version: string;
  readonly capabilities: AdapterCapabilities;

  /**
   * Initialize the adapter with configuration
   */
  initialize(config: AdapterConfig): Promise<void>;

  /**
   * Execute an operation
   */
  execute<T = any>(
    operation: Operation,
    context: OperationContext
  ): Promise<OperationResult<T>>;

  /**
   * Subscribe to real-time updates
   */
  subscribe<T = any>(
    operation: SubscriptionOperation,
    context: OperationContext
  ): Observable<T>;

  /**
   * Introspect schema information
   */
  introspect?(): Promise<Schema[]>;

  /**
   * Get schema for a specific entity
   */
  getSchema?(entity: string): Promise<Schema>;

  /**
   * Dispose and cleanup resources
   */
  dispose(): Promise<void>;
}

/**
 * Base adapter class with common functionality
 */
export abstract class BaseAdapter implements IAdapter {
  abstract readonly name: string;
  abstract readonly version: string;
  abstract readonly capabilities: AdapterCapabilities;

  protected config?: AdapterConfig;
  protected initialized = false;

  async initialize(config: AdapterConfig): Promise<void> {
    if (this.initialized) {
      throw new Error('Adapter already initialized');
    }
    this.config = config;
    this.initialized = true;
  }

  abstract execute<T = any>(
    operation: Operation,
    context: OperationContext
  ): Promise<OperationResult<T>>;

  abstract subscribe<T = any>(
    operation: SubscriptionOperation,
    context: OperationContext
  ): Observable<T>;

  async dispose(): Promise<void> {
    this.initialized = false;
    this.config = undefined;
  }

  protected ensureInitialized(): void {
    if (!this.initialized) {
      throw new Error('Adapter not initialized. Call initialize() first.');
    }
  }
}
