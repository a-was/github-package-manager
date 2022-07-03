package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	repoDir = "tmp"
	binDir  = "bin"
)

func Install(repo string) error {
	fmt.Println("Using repo", repo)
	release, err := getLatestRelease(repo)
	if err != nil {
		return err
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

	if err := uncompressFile(filePath); err != nil {
		return err
	}

	return nil
}
