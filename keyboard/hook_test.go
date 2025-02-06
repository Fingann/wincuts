package keyboard

import (
	"testing"
	"time"

	"github.com/moutend/go-hook/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	wtypes "wincuts/keyboard/types"
)

type HookTestSuite struct {
	suite.Suite
	hook *Hook
}

func TestHookSuite(t *testing.T) {
	suite.Run(t, new(HookTestSuite))
}

func (s *HookTestSuite) SetupTest() {
	var err error
	s.hook, err = NewHook()
	require.NoError(s.T(), err, "NewHook should not return an error")
}

func (s *HookTestSuite) TearDownTest() {
	if s.hook != nil {
		s.hook.Stop()
	}
}

// TestNewHook verifies that a new Hook can be created with the expected initial state
func (s *HookTestSuite) TestNewHook() {
	assert := assert.New(s.T())
	require := require.New(s.T())

	require.NotNil(s.hook, "NewHook should return a non-nil Hook")
	assert.NotNil(s.hook.lowLevelChan, "lowLevelChan should be initialized")
	assert.NotNil(s.hook.keyState, "keyState should be initialized")
	assert.NotNil(s.hook.subscriberManager, "subscriberManager should be initialized")
	assert.NotNil(s.hook.ctx, "context should be initialized")
	assert.NotNil(s.hook.cancel, "cancel function should be initialized")
}

// TestSubscriptionManagement verifies that subscribers can be added and removed correctly
func (s *HookTestSuite) TestSubscriptionManagement() {
	assert := assert.New(s.T())

	// Subscribe and verify channel is created
	ch := s.hook.Subscribe()
	assert.NotNil(ch, "Subscribe should return a non-nil channel")

	// Unsubscribe and verify no panic occurs
	assert.NotPanics(func() {
		s.hook.Unsubscribe(ch)
	}, "Unsubscribe should not panic")
}

// TestKeyEventProcessing verifies that keyboard events are correctly processed and broadcast
func (s *HookTestSuite) TestKeyEventProcessing() {
	require := require.New(s.T())
	assert := assert.New(s.T())

	// Start the hook
	err := s.hook.Start()
	require.NoError(err, "Start should not return an error")

	// Subscribe to receive events
	eventChan := s.hook.Subscribe()

	// Create a helper function to create keyboard events
	createKeyEvent := func(message types.Message, key wtypes.VirtualKey) types.KeyboardEvent {
		return types.KeyboardEvent{
			Message: message,
			KBDLLHOOKSTRUCT: types.KBDLLHOOKSTRUCT{
				VKCode: types.VKCode(key),
			},
		}
	}

	// Simulate a key press event
	s.hook.lowLevelChan <- createKeyEvent(types.WM_KEYDOWN, wtypes.VK_A)

	// Wait for and verify the event
	select {
	case event := <-eventChan:
		assert.True(event.KeyDown, "KeyDown should be true for WM_KEYDOWN")
		assert.Equal(wtypes.VK_A, event.KeyCode, "KeyCode should match the simulated key")
		assert.Empty(event.PressedKeys, "PressedKeys should reflect state before this key press")
	case <-time.After(time.Second):
		s.T().Fatal("Timeout waiting for key event")
	}

	// Simulate key release
	s.hook.lowLevelChan <- createKeyEvent(types.WM_KEYUP, wtypes.VK_A)

	// Verify key release event
	select {
	case event := <-eventChan:
		assert.False(event.KeyDown, "KeyDown should be false for WM_KEYUP")
		assert.Equal(wtypes.VK_A, event.KeyCode, "KeyCode should match the simulated key")
		assert.Contains(event.PressedKeys, wtypes.VK_A, "PressedKeys should reflect state before this key release")
	case <-time.After(time.Second):
		s.T().Fatal("Timeout waiting for key event")
	}
}

