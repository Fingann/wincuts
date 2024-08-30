package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test KeybindSort
func TestKeybindSort(t *testing.T) {
	// Create a slice of VirtualKey objects
	vks := []VirtualKey{
		VK_CONTROL,
		VK_A,
		VK_MENU,
		VK_B,
		VK_C,

	}
	// Call the function
	KeybindSort(vks)
	// Check the result
	assert.Equal(t, VK_CONTROL, vks[0], fmt.Sprintf("Expected %s, got %s", VK_MENU.String(), vks[0].String()))
	assert.Equal(t, VK_MENU, vks[1], fmt.Sprintf("Expected %s, got %s", VK_CONTROL.String(), vks[1].String()))
	assert.Equal(t, VK_A, vks[2], fmt.Sprintf("Expected %s, got %s", VK_A.String(), vks[2].String()))
	assert.Equal(t, VK_B, vks[3], fmt.Sprintf("Expected %s, got %s", VK_B.String(), vks[3].String()))
	assert.Equal(t, VK_C, vks[4], fmt.Sprintf("Expected %s, got %s", VK_C.String(), vks[4].String()))

} 

