package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "github-package-manager",
	Short: "A GitHub package manager",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
