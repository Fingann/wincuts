// Package config provides configuration management for the application.
package config

import "wincuts/keyboard/types"

// ActionProvider defines the contract for providing actions.
// This follows the Interface Segregation Principle by keeping the interface focused on a single responsibility.
type ActionProvider interface {
	// GetActions returns a map of action names to Action objects.
	GetActions() map[string]Action
}

// KeyProvider defines the contract for providing valid keys.
// This follows the Interface Segregation Principle by keeping the interface focused on a single responsibility.
type KeyProvider interface {
	// GetValidKeys returns a map of key names to virtual key codes.
	GetValidKeys() map[string]types.VirtualKey
}

// ConfigLoader defines the contract for loading configuration from various sources.
// This follows the Interface Segregation Principle by keeping the interface focused on a single responsibility.
type ConfigLoader interface {
	// Load loads the configuration from a source and returns it.
	Load() (*Config, error)
}

// ConfigValidator defines the contract for validating configuration.
// This follows the Interface Segregation Principle by keeping the interface focused on a single responsibility.
type ConfigValidator interface {
	// Validate validates the configuration and returns an error if invalid.
	Validate() error
}

// ConfigMerger defines the contract for merging configurations.
// This follows the Interface Segregation Principle by keeping the interface focused on a single responsibility.
type ConfigMerger interface {
	// Merge merges the override configuration into the base configuration.
	Merge(override *Config) error
}
