package repos

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type RepoOptions struct {
	NumberOfResults int
}

func GetRepos(githubAccessToken string, repoOpts RepoOptions, publicOnly bool) ([]*github.Repository, error) {
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

	if publicOnly {
		opts.Visibility = "public"
	}

	repos, resp, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return nil, fmt.Errorf("❌ Invalid GitHub token\n   Your token has expired or is invalid.\n   Please run 'repo login' to re-authenticate")
			case http.StatusForbidden:
				return nil, fmt.Errorf("❌ Insufficient permissions\n   Your GitHub token lacks required permissions.\n   Please run 'repo login' to authenticate with proper scopes")
			}
		}
		return nil, fmt.Errorf("❌ GitHub API error: %w\n   Please check your internet connection or try 'repo login' if the issue persists", err)
	}

	return repos, nil
}
