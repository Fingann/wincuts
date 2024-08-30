package shortcut

import (
	"fmt"
	"sync"
	"wincuts/keyboard"
)

// KeybindingService struct
type KeybindingService struct {
    eventChan    <- chan keyboard.KeyEvent 
    matcher      *Matcher
    wg          sync.WaitGroup
    stopChan   chan struct{}
}

// NewKeybindingService creates a new KeybindingService
func NewKeybindingService(eventChan <- chan keyboard.KeyEvent) *KeybindingService {
    return &KeybindingService{
        eventChan:   eventChan,
        matcher:    NewMatcher(),
    }
}

// RegisterKeyCombo registers a key combination with an action
func (kbs *KeybindingService) RegisterKeyBindingActions(bindings ...KeyBindingAction) *KeybindingService {
    kbs.matcher.AddBindings(bindings...)
    return kbs
}

// Start starts the keybinding service to listen for events
func (kbs *KeybindingService) Start() {
    kbs.wg.Add(1)
    go func() {
        select {
        case <- kbs.stopChan:
            kbs.wg.Done()
            return
        case event := <- kbs.eventChan:
            kbs.matcher.Match(event)
        }
    }()
}

// Stop stops the keybinding service
func (kbs *KeybindingService) Stop() {
    close(kbs.stopChan)
    kbs.wg.Wait()
}
