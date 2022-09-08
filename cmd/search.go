package cmd

import (
	"errors"
	"fmt"

	"github.com/a-was/github-package-manager/config"
	"github.com/a-was/github-package-manager/github"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search {pattern}",
	Short: "Search for repos",
	Long:  `Search shows first 10 best {pattern} matching results`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expected search pattern")
		}

		if !config.RegexSearchPattern.MatchString(args[0]) {
			return errors.New("invalid pattern")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		repos, err := github.Search(args[0])
		if err != nil {
			return err
		}

		if len(repos) == 0 {
			fmt.Println("Nothing found")
			return nil
		}

		for i, repo := range repos {
			var lang string
			if repo.Language != "" {
				lang = fmt.Sprintf(", %s", repo.Language)
			}
			fmt.Printf("%2d) %s - %d stars%s\n", i+1, repo.FullName, repo.Stars, lang)
			fmt.Printf("\t%s\n", repo.Description)
		}

		return nil
	},
}
