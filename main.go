package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Frame represents a single ASCII animation frame
type Frame struct {
	Content string
	Color   lipgloss.Color
}

// CastHeader represents the header of an asciinema .cast file
type CastHeader struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// CastEvent represents a single event in an asciinema .cast file
type CastEvent struct {
	Timestamp float64 `json:"timestamp"`
	EventType string  `json:"event_type"`
	Data      string  `json:"data"`
}

// SystemInfo holds system information to display
type SystemInfo struct {
	OS           string
	Architecture string
	CPUCount     int
	GoVersion    string
	Memory       string
	Uptime       time.Duration
	DiskUsage    string
	Processes    int
	LoadAvg      string
	Username     string
	Weather      string
}

// Config holds all configuration options
type Config struct {
	// Display settings
	FPS          int    `json:"fps"`
	ColorScheme  string `json:"color_scheme"`
	ShowCPU      bool   `json:"show_cpu"`
	ShowMemory   bool   `json:"show_memory"`
	ShowDisk     bool   `json:"show_disk"`
	ShowUptime   bool   `json:"show_uptime"`
	ShowKernel   bool   `json:"show_kernel"`
	ShowOS       bool   `json:"show_os"`
	ShowHostname bool   `json:"show_hostname"`

	// Frame / animation settings
	FrameFile     string `json:"frame_file"`
	LoopAnimation bool   `json:"loop_animation"`

	// Output mode
	StaticMode    bool `json:"static_mode"`
	HideAnimation bool `json:"hide_animation"`

	// Misc
	ShowFPSCounter bool `json:"show_fps_counter"`
	ShowWeather    bool `json:"show_weather"`
}

// Model represents the Bubble Tea model
type Model struct {
	frames       []Frame
	currentFrame int
	frameRate    time.Duration
	startTime    time.Time
	sysInfo      SystemInfo
	config       Config
	ctx          context.Context
	cancel       context.CancelFunc
	width        int
	height       int
	mutex        *sync.RWMutex
}

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	containerStyle = lipgloss.NewStyle().
			Padding(1)
)

// Rain animation constants and variables
const (
	white = "\033[97m"
	blue  = "\033[34m"
	reset = "\033[0m"
)

var cloud = []string{
	"  (   ).  ",
	" (___(__) ",
}

var rainChars = []rune{'\'', '`', '|', '.', '˙'}

// Compiled regex patterns for better performance
var (
	// ANSI escape sequence patterns
	ansiColorRegex       = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	ansiCursorRegex      = regexp.MustCompile(`\x1b\[[0-9]*[ABCDFGHK]`)
	ansiClearRegex       = regexp.MustCompile(`\x1b\[[0-9]*[JK]`)
	ansiComplexRegex     = regexp.MustCompile(`\x1b\[[?0-9;]*[hlnpqr]`)
	ansiOSCRegex         = regexp.MustCompile(`\x1b\][0-9]*;[^\x07]*\x07`)
	ansiPrivateRegex     = regexp.MustCompile(`\x1b\[[?0-9;]*[a-zA-Z]`)
	ansiDeviceRegex      = regexp.MustCompile(`\x1b\[[0-9]*n`)
	ansiApplicationRegex = regexp.MustCompile(`\x1b\[[?0-9;]*[hl]`)
)

// generateCloudWithRain creates a single cloud with rain animation
func generateCloudWithRain(animated bool) []string {
	lines := make([]string, 8) // Extended to 8 lines for better height matching

	// cloud lines
	for i, c := range cloud {
		lines[i] = white + c + reset
	}

	width := len(cloud[0])

	// rain lines
	for i := 2; i < 8; i++ {
		line := ""
		for j := 0; j < width; j++ {
			if animated && rand.Float64() < 0.6 {
				line += blue + string(rainChars[rand.Intn(len(rainChars))]) + reset
			} else if !animated {
				// Static rain - use a fixed pattern
				if (i+j)%3 == 0 {
					line += blue + string(rainChars[0]) + reset
				} else {
					line += " "
				}
			} else {
				line += " "
			}
		}
		lines[i] = line
	}

	return lines
}

