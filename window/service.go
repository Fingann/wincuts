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
	WS_DISABLED      = 0x08000000

	WS_EX_NOACTIVATE = 0x08000000
	WS_EX_LAYERED    = 0x00080000
	// ShowWindow commands
	SW_HIDE            = 0
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8

	// SetWindowPos flags
	SWP_NOMOVE         = 0x0002
	SWP_NOSIZE         = 0x0001
	SWP_NOZORDER       = 0x0004
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOACTIVATE     = 0x0010
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400

	// GetWindowLongPtr indexes using two's complement
	GWL_STYLE_PTR   = ^uintptr(15) // -16 in two's complement
	GWL_EXSTYLE_PTR = ^uintptr(19) // -20 in two's complement

	// Window property names
	PROP_DESKTOP_NUMBER   = "WinCuts_DesktopNumber"
	PROP_ORIGINAL_STYLE   = "WinCuts_OriginalStyle"
	PROP_ORIGINAL_EXSTYLE = "WinCuts_OriginalExStyle"

	// Shell window classes
	SHELL_TRAY_WND = "Shell_TrayWnd"
	SHELL_DEFVIEW  = "SHELLDLL_DefView"

	// GetWindow constants
	GW_OWNER = 4

	// Property prefixes to identify our properties
	PROP_PREFIX = "WinCuts_"
	PROP_MARKER = "WinCuts_Marker"
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
	IsHidden    bool
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
	setProp                *windows.LazyProc
	getProp                *windows.LazyProc
	removeProp             *windows.LazyProc
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
		setProp:                user32.NewProc("SetPropW"),
		getProp:                user32.NewProc("GetPropW"),
		removeProp:             user32.NewProc("RemovePropW"),
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

// getClassName gets the class name of a window
func (s *Service) getClassName(hwnd syscall.Handle) string {
	var className [256]uint16
	ret, _, _ := s.user32.NewProc("GetClassNameW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&className[0])),
		uintptr(len(className)),
	)
	if ret == 0 {
		return ""
	}
	return windows.UTF16ToString(className[:])
}

// getWindow gets a window handle by command
func (s *Service) getWindow(hwnd syscall.Handle, cmd uint32) syscall.Handle {
	ret, _, _ := s.user32.NewProc("GetWindow").Call(
		uintptr(hwnd),
		uintptr(cmd),
	)
	return syscall.Handle(ret)
}




func (s *Service) GetAllWindows() ([]WindowInfo, error) {
	var windows []WindowInfo
	var totalWindows int

	callback := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		totalWindows++

		// Get window title
		title, _ := s.GetWindowTitle(hwnd)
		if title == "" {
			return 1 // Skip windows with no title
		}
		if strings.Contains(title, "Excel") {
			fmt.Printf("  Found window: %s (hwnd: %x)\n", title, hwnd)
		}

		// Check if this is one of our hidden windows first
		if s.isOurHiddenWindow(hwnd) {
			storedDesktop, err := s.getWindowProp(hwnd, PROP_DESKTOP_NUMBER)
			if err == nil  {
				// Get window states for our hidden window
				style, _, _ := s.getWinLongPtr.Call(uintptr(hwnd), GWL_STYLE_PTR)
				isMinimized := (style & WS_MINIMIZE) != 0
				isMaximized := (style & WS_MAXIMIZE) != 0

				windows = append(windows, WindowInfo{
					Handle:      hwnd,
					Title:       title,
					DesktopNum:  int(storedDesktop) + 1,
					IsVisible:   false,
					IsMinimized: isMinimized,
					IsMaximized: isMaximized,
					IsHidden:    true,
				})
				fmt.Printf("  Found hidden window: %s (hwnd: %x, desktop: %d, hidden: %v)\n", title, hwnd, int(storedDesktop)+1, true)
			}
			return 1
		}
	

		// Check desktop number for visible windows
		windowDesktop := s.GetWindowDesktopNumber(hwnd)
		// If the window is not on any desktop, skip it
		if windowDesktop == 4294967295 {
			return 1
		}

		// Get window states
		style, _, _ := s.getWinLongPtr.Call(uintptr(hwnd), GWL_STYLE_PTR)
		isVisible := s.IsWindowVisible(hwnd)
		isMinimized := (style & WS_MINIMIZE) != 0
		isMaximized := (style & WS_MAXIMIZE) != 0

		windows = append(windows, WindowInfo{
			Handle:      hwnd,
			Title:       title,
			DesktopNum:  windowDesktop + 1,
			IsVisible:   isVisible,
			IsMinimized: isMinimized,
			IsMaximized: isMaximized,
			IsHidden:    false,
		})
		fmt.Printf("  Added window: %s (hwnd: %x, desktop: %d, hidden: %v)\n", title, hwnd, windowDesktop+1, false)
		return 1
	})

	s.enumWindows.Call(callback, 0)
	return windows, nil
}

