package config

import (
	"image/color"
	"log/slog"
)

// mergeConfigs merges two configurations, with the override taking precedence
// while preserving default values for unspecified fields
func mergeConfigs(base, override *Config) *Config {
	if override == nil {
		return base
	}

	result := *base

	// Merge logging config
	if override.Logging.Level != slog.LevelDebug {
		result.Logging.Level = override.Logging.Level
	}

	// Merge UI config
	if override.UI.TrayIcon.Size != 0 {
		result.UI.TrayIcon.Size = override.UI.TrayIcon.Size
	}
	if override.UI.TrayIcon.CornerRadius != 0 {
		result.UI.TrayIcon.CornerRadius = override.UI.TrayIcon.CornerRadius
	}
	if override.UI.TrayIcon.Padding != 0 {
		result.UI.TrayIcon.Padding = override.UI.TrayIcon.Padding
	}
	if override.UI.TrayIcon.BgOpacity != 0 {
		result.UI.TrayIcon.BgOpacity = override.UI.TrayIcon.BgOpacity
	}
	if override.UI.TrayIcon.BgColor != (color.RGBA{}) {
		result.UI.TrayIcon.BgColor = override.UI.TrayIcon.BgColor
	}
	if override.UI.TrayIcon.TextColor != (color.RGBA{}) {
		result.UI.TrayIcon.TextColor = override.UI.TrayIcon.TextColor
	}
	if override.UI.TrayIcon.ShadowColor != (color.RGBA{}) {
		result.UI.TrayIcon.ShadowColor = override.UI.TrayIcon.ShadowColor
	}
	if override.UI.TrayIcon.ShadowOpacity != 0 {
		result.UI.TrayIcon.ShadowOpacity = override.UI.TrayIcon.ShadowOpacity
	}

	// Merge virtual desktops config
	if override.VirtualDesktops.MinimumCount != 0 {
		result.VirtualDesktops.MinimumCount = override.VirtualDesktops.MinimumCount
	}

	// Merge shortcuts
	if len(override.Shortcuts.Bindings) > 0 {
		result.Shortcuts.Bindings = override.Shortcuts.Bindings
	}

	return &result
}
