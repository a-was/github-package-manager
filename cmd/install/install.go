package install

import (
	"fmt"
	"github-package-manager/db"
	"github-package-manager/github"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	repoDir = "tmp"
	binDir  = "bin"
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
	selectedAsset := release.Assets[selectedIdx-1]

	filePath := filepath.Join(repoDir, selectedAsset.Name)

	os.MkdirAll(repoDir, 0755)
	os.MkdirAll(binDir, 0755)

	fmt.Printf("Selected file: %s, downloading...\n", selectedAsset.Name)
	if err := downloadFile(selectedAsset.URL, filePath); err != nil {
		return err
	}
	fmt.Println("Downloaded successfully.")

	if err := uncompressFile(repoDir, selectedAsset.Name); err != nil {
		return err
	}

	i := 1
	filesMap := map[int]string{}
	filepath.Walk(repoDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		p := strings.TrimPrefix(path, fmt.Sprintf("%s/", repoDir))
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
	selectedFile := filesMap[selectedIdx]

	// os.Rename(
	// 	filepath.Join(repoDir, selectedFile),
	// 	filepath.Join(binDir, filepath.Base(selectedFile)),
	// )
	copyFile(
		filepath.Join(repoDir, selectedFile),
		filepath.Join(binDir, filepath.Base(selectedFile)),
	)
	fmt.Println("File installed into bin folder")

	// os.RemoveAll(repoDir)

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

	file.Close()
	os.Rename(tmpFilepath, filepath)
	return nil
}

func copyFile(source, destination string) {
	input, err := os.ReadFile(source)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(destination, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destination)
		fmt.Println(err)
		return
	}
}
