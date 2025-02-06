package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

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
// The use of closures to capture the desktop index prevents common closure pitfalls and links user actions to the intended operations.
func setupKeyBindings(hook *keyboard.Hook, dm DesktopManager, traySvc *systray.Service) *shortcut.Service {
	subscription := hook.Subscribe()
	svc := shortcut.NewService(subscription)

	// Map number keys 1-9 to virtual desktops
	desktopKeys := []types.VirtualKey{
		types.VK_1, types.VK_2, types.VK_3, types.VK_4, types.VK_5,
		types.VK_6, types.VK_7, types.VK_8, types.VK_9,
	}

	// Get current desktop count to know how many shortcuts to register
	currentCount := dm.GetCurrentDesktopCount()
	if currentCount > len(desktopKeys) {
		currentCount = len(desktopKeys) // Cap at maximum supported keys
	}

	for index := 0; index < currentCount; index++ {
		// Capture the current index to avoid closure issues.
		switchDesktopAction := func(desktop int) func() error {
			return func() error {
				dm.SwitchToDesktop(desktop)
				// Update system tray with current desktop
				if err := traySvc.UpdateDesktop(desktop + 1); err != nil {
					slog.Error("failed to update system tray", "error", err)
				}
				return nil
			}
		}(index)

		moveAndSwitchAction := func(desktop int) func() error {
			return func() error {
				// Retrieve the current foreground window from the OS and perform the move and switch.
				foregroundW := user.GetForegroundWindow()
				dm.MoveWindowToDesktop(foregroundW, desktop)
				dm.SwitchToDesktop(desktop)
				// Update system tray with current desktop
				if err := traySvc.UpdateDesktop(desktop + 1); err != nil {
					slog.Error("failed to update system tray", "error", err)
				}
				return nil
			}
		}(index)

		svc.RegisterKeyBindingActions(
			shortcut.NewBindingAction([]types.VirtualKey{types.VK_LMENU, desktopKeys[index]}, switchDesktopAction),
			shortcut.NewBindingAction([]types.VirtualKey{types.VK_LMENU, types.VK_LSHIFT, desktopKeys[index]}, moveAndSwitchAction),
		)
	}
	return svc
}

// captureEvents centralizes OS signal handling and keyboard event logging.
// By consolidating shutdown signal capture here, we ensure the application can terminate gracefully when needed.
func captureEvents(hook *keyboard.Hook) error {
	loggingSubscription := hook.Subscribe()
	go func() {
		for ev := range loggingSubscription {
			keyB := types.NewKeybinding(ev.PressedKeys...)

			// Simple event logging
			if ev.KeyDown {
				slog.Debug("key press", "key", ev.KeyCode.KeybindName(), "state", keyB.PrettyString())
			} else {
				slog.Debug("key release", "key", ev.KeyCode.KeybindName(), "state", keyB.PrettyString())
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	slog.Info("started")
	<-signalChan
	slog.Info("stopping")
	return nil
}

// Run aggregates the initialization of system components (desktop environment, keyboard hook, key bindings)
// and starts the user event loop. This separation of startup functionality enhances testability and maintainability.
func Run() error {
	// Load configuration
	cfg := config.LoadConfig()
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

	// Initialize the keyboard hook; early exit if setup fails to ensure proper system state.
	hook, err := keyboard.NewHook()
	if err != nil {
		return fmt.Errorf("failed to create keyboard hook: %w", err)
	}
	hook.Start()
	slog.Info("keyboard hook initialized")

	// Register keyboard shortcuts to facilitate rapid desktop management.
	keybindService := setupKeyBindings(hook, dm, traySvc)
	keybindService.Start()
	slog.Info("keyboard shortcuts registered")

	// Start capturing events; this blocks until an OS shutdown signal is received, ensuring graceful termination.
	return captureEvents(hook)
}
