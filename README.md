# gobackup

![Release](https://github.com/maliceatthepalace/gobackup/actions/workflows/release.yml/badge.svg)

A simple command-line backup tool written in Go that creates timestamped backups of files.

## Demo

![gobackup demo](assets/demo.gif)

*Note: To generate the demo GIF, install [VHS](https://github.com/charmbracelet/vhs) and run `vhs assets/demo.tape`*

## Installation

### From Source

```bash
go build -o gb main.go
```

Or install globally:
```bash
go install
# Make sure $GOPATH/bin is in your PATH
# The binary will be named 'gb'
```

### Pre-built Binaries

Download pre-built binaries from the [Releases](https://github.com/YOUR_USERNAME/gobackup/releases) page for:
- Linux (AMD64, ARM64)
- Windows (AMD64)
- macOS (Intel, Apple Silicon)

### Build for All Platforms

Use the build script to create binaries for all platforms:

```bash
./build-all.sh
```

This will create binaries in the `releases/` directory for all supported platforms.

## Usage

The gobackup tool supports three modes:

### 1. Backup in same directory
Creates a backup in the same directory as the source file:
```bash
gb /path/to/file.txt
```
Result: `/path/to/file.txt.20231216-120000.backup`

### 2. Backup in current directory
Creates a backup in your current working directory:
```bash
gb /path/to/file.txt .
```

### 3. Backup to target directory
Creates a backup in a specified target directory:
```bash
gb /path/to/file.txt to /target/directory
```

## Backup Format

Backups are created with the format: `filename.YYYYMMDD-HHMMSS.backup`

Example: `document.txt.20231216-143022.backup`

## Requirements

- Go 1.16 or later

## License

MIT