// GetWindowsOnDesktop returns all windows on the specified virtual desktop
// desktopNum is 1-based (as shown in Windows UI)
func (s *Service) GetWindowsOnDesktop(desktopNum int) ([]WindowInfo, error) {
	windows, err := s.GetAllWindows()
	if err != nil {
		return nil, err
	}
	wantedWindows:= make([]WindowInfo, 0, len(windows))
	for _, window := range windows {
		if window.DesktopNum == desktopNum {
			wantedWindows = append(wantedWindows, window)
		}
	}

	return wantedWindows, nil
}

// isOurHiddenWindow checks if this is a window that we've hidden
func (s *Service) isOurHiddenWindow(hwnd syscall.Handle) bool {
	// Check for our marker
	markerPtr, err := windows.UTF16PtrFromString(PROP_MARKER)
	if err != nil {
		return false
	}

	marker, _, err := s.getProp.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(markerPtr)),
	)
	if err != nil  && err != windows.NOERROR {
		return false
	}

	// Verify it's our process ID
	if marker == 0  {
		return false
	}

	// Check if the desktop property exists
	_, err = s.getWindowProp(hwnd, PROP_DESKTOP_NUMBER)
	return err == nil
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

// getWindowProp gets a window property
func (s *Service) getWindowProp(hwnd syscall.Handle, name string) (uintptr, error) {
	namePtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0, fmt.Errorf("failed to convert property name: %w", err)
	}

	// First check if this window has our marker
	markerPtr, err := windows.UTF16PtrFromString(PROP_MARKER)
	if err != nil {
		return 0, fmt.Errorf("failed to convert marker name: %w", err)
	}

	marker, _, err := s.getProp.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(markerPtr)),
	)
	if err != nil  && err != windows.NOERROR {
		return 0, fmt.Errorf("failed to get property: %w", err)
	}

	// If no marker or wrong process, this isn't our window
	if marker == 0 {
		return 0, fmt.Errorf("not our window")
	}

	// Get the actual property
	ret, _, err := s.getProp.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(namePtr)),
	)
	if err != nil  && err != windows.NOERROR {
		return 0, fmt.Errorf("failed to get property: %w", err)
	}

	if ret == 0 {
		return 0, fmt.Errorf("property not found")
	}

	return ret, nil
}

// setWindowProp sets a window property
func (s *Service) setWindowProp(hwnd syscall.Handle, name string, value uintptr) error {
	// First set our marker with our process ID
	markerPtr, err := windows.UTF16PtrFromString(PROP_MARKER)
	if err != nil {
		return fmt.Errorf("failed to convert marker name: %w", err)
	}

	pid := windows.GetCurrentProcessId()
	ret, _, err := s.setProp.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(markerPtr)),
		uintptr(pid),
	)
	if ret == 0 {
		return fmt.Errorf("failed to set marker: %w", err)
	}

	// Then set the actual property
	namePtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return fmt.Errorf("failed to convert property name: %w", err)
	}

	ret, _, err = s.setProp.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(namePtr)),
		value,
	)
	if ret == 0 {
		return fmt.Errorf("failed to set window property: %w", err)
	}

	return nil
}

