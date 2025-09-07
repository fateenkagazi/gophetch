# PowerShell build script for Gophetch
param(
    [string]$Output = "gophetch.exe",
    [string]$Platform = "windows/amd64",
    [switch]$Release = $false
)

Write-Host "Building Gophetch..." -ForegroundColor Green

# Set build flags
$ldflags = ""
if ($Release) {
    $ldflags = "-ldflags=`"-s -w`""
    Write-Host "Building in release mode (stripped binaries)" -ForegroundColor Yellow
}

# Set environment variables
$env:GOOS = $Platform.Split('/')[0]
$env:GOARCH = $Platform.Split('/')[1]

# Build the application
$buildCmd = "go build $ldflags -o $Output"
Write-Host "Running: $buildCmd" -ForegroundColor Cyan

try {
    Invoke-Expression $buildCmd
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build successful! Output: $Output" -ForegroundColor Green
        
        # Show file info
        if (Test-Path $Output) {
            $fileInfo = Get-Item $Output
            Write-Host "File size: $([math]::Round($fileInfo.Length / 1MB, 2)) MB" -ForegroundColor Cyan
        }
    } else {
        Write-Host "Build failed with exit code: $LASTEXITCODE" -ForegroundColor Red
        exit $LASTEXITCODE
    }
} catch {
    Write-Host "Build error: $_" -ForegroundColor Red
    exit 1
}
