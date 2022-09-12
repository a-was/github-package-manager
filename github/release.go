package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-was/github-package-manager/config"
)

type Asset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type Release struct {
	ID     int     `json:"id"`
	Tag    string  `json:"tag_name"`
	Assets []Asset `json:"assets"`
	Repo   string  `json:"-"`
}

func GetLatestRelease(repo string) (*Release, error) {
	// get the latest release info
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, config.ErrorApiWrongStatus
	}

	// parse the latest release
	release := &Release{
		Repo: repo,
	}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		return nil, err
	}
	return release, nil
}
