package cmd

import (
	"fmt"
	"os"

	"github.com/ItsAdam01/Lynx/internal/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the default configuration for Lynx FIM",
	Long: `Creates a default config.yaml file in the current directory. 
This file defines which paths and files the agent should monitor.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := "config.yaml"
		if err := config.InitConfig(path); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully initialized configuration: %s\n", path)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
