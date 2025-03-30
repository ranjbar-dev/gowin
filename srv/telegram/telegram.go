package telegram

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (t *Telegram) Log(message string) {
	// Get the directory of the running executable
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		return
	}
	execDir := filepath.Dir(execPath)
	logFile := filepath.Join(execDir, "logs.log")

	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	// Create timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Write to file
	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
		return
	}

	// Also print to console
	fmt.Print(logMessage)
}
