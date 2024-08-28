package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	gohook
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	procGetModuleHandle     = kernel32.NewProc("GetModuleHandleW")
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	VK_ESCAPE      = 0x1B
)

type KBDLLHOOKSTRUCT struct {
	VKCode     uint32
	ScanCode   uint32
	Flags      uint32
	Time       uint32
	ExtraInfo  uintptr
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

var hookID windows.Handle

func LowLevelKeyboardProc(nCode int32, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 && wParam == WM_KEYDOWN {
		kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		if kbdstruct.VKCode == VK_ESCAPE {
			fmt.Println("Escape key was pressed, exiting...")
			procUnhookWindowsHookEx.Call(uintptr(hookID))
		} else {
			fmt.Printf("Key pressed: %d\n", kbdstruct.VKCode)
		}
	}

	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

func main() {
	hInstance, _, err := procGetModuleHandle.Call(0)
	if hInstance == 0 {
		fmt.Printf("GetModuleHandle failed: %v\n", err)
		return
	}

	r1, _, err := procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		windows.NewCallback(LowLevelKeyboardProc),
		hInstance,
		0,
	)

	if r1 == 0 {
		fmt.Printf("SetWindowsHookEx failed: %v\n", err)
		return
	}

	hookID = windows.Handle(r1)

	fmt.Println("Global hook set, press ESC to exit")

	var msg MSG
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
	}
}
