package install

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-was/github-package-manager/config"
	"github.com/a-was/github-package-manager/db"
	"github.com/a-was/github-package-manager/github"
	"github.com/a-was/github-package-manager/prompt"
)

var (
	repoFolder = config.RepoFolder
	binFolder  = config.BinFolder
)

func Install(repo string, force bool) error {
	fmt.Println("Using repo", repo)
	release, err := github.GetLatestRelease(repo)
	if err != nil {
		return err
	}

	if !force && db.CheckIfInstalled(release) {
		fmt.Println("Newest version already installed")
		return nil
	}
	return installRelease(repo, release)
}

func Update(repo string, release *github.Release) error {
	fmt.Println("Updating", repo)
	return installRelease(repo, release)
}

func installRelease(repo string, release *github.Release) error {
	fmt.Println("Select download file:")
	for i, a := range release.Assets {
		fmt.Printf("%2d) %s\n", i+1, a.Name)
	}
	var selectedIdx int
	prompt.Get("Your selection: ", &selectedIdx)
	if 0 > selectedIdx || selectedIdx > len(release.Assets) {
		return errors.New("invalid selection")
	}
	selectedAsset := release.Assets[selectedIdx-1]

	filePath := filepath.Join(repoFolder, selectedAsset.Name)

	var err error
	os.RemoveAll(repoFolder)
	os.MkdirAll(repoFolder, 0755)
	os.MkdirAll(binFolder, 0755)

	fmt.Printf("Selected file: %s, downloading...\n", selectedAsset.Name)
	if err = downloadFile(selectedAsset.URL, filePath); err != nil {
		return err
	}
	fmt.Println("Downloaded successfully.")

	if err = uncompressFile(repoFolder, selectedAsset.Name); err != nil {
		return err
	}

	fmt.Println("Select which file to install:")
	i := 1
	filesMap := map[int]string{}
	filepath.Walk(repoFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		p := strings.TrimPrefix(path, fmt.Sprintf("%s/", repoFolder))
		filesMap[i] = p
		fmt.Printf("%2d) %s\n", i, p)
		i++
		return nil
	})

	if len(filesMap) == 1 {
		selectedIdx = 1
	} else {
		prompt.Get("Your selection (empty to skip installation): ", &selectedIdx)
	}

	if selectedIdx == -1 {
		// do not install file
		fmt.Println()
		fmt.Println("Skipped file installation")
		fmt.Printf("Repo cloned into %s folder, you can install it manually\n", repoFolder)
		return nil
	}

	selectedFile, ok := filesMap[selectedIdx]
	if !ok {
		return errors.New("invalid selection")
	}

	baseSelected := filepath.Base(selectedFile)
	if err = copyFile(
		filepath.Join(repoFolder, selectedFile),
		filepath.Join(binFolder, baseSelected),
	); err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("Repo cloned into %s folder\n", repoFolder)
	fmt.Printf("File %s installed into %s folder\n", baseSelected, binFolder)

	db.SaveRelease(release)

	return nil
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

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	_ = file.Close()
	_ = os.Rename(tmpFilepath, filepath)
	return nil
}

func copyFile(source, destination string) error {
	input, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	if err = os.WriteFile(destination, input, 0744); err != nil {
		return err
	}
	return nil
}
