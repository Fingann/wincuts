package shortcut

import "wincuts/keyboard"

type Matcher struct {
	Bindings []*BindingAction
}

func NewMatcher(bindings ...*BindingAction) *Matcher {
	return &Matcher{
		Bindings: bindings,
	}
}

func (kbm *Matcher) AddBindings(binding ...*BindingAction) {
	kbm.Bindings = append(kbm.Bindings, binding...)
}


// Match returns the keybinding that matches the event
// or nil if no match is found.
func (kbm *Matcher) Match(event keyboard.KeyEvent) (*BindingAction) {
	for _, binding := range kbm.Bindings {
		if len(event.PressedKeys) != len(binding.Binding.VKMap) {
			return nil		}
		for _, key := range event.PressedKeys {
			if _, exists := binding.Binding.VKMap[key]; !exists {
				return nil
			}
		}
		return binding
	}
	return nil
}