// generateColorPalette creates neofetch-style color squares
func generateColorPalette(startTime time.Time) string {
	var palette strings.Builder

	// Calculate animation phase for color cycling
	elapsed := time.Since(startTime)
	phase := math.Sin(elapsed.Seconds()*1.5)*0.5 + 0.5 // Oscillate between 0 and 1

	// Create a rainbow wave effect that cycles through colors
	// First row - create a wave of colors that moves across
	// Use only visible colors (avoid black/0 and white/7)
	baseColors := []lipgloss.Color{"1", "2", "3", "4", "5", "6", "8", "9"}
	for i := 0; i < 8; i++ {
		// Calculate wave position for this square
		wavePos := (float64(i)/8.0 + phase) * 2 * math.Pi
		waveIntensity := (math.Sin(wavePos) + 1) / 2 // 0 to 1

		// Cycle through different color sets based on wave intensity
		var animatedColor lipgloss.Color
		if waveIntensity < 0.33 {
			// Low intensity - use dark but visible colors
			animatedColor = baseColors[i]
		} else if waveIntensity < 0.66 {
			// Medium intensity - use bright colors
			animatedColor = lipgloss.Color(fmt.Sprintf("%d", 8+i))
		} else {
			// High intensity - use vibrant colors
			vibrantColors := []lipgloss.Color{"10", "11", "12", "13", "14", "15", "9", "10"}
			animatedColor = vibrantColors[i]
		}

		palette.WriteString(lipgloss.NewStyle().
			Background(animatedColor).
			Render("   "))
	}
	palette.WriteString("\n")

	// Second row - create a different wave pattern with different timing
	brightColors := []lipgloss.Color{"8", "9", "10", "11", "12", "13", "14", "15"}
	for i := 0; i < 8; i++ {
		// Use a different wave calculation - reverse direction and different frequency
		wavePos := (float64(7-i)/8.0 + phase*1.3 + 0.7) * 2 * math.Pi
		waveIntensity := (math.Sin(wavePos) + 1) / 2 // 0 to 1

		// Cycle through different color sets based on wave intensity
		var animatedColor lipgloss.Color
		if waveIntensity < 0.33 {
			// Low intensity - use bright colors
			animatedColor = brightColors[i]
		} else if waveIntensity < 0.66 {
			// Medium intensity - use vibrant colors (different set from first row)
			vibrantColors := []lipgloss.Color{"11", "12", "13", "14", "15", "9", "10", "11"}
			animatedColor = vibrantColors[i]
		} else {
			// High intensity - use base colors (different set from first row)
			baseColors2 := []lipgloss.Color{"2", "3", "4", "5", "6", "8", "9", "1"}
			animatedColor = baseColors2[i]
		}

		palette.WriteString(lipgloss.NewStyle().
			Background(animatedColor).
			Render("   "))
	}

	// Create a properly bordered palette using lipgloss
	paletteContent := palette.String()

	// Use lipgloss to create a bordered box
	borderedPalette := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Render(paletteContent)

	return borderedPalette
}

// Messages for Bubble Tea
type tickMsg time.Time
type sysInfoMsg SystemInfo

