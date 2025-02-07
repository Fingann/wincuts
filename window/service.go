//go:build windows

package window

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// Window styles
	WS_VISIBLE = 0x10000000

	// Extended window styles
	WS_EX_APPWINDOW  = 0x00040000
	WS_EX_TOOLWINDOW = 0x00000080
	WS_MINIMIZE      = 0x00020000
	WS_MAXIMIZE      = 0x00010000

	// SetWindowPos flags
	SWP_NOMOVE       = 0x0002
	SWP_NOSIZE       = 0x0001
	SWP_NOZORDER     = 0x0004
	SWP_FRAMECHANGED = 0x0020

	// GetWindowLongPtr indexes using two's complement
	GWL_STYLE_PTR   = ^uintptr(15) // -16 in two's complement
	GWL_EXSTYLE_PTR = ^uintptr(19) // -20 in two's complement
)

// GetWindowLongPtrIndex represents the index values for GetWindowLongPtr
const (
	GWLP_WNDPROC    = -4
	GWLP_HINSTANCE  = -6
	GWLP_HWNDPARENT = -8
	GWLP_USERDATA   = -21
	GWLP_ID         = -12
)

// GetWindowLongPtr indices for 64-bit Windows
var (
	gwlStyle   = -16
	gwlExStyle = -20
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
	getWinLongPtr          *windows.LazyProc
	setWinLongPtr          *windows.LazyProc
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
		getWinLongPtr:          user32.NewProc("GetWindowLongPtrW"),
		setWinLongPtr:          user32.NewProc("SetWindowLongPtrW"),
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
		style, _, _ := s.getWinLongPtr.Call(uintptr(hwnd), GWL_STYLE_PTR)
		exStyle, _, _ := s.getWinLongPtr.Call(uintptr(hwnd), GWL_EXSTYLE_PTR)

		// Check window states
		isVisible := s.IsWindowVisible(hwnd)
		isMinimized := (style & WS_MINIMIZE) != 0
		isMaximized := (style & WS_MAXIMIZE) != 0
		isTool := (exStyle & WS_EX_TOOLWINDOW) != 0

		// Check if window is on the requested desktop
		windowDesktop := s.GetWindowDesktopNumber(hwnd)

		// Skip if not on the requested desktop (compare with 0-based index)
		if windowDesktop != zeroBasedDesktop {
			return 1
		}

		// Skip only system tool windows, not our hidden windows
		if isTool && title != "" && !strings.Contains(title, "WinCuts") {
			// Check if this is a window we hid (it will have both TOOLWINDOW and APPWINDOW)
			isApp := (exStyle & WS_EX_APPWINDOW) != 0
			if !isApp {
				return 1 // Skip only if it's a real tool window
			}
		}

		windows = append(windows, WindowInfo{
			Handle:      hwnd,
			Title:       title,
			DesktopNum:  windowDesktop + 1, // Store 1-based desktop number
			IsVisible:   isVisible,
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
	title, _ := s.GetWindowTitle(hwnd)
	windowDesktop := s.GetWindowDesktopNumber(hwnd)

	if hide {
		fmt.Printf("Attempting to hide taskbar icon: '%s' (hwnd: %x, desktop: %d)\n", title, hwnd, windowDesktop+1)

		// Get current extended style
		ret, _, err := s.getWinLongPtr.Call(uintptr(hwnd), GWL_EXSTYLE_PTR)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to get window extended style: %w", err)
		}
		exStyle := uint32(ret)

		// Get current style
		ret, _, err = s.getWinLongPtr.Call(uintptr(hwnd), GWL_STYLE_PTR)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to get window style: %w", err)
		}
		style := uint32(ret)

		fmt.Printf("  Current styles - exStyle: %x, style: %x\n", exStyle, style)

		// Add WS_EX_TOOLWINDOW and remove WS_EX_APPWINDOW to hide from taskbar
		newExStyle := (exStyle | WS_EX_TOOLWINDOW) &^ WS_EX_APPWINDOW
		// Keep the window visible
		newStyle := style | WS_VISIBLE

		fmt.Printf("  New styles - exStyle: %x, style: %x\n", newExStyle, style)

		// Apply the new extended style
		ret, _, err = s.setWinLongPtr.Call(
			uintptr(hwnd),
			GWL_EXSTYLE_PTR,
			uintptr(newExStyle),
		)
		// SetWindowLongPtr returns the previous value, 0 might be valid
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to set window extended style: %w", err)
		}

		// Apply the new style
		ret, _, err = s.setWinLongPtr.Call(
			uintptr(hwnd),
			GWL_STYLE_PTR,
			uintptr(newStyle),
		)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to set window style: %w", err)
		}

		// Update the window to reflect the style changes
		ret, _, err = s.setWinPos.Call(
			uintptr(hwnd),
			0,
			0, 0, 0, 0,
			uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED),
		)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to update window position: %w", err)
		}

		// Check desktop number after hiding
		newDesktop := s.GetWindowDesktopNumber(hwnd)
		fmt.Printf("  Desktop number before: %d, after: %d\n", windowDesktop+1, newDesktop+1)

		fmt.Printf("  Taskbar icon hidden\n")
	} else {
		fmt.Printf("Attempting to show taskbar icon: '%s' (hwnd: %x, desktop: %d)\n", title, hwnd, windowDesktop+1)

		// Get current extended style
		ret, _, err := s.getWinLongPtr.Call(uintptr(hwnd), GWL_EXSTYLE_PTR)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to get window extended style: %w", err)
		}
		exStyle := uint32(ret)

		// Get current style
		ret, _, err = s.getWinLongPtr.Call(uintptr(hwnd), GWL_STYLE_PTR)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to get window style: %w", err)
		}
		style := uint32(ret)

		fmt.Printf("  Current styles - exStyle: %x, style: %x\n", exStyle, style)

		// Remove WS_EX_TOOLWINDOW and add WS_EX_APPWINDOW to show in taskbar
		newExStyle := (exStyle &^ WS_EX_TOOLWINDOW) | WS_EX_APPWINDOW
		// Keep the window visible
		newStyle := style | WS_VISIBLE

		fmt.Printf("  New styles - exStyle: %x, style: %x\n", newExStyle, style)

		// Apply the new extended style
		ret, _, err = s.setWinLongPtr.Call(
			uintptr(hwnd),
			GWL_EXSTYLE_PTR,
			uintptr(newExStyle),
		)
		// SetWindowLongPtr returns the previous value, 0 might be valid
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to set window extended style: %w", err)
		}

		// Apply the new style
		ret, _, err = s.setWinLongPtr.Call(
			uintptr(hwnd),
			GWL_STYLE_PTR,
			uintptr(newStyle),
		)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to set window style: %w", err)
		}

		// Update the window to reflect the style changes
		ret, _, err = s.setWinPos.Call(
			uintptr(hwnd),
			0,
			0, 0, 0, 0,
			uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED),
		)
		if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
			return fmt.Errorf("failed to update window position: %w", err)
		}

		// Check desktop number after showing
		newDesktop := s.GetWindowDesktopNumber(hwnd)
		fmt.Printf("  Desktop number before: %d, after: %d\n", windowDesktop+1, newDesktop+1)

		fmt.Printf("  Taskbar icon shown\n")
	}

	return nil
}

