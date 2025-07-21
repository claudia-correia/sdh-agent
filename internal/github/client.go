package github

import (
	"fmt"

	"github.com/google/go-github/v63/github"
)

// GetIssueContent fetches an issue and all its comments
func (c *Client) GetIssueContent(owner, repo string, issueNumber int) (*GitHubIssueContent, error) {
	// Fetch the issue
	issue, _, err := c.client.Issues.Get(c.ctx, owner, repo, issueNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	// Fetch comments
	comments, err := c.GetIssueComments(owner, repo, issue)
	if err != nil {
		return nil, err
	}

	return &GitHubIssueContent{
		IssueNumber: issueNumber,
		Issue:       issue,
		Comments:    comments,
	}, nil
}

// GetIssueContent fetches an issue and all its comments
func (c *Client) GetIssueComments(owner, repo string, issue *github.Issue) ([]*github.IssueComment, error) {
	if issue == nil || issue.Number == nil {
		return nil, fmt.Errorf("issue is nil or does not have a number")
	}

	comments, _, err := c.client.Issues.ListComments(c.ctx, owner, repo, *issue.Number, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments for issue #%d: %w", *issue.Number, err)
	}

	return comments, nil
}

// SearchIssues searches for issues in the repository.
func (c *Client) SearchIssues(owner, repo, query string) ([]*github.Issue, error) {
	fullQuery := fmt.Sprintf("%s repo:%s/%s is:issue is:closed", query, owner, repo)
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			// Use GitHub's default "best-match" search algorithm
			PerPage: 20, // Limit to top 50 relevant issues
		},
	}

	result, _, err := c.client.Search.Issues(c.ctx, fullQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %w", err)
	}
	return result.Issues, nil
}

// PostComment posts a comment to a GitHub issue.
func (c *Client) PostComment(owner, repo string, issueNumber int, body string) error {
	comment := &github.IssueComment{Body: &body}
	_, _, err := c.client.Issues.CreateComment(c.ctx, owner, repo, issueNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}
	return nil
}
