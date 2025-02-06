// main delegates application startup to the app package to maintain a clear separation of concerns.
// This ensures that main remains minimal and focused solely on bootstrapping the application.
package main

import (
	"log/slog"
	"os"
	"wincuts/app"
	"wincuts/config"
)

// main is the entry point of the application.
func main() {
	// Load and setup configuration
	cfg := config.LoadConfig()
	config.SetupLogging(cfg)

	if err := app.Run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
