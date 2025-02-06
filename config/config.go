// Package config handles all application configuration.
// This includes logging, user preferences, and other configurable aspects of the application.
package config

import (
	"log/slog"
	"os"
)

// Config holds all application configuration
type Config struct {
	Logging LogConfig
	// Add other configuration sections here as needed
	// UserPreferences UserPreferencesConfig
	// Shortcuts       ShortcutConfig
	// etc...
}

// LogConfig holds logging related configuration
type LogConfig struct {
	Level slog.Level
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Logging: LogConfig{
			Level: slog.LevelDebug, // Default to DEBUG level
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
