package cmd

import (
	"fmt"

	"github.com/a-was/github-package-manager/cmd/install"
	"github.com/a-was/github-package-manager/db"
	"github.com/a-was/github-package-manager/github"
	"github.com/a-was/github-package-manager/prompt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

type update struct {
	Latest *github.Release
	From   string
	To     string
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all installed repos",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		toUpdate := map[string]update{}
		for repo, version := range db.GetInstalled() {
			fmt.Printf("Checking %s...\n", repo)
			latest, _ := github.GetLatestRelease(repo)
			if version.Tag != latest.Tag {
				toUpdate[repo] = update{
					Latest: latest,
					From:   version.Tag,
					To:     latest.Tag,
				}
			}
		}

		if len(toUpdate) == 0 {
			fmt.Println("Nothing to update")
			return
		}

		table := [][]string{
			{"Repository", "Old version", "New version"},
		}
		fmt.Println()
		for repo, release := range toUpdate {
			table = append(table, []string{repo, release.From, release.To})
		}
		prompt.PrintTable(table)
		fmt.Println()

		var input string
		prompt.Get("Proceed? [y/N]\n", &input)
		if input != "y" {
			fmt.Println("Aborted")
			return
		}
		for repo, release := range toUpdate {
			fmt.Println()
			install.Update(repo, release.Latest)
		}
	},
}
