package app

import (
	"fmt"
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
		PathsToWatch:    []string{tmpDir},
		HmacSecretEnv:   "LYNX_HMAC_SECRET",
		IgnoredPatterns: []string{"config.yaml"},
	}

	secret := "test-secret"
	os.Setenv("LYNX_HMAC_SECRET", secret)
	defer os.Unsetenv("LYNX_HMAC_SECRET")

	baselinePath := filepath.Join(tmpDir, "baseline.json")
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte("dummy config"), 0644)

	// 3. Run the coordination logic
	err := CreateBaseline(cfg, configPath, baselinePath)
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

func BenchmarkCreateBaseline_100Files(b *testing.B) {
	tmpDir := b.TempDir()
	for i := 0; i < 100; i++ {
		path := filepath.Join(tmpDir, fmt.Sprintf("file-%d.txt", i))
		os.WriteFile(path, []byte("some repetitive content for hashing"), 0644)
	}

	cfg := &config.Config{
		PathsToWatch:  []string{tmpDir},
		HmacSecretEnv: "LYNX_HMAC_SECRET",
	}
	os.Setenv("LYNX_HMAC_SECRET", "bench-secret")
	baselinePath := filepath.Join(tmpDir, "baseline.json")
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte("bench config"), 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateBaseline(cfg, configPath, baselinePath)
	}
}
