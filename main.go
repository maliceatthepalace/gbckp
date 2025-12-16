package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gb <source_file> [. | to <target_dir>]\n")
		os.Exit(1)
	}

	sourceFile := os.Args[1]
	
	// Check if source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: source file '%s' does not exist\n", sourceFile)
		os.Exit(1)
	}

	var targetDir string
	var backupInSameDir bool

	if len(os.Args) == 2 {
		// Case 1: backup /dir/file.txt -> backup in same directory
		backupInSameDir = true
	} else if len(os.Args) == 3 && os.Args[2] == "." {
		// Case 2: backup /dir/file.txt . -> backup in current directory
		var err error
		targetDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not get current directory: %v\n", err)
			os.Exit(1)
		}
	} else if len(os.Args) >= 3 {
		// Case 3: backup /dir/file.txt to /target/dir
		if os.Args[2] == "to" {
			if len(os.Args) < 4 {
				fmt.Fprintf(os.Stderr, "Error: target directory required after 'to'\n")
				os.Exit(1)
			}
			targetDir = strings.Join(os.Args[3:], " ")
		} else {
			fmt.Fprintf(os.Stderr, "Error: invalid syntax. Use: gb <source_file> [. | to <target_dir>]\n")
			os.Exit(1)
		}
	}

	// Determine backup location
	var backupPath string
	if backupInSameDir {
		// Backup in same directory as source file
		sourceDir := filepath.Dir(sourceFile)
		fileName := filepath.Base(sourceFile)
		dateStr := time.Now().Format("20060102-150405")
		backupPath = filepath.Join(sourceDir, fmt.Sprintf("%s.%s.backup", fileName, dateStr))
	} else {
		// Backup in target directory
		// Check if target directory exists
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: target directory '%s' does not exist\n", targetDir)
			os.Exit(1)
		}

		fileName := filepath.Base(sourceFile)
		dateStr := time.Now().Format("20060102-150405")
		backupPath = filepath.Join(targetDir, fmt.Sprintf("%s.%s.backup", fileName, dateStr))
	}

	// Perform the backup
	if err := copyFile(sourceFile, backupPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create backup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Backup created: %s\n", backupPath)
}

func copyFile(src, dst string) error {
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
	return err
}

