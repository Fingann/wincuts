package shortcut


// Matcher handles matching key events to registered shortcuts
type Matcher struct {
	bindings []KeyBindingAction
}

// NewMatcher creates a new Matcher instance
func NewMatcher() *Matcher {
	return &Matcher{
		bindings: make([]KeyBindingAction, 0),
	}
}

// AddBindings registers new key binding actions
func (m *Matcher) AddBindings(bindings ...KeyBindingAction) {
	m.bindings = append(m.bindings, bindings...)
}

// Match checks if the event matches any registered shortcut
func (m *Matcher) Match(event KeyEvent) (*KeyBindingAction, bool) {
	for _, binding := range m.bindings {
		if binding.Match(event) {
			return &binding, true
		}
	}
	return nil, false
}
