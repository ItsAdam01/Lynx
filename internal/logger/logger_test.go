package logger

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitLogger(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "lynx.log")

	// 1. Initialize the logger (should fail initially since it's not defined)
	err := InitLogger(logFile)
	if err != nil {
		t.Fatalf("InitLogger failed: %v", err)
	}

	// 2. Write a test log entry
	Info("test info log", "key", "value")

	// 3. Verify the log file was created and contains valid JSON
	file, err := os.Open(logFile)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		t.Fatal("Expected at least one line in log file")
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal(scanner.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log line as JSON: %v", err)
	}

	if logEntry["msg"] != "test info log" {
		t.Errorf("Expected msg 'test info log', got '%v'", logEntry["msg"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("Expected key 'value', got '%v'", logEntry["key"])
	}
}
