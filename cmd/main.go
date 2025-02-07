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
	windowTitle := flag.String("window", "", "Window title to find")
	hide := flag.Bool("hide", false, "Hide the window (true = hide, false = show)")
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
			state := ""
			if win.IsMinimized {
				state = " (Minimized)"
			} else if win.IsMaximized {
				state = " (Maximized)"
			}
			visibility := ""
			if !win.IsVisible {
				visibility = " [Hidden]"
			}
			fmt.Printf("- %s%s%s\n", win.Title, state, visibility)
		}
		return
	}

	// Handle window visibility if a window title is provided
	if *windowTitle != "" {
		// Find the window
		hwnd, err := svc.FindWindow(*windowTitle)
		if err != nil {
			log.Fatalf("Failed to find window: %v", err)
		}
		fmt.Printf("Found window: 0x%x\n", hwnd)

		// Set window visibility
		if err := svc.SetWindowVisibility(hwnd, *hide); err != nil {
			log.Fatalf("Failed to change window visibility: %v", err)
		}

		fmt.Printf("Window visibility changed successfully. Hidden: %v\n", *hide)
	} else if *listDesktop < 0 && *hideDesktop < 0 && *showDesktop < 0 {
		// No operation specified
		fmt.Println("Please specify one of: -window, -list, -hidedesktop, or -showdesktop flags")
		flag.Usage()
	}
}
