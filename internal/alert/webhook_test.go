package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendWebhook(t *testing.T) {
	// 1. Setup a test server to capture the webhook POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Failed to decode webhook payload: %v", err)
		}

		if payload["agent"] != "test-agent" {
			t.Errorf("Expected agent 'test-agent', got '%v'", payload["agent"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 2. Define the alert (currently undefined)
	a := Alert{
		Agent:     "test-agent",
		Severity:  "CRITICAL",
		Event:     "FILE_MODIFIED",
		File:      "/etc/passwd",
		Message:   "Unauthorized change detected",
		Timestamp: "2025-08-01T10:00:00Z",
	}

	// 3. Send the webhook (currently undefined)
	err := SendWebhook(server.URL, a)
	if err != nil {
		t.Fatalf("SendWebhook failed: %v", err)
	}
}
