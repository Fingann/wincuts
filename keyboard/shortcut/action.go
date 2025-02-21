package shortcut

import (
	"wincuts/keyboard/types"
)

type KeyBindingFunc func() error

type KeyBindingAction struct {
	Binding types.KeyBinding 
	Action 	KeyBindingFunc
	onKeyDown bool
	ShouldBlock bool
}

func (kba *KeyBindingAction) Execute() error {
	return kba.Action()
}

func (kba *KeyBindingAction) Match(event KeyEvent) bool {
	if event.KeyDown != kba.onKeyDown{
		return false
	}
	
	return kba.Binding.Match(event.PressedKeys)
} 

func NewBindingActionFromBinding(binding types.KeyBinding, action KeyBindingFunc) (KeyBindingAction) {
	return KeyBindingAction{
		Binding: binding,
		Action:      action,
		onKeyDown: false,
		ShouldBlock: false,
	}

}

func NewBindingAction(keys []types.VirtualKey, action KeyBindingFunc, shouldBlock bool) (KeyBindingAction) {
	return NewBindingActionFromBinding(types.NewKeybinding(keys...),action)
}
