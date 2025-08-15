package cmd

import (
	"fmt"
	"os"

	"github.com/ItsAdam01/Lynx/internal/app"
	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/ItsAdam01/Lynx/internal/fs"
	"github.com/spf13/cobra"
)

var verifyBaselineInput string

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Perform a manual integrity audit",
	Long: `Loads the configured baseline and performs a one-off, comprehensive 
comparison of the entire file system against the stored hashes.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		secret, err := cfg.GetHmacSecret()
		if err != nil {
			fmt.Printf("Error: HMAC secret not found in environment\n")
			os.Exit(1)
		}

		// Load and verify the baseline
		baseline, err := fs.LoadBaseline(verifyBaselineInput, secret)
		if err != nil {
			fmt.Printf("Error: Failed to load/verify baseline: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Starting manual integrity audit for agent: %s\n", cfg.AgentName)
		fmt.Printf("Comparing against baseline from: %s\n", baseline.Metadata.GeneratedAt.Format("2006-01-02 15:04:05"))

		reports, err := app.VerifyIntegrity(cfg, baseline)
		if err != nil {
			fmt.Printf("Error during verification: %v\n", err)
			os.Exit(1)
		}

		if len(reports) == 0 {
			fmt.Println("✅ Integrity Verified: No discrepancies found.")
		} else {
			fmt.Printf("❌ %d Anomaly(s) Detected:\n", len(reports))
			for _, report := range reports {
				fmt.Printf("  - %s\n", report)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&verifyBaselineInput, "baseline", "b", "baseline.json", "path to the verified baseline file")
}
