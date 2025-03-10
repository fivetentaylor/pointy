package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	startMarker = "### START REVISO MANAGED BLOCK"
	endMarker   = "### END REVISO MANAGED BLOCK"
	hostsPath   = "/etc/hosts"
)

// This is your updated content for the managed block
var updatedBlockContent = `
127.0.0.1 reviso.dev
127.0.0.1 app.reviso.dev
127.0.0.1 www.reviso.dev

127.0.0.1 dev.pointy.ai
127.0.0.1 app.dev.pointy.ai
127.0.0.1 www.dev.pointy.ai
`

func main() {
	// Create a temporary file in the same directory as the hosts file to ensure they're on the same filesystem
	tempFilePath, tempFile, err := createTempFile()
	if err != nil {
		fmt.Printf("Failed to create a temporary file: %v\n", err)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFilePath) // Ensure the temp file is removed after execution

	// Open the /etc/hosts file
	hostsFile, err := os.Open(hostsPath)
	if err != nil {
		fmt.Printf("Failed to open the /etc/hosts file: %v\n", err)
		return
	}
	defer hostsFile.Close()

	scanner := bufio.NewScanner(hostsFile)
	inBlock := false
	for scanner.Scan() {
		line := scanner.Text()

		// Check for the start of the managed block
		if line == startMarker {
			inBlock = true
			continue
		}

		// Check for the end of the managed block and skip writing this line
		if line == endMarker {
			inBlock = false
			continue
		}

		// Write lines outside the managed block to the temp file
		if !inBlock {
			if _, err := tempFile.WriteString(line + "\n"); err != nil {
				fmt.Printf("Failed to write to the temporary file: %v\n", err)
				return
			}
		}
	}

	// Write the updated block content to the temp file
	newBlock := fmt.Sprintf("%s\n%s\n%s\n", startMarker, strings.TrimSpace(updatedBlockContent), endMarker)
	if _, err := tempFile.WriteString(newBlock); err != nil {
		fmt.Printf("Failed to write the updated block to the temporary file: %v\n", err)
		return
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading the /etc/hosts file: %v\n", err)
		return
	}

	// Ensure the temp file is properly closed before moving it to replace the /etc/hosts file
	tempFile.Close()

	// Replace the original /etc/hosts file with the updated temp file
	// Note: This operation requires elevated privileges
	if err := os.Rename(tempFilePath, hostsPath); err != nil {
		fmt.Printf("Failed to replace the /etc/hosts file: %v\n", err)
		// Optionally, you could attempt to copy the file manually if Rename fails due to cross-filesystem issues
		return
	}

	if err := os.Chmod(hostsPath, 0o644); err != nil {
		fmt.Printf("Failed to set permissions on the /etc/hosts file: %v\n", err)
		return
	}

	fmt.Println("Successfully updated the /etc/hosts file.")
}

// createTempFile creates a temporary file in the same directory as the target file
func createTempFile() (string, *os.File, error) {
	tempFile, err := os.CreateTemp("/tmp", "hosts_*.tmp")
	if err != nil {
		return "", nil, err
	}
	return tempFile.Name(), tempFile, nil
}
