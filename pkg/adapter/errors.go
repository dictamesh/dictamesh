// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package adapter

import (
	"errors"
	"fmt"
	"net/http"
)

// Common adapter errors
var (
	// ErrNotInitialized indicates the adapter has not been initialized
	ErrNotInitialized = errors.New("adapter not initialized")

	// ErrAlreadyInitialized indicates the adapter is already initialized
	ErrAlreadyInitialized = errors.New("adapter already initialized")

	// ErrInvalidConfig indicates the configuration is invalid
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrNotSupported indicates the operation is not supported
	ErrNotSupported = errors.New("operation not supported")

	// ErrResourceNotFound indicates the resource was not found
	ErrResourceNotFound = errors.New("resource not found")

	// ErrUnauthorized indicates the request is unauthorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates the request is forbidden
	ErrForbidden = errors.New("forbidden")

	// ErrRateLimited indicates the request was rate limited
	ErrRateLimited = errors.New("rate limited")

	// ErrBadRequest indicates the request is invalid
	ErrBadRequest = errors.New("bad request")

	// ErrServerError indicates a server error occurred
	ErrServerError = errors.New("server error")

	// ErrTimeout indicates the request timed out
	ErrTimeout = errors.New("request timeout")

	// ErrConnectionFailed indicates the connection to the external system failed
	ErrConnectionFailed = errors.New("connection failed")
)

// AdapterError represents a detailed error from an adapter
type AdapterError struct {
	// Code is the error code
	Code ErrorCode

	// Message is the error message
	Message string

	// StatusCode is the HTTP status code (if applicable)
	StatusCode int

	// Err is the underlying error
	Err error

	// Details contains additional error details
	Details map[string]interface{}

	// Retryable indicates if the operation can be retried
	Retryable bool
}

// Error implements the error interface
func (e *AdapterError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AdapterError) Unwrap() error {
	return e.Err
}

// ErrorCode represents an error code
type ErrorCode string

const (
	// ErrorCodeUnknown represents an unknown error
	ErrorCodeUnknown ErrorCode = "unknown"

	// ErrorCodeInvalidConfig represents an invalid configuration error
	ErrorCodeInvalidConfig ErrorCode = "invalid_config"

	// ErrorCodeNotInitialized represents a not initialized error
	ErrorCodeNotInitialized ErrorCode = "not_initialized"

	// ErrorCodeNotFound represents a not found error
	ErrorCodeNotFound ErrorCode = "not_found"

	// ErrorCodeUnauthorized represents an unauthorized error
	ErrorCodeUnauthorized ErrorCode = "unauthorized"

	// ErrorCodeForbidden represents a forbidden error
	ErrorCodeForbidden ErrorCode = "forbidden"

	// ErrorCodeBadRequest represents a bad request error
	ErrorCodeBadRequest ErrorCode = "bad_request"

	// ErrorCodeRateLimited represents a rate limit error
	ErrorCodeRateLimited ErrorCode = "rate_limited"

	// ErrorCodeServerError represents a server error
	ErrorCodeServerError ErrorCode = "server_error"

	// ErrorCodeTimeout represents a timeout error
	ErrorCodeTimeout ErrorCode = "timeout"

	// ErrorCodeConnectionFailed represents a connection failed error
	ErrorCodeConnectionFailed ErrorCode = "connection_failed"

	// ErrorCodeValidation represents a validation error
	ErrorCodeValidation ErrorCode = "validation"

	// ErrorCodeNotSupported represents an unsupported operation error
	ErrorCodeNotSupported ErrorCode = "not_supported"
)

// NewAdapterError creates a new AdapterError
func NewAdapterError(code ErrorCode, message string, err error) *AdapterError {
	return &AdapterError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

// NewAdapterErrorWithStatus creates a new AdapterError with an HTTP status code
func NewAdapterErrorWithStatus(code ErrorCode, message string, statusCode int, err error) *AdapterError {
	return &AdapterError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
		Details:    make(map[string]interface{}),
		Retryable:  isRetryableStatusCode(statusCode),
	}
}

// WithDetail adds a detail to the error
func (e *AdapterError) WithDetail(key string, value interface{}) *AdapterError {
	e.Details[key] = value
	return e
}

// WithRetryable sets the retryable flag
func (e *AdapterError) WithRetryable(retryable bool) *AdapterError {
	e.Retryable = retryable
	return e
}

// isRetryableStatusCode checks if an HTTP status code indicates a retryable error
func isRetryableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusBadGateway:
		return true
	default:
		return statusCode >= 500
	}
}

// HTTPErrorToAdapterError converts an HTTP error to an AdapterError
func HTTPErrorToAdapterError(statusCode int, message string, err error) *AdapterError {
	var code ErrorCode
	var retryable bool

	switch statusCode {
	case http.StatusNotFound:
		code = ErrorCodeNotFound
		retryable = false
	case http.StatusUnauthorized:
		code = ErrorCodeUnauthorized
		retryable = false
	case http.StatusForbidden:
		code = ErrorCodeForbidden
		retryable = false
	case http.StatusBadRequest:
		code = ErrorCodeBadRequest
		retryable = false
	case http.StatusTooManyRequests:
		code = ErrorCodeRateLimited
		retryable = true
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		code = ErrorCodeTimeout
		retryable = true
	default:
		if statusCode >= 500 {
			code = ErrorCodeServerError
			retryable = true
		} else {
			code = ErrorCodeUnknown
			retryable = false
		}
	}

	return &AdapterError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
		Details:    make(map[string]interface{}),
		Retryable:  retryable,
	}
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	var adapterErr *AdapterError
	if errors.As(err, &adapterErr) {
		return adapterErr.Retryable
	}
	return false
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	if errors.Is(err, ErrResourceNotFound) {
		return true
	}
	var adapterErr *AdapterError
	if errors.As(err, &adapterErr) {
		return adapterErr.Code == ErrorCodeNotFound
	}
	return false
}

// IsAuthError checks if an error is an authentication/authorization error
func IsAuthError(err error) bool {
	if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrForbidden) {
		return true
	}
	var adapterErr *AdapterError
	if errors.As(err, &adapterErr) {
		return adapterErr.Code == ErrorCodeUnauthorized || adapterErr.Code == ErrorCodeForbidden
	}
	return false
}

// IsRateLimitError checks if an error is a rate limit error
func IsRateLimitError(err error) bool {
	if errors.Is(err, ErrRateLimited) {
		return true
	}
	var adapterErr *AdapterError
	if errors.As(err, &adapterErr) {
		return adapterErr.Code == ErrorCodeRateLimited
	}
	return false
}
