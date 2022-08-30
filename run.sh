#!/bin/bash
set -e

export GHPM_DATABASE_PATH=database.json
export GHPM_BIN_FOLDER=bin
export GHPM_REPO_FOLDER=tmp

go run . "$@"
