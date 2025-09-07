#!/bin/bash
# Bash build script for Gophetch

set -e

# Default values
OUTPUT="gophetch"
PLATFORM="linux/amd64"
RELEASE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -o|--output)
            OUTPUT="$2"
            shift 2
            ;;
        -p|--platform)
            PLATFORM="$2"
            shift 2
            ;;
        -r|--release)
            RELEASE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  -o, --output OUTPUT    Output binary name (default: gophetch)"
            echo "  -p, --platform PLATFORM Platform to build for (default: linux/amd64)"
            echo "  -r, --release          Build in release mode (stripped binaries)"
            echo "  -h, --help             Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

echo "Building Gophetch..."

# Set build flags
LDFLAGS=""
if [ "$RELEASE" = true ]; then
    LDFLAGS="-ldflags=\"-s -w\""
    echo "Building in release mode (stripped binaries)"
fi

# Set environment variables
GOOS=$(echo $PLATFORM | cut -d'/' -f1)
GOARCH=$(echo $PLATFORM | cut -d'/' -f2)

export GOOS
export GOARCH

# Build the application
echo "Running: go build $LDFLAGS -o $OUTPUT"

if go build $LDFLAGS -o "$OUTPUT"; then
    echo "Build successful! Output: $OUTPUT"
    
    # Show file info
    if [ -f "$OUTPUT" ]; then
        FILE_SIZE=$(du -h "$OUTPUT" | cut -f1)
        echo "File size: $FILE_SIZE"
    fi
else
    echo "Build failed"
    exit 1
fi
