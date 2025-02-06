package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VirtualKeyTestSuite groups related tests for the VirtualKey type
type VirtualKeyTestSuite struct {
	suite.Suite
}

func TestVirtualKeySuite(t *testing.T) {
	suite.Run(t, new(VirtualKeyTestSuite))
}

func (s *VirtualKeyTestSuite) TestVirtualKeyString() {
	// Test a known key with a mapping
	s.Equal("VK_LMENU (0xa4)", VK_LMENU.String(), "String representation of VK_LMENU should match expected format")

	// Test a key without a mapping (using a value that doesn't exist in the map)
	unknownKey := VirtualKey(0xFF)
	s.Equal(fmt.Sprintf("VK_%x", int(unknownKey)), unknownKey.String(), "Unknown key should use hex representation")
}

func (s *VirtualKeyTestSuite) TestVirtualKeyName() {
	// Test a known key with a mapping
	s.Equal("VK_LMENU", VK_LMENU.Name(), "Name of VK_LMENU should match expected format")

	// Test a key without a mapping
	unknownKey := VirtualKey(0xFF)
	s.Equal(fmt.Sprintf("VK_%#x", unknownKey), unknownKey.Name(), "Unknown key should use hex representation with VK_ prefix")
}

func (s *VirtualKeyTestSuite) TestVirtualKeyKeybindName() {
	// Test a known key with a mapping
	s.Equal("LMENU", VK_LMENU.KeybindName(), "KeybindName should return the name without VK_ prefix")

	// Test a key without a mapping
	unknownKey := VirtualKey(0xFF)
	s.Equal(fmt.Sprintf("%#x", unknownKey), unknownKey.KeybindName(), "Unknown key should use hex representation without prefix")
}

func (s *VirtualKeyTestSuite) TestIsModifier() {
	// Test known modifier keys
	s.True(VK_LMENU.IsModifier(), "LMENU should be recognized as a modifier")
	s.True(VK_LSHIFT.IsModifier(), "LSHIFT should be recognized as a modifier")
	s.True(VK_LCONTROL.IsModifier(), "LCONTROL should be recognized as a modifier")

	// Test non-modifier keys
	s.False(VK_A.IsModifier(), "A key should not be recognized as a modifier")
	s.False(VK_1.IsModifier(), "1 key should not be recognized as a modifier")
}

// KeyBindingTestSuite groups related tests for the KeyBinding type
type KeyBindingTestSuite struct {
	suite.Suite
}

func TestKeyBindingSuite(t *testing.T) {
	suite.Run(t, new(KeyBindingTestSuite))
}

func (s *KeyBindingTestSuite) TestNewKeybinding() {
	keys := []VirtualKey{VK_LMENU, VK_1}
	binding := NewKeybinding(keys...)

	s.Equal(len(keys), len(binding), "NewKeybinding should create a binding with the same number of keys")
	s.Equal(keys[0], binding[0], "First key should match")
	s.Equal(keys[1], binding[1], "Second key should match")
}

func (s *KeyBindingTestSuite) TestKeyBindingMatch() {
	binding := NewKeybinding(VK_LMENU, VK_1)

	// Test exact match
	s.True(binding.Match([]VirtualKey{VK_LMENU, VK_1}), "Should match identical key combination")

	// Test different order
	s.True(binding.Match([]VirtualKey{VK_1, VK_LMENU}), "Should match regardless of order")

	// Test non-match
	s.False(binding.Match([]VirtualKey{VK_LMENU, VK_2}), "Should not match different key combination")
	s.False(binding.Match([]VirtualKey{VK_LMENU}), "Should not match subset of keys")
}

func (s *KeyBindingTestSuite) TestKeyBindingSubsetOf() {
	subset := NewKeybinding(VK_LMENU)
	fullSet := NewKeybinding(VK_LMENU, VK_1)

	s.True(subset.SubsetOf(fullSet), "Single key should be recognized as subset")
	s.False(fullSet.SubsetOf(subset), "Larger set should not be recognized as subset of smaller set")
}

func (s *KeyBindingTestSuite) TestKeyBindingContains() {
	binding := NewKeybinding(VK_LMENU, VK_1)

	s.True(binding.Contains(VK_LMENU), "Should contain LMENU key")
	s.True(binding.Contains(VK_1), "Should contain 1 key")
	s.False(binding.Contains(VK_2), "Should not contain 2 key")
}

func (s *KeyBindingTestSuite) TestKeyBindingPrettyString() {
	binding := NewKeybinding(VK_1, VK_LMENU) // Intentionally out of order
	expected := "LMENU + 1"                  // Should be sorted with modifier first

	s.Equal(expected, binding.PrettyString(), "PrettyString should return sorted, human-readable representation")
}

// TestKeybindSort verifies that virtual keys are correctly sorted with modifiers first
func TestKeybindSort(t *testing.T) {
	assert := assert.New(t)

	// Create an unsorted slice of keys
	keys := []VirtualKey{
		VK_A,        // Non-modifier
		VK_LMENU,    // Modifier
		VK_1,        // Non-modifier
		VK_LCONTROL, // Modifier
	}

	KeybindSort(keys)

	// Verify modifiers come first and are sorted among themselves
	assert.True(keys[0].IsModifier(), "First key should be a modifier")
	assert.True(keys[1].IsModifier(), "Second key should be a modifier")
	assert.Equal(VK_LCONTROL, keys[0], "LCONTROL should come before LMENU")
	assert.Equal(VK_LMENU, keys[1], "LMENU should come after LCONTROL")

	// Verify non-modifiers are sorted
	assert.Equal(VK_1, keys[2], "Non-modifier keys should be sorted")
	assert.Equal(VK_A, keys[3], "Non-modifier keys should be sorted")
}
