package cmd

import (
	"fmt"

	"github.com/a-was/github-package-manager/cmd/install"
	"github.com/a-was/github-package-manager/config"
	"github.com/a-was/github-package-manager/db"
	"github.com/a-was/github-package-manager/github"
	"github.com/a-was/github-package-manager/prompt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

type update struct {
	latest           *github.Release
	from             string
	to               string
	recommendedAsset int
	recommendedFile  int
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all installed repos",
	Long:  "Update checks all installed repositories to see if they have a new release.",
	RunE: func(cmd *cobra.Command, args []string) error {
		toUpdate := map[string]update{}
		for repo, version := range db.GetInstalled() {
			fmt.Printf("Checking %s...\n", repo)
			latest, err := github.GetLatestRelease(repo)
			if err != nil {
				return err
			}
			if version.Tag != latest.Tag {
				toUpdate[repo] = update{
					latest:           latest,
					from:             version.Tag,
					to:               latest.Tag,
					recommendedAsset: version.SelectedAsset,
					recommendedFile:  version.SelectedFile,
				}
			}
		}

		if len(toUpdate) == 0 {
			fmt.Println("Nothing to update")
			return nil
		}

		table := [][]string{
			{"Repository", "Old version", "New version"},
		}
		for repo, release := range toUpdate {
			table = append(table, []string{repo, release.from, release.to})
		}
		fmt.Println()
		prompt.PrintTable(table)
		fmt.Println()

		var input string
		prompt.Get("Proceed? [y/N] ", &input)
		if input != "y" {
			return config.ErrorAborted
		}
		for repo, release := range toUpdate {
			fmt.Println()
			c := install.Config{
				Repo:             repo,
				Release:          release.latest,
				RecommendedAsset: release.recommendedAsset,
				RecommendedFile:  release.recommendedFile,
			}
			if err := install.Update(c); err != nil {
				fmt.Println(err)
				prompt.Get("Continue? [y/N]", &input)
				if input != "y" {
					return config.ErrorAborted
				}
			}
		}
		return nil
	},
}
