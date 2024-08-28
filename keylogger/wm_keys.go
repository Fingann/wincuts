package keylogger


const (
	WM_ACTIVATE    uintptr = 0x0006
	WM_APPCOMMAND  uintptr = 0x0319
	WM_CHAR        uintptr = 0x0102
	WM_DEADCHAR    uintptr = 0x0103
	WM_HOTKEY      uintptr = 0x0312
	WM_KEYDOWN     uintptr = 0x0100
	WM_KEYUP       uintptr = 0x0101
	WM_KILLFOCUS   uintptr = 0x0008
	WM_SETFOCUS    uintptr = 0x0007
	WM_SYSDEADCHAR uintptr = 0x0107
	WM_SYSKEYDOWN  uintptr = 0x0104
	WM_SYSKEYUP    uintptr = 0x0105
	WM_UNICHAR     uintptr = 0x0109
)