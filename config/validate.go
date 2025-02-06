package config

import "fmt"

// validateConfig checks if the loaded configuration is valid
func validateConfig(cfg *Config) error {
	// Validate virtual desktops configuration
	if err := cfg.VirtualDesktops.Validate(); err != nil {
		return fmt.Errorf("invalid virtual desktops config: %w", err)
	}

	// Validate shortcuts
	for _, binding := range cfg.Shortcuts.Bindings {
		if err := binding.Validate(); err != nil {
			return fmt.Errorf("invalid binding: %w", err)
		}
	}

	return nil
}
