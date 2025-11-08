// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package adapter

import (
	"fmt"

	"github.com/click2-run/dictamesh/pkg/events"
)

// Config holds adapter configuration
type Config struct {
	// Adapter metadata
	Name          string
	Version       string
	Description   string
	SourceSystem  string
	Domain        string

	// Feature flags
	EnableCache  bool
	EnableEvents bool

	// External configurations
	EventsConfig *events.Config

	// Custom settings (adapter-specific)
	Settings map[string]interface{}
}

// NewConfig creates a new adapter configuration
func NewConfig(name, version, sourceSystem, domain string) *Config {
	return &Config{
		Name:          name,
		Version:       version,
		SourceSystem:  sourceSystem,
		Domain:        domain,
		EnableCache:   true,
		EnableEvents:  true,
		EventsConfig:  events.DefaultConfig(),
		Settings:      make(map[string]interface{}),
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("adapter name cannot be empty")
	}
	if c.Version == "" {
		return fmt.Errorf("adapter version cannot be empty")
	}
	if c.SourceSystem == "" {
		return fmt.Errorf("source system cannot be empty")
	}
	if c.Domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}
	return nil
}

// WithCache enables or disables caching
func (c *Config) WithCache(enabled bool) *Config {
	c.EnableCache = enabled
	return c
}

// WithEvents enables or disables event publishing
func (c *Config) WithEvents(enabled bool) *Config {
	c.EnableEvents = enabled
	return c
}

// WithSetting sets a custom setting
func (c *Config) WithSetting(key string, value interface{}) *Config {
	c.Settings[key] = value
	return c
}

// GetSetting retrieves a custom setting
func (c *Config) GetSetting(key string) (interface{}, bool) {
	val, ok := c.Settings[key]
	return val, ok
}

// GetStringSeating retrieves a string setting
func (c *Config) GetStringSetting(key string) (string, error) {
	val, ok := c.GetSetting(key)
	if !ok {
		return "", fmt.Errorf("setting %s not found", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("setting %s is not a string", key)
	}
	return str, nil
}
