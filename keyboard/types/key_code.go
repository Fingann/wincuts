package types

var (
	VKModifierMap = map[VirtualKey]bool{
		VK_LSHIFT:   true,
		VK_RSHIFT:   true,
		VK_LCONTROL: true,
		VK_RCONTROL: true,
		VK_LMENU:    true,
		VK_RMENU:    true,
		VK_LWIN:     true,
		VK_RWIN:     true,
	}
)

const (
	VK_LBUTTON             VirtualKey = 0x01 // Left mouse button
	VK_RBUTTON             VirtualKey = 0x02 // Right mouse button
	VK_CANCEL              VirtualKey = 0x03 // Control-break processing
	VK_MBUTTON             VirtualKey = 0x04 // Middle mouse button (three-button mouse)
	VK_XBUTTON1            VirtualKey = 0x05 // X1 mouse button
	VK_XBUTTON2            VirtualKey = 0x06 // X2 mouse button
	VK_BACK                VirtualKey = 0x08 // BACKSPACE key
	VK_TAB                 VirtualKey = 0x09 // TAB key
	VK_CLEAR               VirtualKey = 0x0C // CLEAR key
	VK_RETURN              VirtualKey = 0x0D // ENTER key
	VK_SHIFT               VirtualKey = 0x10 // SHIFT key
	VK_CONTROL             VirtualKey = 0x11 // CTRL key
	VK_MENU                VirtualKey = 0x12 // ALT key
	VK_PAUSE               VirtualKey = 0x13 // PAUSE key
	VK_CAPITAL             VirtualKey = 0x14 // CAPS LOCK key
	VK_KANA                VirtualKey = 0x15 // IME Kana mode
	VK_HANGUEL             VirtualKey = 0x15 // IME Hanguel mode (maintained for compatibility; use VK_HANGUL)
	VK_HANGUL              VirtualKey = 0x15 // IME Hangul mode
	VK_IME_ON              VirtualKey = 0x16 // IME On
	VK_JUNJA               VirtualKey = 0x17 // IME Junja mode
	VK_FINAL               VirtualKey = 0x18 // IME final mode
	VK_HANJA               VirtualKey = 0x19 // IME Hanja mode
	VK_KANJI               VirtualKey = 0x19 // IME Kanji mode
	VK_IME_OFF             VirtualKey = 0x1A // IME Off
	VK_ESCAPE              VirtualKey = 0x1B // ESC key
	VK_CONVERT             VirtualKey = 0x1C // IME convert
	VK_NONCONVERT          VirtualKey = 0x1D // IME nonconvert
	VK_ACCEPT              VirtualKey = 0x1E // IME accept
	VK_MODECHANGE          VirtualKey = 0x1F // IME mode change request
	VK_SPACE               VirtualKey = 0x20 // SPACEBAR
	VK_PRIOR               VirtualKey = 0x21 // PAGE UP key
	VK_NEXT                VirtualKey = 0x22 // PAGE DOWN key
	VK_END                 VirtualKey = 0x23 // END key
	VK_HOME                VirtualKey = 0x24 // HOME key
	VK_LEFT                VirtualKey = 0x25 // LEFT ARROW key
	VK_UP                  VirtualKey = 0x26 // UP ARROW key
	VK_RIGHT               VirtualKey = 0x27 // RIGHT ARROW key
	VK_DOWN                VirtualKey = 0x28 // DOWN ARROW key
	VK_SELECT              VirtualKey = 0x29 // SELECT key
	VK_PRINT               VirtualKey = 0x2A // PRINT key
	VK_EXECUTE             VirtualKey = 0x2B // EXECUTE key
	VK_SNAPSHOT            VirtualKey = 0x2C // PRINT SCREEN key
	VK_INSERT              VirtualKey = 0x2D // INS key
	VK_DELETE              VirtualKey = 0x2E // DEL key
	VK_HELP                VirtualKey = 0x2F // HELP key
	VK_0                   VirtualKey = 0x30 // 0 key
	VK_1                   VirtualKey = 0x31 // 1 key
	VK_2                   VirtualKey = 0x32 // 2 key
	VK_3                   VirtualKey = 0x33 // 3 key
	VK_4                   VirtualKey = 0x34 // 4 key
	VK_5                   VirtualKey = 0x35 // 5 key
	VK_6                   VirtualKey = 0x36 // 6 key
	VK_7                   VirtualKey = 0x37 // 7 key
	VK_8                   VirtualKey = 0x38 // 8 key
	VK_9                   VirtualKey = 0x39 // 9 key
	VK_A                   VirtualKey = 0x41 // A key
	VK_B                   VirtualKey = 0x42 // B key
	VK_C                   VirtualKey = 0x43 // C key
	VK_D                   VirtualKey = 0x44 // D key
	VK_E                   VirtualKey = 0x45 // E key
	VK_F                   VirtualKey = 0x46 // F key
	VK_G                   VirtualKey = 0x47 // G key
	VK_H                   VirtualKey = 0x48 // H key
	VK_I                   VirtualKey = 0x49 // I key
	VK_J                   VirtualKey = 0x4A // J key
	VK_K                   VirtualKey = 0x4B // K key
	VK_L                   VirtualKey = 0x4C // L key
	VK_M                   VirtualKey = 0x4D // M key
	VK_N                   VirtualKey = 0x4E // N key
	VK_O                   VirtualKey = 0x4F // O key
	VK_P                   VirtualKey = 0x50 // P key
	VK_Q                   VirtualKey = 0x51 // Q key
	VK_R                   VirtualKey = 0x52 // R key
	VK_S                   VirtualKey = 0x53 // S key
	VK_T                   VirtualKey = 0x54 // T key
	VK_U                   VirtualKey = 0x55 // U key
	VK_V                   VirtualKey = 0x56 // V key
	VK_W                   VirtualKey = 0x57 // W key
	VK_X                   VirtualKey = 0x58 // X key
	VK_Y                   VirtualKey = 0x59 // Y key
	VK_Z                   VirtualKey = 0x5A // Z key
	VK_LWIN                VirtualKey = 0x5B // Left Windows key (Natural keyboard)
	VK_RWIN                VirtualKey = 0x5C // Right Windows key (Natural keyboard)
	VK_APPS                VirtualKey = 0x5D // Applications key (Natural keyboard)
	VK_SLEEP               VirtualKey = 0x5F // Computer Sleep key
	VK_NUMPAD0             VirtualKey = 0x60 // Numeric keypad 0 key
	VK_NUMPAD1             VirtualKey = 0x61 // Numeric keypad 1 key
	VK_NUMPAD2             VirtualKey = 0x62 // Numeric keypad 2 key
	VK_NUMPAD3             VirtualKey = 0x63 // Numeric keypad 3 key
	VK_NUMPAD4             VirtualKey = 0x64 // Numeric keypad 4 key
	VK_NUMPAD5             VirtualKey = 0x65 // Numeric keypad 5 key
	VK_NUMPAD6             VirtualKey = 0x66 // Numeric keypad 6 key
	VK_NUMPAD7             VirtualKey = 0x67 // Numeric keypad 7 key
	VK_NUMPAD8             VirtualKey = 0x68 // Numeric keypad 8 key
	VK_NUMPAD9             VirtualKey = 0x69 // Numeric keypad 9 key
	VK_MULTIPLY            VirtualKey = 0x6A // Multiply key
	VK_ADD                 VirtualKey = 0x6B // Add key
	VK_SEPARATOR           VirtualKey = 0x6C // Separator key
	VK_SUBTRACT            VirtualKey = 0x6D // Subtract key
	VK_DECIMAL             VirtualKey = 0x6E // Decimal key
	VK_DIVIDE              VirtualKey = 0x6F // Divide key
	VK_F1                  VirtualKey = 0x70 // F1 key
	VK_F2                  VirtualKey = 0x71 // F2 key
	VK_F3                  VirtualKey = 0x72 // F3 key
	VK_F4                  VirtualKey = 0x73 // F4 key
	VK_F5                  VirtualKey = 0x74 // F5 key
	VK_F6                  VirtualKey = 0x75 // F6 key
	VK_F7                  VirtualKey = 0x76 // F7 key
	VK_F8                  VirtualKey = 0x77 // F8 key
	VK_F9                  VirtualKey = 0x78 // F9 key
	VK_F10                 VirtualKey = 0x79 // F10 key
	VK_F11                 VirtualKey = 0x7A // F11 key
	VK_F12                 VirtualKey = 0x7B // F12 key
	VK_F13                 VirtualKey = 0x7C // F13 key
	VK_F14                 VirtualKey = 0x7D // F14 key
	VK_F15                 VirtualKey = 0x7E // F15 key
	VK_F16                 VirtualKey = 0x7F // F16 key
	VK_F17                 VirtualKey = 0x80 // F17 key
	VK_F18                 VirtualKey = 0x81 // F18 key
	VK_F19                 VirtualKey = 0x82 // F19 key
	VK_F20                 VirtualKey = 0x83 // F20 key
	VK_F21                 VirtualKey = 0x84 // F21 key
	VK_F22                 VirtualKey = 0x85 // F22 key
	VK_F23                 VirtualKey = 0x86 // F23 key
	VK_F24                 VirtualKey = 0x87 // F24 key
	VK_NUMLOCK             VirtualKey = 0x90 // NUM LOCK key
	VK_SCROLL              VirtualKey = 0x91 // SCROLL LOCK key
	VK_LSHIFT              VirtualKey = 0xA0 // Left SHIFT key
	VK_RSHIFT              VirtualKey = 0xA1 // Right SHIFT key
	VK_LCONTROL            VirtualKey = 0xA2 // Left CONTROL key
	VK_RCONTROL            VirtualKey = 0xA3 // Right CONTROL key
	VK_LMENU               VirtualKey = 0xA4 // Left MENU key
	VK_RMENU               VirtualKey = 0xA5 // Right MENU key
	VK_BROWSER_BACK        VirtualKey = 0xA6 // Browser Back key
	VK_BROWSER_FORWARD     VirtualKey = 0xA7 // Browser Forward key
	VK_BROWSER_REFRESH     VirtualKey = 0xA8 // Browser Refresh key
	VK_BROWSER_STOP        VirtualKey = 0xA9 // Browser Stop key
	VK_BROWSER_SEARCH      VirtualKey = 0xAA // Browser Search key
	VK_BROWSER_FAVORITES   VirtualKey = 0xAB // Browser Favorites key
	VK_BROWSER_HOME        VirtualKey = 0xAC // Browser Start and Home key
	VK_VOLUME_MUTE         VirtualKey = 0xAD // Volume Mute key
	VK_VOLUME_DOWN         VirtualKey = 0xAE // Volume Down key
	VK_VOLUME_UP           VirtualKey = 0xAF // Volume Up key
	VK_MEDIA_NEXT_TRACK    VirtualKey = 0xB0 // Next Track key
	VK_MEDIA_PREV_TRACK    VirtualKey = 0xB1 // Previous Track key
	VK_MEDIA_STOP          VirtualKey = 0xB2 // Stop Media key
	VK_MEDIA_PLAY_PAUSE    VirtualKey = 0xB3 // Play/Pause Media key
	VK_LAUNCH_MAIL         VirtualKey = 0xB4 // Start Mail key
	VK_LAUNCH_MEDIA_SELECT VirtualKey = 0xB5 // Select Media key
	VK_LAUNCH_APP1         VirtualKey = 0xB6 // Start Application 1 key
	VK_LAUNCH_APP2         VirtualKey = 0xB7 // Start Application 2 key
	VK_OEM_1               VirtualKey = 0xBA // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ';:' key
	VK_OEM_PLUS            VirtualKey = 0xBB // For any country/region, the '+' key
	VK_OEM_COMMA           VirtualKey = 0xBC // For any country/region, the ',' key
	VK_OEM_MINUS           VirtualKey = 0xBD // For any country/region, the '-' key
	VK_OEM_PERIOD          VirtualKey = 0xBE // For any country/region, the '.' key
	VK_OEM_2               VirtualKey = 0xBF // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '/?' key
	VK_OEM_3               VirtualKey = 0xC0 // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '`~' key
	VK_OEM_4               VirtualKey = 0xDB // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '[{' key
	VK_OEM_5               VirtualKey = 0xDC // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '\|' key
	VK_OEM_6               VirtualKey = 0xDD // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ']}' key
	VK_OEM_7               VirtualKey = 0xDE // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the 'single-quote/double-quote' key
	VK_OEM_8               VirtualKey = 0xDF // Used for miscellaneous characters; it can vary by keyboard.
	VK_OEM_102             VirtualKey = 0xE2 // Either the angle bracket key or the backslash key on the RT 102-key keyboard
	VK_PROCESSKEY          VirtualKey = 0xE5 // IME PROCESS key
	VK_PACKET              VirtualKey = 0xE7 // Used to pass Unicode characters as if they were keystrokes. The VK_PACKET key is the low word of a 32-bit Virtual Key value used for non-keyboard input methods. For more information, see Remark in KEYBDINPUT, SendInput, WM_KEYDOWN, and WM_KEYUP
	VK_ATTN                VirtualKey = 0xF6 // Attn key
	VK_CRSEL               VirtualKey = 0xF7 // CrSel key
	VK_EXSEL               VirtualKey = 0xF8 // ExSel key
	VK_EREOF               VirtualKey = 0xF9 // Erase EOF key
	VK_PLAY                VirtualKey = 0xFA // Play key
	VK_ZOOM                VirtualKey = 0xFB // Zoom key
	VK_NONAME              VirtualKey = 0xFC // Reserved
	VK_PA1                 VirtualKey = 0xFD // PA1 key
	VK_OEM_CLEAR           VirtualKey = 0xFE // Clear key
)

