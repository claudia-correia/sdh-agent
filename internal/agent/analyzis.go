package agent

import (
	"fmt"
	"log"
	"sort"
	"time"

	"sdh-agent/internal/github"
	"sdh-agent/internal/prompts"
)

// analyzeSimilarIssues analyzes each similar issue for relevance
func (agent *SDHAgent) analyzeSimilarIssues(mainIssue *github.GitHubIssueContent, mainSummary string) ([]AnalyzisResult, error) {
	var results []AnalyzisResult

	// Identify and filter similar issues
	similarIssues, err := agent.findSimilarIssues(mainIssue, mainSummary)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar issues: %w", err)
	}

	// Sort the fetched issues by metadata score
	sort.Slice(similarIssues, func(i, j int) bool {
		return scoreIssueByMetadata(mainIssue, similarIssues[i]) > scoreIssueByMetadata(mainIssue, similarIssues[j])
	})

	for _, issue := range similarIssues {
		// Limit analysis to prevent token overflow
		if len(results) >= 10 {
			break
		}

		// Analyze relevance
		relevance, resolution, err := agent.analyzeIssueRelevance(mainSummary, mainIssue, issue)
		if err != nil {
			log.Printf("Error analyzing issue #%d: %v", issue.IssueNumber, err)
			continue
		}

		if relevance {
			log.Printf("Issue #%d is relevant: %s", issue.IssueNumber, resolution)
			results = append(results, AnalyzisResult{
				IssueContent: issue,
				Resolution:   resolution,
			})
		}
	}

	return results, nil
}

// analyzeIssueRelevance determines if an issue is relevant
func (agent *SDHAgent) analyzeIssueRelevance(mainSummary string, mainIssue, similarIssue *github.GitHubIssueContent) (bool, string, error) {
	log.Printf("Analyzing relevance for issue #%d", similarIssue.IssueNumber)

	var messages []string

	// Create a prompt for relevance analysis
	prompt := prompts.CreateRelevanceAnalysisPrompt(mainIssue.IssueNumber, similarIssue.IssueNumber)
	messages = append(messages, prompt)

	// Add main issue summary
	messages = append(messages, fmt.Sprintf("Main Issue Summary:\n%s", mainSummary))

	// Add similar issue content
	messages = append(messages, fmt.Sprintf("Similar Issue Content:\n%s", formatIssueContent(similarIssue)))

	response, err := agent.llmClient.GenerateText(messages)
	if err != nil {
		return false, "", err
	}

	log.Printf("Relevance analysis response for issue #%d: %s", similarIssue.IssueNumber, response)

	// Parse the response
	relevant, resolution := parseRelevanceResponse(response)

	return relevant, resolution, nil
}

// scoreIssueByMetadata Provides a basic scoring mechanism based on issue metadata
func scoreIssueByMetadata(mainIssue, otherIssue *github.GitHubIssueContent) float64 {
	score := 0.0

	// Return 0 if either issue or its content is nil
	if mainIssue == nil || otherIssue == nil || mainIssue.Issue == nil || otherIssue.Issue == nil {
		return score // Return 0 if either issue is nil
	}

	// Label similarity score: Adds 0.2 points for each matching label
	if mainIssue.Issue.Labels != nil && otherIssue.Issue.Labels != nil {
		mainLabels := mainIssue.GetLabels()
		for _, label := range otherIssue.Issue.Labels {
			if mainLabels[label.GetName()] {
				score += 0.2
			}
		}
	}

	// Engagement score: Adds 1.0 point if the issue has more than 5 comments
	if otherIssue.Issue.Comments != nil && *otherIssue.Issue.Comments > 5 {
		score += 1.0
	}

	// Recency bonus: Adds 1.0 point if the issue was closed within the last year
	if otherIssue.Issue.ClosedAt != nil {
		oneYearAgo := time.Now().AddDate(-1, 0, 0)
		if otherIssue.Issue.ClosedAt.Time.After(oneYearAgo) {
			score += 1.0
		}
	}

	return score
}

// findSimilarIssues searches for related issues
func (agent *SDHAgent) findSimilarIssues(mainIssue *github.GitHubIssueContent, summary string) ([]*github.GitHubIssueContent, error) {
	log.Printf("Searching for similar issues")

	// Extract search terms from the issue
	searchQueries := agent.extractSearchQueries(summary)

	var allIssues []*github.GitHubIssueContent
	seenIssues := make(map[int]bool)

	for _, query := range searchQueries {
		log.Printf("Searching with query '%s'", query)
		results, err := agent.githubClient.SearchIssues(agent.config.GitHubRepoOwner, agent.config.GitHubRepoName, query)
		if err != nil {
			log.Printf("Error searching with query '%s': %v", query, err)
			continue
		}

		for _, issue := range results {
			if issue.Number != nil && *issue.Number != mainIssue.IssueNumber && !seenIssues[*issue.Number] {
				// Ingest similar issue
				comments, err := agent.githubClient.GetIssueComments(agent.config.GitHubRepoOwner, agent.config.GitHubRepoName, issue)

				if err != nil {
					log.Printf("Error ingesting comments for issue #%d: %v", *issue.Number, err)
					continue
				}

				issueContent := &github.GitHubIssueContent{
					IssueNumber: *issue.Number,
					Issue:       issue,
					Comments:    comments,
				}

				// Ensure the issue is not already processed
				seenIssues[*issue.Number] = true
				allIssues = append(allIssues, issueContent)
			}
		}
	}

	return allIssues, nil
}

// extractSearchQueries generates search queries from the issue using LLM
func (agent *SDHAgent) extractSearchQueries(summary string) []string {
	// Create prompt for LLM to generate search queries
	prompt := prompts.CreateSearchQueriesPrompt(summary)

	// Get response from LLM
	response, err := agent.llmClient.GenerateText([]string{prompt})
	if err != nil {
		log.Printf("Error generating search queries: %v", err)
		return []string{} // Return empty slice if LLM fails
	}

	// Parse the response into individual queries
	queries := parseSearchQueries(response)

	// Limit queries to prevent excessive API calls
	if len(queries) > 5 {
		queries = queries[:5]
	}

	return queries
}

// summarizeIssueContent uses an LLM to summarize the issue
func (agent *SDHAgent) summarizeIssueContent(issueContent *github.GitHubIssueContent) (string, error) {
	var messages []string

	// Create a prompt for summarization
	prompt := prompts.CreateSummaryPrompt()
	messages = append(messages, prompt)
	messages = append(messages, formatIssueContent(issueContent)...)

	response, err := agent.llmClient.GenerateText(messages)
	if err != nil {
		return "", err
	}

	return response, nil
}
