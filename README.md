# WinCuts 🪟✂️

WinCuts is a lightweight, customizable Windows utility that enhances your virtual desktop experience with powerful keyboard shortcuts and a sleek system tray interface.

## Quick Install 🚀

Open PowerShell and run:
```powershell
irm https://raw.githubusercontent.com/fingann/WinCuts/main/install.ps1 | iex
```

That's it! For detailed instructions, see [INSTALL.md](INSTALL.md).

## Features 🌟

- **Intuitive Keyboard Shortcuts**: Quickly switch between virtual desktops and move windows
- **Smart System Tray**: Always know which desktop you're on with our minimalist, modern indicator
- **Window Management**: Maximize, minimize, or toggle window states with keyboard shortcuts
- **Highly Configurable**: Customize everything from keyboard shortcuts to the appearance
- **Resource Efficient**: Built in Go for minimal system impact
- **Windows Native**: Seamlessly integrates with Windows virtual desktops

## Default Shortcuts ⌨️

- `Alt + [1-9]`: Switch to desktop 1-9
- `Alt + Shift + [1-9]`: Move current window to desktop 1-9
- `Alt + Up`: Maximize current window
- `Alt + Down`: Minimize current window
- `Alt + Space`: Toggle window state

## Configuration 🔧

Edit `%APPDATA%\WinCuts\config.yaml` to customize:
- Keyboard shortcuts
- System tray appearance
- Virtual desktop settings
- Logging levels

See [example.yaml](config/example.yaml) for all available options.

## Updating ⬆️

Run the installation command again to update to the latest version:
```powershell
irm https://raw.githubusercontent.com/fingann/WinCuts/main/install.ps1 | iex
```
Your configuration will be preserved during updates.

## Development 🛠️

### Prerequisites
- Windows 10 or later
- Go 1.21 or later

### Building from Source
```bash
git clone https://github.com/fingann/WinCuts.git
cd WinCuts
go build
```

### Running Tests
```bash
go test ./...
```

## Contributing 🤝

We welcome contributions! Whether it's bug reports, feature requests, or code contributions:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Support 💪

If you find WinCuts useful, please:
- Star the repository ⭐
- Report bugs 🐛
- Contribute code or documentation 📝
- Share with others 🌟

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments 🙏

- Thanks to all contributors who have helped shape WinCuts
- Built with Go and ❤️

## Roadmap 🗺️

- [ ] Additional keyboard shortcut actions
- [ ] Custom themes for system tray icon
- [ ] Window management features
- [ ] Multiple monitor support improvements
- [ ] Configuration GUI

---