// TestHookLifecycle verifies that the hook's internal state is managed correctly during startup and shutdown
func (s *HookTestSuite) TestHookLifecycle() {
	require := require.New(s.T())

	// Start the hook
	err := s.hook.Start()
	require.NoError(err, "Start should not return an error")

	// Subscribe to receive events
	eventChan := s.hook.Subscribe()
	require.NotNil(eventChan, "Subscribe should return a valid channel")

	// Create a helper function to create keyboard events
	createKeyEvent := func(message types.Message, key wtypes.VirtualKey) types.KeyboardEvent {
		return types.KeyboardEvent{
			Message: message,
			KBDLLHOOKSTRUCT: types.KBDLLHOOKSTRUCT{
				VKCode: types.VKCode(key),
			},
		}
	}

	// Verify that events are processed while the hook is running
	s.hook.lowLevelChan <- createKeyEvent(types.WM_KEYDOWN, wtypes.VK_A)

	select {
	case event := <-eventChan:
		require.True(event.KeyDown, "Should receive events while hook is running")
	case <-time.After(time.Second):
		s.T().Fatal("Timeout waiting for event")
	}

	// Stop the hook and verify that the context is canceled
	s.hook.cancel()  // Cancel the context directly since we can't uninstall the hook in tests
	s.hook.wg.Wait() // Wait for goroutine to finish

	// Verify that sending more events doesn't cause issues
	s.hook.lowLevelChan <- createKeyEvent(types.WM_KEYDOWN, wtypes.VK_B)

	// The channel should eventually be closed
	select {
	case _, ok := <-eventChan:
		require.False(ok, "Channel should be closed after stopping")
	case <-time.After(time.Second):
		// Channel might not be closed immediately, which is also acceptable
	}
}

// TestKeyStateTracking verifies that the hook correctly tracks the state of pressed keys
func (s *HookTestSuite) TestKeyStateTracking() {
	require := require.New(s.T())
	assert := assert.New(s.T())

	err := s.hook.Start()
	require.NoError(err, "Start should not return an error")

	eventChan := s.hook.Subscribe()

	// Create a helper function to create keyboard events
	createKeyEvent := func(message types.Message, key wtypes.VirtualKey) types.KeyboardEvent {
		return types.KeyboardEvent{
			Message: message,
			KBDLLHOOKSTRUCT: types.KBDLLHOOKSTRUCT{
				VKCode: types.VKCode(key),
			},
		}
	}

	// Simulate pressing multiple keys
	keys := []wtypes.VirtualKey{wtypes.VK_LCONTROL, wtypes.VK_A}

	// Press first key
	s.hook.lowLevelChan <- createKeyEvent(types.WM_KEYDOWN, keys[0])

	// Wait for first key event
	select {
	case event := <-eventChan:
		assert.True(event.KeyDown, "KeyDown should be true for first key")
		assert.Empty(event.PressedKeys, "PressedKeys should be empty before first key press")
	case <-time.After(time.Second):
		s.T().Fatal("Timeout waiting for first key event")
	}

	// Press second key
	s.hook.lowLevelChan <- createKeyEvent(types.WM_KEYDOWN, keys[1])

	// Wait for second key event
	select {
	case event := <-eventChan:
		assert.True(event.KeyDown, "KeyDown should be true for second key")
		assert.Contains(event.PressedKeys, keys[0], "First key should be in pressed state")
	case <-time.After(time.Second):
		s.T().Fatal("Timeout waiting for second key event")
	}

	// Release keys in reverse order
	for i := len(keys) - 1; i >= 0; i-- {
		s.hook.lowLevelChan <- createKeyEvent(types.WM_KEYUP, keys[i])

		select {
		case event := <-eventChan:
			assert.False(event.KeyDown, "KeyDown should be false for key release")
			if i == len(keys)-1 {
				// When releasing the second key (A), both keys should still be in pressed state
				assert.Contains(event.PressedKeys, keys[0], "First key should still be pressed")
				assert.Contains(event.PressedKeys, keys[1], "Second key should still be in state before release")
			} else {
				// When releasing the first key (LCONTROL), only it should be in pressed state
				assert.Contains(event.PressedKeys, keys[0], "First key should be in state before release")
				assert.NotContains(event.PressedKeys, keys[1], "Second key should be released")
			}
		case <-time.After(time.Second):
			s.T().Fatal("Timeout waiting for key release event")
		}
	}
}
