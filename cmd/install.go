package cmd

import (
	"errors"

	"github.com/a-was/github-package-manager/cmd/install"
	"github.com/a-was/github-package-manager/config"

	"github.com/spf13/cobra"
)

var force bool

func init() {
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip database check")
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install {repo}",
	Short: "Install repositories",
	Long: `Install requested GitHub repository.
It lets You select asset from latest release, unpack it (if it is an archive)
and select which binary to copy to configurable bin folder.`,
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
		return install.Install(install.Config{
			Repo:  args[0],
			Force: force,
		})
	},
}
