package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBaselineSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	baselinePath := filepath.Join(tmpDir, "baseline.json")
	secret := "test-hmac-secret"

	hashes := map[string]string{
		"/etc/passwd": "hash1",
		"/etc/hosts":  "hash2",
	}

	b := NewBaseline(hashes)
	if err := b.Save(baselinePath, secret); err != nil {
		t.Fatalf("Failed to save baseline: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(baselinePath); os.IsNotExist(err) {
		t.Fatal("Baseline file was not created")
	}

	// Load and verify
	loaded, err := LoadBaseline(baselinePath, secret)
	if err != nil {
		t.Fatalf("Failed to load baseline: %v", err)
	}

	if len(loaded.Hashes) != 2 {
		t.Errorf("Expected 2 hashes, got %d", len(loaded.Hashes))
	}

	if loaded.Signature == "" {
		t.Error("Loaded baseline has empty signature")
	}

	// Test tampering detection
	// 1. Wrong secret
	_, err = LoadBaseline(baselinePath, "wrong-secret")
	if err == nil {
		t.Error("LoadBaseline should have failed with wrong secret")
	}

	// 2. Tampered content
	data, _ := os.ReadFile(baselinePath)
	// Modify a character in the JSON (that isn't the signature)
	tamperedData := []byte(string(data)[:50] + "X" + string(data)[51:])
	os.WriteFile(baselinePath, tamperedData, 0600)

	_, err = LoadBaseline(baselinePath, secret)
	if err == nil {
		t.Error("LoadBaseline should have failed with tampered content")
	}
}
