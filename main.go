package main

import (
	"flag"
	"fmt"
	"os"
	"repo/repos"
	"repo/ui"
)

func main() {
	token, err := getGithubAccessToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	onlyPublicRepository := flag.Bool("p", false, "Set to true to only list public repositories")

	flag.Parse()
	repos := repos.GetRepos(token, repos.RepoOptions{NumberOfResults: 20}, *onlyPublicRepository)

	ui.List(repos, token)
}
