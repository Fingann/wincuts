package shortcut

import (
	"fmt"
	"wincuts/keyboard"
)

type Matcher struct {
	Bindings []KeyBindingAction
}

func NewMatcher(bindings ...KeyBindingAction) *Matcher {
	return &Matcher{
		Bindings: bindings,
	}
}

func (kbm *Matcher) AddBindings(binding ...KeyBindingAction) {
	kbm.Bindings = append(kbm.Bindings, binding...)
}


// Match returns the keybinding that matches the event
// or nil if no match is found.
func (kbm *Matcher) Match(event keyboard.KeyEvent) {
	for _, binding := range kbm.Bindings {
		if binding.Match(event) {
			if err:= binding.Execute(); err != nil {
				fmt.Println(fmt.Errorf("failed to execute action: %v", err))
			}
		}
	}
}
