package cmd

import (
	"fmt"

	"github.com/a-was/github-package-manager/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Github Package Manager",
	Long:  `All software has versions. This is Github Package Manager's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("github-package-manager %s\n", config.Version)
	},
}
