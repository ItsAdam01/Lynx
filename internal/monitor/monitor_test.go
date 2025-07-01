package monitor

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestMonitor_BasicEvent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "monitor-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	m, err := NewMonitor([]string{tmpDir})
	if err != nil {
		t.Fatalf("Failed to create monitor: %v", err)
	}
	defer m.Close()

	events := make(chan fsnotify.Event, 10)
	
	go func() {
		m.Start(events)
	}()

	time.Sleep(200 * time.Millisecond)

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-events:
		if event.Op&fsnotify.Create != fsnotify.Create {
			t.Errorf("Expected Create event, got %v", event.Op)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for file creation event")
	}
}

func TestMonitor_RecursiveWatching(t *testing.T) {
	// 1. Setup
	tmpDir, err := os.MkdirTemp("", "recursive-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	m, err := NewMonitor([]string{tmpDir})
	if err != nil {
		t.Fatalf("Failed to create monitor: %v", err)
	}
	defer m.Close()

	events := make(chan fsnotify.Event, 10)
	go m.Start(events)

	time.Sleep(200 * time.Millisecond)

	// 2. Create a new subdirectory - this should be automatically watched
	subDir := filepath.Join(tmpDir, "new_folder")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Wait for the directory creation event to be processed
	time.Sleep(200 * time.Millisecond)

	// 3. Create a file inside that new subdirectory
	nestedFile := filepath.Join(subDir, "nested.txt")
	if err := os.WriteFile(nestedFile, []byte("nested"), 0644); err != nil {
		t.Fatal(err)
	}

	// 4. Assert that the event from the nested file was captured
	select {
	case event := <-events:
		if event.Name != nestedFile {
			// We might get the directory create event first, let's drain it
			if event.Name == subDir {
				select {
				case nestedEvent := <-events:
					if nestedEvent.Name != nestedFile {
						t.Errorf("Expected event for %s, got %s", nestedFile, nestedEvent.Name)
					}
				case <-time.After(2 * time.Second):
					t.Fatal("Timed out waiting for nested file event")
				}
			} else {
				t.Errorf("Expected event for %s or %s, got %s", subDir, nestedFile, event.Name)
			}
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for nested file event")
	}
}
