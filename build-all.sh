#!/bin/bash
# Build script for gbckp - builds for all platforms

set -e

echo "Building gbckp for all platforms..."

# Create releases directory
mkdir -p releases

# Linux 64-bit
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o releases/gbckp-linux-amd64 main.go

# Linux ARM64
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -o releases/gbckp-linux-arm64 main.go

# Windows 64-bit
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -o releases/gbckp-windows-amd64.exe main.go

# macOS Intel
echo "Building for macOS Intel..."
GOOS=darwin GOARCH=amd64 go build -o releases/gbckp-darwin-amd64 main.go

# macOS Apple Silicon
echo "Building for macOS Apple Silicon..."
GOOS=darwin GOARCH=arm64 go build -o releases/gbckp-darwin-arm64 main.go

echo ""
echo "✅ Build complete! Binaries are in the releases/ directory:"
ls -lh releases/

