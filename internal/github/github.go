package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func (gc *GithubClient) CreateBranch(ctx context.Context, owner, repo, baseBranch, newBranch string) error {
	ref, _, err := gc.client.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		return err
	}

	newRef := &github.Reference{
		Ref: github.String("refs/heads/" + newBranch),
		Object: &github.GitObject{
			SHA: ref.Object.SHA,
		},
	}

	_, _, err = gc.client.Git.CreateRef(ctx, owner, repo, newRef)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GithubClient) CreateBlob(ctx context.Context, owner, repo string, content []byte) (string, error) {
	strContent := string(content)
	blob, _, err := gc.client.Git.CreateBlob(ctx, owner, repo, &github.Blob{
		Content:  &strContent,
		Encoding: github.String("utf-8"),
	})
	if err != nil {
		return "", err
	}

	return blob.GetSHA(), nil
}

func (gc *GithubClient) CreateTree(ctx context.Context, owner, repo, branch string, filepaths []string) error {
	ref, _, err := gc.client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		return err
	}

	tree := []*github.TreeEntry{}
	for _, file := range filepaths {
		fileContent, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var path string
		if filepath.Ext(file) == ".svg" {
			path = "svg/" + filepath.Base(file)
		} else {
			path = filepath.Base(file)
		}

		treeEntry := &github.TreeEntry{
			Path:    github.String(path),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String(string(fileContent)),
		}
		tree = append(tree, treeEntry)
	}

	newTree, _, err := gc.client.Git.CreateTree(ctx, owner, repo, *ref.Object.SHA, tree)
	if err != nil {
		return err
	}

	newCommit, _, err := gc.client.Git.CreateCommit(ctx, owner, repo, &github.Commit{
		Message: github.String(GenerateCommitMessageFromUploadedFiles(filepaths)),
		Parents: []*github.Commit{{SHA: ref.Object.SHA}},
		Tree:    newTree,
	}, nil)
	if err != nil {
		return err
	}

	// Update the branch to point to the new commit
	ref = &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", branch)),
		Object: &github.GitObject{SHA: newCommit.SHA},
	}
	_, _, err = gc.client.Git.UpdateRef(ctx, owner, repo, ref, false)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GithubClient) CreatePullRequest(ctx context.Context, owner, repo, baseBranch, newBranch, title string) error {
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(baseBranch),
		Base:  github.String(newBranch),
	}

	_, _, err := gc.client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		return err
	}

	return nil
}

func GenerateBranchNameFromUploadedFiles(filePaths []string) string {
	builder := strings.Builder{}
	builder.WriteString("add")
	for _, path := range filePaths {
		fileName := filepath.Base(path)
		builder.WriteString("-" + strings.ToLower(strings.TrimSuffix(fileName, filepath.Ext(fileName))))
	}
	return builder.String()
}

func GenerateCommitMessageFromUploadedFiles(filePaths []string) string {
	builder := strings.Builder{}
	builder.WriteString("Add")
	for _, path := range filePaths {
		fileName := filepath.Base(path)
		if filepath.Ext(fileName) == ".svg" {
			builder.WriteString(" " + strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ",")
		} else {
			builder.WriteString(" update " + fileName)
		}
	}
	return builder.String()
}
