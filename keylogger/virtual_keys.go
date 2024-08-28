package keylogger

var CodeToVKMap = map[uint32]string{
	0x01: "VK_LBUTTON",    // Left mouse button
	0x02: "VK_RBUTTON",    // Right mouse button
	0x03: "VK_CANCEL",     // Control-break processing
	0x04: "VK_MBUTTON",    // Middle mouse button
	0x05: "VK_XBUTTON1",   // X1 mouse button
	0x06: "VK_XBUTTON2",   // X2 mouse button
	0x07: "Reserved",      // Reserved
	0x08: "VK_BACK",       // BACKSPACE key
	0x09: "VK_TAB",        // TAB key
	0x0A: "Reserved",      // Reserved
	0x0B: "Reserved",      // Reserved
	0x0C: "VK_CLEAR",      // CLEAR key
	0x0D: "VK_RETURN",     // ENTER key
	0x0E: "Unassigned",    // Unassigned
	0x0F: "Unassigned",    // Unassigned
	0x10: "VK_SHIFT",      // SHIFT key
	0x11: "VK_CONTROL",    // CTRL key
	0x12: "VK_MENU",       // ALT key
	0x13: "VK_PAUSE",      // PAUSE key
	0x14: "VK_CAPITAL",    // CAPS LOCK key
	0x15: "VK_KANA",       // IME Kana mode / Hangul mode
	0x16: "VK_IME_ON",     // IME On
	0x17: "VK_JUNJA",      // IME Junja mode
	0x18: "VK_FINAL",      // IME final mode
	0x19: "VK_HANJA",      // IME Hanja mode / Kanji mode
	0x1A: "VK_IME_OFF",    // IME Off
	0x1B: "VK_ESCAPE",     // ESC key
	0x1C: "VK_CONVERT",    // IME convert
	0x1D: "VK_NONCONVERT", // IME nonconvert
	0x1E: "VK_ACCEPT",     // IME accept
	0x1F: "VK_MODECHANGE", // IME mode change request
	0x20: "VK_SPACE",      // SPACEBAR
	0x21: "VK_PRIOR",      // PAGE UP key
	0x22: "VK_NEXT",       // PAGE DOWN key
	0x23: "VK_END",        // END key
	0x24: "VK_HOME",       // HOME key
	0x25: "VK_LEFT",       // LEFT ARROW key
	0x26: "VK_UP",         // UP ARROW key
	0x27: "VK_RIGHT",      // RIGHT ARROW key
	0x28: "VK_DOWN",       // DOWN ARROW key
	0x29: "VK_SELECT",     // SELECT key
	0x2A: "VK_PRINT",      // PRINT key
	0x2B: "VK_EXECUTE",    // EXECUTE key
	0x2C: "VK_SNAPSHOT",   // PRINT SCREEN key
	0x2D: "VK_INSERT",     // INS key
	0x2E: "VK_DELETE",     // DEL key
	0x2F: "VK_HELP",       // HELP key
	0x30: "VK_0",          // 0 key
	0x31: "VK_1",          // 1 key
	0x32: "VK_2",          // 2 key
	0x33: "VK_3",          // 3 key
	0x34: "VK_4",          // 4 key
	0x35: "VK_5",          // 5 key
	0x36: "VK_6",          // 6 key
	0x37: "VK_7",          // 7 key
	0x38: "VK_8",          // 8 key
	0x39: "VK_9",          // 9 key
	// 0x3A-40 Undefined
	0x41: "VK_A",    // A key
	0x42: "VK_B",    // B key
	0x43: "VK_C",    // C key
	0x44: "VK_D",    // D key
	0x45: "VK_E",    // E key
	0x46: "VK_F",    // F key
	0x47: "VK_G",    // G key
	0x48: "VK_H",    // H key
	0x49: "VK_I",    // I key
	0x4A: "VK_J",    // J key
	0x4B: "VK_K",    // K key
	0x4C: "VK_L",    // L key
	0x4D: "VK_M",    // M key
	0x4E: "VK_N",    // N key
	0x4F: "VK_O",    // O key
	0x50: "VK_P",    // P key
	0x51: "VK_Q",    // Q key
	0x52: "VK_R",    // R key
	0x53: "VK_S",    // S key
	0x54: "VK_T",    // T key
	0x55: "VK_U",    // U key
	0x56: "VK_V",    // V key
	0x57: "VK_W",    // W key
	0x58: "VK_X",    // X key
	0x59: "VK_Y",    // Y key
	0x5A: "VK_Z",    // Z key
	0x5B: "VK_LWIN", // Left Windows key
	0x5C: "VK_RWIN", // Right Windows key
	0x5D: "VK_APPS", // Applications key
	// 0x5E Reserved
	0x5F: "VK_SLEEP",     // Computer Sleep key
	0x60: "VK_NUMPAD0",   // Numeric keypad 0 key
	0x61: "VK_NUMPAD1",   // Numeric keypad 1 key
	0x62: "VK_NUMPAD2",   // Numeric keypad 2 key
	0x63: "VK_NUMPAD3",   // Numeric keypad 3 key
	0x64: "VK_NUMPAD4",   // Numeric keypad 4 key
	0x65: "VK_NUMPAD5",   // Numeric keypad 5 key
	0x66: "VK_NUMPAD6",   // Numeric keypad 6 key
	0x67: "VK_NUMPAD7",   // Numeric keypad 7 key
	0x68: "VK_NUMPAD8",   // Numeric keypad 8 key
	0x69: "VK_NUMPAD9",   // Numeric keypad 9 key
	0x6A: "VK_MULTIPLY",  // Multiply key
	0x6B: "VK_ADD",       // Add key
	0x6C: "VK_SEPARATOR", // Separator key
	0x6D: "VK_SUBTRACT",  // Subtract key
	0x6E: "VK_DECIMAL",   // Decimal key
	0x6F: "VK_DIVIDE",    // Divide key
	0x70: "VK_F1",        // F1 key
	0x71: "VK_F2",        // F2 key
	0x72: "VK_F3",        // F3 key
	0x73: "VK_F4",        // F4 key
	0x74: "VK_F5",        // F5 key
	0x75: "VK_F6",        // F6 key
	0x76: "VK_F7",        // F7 key
	0x77: "VK_F8",        // F8 key
	0x78: "VK_F9",        // F9 key
	0x79: "VK_F10",       // F10 key
	0x7A: "VK_F11",       // F11 key
	0x7B: "VK_F12",       // F12 key
	0x7C: "VK_F13",       // F13 key
	0x7D: "VK_F14",       // F14 key
	0x7E: "VK_F15",       // F15 key
	0x7F: "VK_F16",       // F16 key
	0x80: "VK_F17",       // F17 key
	0x81: "VK_F18",       // F18 key
	0x82: "VK_F19",       // F19 key
	0x83: "VK_F20",       // F20 key
	0x84: "VK_F21",       // F21 key
	0x85: "VK_F22",       // F22 key
	0x86: "VK_F23",       // F23 key
	0x87: "VK_F24",       // F24 key
	// 0x88-8F Reserved
	0x90: "VK_NUMLOCK", // NUM LOCK key
	0x91: "VK_SCROLL",  // SCROLL LOCK key
	// 0x92-96 OEM specific
	// 0x97-9F Unassigned
	0xA0: "VK_LSHIFT",              // Left SHIFT key
	0xA1: "VK_RSHIFT",              // Right SHIFT key
	0xA2: "VK_LCONTROL",            // Left CONTROL key
	0xA3: "VK_RCONTROL",            // Right CONTROL key
	0xA4: "VK_LMENU",               // Left ALT key
	0xA5: "VK_RMENU",               // Right ALT key
	0xA6: "VK_BROWSER_BACK",        // Browser Back key
	0xA7: "VK_BROWSER_FORWARD",     // Browser Forward key
	0xA8: "VK_BROWSER_REFRESH",     // Browser Refresh key
	0xA9: "VK_BROWSER_STOP",        // Browser Stop key
	0xAA: "VK_BROWSER_SEARCH",      // Browser Search key
	0xAB: "VK_BROWSER_FAVORITES",   // Browser Favorites key
	0xAC: "VK_BROWSER_HOME",        // Browser Start and Home key
	0xAD: "VK_VOLUME_MUTE",         // Volume Mute key
	0xAE: "VK_VOLUME_DOWN",         // Volume Down key
	0xAF: "VK_VOLUME_UP",           // Volume Up key
	0xB0: "VK_MEDIA_NEXT_TRACK",    // Next Track key
	0xB1: "VK_MEDIA_PREV_TRACK",    // Previous Track key
	0xB2: "VK_MEDIA_STOP",          // Stop Media key
	0xB3: "VK_MEDIA_PLAY_PAUSE",    // Play/Pause Media key
	0xB4: "VK_LAUNCH_MAIL",         // Start Mail key
	0xB5: "VK_LAUNCH_MEDIA_SELECT", // Select Media key
	0xB6: "VK_LAUNCH_APP1",         // Start Application 1 key
	0xB7: "VK_LAUNCH_APP2",         // Start Application 2 key
	// 0xB8-B9 Reserved
	0xBA: "VK_OEM_1",      // Used for miscellaneous characters; varies by keyboard. For the US standard keyboard, the `;:` key
	0xBB: "VK_OEM_PLUS",   // For any country/region, the `+` key
	0xBC: "VK_OEM_COMMA",  // For any country/region, the `,` key
	0xBD: "VK_OEM_MINUS",  // For any country/region, the `-` key
	0xBE: "VK_OEM_PERIOD", // For any country/region, the `.` key
	0xBF: "VK_OEM_2",      // Used for miscellaneous characters; varies by keyboard. For the US standard keyboard, the `/?` key
	0xC0: "VK_OEM_3",      // Used for miscellaneous characters; varies by keyboard. For the US standard keyboard, the ``~` key
	// 0xC1-DA Reserved
	0xDB: "VK_OEM_4",     // Used for miscellaneous characters; varies by keyboard. For the US standard keyboard, the `[{` key
	0xDC: "VK_OEM_5",     // Used for miscellaneous characters; varies by keyboard. For the US standard keyboard, the `\\|` key
	0xDD: "VK_OEM_6",     // Used for miscellaneous characters; varies by keyboard. For the US standard keyboard, the `]}` key
	0xDE: "VK_OEM_7",     // Used for miscellaneous characters; varies by keyboard. For the US standard keyboard, the `'\"` key
	0xDF: "VK_OEM_8",     // Used for miscellaneous characters; varies by keyboard.
	0xE0: "Reserved",     // Reserved
	0xE1: "OEM specific", // OEM specific
	0xE2: "VK_OEM_102",   // The `<>` keys on the US standard keyboard, or the `\\|` key on the non-US 102-key keyboard
	// 0xE3-E4 OEM specific
	0xE5: "VK_PROCESSKEY", // IME PROCESS key
	// 0xE6 OEM specific
	0xE7: "VK_PACKET", // Used to pass Unicode characters as if they were keystrokes
	// 0xE8 Unassigned
	// 0xE9-F5 OEM specific
	0xF6: "VK_ATTN",      // Attn key
	0xF7: "VK_CRSEL",     // CrSel key
	0xF8: "VK_EXSEL",     // ExSel key
	0xF9: "VK_EREOF",     // Erase EOF key
	0xFA: "VK_PLAY",      // Play key
	0xFB: "VK_ZOOM",      // Zoom key
	0xFC: "VK_NONAME",    // Reserved
	0xFD: "VK_PA1",       // PA1 key
	0xFE: "VK_OEM_CLEAR", // Clear key
}

var VKToCodeMap = swapMap(CodeToVKMap)

func swapMap(vkMap map[uint32]string) map[string]uint32 {
	vkNameMap := make(map[string]uint32)
	for k, v := range vkMap {
		vkNameMap[v] = k
	}
	return vkNameMap
}
