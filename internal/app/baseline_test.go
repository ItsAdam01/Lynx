package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/fs"
)

func TestCreateBaseline(t *testing.T) {
	tmpDir := t.TempDir()
	
	// 1. Create a dummy file to monitor
	testFile := filepath.Join(tmpDir, "monitored.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. Setup a dummy config
	cfg := &config.Config{
		PathsToWatch: []string{tmpDir},
		HmacSecretEnv: "LYNX_HMAC_SECRET",
	}
	
	secret := "test-secret"
	os.Setenv("LYNX_HMAC_SECRET", secret)
	defer os.Unsetenv("LYNX_HMAC_SECRET")

	baselinePath := filepath.Join(tmpDir, "baseline.json")

	// 3. Run the coordination logic (this should fail initially as undefined)
	err := CreateBaseline(cfg, baselinePath)
	if err != nil {
		t.Fatalf("CreateBaseline failed: %v", err)
	}

	// 4. Verify the baseline exists and is valid
	b, err := fs.LoadBaseline(baselinePath, secret)
	if err != nil {
		t.Fatalf("Failed to load/verify created baseline: %v", err)
	}

	if b.Metadata.TotalFiles != 1 {
		t.Errorf("Expected 1 file in baseline, got %d", b.Metadata.TotalFiles)
	}
}
