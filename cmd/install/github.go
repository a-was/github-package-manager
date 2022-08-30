package install

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ghAsset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type ghRelease struct {
	ID     int       `json:"id"`
	Name   string    `json:"tag_name"`
	Assets []ghAsset `json:"assets"`
}

func getLatestRelease(repo string) (*ghRelease, error) {
	// get the latest release info
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse the latest release
	release := new(ghRelease)
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFilepath := filepath + ".tmp"

	file, err := os.Create(tmpFilepath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	file.Close()
	os.Rename(tmpFilepath, filepath)
	return nil
}
