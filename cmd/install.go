package cmd

import (
	"errors"
	"github-package-manager/cmd/install"
	"regexp"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install {repo}",
	Short: "Install repositories",
	Long:  `Install GitHub repository.`, // TODO:
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expected a github repo")
		}

		if !regexp.MustCompile(`^[a-zA-Z0-9\-]+/[a-zA-Z0-9\-]+$`).MatchString(args[0]) {
			return errors.New("invalid github repo")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return install.Install(args[0])
	},
}
