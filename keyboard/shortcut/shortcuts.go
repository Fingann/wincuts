package shortcut

import (
	"fmt"
	"wincuts/keyboard"
)

// KeybindingService struct
type KeybindingService struct {
    eventChan    <- chan keyboard.KeyEvent 
    matcher      *Matcher
}

// NewKeybindingService creates a new KeybindingService
func NewKeybindingService(eventChan <- chan keyboard.KeyEvent) *KeybindingService {
    return &KeybindingService{
        eventChan:   eventChan,
        matcher:    NewMatcher(),
    }
}

// RegisterKeyCombo registers a key combination with an action
func (kbs *KeybindingService) RegisterKeyBindingActions(bindings ...*BindingAction) *KeybindingService {
    kbs.matcher.AddBindings(bindings...)
    return kbs
}

// Start starts the keybinding service to listen for events
func (kbs *KeybindingService) Start() {
    go func() {
        for event := range kbs.eventChan {
            fmt.Println("Received event")
            if !event.KeyDown {
                continue
            }
            prettyString, err := keyMapper.PrettyPrint(event.PressedKeys)
            if err != nil {
                fmt.Println(fmt.Errorf("failed to get pretty string: %v", err))
                continue
            }
            fmt.Println("Pressed keys:", prettyString)

            if binding := kbs.matcher.Match(event); binding != nil {
                fmt.Println("Matched binding")
                err := binding.Action()
                if err != nil {
                    fmt.Println(fmt.Errorf("failed to execute action: %v", err))
                }
            }
        }
    }()
}
