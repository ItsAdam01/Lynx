package app

import (
	"fmt"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/crypto"
	"github.com/ItsAdam01/Lynx/internal/fs"
	"github.com/ItsAdam01/Lynx/internal/monitor"
	"github.com/fsnotify/fsnotify"
)

// StartMonitoring initializes the real-time event monitor and begins the 
// anomaly detection loop.
func StartMonitoring(cfg *config.Config, b *fs.Baseline, anomalies chan<- string, stop <-chan struct{}) error {
	m, err := monitor.NewMonitor(cfg.PathsToWatch)
	if err != nil {
		return fmt.Errorf("failed to start monitor: %w", err)
	}
	defer m.Close()

	events := make(chan fsnotify.Event, 100)
	go func() {
		m.Start(events)
	}()

	for {
		select {
		case event := <-events:
			// 1. Identify the event type and handle accordingly
			handleEvent(event, b, anomalies)

		case <-stop:
			return nil
		}
	}
}

// handleEvent compares a file event with the baseline and reports anomalies.
func handleEvent(event fsnotify.Event, b *fs.Baseline, anomalies chan<- string) {
	// 1. Handle deletion
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		if _, exists := b.Hashes[event.Name]; exists {
			anomalies <- fmt.Sprintf("CRITICAL: File deleted: %s", event.Name)
		}
		return
	}

	// 2. Handle creation or modification
	if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
		// Calculate the new hash
		newHash, err := crypto.HashFile(event.Name)
		if err != nil {
			// Might be a directory or permission issue, skip for now
			return
		}

		oldHash, exists := b.Hashes[event.Name]
		if !exists {
			anomalies <- fmt.Sprintf("WARNING: New file created: %s", event.Name)
			return
		}

		if newHash != oldHash {
			anomalies <- fmt.Sprintf("CRITICAL: File modified: %s", event.Name)
		}
	}
}
