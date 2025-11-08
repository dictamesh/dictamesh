// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * HTTP Connector for DictaMesh SDK
 */

import { BaseConnector, type ConnectorConfig, type ProtocolMessage } from '../../abstractions/connector';
import { NetworkError, TimeoutError } from '../../core/error';

/**
 * HTTP Connector implementation
 */
export class HTTPConnector extends BaseConnector {
  private abortControllers: Map<string, AbortController> = new Map();

  constructor(config: ConnectorConfig) {
    super(config);
  }

  async send(message: ProtocolMessage): Promise<ProtocolMessage> {
    const controller = new AbortController();
    this.abortControllers.set(message.id, controller);

    try {
      const timeout = this.config.timeout || 30000;
      const timeoutId = setTimeout(() => controller.abort(), timeout);

      const response = await fetch(this.config.endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...this.config.headers,
        },
        body: JSON.stringify(message),
        signal: controller.signal,
      });

      clearTimeout(timeoutId);
      this.abortControllers.delete(message.id);

      if (!response.ok) {
        throw new NetworkError(
          `HTTP ${response.status}: ${response.statusText}`,
          { status: response.status, statusText: response.statusText }
        );
      }

      const responseMessage: ProtocolMessage = await response.json();
      return responseMessage;
    } catch (error) {
      this.abortControllers.delete(message.id);

      if (error instanceof Error) {
        if (error.name === 'AbortError') {
          throw new TimeoutError('Request timeout', { originalError: error });
        }
        throw new NetworkError(error.message, { originalError: error });
      }

      throw error;
    }
  }

  async connect(): Promise<void> {
    this.setState('connected');
  }

  async disconnect(): Promise<void> {
    // Cancel all pending requests
    this.abortControllers.forEach(controller => controller.abort());
    this.abortControllers.clear();
    this.setState('disconnected');
  }
}
