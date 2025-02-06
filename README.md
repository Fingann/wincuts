# WinCuts 🪟✂️

WinCuts is a lightweight, customizable Windows utility that enhances your virtual desktop experience with powerful keyboard shortcuts and a sleek system tray interface.

## Features 🌟

- **Intuitive Keyboard Shortcuts**: Quickly switch between virtual desktops and move windows using customizable keyboard combinations
- **Smart System Tray**: Always know which desktop you're on with our minimalist, modern system tray indicator
- **Highly Configurable**: Customize everything from keyboard shortcuts to the appearance of the system tray icon
- **Resource Efficient**: Built in Go for minimal system impact
- **Windows Native**: Seamlessly integrates with Windows virtual desktops

## Installation 🚀

### Prerequisites

- Windows 10 or later
- Go 1.21 or later (for building from source)

### Quick Start

1. Download the latest release from our [Releases](https://github.com/yourusername/WinCuts/releases) page
2. Extract the zip file to your desired location
3. Run `WinCuts.exe`

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/WinCuts.git
cd WinCuts

# Build the project
go build

# Run WinCuts
./WinCuts
```

## Configuration ⚙️

WinCuts is highly configurable through a YAML configuration file. You can generate a default configuration file with:

```bash
WinCuts.exe --generate-config config.yaml
```

### Default Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Alt + [1-9] | Switch to desktop 1-9 |
| Alt + Shift + [1-9] | Move current window to desktop 1-9 |
| Alt + N | Create new desktop |

### Customizing Configuration

The configuration file (`config.yaml`) allows you to customize:

- Keyboard shortcuts
- System tray appearance
- Logging levels
- Minimum number of virtual desktops

Example configuration:

```yaml
logging:
  level: INFO
ui:
  tray_icon:
    size: 22
    corner_radius: 4
    padding: 2
    bg_opacity: 230
    bg_color:
      r: 0
      g: 120
      b: 215
      a: 255
shortcuts:
  bindings:
    - keys: ["LAlt", "1"]
      action: SwitchDesktop
      params: ["1"]
    - keys: ["LAlt", "LShift", "1"]
      action: MoveWindowToDesktop
      params: ["1"]
```

## Command Line Options 🖥️

- `--config <path>`: Load configuration from specified file
- `--generate-config <path>`: Generate default configuration file
- `--log-level <level>`: Set logging level (DEBUG, INFO, WARN, ERROR)
- `--min-desktops <number>`: Set minimum number of virtual desktops

## Contributing 🤝

We welcome contributions! Whether it's bug reports, feature requests, or code contributions, please feel free to help make WinCuts better.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Setup

1. Install Go 1.21 or later
2. Clone the repository
3. Install dependencies: `go mod download`
4. Run tests: `go test ./...`

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments 🙏

- Thanks to all contributors who have helped shape WinCuts
- Built with Go and ❤️

## Support 💪

If you find WinCuts useful, please consider:
- Starring the repository
- Reporting bugs
- Contributing code or documentation
- Sharing with others

## Roadmap 🗺️

- [ ] Additional keyboard shortcut actions
- [ ] Custom themes for system tray icon
- [ ] Window management features
- [ ] Multiple monitor support improvements
- [ ] Configuration GUI

---
