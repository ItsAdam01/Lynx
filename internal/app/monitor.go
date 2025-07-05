package app

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/crypto"
	"github.com/ItsAdam01/Lynx/internal/fs"
	"github.com/ItsAdam01/Lynx/internal/monitor"
	"github.com/fsnotify/fsnotify"
)

// debounceInterval defines how long to wait for more events on the same file.
const debounceInterval = 500 * time.Millisecond

// StartMonitoring initializes the real-time event monitor and begins the 
// anomaly detection loop with event debouncing.
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

	// Map to track timers for each file path
	pendingEvents := make(map[string]*time.Timer)
	var mu sync.Mutex

	for {
		select {
		case event := <-events:
			mu.Lock()
			// If a timer already exists for this file, stop it
			if timer, exists := pendingEvents[event.Name]; exists {
				timer.Stop()
			}

			// Start a new timer for this file
			pendingEvents[event.Name] = time.AfterFunc(debounceInterval, func() {
				mu.Lock()
				delete(pendingEvents, event.Name)
				mu.Unlock()
				
				// Process the final state of the file after the interval
				handleEvent(event, b, incidents)
			})
			mu.Unlock()

		case <-stop:
			// Cleanup timers on shutdown
			mu.Lock()
			for _, t := range pendingEvents {
				t.Stop()
			}
			mu.Unlock()
			return nil
		}
	}
}

// handleEvent compares a file event with the baseline and reports incidents.
func handleEvent(event fsnotify.Event, b *fs.Baseline, incidents chan<- Incident) {
	// 1. Handle deletion or rename (CRITICAL)
	// After debouncing, if the file is gone, report deletion.
	_, err := os.Stat(event.Name)
	if os.IsNotExist(err) {
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
