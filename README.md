# GitHub Package Manager

Based on https://gist.github.com/redraw/13ff169741d502b6616dd05dccaa5554

# Installation

### Using asset from release
Just download binary from latest release

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

### Using github-package-manager
```bash
github-package-manager install a-was/github-package-manager
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

### Search for repo
For example to search for a [bat](https://github.com/sharkdp/bat) (to see which username to use)
```bash
github-package-manager search bat
```

### Install repo
For example to install [bat](https://github.com/sharkdp/bat)
```bash
github-package-manager install sharkdp/bat
```

### Update all installed repos
```bash
github-package-manager update
```

# Todo
- Handle .deb / .rpm / .Appimage files
- docs
    - Long command descriptions

# License
[MIT](https://choosealicense.com/licenses/mit/)
