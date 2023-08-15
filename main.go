package main

import (
	"repo/repos"
	"repo/ui"
)

func main() {
	token := getGithubAccessToken()

	repos := repos.GetRepos(token, repos.RepoOptions{NumberOfResults: 20})

	ui.List(repos)
}
