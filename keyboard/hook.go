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

type Hook struct {
	lowLevelChan      chan types.KeyboardEvent
	keyState          map[wtypes.VirtualKey]bool
	subscriberManager *SubscriberManager[KeyEvent]
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
}

type KeyEvent struct {
	PressedKeys []wtypes.VirtualKey
	KeyCode     wtypes.VirtualKey
	KeyDown     bool
}

func NewHook() (*Hook, error) {
	lowLevelChan := make(chan types.KeyboardEvent, 100)
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

func (h *Hook) Start() error {
	if err := keyboard.Install(nil, h.lowLevelChan); err != nil {
		return fmt.Errorf("failed to install keyboard hook: %v", err)
	}

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		for {
			select {
			case <-h.ctx.Done():
				// Clean up and exit the goroutine
				return
			case k := <-h.lowLevelChan:
				if k.Message == 0x0312 {
					fmt.Println("Key kombo")
				}
				vCode := wtypes.VirtualKey(k.VKCode)
				isKeyDown := k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN

				event := KeyEvent{
					PressedKeys: slices.Collect(maps.Keys(h.keyState)),
					KeyCode:     vCode,
					KeyDown:     isKeyDown,
				}

				h.subscriberManager.Broadcast(event)
				if isKeyDown {
					h.keyState[vCode] = true
				} else {
					delete(h.keyState, vCode)
				}
			}
		}
	}()

	return nil
}

func (h *Hook) Stop() error {
	h.cancel()  // Signal the goroutine to stop
	h.wg.Wait() // Wait for the goroutine to finish

	if err := keyboard.Uninstall(); err != nil {
		return err
	}

	close(h.lowLevelChan)          // Close the low-level channel
	h.subscriberManager.CloseAll() // Close all subscriber channels

	return nil
}

func (h *Hook) Subscribe() chan KeyEvent {
	return h.subscriberManager.AddSubscriber()
}

func (h *Hook) Unsubscribe(ch chan KeyEvent) {
	h.subscriberManager.RemoveSubscriber(ch)
}
