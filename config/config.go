package config

import "os"

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
