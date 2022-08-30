package install

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/a-was/github-package-manager/config"
	"github.com/a-was/github-package-manager/db"
	"github.com/a-was/github-package-manager/github"
)

var (
	repoFolder = config.RepoFolder
	binFolder  = config.BinFolder
)

func Install(repo string) error {
	fmt.Println("Using repo", repo)
	release, err := github.GetLatestRelease(repo)
	if err != nil {
		return err
	}

	if db.CheckIfInstalled(release) {
		fmt.Println("Newest version already installed")
		return nil
	}

	fmt.Println("Select download file:")
	for i, a := range release.Assets {
		fmt.Printf("%2d) %s\n", i+1, a.Name)
	}
	var selectedStr string
	fmt.Print("Your selection: ")
	fmt.Scanln(&selectedStr)
	selectedIdx, err := strconv.Atoi(selectedStr)
	if err != nil {
		return err
	}
	if 0 > selectedIdx || selectedIdx > len(release.Assets) {
		return errors.New("invalid selection")
	}
	selectedAsset := release.Assets[selectedIdx-1]

	filePath := filepath.Join(repoFolder, selectedAsset.Name)

	os.RemoveAll(repoFolder)
	os.MkdirAll(repoFolder, 0755)
	os.MkdirAll(binFolder, 0755)

	fmt.Printf("Selected file: %s, downloading...\n", selectedAsset.Name)
	if err := downloadFile(selectedAsset.URL, filePath); err != nil {
		return err
	}
	fmt.Println("Downloaded successfully.")

	if err := uncompressFile(repoFolder, selectedAsset.Name); err != nil {
		return err
	}

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
		fmt.Print("Your selection: ")
		fmt.Scanln(&selectedStr)
		selectedIdx, err = strconv.Atoi(selectedStr)
		if err != nil {
			return err
		}
	}
	selectedFile, ok := filesMap[selectedIdx]
	if !ok {
		return errors.New("invalid selection")
	}

	err = copyFile(
		filepath.Join(repoFolder, selectedFile),
		filepath.Join(binFolder, filepath.Base(selectedFile)),
	)
	if err != nil {
		return err
	}
	fmt.Println("File installed into bin folder")

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

	_, err = io.Copy(file, resp.Body)
	if err != nil {
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

	err = os.WriteFile(destination, input, 0744)
	if err != nil {
		return err
	}
	return nil
}
