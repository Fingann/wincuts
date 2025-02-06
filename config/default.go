// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"image/color"
	"log/slog"
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
			Keys:     []string{"LAlt", desktop},
			Action:   "SwitchDesktop",
			Params:   []string{desktop},
			Category: "Desktop",
		})

		// Move window to desktop binding (Alt + Shift + Number)
		bindings = append(bindings, KeyBinding{
			Keys:     []string{"LAlt", "LShift", desktop},
			Action:   "MoveWindowToDesktop",
			Params:   []string{desktop},
			Category: "Window",
		})
	}

	// Add create desktop binding (Alt + N)
	bindings = append(bindings, KeyBinding{
		Keys:     []string{"LAlt", "N"},
		Action:   "CreateDesktop",
		Params:   []string{},
		Category: "Desktop",
	})

	return bindings
}
