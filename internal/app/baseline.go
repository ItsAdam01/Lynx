package app

import (
	"fmt"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/crypto"
	"github.com/ItsAdam01/Lynx/internal/fs"
)

// CreateBaseline coordinates the scanning, hashing, and saving of a new file integrity baseline.
func CreateBaseline(cfg *config.Config, configPath, baselinePath string) error {
	// 1. Get the HMAC secret from environment
	secret, err := cfg.GetHmacSecret()
	if err != nil {
		return err
	}

	// 2. Hash the configuration file itself to protect the ignore list and watch paths
	cfgHash, err := crypto.HashFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to hash config file: %w", err)
	}

	// 3. Scan targets to find all files, applying ignore patterns
	targets := append(cfg.PathsToWatch, cfg.FilesToWatch...)
	files, err := fs.ScanTargets(targets, cfg.IgnoredPatterns)
	if err != nil {
		return fmt.Errorf("failed to scan targets: %w", err)
	}

	// 4. Hash each file
	hashes := make(map[string]string)
	for _, path := range files {
		hash, err := crypto.HashFile(path)
		if err != nil {
			return fmt.Errorf("failed to hash file %s: %w", path, err)
		}
		hashes[path] = hash
	}

	// 5. Create and save the baseline with the config fingerprint
	b := fs.NewBaseline(hashes, cfgHash)
	if err := b.Save(baselinePath, secret); err != nil {
		return fmt.Errorf("failed to save baseline: %w", err)
	}

	return nil
}
