package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"repo/repos"
	"repo/ui"
	"runtime"
	"strings"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("Version: %s\n", version)
			return
		case "login":
			if err := runLogin(); err != nil {
				fmt.Fprintln(os.Stderr, "Login failed:", err)
				os.Exit(1)
			}
			return
		case "update":
			if err := runUpdate(); err != nil {
				fmt.Fprintln(os.Stderr, "Update failed:", err)
				os.Exit(1)
			}
			return
		}
	}

	token, err := getGithubAccessToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	allowPrivate := flag.Bool("p", false, "Include private repositories in the list")
	allowPrivateAlias := flag.Bool("private", false, "Include private repositories in the list")
	includeOrgRepos := flag.Bool("A", false, "Include repositories from organizations you're a member of")
	includeOrgReposAlias := flag.Bool("all", false, "Include repositories from organizations you're a member of")

	flag.Parse()
	
	// Use either flag for allowPrivate
	showPrivate := *allowPrivate || *allowPrivateAlias
	// Use either flag for includeOrgRepos
	showOrgRepos := *includeOrgRepos || *includeOrgReposAlias
	
	repoList, err := repos.GetRepos(token, repos.RepoOptions{NumberOfResults: 20}, showPrivate, showOrgRepos)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ui.List(repoList, token)
}

type Release struct {
	TagName string `json:"tag_name"`
}

func runUpdate() error {
	fmt.Println("Checking for updates...")
	
	// Get the latest release version
	resp, err := http.Get("https://api.github.com/repos/kandros/repo/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release info: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(version, "v")

	if latestVersion == currentVersion || currentVersion == "dev" {
		if currentVersion == "dev" {
			fmt.Println("You are running a development version.")
		} else {
			fmt.Printf("You are already running the latest version (%s)\n", currentVersion)
		}
		return nil
	}

	fmt.Printf("Updating from %s to %s...\n", currentVersion, latestVersion)

	// Determine OS and architecture
	osName := strings.ToLower(runtime.GOOS)
	arch := runtime.GOARCH

	// Convert Go arch names to match release names
	switch arch {
	case "amd64":
		// keep as is
	case "386":
		// keep as is  
	case "arm64":
		// keep as is
	case "arm":
		// keep as is
	default:
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	// Download the new version
	downloadURL := fmt.Sprintf("https://github.com/kandros/repo/releases/download/v%s/repo_%s_%s_%s.tar.gz", 
		latestVersion, latestVersion, osName, arch)

	fmt.Printf("Downloading %s...\n", downloadURL)
	
	resp, err = http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "repo-update")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract the tar.gz
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var newBinaryPath string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar archive: %w", err)
		}

		if header.Name == "repo" || header.Name == "repo.exe" {
			newBinaryPath = filepath.Join(tempDir, header.Name)
			file, err := os.Create(newBinaryPath)
			if err != nil {
				return fmt.Errorf("failed to create temp binary: %w", err)
			}

			_, err = io.Copy(file, tr)
			file.Close()
			if err != nil {
				return fmt.Errorf("failed to extract binary: %w", err)
			}

			// Make executable
			if err := os.Chmod(newBinaryPath, 0755); err != nil {
				return fmt.Errorf("failed to make binary executable: %w", err)
			}
			break
		}
	}

	if newBinaryPath == "" {
		return fmt.Errorf("binary not found in release archive")
	}

	// Replace current binary
	fmt.Println("Installing update...")
	
	// On Windows, we can't replace a running executable directly
	if runtime.GOOS == "windows" {
		// Create a batch script to replace the binary after this process exits
		batchScript := filepath.Join(tempDir, "update.bat")
		batchContent := fmt.Sprintf(`@echo off
timeout /t 1 /nobreak >nul
move "%s" "%s"
echo Update completed successfully!
pause
del "%%~f0"
`, newBinaryPath, execPath)
		
		if err := os.WriteFile(batchScript, []byte(batchContent), 0755); err != nil {
			return fmt.Errorf("failed to create update script: %w", err)
		}
		
		cmd := exec.Command("cmd", "/C", "start", "/MIN", batchScript)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start update script: %w", err)
		}
		
		fmt.Println("Update will complete after this process exits.")
		return nil
	}

	// On Unix-like systems, replace the binary directly
	if err := os.Rename(newBinaryPath, execPath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	fmt.Printf("Successfully updated to version %s!\n", latestVersion)
	return nil
}
