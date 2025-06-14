package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanTargets(t *testing.T) {
	// Setup a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a nested file structure:
	// tmpDir/
	//   file1.txt
	//   subdir/
	//     file2.txt
	//     subsubdir/
	//       file3.txt

	file1 := filepath.Join(tmpDir, "file1.txt")
	subdir := filepath.Join(tmpDir, "subdir")
	file2 := filepath.Join(subdir, "file2.txt")
	subsubdir := filepath.Join(subdir, "subsubdir")
	file3 := filepath.Join(subsubdir, "file3.txt")

	if err := os.MkdirAll(subsubdir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file1, []byte("file1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("file2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("file3"), 0644); err != nil {
		t.Fatal(err)
	}

	targets := []string{tmpDir}
	files, err := ScanTargets(targets)
	if err != nil {
		t.Fatalf("ScanTargets failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}

	// Verify all expected files are present
	fileMap := make(map[string]bool)
	for _, f := range files {
		fileMap[f] = true
	}

	expectedFiles := []string{file1, file2, file3}
	for _, ef := range expectedFiles {
		absEf, _ := filepath.Abs(ef)
		if !fileMap[absEf] {
			t.Errorf("Expected file %s not found in ScanTargets results", absEf)
		}
	}
}
