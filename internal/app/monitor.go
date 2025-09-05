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
func StartMonitoring(cfg *config.Config, b *fs.Baseline, incidents chan<- Incident, stop <-chan struct{}) error {
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
			handleEvent(event, b, incidents)

		case <-stop:
			return nil
		}
	}
}

// handleEvent compares a file event with the baseline and reports incidents.
func handleEvent(event fsnotify.Event, b *fs.Baseline, incidents chan<- Incident) {
	// 1. Handle deletion or rename (CRITICAL)
	// Note: A rename often sends a 'Rename' event for the old path.
	if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
		if _, exists := b.Hashes[event.Name]; exists {
			incidents <- Incident{
				Severity:  "CRITICAL",
				EventType: "FILE_DELETED",
				FilePath:  event.Name,
				Message:   fmt.Sprintf("Monitored file was deleted or renamed: %s", event.Name),
			}
		}
		return
	}

	// 2. Handle creation or modification
	if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
		newHash, err := crypto.HashFile(event.Name)
		if err != nil {
			return
		}

		oldHash, exists := b.Hashes[event.Name]
		if !exists {
			// New file (WARNING)
			incidents <- Incident{
				Severity:  "WARNING",
				EventType: "FILE_CREATED",
				FilePath:  event.Name,
				Message:   fmt.Sprintf("New file created in monitored path: %s", event.Name),
			}
			return
		}

		if newHash != oldHash {
			// Modification (CRITICAL)
			incidents <- Incident{
				Severity:  "CRITICAL",
				EventType: "FILE_MODIFIED",
				FilePath:  event.Name,
				Message:   fmt.Sprintf("Monitored file was modified: %s", event.Name),
			}
		}
	}
}
