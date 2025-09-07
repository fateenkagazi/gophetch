# Gophetch

[![Go Report Card](https://goreportcard.com/badge/github.com/Cod-e-Codes/gophetch)](https://goreportcard.com/report/github.com/Cod-e-Codes/gophetch)
[![CI](https://github.com/Cod-e-Codes/gophetch/workflows/CI/badge.svg)](https://github.com/Cod-e-Codes/gophetch/actions)
[![Release](https://img.shields.io/github/v/release/Cod-e-Codes/gophetch)](https://github.com/Cod-e-Codes/gophetch/releases)

A terminal-based system monitor with ASCII animation built in Go using Bubble Tea.

<img src="gophetch-demo.gif" width="600" alt="Gophetch Demo">

*Demo showing the animated rain cloud, color palette, and real-time system information.*

## Features

- Real-time system information display
- Animated rain cloud ASCII art (default)
- Custom ASCII frame file support for personalized animations
- Cross-platform compatibility (Windows, Linux, macOS, Android/Termux)
- Color palette with animated wave effects
- System metrics including CPU, memory, disk usage, and load average
- Responsive terminal UI with proper cleanup

## Requirements

- Go 1.21 or later
- Terminal with ANSI color support

## Installation

### From Source
```bash
git clone https://github.com/Cod-e-Codes/gophetch.git
cd gophetch
go build
```

### From Releases
Download the latest release for your platform from the [Releases](https://github.com/Cod-e-Codes/gophetch/releases) page.

### Build Scripts
The project includes build scripts for different platforms:

**PowerShell (Windows):**
```powershell
.\build.ps1                    # Build for current platform
.\build.ps1 -Release           # Build optimized release
.\build.ps1 -Platform "linux/amd64"  # Cross-compile
```

**Bash (Linux/macOS):**
```bash
./build.sh                     # Build for current platform
./build.sh --release           # Build optimized release
./build.sh --platform "windows/amd64"  # Cross-compile
```

**Release Build (All Platforms):**
```bash
./release.sh 1.0.0             # Build release for all platforms
```

## Usage

```bash
# Run with default rain animation
./gophetch

# Run with custom frame rate
./gophetch 100ms

# Run with custom ASCII frames file
./gophetch frames.txt

# Run with custom frames file and frame rate
./gophetch frames.txt 500ms
```

## Frame File Format

Custom ASCII frames can be loaded from text files. Each frame is separated by `---FRAME---`:

```
┌─────────────┐
│   FRAME 1   │
│  ┌───────┐  │
│  │ ● ● ● │  │
│  └───────┘  │
└─────────────┘
---FRAME---
┌─────────────┐
│   FRAME 2   │
│  ┌───────┐  │
│  │ ○ ● ○ │  │
│  └───────┘  │
└─────────────┘
```

## Controls

- `q` or `Ctrl+C` - Exit application

## System Information Displayed

- Operating system and architecture
- Username
- CPU core count
- Memory allocation and garbage collection stats
- Disk usage and permissions
- Process count
- Load average (estimated on Windows)
- Runtime information (uptime, FPS, Go version)

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library

## License

See [LICENSE](LICENSE) file for details.
