#!/bin/bash
# Release build script for Gophetch

set -e

VERSION=${1:-"1.0.0"}
RELEASE_DIR="releases"

echo "Building Gophetch v$VERSION for all platforms..."

# Create release directory
mkdir -p "$RELEASE_DIR"

# Platforms to build for
PLATFORMS=(
    "windows/amd64"
    "windows/386"
    "linux/amd64"
    "linux/386"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "android/arm64"
)

# Build for each platform
for platform in "${PLATFORMS[@]}"; do
    GOOS=$(echo $platform | cut -d'/' -f1)
    GOARCH=$(echo $platform | cut -d'/' -f2)
    
    echo "Building for $GOOS/$GOARCH..."
    
    # Set output name
    OUTPUT="gophetch"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT="gophetch.exe"
    fi
    
    # Set environment variables
    export GOOS
    export GOARCH
    
    # Build with release flags
    if go build -ldflags="-s -w" -o "$RELEASE_DIR/gophetch-${VERSION}-${GOOS}-${GOARCH}${OUTPUT##gophetch}"; then
        echo "✓ Built gophetch-${VERSION}-${GOOS}-${GOARCH}${OUTPUT##gophetch}"
    else
        echo "✗ Failed to build for $GOOS/$GOARCH"
        exit 1
    fi
done

# Create checksums
echo "Creating checksums..."
cd "$RELEASE_DIR"
sha256sum gophetch-${VERSION}-* > gophetch-${VERSION}-checksums.txt
cd ..

echo ""
echo "Release v$VERSION built successfully!"
echo "Files created in $RELEASE_DIR/:"
ls -la "$RELEASE_DIR/"

echo ""
echo "To create a GitHub release:"
echo "1. Create a new release on GitHub with tag v$VERSION"
echo "2. Upload all files from $RELEASE_DIR/"
echo "3. Use the checksums file for verification"
