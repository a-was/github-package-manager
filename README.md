# github-package-manager

Based on https://gist.github.com/redraw/13ff169741d502b6616dd05dccaa5554

# Installation

### Using `go`
```bash
go install github.com/a-was/github-package-manager@latest
```

### Manually
```bash
git clone https://github.com/a-was/github-package-manager.git
cd github-package-manager
go build .
```

# Configuration

GHPM is configured by environment variables

#### `GHPM_DATABASE_PATH`
Path to database file \
Default `$HOME/.ghpm.json`

#### `GHPM_BIN_FOLDER`
Folder to install downloaded binaries \
It is created by GHPM \
Default `$HOME/bin/`

#### `GHPM_REPO_FOLDER`
Folder to download assets \
It is created by GHPM and cleared on every installation \
Default `$HOME/tmp/`

# Usage

### Get help
```bash
github-package-manager
```

### Install repo
For example to install [bat](https://github.com/sharkdp/bat)
```bash
github-package-manager install sharkdp/bat
```

# Todo
- Database
    - Selected asset?
- Handle .deb / .rpm / .Appimage files
- bin folder
    - ~/.bin?
- `update` command
