package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/fs"
)

func TestVerifyIntegrity(t *testing.T) {
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

	// 2. Modify the file manually
	if err := os.WriteFile(testFile, []byte("tampered content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 3. Run the verification logic (currently undefined)
	reports, err := VerifyIntegrity(cfg, b)
	if err != nil {
		t.Fatalf("VerifyIntegrity failed: %v", err)
	}

	// 4. Assert that the anomaly was detected in the manual audit
	if len(reports) == 0 {
		t.Error("VerifyIntegrity did not detect the modification")
	}

	found := false
	for _, report := range reports {
		if report == "CRITICAL: File modified: "+testFile {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected report for modified file %s, but was not found", testFile)
	}
}
