// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Connector abstractions for DictaMesh SDK
 *
 * Connectors handle the transport layer (HTTP, WebSocket, SSE, etc.)
 */

import type { ConnectionState, MessageHandler, Unsubscribe } from '../core/types';

/**
 * Protocol message structure
 */
export interface ProtocolMessage {
  version: string;
  id: string;
  type: string;
  timestamp: number;
  [key: string]: any;
}

/**
 * Connector configuration
 */
export interface ConnectorConfig {
  /**
   * Endpoint URL
   */
  endpoint: string;

  /**
   * Connection timeout in milliseconds
   */
  timeout?: number;

  /**
   * Custom headers
   */
  headers?: Record<string, string>;

  /**
   * Retry configuration
   */
  retry?: {
    enabled: boolean;
    maxAttempts: number;
    delay: number;
    backoff?: 'linear' | 'exponential';
  };

  /**
   * Additional options
   */
  options?: Record<string, any>;
}

/**
 * Connector interface
 *
 * All connectors must implement this interface to handle transport
 */
export interface IConnector {
  /**
   * Send a message and await response
   */
  send(message: ProtocolMessage): Promise<ProtocolMessage>;

  /**
   * Connect to the server (for persistent connections)
   */
  connect(): Promise<void>;

  /**
   * Disconnect from the server
   */
  disconnect(): Promise<void>;

  /**
   * Register a message handler
   */
  onMessage(handler: MessageHandler): Unsubscribe;

  /**
   * Check if connected
   */
  isConnected(): boolean;

  /**
   * Get current connection state
   */
  getConnectionState(): ConnectionState;

  /**
   * Register connection state change handler
   */
  onStateChange(handler: (state: ConnectionState) => void): Unsubscribe;
}

/**
 * Base connector class with common functionality
 */
export abstract class BaseConnector implements IConnector {
  protected config: ConnectorConfig;
  protected state: ConnectionState = 'disconnected';
  protected messageHandlers: MessageHandler[] = [];
  protected stateHandlers: Array<(state: ConnectionState) => void> = [];

  constructor(config: ConnectorConfig) {
    this.config = config;
  }

  abstract send(message: ProtocolMessage): Promise<ProtocolMessage>;
  abstract connect(): Promise<void>;
  abstract disconnect(): Promise<void>;

  onMessage(handler: MessageHandler): Unsubscribe {
    this.messageHandlers.push(handler);
    return () => {
      const index = this.messageHandlers.indexOf(handler);
      if (index > -1) {
        this.messageHandlers.splice(index, 1);
      }
    };
  }

  isConnected(): boolean {
    return this.state === 'connected';
  }

  getConnectionState(): ConnectionState {
    return this.state;
  }

  onStateChange(handler: (state: ConnectionState) => void): Unsubscribe {
    this.stateHandlers.push(handler);
    return () => {
      const index = this.stateHandlers.indexOf(handler);
      if (index > -1) {
        this.stateHandlers.splice(index, 1);
      }
    };
  }

  protected setState(newState: ConnectionState): void {
    if (this.state !== newState) {
      this.state = newState;
      this.stateHandlers.forEach(handler => handler(newState));
    }
  }

  protected notifyMessage(message: any): void {
    this.messageHandlers.forEach(handler => handler(message));
  }
}
