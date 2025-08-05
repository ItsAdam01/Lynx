package alert

import (
	"log"

	"github.com/ItsAdam01/Lynx/internal/logger"
)

// StartDispatcher runs a non-blocking loop to handle outgoing alert delivery.
// This function consumes alerts from the channel and sends them to the webhook URL.
func StartDispatcher(webhookURL string, alertChan <-chan Alert, stop <-chan struct{}) {
	for {
		select {
		case a := <-alertChan:
			// Send the alert asynchronously
			go func(alert Alert) {
				if err := SendWebhook(webhookURL, alert); err != nil {
					// We'll log the error but keep the dispatcher running
					logger.Error("Failed to deliver alert via webhook", "error", err.Error())
					log.Printf("Alert Delivery Error: %v\n", err)
				}
			}(a)

		case <-stop:
			logger.Info("Stopping alert dispatcher")
			return
		}
	}
}
