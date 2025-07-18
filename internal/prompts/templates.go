package prompts

import (
	"fmt"
)

// Common context to include at the beginning of all prompts
const SDHContext = `You are a Software Engineer working on the Elastic Cloud offering.
Elastic Cloud is a managed Elasticsearch service that provides Elasticsearch, Kibana, and related services in the cloud.
GitHub SDH (Support Development Help) issues are opened by Support Engineers to get help from the Engineering team on a specific customer issue related to the Elastic Cloud offering.

You are a Software Engineer that is part of the Engineering team, therefore you can help Support Engineers resolve these issues.
Your goal is to analyze the issues you are assigned to and find relevant information in other SDH issues that can help resolve the current issue efficiently.
Focus on technical details and be precise in your analysis.`

// CreateSummaryPrompt creates a prompt to summarize an SDH issue.
func CreateSummaryPrompt() string {
	return fmt.Sprintf(`Analyze the following GitHub SDH issue which you have been assigned to. The issue content includes the initial description posted by Support and follow-up comments.
Provide a concise summary with three specific sections:

1.  **Investigation So Far:** What steps have already been taken to diagnose or fix the problem?
2.  **Established Conclusions:** What facts have been confirmed or ruled out?
3.  **Open Questions:** What specific questions or problems remain unresolved?

Include error messages and any relevant technical details that can help identify similar issues.

The SDH issue content will be provided in follow-up messages.`)
}

// CreateSearchQueriesPrompt creates a prompt to generate GitHub search queries.
func CreateSearchQueriesPrompt(summary string) string {
	return fmt.Sprintf(`Analyze the summary provided below, which corresponds to the GitHub SDH issue you have been assigned to. 
Generate 5 distinct and effective GitHub search query strings that will help find similar SDH issues.
The queries should focus on key error messages, technical components, and problem descriptions.
Do not include explicit IDs in the queries (e.g., allocator IDs or instance IDs).
Do not include "repo:" or "is:closed" filters. 
Do not add quotation marks around the queries.
Output ONLY the queries. 

Format your response with one query per line, with no additional text or formatting.
Each query should be on its own line with no prefixes, numbers, or bullet points.

Example output:
"database connection timeout"
"authentication failure" ORA-12545
"connection pool exhausted" JBoss

SDH issue summary:
%s`, summary)
}

// CreateRelevanceAnalysisPrompt creates a prompt for comparing issues
func CreateRelevanceAnalysisPrompt(mainIssueNumber, otherIssueNumber int) string {
	return fmt.Sprintf(`You have been assigned GitHub SDH issue #%d. 
There is another SDH issue (#%d) that can potentially be related to the current issue.
Analyze the content of both issues and determine if the other issue contains relevant information to help resolve the current issue.

Answer with:
1. RELEVANT: true/false
2. If relevant, provide a summary of how this issue was resolved and what insights it provides.

The content of both SDH issues will be provided in the next two messages.

Format your response as follows with no additional text or formatting:
RELEVANT: [true/false]
RESOLUTION: [result of your analyzis if relevant, or "N/A" if not relevant]`, mainIssueNumber, otherIssueNumber)
}

// CreateReportGenerationPrompt creates the final prompt to generate the full analysis report.
func CreateReportGenerationPrompt(mainIssueNumber int) string {
	return fmt.Sprintf(`You have been assigned GitHub SDH issue #%d. 
You have also collected information from similar resolved issues.
Based on the summary of the current issue plus the information about the remaining similar issues, generate a final report to be posted as a comment on the current GitHub issue. 
The report must be in Markdown format and contain exactly these three sections:

**A. Summary Of Current Issue:**
A summary of the current issue (you can use the same summary that I'll provide you).

**B. Findings From Similar Issues:**
Consolidate the key findings from the similar issues. For each finding, state the information and reference the source GitHub issue URL (e.g., "In issue #123, it was found that...").

**C. Plausible Cause:**
If possible, formulate a clear hypothesis about the likely root cause of the current issue. Base this hypothesis on the outcomes of the similar past issues.

**D. Recommended Actions:**
Provide a clear, actionable, and ordered list of steps to investigate or resolve the issue. These should be concrete actions, such as commands to run, logs to check, specific configurations to verify, or questions for the customer.

Generate only the report content, starting with the first heading.

The main SDH issue summary and the information about the similar issues will be provided in follow-up messages.`, mainIssueNumber)
}
