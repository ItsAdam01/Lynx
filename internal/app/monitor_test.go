package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/fs"
)

func TestMonitorAndDetect(t *testing.T) {
	tmpDir := t.TempDir()
	
	// 1. Setup a baseline
	testFile := filepath.Join(tmpDir, "monitored.txt")
	content := []byte("original content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		PathsToWatch: []string{tmpDir},
		HmacSecretEnv: "LYNX_HMAC_SECRET",
	}
	secret := "test-secret"
	os.Setenv("LYNX_HMAC_SECRET", secret)
	defer os.Unsetenv("LYNX_HMAC_SECRET")

	baselinePath := filepath.Join(tmpDir, "baseline.json")
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte("dummy config"), 0644)

	if err := CreateBaseline(cfg, configPath, baselinePath); err != nil {
		t.Fatal(err)
	}

	b, _ := fs.LoadBaseline(baselinePath, secret)

	// 2. Start the monitor and detection loop
	incidents := make(chan Incident, 10)
	stop := make(chan struct{})
	
	go func() {
		StartMonitoring(cfg, b, incidents, stop)
	}()

	time.Sleep(200 * time.Millisecond)

	// 3. Trigger a modification (should be CRITICAL)
	if err := os.WriteFile(testFile, []byte("tampered content"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case inc := <-incidents:
		if inc.Severity != "CRITICAL" {
			t.Errorf("Expected CRITICAL severity for modification, got %s", inc.Severity)
		}
		if inc.EventType != "FILE_MODIFIED" {
			t.Errorf("Expected FILE_MODIFIED event, got %s", inc.EventType)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for modification detection")
	}

	// 4. Drain any duplicate events (common with fsnotify during rapid writes)
	time.Sleep(100 * time.Millisecond)
loop:
	for {
		select {
		case <-incidents:
		default:
			break loop
		}
	}

	// 5. Trigger a new file (should be WARNING)
	newFile := filepath.Join(tmpDir, "new.txt")
	if err := os.WriteFile(newFile, []byte("new content"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case inc := <-incidents:
		if inc.Severity != "WARNING" {
			t.Errorf("Expected WARNING severity for new file, got %s. Event: %s, Message: %s", inc.Severity, inc.EventType, inc.Message)
		}
		if inc.EventType != "FILE_CREATED" {
			t.Errorf("Expected FILE_CREATED event, got %s", inc.EventType)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for new file detection")
	}
	
	close(stop)
}
