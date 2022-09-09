package config

import (
	"errors"
	"os"
	"regexp"
)

var (
	DatabasePath = getenv("GHPM_DATABASE_PATH", "$HOME/.ghpm.json")
	BinFolder    = getenv("GHPM_BIN_FOLDER", "$HOME/bin")
	RepoFolder   = getenv("GHPM_REPO_FOLDER", "$HOME/tmp")
)

func getenv(env, fallback string) string {
	v := os.Getenv(env)
	if v == "" {
		v = fallback
	}
	return os.ExpandEnv(v)
}

var (
	RegexRepo          = regexp.MustCompile(`^[a-zA-Z0-9\-]+/[a-zA-Z0-9\-]+$`)
	RegexSearchPattern = regexp.MustCompile(`^[a-zA-Z0-9\-/]+$`)

	ErrorAborted        = errors.New("aborted")
	ErrorApiWrongStatus = errors.New("API returned wrong status")
)
