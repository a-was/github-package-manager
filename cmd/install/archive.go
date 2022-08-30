package install

import (
	"os/exec"
	"path/filepath"
)

func uncompressFile(baseDir, filename string) error {
	ext := filepath.Ext(filename)
	u := uncompress{
		filename: filename,
	}

	switch ext {
	case ".zip":
		u.cmd = []string{"unzip"}
	case ".tar.gz", ".gz", ".tgz", ".tar.xz", ".xz":
		u.cmd = []string{"tar", "xzf"}
	// case ".gz":
	// 	u.cmd = []string{"gunzip"}
	case ".tar.bz2", ".bz2", ".tbz2":
		u.cmd = []string{"tar", "xjf"}
	// case ".bz2":
	// 	u.cmd = []string{"bunzip2"}
	case ".tar":
		u.cmd = []string{"tar", "xf"}
	case ".Z":
		u.cmd = []string{"uncompress"}
	case ".rar":
		u.cmd = []string{"unrar", "x"}
	case ".jar":
		u.cmd = []string{"jar", "-xvf"}
	default:
		// not an archive
		return nil
		// return errors.New("unsupported file extension")
	}

	if err := u.exec(); err != nil {
		return err
	}
	return nil
}

type uncompress struct {
	filename string
	cmd      []string
}

func (u *uncompress) exec() error {
	u.cmd = append(u.cmd, u.filename)
	cmd := exec.Command(u.cmd[0], u.cmd[1:]...)
	cmd.Dir = repoDir
	return cmd.Run()
	// d, err := cmd.CombinedOutput()
	// fmt.Println("uncompress:", string(d), err)
	// return err
}
