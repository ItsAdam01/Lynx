package app

import (
	"fmt"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/crypto"
	"github.com/ItsAdam01/Lynx/internal/fs"
)

// CreateBaseline coordinates the scanning, hashing, and saving of a new file integrity baseline.
func CreateBaseline(cfg *config.Config, baselinePath string) error {
	// 1. Get the HMAC secret from environment
	secret, err := cfg.GetHmacSecret()
	if err != nil {
		return err
	}

	// 2. Scan targets to find all files
	targets := append(cfg.PathsToWatch, cfg.FilesToWatch...)
	files, err := fs.ScanTargets(targets)
	if err != nil {
		return fmt.Errorf("failed to scan targets: %w", err)
	}

	// 3. Hash each file
	hashes := make(map[string]string)
	for _, path := range files {
		hash, err := crypto.HashFile(path)
		if err != nil {
			// In a real-world scenario, we might want to log this and continue
			// but for now, we'll treat it as a hard failure for safety.
			return fmt.Errorf("failed to hash file %s: %w", path, err)
		}
		hashes[path] = hash
	}

	// 4. Create and save the baseline
	b := fs.NewBaseline(hashes)
	if err := b.Save(baselinePath, secret); err != nil {
		return fmt.Errorf("failed to save baseline: %w", err)
	}

	return nil
}
