package shortcut

import (
	"sync"
	"wincuts/keyboard"
)

// KeybindingService struct
type Service struct {
    eventChan    <- chan keyboard.KeyEvent 
    matcher      *Matcher
    wg          sync.WaitGroup
    stopChan   chan struct{}
}

// NewKeybindingService creates a new KeybindingService
func NewService(eventChan <- chan keyboard.KeyEvent) *Service {

    return &Service{
        eventChan:   eventChan,
        matcher:    NewMatcher(),
    }
}


// RegisterKeyCombo registers a key combination with an action
func (s *Service) RegisterKeyBindingActions(bindings ...KeyBindingAction) *Service {
    s.matcher.AddBindings(bindings...)
    return s
}


// Start starts the keybinding service to listen for events
func (s *Service) Start() {
    s.wg.Add(1)
    go func() {

        for {
        select {
        case <- s.stopChan:
            s.wg.Done()
            return

        case event := <- s.eventChan:
            s.matcher.Match(event)

        }
    }
    }()
}

// Stop stops the keybinding service
func (s *Service) Stop() {
    close(s.stopChan)
    s.wg.Wait()
}

