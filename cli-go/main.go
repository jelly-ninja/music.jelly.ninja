package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	// Read all stdin
	code, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	codeStr := strings.TrimSpace(string(code))

	// Filter welcome message if present
	for _, line := range strings.Split(codeStr, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && 
		   !strings.Contains(line, "Strudel") && 
		   !strings.Contains(line, "===") &&
		   !strings.Contains(line, "Run locally") &&
		   !strings.Contains(line, "Now playing") &&
		   !strings.Contains(line, "Ctrl+C") &&
		   !strings.Contains(line, "<host>") {
			codeStr = line
			break
		}
	}

	codeStr = strings.TrimSpace(codeStr)
	if codeStr == "" {
		fmt.Fprintf(os.Stderr, "Error: No valid Strudel code found\n")
		os.Exit(1)
	}

	// Base64 encode the code for URL
	encoded := base64.URLEncoding.EncodeToString([]byte(codeStr))
	
	// Build strudel.cc URL with code in hash
	strudelURL := "https://strudel.cc/#" + encoded

	fmt.Println("Opening:", strudelURL)

	// Open browser
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", strudelURL)
	case "linux":
		cmd = exec.Command("xdg-open", strudelURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", strudelURL)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported platform: %s\n", runtime.GOOS)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening browser: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Strudel opened in your browser!")
}
