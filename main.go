// main delegates application startup to the app package to maintain a clear separation of concerns.
// This ensures that main remains minimal and focused solely on bootstrapping the application.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"wincuts/app"
	"wincuts/config"
)

// Version is set during build using -ldflags
var Version = "development"

// main is the entry point of the application.
func main() {
	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("WinCuts %s\n", Version)
		os.Exit(0)
	}

	// Load and setup configuration
	cfg, err := config.LoadConfigFromArgs(os.Args)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}
	config.SetupLogging(cfg)

	slog.Info("starting WinCuts", "version", Version)

	if err := app.Run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
