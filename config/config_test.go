package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfigFromArgs tests loading configuration from command line arguments
func TestLoadConfigFromArgs(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a test config file
	configData := []byte(`
logging:
  level: INFO
virtual_desktops:
  minimum_count: 6
ui:
  trayIcon:
    size: 24
    corner_radius: 4
    padding: 2
    bg_opacity: 230
`)
	err := os.WriteFile(configPath, configData, 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		args     []string
		expected *Config
		wantErr  bool
	}{
		{
			name: "default config when no args",
			args: []string{},
			expected: &Config{
				Logging: LogConfig{Level: slog.LevelDebug},
				UI: UIConfig{TrayIcon: TrayIconConfig{
					Size:         22,
					CornerRadius: 4,
					Padding:      2,
					BgOpacity:    230,
				}},
				VirtualDesktops: VirtualDesktopsConfig{MinimumCount: 9},
				Shortcuts:       ShortcutsConfig{Bindings: []KeyBinding{}},
			},
		},
		{
			name: "load from config file",
			args: []string{"--config", configPath},
			expected: &Config{
				Logging: LogConfig{Level: slog.LevelInfo},
				UI: UIConfig{TrayIcon: TrayIconConfig{
					Size:         24,
					CornerRadius: 4,
					Padding:      2,
					BgOpacity:    230,
				}},
				VirtualDesktops: VirtualDesktopsConfig{MinimumCount: 6},
				Shortcuts:       ShortcutsConfig{Bindings: []KeyBinding{}},
			},
		},
		{
			name: "override with command line args",
			args: []string{
				"--config", configPath,
				"--log-level", "DEBUG",
				"--min-desktops", "8",
			},
			expected: &Config{
				Logging: LogConfig{Level: slog.LevelDebug},
				UI: UIConfig{TrayIcon: TrayIconConfig{
					Size:         24,
					CornerRadius: 4,
					Padding:      2,
					BgOpacity:    230,
				}},
				VirtualDesktops: VirtualDesktopsConfig{MinimumCount: 8},
				Shortcuts:       ShortcutsConfig{Bindings: []KeyBinding{}},
			},
		},
		{
			name:    "error on missing config file path",
			args:    []string{"--config"},
			wantErr: true,
		},
		{
			name:    "error on missing log level",
			args:    []string{"--log-level"},
			wantErr: true,
		},
		{
			name:    "error on missing min desktops",
			args:    []string{"--min-desktops"},
			wantErr: true,
		},
		{
			name:    "error on invalid log level",
			args:    []string{"--log-level", "INVALID"},
			wantErr: true,
		},
		{
			name:    "error on invalid min desktops",
			args:    []string{"--min-desktops", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewArgsConfigLoader(tt.args)
			cfg, err := loader.Load()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

// TestMergeConfigs tests configuration merging functionality
func TestMergeConfigs(t *testing.T) {
	tests := []struct {
		name     string
		base     *Config
		override *Config
		expected *Config
	}{
		{
			name: "merge all fields",
			base: &Config{
				Logging: LogConfig{Level: slog.LevelDebug},
				UI: UIConfig{TrayIcon: TrayIconConfig{
					Size:         22,
					CornerRadius: 4,
					Padding:      2,
					BgOpacity:    230,
				}},
				VirtualDesktops: VirtualDesktopsConfig{MinimumCount: 4},
				Shortcuts: ShortcutsConfig{
					Bindings: []KeyBinding{
						{
							Keys:     []string{"LAlt", "1"},
							Action:   "SwitchDesktop",
							Params:   []string{"1"},
						},
					},
				},
			},
			override: &Config{
				Logging: LogConfig{Level: slog.LevelInfo},
				UI: UIConfig{TrayIcon: TrayIconConfig{
					Size:         24,
					CornerRadius: 6,
					Padding:      3,
					BgOpacity:    200,
				}},
				VirtualDesktops: VirtualDesktopsConfig{MinimumCount: 6},
				Shortcuts: ShortcutsConfig{
					Bindings: []KeyBinding{
						{
							Keys:     []string{"LAlt", "2"},
							Action:   "SwitchDesktop",
							Params:   []string{"2"},
						},
					},
				},
			},
			expected: &Config{
				Logging: LogConfig{Level: slog.LevelInfo},
				UI: UIConfig{TrayIcon: TrayIconConfig{
					Size:         24,
					CornerRadius: 6,
					Padding:      3,
					BgOpacity:    200,
				}},
				VirtualDesktops: VirtualDesktopsConfig{MinimumCount: 6},
				Shortcuts: ShortcutsConfig{
					Bindings: []KeyBinding{
						{
							Keys:     []string{"LAlt", "2"},
							Action:   "SwitchDesktop",
							Params:   []string{"2"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeConfigs(tt.base, tt.override)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// TestFileConfigLoader tests loading configuration from different file types
func TestFileConfigLoader(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a test config file
	configData := []byte(`
logging:
  level: INFO
virtual_desktops:
  minimum_count: 6
ui:
  trayIcon:
    size: 24
    corner_radius: 6
    padding: 3
    bg_opacity: 200
shortcuts:
  bindings:
    - keys: ["LAlt", "1"]
      action: "SwitchDesktop"
      params: ["1"]
      category: "Desktop"
`)
	err := os.WriteFile(configPath, configData, 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filePath string
		expected *Config
		wantErr  bool
	}{
		{
			name:     "load valid config file",
			filePath: configPath,
			expected: &Config{
				Logging: LogConfig{Level: slog.LevelInfo},
				UI: UIConfig{TrayIcon: TrayIconConfig{
					Size:         24,
					CornerRadius: 6,
					Padding:      3,
					BgOpacity:    200,
				}},
				VirtualDesktops: VirtualDesktopsConfig{MinimumCount: 6},
				Shortcuts: ShortcutsConfig{
					Bindings: []KeyBinding{
						{
							Keys:     []string{"LAlt", "1"},
							Action:   "SwitchDesktop",
							Params:   []string{"1"},
						},
					},
				},
			},
		},
		{
			name:     "default config for non-existent file",
			filePath: "nonexistent.yaml",
			expected: DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewFileConfigLoader(tt.filePath)
			cfg, err := loader.Load()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

// TestKeyBindingValidation tests the validation of key bindings
func TestKeyBindingValidation(t *testing.T) {
	tests := []struct {
		name       string
		keyBinding KeyBinding
		wantErr    bool
	}{
		{
			name: "valid single key",
			keyBinding: KeyBinding{
				Keys:     []string{"1"},
				Action:   "SwitchDesktop",
				Params:   []string{"1"},
			},
			wantErr: false,
		},
		{
			name: "valid modifier + key",
			keyBinding: KeyBinding{
				Keys:     []string{"LCtrl", "1"},
				Action:   "SwitchDesktop",
				Params:   []string{"1"},
			},
			wantErr: false,
		},
		{
			name: "invalid key",
			keyBinding: KeyBinding{
				Keys:     []string{"invalid"},
				Action:   "SwitchDesktop",
				Params:   []string{"1"},
			},
			wantErr: true,
		},
		{
			name: "empty binding",
			keyBinding: KeyBinding{
				Keys:     []string{},
				Action:   "SwitchDesktop",
				Params:   []string{"1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.keyBinding.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDefaultConfig tests the default configuration values
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, slog.LevelDebug, cfg.Logging.Level)
	assert.Equal(t, 22, cfg.UI.TrayIcon.Size)
	assert.Equal(t, 4, cfg.UI.TrayIcon.CornerRadius)
	assert.Equal(t, 2, cfg.UI.TrayIcon.Padding)
	assert.Equal(t, uint8(230), cfg.UI.TrayIcon.BgOpacity)
	assert.Equal(t, 9, cfg.VirtualDesktops.MinimumCount)
	assert.NotEmpty(t, cfg.Shortcuts.Bindings)
}
