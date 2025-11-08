// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Middleware abstractions for DictaMesh SDK
 *
 * Middleware for request/response transformation and cross-cutting concerns
 */

import type { Operation, OperationContext, OperationResult } from '../core/types';

/**
 * Middleware context with request and response
 */
export interface MiddlewareContext {
  operation: Operation;
  context: OperationContext;
  result?: OperationResult;
  error?: Error;
}

/**
 * Next function to call next middleware in chain
 */
export type MiddlewareNext = () => Promise<OperationResult>;

/**
 * Middleware interface
 */
export interface IMiddleware {
  /**
   * Middleware name for debugging
   */
  readonly name: string;

  /**
   * Execute middleware logic
   */
  execute(
    ctx: MiddlewareContext,
    next: MiddlewareNext
  ): Promise<OperationResult>;
}

/**
 * Middleware function type
 */
export type MiddlewareFunction = (
  ctx: MiddlewareContext,
  next: MiddlewareNext
) => Promise<OperationResult>;

/**
 * Create middleware from function
 */
export function createMiddleware(
  name: string,
  fn: MiddlewareFunction
): IMiddleware {
  return {
    name,
    execute: fn,
  };
}
