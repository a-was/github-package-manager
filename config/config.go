package config

import "os"

var (
	DatabasePath = getenv("GHPM_DATABASE_PATH", "~/.ghpm.json")
	BinFolder    = getenv("GHPM_BIN_FOLDER", "~/.bin")
	RepoFolder   = getenv("GHPM_REPO_FOLDER", "~/tmp")
)

func getenv(env, fallback string) string {
	v := os.Getenv(env)
	if v == "" {
		return fallback
	}
	return v
}
