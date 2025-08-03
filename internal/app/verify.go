package app

import (
	"fmt"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/crypto"
	"github.com/ItsAdam01/Lynx/internal/fs"
)

// VerifyIntegrity performs a one-off, manual comparison of the entire configured 
// file system against the memory-loaded baseline.
func VerifyIntegrity(cfg *config.Config, b *fs.Baseline) ([]string, error) {
	var reports []string

	// 1. Scan targets to find all files
	targets := append(cfg.PathsToWatch, cfg.FilesToWatch...)
	files, err := fs.ScanTargets(targets)
	if err != nil {
		return nil, fmt.Errorf("failed to scan targets: %w", err)
	}

	// 2. Track seen files to identify deletions later
	seenFiles := make(map[string]bool)

	// 3. Hash each file and compare with the baseline
	for _, path := range files {
		seenFiles[path] = true
		newHash, err := crypto.HashFile(path)
		if err != nil {
			reports = append(reports, fmt.Sprintf("ERROR: Failed to hash file: %s", path))
			continue
		}

		oldHash, exists := b.Hashes[path]
		if !exists {
			reports = append(reports, fmt.Sprintf("WARNING: New file created: %s", path))
			continue
		}

		if newHash != oldHash {
			reports = append(reports, fmt.Sprintf("CRITICAL: File modified: %s", path))
		}
	}

	// 4. Identify deleted files (in baseline but not on disk)
	for path := range b.Hashes {
		if !seenFiles[path] {
			reports = append(reports, fmt.Sprintf("CRITICAL: File deleted: %s", path))
		}
	}

	return reports, nil
}
