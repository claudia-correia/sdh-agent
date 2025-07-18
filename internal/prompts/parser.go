package prompts

import "strings"

// ParseRelevanceResponse extracts relevant status and resolution from LLM's response
func ParseRelevanceResponse(response string) (bool, string) {
	lines := strings.Split(response, "\n")
	relevant := false
	resolution := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "RELEVANT:") {
			relevant = strings.Contains(strings.ToLower(line), "true")
		} else if strings.HasPrefix(line, "RESOLUTION:") {
			resolution = strings.TrimSpace(strings.TrimPrefix(line, "RESOLUTION:"))
		}
	}

	return relevant, resolution
}
