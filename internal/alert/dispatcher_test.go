package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestAlertDispatcher(t *testing.T) {
	// 1. Setup a test server to count received alerts
	var mu sync.Mutex
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()

		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 2. Initialize the dispatcher (currently undefined)
	// We'll use a channel to send alerts to the dispatcher
	alertChan := make(chan Alert, 10)
	stop := make(chan struct{})

	go func() {
		// Mocked or actual start logic
		StartDispatcher(server.URL, alertChan, stop)
	}()

	// 3. Send multiple alerts asynchronously
	for i := 0; i < 5; i++ {
		alertChan <- NewAlert("test-agent", "INFO", "TEST", "/test/path", "test message")
	}

	// Wait briefly for the dispatcher to process
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	if count != 5 {
		t.Errorf("Expected 5 alerts received, got %d", count)
	}
	mu.Unlock()

	close(stop)
}
