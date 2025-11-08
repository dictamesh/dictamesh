// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Core type definitions for DictaMesh SDK
 */

/**
 * Operation types supported by the SDK
 */
export type OperationType =
  | 'get'         // Get single entity
  | 'list'        // List entities
  | 'create'      // Create entity
  | 'update'      // Update entity
  | 'delete'      // Delete entity
  | 'search'      // Full-text search
  | 'aggregate'   // Aggregation query
  | 'batch'       // Batch operation
  | 'transaction'; // Transactional operation

/**
 * Message types for protocol communication
 */
export type MessageType =
  | 'query'
  | 'mutation'
  | 'subscription'
  | 'subscription_data'
  | 'subscription_complete'
  | 'response'
  | 'error'
  | 'ping'
  | 'pong';

/**
 * Filter operators for query DSL
 */
export type FilterOperator =
  | '$eq'         // Equal
  | '$ne'         // Not equal
  | '$gt'         // Greater than
  | '$gte'        // Greater than or equal
  | '$lt'         // Less than
  | '$lte'        // Less than or equal
  | '$in'         // In array
  | '$nin'        // Not in array
  | '$contains'   // Contains substring
  | '$startsWith' // Starts with
  | '$endsWith'   // Ends with
  | '$regex';     // Regex match

/**
 * Sort direction
 */
export type SortDirection = 'asc' | 'desc';

/**
 * Cache type options
 */
export type CacheType = 'memory' | 'storage' | 'indexeddb' | 'custom';

/**
 * Connection state
 */
export type ConnectionState =
  | 'disconnected'
  | 'connecting'
  | 'connected'
  | 'reconnecting'
  | 'error';

/**
 * Client status
 */
export type ClientStatus = {
  connected: boolean;
  adapter: string | null;
  connectionState: ConnectionState;
  lastError?: Error;
  metrics?: ClientMetrics;
};

/**
 * Client metrics
 */
export type ClientMetrics = {
  requestCount: number;
  errorCount: number;
  cacheHitRate: number;
  avgResponseTime: number;
  activeSubscriptions: number;
};

/**
 * Filter expression for queries
 */
export type FilterExpression = {
  [field: string]:
    | any // Exact match
    | {
        [K in FilterOperator]?: any;
      };
};

/**
 * Sort expression
 */
export type SortExpression = {
  field: string;
  direction: SortDirection;
};

/**
 * Cache options
 */
export type CacheOptions = {
  enabled?: boolean;
  ttl?: number; // Time to live in milliseconds
  key?: string; // Custom cache key
  strategy?: 'cache-first' | 'network-first' | 'cache-and-network';
};

/**
 * Query options
 */
export type QueryOptions = {
  // Filtering
  where?: FilterExpression;

  // Sorting
  orderBy?: SortExpression[];

  // Pagination
  limit?: number;
  offset?: number;
  cursor?: string;

  // Field selection
  select?: string[];
  include?: string[];
  exclude?: string[];

  // Relationships
  expand?: string[];

  // Caching
  cache?: CacheOptions;

  // Misc
  timeout?: number;
  signal?: AbortSignal;
};

/**
 * Operation context
 */
export type OperationContext = {
  requestId: string;
  timestamp: number;
  headers?: Record<string, string>;
  auth?: {
    type: 'bearer' | 'apikey' | 'basic' | 'custom';
    token?: string;
    credentials?: any;
  };
  trace?: {
    traceId: string;
    spanId: string;
    parentSpanId?: string;
  };
  metadata?: Record<string, any>;
};

/**
 * Base operation
 */
export type Operation = {
  type: OperationType;
  entity: string;
  params?: Record<string, any>;
  options?: QueryOptions;
  context?: OperationContext;
};

/**
 * Query operation
 */
export type QueryOperation = Operation & {
  type: 'get' | 'list' | 'search' | 'aggregate';
};

/**
 * Mutation operation
 */
export type MutationOperation = Operation & {
  type: 'create' | 'update' | 'delete';
  data?: any;
};

/**
 * Subscription operation
 */
export type SubscriptionOperation = Operation & {
  type: 'list' | 'get';
  onData?: (data: any) => void;
  onError?: (error: Error) => void;
  onComplete?: () => void;
};

/**
 * Operation result
 */
export type OperationResult<T = any> = {
  data?: T;
  meta?: {
    count?: number;
    total?: number;
    page?: number;
    hasMore?: boolean;
    took?: number; // Execution time in ms
  };
  error?: {
    code: string;
    message: string;
    details?: any;
    stack?: string;
  };
};

/**
 * Query result
 */
export type QueryResult<T = any> = OperationResult<T>;

/**
 * Mutation result
 */
export type MutationResult<T = any> = OperationResult<T>;

/**
 * Batch result
 */
export type BatchResult = {
  results: OperationResult[];
  errors: Error[];
  meta?: {
    total: number;
    successful: number;
    failed: number;
    took: number;
  };
};

/**
 * Validation result
 */
export type ValidationResult = {
  valid: boolean;
  errors?: Array<{
    field: string;
    message: string;
    code: string;
  }>;
};

/**
 * Schema field definition
 */
export type SchemaField = {
  name: string;
  type: string;
  required?: boolean;
  nullable?: boolean;
  description?: string;
  validation?: Record<string, any>;
};

/**
 * Schema definition
 */
export type Schema = {
  entity: string;
  version: string;
  fields: SchemaField[];
  relationships?: Array<{
    name: string;
    type: 'one-to-one' | 'one-to-many' | 'many-to-many';
    entity: string;
  }>;
  metadata?: Record<string, any>;
};

/**
 * Adapter capabilities
 */
export type AdapterCapabilities = {
  query: boolean;
  mutation: boolean;
  subscription: boolean;
  batch: boolean;
  transaction: boolean;
  realtime: boolean;
  caching: boolean;
  offline: boolean;
};

/**
 * Unsubscribe function type
 */
export type Unsubscribe = () => void;

/**
 * Message handler type
 */
export type MessageHandler = (message: any) => void;

/**
 * Observer interface
 */
export interface Observer<T> {
  next?: (value: T) => void;
  error?: (error: Error) => void;
  complete?: () => void;
}

/**
 * Observable interface
 */
export interface Observable<T> {
  subscribe(observer: Observer<T>): Unsubscribe;
}

/**
 * Subscription interface
 */
export interface Subscription<T> extends Observable<T> {
  unsubscribe(): void;
  closed: boolean;
}

/**
 * Clear cache options
 */
export type ClearCacheOptions = {
  entity?: string;
  id?: string;
  pattern?: string;
};
