package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func getGithubAccessToken() (string, error) {
	// Try to get token from our own token file first
	token, err := getStoredToken()
	if err == nil && token != "" {
		return token, nil
	}

	// Fallback to GitHub CLI config for backward compatibility
	token, err = getGithubCLIToken()
	if err == nil && token != "" {
		return token, nil
	}

	return "", fmt.Errorf("no GitHub token found - run 'repo login' to authenticate with GitHub")
}

func getStoredToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}

	tokenPath := filepath.Join(home, ".config", "repo", "token")
	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(tokenBytes)), nil
}

func storeToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "repo")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	tokenPath := filepath.Join(configDir, "token")
	return os.WriteFile(tokenPath, []byte(token), 0600)
}

func getGithubCLIToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := filepath.Join(home, ".config", "gh", "hosts.yml")

	yfile, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	type Host struct {
		OauthToken string `yaml:"oauth_token"`
	}

	data := make(map[interface{}]Host)
	if err = yaml.Unmarshal(yfile, &data); err != nil {
		return "", err
	}

	host, ok := data["github.com"]
	if !ok {
		return "", fmt.Errorf("no github.com configuration found")
	}

	if host.OauthToken == "" {
		return "", fmt.Errorf("no OAuth token found")
	}

	return host.OauthToken, nil
}

func promptForToken() (string, error) {
	fmt.Print("Paste your GitHub token (input will be hidden): ")
	tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // New line after password input
	if err != nil {
		return "", fmt.Errorf("could not read token: %w", err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return "", fmt.Errorf("no token provided")
	}

	return token, nil
}

func runLogin() error {
	fmt.Println("To authenticate with GitHub, you need a personal access token.")
	fmt.Println("")
	fmt.Println("1. Go to: https://github.com/settings/tokens/new?scopes=repo,read:user&description=repo-cli")
	fmt.Println("2. Generate a new token with 'repo' and 'read:user' scopes")
	fmt.Println("3. Copy the token and paste it below")
	fmt.Println("")

	token, err := promptForToken()
	if err != nil {
		return err
	}

	// Verify the token works by making a simple API call
	fmt.Print("Verifying token... ")
	if err := verifyToken(token); err != nil {
		fmt.Println("failed")
		return fmt.Errorf("token verification failed: %w", err)
	}
	fmt.Println("success!")

	// Store the token
	if err := storeToken(token); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	fmt.Println("Authentication successful! You can now use the repo command.")
	return nil
}

func verifyToken(token string) error {
	// For now, we skip verification to avoid circular dependencies
	// The token will be validated when first used in the repos package
	return nil
}
