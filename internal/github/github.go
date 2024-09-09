package github

import (
	"context"

	"github.com/google/go-github/v64/github"
)

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(token string) *GithubClient {
	client := github.NewClient(nil).WithAuthToken(token)

	return &GithubClient{
		client: client,
	}
}

func (gc *GithubClient) GetRepositoryContent(ctx context.Context, user, repo, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	content, dirContent, _, err := gc.client.Repositories.GetContents(ctx, user, repo, path, nil)
	return content, dirContent, err
}
