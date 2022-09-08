package install

import (
	"fmt"
	"os/exec"
	"strings"
)

func uncompressFile(baseFolder, filename string) error {
	var cmd []string
	var options []string
	parts := strings.Split(filename, ".")
	ln := len(parts)
	switch {
	case ln >= 3:
		options = []string{
			fmt.Sprintf(".%s.%s", parts[ln-2], parts[ln-1]),
			fmt.Sprintf(".%s", parts[ln-1]),
		}
	case ln >= 2:
		options = []string{
			fmt.Sprintf(".%s", parts[ln-1]),
		}
	}

	for _, option := range options {
		switch option {
		case ".zip":
			cmd = []string{"unzip"}
		case ".tar.gz", ".tgz", ".tar.xz", ".xz":
			cmd = []string{"tar", "xzf"}
		case ".gz":
			cmd = []string{"gunzip"}
		case ".tar.bz", ".tbz", ".tar.bz2", ".tbz2":
			cmd = []string{"tar", "xjf"}
		case ".bz", ".bz2":
			cmd = []string{"bunzip2"}
		case ".tar":
			cmd = []string{"tar", "xf"}
		case ".Z":
			cmd = []string{"uncompress"}
		case ".rar":
			cmd = []string{"unrar", "x"}
		case ".jar":
			cmd = []string{"jar", "-xvf"}
		default:
			continue
		}
		break
	}
	if cmd == nil {
		return nil
	}
	cmd = append(cmd, filename)

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Dir = repoFolder
	err := command.Run()
	// d, err := command.CombinedOutput()
	// fmt.Println("uncompress:", string(d), err)
	if err != nil {
		return err
	}

	return nil
}
