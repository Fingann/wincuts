//go:build windows

package window

import (
	"errors"

	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Property types for window properties
const (
	PROP_DATA   = "WinCuts_Data"

)
var (
	ERR_NO_DATA = errors.New("No data found")
	ERR_INVALID_DATA = errors.New("Invalid data")

)



// PropService handles window property management using TLV encoding

type PropService struct {
	setProp    *windows.LazyProc
	getProp    *windows.LazyProc
	removeProp *windows.LazyProc
}

// NewPropService creates a new property service
func NewPropService() *PropService {
	user32 := windows.NewLazyDLL("user32.dll")
	return &PropService{
		setProp:    user32.NewProc("SetPropW"),
		getProp:    user32.NewProc("GetPropW"),
		removeProp: user32.NewProc("RemovePropW"),
	}
}
// SetIntegerProperty sets an integer as the property directly, no memory allocation.
func (p *PropService) setIntegerProperty(hwnd syscall.Handle, propName string, value uint32) error {
    nameUTF16, err := windows.UTF16PtrFromString(propName)
    if err != nil {
        return fmt.Errorf("failed to convert property name: %w", err)
    }

    // Cast value to uintptr and call SetPropW directly
    ret, _, callErr := p.setProp.Call(
        uintptr(hwnd),
        uintptr(unsafe.Pointer(nameUTF16)),
        uintptr(value+1),
    )
    if ret == 0 {
        return fmt.Errorf("SetPropW failed: %v", callErr)
    }
    return nil
}

// GetIntegerProperty retrieves the integer from the property. 
// If GetPropW returns 0, this might mean the property doesn't exist or its value is truly 0.
func (p *PropService) getIntegerProperty(hwnd syscall.Handle, propName string) (uint32, error) {
    nameUTF16, err := windows.UTF16PtrFromString(propName)
    if err != nil {

        return 0, fmt.Errorf("failed to convert property name: %w", err)
    }

    ptr, _, callErr := p.getProp.Call(
        uintptr(hwnd),
        uintptr(unsafe.Pointer(nameUTF16)),
    )
    if ptr == 0 {
        // Could be an actual 0 value OR the property might not exist
        return 0, fmt.Errorf("property may be missing, or is zero: %v", callErr)
    }

    // Convert the uintptr to a uint32
    return uint32(ptr)-1, nil
}


// RemoveIntegerProperty removes the integer property (no memory to free).
func (p *PropService) RemoveIntegerProperty(hwnd syscall.Handle, propName string) error {
    nameUTF16, err := windows.UTF16PtrFromString(propName)
    if err != nil {
        return fmt.Errorf("failed to convert property name: %w", err)

    }
    ptr, _, callErr := p.removeProp.Call(
        uintptr(hwnd),
        uintptr(unsafe.Pointer(nameUTF16)),
    )
    if ptr == 0 {
        return fmt.Errorf("RemovePropW failed or property not found: %v", callErr)
    }
    // Because we never allocated memory with LocalAlloc, there's nothing to free here.
    return nil
}

// SetDesktopNumber sets the desktop number property for a window
func (p *PropService) SetDesktopNumber(hwnd syscall.Handle, desktop int) error {
	// Store the encoded data as a property
	err := p.setIntegerProperty(hwnd, PROP_DATA, uint32(desktop))
	if err != nil {
		return fmt.Errorf("failed to set desktop property: %w", err)
	}




	return nil
}


// GetDesktopNumber gets the desktop number property for a window
func (p *PropService) GetDesktopNumber(hwnd syscall.Handle) (int, error) {
	// Get the property data
	data, err := p.getIntegerProperty(hwnd, PROP_DATA)
	if err != nil {
		return 0, fmt.Errorf("failed to get desktop property: %w", err)

	}

	return int(data), nil
}




// RemoveProps removes all properties from a window
func (p *PropService) RemoveProps(hwnd syscall.Handle) error {
	namePtr, err := windows.UTF16PtrFromString(PROP_DATA)
	if err != nil {
		return fmt.Errorf("failed to convert property name: %w", err)
	}
	_, _, err = p.removeProp.Call(uintptr(hwnd), uintptr(unsafe.Pointer(namePtr)))
	if err != nil {
		return fmt.Errorf("failed to remove property: %w", err)
	}

	return nil
}
