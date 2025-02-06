# Installing WinCuts

## Quick Install (Recommended)
Open PowerShell and run:
```powershell
irm https://raw.githubusercontent.com/fingann/WinCuts/main/install.ps1 | iex
```

That's it! The script will:
- Download the latest version
- Install to the correct locations
- Create a startup shortcut
- Start WinCuts automatically
- Preserve your configuration during updates

## Manual Installation
If you prefer to install manually:
1. Download the latest release zip file from [GitHub Releases](https://github.com/fingann/WinCuts/releases)
2. Extract the zip file to a permanent location (e.g., `C:\Program Files\WinCuts`)
3. Double-click `WinCuts.exe` to start the application
4. (Optional) Create a shortcut in your startup folder to run WinCuts automatically at login:

   - Press `Win + R`
   - Type `shell:startup`
   - Create a shortcut to `WinCuts.exe` in this folder

## Configuration
1. On first run, WinCuts will use the default configuration
2. To customize:
   - Edit `%APPDATA%\WinCuts\config.yaml`
   - Restart WinCuts to apply changes

## Default Keyboard Shortcuts
- `Alt + [1-9]`: Switch to desktop 1-9
- `Alt + Shift + [1-9]`: Move current window to desktop 1-9
- `Alt + Up`: Maximize current window
- `Alt + Down`: Minimize current window
- `Alt + Space`: Toggle window state

## Command Line Options
```powershell
# Show version
WinCuts.exe --version

# Use custom config file
WinCuts.exe --config path/to/config.yaml

# Set logging level
WinCuts.exe --log-level DEBUG
```

## Updating
Run the installation command again to update to the latest version. Your configuration will be preserved.

## Uninstalling
Run in PowerShell:
```powershell
Remove-Item "$env:LOCALAPPDATA\WinCuts","$env:APPDATA\WinCuts","$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup\WinCuts.lnk" -Recurse -Force
```

## Troubleshooting
- If the system tray icon doesn't appear, check the Windows Event Viewer for errors
- For support, please [open an issue](https://github.com/yourusername/WinCuts/issues) 