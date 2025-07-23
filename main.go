package main

import (
	"flag"
	"fmt"
	"os"
	"repo/repos"
	"repo/ui"
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
		}
	}

	token, err := getGithubAccessToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	onlyPublicRepository := flag.Bool("p", false, "Set to true to only list public repositories")
	includeOrgRepos := flag.Bool("org", false, "Set to true to include repositories from organizations you're a member of")

	flag.Parse()
	repoList, err := repos.GetRepos(token, repos.RepoOptions{NumberOfResults: 20}, *onlyPublicRepository, *includeOrgRepos)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ui.List(repoList, token)
}
