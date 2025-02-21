package shortcut

import "wincuts/keyboard/types"

type KeyEvent struct {
	PressedKeys []types.VirtualKey // Current state of all pressed keys
	KeyCode     types.VirtualKey   // The key involved in this event
	KeyDown     bool                // Whether this is a key press (true) or release (false)
}

func (ke *KeyEvent) CurrentState() []types.VirtualKey {
	return append(ke.PressedKeys, ke.KeyCode)
}

func (ke *KeyEvent) PressedKeyName() string {
	return types.VKCodeToKeyName[ke.KeyCode]
}