// Commands
func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func updateSysInfo() tea.Cmd {
	return func() tea.Msg {
		return sysInfoMsg(GetSystemInfo())
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickEvery(m.frameRate),
		updateSysInfo(),
		tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return sysInfoMsg(GetSystemInfo())
		}),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		}

	case tickMsg:
		m.mutex.Lock()
		// Only cycle through frames if we have frames loaded
		// For rain animation, we don't need to cycle frames
		if len(m.frames) > 0 {
			if m.config.LoopAnimation {
				m.currentFrame = (m.currentFrame + 1) % len(m.frames)
			} else {
				// Don't loop - stay on last frame
				if m.currentFrame < len(m.frames)-1 {
					m.currentFrame++
				}
			}
		}
		m.mutex.Unlock()

		// In static mode, don't continue ticking after first update
		if m.config.StaticMode {
			return m, nil
		}
		return m, tickEvery(m.frameRate)

	case sysInfoMsg:
		m.mutex.Lock()
		m.sysInfo = SystemInfo(msg)
		m.mutex.Unlock()
		// Only schedule next system info update if not in static mode
		if !m.config.StaticMode {
			return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return sysInfoMsg(GetSystemInfo())
			})
		}
		return m, nil
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create system info panel
	sysInfoPanel := m.renderSystemInfo()

	var content string
	if m.config.HideAnimation {
		// Hide animation - just show system info
		content = lipgloss.NewStyle().
			Width(m.width).
			Padding(1, 2).
			Render(sysInfoPanel)
	} else {
		var animationPanel string
		var actualAnimationWidth int

		if len(m.frames) > 0 {
			// Use custom frames
			if m.currentFrame < len(m.frames) {
				animationPanel = m.frames[m.currentFrame].Content
			} else {
				animationPanel = m.frames[0].Content // Fallback to first frame
			}
			// Calculate width from the actual frame content
			lines := strings.Split(animationPanel, "\n")
			if len(lines) > 0 {
				actualAnimationWidth = len(lines[0]) + 4 // +4 for padding
			} else {
				actualAnimationWidth = 20 // Default width
			}
		} else {
			// Generate rain animation
			rainAnimation := generateCloudWithRain(!m.config.StaticMode)
			animationPanel = strings.Join(rainAnimation, "\n")
			actualAnimationWidth = len(rainAnimation[0]) + 4 // +4 for padding
		}

		infoWidth := m.width - actualAnimationWidth - 3 // -3 for spacing

		// Style the panels without borders
		animationStyled := lipgloss.NewStyle().
			Width(actualAnimationWidth).
			Height(10). // Reduced to match the 8-line animation + padding
			Padding(1, 2).
			Render(animationPanel)

		infoStyled := lipgloss.NewStyle().
			Width(infoWidth).
			Padding(1, 2).
			Render(sysInfoPanel)

		// Join panels side by side with proper spacing
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			animationStyled,
			lipgloss.NewStyle().Width(3).Render(""), // 3-space gap
			infoStyled,
		)
	}

	// Generate color palette with animation (will update on each tick)
	var colorPalette string
	if m.config.StaticMode {
		// In static mode, use a static color palette
		colorPalette = generateStaticColorPalette()
	} else {
		colorPalette = generateColorPalette(m.startTime)
	}

	// Add title and controls
	title := titleStyle.Render("Gophetch - System Monitor")
	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Press 'q' or Ctrl+C to quit")

	// Add just a small amount of top padding for better visual balance
	// without cutting off the top
	topPadding := "\n"

	// Combine everything with minimal vertical spacing
	return containerStyle.Render(
		topPadding +
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				content,
				"",
				colorPalette,
				"",
				controls,
			),
	)
}

// renderSystemInfo creates the system information display
func (m Model) renderSystemInfo() string {
	var info strings.Builder

	// System information
	info.WriteString(infoStyle.Render("System Information"))
	info.WriteString("\n")
	info.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("─────────────────────"))
	info.WriteString("\n\n")

	if m.config.ShowOS {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("OS:"),
			valueStyle.Render(fmt.Sprintf("%s (%s)", m.sysInfo.OS, m.sysInfo.Architecture))))
	}

	if m.config.ShowHostname {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("User:"),
			valueStyle.Render(m.sysInfo.Username)))
	}

	if m.config.ShowCPU {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("CPU:"),
			valueStyle.Render(fmt.Sprintf("%d cores", m.sysInfo.CPUCount))))
	}

	if m.config.ShowMemory {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("Memory:"),
			valueStyle.Render(m.sysInfo.Memory)))
	}

	if m.config.ShowKernel {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("Go Version:"),
			valueStyle.Render(m.sysInfo.GoVersion)))
	}

	if m.sysInfo.Processes > 0 {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("Processes:"),
			valueStyle.Render(fmt.Sprintf("%d", m.sysInfo.Processes))))
	}

	if m.sysInfo.LoadAvg != "Load: N/A" && m.sysInfo.LoadAvg != "" {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("Load:"),
			valueStyle.Render(strings.TrimPrefix(m.sysInfo.LoadAvg, "Load: "))))
	}

	if m.config.ShowDisk {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("Disk:"),
			valueStyle.Render(strings.TrimPrefix(m.sysInfo.DiskUsage, "Disk: "))))
	}

	// Runtime information
	info.WriteString("\n")
	info.WriteString(infoStyle.Render("Runtime Information"))
	info.WriteString("\n")
	info.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("─────────────────────"))
	info.WriteString("\n\n")

	if m.config.ShowUptime {
		uptime := time.Since(m.startTime)
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("Uptime:"),
			valueStyle.Render(uptime.Truncate(time.Second).String())))
	}

	info.WriteString(fmt.Sprintf("%s %s\n",
		infoStyle.Render("Time:"),
		valueStyle.Render(time.Now().Format("15:04:05"))))

	if m.config.ShowWeather {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("Weather:"),
			valueStyle.Render(strings.TrimPrefix(m.sysInfo.Weather, "Weather: "))))
	}

	if m.config.ShowFPSCounter {
		info.WriteString(fmt.Sprintf("%s %s\n",
			infoStyle.Render("FPS:"),
			valueStyle.Render(fmt.Sprintf("%.1f", float64(time.Second)/float64(m.frameRate)))))
	}

	return info.String()
}

