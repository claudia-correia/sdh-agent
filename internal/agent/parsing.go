package agent

import "strings"

// parseRelevanceResponse extracts relevance and resolution from LLM response
func parseRelevanceResponse(response string) (bool, string) {
	// Split the response into sections based on newlines
	parts := strings.SplitN(response, "\n", 2)

	// Default values
	relevant := false
	resolution := ""

	// Parse relevance from first line
	if len(parts) > 0 && strings.HasPrefix(parts[0], "RELEVANT:") {
		relevant = strings.Contains(strings.ToLower(parts[0]), "true")
	}

	// Parse resolution from second part (which may contain multiple lines)
	if len(parts) > 1 && strings.HasPrefix(parts[1], "RESOLUTION:") {
		// Remove the "RESOLUTION:" prefix
		resolution = strings.TrimSpace(strings.TrimPrefix(parts[1], "RESOLUTION:"))
	}

	return relevant, resolution
}

// parseSearchQueries extracts individual search queries from LLM response
func parseSearchQueries(response string) []string {
	var queries []string
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines
		if line == "" {
			continue
		}
		queries = append(queries, line)
	}

	return queries
}
