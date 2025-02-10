//go:build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"wincuts/window"
)

func main() {
	// Parse command line flags
	listDesktop := flag.Int("list", -1, "List all windows on the specified desktop number")
	hideDesktop := flag.Int("hidedesktop", -1, "Hide all windows on the specified desktop number")
	showDesktop := flag.Int("showdesktop", -1, "Show all windows on the specified desktop number")
	flag.Parse()

	// Create a cancellable context

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Println("\nCancelling operation...")
	}()

	// Create window service
	svc, err := window.NewService()
	if err != nil {
		log.Fatalf("Failed to create window service: %v", err)
	}

	// Handle show windows on desktop if specified
	if *showDesktop >= 0 {
		if err := svc.ShowWindowsOnDesktop(*showDesktop); err != nil {
			log.Fatalf("Failed to show windows on desktop %d: %v", *showDesktop, err)
		}
		return
	}

	// Handle hide windows on desktop if specified
	if *hideDesktop >= 0 {
		if err := svc.HideWindowsOnDesktop(*hideDesktop); err != nil {
			log.Fatalf("Failed to hide windows on desktop %d: %v", *hideDesktop, err)
		}
		return
	}

	// List windows on desktop if requested
	if *listDesktop >= 0 {
		windows, err := svc.GetWindowsOnDesktop(*listDesktop)
		if err != nil {
			log.Fatalf("Failed to get windows on desktop %d: %v", *listDesktop, err)
		}
		fmt.Printf("Windows on desktop %d:\n", *listDesktop)
		for _, win := range windows {
			fmt.Printf("- %s: %d, %v\n", win.Title, win.DesktopNum, win.IsHidden)
		}
		return
	}

	 if *listDesktop < 0 && *hideDesktop < 0 && *showDesktop < 0 {
		// No operation specified
		fmt.Println("Please specify one of: -window, -list, -hidedesktop, or -showdesktop flags")
		flag.Usage()
	}
}
