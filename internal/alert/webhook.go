package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Alert represents the structured payload for a security notification.
type Alert struct {
	Agent     string `json:"agent"`
	Timestamp string `json:"timestamp"`
	Severity  string `json:"severity"`
	Event     string `json:"event"`
	File      string `json:"file"`
	Message   string `json:"message"`
}

// NewAlert creates a new Alert with the current timestamp.
func NewAlert(agent, severity, event, file, message string) Alert {
	return Alert{
		Agent:     agent,
		Timestamp: time.Now().Format(time.RFC3339),
		Severity:  severity,
		Event:     event,
		File:      file,
		Message:   message,
	}
}

// SendWebhook performs an HTTP POST request to the specified URL with the alert payload.
func SendWebhook(url string, a Alert) error {
	if url == "" {
		return nil // No webhook configured, skip alert delivery
	}

	payload, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send webhook POST: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook alert failed with status: %s", resp.Status)
	}

	return nil
}
