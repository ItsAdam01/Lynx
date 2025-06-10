package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lynx",
	Short: "Lynx FIM is a host-based intrusion detection agent",
	Long: `A lightweight File Integrity Monitor (FIM) built in Go to monitor 
critical system files for unauthorized changes. 

This project is part of a self-directed study in cybersecurity and 
real-time system monitoring.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file (default is config.yaml)")
}
