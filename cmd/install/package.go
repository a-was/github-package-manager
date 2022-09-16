package install

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func installPackage(path string) (bool, error) {
	p, _ := filepath.Abs(fmt.Sprintf("%s/%s", repoFolder, path))

	var cmd []string
	switch filepath.Ext(path) {
	case ".deb":
		cmd = []string{"sudo", "apt", "install", p}
	case ".rpm":
		cmd = []string{"sudo", "dnf", "install", p}
	default:
		return false, nil
	}

	fmt.Println("Installing...")

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Dir = repoFolder
	return true, command.Run()
}
