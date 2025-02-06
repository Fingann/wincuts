// Package keyboard abstracts the setup and management of low-level keyboard hooks.
// We structure this package to isolate system-specific input handling from our core application logic,
// thereby improving testability and modularity.

// The NewHook function configures a keyboard hook to capture input events essential for key binding operations.

package keyboard

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"

	wtypes "wincuts/keyboard/types"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// Hook manages keyboard event capturing and distribution to subscribers.
// It maintains thread-safe state tracking of currently pressed keys.
type Hook struct {
	lowLevelChan      chan types.KeyboardEvent
	keyState          map[wtypes.VirtualKey]bool
	subscriberManager *SubscriberManager[KeyEvent]
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	stateMutex        sync.RWMutex // Protects keyState map
}

// KeyEvent represents a processed keyboard event with additional context.
type KeyEvent struct {
	PressedKeys []wtypes.VirtualKey // Current state of all pressed keys
	KeyCode     wtypes.VirtualKey   // The key involved in this event
	KeyDown     bool                // Whether this is a key press (true) or release (false)
}

const (
	// DefaultEventBuffer is the size of the channel buffer for keyboard events
	DefaultEventBuffer = 100
)

// NewHook creates a new keyboard hook with default configuration.
func NewHook() (*Hook, error) {
	lowLevelChan := make(chan types.KeyboardEvent, DefaultEventBuffer)
	ctx, cancel := context.WithCancel(context.Background())
	subscriberManager := NewSubscriberManager[KeyEvent]()

	return &Hook{
		lowLevelChan:      lowLevelChan,
		keyState:          make(map[wtypes.VirtualKey]bool),
		subscriberManager: subscriberManager,
		ctx:               ctx,
		cancel:            cancel,
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
				event := KeyEvent{
					PressedKeys: currentState,
					KeyCode:     vCode,
					KeyDown:     isKeyDown,
				}

				// Update state after creating the event
				h.updateKeyState(vCode, isKeyDown)

				// Broadcast the event
				h.subscriberManager.Broadcast(event)
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
	h.subscriberManager.CloseAll()

	// Clear the key state
	h.stateMutex.Lock()
	h.keyState = make(map[wtypes.VirtualKey]bool)
	h.stateMutex.Unlock()

	return nil
}

// Subscribe returns a new channel that will receive keyboard events.
func (h *Hook) Subscribe() chan KeyEvent {
	return h.subscriberManager.AddSubscriber()
}

// Unsubscribe removes a subscriber and closes their channel.
func (h *Hook) Unsubscribe(ch chan KeyEvent) {
	h.subscriberManager.RemoveSubscriber(ch)
}
