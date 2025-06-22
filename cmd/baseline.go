package cmd

import (
	"fmt"
	"os"

	"github.com/ItsAdam01/Lynx/internal/app"
	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/spf13/cobra"
)

var baselineOutput string

// baselineCmd represents the baseline command
var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Create a new cryptographic baseline of the configured paths",
	Long: `Scans all configured files and directories, calculates SHA-256 hashes, 
and saves a signed baseline.json for future integrity verification.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("Error loading config: %v
", err)
			os.Exit(1)
		}

		if err := app.CreateBaseline(cfg, baselineOutput); err != nil {
			fmt.Printf("Error creating baseline: %v
", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created baseline: %s
", baselineOutput)
	},
}

func init() {
	RootCmd.AddCommand(baselineCmd)
	baselineCmd.Flags().StringVarP(&baselineOutput, "output", "o", "baseline.json", "output path for the baseline file")
}