var KeyNameToVKCode = map[string]VirtualKey{
	"LBUTTON":             VK_LBUTTON,
	"RBUTTON":             VK_RBUTTON,
	"CANCEL":              VK_CANCEL,
	"MBUTTON":             VK_MBUTTON,
	"XBUTTON1":            VK_XBUTTON1,
	"XBUTTON2":            VK_XBUTTON2,
	"BACK":                VK_BACK,
	"TAB":                 VK_TAB,
	"CLEAR":               VK_CLEAR,
	"RETURN":              VK_RETURN,
	"SHIFT":               VK_SHIFT,
	"CONTROL":             VK_CONTROL,
	"MENU":                VK_MENU,
	"PAUSE":               VK_PAUSE,
	"CAPITAL":             VK_CAPITAL,
	"KANA":                VK_KANA,
	"HANGUEL":             VK_HANGUEL,
	"HANGUL":              VK_HANGUL,
	"IME_ON":              VK_IME_ON,
	"JUNJA":               VK_JUNJA,
	"FINAL":               VK_FINAL,
	"HANJA":               VK_HANJA,
	"KANJI":               VK_KANJI,
	"IME_OFF":             VK_IME_OFF,
	"ESCAPE":              VK_ESCAPE,
	"CONVERT":             VK_CONVERT,
	"NONCONVERT":          VK_NONCONVERT,
	"ACCEPT":              VK_ACCEPT,
	"MODECHANGE":          VK_MODECHANGE,
	"SPACE":               VK_SPACE,
	"PRIOR":               VK_PRIOR,
	"NEXT":                VK_NEXT,
	"END":                 VK_END,
	"HOME":                VK_HOME,
	"LEFT":                VK_LEFT,
	"UP":                  VK_UP,
	"RIGHT":               VK_RIGHT,
	"DOWN":                VK_DOWN,
	"SELECT":              VK_SELECT,
	"PRINT":               VK_PRINT,
	"EXECUTE":             VK_EXECUTE,
	"SNAPSHOT":            VK_SNAPSHOT,
	"INSERT":              VK_INSERT,
	"DELETE":              VK_DELETE,
	"HELP":                VK_HELP,
	"0":                   VK_0,
	"1":                   VK_1,
	"2":                   VK_2,
	"3":                   VK_3,
	"4":                   VK_4,
	"5":                   VK_5,
	"6":                   VK_6,
	"7":                   VK_7,
	"8":                   VK_8,
	"9":                   VK_9,
	"A":                   VK_A,
	"B":                   VK_B,
	"C":                   VK_C,
	"D":                   VK_D,
	"E":                   VK_E,
	"F":                   VK_F,
	"G":                   VK_G,
	"H":                   VK_H,
	"I":                   VK_I,
	"J":                   VK_J,
	"K":                   VK_K,
	"L":                   VK_L,
	"M":                   VK_M,
	"N":                   VK_N,
	"O":                   VK_O,
	"P":                   VK_P,
	"Q":                   VK_Q,
	"R":                   VK_R,
	"S":                   VK_S,
	"T":                   VK_T,
	"U":                   VK_U,
	"V":                   VK_V,
	"W":                   VK_W,
	"X":                   VK_X,
	"Y":                   VK_Y,
	"Z":                   VK_Z,
	"LWIN":                VK_LWIN,
	"RWIN":                VK_RWIN,
	"APPS":                VK_APPS,
	"SLEEP":               VK_SLEEP,
	"NUMPAD0":             VK_NUMPAD0,
	"NUMPAD1":             VK_NUMPAD1,
	"NUMPAD2":             VK_NUMPAD2,
	"NUMPAD3":             VK_NUMPAD3,
	"NUMPAD4":             VK_NUMPAD4,
	"NUMPAD5":             VK_NUMPAD5,
	"NUMPAD6":             VK_NUMPAD6,
	"NUMPAD7":             VK_NUMPAD7,
	"NUMPAD8":             VK_NUMPAD8,
	"NUMPAD9":             VK_NUMPAD9,
	"MULTIPLY":            VK_MULTIPLY,
	"ADD":                 VK_ADD,
	"SEPARATOR":           VK_SEPARATOR,
	"SUBTRACT":            VK_SUBTRACT,
	"DECIMAL":             VK_DECIMAL,
	"DIVIDE":              VK_DIVIDE,
	"F1":                  VK_F1,
	"F2":                  VK_F2,
	"F3":                  VK_F3,
	"F4":                  VK_F4,
	"F5":                  VK_F5,
	"F6":                  VK_F6,
	"F7":                  VK_F7,
	"F8":                  VK_F8,
	"F9":                  VK_F9,
	"F10":                 VK_F10,
	"F11":                 VK_F11,
	"F12":                 VK_F12,
	"F13":                 VK_F13,
	"F14":                 VK_F14,
	"F15":                 VK_F15,
	"F16":                 VK_F16,
	"F17":                 VK_F17,
	"F18":                 VK_F18,
	"F19":                 VK_F19,
	"F20":                 VK_F20,
	"F21":                 VK_F21,
	"F22":                 VK_F22,
	"F23":                 VK_F23,
	"F24":                 VK_F24,
	"NUMLOCK":             VK_NUMLOCK,
	"SCROLL":              VK_SCROLL,
	"LSHIFT":              VK_LSHIFT,
	"RSHIFT":              VK_RSHIFT,
	"LCONTROL":            VK_LCONTROL,
	"RCONTROL":            VK_RCONTROL,
	"LMENU":               VK_LMENU,
	"RMENU":               VK_RMENU,
	"BROWSER_BACK":        VK_BROWSER_BACK,
	"BROWSER_FORWARD":     VK_BROWSER_FORWARD,
	"BROWSER_REFRESH":     VK_BROWSER_REFRESH,
	"BROWSER_STOP":        VK_BROWSER_STOP,
	"BROWSER_SEARCH":      VK_BROWSER_SEARCH,
	"BROWSER_FAVORITES":   VK_BROWSER_FAVORITES,
	"BROWSER_HOME":        VK_BROWSER_HOME,
	"VOLUME_MUTE":         VK_VOLUME_MUTE,
	"VOLUME_DOWN":         VK_VOLUME_DOWN,
	"VOLUME_UP":           VK_VOLUME_UP,
	"MEDIA_NEXT_TRACK":    VK_MEDIA_NEXT_TRACK,
	"MEDIA_PREV_TRACK":    VK_MEDIA_PREV_TRACK,
	"MEDIA_STOP":          VK_MEDIA_STOP,
	"MEDIA_PLAY_PAUSE":    VK_MEDIA_PLAY_PAUSE,
	"LAUNCH_MAIL":         VK_LAUNCH_MAIL,
	"LAUNCH_MEDIA_SELECT": VK_LAUNCH_MEDIA_SELECT,
	"LAUNCH_APP1":         VK_LAUNCH_APP1,
	"LAUNCH_APP2":         VK_LAUNCH_APP2,
	"OEM_1":               VK_OEM_1,
	"OEM_PLUS":            VK_OEM_PLUS,
	"OEM_COMMA":           VK_OEM_COMMA,
	"OEM_MINUS":           VK_OEM_MINUS,
	"OEM_PERIOD":          VK_OEM_PERIOD,
	"OEM_2":               VK_OEM_2,
	"OEM_3":               VK_OEM_3,
	"OEM_4":               VK_OEM_4,
	"OEM_5":               VK_OEM_5,
	"OEM_6":               VK_OEM_6,
	"OEM_7":               VK_OEM_7,
	"OEM_8":               VK_OEM_8,
	"OEM_102":             VK_OEM_102,
	"PROCESSKEY":          VK_PROCESSKEY,
	"PACKET":              VK_PACKET,
	"ATTN":                VK_ATTN,
	"CRSEL":               VK_CRSEL,
	"EXSEL":               VK_EXSEL,
	"EREOF":               VK_EREOF,
	"PLAY":                VK_PLAY,
	"ZOOM":                VK_ZOOM,
	"NONAME":              VK_NONAME,
	"PA1":                 VK_PA1,
	"OEM_CLEAR":           VK_OEM_CLEAR,
}

