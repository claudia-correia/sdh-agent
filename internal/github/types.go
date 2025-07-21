package github

import (
	"context"

	"github.com/google/go-github/v63/github"
	"golang.org/x/oauth2"
)

// Client is a wrapper around the go-github client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// NewClient creates a new GitHub API client
func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &Client{
		client: github.NewClient(tc),
		ctx:    ctx,
	}
}

// GitHubIssueContent represents a GitHub issue along with its comments
type GitHubIssueContent struct {
	IssueNumber int
	Issue       *github.Issue
	Comments    []*github.IssueComment
}

// GetLabels returns the labels of the issue as a map
func (content *GitHubIssueContent) GetLabels() map[string]bool {
	mainLabels := make(map[string]bool)
	for _, label := range content.Issue.Labels {
		mainLabels[label.GetName()] = true
	}

	return mainLabels
}
