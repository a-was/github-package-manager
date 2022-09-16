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

type Config struct {
	Repo    string
	Release *github.Release

	Force bool

	RecommendedAsset int
	RecommendedFile  int
}

func Install(c Config) error {
	fmt.Println("Using repo", c.Repo)
	release, err := github.GetLatestRelease(c.Repo)
	if err != nil {
		return err
	}

	if !c.Force && db.CheckIfInstalled(release) {
		fmt.Println("Newest version already installed")
		return nil
	}

	c.Release = release

	return installRelease(c)
}

func Update(c Config) error {
	fmt.Println("Updating", c.Repo)
	return installRelease(c)
}

func installRelease(c Config) error {
	fmt.Println("Select download file:")
	for i, a := range c.Release.Assets {
		if i+1 == c.RecommendedAsset {
			fmt.Printf("%2d) (recommended) %s\n", i+1, a.Name)
		} else {
			fmt.Printf("%2d)  %s\n", i+1, a.Name)
		}
	}
	var selectedAssetIdx int
	prompt.Get("Your selection: ", &selectedAssetIdx)
	if 0 > selectedAssetIdx || selectedAssetIdx > len(c.Release.Assets) {
		return errors.New("invalid selection")
	}
	selectedAsset := c.Release.Assets[selectedAssetIdx-1]

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
	var skipped bool
	filesMap := map[int]string{}
	filepath.Walk(repoFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if i > 50 {
			skipped = true
			return nil
		}
		if info.IsDir() {
			return nil
		}

		p := strings.TrimPrefix(path, fmt.Sprintf("%s/", repoFolder))
		filesMap[i] = p
		if i == c.RecommendedFile {
			fmt.Printf("%2d) (recommended) %s\n", i, p)
		} else {
			fmt.Printf("%2d)  %s\n", i, p)
		}
		i++
		return nil
	})

	if skipped {
		fmt.Println("Warning! Only first 50 elements are shown")
	}
	var selectedFileIdx int
	prompt.Get("Your selection (empty to skip installation): ", &selectedFileIdx)
	if selectedFileIdx == -1 {
		// do not install file
		fmt.Println()
		fmt.Println("Skipped file installation")
		fmt.Printf("Repo cloned into %s folder, you can install it manually\n", repoFolder)

		db.SaveVersion(c.Release, selectedAssetIdx, selectedFileIdx)
		return nil
	}

	selectedFile, ok := filesMap[selectedFileIdx]
	if !ok {
		return errors.New("invalid selection")
	}
	baseSelected := filepath.Base(selectedFile)

	installed, err := installPackage(selectedFile)
	if err != nil {
		return err
	}
	if installed {
		fmt.Println()
		fmt.Printf("File %s installed via package manager\n", baseSelected)

		db.SaveVersion(c.Release, selectedAssetIdx, selectedFileIdx)
		return nil
	}

	if err = copyFile(
		filepath.Join(repoFolder, selectedFile),
		filepath.Join(binFolder, baseSelected),
	); err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("Repo cloned into %s folder\n", repoFolder)
	fmt.Printf("File %s installed into %s folder\n", baseSelected, binFolder)

	db.SaveVersion(c.Release, selectedAssetIdx, selectedFileIdx)

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
