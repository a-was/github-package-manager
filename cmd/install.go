package cmd

import (
	"errors"

	"github.com/a-was/github-package-manager/cmd/install"
	"github.com/a-was/github-package-manager/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install {repo}",
	Short: "Install repositories",
	Long:  `Install GitHub repository.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expected a github repo")
		}

		if !config.RegexRepo.MatchString(args[0]) {
			return errors.New("invalid github repo")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return install.Install(args[0])
	},
}
