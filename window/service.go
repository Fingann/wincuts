//go:build windows

package window

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	GWL_EXSTYLE      = ^(^uint32(0) >> 1) - 19 // -20 in two's complement
	GWL_STYLE        = ^(^uint32(0) >> 1) - 15 // -16 in two's complement
	WS_EX_APPWINDOW  = 0x00040000
	WS_EX_TOOLWINDOW = 0x00000080
	WS_VISIBLE       = 0x10000000
	WS_MINIMIZE      = 0x20000000
	WS_MAXIMIZE      = 0x01000000
	SWP_NOSIZE       = 0x0001
	SWP_NOMOVE       = 0x0002
	SWP_NOZORDER     = 0x0004
	SWP_FRAMECHANGED = 0x0020
	SW_HIDE          = 0
	SW_SHOW          = 5
)

// WindowInfo contains information about a window
type WindowInfo struct {
	Handle      syscall.Handle
	Title       string
	DesktopNum  int
	IsVisible   bool
	IsMinimized bool
	IsMaximized bool
}

// Service provides methods to control window visibility
type Service struct {
	user32                 *windows.LazyDLL
	vdapi                  *windows.LazyDLL
	findWindow             *windows.LazyProc
	getWinLong             *windows.LazyProc
	setWinLong             *windows.LazyProc
	setWinPos              *windows.LazyProc
	showWindow             *windows.LazyProc
	enumWindows            *windows.LazyProc
	getWindowText          *windows.LazyProc
	isWindowVisible        *windows.LazyProc
	getWindowDesktopNumber *windows.LazyProc
}

// NewService creates a new window management service
func NewService() (*Service, error) {
	user32 := windows.NewLazyDLL("user32.dll")

	// Try to find VirtualDesktopAccessor.dll in different locations
	dllPaths := []string{
		"VirtualDesktopAccessor.dll", // Current directory
		filepath.Join(os.Getenv("LOCALAPPDATA"), "WinCuts", "VirtualDesktopAccessor.dll"), // Installation directory
	}

	var vdapi *windows.LazyDLL
	var dllErr error
	for _, path := range dllPaths {
		vdapi = windows.NewLazyDLL(path)
		err := vdapi.Load()
		if err == nil {
			dllErr = nil
			break
		}
		dllErr = err
	}

	if dllErr != nil {
		return nil, fmt.Errorf("failed to load VirtualDesktopAccessor.dll: %w", dllErr)
	}

	return &Service{
		user32:                 user32,
		vdapi:                  vdapi,
		findWindow:             user32.NewProc("FindWindowW"),
		getWinLong:             user32.NewProc("GetWindowLongW"),
		setWinLong:             user32.NewProc("SetWindowLongW"),
		setWinPos:              user32.NewProc("SetWindowPos"),
		showWindow:             user32.NewProc("ShowWindow"),
		enumWindows:            user32.NewProc("EnumWindows"),
		getWindowText:          user32.NewProc("GetWindowTextW"),
		isWindowVisible:        user32.NewProc("IsWindowVisible"),
		getWindowDesktopNumber: vdapi.NewProc("GetWindowDesktopNumber"),
	}, nil
}

// GetWindowDesktopNumber gets the desktop number for a window
func (s *Service) GetWindowDesktopNumber(hwnd syscall.Handle) int {
	ret, _, _ := s.getWindowDesktopNumber.Call(uintptr(hwnd))
	return int(ret)
}

// GetWindowTitle gets the title of a window
func (s *Service) GetWindowTitle(hwnd syscall.Handle) (string, error) {
	var buffer [256]uint16
	_, _, err := s.getWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	if err != nil && err != windows.ERROR_SUCCESS {
		return "", fmt.Errorf("failed to get window text: %w", err)
	}
	return windows.UTF16ToString(buffer[:]), nil
}

// IsWindowVisible checks if a window is visible
func (s *Service) IsWindowVisible(hwnd syscall.Handle) bool {
	ret, _, _ := s.isWindowVisible.Call(uintptr(hwnd))
	return ret != 0
}

// GetWindowsOnDesktop returns all windows on the specified virtual desktop
// desktopNum is 1-based (as shown in Windows UI)
func (s *Service) GetWindowsOnDesktop(desktopNum int) ([]WindowInfo, error) {
	// Convert to 0-based for internal Windows API
	zeroBasedDesktop := desktopNum - 1
	var windows []WindowInfo
	var totalWindows int

	callback := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		totalWindows++

		// Get window title
		title, _ := s.GetWindowTitle(hwnd)
		if title == "" {
			return 1 // Skip windows with no title
		}

		// Get window styles
		style, _, _ := s.getWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_STYLE)))
		exStyle, _, _ := s.getWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_EXSTYLE)))

		// Check window states
		isVisible := s.IsWindowVisible(hwnd)
		hasVisibleStyle := (style & WS_VISIBLE) != 0
		isMinimized := (style & WS_MINIMIZE) != 0
		isMaximized := (style & WS_MAXIMIZE) != 0
		isTool := (exStyle & WS_EX_TOOLWINDOW) != 0

		// Skip tool windows
		if isTool {
			return 1
		}

		// Check if window is on the requested desktop
		windowDesktop := s.GetWindowDesktopNumber(hwnd)

		// Skip if not on the requested desktop (compare with 0-based index)
		if windowDesktop != zeroBasedDesktop {
			return 1
		}

		windows = append(windows, WindowInfo{
			Handle:      hwnd,
			Title:       title,
			DesktopNum:  windowDesktop + 1, // Store 1-based desktop number
			IsVisible:   isVisible || hasVisibleStyle,
			IsMinimized: isMinimized,
			IsMaximized: isMaximized,
		})
		return 1
	})

	s.enumWindows.Call(callback, 0)
	return windows, nil
}

// FindWindow finds a window by its title
func (s *Service) FindWindow(title string) (syscall.Handle, error) {
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0, fmt.Errorf("error converting window title: %w", err)
	}

	hwnd, _, err := s.findWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return 0, fmt.Errorf("window not found: %w", err)
	}

	return syscall.Handle(hwnd), nil
}

// SetWindowVisibility changes the window visibility state
func (s *Service) SetWindowVisibility(hwnd syscall.Handle, hide bool) error {
	if hide {
		// Hide the window first
		s.showWindow.Call(uintptr(hwnd), SW_HIDE)

		// Remove the window from taskbar by changing styles
		exStyle, _, _ := s.getWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_EXSTYLE)))
		style, _, _ := s.getWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_STYLE)))

		newExStyle := (exStyle &^ WS_EX_APPWINDOW) | WS_EX_TOOLWINDOW
		newStyle := style &^ WS_VISIBLE

		s.setWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_EXSTYLE)), newExStyle)
		s.setWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_STYLE)), newStyle)
	} else {
		// Restore window styles
		exStyle, _, _ := s.getWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_EXSTYLE)))
		style, _, _ := s.getWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_STYLE)))

		newExStyle := (exStyle &^ WS_EX_TOOLWINDOW) | WS_EX_APPWINDOW
		newStyle := style | WS_VISIBLE

		s.setWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_EXSTYLE)), newExStyle)
		s.setWinLong.Call(uintptr(hwnd), uintptr(uint32(GWL_STYLE)), newStyle)

		// Show the window
		s.showWindow.Call(uintptr(hwnd), SW_SHOW)
	}

	// Update the window to reflect the style changes
	ret, _, err := s.setWinPos.Call(
		uintptr(hwnd),
		0,
		0, 0, 0, 0,
		uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED),
	)
	if ret == 0 {
		return fmt.Errorf("failed to update window position: %w", err)
	}

	return nil
}
