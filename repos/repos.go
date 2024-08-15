package repos

import (
	"context"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type RepoOptions struct {
	NumberOfResults int
}

func GetRepos(githubAccessToken string, repoOpts RepoOptions) []*github.Repository {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryListOptions{
		Affiliation: "owner,organization_member",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{Page: 1, PerPage: repoOpts.NumberOfResults},
	}

	repos, _, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		panic(err)
	}

	return repos
}
