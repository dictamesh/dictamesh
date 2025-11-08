// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * DictaMesh Core SDK - Main Entry Point
 *
 * A framework-agnostic SDK for interacting with DictaMesh data mesh infrastructure.
 * Supports real-time queries, subscriptions, CRUD operations, and multiple backend adapters.
 */

// Core exports
export { DictaMeshClient } from './core/client';
export type { ClientConfig, ClientStatus } from './core/client';

// Core types
export type {
  OperationType,
  MessageType,
  FilterOperator,
  SortDirection,
  CacheType,
  ConnectionState,
  FilterExpression,
  SortExpression,
  CacheOptions,
  QueryOptions,
  OperationContext,
  Operation,
  QueryOperation,
  MutationOperation,
  SubscriptionOperation,
  OperationResult,
  QueryResult,
  MutationResult,
  BatchResult,
  ValidationResult,
  Schema,
  SchemaField,
  AdapterCapabilities,
  Observable,
  Subscription,
  Observer,
  Unsubscribe,
  MessageHandler,
} from './core/types';

// Error classes
export {
  DictaMeshError,
  NetworkError,
  TimeoutError,
  ValidationError,
  NotFoundError,
  AuthorizationError,
  AuthenticationError,
  ProtocolError,
  AdapterError,
  ConfigurationError,
  CacheError,
  SubscriptionError,
} from './core/error';

// Abstractions
export type { IAdapter, AdapterConfig } from './abstractions/adapter';
export { BaseAdapter } from './abstractions/adapter';

export type { IConnector, ConnectorConfig, ProtocolMessage } from './abstractions/connector';
export { BaseConnector } from './abstractions/connector';

export type { ICache, CacheEntry, CacheStats } from './abstractions/cache';
export { BaseCache } from './abstractions/cache';

export type { IMiddleware, MiddlewareContext, MiddlewareNext, MiddlewareFunction } from './abstractions/middleware';
export { createMiddleware } from './abstractions/middleware';

// Protocol
export {
  createRequestMessage,
  createResponseMessage,
  parseResponseMessage,
  generateMessageId,
} from './protocol/message';

// Adapters
export { DictaMeshAdapter } from './adapters/dictamesh/adapter';
export type { DictaMeshAdapterConfig } from './adapters/dictamesh/adapter';

// Connectors
export { HTTPConnector } from './connectors/http/connector';

// Cache
export { MemoryCache } from './cache/memory-cache';

// Version
export const VERSION = '0.1.0';
