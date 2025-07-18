package agent

import (
	"fmt"
	"log"

	"sdh-agent/internal/config"
	"sdh-agent/internal/github"
	"sdh-agent/internal/llm"
)

// NewSDHAgent creates a new SDH agent instance
func NewSDHAgent(config config.Configuration) *SDHAgent {
	// Initialize API clients
	githubClient := github.NewClient(config.GitHubToken)
	llmClient := llm.NewClient(config.LlmApiKey)

	return &SDHAgent{
		config:       config,
		llmClient:    llmClient,
		githubClient: githubClient,
	}
}

// ProcessIssue executes the full workflow for a given SDH issue
func (agent *SDHAgent) ProcessIssue(issueNumber int) (string, error) {
	log.Printf("Starting to process SDH issue #%d", issueNumber)

	// Ingest active SDH issue
	issueContent, err := agent.githubClient.GetIssueContent(agent.config.GitHubRepoOwner, agent.config.GitHubRepoName, issueNumber)
	if err != nil {
		return "", fmt.Errorf("failed to ingest main SDH issue #%d: %w", issueNumber, err)
	}

	// Summarize the issue
	log.Printf("Summarizing content for SDH issue")
	summary, err := agent.summarizeIssueContent(issueContent)
	if err != nil {
		return "", fmt.Errorf("failed to summarize issue content: %w", err)
	}
	log.Printf("Summary for SDH issue #%d:\n %s", issueNumber, summary)

	// Analyze similar issues
	log.Printf("Analyzing similar issues")
	analysisResults, err := agent.analyzeSimilarIssues(issueContent, summary)
	if err != nil {
		return "", fmt.Errorf("failed to analyze similar issues: %w", err)
	}

	// Generate report
	log.Printf("Generating final report")
	report, err := agent.generateReport(issueContent, summary, analysisResults)
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	log.Printf("Successfully processed SDH issue #%d", issueNumber)
	return report, nil
}
