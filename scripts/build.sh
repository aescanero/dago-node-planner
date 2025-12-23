#!/bin/bash
# Build script for dago-node-planner

set -e

VERSION="${VERSION:-dev}"
BUILD_DIR="./bin"
APP_NAME="node-planner"

echo "Building dago-node-planner version: $VERSION"

# Create build directory
mkdir -p "$BUILD_DIR"

# Build for current platform
echo "Building for $(go env GOOS)/$(go env GOARCH)..."
go build -ldflags "-X main.version=$VERSION" -o "$BUILD_DIR/$APP_NAME" ./cmd/node-planner

echo "Build complete: $BUILD_DIR/$APP_NAME"

# Make executable
chmod +x "$BUILD_DIR/$APP_NAME"

# Display version
"$BUILD_DIR/$APP_NAME" -version 2>/dev/null || true

echo "Build successful!"
