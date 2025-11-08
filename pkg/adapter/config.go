// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package adapter

import (
	"fmt"
	"time"
)

// MapConfig is a simple map-based configuration implementation
type MapConfig struct {
	data map[string]interface{}
}

// NewMapConfig creates a new MapConfig
func NewMapConfig(data map[string]interface{}) *MapConfig {
	if data == nil {
		data = make(map[string]interface{})
	}
	return &MapConfig{data: data}
}

// GetString retrieves a string value
func (c *MapConfig) GetString(key string) (string, error) {
	val, ok := c.data[key]
	if !ok {
		return "", fmt.Errorf("configuration key not found: %s", key)
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("configuration value for key %s is not a string", key)
	}

	return str, nil
}

// GetInt retrieves an integer value
func (c *MapConfig) GetInt(key string) (int, error) {
	val, ok := c.data[key]
	if !ok {
		return 0, fmt.Errorf("configuration key not found: %s", key)
	}

	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("configuration value for key %s is not an integer", key)
	}
}

// GetBool retrieves a boolean value
func (c *MapConfig) GetBool(key string) (bool, error) {
	val, ok := c.data[key]
	if !ok {
		return false, fmt.Errorf("configuration key not found: %s", key)
	}

	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("configuration value for key %s is not a boolean", key)
	}

	return b, nil
}

// GetDuration retrieves a duration value
func (c *MapConfig) GetDuration(key string) (time.Duration, error) {
	val, ok := c.data[key]
	if !ok {
		return 0, fmt.Errorf("configuration key not found: %s", key)
	}

	switch v := val.(type) {
	case time.Duration:
		return v, nil
	case string:
		return time.ParseDuration(v)
	case int64:
		return time.Duration(v), nil
	default:
		return 0, fmt.Errorf("configuration value for key %s is not a duration", key)
	}
}

// Set sets a configuration value
func (c *MapConfig) Set(key string, value interface{}) {
	c.data[key] = value
}

// Validate performs basic validation
func (c *MapConfig) Validate() error {
	// Base implementation - adapters should override this
	return nil
}

// GetAll returns all configuration data
func (c *MapConfig) GetAll() map[string]interface{} {
	return c.data
}

// Has checks if a key exists
func (c *MapConfig) Has(key string) bool {
	_, ok := c.data[key]
	return ok
}

// GetStringDefault retrieves a string value with a default
func (c *MapConfig) GetStringDefault(key, defaultValue string) string {
	val, err := c.GetString(key)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetIntDefault retrieves an integer value with a default
func (c *MapConfig) GetIntDefault(key string, defaultValue int) int {
	val, err := c.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetBoolDefault retrieves a boolean value with a default
func (c *MapConfig) GetBoolDefault(key string, defaultValue bool) bool {
	val, err := c.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetDurationDefault retrieves a duration value with a default
func (c *MapConfig) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	val, err := c.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return val
}
