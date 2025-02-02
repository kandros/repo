# repo

This is a command-line application built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea) that allows you to see your recent repositories on GitHub and open, clone, or copy the URL of the repositories.

![image](./image.gif)

## Installation

You can install using one of these methods:

### Option 1: Using install script

```bash
curl -fsSL https://raw.githubusercontent.com/kandros/repo/main/install.sh | bash
```

### Option 2: Manual installation
clone repo and run `make install` (requires Go and make)

## Setup
This app uses the same auth token used by the [official GitHub CLI (gh)](https://cli.github.com/)
