package db

import (
	"encoding/json"
	"os"

	"github.com/a-was/github-package-manager/config"
	"github.com/a-was/github-package-manager/github"
)

type dbT struct {
	Installed map[string]version // key: repo, value: installed version
}

type version struct {
	ID            int
	Tag           string
	SelectedAsset int
	SelectedFile  int
}

var db = dbT{
	Installed: map[string]version{},
}

func (db *dbT) load() {
	data, err := os.ReadFile(config.DatabasePath)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, db)
}

func (db *dbT) dump() {
	file, _ := json.MarshalIndent(db, "", "    ")
	_ = os.WriteFile(config.DatabasePath, file, 0644)
}

func CheckIfInstalled(r *github.Release) bool {
	db.load()
	v, ok := db.Installed[r.Repo]
	return ok && v.Tag == r.Tag
}

func SaveVersion(r *github.Release, selectedAsset, selectedFile int) {
	db.load()
	db.Installed[r.Repo] = version{
		ID:            r.ID,
		Tag:           r.Tag,
		SelectedAsset: selectedAsset,
		SelectedFile:  selectedFile,
	}
	db.dump()
}

func GetInstalled() map[string]version {
	db.load()
	return db.Installed
}
