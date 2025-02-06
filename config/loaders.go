// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// FileConfigLoader loads configuration from a file.
// This follows the Single Responsibility Principle by focusing only on file-based configuration loading.
type FileConfigLoader struct {
	filePath string
}

// NewFileConfigLoader creates a new FileConfigLoader.
func NewFileConfigLoader(filePath string) *FileConfigLoader {
	return &FileConfigLoader{filePath: filePath}
}

// Load implements ConfigLoader.
func (f *FileConfigLoader) Load() (*Config, error) {
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// ArgsConfigLoader loads configuration from command line arguments.
// This follows the Single Responsibility Principle by focusing only on argument-based configuration loading.
type ArgsConfigLoader struct {
	args []string
}

// NewArgsConfigLoader creates a new ArgsConfigLoader.
func NewArgsConfigLoader(args []string) *ArgsConfigLoader {
	return &ArgsConfigLoader{args: args}
}

// Load implements ConfigLoader.
func (a *ArgsConfigLoader) Load() (*Config, error) {
	cfg := DefaultConfig()

	for i := 0; i < len(a.args); i++ {
		switch a.args[i] {
		case "--config":
			if i+1 >= len(a.args) {
				return nil, fmt.Errorf("--config requires a file path")
			}
			filePath := a.args[i+1]
			i++ // Skip next arg

			fileLoader := NewFileConfigLoader(filePath)
			fileCfg, err := fileLoader.Load()
			if err != nil {
				return nil, err
			}
			cfg = mergeConfigs(cfg, fileCfg)

		case "--log-level":
			if i+1 >= len(a.args) {
				return nil, fmt.Errorf("--log-level requires a level")
			}
			level := a.args[i+1]
			i++ // Skip next arg

			// Parse log level
			switch level {
			case "DEBUG":
				cfg.Logging.Level = slog.LevelDebug
			case "INFO":
				cfg.Logging.Level = slog.LevelInfo
			case "WARN":
				cfg.Logging.Level = slog.LevelWarn
			case "ERROR":
				cfg.Logging.Level = slog.LevelError
			default:
				return nil, fmt.Errorf("invalid log level: %s", level)
			}

		case "--min-desktops":
			if i+1 >= len(a.args) {
				return nil, fmt.Errorf("--min-desktops requires a number")
			}
			var count int
			if _, err := fmt.Sscanf(a.args[i+1], "%d", &count); err != nil {
				return nil, fmt.Errorf("invalid min desktops count: %s", a.args[i+1])
			}
			i++ // Skip next arg
			cfg.VirtualDesktops.MinimumCount = count
		}
	}

	return cfg, nil
}

// DefaultConfigLoader loads the default configuration.
// This follows the Single Responsibility Principle by focusing only on providing default configuration.
type DefaultConfigLoader struct{}

// NewDefaultConfigLoader creates a new DefaultConfigLoader.
func NewDefaultConfigLoader() *DefaultConfigLoader {
	return &DefaultConfigLoader{}
}

// Load implements ConfigLoader.
func (d *DefaultConfigLoader) Load() (*Config, error) {
	return DefaultConfig(), nil
}
