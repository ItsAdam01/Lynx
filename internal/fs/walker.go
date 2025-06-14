package fs

import (
	"os"
	"path/filepath"
)

// ScanTargets takes a list of paths (directories or files) and returns a slice 
// of all unique absolute file paths found.
func ScanTargets(targets []string) ([]string, error) {
	fileMap := make(map[string]struct{})

	for _, target := range targets {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			continue // Skip paths that can't be resolved
		}

		err = filepath.Walk(absTarget, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// We only care about files, not directories
			if !info.IsDir() {
				fileMap[path] = struct{}{}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	files := make([]string, 0, len(fileMap))
	for f := range fileMap {
		files = append(files, f)
	}

	return files, nil
}
