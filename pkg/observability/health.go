// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthStatus represents the health status
type HealthStatus string

const (
	// HealthStatusHealthy indicates the service is healthy
	HealthStatusHealthy HealthStatus = "healthy"

	// HealthStatusDegraded indicates the service is degraded but operational
	HealthStatusDegraded HealthStatus = "degraded"

	// HealthStatusUnhealthy indicates the service is unhealthy
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// HealthCheck is a function that checks the health of a component
type HealthCheck func(ctx context.Context) error

// HealthChecker manages health checks for the service
type HealthChecker struct {
	config     *HealthConfig
	liveness   map[string]HealthCheck
	readiness  map[string]HealthCheck
	startup    map[string]HealthCheck
	mu         sync.RWMutex
	server     *http.Server
	startupOk  bool
	startupMu  sync.RWMutex
}

// HealthResponse is the JSON response for health checks
type HealthResponse struct {
	Status    HealthStatus           `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// CheckResult represents the result of a single health check
type CheckResult struct {
	Status    HealthStatus `json:"status"`
	Error     string       `json:"error,omitempty"`
	Duration  string       `json:"duration"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(cfg *HealthConfig) *HealthChecker {
	return &HealthChecker{
		config:    cfg,
		liveness:  make(map[string]HealthCheck),
		readiness: make(map[string]HealthCheck),
		startup:   make(map[string]HealthCheck),
		startupOk: false,
	}
}

// RegisterLivenessCheck registers a liveness check
// Liveness checks determine if the application is running
func (h *HealthChecker) RegisterLivenessCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.liveness[name] = check
}

// RegisterReadinessCheck registers a readiness check
// Readiness checks determine if the application can serve traffic
func (h *HealthChecker) RegisterReadinessCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.readiness[name] = check
}

// RegisterStartupCheck registers a startup check
// Startup checks determine if the application has started successfully
func (h *HealthChecker) RegisterStartupCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.startup[name] = check
}

// Start starts the health check HTTP server
func (h *HealthChecker) Start() error {
	if !h.config.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc(h.config.LivenessPath, h.handleLiveness)
	mux.HandleFunc(h.config.ReadinessPath, h.handleReadiness)
	mux.HandleFunc(h.config.StartupPath, h.handleStartup)

	h.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", h.config.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log error
			fmt.Printf("health check server error: %v\n", err)
		}
	}()

	// Run startup checks once
	go h.runStartupChecks()

	return nil
}

// Shutdown gracefully shuts down the health check server
func (h *HealthChecker) Shutdown(ctx context.Context) error {
	if h.server == nil {
		return nil
	}
	return h.server.Shutdown(ctx)
}

// handleLiveness handles liveness probe requests
func (h *HealthChecker) handleLiveness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	result := h.runChecks(ctx, h.liveness)
	h.writeResponse(w, result)
}

// handleReadiness handles readiness probe requests
func (h *HealthChecker) handleReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	// Check startup first
	h.startupMu.RLock()
	startupOk := h.startupOk
	h.startupMu.RUnlock()

	if !startupOk {
		h.writeResponse(w, &HealthResponse{
			Status:    HealthStatusUnhealthy,
			Timestamp: time.Now(),
			Error:     "startup checks not completed",
		})
		return
	}

	result := h.runChecks(ctx, h.readiness)
	h.writeResponse(w, result)
}

// handleStartup handles startup probe requests
func (h *HealthChecker) handleStartup(w http.ResponseWriter, r *http.Request) {
	h.startupMu.RLock()
	startupOk := h.startupOk
	h.startupMu.RUnlock()

	if startupOk {
		h.writeResponse(w, &HealthResponse{
			Status:    HealthStatusHealthy,
			Timestamp: time.Now(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	result := h.runChecks(ctx, h.startup)
	h.writeResponse(w, result)
}

// runStartupChecks runs startup checks once and marks startup as complete
func (h *HealthChecker) runStartupChecks() {
	ctx, cancel := context.WithTimeout(context.Background(), h.config.Timeout)
	defer cancel()

	result := h.runChecks(ctx, h.startup)

	h.startupMu.Lock()
	defer h.startupMu.Unlock()

	if result.Status == HealthStatusHealthy {
		h.startupOk = true
	}
}

// runChecks runs a set of health checks
func (h *HealthChecker) runChecks(ctx context.Context, checks map[string]HealthCheck) *HealthResponse {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(checks) == 0 {
		return &HealthResponse{
			Status:    HealthStatusHealthy,
			Timestamp: time.Now(),
		}
	}

	results := make(map[string]CheckResult)
	overallStatus := HealthStatusHealthy

	for name, check := range checks {
		start := time.Now()
		err := check(ctx)
		duration := time.Since(start)

		status := HealthStatusHealthy
		errorMsg := ""

		if err != nil {
			status = HealthStatusUnhealthy
			errorMsg = err.Error()
			overallStatus = HealthStatusUnhealthy
		}

		results[name] = CheckResult{
			Status:   status,
			Error:    errorMsg,
			Duration: duration.String(),
		}
	}

	return &HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Checks:    results,
	}
}

// writeResponse writes the health check response
func (h *HealthChecker) writeResponse(w http.ResponseWriter, response *HealthResponse) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusOK
	if response.Status == HealthStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	} else if response.Status == HealthStatusDegraded {
		statusCode = http.StatusOK // 200 but degraded
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error
		fmt.Printf("failed to encode health response: %v\n", err)
	}
}

// Common health checks

// PingCheck is a simple ping check that always succeeds
func PingCheck() HealthCheck {
	return func(ctx context.Context) error {
		return nil
	}
}

// TimeoutCheck creates a check that fails if it takes too long
func TimeoutCheck(timeout time.Duration, check HealthCheck) HealthCheck {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- check(ctx)
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return fmt.Errorf("health check timeout after %v", timeout)
		}
	}
}
