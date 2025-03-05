package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// LoadEnv loads the environment variables from a .env file
func LoadEnv(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	vars := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			varName := parts[0]
			varValue := parts[1]
			vars[varName] = varValue
		}
	}
	return vars, scanner.Err()
}

// RunCommand runs the specified command with the temporary environment
func RunCommand(command string, args []string, tempEnv map[string]string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin // This line enables interactive input
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Save current environment and defer restoration
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range originalEnv {
			parts := strings.SplitN(e, "=", 2)
			os.Setenv(parts[0], parts[1])
		}
	}()

	// Set temporary environment variables
	for key, value := range tempEnv {
		os.Setenv(key, value)
	}

	return cmd.Run()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run script.go /path/to/.env command [args...]")
		os.Exit(1)
	}

	envPath := os.Args[1]
	command := os.Args[2]
	args := os.Args[3:]

	// Load environment variables from .env file
	envVars, err := LoadEnv(envPath)
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Run the command
	err = RunCommand(command, args, envVars)
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1)
	}
}
