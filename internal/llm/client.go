// internal/llm/client.go
package llm

import (
	"sdh-agent/internal/llm/anthropic"
)

// Client defines the interface for any LLM provider
type Client interface {
	// GenerateText sends a prompt with `messages` to the LLM and returns generated text
	GenerateText(messages []string) (string, error)
}

// NewClient creates a new LLM client based on the provider type
func NewClient(apiKey string) Client {
	// For now, default to Anthropic
	return anthropic.NewClient(apiKey)
}
