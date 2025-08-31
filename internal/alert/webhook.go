package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Alert represents the structured payload for a security notification.
// We include 'content' for Discord and 'text' for Slack compatibility.
type Alert struct {
	Content   string `json:"content"`   // For Discord
	Text      string `json:"text"`      // For Slack
	Agent     string `json:"agent"`
	Timestamp string `json:"timestamp"`
	Severity  string `json:"severity"`
	Event     string `json:"event"`
	File      string `json:"file"`
	Message   string `json:"message"`
}

// NewAlert creates a new Alert with the current timestamp and formats the content for webhooks.
func NewAlert(agent, severity, event, file, message string) Alert {
	timestamp := time.Now().Format(time.RFC3339)
	
	// Create a readable summary for Discord/Slack
	summary := fmt.Sprintf("🚨 **Lynx FIM Alert** 🚨\n**Agent:** %s\n**Severity:** %s\n**Event:** %s\n**File:** %s\n**Details:** %s\n*Timestamp: %s*", 
		agent, severity, event, file, message, timestamp)

	return Alert{
		Content:   summary,
		Text:      summary,
		Agent:     agent,
		Timestamp: timestamp,
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
