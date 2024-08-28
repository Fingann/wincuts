package keylogger 

import (
	"fmt"
	"maps"
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
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
)

const (
	WH_KEYBOARD_LL = 13
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

// KeyboardHook holds the state of the modifier keys and manages key events
type KeyboardHook struct {
	keyState    map[uint32]bool
	events      chan KeyEvent
	stopChannel chan struct{}
	hookHandle  uintptr
}

func NewKeyboardHook() *KeyboardHook {
	return &KeyboardHook{
		keyState:    make(map[uint32]bool),
		events:      make(chan KeyEvent, 10), // Buffer added to avoid blocking
		stopChannel: make(chan struct{}),
	}
}

type KeyEvent struct {
	PressedKeys map[uint32]bool
	KeyCode     uint32
	KeyDown     bool
}

// HandleKeyPress handles key press events and sends them to the event channel
func (ss *KeyboardHook) handleKeyPress(vkCode uint32, keyDown bool) {
	if keyDown {
		ss.keyState[vkCode] = true
	} else {
		delete(ss.keyState, vkCode)
	}

	select {
	case ss.events <- KeyEvent{
		PressedKeys: maps.Clone(ss.keyState),
		KeyCode:     vkCode,
		KeyDown:     keyDown,
	}:
		// Event successfully sent
	default:
		// If the channel is full, we might decide to log an error, or drop the event
		fmt.Println("Warning: Event channel is full, dropping key event")
	}
}

// lowLevelKeyboardProc is the callback method for the keyboard hook
// wParam: This parameter can be one of the following messages: WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, or WM_SYSKEYUP
// lParam: A pointer to a KBDLLHOOKSTRUCT structure.
// https://learn.microsoft.com/en-us/windows/win32/winmsg/lowlevelkeyboardproc
func (ss *KeyboardHook) lowLevelKeyboardProc(nCode int32, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		keyDown := wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN 
		ss.handleKeyPress(kbdstruct.VKCode, keyDown)
	}
	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

// Subscribe allows external clients to receive key events from the hook
func (ss *KeyboardHook) Subscribe() <-chan KeyEvent {
		return ss.events
}

// Start sets up the global keyboard hook and processes key events
func (ss *KeyboardHook) Start() error {
	hInstance, _, err := procGetModuleHandle.Call(0)
	if hInstance == 0 {
		return fmt.Errorf("GetModuleHandle failed: %v", err)
	}

	r1, _, err := procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		windows.NewCallback(ss.lowLevelKeyboardProc),
		hInstance,
		0,
	)

	if r1 == 0 {
		return fmt.Errorf("SetWindowsHookEx failed: %v", err)
	}
	ss.hookHandle = r1

	fmt.Println("Global hook set, press Shift, Ctrl, or Alt to test")

	go func() {
		var msg MSG
		for {
			select {
			case <-ss.stopChannel:
				return
			default:
				ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
				if ret == 0 {
					break
				}
			}
		}
	}()

	return nil
}

// Stop removes the keyboard hook and stops processing key events
func (ss *KeyboardHook) Stop() {
	close(ss.stopChannel)
	if ss.hookHandle != 0 {
		procUnhookWindowsHookEx.Call(ss.hookHandle)
		ss.hookHandle = 0
	}
	fmt.Println("Keyboard hook stopped")
}