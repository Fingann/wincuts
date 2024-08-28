package keyboard

import (
	"fmt"
	"maps"
	"os/signal"
    "os"	
	"slices"
	"syscall"
	"unsafe"
	"wincuts/keyboard/code"

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
	procTranslateMessage  = user32.NewProc("TranslateMessage")
    procDispatchMessage   = user32.NewProc("DispatchMessageW")
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

func NewHook() *KeyboardHook {
	return &KeyboardHook{
		keyState:    make(map[uint32]bool),
		events:      make(chan KeyEvent, 10), // Buffer added to avoid blocking
		stopChannel: make(chan struct{}),
	}
}

type KeyEvent struct {
	PressedKeys []uint32 
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
		PressedKeys: slices.Collect(maps.Keys(ss.keyState)),
		KeyCode:     vkCode,
		KeyDown:     keyDown,
	}:
		// Event successfully sent
	default:
		// If the channel is full, we might decide to log an error, or drop the event
		fmt.Println("Warning: Event channel is full, dropping key event") //TODO: implement proper logging levels
	}
}

// lowLevelKeyboardProc is the callback method for the keyboard hook
func (ss *KeyboardHook) lowLevelKeyboardProc(nCode int32, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		fmt.Println("lowLevelKeyboardProc invoked") // Debugging log
		kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		keyDown := wParam == code.WM_KEYDOWN || wParam == code.WM_SYSKEYDOWN 
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
	// Get the handle of the current module
	hInstance, _, err := procGetModuleHandle.Call(0)
	if hInstance == 0 {
		return fmt.Errorf("GetModuleHandle failed: %v", err)
	}

	// Set the hook on the keyboard events, with a callback to lowLevelKeyboardProc
	r1, r2, err := procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		windows.NewCallback(ss.lowLevelKeyboardProc),
		hInstance,
		0,
	)

	if r2 != 0 {
		return fmt.Errorf("SetWindowsHookEx failed: %v", err)
	}

	if r1 == 0 {
		return fmt.Errorf("SetWindowsHookEx failed: %v", err)
	}
	// Save the hook handle so we can remove it later
	ss.hookHandle = r1

	go ss.messageLoop()
	// Ensure that cleanup happens on program exit
	go ss.handleShutdown()

	fmt.Println("Keyboard hook started")
	return nil
}

func (ss *KeyboardHook) handleShutdown() {
	// Capture OS signals to ensure proper cleanup on exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	<-sigs
	ss.Stop()
}

func (ss *KeyboardHook) messageLoop() {
	var msg MSG
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		fmt.Println("Message received:", msg.Message) // Debugging log
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// Stop removes the keyboard hook and stops processing key events
func (ss *KeyboardHook) Stop() {
	// check if the hook is already stopped
	if ss.hookHandle == 0 {
		return
	}	
	close(ss.stopChannel)
	if ss.hookHandle != 0 {
		procUnhookWindowsHookEx.Call(ss.hookHandle)
		ss.hookHandle = 0
	}
	fmt.Println("Keyboard hook stopped")
}
jj