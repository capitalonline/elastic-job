package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "mongodb-job",
	Short:   "mongodb-job is Elastic Crontab System.",
	Long:    `mongodb-job is Elastic Crontab System.`,
	Version: "0.0.1",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
