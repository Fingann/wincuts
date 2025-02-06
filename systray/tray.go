// Package systray handles the system tray icon and its interactions.
// It provides functionality to display and update the current virtual desktop number
// in the Windows system tray.
package systray

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log/slog"
	"sync"
	"syscall"
	"unsafe"
	"wincuts/config"

	"github.com/lxn/win"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Icon represents a system tray icon instance
type Icon struct {
	hwnd        win.HWND
	nid         *win.NOTIFYICONDATA
	currentText string
	iconCache   map[int]win.HICON // Cache for rendered icons
	mu          sync.Mutex
	config      config.TrayIconConfig
}

const (
	wmApp = win.WM_APP + 1
	// Windows message constants
	wmSysCommand   = 0x0112
	wmDestroy      = 0x0002
	wmClose        = 0x0010
	wmUser         = 0x0400
	wmTrayCallback = wmUser + 1
	// Icon constants
	niifNone   = 0x00000000
	nifMessage = win.NIF_MESSAGE
	nifIcon    = win.NIF_ICON
	nifTip     = win.NIF_TIP
	// System commands
	scMinimize = 0xF020
	// Icon size
	iconSize     = 22  // Slightly larger for better visibility
	cornerRadius = 4   // Rounded corners
	padding      = 2   // Space between edge and content
	bgOpacity    = 230 // Slight transparency for modern look
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
)

// drawRoundedRect draws a rounded rectangle with anti-aliased edges
func drawRoundedRect(img *image.RGBA, col color.Color, radius int) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Create a temporary mask image with 8x supersampling
	scale := 8
	maskW, maskH := w*scale, h*scale
	mask := image.NewAlpha(image.Rect(0, 0, maskW, maskH))
	radius *= scale

	// Draw the supersampled rounded rectangle on the mask
	for y := 0; y < maskH; y++ {
		for x := 0; x < maskW; x++ {
			// Convert to original coordinates
			ox, oy := float64(x)/float64(scale), float64(y)/float64(scale)

			// Calculate distance from nearest corner
			var d float64
			switch {
			// Top-left corner
			case ox < float64(radius/scale) && oy < float64(radius/scale):
				dx := float64(radius/scale) - ox
				dy := float64(radius/scale) - oy
				d = dx*dx + dy*dy
			// Top-right corner
			case ox >= float64(w-radius/scale) && oy < float64(radius/scale):
				dx := ox - float64(w-radius/scale)
				dy := float64(radius/scale) - oy
				d = dx*dx + dy*dy
			// Bottom-left corner
			case ox < float64(radius/scale) && oy >= float64(h-radius/scale):
				dx := float64(radius/scale) - ox
				dy := oy - float64(h-radius/scale)
				d = dx*dx + dy*dy
			// Bottom-right corner
			case ox >= float64(w-radius/scale) && oy >= float64(h-radius/scale):
				dx := ox - float64(w-radius/scale)
				dy := oy - float64(h-radius/scale)
				d = dx*dx + dy*dy
			default:
				mask.Set(x, y, color.Alpha{A: 255})
				continue
			}

			// Calculate anti-aliased alpha for corners
			r := float64(radius / scale)
			if d > r*r {
				mask.Set(x, y, color.Alpha{A: 0})
			} else if d > (r-1)*(r-1) {
				alpha := uint8(255 * (1 - (d-(r-1)*(r-1))/(2*r-1)))
				mask.Set(x, y, color.Alpha{A: alpha})
			} else {
				mask.Set(x, y, color.Alpha{A: 255})
			}
		}
	}

	// Scale down the mask to the original size with anti-aliasing
	finalMask := image.NewAlpha(bounds)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Average the supersampled pixels
			var sum uint32
			for sy := 0; sy < scale; sy++ {
				for sx := 0; sx < scale; sx++ {
					sum += uint32(mask.AlphaAt(x*scale+sx, y*scale+sy).A)
				}
			}
			avg := uint8(sum / uint32(scale*scale))
			finalMask.Set(x, y, color.Alpha{A: avg})
		}
	}

	// Draw the colored rectangle using the anti-aliased mask
	draw.DrawMask(img, bounds, image.NewUniform(col), image.Point{}, finalMask, image.Point{}, draw.Over)
}

// createIconWithNumber creates an HICON with the desktop number drawn on it
func (i *Icon) createIconWithNumber(number int) (win.HICON, error) {
	cfg := i.config
	// Create a new RGBA image with slightly larger size for better text rendering
	size := cfg.Size
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Fill background with transparent
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// Draw solid background
	bgColor := cfg.BgColor
	bgColor.A = cfg.BgOpacity
	draw.Draw(img, img.Bounds(), image.NewUniform(bgColor), image.Point{}, draw.Over)

	// Draw the number - draw it multiple times with slight offsets to create a bold effect
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(cfg.TextColor),
		Face: basicfont.Face7x13,
	}

	// Center the number
	numStr := fmt.Sprintf("%d", number)
	textWidth := d.MeasureString(numStr).Round()
	x := (size - textWidth) / 2
	y := (size + 9) / 2 // Adjusted for better vertical centering

	// Draw the number multiple times with slight offsets to create bold effect
	offsets := []struct{ dx, dy int }{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
	}

	for _, offset := range offsets {
		d.Dot = fixed.Point26_6{
			X: fixed.I(x + offset.dx),
			Y: fixed.I(y + offset.dy),
		}
		d.DrawString(numStr)
	}

	// Convert to HICON
	return createHICONFromImage(img), nil
}

