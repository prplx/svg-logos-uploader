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

func (gc *GithubClient) GetRepositories(ctx context.Context) (*github.RepositoryContent, error) {
	content, _, _, err := gc.client.Repositories.GetContents(ctx, "prplx", "svg-logos", "README.md", nil)
	return content, err
}