var VKCodeToKeyName = map[VirtualKey]string{
	VK_LBUTTON:             "LBUTTON",
	VK_RBUTTON:             "RBUTTON",
	VK_CANCEL:              "CANCEL",
	VK_MBUTTON:             "MBUTTON",
	VK_XBUTTON1:            "XBUTTON1",
	VK_XBUTTON2:            "XBUTTON2",
	VK_BACK:                "BACK",
	VK_TAB:                 "TAB",
	VK_CLEAR:               "CLEAR",
	VK_RETURN:              "RETURN",
	VK_SHIFT:               "SHIFT",
	VK_CONTROL:             "CONTROL",
	VK_MENU:                "MENU",
	VK_PAUSE:               "PAUSE",
	VK_CAPITAL:             "CAPITAL",
	VK_KANA:                "KANA",
	VK_IME_ON:              "IME_ON",
	VK_JUNJA:               "JUNJA",
	VK_FINAL:               "FINAL",
	VK_HANJA:               "HANJA",
	VK_IME_OFF:             "IME_OFF",
	VK_ESCAPE:              "ESCAPE",
	VK_CONVERT:             "CONVERT",
	VK_NONCONVERT:          "NONCONVERT",
	VK_ACCEPT:              "ACCEPT",
	VK_MODECHANGE:          "MODECHANGE",
	VK_SPACE:               "SPACE",
	VK_PRIOR:               "PRIOR",
	VK_NEXT:                "NEXT",
	VK_END:                 "END",
	VK_HOME:                "HOME",
	VK_LEFT:                "LEFT",
	VK_UP:                  "UP",
	VK_RIGHT:               "RIGHT",
	VK_DOWN:                "DOWN",
	VK_SELECT:              "SELECT",
	VK_PRINT:               "PRINT",
	VK_EXECUTE:             "EXECUTE",
	VK_SNAPSHOT:            "SNAPSHOT",
	VK_INSERT:              "INSERT",
	VK_DELETE:              "DELETE",
	VK_HELP:                "HELP",
	VK_0:                   "0",
	VK_1:                   "1",
	VK_2:                   "2",
	VK_3:                   "3",
	VK_4:                   "4",
	VK_5:                   "5",
	VK_6:                   "6",
	VK_7:                   "7",
	VK_8:                   "8",
	VK_9:                   "9",
	VK_A:                   "A",
	VK_B:                   "B",
	VK_C:                   "C",
	VK_D:                   "D",
	VK_E:                   "E",
	VK_F:                   "F",
	VK_G:                   "G",
	VK_H:                   "H",
	VK_I:                   "I",
	VK_J:                   "J",
	VK_K:                   "K",
	VK_L:                   "L",
	VK_M:                   "M",
	VK_N:                   "N",
	VK_O:                   "O",
	VK_P:                   "P",
	VK_Q:                   "Q",
	VK_R:                   "R",
	VK_S:                   "S",
	VK_T:                   "T",
	VK_U:                   "U",
	VK_V:                   "V",
	VK_W:                   "W",
	VK_X:                   "X",
	VK_Y:                   "Y",
	VK_Z:                   "Z",
	VK_LWIN:                "LWIN",
	VK_RWIN:                "RWIN",
	VK_APPS:                "APPS",
	VK_SLEEP:               "SLEEP",
	VK_NUMPAD0:             "NUMPAD0",
	VK_NUMPAD1:             "NUMPAD1",
	VK_NUMPAD2:             "NUMPAD2",
	VK_NUMPAD3:             "NUMPAD3",
	VK_NUMPAD4:             "NUMPAD4",
	VK_NUMPAD5:             "NUMPAD5",
	VK_NUMPAD6:             "NUMPAD6",
	VK_NUMPAD7:             "NUMPAD7",
	VK_NUMPAD8:             "NUMPAD8",
	VK_NUMPAD9:             "NUMPAD9",
	VK_MULTIPLY:            "MULTIPLY",
	VK_ADD:                 "ADD",
	VK_SEPARATOR:           "SEPARATOR",
	VK_SUBTRACT:            "SUBTRACT",
	VK_DECIMAL:             "DECIMAL",
	VK_DIVIDE:              "DIVIDE",
	VK_F1:                  "F1",
	VK_F2:                  "F2",
	VK_F3:                  "F3",
	VK_F4:                  "F4",
	VK_F5:                  "F5",
	VK_F6:                  "F6",
	VK_F7:                  "F7",
	VK_F8:                  "F8",
	VK_F9:                  "F9",
	VK_F10:                 "F10",
	VK_F11:                 "F11",
	VK_F12:                 "F12",
	VK_F13:                 "F13",
	VK_F14:                 "F14",
	VK_F15:                 "F15",
	VK_F16:                 "F16",
	VK_F17:                 "F17",
	VK_F18:                 "F18",
	VK_F19:                 "F19",
	VK_F20:                 "F20",
	VK_F21:                 "F21",
	VK_F22:                 "F22",
	VK_F23:                 "F23",
	VK_F24:                 "F24",
	VK_NUMLOCK:             "NUMLOCK",
	VK_SCROLL:              "SCROLL",
	VK_LSHIFT:              "LSHIFT",
	VK_RSHIFT:              "RSHIFT",
	VK_LCONTROL:            "LCONTROL",
	VK_RCONTROL:            "RCONTROL",
	VK_LMENU:               "LMENU",
	VK_RMENU:               "RMENU",
	VK_BROWSER_BACK:        "BROWSER_BACK",
	VK_BROWSER_FORWARD:     "BROWSER_FORWARD",
	VK_BROWSER_REFRESH:     "BROWSER_REFRESH",
	VK_BROWSER_STOP:        "BROWSER_STOP",
	VK_BROWSER_SEARCH:      "BROWSER_SEARCH",
	VK_BROWSER_FAVORITES:   "BROWSER_FAVORITES",
	VK_BROWSER_HOME:        "BROWSER_HOME",
	VK_VOLUME_MUTE:         "VOLUME_MUTE",
	VK_VOLUME_DOWN:         "VOLUME_DOWN",
	VK_VOLUME_UP:           "VOLUME_UP",
	VK_MEDIA_NEXT_TRACK:    "MEDIA_NEXT_TRACK",
	VK_MEDIA_PREV_TRACK:    "MEDIA_PREV_TRACK",
	VK_MEDIA_STOP:          "MEDIA_STOP",
	VK_MEDIA_PLAY_PAUSE:    "MEDIA_PLAY_PAUSE",
	VK_LAUNCH_MAIL:         "LAUNCH_MAIL",
	VK_LAUNCH_MEDIA_SELECT: "LAUNCH_MEDIA_SELECT",
	VK_LAUNCH_APP1:         "LAUNCH_APP1",
	VK_LAUNCH_APP2:         "LAUNCH_APP2",
	VK_OEM_1:               "OEM_1",
	VK_OEM_PLUS:            "OEM_PLUS",
	VK_OEM_COMMA:           "OEM_COMMA",
	VK_OEM_MINUS:           "OEM_MINUS",
	VK_OEM_PERIOD:          "OEM_PERIOD",
	VK_OEM_2:               "OEM_2",
	VK_OEM_3:               "OEM_3",
	VK_OEM_4:               "OEM_4",
	VK_OEM_5:               "OEM_5",
	VK_OEM_6:               "OEM_6",
	VK_OEM_7:               "OEM_7",
	VK_OEM_8:               "OEM_8",
	VK_OEM_102:             "OEM_102",
	VK_PROCESSKEY:          "PROCESSKEY",
	VK_PACKET:              "PACKET",
	VK_ATTN:                "ATTN",
	VK_CRSEL:               "CRSEL",
	VK_EXSEL:               "EXSEL",
	VK_EREOF:               "EREOF",
	VK_PLAY:                "PLAY",
	VK_ZOOM:                "ZOOM",
	VK_NONAME:              "NONAME",
	VK_PA1:                 "PA1",
	VK_OEM_CLEAR:           "OEM_CLEAR",
}
