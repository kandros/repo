package main

import (
	"repo/repos"
	"repo/ui"
)

func main() {
	token := getGithubAccessToken()

	repos := repos.GetRepos(token)

	ui.List(repos)
}
