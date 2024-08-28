package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procSetWindowsHookEx = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx   = user32.NewProc("CallNextHookEx")
	procGetMessage       = user32.NewProc("GetMessageW")
	procGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101

	// Virtual key codes
	VK_SHIFT uint32 = 160
	VK_ALT   uint32 = 164
	VK_CTRL  uint32 = 162
	VK_A     uint32 = 65
)

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

// SystemState holds the state of the modifier keys and manages key events
type SystemState struct {
	keyState    map[uint32]bool
	EventChannel chan KeyBoardEvent
}

func NewSystemState() *SystemState {
	return &SystemState{
		keyState:    make(map[uint32]bool),
		EventChannel: make(chan KeyBoardEvent),
	}
}

type KeyBoardEvent struct {
	PressedKeys map[uint32]bool
	KeyCode     uint32
	KeyDown     bool
}

// HandleKeyPress handles key press events and sends them to the event channel
func (ss *SystemState) HandleKeyPress(vkCode uint32, keyDown bool) {
	if keyDown {
		ss.keyState[vkCode] = true
	} else {
		delete(ss.keyState, vkCode)
	}

	ss.EventChannel <- KeyBoardEvent{
		PressedKeys: copyMap(ss.keyState),
		KeyCode:     vkCode,
		KeyDown:     keyDown,
	}
}

// copyMap creates a copy of the key state map
func copyMap(original map[uint32]bool) map[uint32]bool {
	copy := make(map[uint32]bool)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// KeyboardProc is the callback method for the keyboard hook
func (ss *SystemState) KeyboardProc(nCode int32, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 {
		kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		keyDown := wParam == WM_KEYDOWN
		ss.HandleKeyPress(kbdstruct.VKCode, keyDown)
	}
	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

// start sets up the global keyboard hook and processes key events
func (ss *SystemState) start() error {
	hInstance, _, err := procGetModuleHandle.Call(0)
	if hInstance == 0 {
		return fmt.Errorf("GetModuleHandle failed: %v", err)
	}

	r1, _, err := procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		windows.NewCallback(ss.KeyboardProc),
		hInstance,
		0,
	)

	if r1 == 0 {
		return fmt.Errorf("SetWindowsHookEx failed: %v", err)
	}

	fmt.Println("Global hook set, press Shift, Ctrl, or Alt to test")

	go func() {
		var msg MSG
		for {
			ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if ret == 0 {
				break
			}
		}
	}()

	return nil
}

func main() {
	systemState := NewSystemState()

	// Start the system state which sets the hook and processes key events
	if err := systemState.start(); err != nil {
		fmt.Println(err)
		return
	}

	// Process key events in a separate goroutine
	go func() {
		for event := range systemState.EventChannel {
			fmt.Printf("Keys: %v, KeyCode: %v, KeyDown: %v\n", event.PressedKeys, event.KeyCode, event.KeyDown)
		}
	}()

	// Prevent main from exiting
	select {}
}
