package shortcut

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wincuts/keyboard/types"
)

// TestMatcherMatch_ExecutesActionOnMatch verifies that the Matcher executes the action when a matching key event is received.
func TestMatcherMatch_ExecutesActionOnMatch(t *testing.T) {
	assert := assert.New(t)

	var executed bool
	action := func() error {
		executed = true
		return nil
	}
	keys := []types.VirtualKey{types.VK_LMENU, types.VK_1}
	binding := NewBindingAction(keys, action, false)

	matcher := NewMatcher()
	matcher.AddBindings(binding)

	// Create a key event that matches the binding. Since NewBindingAction sets onKeyDown = false by default,
	// the event's KeyDown should be false and the PressedKeys should match keys.
	event := KeyEvent{
		KeyDown:     false,
		PressedKeys: keys,
	}

	matcher.Match(event)
	assert.True(executed, "Expected binding action to be executed for matching event")
}

// TestMatcherMatch_DoesNotExecuteOnNoMatch verifies that the Matcher does not execute the action when a non-matching event is received.
func TestMatcherMatch_DoesNotExecuteOnNoMatch(t *testing.T) {
	assert := assert.New(t)

	var executed bool
	action := func() error {
		executed = true
		return nil
	}
	keys := []types.VirtualKey{types.VK_LMENU, types.VK_1}
	binding := NewBindingAction(keys, action, false)

	matcher := NewMatcher()
	matcher.AddBindings(binding)

	// Create a key event that does not match (different key combination).
	event := KeyEvent{
		KeyDown:     false,
		PressedKeys: []types.VirtualKey{types.VK_LMENU, types.VK_LSHIFT},
	}

	matcher.Match(event)
	assert.False(executed, "Expected binding action not to be executed for non-matching event")
}
