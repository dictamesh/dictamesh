// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Protocol message definitions for DictaMesh SDK
 */

import type {
  MessageType,
  OperationType,
  Operation,
  OperationResult,
} from '../core/types';

/**
 * Protocol message structure
 */
export interface ProtocolMessage {
  version: string;
  id: string;
  type: MessageType;
  timestamp: number;

  // Request fields
  operation?: {
    type: OperationType;
    entity: string;
    params?: Record<string, any>;
    options?: Record<string, any>;
  };

  // Response fields
  data?: any;
  meta?: {
    count?: number;
    total?: number;
    page?: number;
    hasMore?: boolean;
    took?: number;
  };

  // Error fields
  error?: {
    code: string;
    message: string;
    details?: any;
    stack?: string;
  };

  // Tracing
  trace?: {
    traceId: string;
    spanId: string;
    parentSpanId?: string;
  };
}

/**
 * Create a protocol message from an operation
 */
export function createRequestMessage(
  operation: Operation,
  version: string = '1.0.0'
): ProtocolMessage {
  const id = generateMessageId();
  const timestamp = Date.now();

  let messageType: MessageType;
  switch (operation.type) {
    case 'get':
    case 'list':
    case 'search':
    case 'aggregate':
      messageType = 'query';
      break;
    case 'create':
    case 'update':
    case 'delete':
      messageType = 'mutation';
      break;
    default:
      messageType = 'query';
  }

  return {
    version,
    id,
    type: messageType,
    timestamp,
    operation: {
      type: operation.type,
      entity: operation.entity,
      params: operation.params,
      options: operation.options,
    },
    trace: operation.context?.trace,
  };
}

/**
 * Create a response message
 */
export function createResponseMessage(
  requestId: string,
  result: OperationResult,
  version: string = '1.0.0'
): ProtocolMessage {
  const message: ProtocolMessage = {
    version,
    id: requestId,
    type: result.error ? 'error' : 'response',
    timestamp: Date.now(),
  };

  if (result.error) {
    message.error = result.error;
  } else {
    message.data = result.data;
    message.meta = result.meta;
  }

  return message;
}

/**
 * Parse response message into operation result
 */
export function parseResponseMessage(message: ProtocolMessage): OperationResult {
  if (message.error) {
    return {
      error: message.error,
    };
  }

  return {
    data: message.data,
    meta: message.meta,
  };
}

/**
 * Generate unique message ID
 */
export function generateMessageId(): string {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 11)}`;
}
