# PowerShell release build script for Gophetch

param(
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

$ReleaseDir = "releases"

Write-Host "Building Gophetch v$Version for all platforms..." -ForegroundColor Green

# Create release directory
if (!(Test-Path $ReleaseDir)) {
    New-Item -ItemType Directory -Path $ReleaseDir | Out-Null
}

# Platforms to build for
$Platforms = @(
    "windows/amd64",
    "windows/386",
    "linux/amd64",
    "linux/386",
    "linux/arm64",
    "darwin/amd64",
    "darwin/arm64",
    "android/arm64"
)

# Build for each platform
foreach ($platform in $Platforms) {
    $GOOS, $GOARCH = $platform.Split('/')
    
    Write-Host "Building for $GOOS/$GOARCH..." -ForegroundColor Cyan
    
    # Set output name
    $Output = "gophetch"
    if ($GOOS -eq "windows") {
        $Output = "gophetch.exe"
    }
    
    # Set environment variables
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    
    # Build with release flags
    $OutputPath = "$ReleaseDir/gophetch-${Version}-${GOOS}-${GOARCH}$($Output.Substring(7))"
    $BuildCmd = "go build -ldflags=`"-s -w`" -o $OutputPath"
    
    try {
        Invoke-Expression $BuildCmd
        Write-Host "✓ Built gophetch-${Version}-${GOOS}-${GOARCH}$($Output.Substring(7))" -ForegroundColor Green
    }
    catch {
        Write-Host "✗ Failed to build for $GOOS/$GOARCH" -ForegroundColor Red
        Write-Host "Error: $_" -ForegroundColor Red
        exit 1
    }
}

# Create checksums
Write-Host "Creating checksums..." -ForegroundColor Yellow
Push-Location $ReleaseDir
try {
    Get-ChildItem "gophetch-${Version}-*" | ForEach-Object {
        $hash = (Get-FileHash $_.Name -Algorithm SHA256).Hash
        "$hash  $($_.Name)" | Add-Content "gophetch-${Version}-checksums.txt"
    }
}
finally {
    Pop-Location
}

Write-Host ""
Write-Host "Release v$Version built successfully!" -ForegroundColor Green
Write-Host "Files created in $ReleaseDir/:" -ForegroundColor Cyan
Get-ChildItem $ReleaseDir | Format-Table Name, Length -AutoSize

Write-Host ""
Write-Host "To create a GitHub release:" -ForegroundColor Yellow
Write-Host "1. Create a new release on GitHub with tag v$Version" -ForegroundColor White
Write-Host "2. Upload all files from $ReleaseDir/" -ForegroundColor White
Write-Host "3. Use the checksums file for verification" -ForegroundColor White
