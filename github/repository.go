package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
}

type ghApiRepos struct {
	Items []*Repository `json:"items"`
}

func Search(pattern string) ([]*Repository, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/search/repositories?q=%s&per_page=10", pattern))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse best matches
	repos := &ghApiRepos{}
	err = json.NewDecoder(resp.Body).Decode(repos)
	if err != nil {
		return nil, err
	}
	return repos.Items, nil
}