// getUsername gets the current username (cross-platform)
func getUsername() string {
	// Check if we're in Termux environment
	if isTermux() {
		return getTermuxUsername()
	}

	// Try various environment variables for different systems
	envVars := []string{"USER", "USERNAME", "LOGNAME"}

	for _, envVar := range envVars {
		if username := os.Getenv(envVar); username != "" {
			// Skip Termux installation directory
			if !strings.Contains(username, "/data/data/com.termux") {
				return username
			}
		}
	}

	// Try to get user from whoami command as fallback
	if runtime.GOOS == "linux" || runtime.GOOS == "android" {
		if output, err := exec.Command("whoami").Output(); err == nil {
			username := strings.TrimSpace(string(output))
			if username != "" && !strings.Contains(username, "/data/data/com.termux") {
				return username
			}
		}

		// Try id -un as alternative
		if output, err := exec.Command("id", "-un").Output(); err == nil {
			username := strings.TrimSpace(string(output))
			if username != "" && !strings.Contains(username, "/data/data/com.termux") {
				return username
			}
		}
	}

	return "Unknown"
}

// isTermux detects if we're running in Termux environment
func isTermux() bool {
	// Check for Termux-specific environment variables or paths
	termuxIndicators := []string{
		"PREFIX",         // Termux sets this to /data/data/com.termux/files/usr
		"ANDROID_ROOT",   // Android system indicator
		"TERMUX_VERSION", // Termux version variable
	}

	for _, indicator := range termuxIndicators {
		if value := os.Getenv(indicator); value != "" {
			if strings.Contains(value, "termux") || strings.Contains(value, "android") {
				return true
			}
		}
	}

	// Check if PWD contains termux path
	if pwd := os.Getenv("PWD"); pwd != "" {
		if strings.Contains(pwd, "/data/data/com.termux") {
			return true
		}
	}

	return false
}

// getTermuxUsername gets the username specifically for Termux
func getTermuxUsername() string {
	// First try USER environment variable
	if user := os.Getenv("USER"); user != "" && !strings.Contains(user, "/data/data/com.termux") {
		return user
	}

	// Try whoami command
	if output, err := exec.Command("whoami").Output(); err == nil {
		username := strings.TrimSpace(string(output))
		if username != "" && !strings.Contains(username, "/data/data/com.termux") {
			return username
		}
	}

	// Try id -un command
	if output, err := exec.Command("id", "-un").Output(); err == nil {
		username := strings.TrimSpace(string(output))
		if username != "" && !strings.Contains(username, "/data/data/com.termux") {
			return username
		}
	}

	// Try to extract from PREFIX path as last resort
	if prefix := os.Getenv("PREFIX"); prefix != "" {
		// PREFIX is usually /data/data/com.termux/files/usr
		// We can't get the actual username from this, so return a generic termux user
		return "termux"
	}

	return "termux"
}

// GetSystemInfo gathers comprehensive system information
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		CPUCount:     runtime.NumCPU(),
		GoVersion:    runtime.Version(),
		Username:     getUsername(),
	}

	// Get memory information
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	info.Memory = fmt.Sprintf("Alloc: %d MB, Sys: %d MB, GC: %d",
		bToMb(m.Alloc), bToMb(m.Sys), m.NumGC)

	// Get disk usage
	info.DiskUsage = getDiskUsage()

	// Get process count
	info.Processes = getProcessCount()

	// Get load average (Unix-like systems)
	info.LoadAvg = getLoadAverage()

	// Get weather information
	info.Weather = getWeather()

	return info
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// getDiskUsage gets actual disk usage information (cross-platform)
func getDiskUsage() string {
	switch runtime.GOOS {
	case "linux", "darwin":
		return getUnixDiskUsage()
	case "android":
		return getAndroidDiskUsage()
	case "windows":
		return getWindowsDiskUsage()
	default:
		return "N/A"
	}
}

// getUnixDiskUsage gets disk usage on Unix-like systems
func getUnixDiskUsage() string {
	if runtime.GOOS == "linux" {
		if usage := getLinuxDiskUsageFromProc(); usage != "" {
			return usage
		}
	}
	return "Unix filesystem accessible"
}

