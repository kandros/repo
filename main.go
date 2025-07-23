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