// removeWindowProp removes a window property and our marker
func (s *Service) removeWindowProp(hwnd syscall.Handle, name string) error {
	namePtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return fmt.Errorf("failed to convert property name: %w", err)
	}

	_, _, err = s.removeProp.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(namePtr)),
	)
	if err != nil && err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to remove window property: %w", err)
	}

	// Also remove our marker if this was the last property
	markerPtr, err := windows.UTF16PtrFromString(PROP_MARKER)
	if err != nil {
		return fmt.Errorf("failed to convert marker name: %w", err)
	}

	_, _, err = s.removeProp.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(markerPtr)),
	)
	if err != nil && err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to remove marker: %w", err)
	}

	return nil
}

func (s *Service) SetWindowVisabilityHidden(hwnd syscall.Handle) error {
	title, _ := s.GetWindowTitle(hwnd)
	windowDesktop := s.GetWindowDesktopNumber(hwnd)

	fmt.Printf("Attempting to hide window: '%s' (hwnd: %x, desktop: %d)\n", title, hwnd, windowDesktop+1)


	// Save the desktop number before hiding
	if err := s.setWindowProp(hwnd, PROP_DESKTOP_NUMBER, uintptr(windowDesktop)); err != nil {
		return fmt.Errorf("failed to save desktop number: %w", err)
	}

	// First hide the window to prevent flashing
	_, _, err := s.showWindow.Call(uintptr(hwnd), SW_HIDE)
	if err != nil && err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to hide window: %w", err)
	}


	fmt.Printf("  Window hidden\n")
	return nil
}

func (s *Service) SetWindowVisabilityVisible(hwnd syscall.Handle) error {

	title, _ := s.GetWindowTitle(hwnd)
	windowDesktop := s.GetWindowDesktopNumber(hwnd)
	fmt.Printf("Attempting to show window: '%s' (hwnd: %x, desktop: %d)\n", title, hwnd, windowDesktop+1)

	// Get the original desktop number
	origDesktop, err := s.getWindowProp(hwnd, PROP_DESKTOP_NUMBER)
	if err == nil {
		// Move window back to its original desktop before showing it
		if err := s.MoveWindowToDesktop(hwnd, int(origDesktop)); err != nil {
			return fmt.Errorf("failed to restore window to original desktop: %w", err)
		}
		fmt.Printf("  Restored to original desktop: %d\n", int(origDesktop)+1)
	}

	// First hide the window to prevent flashing while we restore styles
	_, _, err = s.showWindow.Call(uintptr(hwnd), SW_SHOW)
	if err != nil && err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to temporarily hide window: %w", err)
	}

		// Move window back to its original desktop before showing it
		if err := s.MoveWindowToDesktop(hwnd, int(origDesktop)); err != nil {
			return fmt.Errorf("failed to restore window to original desktop: %w", err)
		}
		fmt.Printf("  Restored to original desktop: %d\n", int(origDesktop)+1)
	


	fmt.Printf("  Window shown\n")

	return nil
}

// MoveWindowToDesktop moves a window to the specified desktop number (0-based)
func (s *Service) MoveWindowToDesktop(hwnd syscall.Handle, desktopNum int) error {
	ret, _, err := s.vdapi.NewProc("MoveWindowToDesktopNumber").Call(
		uintptr(hwnd),
		uintptr(desktopNum),
	)
	if ret == 0 && err != nil && err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to move window to desktop: %w", err)
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

		if err := s.SetWindowVisabilityHidden(win.Handle); err != nil {
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

		if err := s.SetWindowVisabilityVisible(win.Handle); err != nil {
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
