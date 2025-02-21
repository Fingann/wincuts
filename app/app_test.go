//go:build windows

package app

import (
	"testing"

	winapi "github.com/chrsm/winapi"
)

// fakeDesktopManager is used as a test double to simulate the DesktopManager interface.
// We use it to verify that the EnsureMinimumDesktops function triggers the expected desktop creation behavior.
// This abstraction allows testing without relying on external dependencies or side effects.

type fakeDesktopManager struct {
	count           int
	createdDesktops int
}

// GetCurrentDesktopCount provides the simulated current desktop count so that tests can verify state changes.
func (f *fakeDesktopManager) GetCurrentDesktopCount() int {
	return f.count
}

// CreateNewDesktop simulates the effect of creating a new desktop by updating internal counters.
// We do this to ensure that EnsureMinimumDesktops makes the correct number of creation calls.
func (f *fakeDesktopManager) CreateNewDesktop() {
	f.createdDesktops++
	f.count++
}

// SwitchToDesktop is provided to satisfy the DesktopManager interface but is not required for this test scenario.
func (f *fakeDesktopManager) SwitchToDesktop(desktopNumber int) {}

// MoveWindowToDesktop is provided to satisfy the DesktopManager interface and is unused in the context of this test.
func (f *fakeDesktopManager) MoveWindowToDesktop(window winapi.HWND, desktopNumber int) {}

// TestEnsureMinimumDesktops asserts that EnsureMinimumDesktops triggers the correct number of desktop creation operations.
// This ensures that the application will enforce a required minimum number of desktops at runtime.
func TestEnsureMinimumDesktops(t *testing.T) {
	initialCount := 5
	minimumRequired := 9
	fake := &fakeDesktopManager{count: initialCount}

	EnsureMinimumDesktops(fake, minimumRequired)

	expectedCreated := minimumRequired - initialCount
	if fake.createdDesktops != expectedCreated {
		t.Errorf("Expected %d desktops to be created, got %d", expectedCreated, fake.createdDesktops)
	}

	if fake.count != minimumRequired {
		t.Errorf("Expected final desktop count %d, got %d", minimumRequired, fake.count)
	}
}
