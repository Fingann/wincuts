// Package config handles all application configuration.
// This includes logging, user preferences, and other configurable aspects of the application.
package config

import (
	"image/color"
	"log/slog"
	"os"
)

// Config holds all application configuration
type Config struct {
	Logging         LogConfig
	UI              UIConfig
	VirtualDesktops VirtualDesktopsConfig
}

// LogConfig holds logging related configuration
type LogConfig struct {
	Level slog.Level
}

// UIConfig holds UI related configuration including colors and styling
type UIConfig struct {
	TrayIcon TrayIconConfig
}

// TrayIconConfig holds configuration for the system tray icon
type TrayIconConfig struct {
	Size          int
	CornerRadius  int
	Padding       int
	BgOpacity     uint8
	BgColor       color.RGBA
	TextColor     color.RGBA
	ShadowColor   color.RGBA
	ShadowOpacity uint8
}

// VirtualDesktopsConfig holds configuration for virtual desktops
type VirtualDesktopsConfig struct {
	MinimumCount int // Minimum number of virtual desktops to ensure
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Logging: LogConfig{
			Level: slog.LevelDebug, // Default to DEBUG level
		},
		UI: UIConfig{
			TrayIcon: TrayIconConfig{
				Size:          22,                             // Icon size in pixels
				CornerRadius:  4,                              // Rounded corner radius
				Padding:       2,                              // Padding around content
				BgOpacity:     230,                            // Background opacity (0-255)
				BgColor:       color.RGBA{0, 120, 215, 255},   // Windows blue
				TextColor:     color.RGBA{255, 255, 255, 255}, // White
				ShadowColor:   color.RGBA{0, 0, 0, 255},       // Black
				ShadowOpacity: 40,                             // Shadow opacity (0-255)
			},
		},
		VirtualDesktops: VirtualDesktopsConfig{
			MinimumCount: 9, // Default to 9 virtual desktops
		},
	}
}

// LoadConfig loads the configuration from environment/files
// In the future, this can be expanded to load from config files
func LoadConfig() *Config {
	cfg := DefaultConfig()

	// Override with environment variables if present
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		switch levelStr {
		case "DEBUG":
			cfg.Logging.Level = slog.LevelDebug
		case "INFO":
			cfg.Logging.Level = slog.LevelInfo
		case "WARN":
			cfg.Logging.Level = slog.LevelWarn
		case "ERROR":
			cfg.Logging.Level = slog.LevelError
		}
	}

	return cfg
}

// SetupLogging configures the global logger based on config
func SetupLogging(cfg *Config) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Logging.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize time format
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format("15:04:05")),
				}
			}
			return a
		},
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
