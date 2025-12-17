package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gbckp <source_file(s)> [. | to <target_dir>]\n")
		fmt.Fprintf(os.Stderr, "       gbckp --help\n")
		fmt.Fprintf(os.Stderr, "       gbckp --version\n")
		os.Exit(1)
	}

	// Handle --help flag
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printHelp()
		os.Exit(0)
	}

	// Handle --version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("gbckp version %s\n", version)
		os.Exit(0)
	}

	// Parse arguments to separate source files from backup mode
	var sourceFiles []string
	var targetDir string
	var backupInSameDir bool

	// Find where the mode indicator starts (. or to)
	modeIdx := -1
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "." || os.Args[i] == "to" {
			modeIdx = i
			break
		}
	}

	// Collect source files
	if modeIdx == -1 {
		// No mode indicator, all remaining args are source files
		sourceFiles = os.Args[1:]
		backupInSameDir = true
	} else {
		// Source files are everything before the mode indicator
		sourceFiles = os.Args[1:modeIdx]
		
		// Determine backup mode
		if os.Args[modeIdx] == "." {
			// Case: backup to current directory
			var err error
			targetDir, err = os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: could not get current directory: %v\n", err)
				os.Exit(1)
			}
		} else if os.Args[modeIdx] == "to" {
			// Case: backup to target directory
			if modeIdx+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "Error: target directory required after 'to'\n")
				os.Exit(1)
			}
			targetDir = strings.Join(os.Args[modeIdx+1:], " ")
			
			// Check if target directory exists
			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: target directory '%s' does not exist\n", targetDir)
				os.Exit(1)
			}
		}
	}

	// Validate that we have at least one source file
	if len(sourceFiles) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no source files specified\n")
		os.Exit(1)
	}

	// Process each source file or directory
	successCount := 0
	failCount := 0
	
	for _, source := range sourceFiles {
		// Check if source exists and get info
		info, err := os.Stat(source)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: '%s' does not exist - skipping\n", source)
			failCount++
			continue
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not stat '%s': %v - skipping\n", source, err)
			failCount++
			continue
		}

		var backupPath string
		
		// Automatically detect if it's a file or directory
		if info.IsDir() {
			// Directory → create tar.gz backup
			backupPath, err = backupDirectory(source, targetDir, backupInSameDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to backup directory '%s': %v\n", source, err)
				failCount++
				continue
			}
		} else {
			// File → create .backup file
			backupPath, err = backupFile(source, targetDir, backupInSameDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to backup file '%s': %v\n", source, err)
				failCount++
				continue
			}
		}

		fmt.Printf("Backup created: %s\n", backupPath)
		successCount++
		
		// Small delay to ensure unique timestamps
		time.Sleep(time.Second)
	}

	// Print summary if multiple files
	if len(sourceFiles) > 1 {
		fmt.Printf("\nSummary: %d successful, %d failed\n", successCount, failCount)
	}

	// Exit with error if any backups failed
	if failCount > 0 {
		os.Exit(1)
	}
}

func printHelp() {
	help := `gbckp - A simple command-line backup tool

USAGE:
    gbckp <source_file(s)> [. | to <target_dir>]
    gbckp --help
    gbckp --version

FEATURES:
    - Creates timestamped backups (YYYYMMDD-HHMMSS)
    - Preserves file permissions automatically
    - Supports multiple files and directories at once
    - Automatically detects files vs directories
    - Files → .backup, Directories → .tar.gz

MODES:
    1. Backup in same directory:
       gbckp /path/to/file.txt
       Result: /path/to/file.txt.YYYYMMDD-HHMMSS.backup
       
       gbckp /path/to/mydir
       Result: /path/to/mydir.YYYYMMDD-HHMMSS.tar.gz

    2. Backup in current directory:
       gbckp /path/to/file.txt .
       Result: ./file.txt.YYYYMMDD-HHMMSS.backup
       
       gbckp /path/to/mydir .
       Result: ./mydir.YYYYMMDD-HHMMSS.tar.gz

    3. Backup to target directory:
       gbckp /path/to/file.txt to /target/directory
       Result: /target/directory/file.txt.YYYYMMDD-HHMMSS.backup
       
       gbckp /path/to/mydir to /target/directory
       Result: /target/directory/mydir.YYYYMMDD-HHMMSS.tar.gz

OPTIONS:
    --help, -h      Show this help message
    --version, -v   Show version information

EXAMPLES:
    # Single file backups
    gbckp document.txt
    gbckp /etc/config.conf .
    gbckp important.txt to /backups/
    
    # Directory backups (creates tar.gz)
    gbckp /etc/nginx
    gbckp /var/www to /backups/

    # Multiple files and directories at once
    gbckp file1.txt file2.txt file3.txt
    gbckp file.txt mydir/ to /backups/
    gbckp *.conf .
    
    # Mixed: files AND directories
    gbckp config.txt logs/ data/ to /backups/
`
	fmt.Print(help)
}

// backupFile creates a backup of a single file
func backupFile(sourceFile, targetDir string, backupInSameDir bool) (string, error) {
	fileName := filepath.Base(sourceFile)
	dateStr := time.Now().Format("20060102-150405")
	
	var backupPath string
	if backupInSameDir {
		sourceDir := filepath.Dir(sourceFile)
		backupPath = filepath.Join(sourceDir, fmt.Sprintf("%s.%s.backup", fileName, dateStr))
	} else {
		backupPath = filepath.Join(targetDir, fmt.Sprintf("%s.%s.backup", fileName, dateStr))
	}
	
	// Copy the file with permissions
	if err := copyFile(sourceFile, backupPath); err != nil {
		return "", err
	}
	
	return backupPath, nil
}

// backupDirectory creates a tar.gz backup of a directory
func backupDirectory(sourceDir, targetDir string, backupInSameDir bool) (string, error) {
	dirName := filepath.Base(sourceDir)
	dateStr := time.Now().Format("20060102-150405")
	archiveName := fmt.Sprintf("%s.%s.tar.gz", dirName, dateStr)
	
	var archivePath string
	if backupInSameDir {
		parentDir := filepath.Dir(sourceDir)
		archivePath = filepath.Join(parentDir, archiveName)
	} else {
		archivePath = filepath.Join(targetDir, archiveName)
	}
	
	// Use system tar command to create archive
	// -c: create, -z: gzip, -f: file, -p: preserve permissions
	cmd := exec.Command("tar", "-czpf", archivePath, "-C", filepath.Dir(sourceDir), filepath.Base(sourceDir))
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("tar failed: %v, output: %s", err, output)
	}
	
	return archivePath, nil
}

func copyFile(src, dst string) error {
	// Get source file info to preserve permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	// Preserve file permissions from source
	err = os.Chmod(dst, sourceInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

