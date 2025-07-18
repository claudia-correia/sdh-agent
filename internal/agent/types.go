package agent

import (
	"sdh-agent/internal/config"
	"sdh-agent/internal/github"
	"sdh-agent/internal/llm"
)

// SDHAgent is the main agent structure
type SDHAgent struct {
	config       config.Configuration
	llmClient    llm.Client
	githubClient *github.Client
}

// AnalyzisResult represents the analysis of a similar issue
type AnalyzisResult struct {
	IssueContent *github.GitHubIssueContent
	Resolution   string
}
