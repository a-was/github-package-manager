package install

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	fmt.Println("File moved to bin folder")

	// os.RemoveAll(repoDir)

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
