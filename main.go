package main

import (
	"flag"
	"repo/repos"
	"repo/ui"
)

func main() {
	token := getGithubAccessToken()
	onlyPublicRepository := flag.Bool("p", false, "Set to true to only list public repositories")

	flag.Parse()
	repos := repos.GetRepos(token, repos.RepoOptions{NumberOfResults: 20}, *onlyPublicRepository)

	ui.List(repos)
}
