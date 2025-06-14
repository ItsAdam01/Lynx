package fs

import (
	"os"
	"path/filepath"
)

// ScanTargets takes a list of paths and a list of ignore patterns, returning 
// all unique absolute file paths that do not match the ignored patterns.
func ScanTargets(targets []string, ignoredPatterns []string) ([]string, error) {
	fileMap := make(map[string]struct{})

	for _, target := range targets {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			continue
		}

		err = filepath.Walk(absTarget, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if info.IsDir() {
				return nil
			}

			// Check if the file matches any ignore pattern
			fileName := filepath.Base(path)
			for _, pattern := range ignoredPatterns {
				matched, err := filepath.Match(pattern, fileName)
				if err == nil && matched {
					return nil // Skip this file
				}
			}

			fileMap[path] = struct{}{}
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
