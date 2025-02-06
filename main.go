// main delegates application startup to the app package to maintain a clear separation of concerns.
// This ensures that main remains minimal and focused solely on bootstrapping the application.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"wincuts/app"
	"wincuts/background"
	"wincuts/config"
)

// Version is set during build using -ldflags
var Version = "development"

// main is the entry point of the application.
func main() {
	// Parse command line flags
	showVersion := flag.Bool("v", false, "Show version")
	noWindow := flag.Bool("background", true, "Run in background mode without a window")
	debug := flag.Bool("debug", false, "Run in debug mode")
	flag.Parse()

	// Check for version flag
	if *showVersion {
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

	// Run in background mode by default unless debug is enabled
	if *noWindow && !*debug {
		if err := background.RunInBackground(); err != nil {
			slog.Error("background error", "error", err)
			os.Exit(1)
		}
		return
	}

	// Run as a normal application with console window
	if err := app.Run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
