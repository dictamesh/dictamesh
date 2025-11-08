// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * DictaMesh Protocol Adapter
 *
 * This adapter implements the DictaMesh-specific protocol for communicating
 * with DictaMesh Gateway and backend services.
 */

import { BaseAdapter, type AdapterConfig } from '../../abstractions/adapter';
import type { IConnector } from '../../abstractions/connector';
import type {
  Operation,
  OperationContext,
  OperationResult,
  SubscriptionOperation,
  Observable,
  AdapterCapabilities,
} from '../../core/types';
import { createRequestMessage, parseResponseMessage } from '../../protocol/message';
import { AdapterError } from '../../core/error';

/**
 * DictaMesh adapter configuration
 */
export interface DictaMeshAdapterConfig extends AdapterConfig {
  connector: IConnector;
}

/**
 * DictaMesh Protocol Adapter
 */
export class DictaMeshAdapter extends BaseAdapter {
  readonly name = 'dictamesh';
  readonly version = '1.0.0';
  readonly capabilities: AdapterCapabilities = {
    query: true,
    mutation: true,
    subscription: true,
    batch: true,
    transaction: false,
    realtime: true,
    caching: true,
    offline: false,
  };

  private connector!: IConnector;

  async initialize(config: AdapterConfig): Promise<void> {
    await super.initialize(config);

    const dictameshConfig = config as DictaMeshAdapterConfig;
    if (!dictameshConfig.connector) {
      throw new AdapterError('DictaMesh adapter requires a connector');
    }

    this.connector = dictameshConfig.connector;
    await this.connector.connect();
  }

  async execute<T = any>(
    operation: Operation,
    context: OperationContext
  ): Promise<OperationResult<T>> {
    this.ensureInitialized();

    try {
      // Create protocol message
      const requestMessage = createRequestMessage(operation, '1.0.0');

      // Add context information
      if (context.auth) {
        requestMessage.auth = context.auth;
      }
      if (context.trace) {
        requestMessage.trace = context.trace;
      }

      // Send request through connector
      const responseMessage = await this.connector.send(requestMessage);

      // Parse response
      const result = parseResponseMessage(responseMessage);

      if (result.error) {
        throw new AdapterError(
          result.error.message,
          result.error
        );
      }

      return result;
    } catch (error) {
      if (error instanceof AdapterError) {
        throw error;
      }
      throw new AdapterError(
        `Failed to execute operation: ${(error as Error).message}`,
        { originalError: error }
      );
    }
  }

  subscribe<T = any>(
    operation: SubscriptionOperation,
    context: OperationContext
  ): Observable<T> {
    this.ensureInitialized();

    // TODO: Implement subscription using WebSocket connector
    throw new AdapterError('Subscriptions not yet implemented in DictaMesh adapter');
  }

  async dispose(): Promise<void> {
    if (this.connector) {
      await this.connector.disconnect();
    }
    await super.dispose();
  }
}
