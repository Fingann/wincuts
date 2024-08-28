package keylogger 

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// KeyEvent represents a key event with associated modifiers.
type KeyEvent struct {
	Keys      []string // Keys pressed (e.g., ["A", "P"])
	Modifiers []string // Modifiers pressed (e.g., ["Shift", "Ctrl"])
}

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procSetWindowsHookEx = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx   = user32.NewProc("CallNextHookEx")
	procGetMessage       = user32.NewProc("GetMessageW")
	procGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
)

// Virtual key codes
const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101

	VK_SHIFT uint32 = 160
	VK_ALT   uint32 = 164
	VK_CTRL  uint32 = 162
	VK_A     uint32 = 65
	VK_P     uint32 = 80
)

// SystemState holds the state of the modifier keys and the channel to report events.
type SystemState struct {
	ShiftPressed bool
	CtrlPressed  bool
	AltPressed   bool
	KeyBuffer    []string
	EventChan    chan KeyEvent
}

// HandleKeyPress handles key press and release events.
func (ss *SystemState) HandleKeyPress(vkCode uint32, keyDown bool) {
	var keyName string
	switch vkCode {
	case VK_SHIFT:
		ss.ShiftPressed = keyDown
		keyName = "Shift"
	case VK_CTRL:
		ss.CtrlPressed = keyDown
		keyName = "Ctrl"
	case VK_ALT:
		ss.AltPressed = keyDown
		keyName = "Alt"
	case VK_A:
		keyName = "A"
	case VK_P:
		keyName = "P"
	}

	if keyName != "" && keyDown {
		ss.KeyBuffer = append(ss.KeyBuffer, keyName)
	} else if keyName != "" && !keyDown {
		// On key up, report the event and clear the buffer
		ss.reportEvent()
	}
}

// reportEvent sends the current key combination to the event channel.
func (ss *SystemState) reportEvent() {
	if len(ss.KeyBuffer) == 0 {
		return
	}

	modifiers := []string{}
	if ss.ShiftPressed {
		modifiers = append(modifiers, "Shift")
	}
	if ss.CtrlPressed {
		modifiers = append(modifiers, "Ctrl")
	}
	if ss.AltPressed {
		modifiers = append(modifiers, "Alt")
	}

	event := KeyEvent{
		Keys:      append([]string{}, ss.KeyBuffer...), // Copy the key buffer
		Modifiers: modifiers,
	}

	ss.EventChan <- event
	ss.KeyBuffer = ss.KeyBuffer[:0] // Clear the key buffer after reporting
}

// KeyboardProc is the callback method for the keyboard hook.
func (ss *SystemState) KeyboardProc(nCode int32, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 {
		kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		keyDown := wParam == WM_KEYDOWN
		ss.HandleKeyPress(kbdstruct.VKCode, keyDown)
	}

	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

// StartKeyEventCapture sets up the hook and returns the event channel.
func StartKeyEventCapture() (<-chan KeyEvent, error) {
	systemState := &SystemState{
		EventChan: make(chan KeyEvent, 100), // Buffered channel
	}

	hInstance, _, err := procGetModuleHandle.Call(0)
	if hInstance == 0 {
		return nil, fmt.Errorf("GetModuleHandle failed: %v", err)
	}

	r1, _, err := procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		windows.NewCallback(systemState.KeyboardProc),
		hInstance,
		0,
	)

	if r1 == 0 {
		return nil, fmt.Errorf("SetWindowsHookEx failed: %v", err)
	}

	return systemState.EventChan, nil
}

type KBDLLHOOKSTRUCT struct {
	VKCode    uint32
	ScanCode  uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type MSG struct {
	HWND    uintptr
	Message uint32
	WPARAM  uintptr
	LPARAM  uintptr
	Time    uint32
	Pt      POINT
}

type POINT struct {
	X, Y int32
}
