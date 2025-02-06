// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"image/color"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultConfig creates a new Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Logging: LogConfig{
			Level: slog.LevelDebug,
		},
		UI: UIConfig{
			TrayIcon: TrayIconConfig{
				Size:          22,
				CornerRadius:  4,
				Padding:       2,
				BgOpacity:     230,
				BgColor:       color.RGBA{0, 120, 215, 255},   // Windows blue
				TextColor:     color.RGBA{255, 255, 255, 255}, // White
				ShadowColor:   color.RGBA{0, 0, 0, 255},       // Black
				ShadowOpacity: 40,
			},
		},
		VirtualDesktops: VirtualDesktopsConfig{
			MinimumCount: 9,
		},
		Shortcuts: ShortcutsConfig{
			Bindings: createDefaultDesktopBindings(),
		},
	}
}

// GenerateDefaultConfigFile generates a YAML file with the default configuration.
// This is useful for creating a template configuration file or for documenting the default values.
func GenerateDefaultConfigFile(path string) error {
	// Create default config
	cfg := DefaultConfig()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal config to YAML with comments
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := []byte(`# WinCuts Default Configuration
# This file was automatically generated from the default configuration.
# You can modify these values to customize the application behavior.
# For more information, see the documentation.

`)
	data = append(header, data...)

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// defaultLoggingConfig provides default logging settings
func defaultLoggingConfig() LogConfig {
	return LogConfig{
		Level: slog.LevelDebug,
	}
}

// defaultUIConfig provides default UI settings including tray icon configuration
func defaultUIConfig() UIConfig {
	return UIConfig{
		TrayIcon: defaultTrayIconConfig(),
	}
}

// defaultTrayIconConfig provides default tray icon settings
func defaultTrayIconConfig() TrayIconConfig {
	return TrayIconConfig{
		Size:          22,
		CornerRadius:  4,
		Padding:       2,
		BgOpacity:     230,
		BgColor:       color.RGBA{0, 120, 215, 255},   // Windows blue
		TextColor:     color.RGBA{255, 255, 255, 255}, // White
		ShadowColor:   color.RGBA{0, 0, 0, 255},       // Black
		ShadowOpacity: 40,
	}
}

// defaultVirtualDesktopsConfig provides default virtual desktop settings
func defaultVirtualDesktopsConfig() VirtualDesktopsConfig {
	return VirtualDesktopsConfig{
		MinimumCount: 9,
	}
}

// defaultShortcutsConfig provides default keyboard shortcuts
func defaultShortcutsConfig() ShortcutsConfig {
	return ShortcutsConfig{
		Bindings: createDefaultDesktopBindings(),
	}
}

// createDefaultDesktopBindings creates the default keyboard bindings for desktop operations
func createDefaultDesktopBindings() []KeyBinding {
	var bindings []KeyBinding

	// Create bindings for desktops 1-9
	for i := 1; i <= 9; i++ {
		desktop := fmt.Sprintf("%d", i)

		// Switch to desktop binding (Alt + Number)
		bindings = append(bindings, KeyBinding{
			Keys:   []string{"LAlt", desktop},
			Action: "SwitchDesktop",
			Params: []string{desktop},
		})

		// Move window to desktop binding (Alt + Shift + Number)
		bindings = append(bindings, KeyBinding{
			Keys:   []string{"LAlt", "LShift", desktop},
			Action: "MoveWindowToDesktop",
			Params: []string{desktop},
		})
	}

	// Add create desktop binding (Alt + N)
	bindings = append(bindings, KeyBinding{
		Keys:   []string{"LAlt", "N"},
		Action: "CreateDesktop",
		Params: []string{},
	})

	return bindings
}
