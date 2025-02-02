package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func getGithubAccessToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := home + "/.config/gh/hosts.yml"

	yfile, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("could not read auth file at %s - you need either:\n1. GitHub CLI installed and authenticated ('gh auth login'), or\n2. A manual auth file at this path with your token", configPath)
	}

	type Host struct {
		OauthToken string `yaml:"oauth_token"`
	}

	data := make(map[interface{}]Host)
	if err = yaml.Unmarshal(yfile, &data); err != nil {
		return "", fmt.Errorf("could not parse auth file - file should be in YAML format with a github.com entry containing an oauth_token")
	}

	host, ok := data["github.com"]
	if !ok {
		return "", fmt.Errorf("no github.com configuration found in auth file - file should contain a github.com entry with an oauth_token")
	}

	if host.OauthToken == "" {
		return "", fmt.Errorf("no OAuth token found in auth file - ensure github.com entry contains a valid oauth_token")
	}

	return host.OauthToken, nil
}
