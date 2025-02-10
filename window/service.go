//go:build windows

package window

import (
	"errors"
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

	ERR_WINDOW_NOT_ON_ANY_DESKTOP = errors.New("window is not on any desktop")
)

// WindowInfo contains information about a window
type WindowInfo struct {
	Handle      syscall.Handle
	Title       string
	DesktopNum  int
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
	propService            *PropService
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

	propService := NewPropService()

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
		propService:            propService,
	}, nil
}

// GetWindowDesktopNumber gets the desktop number for a window
func (s *Service) GetWindowDesktopNumber(hwnd syscall.Handle) (int, error) {
	ret, _, err := s.getWindowDesktopNumber.Call(uintptr(hwnd))
	if err != nil && err != windows.ERROR_SUCCESS {
		return 0, fmt.Errorf("failed to get window desktop number: %w", err)
	}
	if ret == 4294967295 {
		return 0, ERR_WINDOW_NOT_ON_ANY_DESKTOP
	}
	return int(ret), nil
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

func (s *Service) GetAllWindows() ([]WindowInfo, error) {
	var windows []WindowInfo
	var totalWindows int

	callback := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		totalWindows++

		// Get window title
		title, err := s.GetWindowTitle(hwnd)
		if err != nil {
			fmt.Printf("  Error getting window title, window: %x, err: %v\n", hwnd, err.Error())
			return 1
		}
		if strings.Contains(title, "Teams") {
			fmt.Printf("  Found window: %s (hwnd: %x)\n", title, hwnd)
		}
		isHidden := false
		// check to see if the window is on any desktop
		desktopNumber, err := s.GetWindowDesktopNumber(hwnd)
		if err != nil {
			if err != ERR_WINDOW_NOT_ON_ANY_DESKTOP {
				return 1
			}
				// get current desktop number from prop service
				desktopNumber, err = s.propService.GetDesktopNumber(hwnd)
				if err != nil {
					if err == ERR_NO_DATA {
						// window is not on any desktop and prop service has no data
						// This window is not managed by WinCuts
						return 1
					}
					// error while getting desktop number from prop service
					return 1
				}
				isHidden = true
		}
		// Get window states for our hidden window

		windows = append(windows, WindowInfo{
			Handle:     hwnd,
			Title:      title,
			DesktopNum: desktopNumber + 1,
			IsHidden:   isHidden,
		})
		fmt.Printf("  Found hidden window: %s (hwnd: %x, desktop: %d, hidden: %v)\n", title, hwnd, desktopNumber+1, isHidden)

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
	wantedWindows := make([]WindowInfo, 0, len(windows))
	for _, window := range windows {
		if window.DesktopNum == desktopNum {
			wantedWindows = append(wantedWindows, window)
		}
	}

	return wantedWindows, nil
}

func (s *Service) SetWindowVisabilityHidden(hwnd syscall.Handle) error {
	title, _ := s.GetWindowTitle(hwnd)
	fmt.Printf("  Setting window visibility to hidden: %s (hwnd: %x)\n", title, hwnd)
	windowDesktop, err := s.GetWindowDesktopNumber(hwnd)
	if err != nil {
		if err == ERR_WINDOW_NOT_ON_ANY_DESKTOP {
			// window is not on any desktop, so we don't need to hide it
			return nil
		}
		return fmt.Errorf("failed to get window desktop number: %w", err)
	}

	// Save the desktop number before hiding
	if err := s.propService.SetDesktopNumber(hwnd, windowDesktop); err != nil {
		return fmt.Errorf("failed to save desktop number: %w", err)
	}
	// hide the window
	_, _, err = s.showWindow.Call(uintptr(hwnd), SW_HIDE)
	if err != nil && err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to hide window: %w", err)
	}

	fmt.Sprintf("  Window hidden: %s (hwnd: %x, desktop: %d)\n", title, hwnd, windowDesktop+1)
	return nil
}

func (s *Service) SetWindowVisabilityVisible(hwnd syscall.Handle) error {

	title, _ := s.GetWindowTitle(hwnd)
	// Get the original desktop number
	origDesktop, err := s.propService.GetDesktopNumber(hwnd)
	if err != nil {
		return fmt.Errorf("failed to get original desktop number: %w", err)
	}

	// show the window
	_, _, err = s.showWindow.Call(uintptr(hwnd), SW_SHOW)
	if err != nil && err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to temporarily hide window: %w", err)
	}

	// Move window back to its original desktop before showing it
	if err := s.MoveWindowToDesktop(hwnd, int(origDesktop)); err != nil {
		return fmt.Errorf("failed to restore window to original desktop: %w", err)
	}
	fmt.Sprintf("  Window shown: %s (hwnd: %x, desktop: %d)\n", title, hwnd, origDesktop+1)

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
		// skip windows that are already hidden
		if win.IsHidden {
			continue
		}
		if err := s.SetWindowVisabilityHidden(win.Handle); err != nil {
			fmt.Printf("  Error hiding window: %v\n", err)
			errors = append(errors, fmt.Errorf("failed to hide window '%s': %w", win.Title, err))
		}

		fmt.Printf("  Successfully hidden\n")
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
		if !win.IsHidden {
			// window is not hidden, so we don't need to show it
			continue
		}
		if err := s.SetWindowVisabilityVisible(win.Handle); err != nil {
			errors = append(errors, fmt.Errorf("failed to show window '%s': %w", win.Title, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered errors while showing windows: %v", errors)
	}

	fmt.Printf("\nSuccessfully showed %d windows on desktop %d\n", len(windows), desktopNum)
	return nil
}
