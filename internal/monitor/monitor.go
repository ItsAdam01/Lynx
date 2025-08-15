package monitor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Monitor handles real-time file system event watching.
type Monitor struct {
	watcher *fsnotify.Watcher
	paths   []string
}

// NewMonitor initializes a new Monitor for the given set of target paths.
func NewMonitor(targets []string) (*Monitor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	m := &Monitor{
		watcher: watcher,
		paths:   targets,
	}

	for _, target := range targets {
		if err := m.addRecursive(target); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// addRecursive recursively adds a path and all its subdirectories to the watcher.
func (m *Monitor) addRecursive(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}

	return filepath.Walk(absPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := m.watcher.Add(p); err != nil {
				return fmt.Errorf("failed to watch directory %s: %w", p, err)
			}
		}
		return nil
	})
}

// Start runs the event monitoring loop and sends events to the provided channel.
// This function blocks until the monitor is closed.
func (m *Monitor) Start(eventChan chan<- fsnotify.Event) error {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok {
				return nil
			}

			// If a new directory is created, we MUST add it to the watcher
			// to ensure recursive watching works in real-time.
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					m.addRecursive(event.Name)
				}
			}

			// Send the event to our internal channel for processing
			eventChan <- event

		case err, ok := <-m.watcher.Errors:
			if !ok {
				return nil
			}
			return fmt.Errorf("watcher error: %w", err)
		}
	}
}

// Close gracefully shuts down the file system watcher.
func (m *Monitor) Close() error {
	return m.watcher.Close()
}
