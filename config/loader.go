package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfigFromArgs loads configuration based on command line arguments
func LoadConfigFromArgs(args []string) (*Config, error) {
	// Check for generate-config first
	for i := 1; i < len(args); i++ {
		if args[i] == "--generate-config" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--generate-config requires a file path")
			}
			path := args[i+1]

			if err := GenerateDefaultConfigFile(path); err != nil {
				return nil, fmt.Errorf("failed to generate config file: %w", err)
			}
			slog.Info("generated default configuration file", "path", path)
			os.Exit(0) // Exit after generating config
		}
	}

	// Start with default configuration
	config := DefaultConfig()

	// Parse command line arguments
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--config":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--config requires a file path")
			}
			configPath := args[i+1]
			i++

			// Load configuration from file
			fileConfig, err := loadConfigFromFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load config file: %w", err)
			}
			config = mergeConfigs(config, fileConfig)

		case "--log-level":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--log-level requires a level")
			}
			level := args[i+1]
			i++

			// Parse log level
			switch strings.ToUpper(level) {
			case "DEBUG":
				config.Logging.Level = slog.LevelDebug
			case "INFO":
				config.Logging.Level = slog.LevelInfo
			case "WARN":
				config.Logging.Level = slog.LevelWarn
			case "ERROR":
				config.Logging.Level = slog.LevelError
			default:
				return nil, fmt.Errorf("invalid log level: %s", level)
			}

		case "--min-desktops":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--min-desktops requires a number")
			}
			var count int
			if _, err := fmt.Sscanf(args[i+1], "%d", &count); err != nil {
				return nil, fmt.Errorf("invalid min desktops count: %s", args[i+1])
			}
			i++
			config.VirtualDesktops.MinimumCount = count
		}
	}

	return config, nil
}

// loadConfigFromFile loads configuration from a file
func loadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config file: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	return &config, nil
}

// SetupLogging configures the global logger based on config
func SetupLogging(cfg *Config) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Logging.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
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
