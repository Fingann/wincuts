package background

import (
	"fmt"
	"log/slog"
	"wincuts/app"

	"github.com/lxn/win"
)

// HideConsoleWindow hides the console window when running in background mode
func HideConsoleWindow() {
	console := win.GetConsoleWindow()
	if console != 0 {
		if !win.ShowWindow(console, win.SW_HIDE) {
			slog.Error("failed to hide console window")
		}
	}
}

// RunInBackground starts the application in background mode
func RunInBackground() error {
	// Hide the console window
	HideConsoleWindow()

	// Run the application
	if err := app.Run(); err != nil {
		return fmt.Errorf("application error: %w", err)
	}

	return nil
}
