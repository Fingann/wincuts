# Installation script for WinCuts
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

try {
    $ErrorActionPreference = "Stop"

    # Download latest release info from GitHub
    $repo = "fingann/WinCuts"
    $releases = "https://api.github.com/repos/$repo/releases/latest"
    $tag = (Invoke-WebRequest $releases -UseBasicParsing | ConvertFrom-Json).tag_name
    
    # Download URLs
    $exeUrl = "https://github.com/$repo/releases/download/$tag/WinCuts.exe"
    $configUrl = "https://raw.githubusercontent.com/$repo/$tag/default_config.yaml"
    
    Write-Host "🔍 Found latest version: $tag"
    
    # Installation directories
    $installDir = "$env:LOCALAPPDATA\WinCuts"
    $configDir = "$env:APPDATA\WinCuts"
    
    Write-Host "📦 Installing WinCuts..."
    
    # Check for existing WinCuts process and stop it
    $existingProcess = Get-Process "WinCuts" -ErrorAction SilentlyContinue
    if ($existingProcess) {
        Write-Host "⏹️ Stopping existing WinCuts process..."
        $existingProcess | Stop-Process -Force
        Start-Sleep -Seconds 1
    }
    
    # Create directories
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    New-Item -ItemType Directory -Force -Path $configDir | Out-Null
    
    # Backup existing config if it exists
    $configPath = "$configDir\config.yaml"
    if (Test-Path $configPath) {
        Write-Host "💾 Backing up existing configuration..."
        Copy-Item $configPath "$configPath.backup"
    }
    
    # Clean up old files
    Write-Host "🧹 Cleaning up old files..."
    Remove-Item "$installDir\*" -Force -Recurse -ErrorAction SilentlyContinue
    
    # Download and install new files
    Write-Host "⬇️ Downloading latest version..."
    Invoke-WebRequest $exeUrl -OutFile "$installDir\WinCuts.exe" -UseBasicParsing
    
    if (-not (Test-Path $configPath)) {
        Write-Host "⚙️ Downloading default configuration..."
        Invoke-WebRequest $configUrl -OutFile $configPath -UseBasicParsing
        Write-Host "Created default configuration file"
    } else {
        Write-Host "ℹ️ Keeping existing configuration file"
    }
    
    # Create or update shortcut in startup folder
    $startupPath = [System.IO.Path]::Combine([Environment]::GetFolderPath("Startup"), "WinCuts.lnk")
    $shell = New-Object -ComObject WScript.Shell
    $shortcut = $shell.CreateShortcut($startupPath)
    $shortcut.TargetPath = "$installDir\WinCuts.exe"
    $shortcut.Save()
    
    # Start WinCuts
    Start-Process -FilePath "$installDir\WinCuts.exe"
    
    Write-Host @"
    
✅ WinCuts $tag installed successfully!
   - Location: $installDir
   - Config: $configPath
   - Autostart: Enabled

🔄 Update Summary:
   - Old version stopped and removed
   - New version installed
   - Configuration preserved
   $(if (Test-Path "$configPath.backup") {"   - Backup created: $configPath.backup"})

🎮 Default shortcuts:
   - Alt + [1-9]: Switch desktop
   - Alt + Shift + [1-9]: Move window
   - Alt + Up/Down: Maximize/Minimize
   - Alt + Space: Toggle window state

🚀 WinCuts is now running!

💡 To uninstall, run in PowerShell:
   Remove-Item '$installDir','$configDir','$startupPath' -Recurse -Force
"@

} catch {
    Write-Host "❌ Installation failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Please report this issue at: https://github.com/$repo/issues" -ForegroundColor Yellow
    exit 1
} 