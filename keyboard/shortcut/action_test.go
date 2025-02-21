package shortcut

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"wincuts/keyboard/types"
)

// TestKeyBindingActionExecute verifies that Execute calls the underlying callback function.
func TestKeyBindingActionExecute(t *testing.T) {
	assert := assert.New(t)

	var executed bool
	action := func() error {
		executed = true
		return nil
	}
	// Create a binding action using the provided keys and action.
	keys := []types.VirtualKey{types.VK_LMENU, types.VK_1}
	bindingAction := NewBindingAction(keys, action, false)

	err := bindingAction.Execute()
	assert.NoError(err, "Execute should not return an error")
	assert.True(executed, "Expected action function to have been executed")
}

// TestKeyBindingActionMatch verifies that Match returns true for a matching event and false for non-matching events.
func TestKeyBindingActionMatch(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create a binding action; by default, NewBindingAction sets onKeyDown to false.
	action := func() error { return nil }
	keys := []types.VirtualKey{types.VK_LMENU, types.VK_1}
	bindingAction := NewBindingAction(keys, action, false)

	// Create a matching event: KeyDown is false and pressed keys exactly match.
	matchingEvent := KeyEvent{
		KeyDown:     false,
		PressedKeys: keys,
	}

	// Create a non-matching event where KeyDown differs.
	nonMatchingEvent := KeyEvent{
		KeyDown:     true,
		PressedKeys: keys,
	}

	// Create a non-matching event where the pressed keys differ.
	nonMatchingKeysEvent := KeyEvent{
		KeyDown:     false,
		PressedKeys: []types.VirtualKey{types.VK_LMENU, types.VK_LSHIFT},
	}

	// Verify matching behavior.
	require.True(bindingAction.Match(matchingEvent), "Expected Match to return true for matching event")
	assert.False(bindingAction.Match(nonMatchingEvent), "Expected Match to return false when KeyDown state differs")
	assert.False(bindingAction.Match(nonMatchingKeysEvent), "Expected Match to return false for non-matching keys")
}
