// Package keyboard abstracts the setup and management of low-level keyboard hooks.
// We structure this package to isolate system-specific input handling from our core application logic,
// thereby improving testability and modularity.

// The NewHook function configures a keyboard hook to capture input events essential for key binding operations.

package keyboard

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"sync"

	"wincuts/keyboard/shortcut"
	wtypes "wincuts/keyboard/types"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// Hook manages keyboard event capturing and distribution to subscribers.
// It maintains thread-safe state tracking of currently pressed keys.
type Hook struct {
	lowLevelChan      chan types.KeyboardEvent
	keyState          map[wtypes.VirtualKey]bool
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	stateMutex        sync.RWMutex                   // Protects keyState map
	shortcutChan      chan *shortcut.KeyBindingAction // Channel for matched shortcuts
	shortcutService   *shortcut.Service
}

// KeyEvent represents a processed keyboard event with additional context.


const (
	// DefaultEventBuffer is the size of the channel buffer for keyboard events
	DefaultEventBuffer = 100
)

// NewHook creates a new keyboard hook with default configuration.
func NewHook(shortcutService *shortcut.Service) (*Hook, error) {
	lowLevelChan := make(chan types.KeyboardEvent, DefaultEventBuffer)
	ctx, cancel := context.WithCancel(context.Background())



	return &Hook{
		lowLevelChan:      lowLevelChan,
		keyState:          make(map[wtypes.VirtualKey]bool),
		ctx:               ctx,
		cancel:            cancel,
		shortcutChan:      shortcutService.GetShortcutChan(), // Buffered channel
		shortcutService:   shortcutService,
	}, nil
}


// updateKeyState safely updates the state of a key and returns a snapshot of all pressed keys.
func (h *Hook) updateKeyState(key wtypes.VirtualKey, isDown bool) []wtypes.VirtualKey {
	h.stateMutex.Lock()
	defer h.stateMutex.Unlock()

	if isDown {
		h.keyState[key] = true
	} else {
		delete(h.keyState, key)
	}

	// Create a snapshot of currently pressed keys
	return slices.Collect(maps.Keys(h.keyState))
}

// getCurrentKeyState safely returns a snapshot of all currently pressed keys.
func (h *Hook) getCurrentKeyState() []wtypes.VirtualKey {
	h.stateMutex.RLock()
	defer h.stateMutex.RUnlock()
	return slices.Collect(maps.Keys(h.keyState))
}

// GetShortcutChan returns the channel for matched shortcuts
func (h *Hook) GetShortcutChan() <-chan *shortcut.KeyBindingAction {
	return h.shortcutChan
}

// Start begins capturing keyboard events and distributing them to subscribers.
func (h *Hook) Start() error {
	if err := keyboard.Install(nil, h.lowLevelChan); err != nil {
		return fmt.Errorf("failed to install keyboard hook: %w", err)
	}

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		for {
			select {
			case <-h.ctx.Done():
				return
			case k := <-h.lowLevelChan:
				if k.Message == 0x0312 {
					fmt.Println("Key combo")
					continue
				}

				vCode := wtypes.VirtualKey(k.VKCode)
				isKeyDown := k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN

				// Get current state before updating
				currentState := h.getCurrentKeyState()

				// Create event with state before the update
				event := shortcut.KeyEvent{
					PressedKeys: currentState,
					KeyCode:     vCode,
					KeyDown:     isKeyDown,
				}
				if isKeyDown {
					slog.Debug("key press", "key", vCode.KeybindName(), "state", currentState)
				} else {
					slog.Debug("key release", "key", vCode.KeybindName(), "state", currentState)
				}

				// Update state after creating the event
				h.updateKeyState(vCode, isKeyDown)

				// Check if this is a registered shortcut
				keyBinding, found := h.shortcutService.Match(event)
				if !found {
					continue
				}

				// Send the matched shortcut through the channel
				select {
				case h.shortcutChan <- keyBinding:
					slog.Debug("sent shortcut", "binding", keyBinding.Binding.PrettyString())
				default:
					// Drop the event if the channel is full
				}

			}
		}
	}()

	return nil
}

// Stop gracefully shuts down the hook and cleans up resources.
func (h *Hook) Stop() error {
	h.cancel()  // Signal the goroutine to stop
	h.wg.Wait() // Wait for goroutine to finish

	if err := keyboard.Uninstall(); err != nil {
		return fmt.Errorf("keyboard: failed to uninstall hook function: %w", err)
	}

	close(h.lowLevelChan)
	// Clear the key state
	h.stateMutex.Lock()
	h.keyState = make(map[wtypes.VirtualKey]bool)
	h.stateMutex.Unlock()

	return nil
}