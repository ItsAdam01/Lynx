package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsAdam01/Lynx/internal/alert"
	"github.com/ItsAdam01/Lynx/internal/app"
	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/fs"
	"github.com/ItsAdam01/Lynx/internal/logger"
	"github.com/spf13/cobra"
)

var startBaselineInput string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start real-time file integrity monitoring",
	Long: `Loads the configured baseline and begins watching the file system 
for any unauthorized changes using inotify. Events are logged in JSON format
and dispatched asynchronously to the configured webhook.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Initialize structured JSON logging
		if err := logger.InitLogger(cfg.LogFile); err != nil {
			fmt.Printf("Error initializing logger: %v\n", err)
			os.Exit(1)
		}

		secret, err := cfg.GetHmacSecret()
		if err != nil {
			logger.Error("Failed to get HMAC secret", "error", err.Error())
			os.Exit(1)
		}

		// Load and verify the baseline
		baseline, err := fs.LoadBaseline(startBaselineInput, secret)
		if err != nil {
			logger.Error("Failed to load baseline. Integrity compromised or file missing.", "error", err.Error())
			os.Exit(1)
		}

		logger.Info("Starting Lynx FIM agent", "agent_name", cfg.AgentName, "total_baseline_files", baseline.Metadata.TotalFiles)

		// Set up synchronization and channels
		stop := make(chan struct{})
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		anomalies := make(chan string, 100)
		alertChan := make(chan alert.Alert, 100)

		// 1. Start the asynchronous alert dispatcher
		go alert.StartDispatcher(cfg.WebhookURL, alertChan, stop)

		// 2. Start the real-time monitoring system
		go func() {
			if err := app.StartMonitoring(cfg, baseline, anomalies, stop); err != nil {
				logger.Error("Monitor failed", "error", err.Error())
				os.Exit(1)
			}
		}()

		// Main event loop
		for {
			select {
			case anomaly := <-anomalies:
				logger.Warn("Anomaly detected", "details", anomaly)
				fmt.Println(anomaly)

				// 3. Dispatch the alert asynchronously
				alertChan <- alert.NewAlert(cfg.AgentName, "CRITICAL", "FILE_CHANGE", "multiple", anomaly)

			case sig := <-sigChan:
				logger.Info("Shutting down Lynx FIM", "signal", sig.String())
				close(stop)
				return
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&startBaselineInput, "baseline", "b", "baseline.json", "path to the verified baseline file")
}
