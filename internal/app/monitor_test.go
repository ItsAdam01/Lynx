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
		PathsToWatch:  []string{tmpDir},
		HmacSecretEnv: "LYNX_HMAC_SECRET",
	}
	secret := "test-secret"
	os.Setenv("LYNX_HMAC_SECRET", secret)
	defer os.Unsetenv("LYNX_HMAC_SECRET")

	baselinePath := filepath.Join(tmpDir, "baseline.json")
	if err := CreateBaseline(cfg, baselinePath); err != nil {
		t.Fatal(err)
	}

	b, _ := fs.LoadBaseline(baselinePath, secret)

	// 2. Start the monitor and detection loop (currently undefined)
	// We'll use a channel to capture detected anomalies
	anomalies := make(chan string, 10)
	stop := make(chan struct{})

	go func() {
		// Mocked or actual start logic
		StartMonitoring(cfg, b, anomalies, stop)
	}()

	// Wait for monitor to settle
	time.Sleep(200 * time.Millisecond)

	// 3. Trigger a modification
	if err := os.WriteFile(testFile, []byte("tampered content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 4. Assert that the anomaly was detected
	select {
	case msg := <-anomalies:
		if msg == "" {
			t.Error("Detected empty anomaly message")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for anomaly detection")
	}

	close(stop)
}
