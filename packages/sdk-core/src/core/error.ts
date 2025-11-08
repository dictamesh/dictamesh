// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Error classes for DictaMesh SDK
 */

/**
 * Base SDK error class
 */
export class DictaMeshError extends Error {
  code: string;
  details?: any;

  constructor(message: string, code: string, details?: any) {
    super(message);
    this.name = 'DictaMeshError';
    this.code = code;
    this.details = details;
    Object.setPrototypeOf(this, DictaMeshError.prototype);
  }
}

/**
 * Network error
 */
export class NetworkError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'NETWORK_ERROR', details);
    this.name = 'NetworkError';
    Object.setPrototypeOf(this, NetworkError.prototype);
  }
}

/**
 * Timeout error
 */
export class TimeoutError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'TIMEOUT_ERROR', details);
    this.name = 'TimeoutError';
    Object.setPrototypeOf(this, TimeoutError.prototype);
  }
}

/**
 * Validation error
 */
export class ValidationError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'VALIDATION_ERROR', details);
    this.name = 'ValidationError';
    Object.setPrototypeOf(this, ValidationError.prototype);
  }
}

/**
 * Not found error
 */
export class NotFoundError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'NOT_FOUND', details);
    this.name = 'NotFoundError';
    Object.setPrototypeOf(this, NotFoundError.prototype);
  }
}

/**
 * Authorization error
 */
export class AuthorizationError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'AUTHORIZATION_ERROR', details);
    this.name = 'AuthorizationError';
    Object.setPrototypeOf(this, AuthorizationError.prototype);
  }
}

/**
 * Authentication error
 */
export class AuthenticationError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'AUTHENTICATION_ERROR', details);
    this.name = 'AuthenticationError';
    Object.setPrototypeOf(this, AuthenticationError.prototype);
  }
}

/**
 * Protocol error
 */
export class ProtocolError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'PROTOCOL_ERROR', details);
    this.name = 'ProtocolError';
    Object.setPrototypeOf(this, ProtocolError.prototype);
  }
}

/**
 * Adapter error
 */
export class AdapterError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'ADAPTER_ERROR', details);
    this.name = 'AdapterError';
    Object.setPrototypeOf(this, AdapterError.prototype);
  }
}

/**
 * Configuration error
 */
export class ConfigurationError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'CONFIGURATION_ERROR', details);
    this.name = 'ConfigurationError';
    Object.setPrototypeOf(this, ConfigurationError.prototype);
  }
}

/**
 * Cache error
 */
export class CacheError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'CACHE_ERROR', details);
    this.name = 'CacheError';
    Object.setPrototypeOf(this, CacheError.prototype);
  }
}

/**
 * Subscription error
 */
export class SubscriptionError extends DictaMeshError {
  constructor(message: string, details?: any) {
    super(message, 'SUBSCRIPTION_ERROR', details);
    this.name = 'SubscriptionError';
    Object.setPrototypeOf(this, SubscriptionError.prototype);
  }
}
