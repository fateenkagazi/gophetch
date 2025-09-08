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
- **Tab System**: Multi-tab interface with specialized information panels
- **Interactive Navigation**: Keyboard-driven tab switching and list navigation

## Tab System

Gophetch features a comprehensive tab system that organizes system information into specialized panels for better usability and information density.

### Available Tabs

1. **Standard Tab** - Traditional system information display
   - Operating system and architecture
   - CPU, memory, and disk usage
   - System uptime and process count
   - Load average and runtime information

2. **Network Tab** - Network connectivity and activity
   - IP addresses (all interfaces)
   - Active network connections count
   - Listening ports and services
   - Network interface statistics

3. **Hardware Tab** - Hardware-specific information
   - GPU information and driver details
   - Temperature monitoring (where available)
   - Fan speed readings (where available)
   - Battery status and level (laptops/mobile devices)

4. **Processes Tab** - Interactive process management
   - Total process count
   - Top processes by resource usage
   - Interactive list with up/down navigation
   - Process details (PID, CPU%, Memory usage)

5. **Weather Tab** - Current weather and forecast
   - Current weather conditions and temperature
   - Today's detailed forecast with ASCII art
   - Auto-detected location or manual location support
   - Real-time weather data from wttr.in

### Tab Navigation

- **Tab/Shift+Tab** - Navigate between tabs
- **Number keys (1-5)** - Jump directly to specific tabs
- **Up/Down arrows or j/k** - Navigate within the Processes tab list
- **q or Ctrl+C** - Exit application

### Tab Configuration

The tab system can be configured in `gophetch.json`:

```json
{
  "enable_tabs": true,
  "visible_tabs": ["standard", "network", "hardware", "processes", "weather"],
  "default_tab": "standard",
  "tab_order": ["standard", "network", "hardware", "processes", "weather"]
}
```

- **enable_tabs**: Enable/disable the tab system (default: true)
- **visible_tabs**: Array of tabs to show (can hide specific tabs)
- **default_tab**: Which tab to show on startup
- **tab_order**: Custom order for tab display

### Performance Features

- **Background Data Fetching**: All system calls run in background goroutines to keep UI responsive
- **Non-blocking Updates**: UI renders immediately while data loads progressively in background
- **Smart Caching**: Data is cached and updated at optimal intervals
- **Efficient Updates**: Network and hardware data updates every 10 seconds
- **Weather Caching**: Weather data updates every 30 seconds to respect API limits
- **Loading Indicators**: Visual feedback shows "(Loading...)" while data is being fetched
- **Responsive UI**: Smooth navigation and real-time updates without blocking

## Requirements

- Terminal with ANSI color support
- Go 1.21 or later (only required for building from source)

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

An example frame file (`example-frames.txt`) is included with a simple character animation. You can use it as a starting point for your own custom animations.

## Asciinema .cast File Support

Gophetch now supports asciinema `.cast` files, allowing you to use terminal recordings as animations. This opens up a world of possibilities for creating dynamic, realistic terminal animations.

### Creating .cast Files

You can create `.cast` files using the `asciinema` tool:

```bash
# Record a terminal session
asciinema rec my-animation.cast

# Play it back to test
asciinema play my-animation.cast

# Use it with Gophetch
./gophetch my-animation.cast
```

### .cast File Features

- **Automatic frame extraction**: Gophetch automatically extracts frames from the continuous terminal output
- **ANSI sequence processing**: Handles terminal colors, cursor movements, and formatting using regex-based processing
- **Timing preservation**: Maintains the original timing relationships from the recording
- **Standard format**: Works with any asciinema-compatible recording
- **Cross-platform**: Works on all supported platforms (Windows, Linux, macOS, Android/Termux)

### Supported File Formats

Gophetch now supports two animation formats:

1. **Custom Frame Files** (`.txt`, `.frames`): Traditional ASCII art frames separated by `---FRAME---`
2. **Asciinema .cast Files** (`.cast`): Terminal recordings with automatic frame extraction

### Example Use Cases

- Terminal demos and tutorials
- Command-line tool showcases
- Interactive script demonstrations
- Real-time data visualization
- Terminal-based games and applications
- Live coding sessions
- System monitoring displays

## Controls

### Basic Controls
- `q` or `Ctrl+C` - Exit application

### Tab System Controls
- `Tab` - Switch to next tab
- `Shift+Tab` - Switch to previous tab
- `1-5` - Jump directly to specific tab (Standard, Network, Hardware, Processes, Weather)
- `Up/Down arrows` or `j/k` - Navigate within the Processes tab list

## System Information Displayed

### Standard Tab
- Operating system and architecture
- Username
- CPU core count and usage
- Memory allocation and garbage collection stats
- Disk usage and permissions
- Process count
- Load average (estimated on Windows)
- Runtime information (uptime, FPS, Go version)

### Network Tab
- IP addresses (all network interfaces)
- Active network connections count
- Listening ports and services
- Network interface statistics

### Hardware Tab
- GPU information and driver details
- Temperature monitoring (where available)
- Fan speed readings (where available)
- Battery status and level (laptops/mobile devices)

### Processes Tab
- Total process count
- Top processes by resource usage
- Interactive process list with navigation
- Process details (PID, CPU%, Memory usage)

### Weather Tab
- Current weather conditions and temperature
- Today's detailed forecast with ASCII art
- Auto-detected location
- Real-time weather data from wttr.in

## Future Enhancements

- Change default config location to OS appropriate locations (e.g., `~/.config/gophetch/` on Linux/macOS, `%APPDATA%\gophetch\` on Windows)
- Enhanced ANSI sequence processing for better .cast file rendering
- Support for custom frame extraction intervals in .cast files

## License

See [LICENSE](LICENSE) file for details.