// getLinuxDiskUsageFromProc reads filesystem info from /proc/mounts
func getLinuxDiskUsageFromProc() string {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[1] == "/" {
			filesystem := fields[0]
			fstype := fields[2]
			return fmt.Sprintf("%s (%s)", filesystem, fstype)
		}
	}

	return "Linux filesystem accessible"
}

// getAndroidDiskUsage gets disk usage on Android/Termux
func getAndroidDiskUsage() string {
	pwd, err := os.Getwd()
	if err != nil {
		return "Cannot access current directory"
	}

	// Try to get basic directory info
	info, err := os.Stat(pwd)
	if err != nil {
		return "Directory not accessible"
	}

	// Check if we can read and write
	readable := true
	writable := true

	// Try to create a temporary file to test write permissions
	tempFile := pwd + "/.gophetch_test"
	if f, err := os.Create(tempFile); err != nil {
		writable = false
	} else {
		f.Close()
		os.Remove(tempFile) // Clean up
	}

	// Try to read directory to test read permissions
	if _, err := os.ReadDir(pwd); err != nil {
		readable = false
	}

	// Format permissions
	perms := ""
	if readable && writable {
		perms = " (R/W)"
	} else if readable {
		perms = " (R)"
	} else if writable {
		perms = " (W)"
	} else {
		perms = " (No access)"
	}

	// Get directory name for display
	dirName := "Termux"
	if info.IsDir() {
		dirName = "Android"
	}

	return fmt.Sprintf("%s filesystem%s", dirName, perms)
}

// getWindowsDiskUsage gets disk usage on Windows
func getWindowsDiskUsage() string {
	pwd, err := os.Getwd()
	if err != nil {
		return "Cannot access current directory"
	}

	if len(pwd) >= 2 && pwd[1] == ':' {
		drive := pwd[:2]

		// Test permissions without needing file info

		// Check if we can read and write
		readable := true
		writable := true

		// Try to create a temporary file to test write permissions
		tempFile := pwd + "/.gophetch_test"
		if f, err := os.Create(tempFile); err != nil {
			writable = false
		} else {
			f.Close()
			os.Remove(tempFile) // Clean up
		}

		// Try to read directory to test read permissions
		if _, err := os.ReadDir(pwd); err != nil {
			readable = false
		}

		// Format permissions
		perms := ""
		if readable && writable {
			perms = " (R/W)"
		} else if readable {
			perms = " (R)"
		} else if writable {
			perms = " (W)"
		} else {
			perms = " (No access)"
		}

		return fmt.Sprintf("Drive %s%s", drive, perms)
	}

	return "Windows filesystem accessible"
}

// getProcessCount attempts to get the number of running processes
func getProcessCount() int {
	switch runtime.GOOS {
	case "linux":
		return getLinuxProcessCount()
	case "android":
		return getAndroidProcessCount()
	case "darwin":
		return getDarwinProcessCount()
	case "windows":
		return getWindowsProcessCount()
	default:
		return -1
	}
}

// getAndroidProcessCount gets process count on Android/Termux
func getAndroidProcessCount() int {
	// Try to use ps command as fallback
	if output, err := exec.Command("ps", "-A").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		// Subtract 1 for the header line
		return len(lines) - 1
	}

	// Fallback to CPU-based estimate
	return runtime.NumCPU() * 30 // Conservative estimate for mobile
}

// getLinuxProcessCount gets process count on Linux from /proc
func getLinuxProcessCount() int {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return -1
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			if _, err := strconv.Atoi(entry.Name()); err == nil {
				count++
			}
		}
	}
	return count
}

// getDarwinProcessCount gets process count estimate on macOS
func getDarwinProcessCount() int {
	return runtime.NumCPU() * 50
}

// getWindowsProcessCount gets process count estimate on Windows
func getWindowsProcessCount() int {
	return runtime.NumCPU() * 40
}

// getLoadAverage gets system load average (cross-platform)
func getLoadAverage() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxLoadAverage()
	case "android":
		return getAndroidLoadAverage()
	case "darwin":
		return "macOS - use Activity Monitor"
	case "windows":
		return getWindowsLoadAverage()
	default:
		return "N/A"
	}
}

