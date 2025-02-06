# Installation script for WinCuts
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

try {
    $ErrorActionPreference = "Stop"

    # Download latest release info from GitHub
    $repo = "fingann/WinCuts"
    $releases = "https://api.github.com/repos/$repo/releases/latest"
    
    Write-Host "🔍 Fetching latest release info..."
    $releaseInfo = Invoke-WebRequest $releases -UseBasicParsing | ConvertFrom-Json
    $tag = $releaseInfo.tag_name
    
    # Download URLs
    $exeUrl = $releaseInfo.assets | Where-Object { $_.name -like "*.zip" } | Select-Object -ExpandProperty browser_download_url
    
    if (-not $exeUrl) {
        throw "No release asset found matching *.zip"
    }
    
    Write-Host "🔍 Found latest version: $tag"
    Write-Host "📥 Download URL: $exeUrl"
    
    # Installation directories
    $installDir = "$env:LOCALAPPDATA\WinCuts"
    $configDir = "$env:APPDATA\WinCuts"
    $tempDir = "$env:TEMP\WinCuts_Install"
    
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
    New-Item -ItemType Directory -Force -Path $tempDir | Out-Null
    
    # Backup existing config if it exists
    $configPath = "$configDir\config.yaml"
    if (Test-Path $configPath) {
        Write-Host "💾 Backing up existing configuration..."
        Copy-Item $configPath "$configPath.backup"
    }
    
    # Clean up old files
    Write-Host "🧹 Cleaning up old files..."
    Remove-Item "$installDir\*" -Force -Recurse -ErrorAction SilentlyContinue
    Remove-Item "$tempDir\*" -Force -Recurse -ErrorAction SilentlyContinue
    
    # Download and extract release
    Write-Host "⬇️ Downloading latest version..."
    $zipPath = "$tempDir\WinCuts.zip"
    Invoke-WebRequest $exeUrl -OutFile $zipPath -UseBasicParsing
    
    Write-Host "📦 Extracting files..."
    Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
    
    # Verify extracted files
    $extractedFiles = Get-ChildItem -Path $tempDir -Recurse
    Write-Host "📋 Extracted files:"
    $extractedFiles | ForEach-Object { Write-Host "   - $($_.FullName)" }
    
    # Find and copy the executable
    $exeFile = Get-ChildItem -Path $tempDir -Filter "WinCuts.exe" -Recurse | Select-Object -First 1
    if (-not $exeFile) {
        throw "WinCuts.exe not found in extracted files"
    }
    
    Copy-Item $exeFile.FullName -Destination "$installDir\WinCuts.exe"
    
    # Find and handle config
    $defaultConfig = Get-ChildItem -Path $tempDir -Filter "default_config.yaml" -Recurse | Select-Object -First 1
    if (-not $defaultConfig) {
        throw "default_config.yaml not found in extracted files"
    }
    
    if (-not (Test-Path $configPath)) {
        Write-Host "⚙️ Creating default configuration..."
        Copy-Item $defaultConfig.FullName $configPath
        Write-Host "Created default configuration file"
    } else {
        Write-Host "ℹ️ Keeping existing configuration file"
    }
    
    # Clean up temp files
    Remove-Item $tempDir -Recurse -Force
    
    # Install and start the service
    Write-Host "🔧 Installing Windows service..."
    & "$installDir\WinCuts.exe" -install
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to install service"
    }
    
    Write-Host "▶️ Starting service..."
    Start-Service -Name "WinCuts"
    
    Write-Host @"
    
✅ WinCuts $tag installed successfully!
   - Location: $installDir
   - Config: $configPath
   - Service: Installed and running

🔄 Update Summary:
   - Old version stopped and removed
   - New version installed as Windows service
   - Configuration preserved
   $(if (Test-Path "$configPath.backup") {"   - Backup created: $configPath.backup"})

🎮 Default shortcuts:
   - Alt + [1-9]: Switch desktop
   - Alt + Shift + [1-9]: Move window
   - Alt + Up/Down: Maximize/Minimize
   - Alt + Space: Toggle window state

🚀 WinCuts service is now running!

💡 To uninstall, run:
   1. Stop-Service WinCuts
   2. & '$installDir\WinCuts.exe' -uninstall
   3. Remove-Item '$installDir','$configDir' -Recurse -Force
"@

} catch {
    Write-Host "❌ Installation failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Debug information:" -ForegroundColor Yellow
    Write-Host "  Temp directory contents (if exists):" -ForegroundColor Yellow
    if (Test-Path $tempDir) {
        Get-ChildItem -Path $tempDir -Recurse | ForEach-Object { Write-Host "   - $($_.FullName)" }
    } else {
        Write-Host "   Temp directory does not exist"
    }
    Write-Host "Please report this issue at: https://github.com/$repo/issues" -ForegroundColor Yellow
    exit 1
} 