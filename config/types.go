// Package config provides configuration management for the application.
package config

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log/slog"
	"wincuts/keyboard/types"
)

// Config holds all application configuration.
// This struct follows the Single Responsibility Principle by being a pure data container.
type Config struct {
	Logging         LogConfig             `yaml:"logging" json:"logging"`
	UI              UIConfig              `yaml:"ui" json:"ui"`
	VirtualDesktops VirtualDesktopsConfig `yaml:"virtual_desktops" json:"virtual_desktops"`
	Shortcuts       ShortcutsConfig       `yaml:"shortcuts" json:"shortcuts"`
}

// LogConfig holds logging related configuration.
type LogConfig struct {
	Level slog.Level `yaml:"level" json:"level"`
}

// UnmarshalYAML implements yaml.Unmarshaler for LogConfig.
func (l *LogConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		Level string
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	return l.parseLevel(raw.Level)
}

// UnmarshalJSON implements json.Unmarshaler for LogConfig.
func (l *LogConfig) UnmarshalJSON(data []byte) error {
	var raw struct {
		Level string
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	return l.parseLevel(raw.Level)
}

// parseLevel converts a string level to slog.Level.
func (l *LogConfig) parseLevel(level string) error {
	switch level {
	case "DEBUG":
		l.Level = slog.LevelDebug
	case "INFO":
		l.Level = slog.LevelInfo
	case "WARN":
		l.Level = slog.LevelWarn
	case "ERROR":
		l.Level = slog.LevelError
	default:
		l.Level = slog.LevelDebug // Default to debug if invalid
	}
	return nil
}

// UIConfig holds UI related configuration including colors and styling.
type UIConfig struct {
	TrayIcon TrayIconConfig `yaml:"tray_icon" json:"tray_icon"`
}

// TrayIconConfig holds configuration for the system tray icon.
type TrayIconConfig struct {
	Size          int        `yaml:"size" json:"size"`
	CornerRadius  int        `yaml:"corner_radius" json:"corner_radius"`
	Padding       int        `yaml:"padding" json:"padding"`
	BgOpacity     uint8      `yaml:"bg_opacity" json:"bg_opacity"`
	BgColor       color.RGBA `yaml:"bg_color" json:"bg_color"`
	TextColor     color.RGBA `yaml:"text_color" json:"text_color"`
	ShadowColor   color.RGBA `yaml:"shadow_color" json:"shadow_color"`
	ShadowOpacity uint8      `yaml:"shadow_opacity" json:"shadow_opacity"`
}

// VirtualDesktopsConfig holds configuration for virtual desktops.
type VirtualDesktopsConfig struct {
	MinimumCount int `yaml:"minimum_count" json:"minimum_count"` // Minimum number of virtual desktops to ensure
}

// Validate implements ConfigValidator for VirtualDesktopsConfig.
func (v *VirtualDesktopsConfig) Validate() error {
	if v.MinimumCount < 0 {
		return fmt.Errorf("minimum_count cannot be negative")
	}
	return nil
}

// ShortcutsConfig holds keyboard shortcut configurations.
type ShortcutsConfig struct {
	Bindings []KeyBinding `yaml:"bindings" json:"bindings"`
}

// KeyBinding represents a single keyboard shortcut and its associated action
type KeyBinding struct {
	Keys   []string `yaml:"keys" json:"keys"`     // List of keys that make up the binding (e.g., ["LAlt", "LShift", "1"])
	Action string   `yaml:"action" json:"action"` // Name of the action to perform (e.g., "SwitchDesktop", "MoveWindowToDesktop")
	Params []string `yaml:"params" json:"params"` // Parameters for the action (e.g., ["1"] for desktop number)
}

// GetVirtualKeys converts a slice of key names to VirtualKeys
func (k *KeyBinding) GetVirtualKeys() []types.VirtualKey {
	provider := &DefaultKeyProvider{}
	validKeys := provider.GetValidKeys()

	var keys []types.VirtualKey
	for _, keyName := range k.Keys {
		if vk, ok := validKeys[keyName]; ok {
			keys = append(keys, vk)
		}
	}
	return keys
}

// Validate checks if a key binding is valid
func (k *KeyBinding) Validate() error {
	actionProvider := &DefaultActionProvider{}
	keyProvider := &DefaultKeyProvider{}

	// Check if action exists
	actions := actionProvider.GetActions()
	action, exists := actions[k.Action]
	if !exists {
		return fmt.Errorf("unknown action: %s", k.Action)
	}

	// Check if all keys are valid
	validKeys := keyProvider.GetValidKeys()
	for _, key := range k.Keys {
		if _, ok := validKeys[key]; !ok {
			return fmt.Errorf("invalid key: %s", key)
		}
	}

	// Validate parameters using the action's validator
	if action.Validator != nil {
		if err := action.Validator(k.Params); err != nil {
			return fmt.Errorf("invalid parameters for %s: %w", k.Action, err)
		}
	}

	return nil
}

// Action represents a function that can be bound to keys.
type Action struct {
	Name        string               // Name of the action
	Description string               // Description of what the action does
	Category    string               // Category for grouping
	ParamTypes  []string             // Types of parameters expected
	Validator   func([]string) error // Optional validator for parameters
}
