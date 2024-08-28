package shortcut

import (
	"wincuts/keyboard/code"
)

type KeyBindingFunc func() error

var keyMapper code.Mapper = code.NewMapper()

func NewKeyCombinations(keys []string, exact bool) (*KeyBinding,error) {
	keyCodes,err := keyMapper.KeysToCodeMap(keys)
	if err != nil {
		return nil,err
	}
	return &KeyBinding{
		VKMap:  keyCodes,
		Exact: exact,
	},nil
}

type KeyBinding struct {
	VKMap map[uint32]string 
	Exact bool
}

func NewBindingActionFromBinding(binding KeyBinding, action KeyBindingFunc) (*BindingAction,error) {
	return &BindingAction{
		Binding: binding,
		Action:      action,
	},nil
}

func NewBindingAction(keys []string, exact bool, action KeyBindingFunc) (*BindingAction,error) {
	binding,err := NewKeyCombinations(keys, exact)
	if err != nil {
		return nil,err
	}
	return NewBindingActionFromBinding(*binding, action)
}

type BindingAction struct {
	Binding KeyBinding 
	Action      KeyBindingFunc 
}


