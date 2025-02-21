//go:build windows

// Package virtd implements APIs from Ciantic's VirtualDesktopAccessor.
// More info about VDA can be found at: https://github.com/Ciantic/VirtualDesktopAccessor
package virtd

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/chrsm/winapi"
)

var (
	// Get the executable's directory to find the DLL
	execDir = func() string {
		exe, err := os.Executable()
		if err != nil {
			return "."
		}
		return filepath.Dir(exe)
	}()

	vdapi = syscall.NewLazyDLL(filepath.Join(execDir, "VirtualDesktopAccessor.dll"))

	pGetCurrentDesktopNumber         = vdapi.NewProc("GetCurrentDesktopNumber")
	pGetDesktopCount                 = vdapi.NewProc("GetDesktopCount")
	pGetDesktopIdByNumber            = vdapi.NewProc("GetDesktopIdByNumber")
	pGetDesktopNumberById            = vdapi.NewProc("GetDesktopNumberById")
	pGetWindowDesktopId              = vdapi.NewProc("GetWindowDesktopId")
	pGetWindowDesktopNumber          = vdapi.NewProc("GetWindowDesktopNumber")
	pIsWindowOnCurrentVirtualDesktop = vdapi.NewProc("IsWindowOnCurrentVirtualDesktop")
	pMoveWindowToDesktopNumber       = vdapi.NewProc("MoveWindowToDesktopNumber")
	pGoToDesktopNumber               = vdapi.NewProc("GoToDesktopNumber")
	pSetDesktopName                  = vdapi.NewProc("SetDesktopName")
	pGetDesktopName                  = vdapi.NewProc("GetDesktopName")
	pRegisterPostMessageHook         = vdapi.NewProc("RegisterPostMessageHook")
	pUnregisterPostMessageHook       = vdapi.NewProc("UnregisterPostMessageHook")
	pIsPinnedWindow                  = vdapi.NewProc("IsPinnedWindow")
	pPinWindow                       = vdapi.NewProc("PinWindow")
	pUnPinWindow                     = vdapi.NewProc("UnPinWindow")
	pIsPinnedApp                     = vdapi.NewProc("IsPinnedApp")
	pPinApp                          = vdapi.NewProc("PinApp")
	pUnPinApp                        = vdapi.NewProc("UnPinApp")
	pIsWindowOnDesktopNumber         = vdapi.NewProc("IsWindowOnDesktopNumber")
	pCreateDesktop                   = vdapi.NewProc("CreateDesktop")
	pRemoveDesktop                   = vdapi.NewProc("RemoveDesktop")
)

func GetCurrentDesktopNumber() int {
	ret, _, _ := pGetCurrentDesktopNumber.Call()
	return int(ret)
}

func GetDesktopCount() int {
	ret, _, _ := pGetDesktopCount.Call()
	return int(ret)
}

func GetDesktopIdByNumber(i int) winapi.GUID {
	var guid winapi.GUID
	pGetDesktopIdByNumber.Call(uintptr(i), uintptr(unsafe.Pointer(&guid)))
	return guid
}

func GetDesktopNumberById(id winapi.GUID) int {
	ret, _, _ := pGetDesktopNumberById.Call(uintptr(unsafe.Pointer(&id)))
	return int(ret)
}

func GetWindowDesktopId(w winapi.HWND) winapi.GUID {
	var guid winapi.GUID
	pGetWindowDesktopId.Call(uintptr(w), uintptr(unsafe.Pointer(&guid)))
	return guid
}

func GetWindowDesktopNumber(w winapi.HWND) int {
	ret, _, _ := pGetWindowDesktopNumber.Call(uintptr(w))
	return int(ret)
}

func IsWindowOnCurrentVirtualDesktop(w winapi.HWND) bool {
	ret, _, _ := pIsWindowOnCurrentVirtualDesktop.Call(uintptr(w))
	return ret == 1
}

func MoveWindowToDesktopNumber(w winapi.HWND, i int) bool {
	ret, _, _ := pMoveWindowToDesktopNumber.Call(uintptr(w), uintptr(i))
	return ret == 1
}

func GoToDesktopNumber(i int) {
	pGoToDesktopNumber.Call(uintptr(i))
}

func SetDesktopName(desktopNumber int, name string) bool {
	namePtr := uintptr(unsafe.Pointer(syscall.StringBytePtr(name)))
	ret, _, _ := pSetDesktopName.Call(uintptr(desktopNumber), namePtr)
	return ret == 0
}

func GetDesktopName(desktopNumber int) string {
	var buffer [256]uint16
	pGetDesktopName.Call(uintptr(desktopNumber), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	return syscall.UTF16ToString(buffer[:])
}

func RegisterPostMessageHook(l winapi.HWND, offset int) {
	pRegisterPostMessageHook.Call(uintptr(l), uintptr(offset))
}

func UnregisterPostMessageHook(l winapi.HWND) {
	pUnregisterPostMessageHook.Call(uintptr(l))
}

func IsPinnedWindow(w winapi.HWND) bool {
	ret, _, _ := pIsPinnedWindow.Call(uintptr(w))
	return ret == 1
}

func PinWindow(w winapi.HWND) {
	pPinWindow.Call(uintptr(w))
}

func UnpinWindow(w winapi.HWND) {
	pUnPinWindow.Call(uintptr(w))
}

func IsPinnedApp(w winapi.HWND) bool {
	ret, _, _ := pIsPinnedApp.Call(uintptr(w))
	return ret == 1
}

func PinApp(w winapi.HWND) {
	pPinApp.Call(uintptr(w))
}

func UnpinApp(w winapi.HWND) {
	pUnPinApp.Call(uintptr(w))
}

func IsWindowOnDesktopNumber(w winapi.HWND, i int) bool {
	ret, _, _ := pIsWindowOnDesktopNumber.Call(uintptr(w), uintptr(i))
	return ret == 1
}

func CreateDesktop() bool {
	ret, _, _ := pCreateDesktop.Call()
	return ret == 0
}

func RemoveDesktop(removeDesktopNumber, fallbackDesktopNumber int) bool {
	ret, _, _ := pRemoveDesktop.Call(uintptr(removeDesktopNumber), uintptr(fallbackDesktopNumber))
	return ret == 0
}
