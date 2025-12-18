# gbckp

![Release](https://github.com/maliceatthepalace/gbckp/actions/workflows/release.yml/badge.svg)

A simple command-line backup tool written in Go that creates timestamped backups of files.

## Demo

![gbckp demo](assets/demo.gif)

*Note: To generate the demo GIF, install [VHS](https://github.com/charmbracelet/vhs) and run `vhs assets/demo.tape`*

## Installation

### via wget on linux

```bash
curl -L https://github.com/maliceatthepalace/gbckp/releases/download/v0.1.1/gb-linux-amd64 \
  -o /tmp/gbckp && \
sudo mv /tmp/gbckp /usr/local/bin/gbckp && \
sudo chmod +x /usr/local/bin/gbckp
```

### From Source

```bash
go build -o gbckp main.go
```

Or install globally:
```bash
go install
# Make sure $GOPATH/bin is in your PATH
# The binary will be named 'gbckp'
```

### Pre-built Binaries

Download pre-built binaries from the [Releases](https://github.com/YOUR_USERNAME/gbckp/releases) page for:
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

```bash
gbckp <source_file(s)> [. | to <target_dir>]
gbckp --help
gbckp --version
```

The gbckp tool supports three modes and can backup multiple files at once:

### 1. Backup in same directory
Creates a backup in the same directory as the source:
```bash
# Files
gbckp /path/to/file.txt
# → /path/to/file.txt.20231216-120000.backup

# Directories
gbckp /etc/nginx
# → /etc/nginx.20231216-120000.tar.gz

# Multiple mixed
gbckp file1.txt file2.txt mydir/
```

### 2. Backup in current directory
Creates a backup in your current working directory:
```bash
gbckp /path/to/file.txt .
gbckp /etc/nginx .
gbckp file1.txt mydir/ .
```

### 3. Backup to target directory
Creates a backup in a specified target directory:
```bash
gbckp /path/to/file.txt to /backups/
gbckp /etc/nginx to /backups/
gbckp file1.txt file2.txt mydir/ to /backups/
```

### Options

- `--help`, `-h`: Show help message with usage information
- `--version`, `-v`: Show version information

## Features

- **Timestamped backups**: Each backup includes a timestamp in the format `YYYYMMDD-HHMMSS`
- **Preserves permissions**: File permissions (like 755, 644, etc.) are automatically preserved
- **Multiple files and directories**: Backup multiple files and directories at once with a single command
- **Auto-detection**: Automatically detects files vs directories - files get `.backup`, directories get `.tar.gz`
- **Multiple backup modes**: Same directory, current directory, or custom target directory

## Backup Format

Backups are created with the format: `filename.YYYYMMDD-HHMMSS.backup`

Example: `document.txt.20231216-143022.backup`

## Requirements

- Go 1.16 or later