// getAndroidLoadAverage calculates a simple load estimate for Android/Termux
func getAndroidLoadAverage() string {
	// For Android/Termux, we'll use a simple CPU usage estimate
	// This is a basic approximation since Android doesn't have traditional load averages
	cpuCount := runtime.NumCPU()

	// Get memory stats as a proxy for system load
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate a simple load estimate based on memory usage and GC activity
	memUsagePercent := float64(m.Alloc) / float64(m.Sys) * 100
	gcLoad := float64(m.NumGC) / 100.0 // Normalize GC count

	// Combine into a simple load estimate (0.0 to cpuCount*2.0)
	estimatedLoad := (memUsagePercent/100.0 + gcLoad) * float64(cpuCount)
	if estimatedLoad > float64(cpuCount)*2.0 {
		estimatedLoad = float64(cpuCount) * 2.0
	}

	return fmt.Sprintf("%.2f (est)", estimatedLoad)
}

// getWindowsLoadAverage calculates a simple load estimate for Windows
func getWindowsLoadAverage() string {
	// For Windows, we'll use a simple CPU usage estimate
	// This is a basic approximation since Windows doesn't have traditional load averages
	cpuCount := runtime.NumCPU()

	// Get memory stats as a proxy for system load
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate a simple load estimate based on memory usage and GC activity
	memUsagePercent := float64(m.Alloc) / float64(m.Sys) * 100
	gcLoad := float64(m.NumGC) / 100.0 // Normalize GC count

	// Combine into a simple load estimate (0.0 to cpuCount*2.0)
	estimatedLoad := (memUsagePercent/100.0 + gcLoad) * float64(cpuCount)
	if estimatedLoad > float64(cpuCount)*2.0 {
		estimatedLoad = float64(cpuCount) * 2.0
	}

	return fmt.Sprintf("%.2f (est)", estimatedLoad)
}

// getLinuxLoadAverage reads load average from /proc/loadavg
func getLinuxLoadAverage() string {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "Error reading"
	}

	fields := strings.Fields(string(data))
	if len(fields) >= 3 {
		return fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2])
	}
	return "Error parsing"
}

// LoadFramesFromFile loads ASCII frames from a file
func LoadFramesFromFile(filename string) ([]Frame, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot access file %s: %w", filename, err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a file", filename)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("file %s is empty", filename)
	}

	if fileInfo.Size() > 50*1024*1024 { // 50MB limit
		return nil, fmt.Errorf("file %s is too large (>50MB)", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", filename, err)
	}
	defer file.Close()

	var frames []Frame
	var currentFrame strings.Builder
	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if line == "---FRAME---" {
			if currentFrame.Len() > 0 {
				frames = append(frames, Frame{
					Content: currentFrame.String(),
					Color:   lipgloss.Color("252"), // Default color
				})
				currentFrame.Reset()
			}
		} else if strings.HasPrefix(line, "---FRAME:") {
			if currentFrame.Len() > 0 {
				color := extractColor(line)
				frames = append(frames, Frame{
					Content: currentFrame.String(),
					Color:   color,
				})
				currentFrame.Reset()
			}
		} else {
			currentFrame.WriteString(line + "\n")
		}

		if lineCount > 100000 {
			return nil, fmt.Errorf("file %s has too many lines (>100,000)", filename)
		}
	}

	if currentFrame.Len() > 0 {
		frames = append(frames, Frame{
			Content: currentFrame.String(),
			Color:   lipgloss.Color("252"),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no frames found in file %s", filename)
	}

	if len(frames) > 10000 {
		return nil, fmt.Errorf("too many frames in file %s (%d > 10,000)", filename, len(frames))
	}

	return frames, nil
}

// extractColor extracts color from frame delimiter line
func extractColor(line string) lipgloss.Color {
	if strings.Contains(line, ":") && strings.Contains(line, "---") {
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			colorName := strings.TrimSuffix(parts[1], "---")
			switch strings.ToUpper(colorName) {
			case "RED":
				return lipgloss.Color("196")
			case "GREEN":
				return lipgloss.Color("82")
			case "BLUE":
				return lipgloss.Color("39")
			case "YELLOW":
				return lipgloss.Color("226")
			case "CYAN":
				return lipgloss.Color("86")
			case "MAGENTA":
				return lipgloss.Color("213")
			case "WHITE":
				return lipgloss.Color("252")
			case "BRIGHTBLUE":
				return lipgloss.Color("75")
			case "BRIGHTGREEN":
				return lipgloss.Color("118")
			case "BRIGHTRED":
				return lipgloss.Color("203")
			}
		}
	}
	return lipgloss.Color("252")
}

