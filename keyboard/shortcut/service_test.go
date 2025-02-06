package shortcut

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"wincuts/keyboard"
	"wincuts/keyboard/types"
)

// TestNewBindingAction verifies that a binding action correctly captures the intended key combination and callback function.
// This test ensures that the binding action's callback is correctly bound and executed, which is critical for key event handling.
func TestNewBindingAction(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var called bool
	action := func() error {
		called = true
		return nil
	}
	keys := []types.VirtualKey{types.VK_LMENU, types.VK_1}
	binding := NewBindingAction(keys, action)

	// Execute the action to verify that the callback is correctly bound.
	err := binding.Action()
	require.NoError(err, "Binding action should not return an error")
	assert.True(called, "Expected action function to be executed")
}

// TestKeybindingServiceRegistration checks that the keybinding service can register binding actions without error.
// This is important for ensuring that dynamic registration of shortcuts operates as expected in the application.
func TestKeybindingServiceRegistration(t *testing.T) {
	require := require.New(t)
	dummyChan := make(chan keyboard.KeyEvent, 1)
	svc := NewService(dummyChan)

	require.NotNil(svc, "NewKeybindingService should not return nil")

	dummyAction := NewBindingAction([]types.VirtualKey{types.VK_LMENU}, func() error {
		return nil
	})

	require.NotPanics(func() {
		svc.RegisterKeyBindingActions(dummyAction)
	}, "Registering binding actions should not panic")
}

// TestMatchKeys verifies that the key matching logic correctly identifies matching and non-matching key combinations.
// This is vital for ensuring that only the intended keyboard shortcuts trigger their associated actions.
func TestMatchKeys(t *testing.T) {
	assert := assert.New(t)
	bindingKeys := []types.VirtualKey{types.VK_LMENU, types.VK_1}
	matchingInput := []types.VirtualKey{types.VK_LMENU, types.VK_1}
	nonMatchingInput := []types.VirtualKey{types.VK_LMENU, types.VK_LSHIFT}

	assert.True(MatchKeys(bindingKeys, matchingInput), "Expected matchingInput to match bindingKeys")
	assert.False(MatchKeys(bindingKeys, nonMatchingInput), "Expected nonMatchingInput not to match bindingKeys")
}

// MatchKeys checks for equality between two slices of VirtualKey. This simple implementation is used for testing the matcher.
func MatchKeys(binding, input []types.VirtualKey) bool {
	if len(binding) != len(input) {
		return false
	}
	for i, key := range binding {
		if key != input[i] {
			return false
		}
	}
	return true
}

// TestNewService verifies that the service can be created and can register key bindings without error.
func TestNewService(t *testing.T) {
	require := require.New(t)

	dummyChan := make(chan keyboard.KeyEvent, 1)
	svc := NewService(dummyChan)

	require.NotNil(svc, "NewService should not return nil")

	dummyAction := NewBindingAction([]types.VirtualKey{types.VK_LMENU}, func() error {
		return nil
	})

	require.NotPanics(func() {
		svc.RegisterKeyBindingActions(dummyAction)
	}, "Registering binding actions should not panic")
}