// HideWindowsOnDesktop hides all windows on the specified desktop number
// desktopNum is 1-based (as shown in Windows UI)
func (s *Service) HideWindowsOnDesktop(desktopNum int) error {
	// First get all windows on the desktop
	windows, err := s.GetWindowsOnDesktop(desktopNum)
	if err != nil {
		return fmt.Errorf("failed to get windows on desktop %d: %w", desktopNum, err)
	}

	fmt.Printf("\nFound %d windows on desktop %d:\n", len(windows), desktopNum)

	// Hide each window
	var errors []error
	for _, win := range windows {
		state := ""
		if win.IsMinimized {
			state = " (Minimized)"
		} else if win.IsMaximized {
			state = " (Maximized)"
		}
		visibility := ""
		if !win.IsVisible {
			visibility = " [Hidden]"
		}
		fmt.Printf("- Hiding: %s%s%s\n", win.Title, state, visibility)

		if err := s.SetWindowVisibility(win.Handle, true); err != nil {
			fmt.Printf("  Error hiding window: %v\n", err)
			errors = append(errors, fmt.Errorf("failed to hide window '%s': %w", win.Title, err))
		} else {
			fmt.Printf("  Successfully hidden\n")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered errors while hiding windows: %v", errors)
	}

	fmt.Printf("\nSuccessfully hid %d windows on desktop %d\n", len(windows), desktopNum)
	return nil
}

// ShowWindowsOnDesktop shows all windows on the specified desktop number
// desktopNum is 1-based (as shown in Windows UI)
func (s *Service) ShowWindowsOnDesktop(desktopNum int) error {
	// First get all windows on the desktop
	windows, err := s.GetWindowsOnDesktop(desktopNum)
	if err != nil {
		return fmt.Errorf("failed to get windows on desktop %d: %w", desktopNum, err)
	}

	fmt.Printf("\nFound %d windows on desktop %d:\n", len(windows), desktopNum)

	// Show each window
	var errors []error
	for _, win := range windows {
		state := ""
		if win.IsMinimized {
			state = " (Minimized)"
		} else if win.IsMaximized {
			state = " (Maximized)"
		}
		visibility := ""
		if !win.IsVisible {
			visibility = " [Hidden]"
		}
		fmt.Printf("- Showing: %s%s%s\n", win.Title, state, visibility)

		if err := s.SetWindowVisibility(win.Handle, false); err != nil {
			fmt.Printf("  Error showing window: %v\n", err)
			errors = append(errors, fmt.Errorf("failed to show window '%s': %w", win.Title, err))
		} else {
			fmt.Printf("  Successfully shown\n")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered errors while showing windows: %v", errors)
	}

	fmt.Printf("\nSuccessfully showed %d windows on desktop %d\n", len(windows), desktopNum)
	return nil
}
