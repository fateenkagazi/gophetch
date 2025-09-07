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
- JSON configuration file for customizable display options
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

## Configuration

Gophetch supports a JSON configuration file (`gophetch.json`) that is automatically created on first run. You can customize all display options:

```json
{
  "fps": 5,
  "color_scheme": "blue",
  "show_cpu": true,
  "show_memory": true,
  "show_disk": true,
  "show_uptime": true,
  "show_kernel": true,
  "show_os": true,
  "show_hostname": true,
  "frame_file": "default",
  "loop_animation": true,
  "center_content": true,
  "static_mode": false,
  "hide_animation": false,
  "show_fps_counter": false,
  "show_weather": false
}
```

### Configuration Options

- **fps**: Animation frame rate (default: 5)
- **color_scheme**: Main color theme (default: "blue")
- **show_cpu**: Display CPU information (default: true)
- **show_memory**: Display memory information (default: true)
- **show_disk**: Display disk usage (default: true)
- **show_uptime**: Display system uptime (default: true)
- **show_kernel**: Display kernel/Go version (default: true)
- **show_os**: Display OS and architecture (default: true)
- **show_hostname**: Display username (default: true)
- **frame_file**: Path to custom ASCII frames file (default: "default")
- **loop_animation**: Loop animation frames (default: true)
- **center_content**: Center ASCII art (default: true)
- **static_mode**: One-shot info dump like Neofetch (default: false)
- **hide_animation**: Skip animation even if frames exist (default: false)
- **show_fps_counter**: Show FPS overlay (default: false)
- **show_weather**: Display weather info (default: false)

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
