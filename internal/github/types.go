package github

import (
	"context"
	"fmt"
	"strings"
	"time"

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

// Add a method to get formatted content on demand
func (content *GitHubIssueContent) GetFormattedContent() []string {
	return formatIssueContent(content.Issue, content.Comments)
}

// FormatIssueContent converts a GitHubIssueContent struct into a list of readable strings
// First string contains the main issue, subsequent strings contain individual comments
func formatIssueContent(issue *github.Issue, comments []*github.IssueComment) []string {
	numberOfMessages := len(comments) + 1
	messages := make([]string, 0, numberOfMessages)

	// Format main issue details (Message 1)
	var issueBuilder strings.Builder
	issueBuilder.WriteString(fmt.Sprintf("[Message 1/%d: Main Issue]\n\n", numberOfMessages))
	issueBuilder.WriteString(FormatMainIssue(issue))

	// Add note about comment count
	if len(comments) > 0 {
		issueBuilder.WriteString(fmt.Sprintf("\n\n---\n\nNote: This issue has %d comments that will follow in subsequent messages.", len(comments)))
	}

	messages = append(messages, issueBuilder.String())

	// Format each comment as a separate string
	for i, comment := range comments {
		var commentBuilder strings.Builder
		commentBuilder.WriteString(fmt.Sprintf("[Message %d/%d: Comment %d]\n\n", i+2, numberOfMessages, i+1))

		// Add comment details
		commentBuilder.WriteString(FormatIssueComment(issue, comment, i+1))

		// Add context about where this comment fits in the sequence
		if i < len(comments)-1 {
			commentBuilder.WriteString(fmt.Sprintf("\n\n---\n\nNote: This is comment %d of %d. More comments follow in subsequent messages.", i+1, len(comments)))
		} else {
			commentBuilder.WriteString(fmt.Sprintf("\n\n---\n\nNote: This is the final comment (%d of %d) for this issue.", i+1, len(comments)))
		}

		messages = append(messages, commentBuilder.String())
	}

	return messages
}

// FormatMainIssue formats the main issue details into a readable string
func FormatMainIssue(issue *github.Issue) string {
	var contentBuilder strings.Builder

	if issue.Number != nil && issue.Title != nil {
		contentBuilder.WriteString(fmt.Sprintf("# Issue #%d: %s\n\n", *issue.Number, *issue.Title))
	}
	contentBuilder.WriteString(fmt.Sprintf("## Issue Details\n\n"))

	if issue.State != nil {
		contentBuilder.WriteString(fmt.Sprintf("**State:** %s\n", *issue.State))
	}

	if len(issue.Labels) > 0 {
		contentBuilder.WriteString(fmt.Sprintf("**Labels:** %v\n\n", issue.Labels))
	}

	if issue.Body != nil {
		contentBuilder.WriteString("**Description:**\n")
		contentBuilder.WriteString(*issue.Body)
		contentBuilder.WriteString("\n\n")
	}

	return contentBuilder.String()
}

// FormatIssueComments formats a comment of an issue into a readable string
func FormatIssueComment(issue *github.Issue, comment *github.IssueComment, commentNumber int) string {
	var commentBuilder strings.Builder

	if issue.Number != nil {
		commentBuilder.WriteString(fmt.Sprintf("## Comment %d on Issue #%d\n\n", commentNumber, *issue.Number))
	}

	if comment.User != nil && comment.User.Login != nil {
		commentBuilder.WriteString(fmt.Sprintf("**Author:** %s\n", *comment.User.Login))
	}

	if comment.CreatedAt != nil {
		commentBuilder.WriteString(fmt.Sprintf("**Posted at:** %s\n\n", comment.CreatedAt.Format(time.RFC1123)))
	}

	if comment.Body != nil {
		commentBuilder.WriteString(fmt.Sprintf("**Content:**\n %s\n\n", *comment.Body))
	}

	return commentBuilder.String()
}