// createHICONFromImage converts an RGBA image to HICON
func createHICONFromImage(img *image.RGBA) win.HICON {
	// Create bitmap info header
	bi := win.BITMAPINFOHEADER{
		BiSize:          uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})),
		BiWidth:         int32(img.Bounds().Dx()),
		BiHeight:        int32(-img.Bounds().Dy()), // Negative height for top-down image
		BiPlanes:        1,
		BiBitCount:      32,
		BiCompression:   win.BI_RGB,
		BiSizeImage:     0,
		BiXPelsPerMeter: 0,
		BiYPelsPerMeter: 0,
		BiClrUsed:       0,
		BiClrImportant:  0,
	}

	// Create device context
	hdc := win.GetDC(0)
	defer win.ReleaseDC(0, hdc)

	// Create DIB section
	var bits unsafe.Pointer
	hBitmap := win.CreateDIBSection(hdc, &bi, win.DIB_RGB_COLORS, &bits, 0, 0)
	if hBitmap == 0 {
		return 0
	}

	// Copy image data
	srcBytes := img.Pix
	dstBytes := (*[1 << 30]byte)(bits)
	for i := 0; i < len(srcBytes); i += 4 {
		dstBytes[i] = srcBytes[i+2]   // Blue
		dstBytes[i+1] = srcBytes[i+1] // Green
		dstBytes[i+2] = srcBytes[i]   // Red
		dstBytes[i+3] = srcBytes[i+3] // Alpha
	}

	// Create icon mask (1-bit bitmap)
	hMonoBitmap := win.CreateBitmap(int32(img.Bounds().Dx()), int32(img.Bounds().Dy()), 1, 1, nil)

	// Create icon info
	iconInfo := win.ICONINFO{
		FIcon:    win.TRUE,
		HbmMask:  hMonoBitmap,
		HbmColor: hBitmap,
	}

	// Create icon
	hIcon := win.CreateIconIndirect(&iconInfo)

	// Clean up
	win.DeleteObject(win.HGDIOBJ(hBitmap))
	win.DeleteObject(win.HGDIOBJ(hMonoBitmap))

	return hIcon
}

// windowProc handles window messages
func windowProc(hwnd win.HWND, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case win.WM_DESTROY:
		win.PostQuitMessage(0)
		return 0
	case wmTrayCallback:
		return 0
	default:
		return win.DefWindowProc(hwnd, msg, wparam, lparam)
	}
}

// New creates a new system tray icon
func New(cfg config.TrayIconConfig) (*Icon, error) {
	// Register window class
	className := syscall.StringToUTF16Ptr("WinCutsSystemTray")
	wc := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		LpfnWndProc:   syscall.NewCallback(windowProc),
		HInstance:     win.GetModuleHandle(nil),
		LpszClassName: className,
	}

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		return nil, fmt.Errorf("failed to register window class")
	}

	// Create window
	hwnd := win.CreateWindowEx(
		0,
		className,
		syscall.StringToUTF16Ptr("WinCuts"),
		0,
		0, 0, 0, 0,
		0,
		0,
		win.GetModuleHandle(nil),
		nil)

	if hwnd == 0 {
		return nil, fmt.Errorf("failed to create window")
	}

	icon := &Icon{
		hwnd:      hwnd,
		iconCache: make(map[int]win.HICON),
		config:    cfg,
	}

	// Create initial icon
	hIcon, err := icon.createIconWithNumber(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create icon: %w", err)
	}

	// Initialize NOTIFYICONDATA
	nid := &win.NOTIFYICONDATA{
		HWnd:             hwnd,
		UFlags:           nifMessage | nifIcon | nifTip,
		UCallbackMessage: wmTrayCallback,
		HIcon:            hIcon,
	}
	nid.CbSize = uint32(unsafe.Sizeof(*nid))

	// Add the icon
	if !win.Shell_NotifyIcon(win.NIM_ADD, nid) {
		return nil, fmt.Errorf("failed to add system tray icon")
	}

	icon.nid = nid
	icon.iconCache[1] = hIcon // Cache the initial icon

	return icon, nil
}

// UpdateText updates the system tray icon tooltip and icon with the current desktop number
func (i *Icon) UpdateText(desktopNum int) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	text := fmt.Sprintf("Desktop %d", desktopNum)
	if text == i.currentText {
		return nil
	}

	// Check if icon is already cached
	hIcon, exists := i.iconCache[desktopNum]
	if !exists {
		// Create and cache new icon
		var err error
		hIcon, err = i.createIconWithNumber(desktopNum)
		if err != nil {
			return fmt.Errorf("failed to create icon: %w", err)
		}
		i.iconCache[desktopNum] = hIcon
		slog.Debug("created and cached new icon", "desktop", desktopNum)
	}

	// Update icon and tooltip
	i.nid.HIcon = hIcon
	copy(i.nid.SzTip[:], syscall.StringToUTF16(text))

	if !win.Shell_NotifyIcon(win.NIM_MODIFY, i.nid) {
		return fmt.Errorf("failed to update system tray icon")
	}

	i.currentText = text
	slog.Debug("updated system tray", "desktop", desktopNum)
	return nil
}

// Close removes the system tray icon and cleans up resources
func (i *Icon) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !win.Shell_NotifyIcon(win.NIM_DELETE, i.nid) {
		return fmt.Errorf("failed to remove system tray icon")
	}

	// Clean up all cached icons
	for _, hIcon := range i.iconCache {
		if hIcon != 0 {
			win.DestroyIcon(hIcon)
		}
	}

	win.DestroyWindow(i.hwnd)
	return nil
}
