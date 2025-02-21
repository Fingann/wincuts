//go:build windows

package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"

	"wincuts/config"
	"wincuts/keyboard"
	"wincuts/keyboard/shortcut"
	"wincuts/keyboard/types"
	"wincuts/systray"
	"wincuts/virtd"

	winapi "github.com/chrsm/winapi"
	"github.com/chrsm/winapi/user"
)

// DesktopManager defines an abstraction for desktop operations.
// This interface allows us to decouple the core application logic from its concrete implementation,
// which improves testability and flexibility to change the underlying desktop handling behavior.

type DesktopManager interface {
	// GetCurrentDesktopCount returns the number of desktops available.
	GetCurrentDesktopCount() int
	// CreateNewDesktop creates a new desktop.
	CreateNewDesktop()
	// SwitchToDesktop switches to the specified desktop number.
	SwitchToDesktop(desktopNumber int)
	// MoveWindowToDesktop moves the given window to the specified desktop.
	MoveWindowToDesktop(window winapi.HWND, desktopNumber int)
}

// VirtdDesktopManager is the concrete implementation of DesktopManager that interacts with the Windows desktop system.
// By centralizing platform-specific calls to the virtd package, we ensure that the rest of the application remains portable and testable.

type VirtdDesktopManager struct{}

func (v VirtdDesktopManager) GetCurrentDesktopCount() int {
	return virtd.GetDesktopCount()
}

func (v VirtdDesktopManager) CreateNewDesktop() {
	virtd.CreateDesktop()
}

func (v VirtdDesktopManager) SwitchToDesktop(desktopNumber int) {
	virtd.GoToDesktopNumber(desktopNumber)
}

func (v VirtdDesktopManager) MoveWindowToDesktop(window winapi.HWND, desktopNumber int) {
	virtd.MoveWindowToDesktopNumber(window, desktopNumber)
}

// EnsureMinimumDesktops enforces a minimum available desktop count at runtime.
// This is crucial for features that depend on several desktops being present, and ensures consistent behavior across environments.
func EnsureMinimumDesktops(dm DesktopManager, minCount int) {
	current := dm.GetCurrentDesktopCount()
	for i := current; i < minCount; i++ {
		dm.CreateNewDesktop()
	}
}

// setupKeyBindings registers keyboard shortcuts for switching desktops and moving windows.
func setupKeyBindings(dm DesktopManager, traySvc *systray.Service) *shortcut.Service {
	keyChan := make(chan *shortcut.KeyBindingAction, 100)
	svc := shortcut.NewService(keyChan, shortcut.NewMatcher())

	// Load configuration
	cfg, err := config.LoadConfigFromArgs(os.Args)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		// Use default config as fallback
		cfg = config.DefaultConfig()
	}

	// Register each configured binding
	for _, binding := range cfg.Shortcuts.Bindings {
		// Validate the binding
		if err := binding.Validate(); err != nil {
			slog.Error("invalid key binding",
				"keys", binding.Keys,
				"action", binding.Action,
				"error", err)
			continue
		}

		// Create the appropriate action based on the binding type
		var action func() error
		var shouldBlock bool

		switch binding.Action {
		case "SwitchDesktop":
			if len(binding.Params) != 1 {
				slog.Error("invalid parameters for SwitchDesktop", "params", binding.Params)
				continue
			}
			desktop := parseDesktopNumber(binding.Params[0]) - 1 // Convert to 0-based index
			action = func() error {
				dm.SwitchToDesktop(desktop)
				if err := traySvc.UpdateDesktop(desktop + 1); err != nil {
					slog.Error("failed to update system tray", "error", err)
				}
				return nil
			}
			shouldBlock = true

		case "MoveWindowToDesktop":
			if len(binding.Params) != 1 {
				slog.Error("invalid parameters for MoveWindowToDesktop", "params", binding.Params)
				continue
			}
			desktop := parseDesktopNumber(binding.Params[0]) - 1 // Convert to 0-based index
			action = func() error {
				foregroundW := user.GetForegroundWindow()
				dm.MoveWindowToDesktop(foregroundW, desktop)
				dm.SwitchToDesktop(desktop)
				if err := traySvc.UpdateDesktop(desktop + 1); err != nil {
					slog.Error("failed to update system tray", "error", err)
				}
				return nil
			}
			shouldBlock = true

		case "CreateDesktop":
			action = func() error {
				dm.CreateNewDesktop()
				return nil
			}
			shouldBlock = true

		default:
			slog.Error("unknown action type", "action", binding.Action)
			continue
		}

		// Register the binding
		svc.RegisterKeyBindingActions(
			shortcut.NewBindingAction(binding.GetVirtualKeys(), action, shouldBlock),
		)

		slog.Debug("registered shortcut",
			"keys", types.NewKeybinding(binding.GetVirtualKeys()...).PrettyString(),
			"action", binding.Action)
	}

	return svc
}

// parseDesktopNumber safely converts a string parameter to a desktop number
func parseDesktopNumber(param string) int {
	num, err := strconv.Atoi(param)
	if err != nil {
		slog.Error("failed to parse desktop number", "param", param, "error", err)
		return 1
	}
	return num
}

// Run aggregates the initialization of system components (desktop environment, keyboard hook, key bindings)
// and starts the user event loop. This separation of startup functionality enhances testability and maintainability.
func Run() error {
	// Load configuration
	cfg, err := config.LoadConfigFromArgs(os.Args)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	config.SetupLogging(cfg)

	// Initialize system tray
	traySvc, err := systray.NewService(cfg.UI.TrayIcon)
	if err != nil {
		return fmt.Errorf("failed to initialize system tray: %w", err)
	}
	defer traySvc.Stop()

	// Enforce minimum number of virtual desktops from config
	dm := VirtdDesktopManager{}
	EnsureMinimumDesktops(dm, cfg.VirtualDesktops.MinimumCount)
	slog.Info("virtual desktops initialized", "count", dm.GetCurrentDesktopCount(), "minimum", cfg.VirtualDesktops.MinimumCount)


	keybindService := setupKeyBindings(dm, traySvc)
	// Initialize the keyboard hook; early exit if setup fails to ensure proper system state.
	hook, err := keyboard.NewHook(keybindService)
	if err != nil {
		return fmt.Errorf("failed to create keyboard hook: %w", err)
	}
	hook.Start()
	slog.Info("keyboard hook initialized")

	// Register keyboard shortcuts to facilitate rapid desktop management.
	keybindService.Start()
	slog.Info("keyboard shortcuts registered")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	slog.Info("started")
	<-signalChan
	slog.Info("stopping")
	return nil
}
