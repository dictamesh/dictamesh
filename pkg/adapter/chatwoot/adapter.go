// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package chatwoot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/click2-run/dictamesh/pkg/adapter"
)

const (
	// AdapterName is the unique name for the Chatwoot adapter
	AdapterName = "chatwoot"

	// AdapterVersion is the current version of the adapter
	AdapterVersion = "1.0.0"
)

// Adapter represents the Chatwoot adapter implementation
type Adapter struct {
	config *Config

	// API clients
	platformClient     *PlatformClient
	applicationClient  *ApplicationClient
	publicClient       *PublicClient

	// State management
	initialized bool
	mu          sync.RWMutex
}

// NewAdapter creates a new Chatwoot adapter
func NewAdapter() *Adapter {
	return &Adapter{
		initialized: false,
	}
}

// Name returns the adapter name
func (a *Adapter) Name() string {
	return AdapterName
}

// Version returns the adapter version
func (a *Adapter) Version() string {
	return AdapterVersion
}

// Initialize initializes the adapter with the given configuration
func (a *Adapter) Initialize(ctx context.Context, config adapter.Config) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.initialized {
		return adapter.ErrAlreadyInitialized
	}

	// Parse configuration
	var cfg *Config
	var err error

	switch c := config.(type) {
	case *Config:
		cfg = c
	case *adapter.MapConfig:
		cfg, err = NewConfig(c.GetAll())
		if err != nil {
			return fmt.Errorf("failed to parse configuration: %w", err)
		}
	default:
		return fmt.Errorf("unsupported configuration type: %T", config)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	a.config = cfg

	// Initialize API clients based on enabled features
	if cfg.EnablePlatformAPI {
		a.platformClient = NewPlatformClient(cfg)
	}

	if cfg.EnableApplicationAPI {
		a.applicationClient = NewApplicationClient(cfg)
	}

	if cfg.EnablePublicAPI {
		a.publicClient = NewPublicClient(cfg)
	}

	a.initialized = true

	return nil
}

// Health checks the health of the adapter and its connection to Chatwoot
func (a *Adapter) Health(ctx context.Context) (*adapter.HealthStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.initialized {
		return &adapter.HealthStatus{
			Status:    adapter.HealthStatusUnhealthy,
			Message:   "adapter not initialized",
			Timestamp: time.Now(),
		}, adapter.ErrNotInitialized
	}

	startTime := time.Now()
	status := &adapter.HealthStatus{
		Status:    adapter.HealthStatusHealthy,
		Message:   "adapter is healthy",
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}

	// Check each enabled API client
	healthyClients := 0
	totalClients := 0

	if a.config.EnablePlatformAPI && a.platformClient != nil {
		totalClients++
		if err := a.platformClient.Ping(ctx); err == nil {
			healthyClients++
			status.Details["platform_api"] = "healthy"
		} else {
			status.Details["platform_api"] = fmt.Sprintf("unhealthy: %v", err)
		}
	}

	if a.config.EnableApplicationAPI && a.applicationClient != nil {
		totalClients++
		if err := a.applicationClient.Ping(ctx); err == nil {
			healthyClients++
			status.Details["application_api"] = "healthy"
		} else {
			status.Details["application_api"] = fmt.Sprintf("unhealthy: %v", err)
		}
	}

	if a.config.EnablePublicAPI && a.publicClient != nil {
		totalClients++
		// Public API doesn't have a ping endpoint, so we assume it's healthy if configured
		healthyClients++
		status.Details["public_api"] = "healthy"
	}

	// Determine overall health status
	if healthyClients == 0 {
		status.Status = adapter.HealthStatusUnhealthy
		status.Message = "all API clients are unhealthy"
	} else if healthyClients < totalClients {
		status.Status = adapter.HealthStatusDegraded
		status.Message = fmt.Sprintf("%d of %d API clients are healthy", healthyClients, totalClients)
	}

	status.Latency = time.Since(startTime)

	return status, nil
}

// Shutdown gracefully shuts down the adapter
func (a *Adapter) Shutdown(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.initialized {
		return adapter.ErrNotInitialized
	}

	// Close all API clients
	if a.platformClient != nil {
		a.platformClient.Close()
	}

	if a.applicationClient != nil {
		a.applicationClient.Close()
	}

	if a.publicClient != nil {
		a.publicClient.Close()
	}

	a.initialized = false

	return nil
}

// GetCapabilities returns the capabilities supported by this adapter
func (a *Adapter) GetCapabilities() []adapter.Capability {
	capabilities := []adapter.Capability{
		adapter.CapabilityRead,
		adapter.CapabilityWrite,
		adapter.CapabilityPagination,
		adapter.CapabilitySearch,
		adapter.CapabilityWebhooks,
	}

	return capabilities
}

// GetPlatformClient returns the Platform API client
func (a *Adapter) GetPlatformClient() (*PlatformClient, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.initialized {
		return nil, adapter.ErrNotInitialized
	}

	if !a.config.EnablePlatformAPI || a.platformClient == nil {
		return nil, adapter.NewAdapterError(
			adapter.ErrorCodeNotSupported,
			"Platform API is not enabled",
			nil,
		)
	}

	return a.platformClient, nil
}

// GetApplicationClient returns the Application API client
func (a *Adapter) GetApplicationClient() (*ApplicationClient, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.initialized {
		return nil, adapter.ErrNotInitialized
	}

	if !a.config.EnableApplicationAPI || a.applicationClient == nil {
		return nil, adapter.NewAdapterError(
			adapter.ErrorCodeNotSupported,
			"Application API is not enabled",
			nil,
		)
	}

	return a.applicationClient, nil
}

// GetPublicClient returns the Public API client
func (a *Adapter) GetPublicClient() (*PublicClient, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.initialized {
		return nil, adapter.ErrNotInitialized
	}

	if !a.config.EnablePublicAPI || a.publicClient == nil {
		return nil, adapter.NewAdapterError(
			adapter.ErrorCodeNotSupported,
			"Public API is not enabled",
			nil,
		)
	}

	return a.publicClient, nil
}

// GetConfig returns the adapter configuration
func (a *Adapter) GetConfig() *Config {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.config
}

// IsInitialized returns whether the adapter is initialized
func (a *Adapter) IsInitialized() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.initialized
}

// Ensure Adapter implements adapter.Adapter interface
var _ adapter.Adapter = (*Adapter)(nil)
