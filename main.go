// main delegates application startup to the app package to maintain a clear separation of concerns.
// This ensures that main remains minimal and focused solely on bootstrapping the application.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"wincuts/app"
	"wincuts/config"
	"wincuts/service"
)

// Version is set during build using -ldflags
var Version = "development"

// main is the entry point of the application.
func main() {
	// Parse command line flags
	isService := flag.Bool("service", false, "Run as Windows service")
	installService := flag.Bool("install", false, "Install Windows service")
	uninstallService := flag.Bool("uninstall", false, "Uninstall Windows service")
	debug := flag.Bool("debug", false, "Run in debug mode")
	showVersion := flag.Bool("v", false, "Show version")
	flag.Parse()

	// Check for version flag
	if *showVersion {
		fmt.Printf("WinCuts %s\n", Version)
		os.Exit(0)
	}

	// Handle service installation/uninstallation
	if *installService {
		if err := service.InstallService(os.Args[0]); err != nil {
			slog.Error("failed to install service", "error", err)
			os.Exit(1)
		}
		fmt.Println("Service installed successfully")
		os.Exit(0)
	}

	if *uninstallService {
		if err := service.UninstallService(); err != nil {
			slog.Error("failed to uninstall service", "error", err)
			os.Exit(1)
		}
		fmt.Println("Service uninstalled successfully")
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

	// Check if we should run as a service
	if *isService {
		if err := service.RunAsService(*debug); err != nil {
			slog.Error("service error", "error", err)
			os.Exit(1)
		}
		return
	}

	// Run as a normal application
	if err := app.Run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