// LoadFramesFromCastFile loads ASCII frames from an asciinema .cast file
func LoadFramesFromCastFile(filename string) ([]Frame, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot access file %s: %w", filename, err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a file", filename)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("file %s is empty", filename)
	}

	if fileInfo.Size() > 50*1024*1024 { // 50MB limit
		return nil, fmt.Errorf("file %s is too large (>50MB)", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", filename, err)
	}
	defer file.Close()

	// Read the first line to get the header
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("file %s appears to be empty or invalid", filename)
	}

	headerLine := scanner.Text()
	var header CastHeader
	if err := json.Unmarshal([]byte(headerLine), &header); err != nil {
		return nil, fmt.Errorf("invalid .cast file header: %w", err)
	}

	// Parse events and extract frames
	var frames []Frame
	var currentContent strings.Builder
	var lastTimestamp float64
	frameInterval := 0.1 // Extract frame every 100ms by default

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Parse event line as JSON array
		var eventArray []interface{}
		if err := json.Unmarshal([]byte(line), &eventArray); err != nil {
			continue // Skip invalid lines
		}

		// Check if it's a valid event array [timestamp, eventType, data]
		if len(eventArray) != 3 {
			continue // Skip invalid event arrays
		}

		// Extract event data
		timestamp, ok1 := eventArray[0].(float64)
		eventType, ok2 := eventArray[1].(string)
		data, ok3 := eventArray[2].(string)

		if !ok1 || !ok2 || !ok3 {
			continue // Skip if we can't extract the data properly
		}

		// Only process output events
		if eventType != "o" {
			continue
		}

		// Accumulate content first
		currentContent.WriteString(data)

		// Check if we should create a new frame based on time interval
		if timestamp-lastTimestamp >= frameInterval {
			if currentContent.Len() > 0 {
				// Process ANSI escape sequences and create frame
				processedContent := processANSISequences(currentContent.String())

				// Accept frames with meaningful content
				if len(strings.TrimSpace(processedContent)) > 5 {
					frames = append(frames, Frame{
						Content: processedContent,
						Color:   lipgloss.Color("252"), // Default color
					})
				}
				currentContent.Reset()
			}
			lastTimestamp = timestamp
		}

		if lineCount > 100000 {
			return nil, fmt.Errorf("file %s has too many lines (>100,000)", filename)
		}
	}

	// Add the last frame if there's content
	if currentContent.Len() > 0 {
		processedContent := processANSISequences(currentContent.String())
		if len(strings.TrimSpace(processedContent)) > 5 {
			frames = append(frames, Frame{
				Content: processedContent,
				Color:   lipgloss.Color("252"),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no frames found in .cast file %s", filename)
	}

	if len(frames) > 10000 {
		return nil, fmt.Errorf("too many frames in .cast file %s (%d > 10,000)", filename, len(frames))
	}

	return frames, nil
}

// processANSISequences processes ANSI escape sequences and returns clean text
func processANSISequences(input string) string {
	// Use regex patterns for efficient ANSI sequence removal
	result := input

	// Remove all ANSI escape sequences using regex patterns
	result = ansiColorRegex.ReplaceAllString(result, "")       // Color codes
	result = ansiCursorRegex.ReplaceAllString(result, "")      // Cursor movement
	result = ansiClearRegex.ReplaceAllString(result, "")       // Clear screen/line
	result = ansiComplexRegex.ReplaceAllString(result, "")     // Complex sequences
	result = ansiOSCRegex.ReplaceAllString(result, "")         // Operating System Command
	result = ansiPrivateRegex.ReplaceAllString(result, "")     // Private sequences
	result = ansiDeviceRegex.ReplaceAllString(result, "")      // Device control
	result = ansiApplicationRegex.ReplaceAllString(result, "") // Application sequences

	// Remove any remaining escape sequences that might have been missed
	result = strings.ReplaceAll(result, "\u001b[", "")

	// Remove bell character
	result = strings.ReplaceAll(result, "\u0007", "")

	return result
}

// generateStaticColorPalette creates a static color palette for static mode
func generateStaticColorPalette() string {
	var palette strings.Builder

	// First row - static colors
	colors1 := []lipgloss.Color{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, color := range colors1 {
		palette.WriteString(lipgloss.NewStyle().
			Background(color).
			Render("   "))
	}

	palette.WriteString("\n")

	// Second row - static colors
	colors2 := []lipgloss.Color{"9", "10", "11", "12", "13", "14", "15", "16"}
	for _, color := range colors2 {
		palette.WriteString(lipgloss.NewStyle().
			Background(color).
			Render("   "))
	}

	// Create a properly bordered palette using lipgloss
	paletteContent := palette.String()

	// Use lipgloss to create a bordered box
	borderedPalette := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Render(paletteContent)

	return borderedPalette
}

// getWeather gets weather information from wttr.in
func getWeather() string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://wttr.in/?format=%C+%t")
	if err != nil {
		return "Weather: N/A"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Weather: N/A"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Weather: N/A"
	}

	weather := strings.TrimSpace(string(body))
	if weather == "" {
		return "Weather: N/A"
	}

	return fmt.Sprintf("Weather: %s", weather)
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() Config {
	return Config{
		// Display settings
		FPS:          5,
		ColorScheme:  "blue",
		ShowCPU:      true,
		ShowMemory:   true,
		ShowDisk:     true,
		ShowUptime:   true,
		ShowKernel:   true,
		ShowOS:       true,
		ShowHostname: true,

		// Frame / animation settings
		FrameFile:     "default",
		LoopAnimation: true,

		// Output mode
		StaticMode:    false,
		HideAnimation: false,

		// Misc
		ShowFPSCounter: false,
		ShowWeather:    false,
	}
}

// loadConfig loads configuration from file or creates default
func loadConfig() (Config, error) {
	configPath := "gophetch.json"

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		defaultConfig := getDefaultConfig()
		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return defaultConfig, fmt.Errorf("failed to marshal default config: %v", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return defaultConfig, fmt.Errorf("failed to write default config: %v", err)
		}

		fmt.Printf("Created default config file: %s\n", configPath)
		return defaultConfig, nil
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return getDefaultConfig(), fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return getDefaultConfig(), fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Warning: %v, using defaults\n", err)
		config = getDefaultConfig()
	}

	var frames []Frame
	frameRate := time.Duration(1000/config.FPS) * time.Millisecond

	// Load frames based on config or command line arguments
	if len(os.Args) > 1 {
		// Command line arguments override config
		if strings.Contains(os.Args[1], ".txt") || strings.Contains(os.Args[1], ".frames") || strings.Contains(os.Args[1], ".cast") {
			// Load frames from file
			filename := os.Args[1]
			fmt.Printf("Loading frames from file: %s\n", filename)

			// Detect file type and use appropriate parser
			if strings.HasSuffix(filename, ".cast") {
				frames, err = LoadFramesFromCastFile(filename)
			} else {
				frames, err = LoadFramesFromFile(filename)
			}

			if err != nil {
				fmt.Printf("Error loading file: %v\n", err)
				fmt.Printf("Falling back to rain animation...\n")
				frames = []Frame{} // Use rain animation as fallback
			} else {
				fmt.Printf("Successfully loaded %d frames\n", len(frames))
			}

			// Check for frame rate as second argument
			if len(os.Args) > 2 {
				if duration, err := time.ParseDuration(os.Args[2]); err == nil {
					frameRate = duration
				}
			}
		} else {
			// First argument is frame rate
			if duration, err := time.ParseDuration(os.Args[1]); err == nil {
				frameRate = duration
			}
		}
	} else {
		// Use config file setting
		if config.FrameFile != "default" && config.FrameFile != "" {
			fmt.Printf("Loading frames from config file: %s\n", config.FrameFile)

			// Detect file type and use appropriate parser
			if strings.HasSuffix(config.FrameFile, ".cast") {
				frames, err = LoadFramesFromCastFile(config.FrameFile)
			} else {
				frames, err = LoadFramesFromFile(config.FrameFile)
			}

			if err != nil {
				fmt.Printf("Error loading config frame file: %v\n", err)
				fmt.Printf("Falling back to rain animation...\n")
				frames = []Frame{} // Use rain animation as fallback
			} else {
				fmt.Printf("Successfully loaded %d frames from config\n", len(frames))
			}
		}
	}

	// If no frames loaded, use rain animation
	if len(frames) == 0 {
		frames = []Frame{} // Empty frames will trigger rain animation
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create model
	model := Model{
		frames:       frames,
		currentFrame: 0,
		frameRate:    frameRate,
		startTime:    time.Now(),
		sysInfo:      GetSystemInfo(),
		config:       config,
		ctx:          ctx,
		cancel:       cancel,
		mutex:        &sync.RWMutex{},
	}

	// Start the program
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
