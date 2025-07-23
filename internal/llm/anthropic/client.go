package anthropic

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"sdh-agent/internal/prompts"
	"sdh-agent/pkg/utils"

	"golang.org/x/time/rate"
)

const (
	apiURL = "https://api.anthropic.com/v1/messages"
	// Using the latest model available at the time of writing.
	// This might need updating in the future.
	model = "claude-3-5-haiku-latest"

	// Rate limiting configurations
	tokensPerMinute  = 19000 // Anthropic's limit is 20000 tokens per minute, we use a conservative estimate
	avgCharsPerToken = 4.0   // Average characters per token for English text
	baseTokens       = 3     // Base tokens per message (conservative estimate)
)

// Client is a wrapper for the Anthropic API
type Client struct {
	apiKey      string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	// Add backoff configuration
	maxRetries  int
	baseBackoff time.Duration
}

// NewClient creates a new Anthropic API client.
func NewClient(apiKey string) *Client {
	// Create a rate limiter with `tokensPerMinute` tokens per minute (slightly under the limit)
	// and burst of `tokensPerMinute` tokens maximum
	limiter := rate.NewLimiter(rate.Limit(tokensPerMinute/60), tokensPerMinute)

	return &Client{
		apiKey:      apiKey,
		httpClient:  utils.CreateDefaultHTTPClient(),
		rateLimiter: limiter,
		maxRetries:  5,
		baseBackoff: 15 * time.Second, // Base backoff of 15 seconds
	}
}

// GenerateText sends a request to the Anthropic API and returns the generated text
func (c *Client) GenerateText(messages []string) (string, error) {
	// Wait for rate limiter (estimate 1 token per character as a conservative approach)
	ctx := context.Background()
	estimatedTokens := estimateTokenCount(messages)
	if err := c.rateLimiter.WaitN(ctx, estimatedTokens); err != nil {
		return "", fmt.Errorf("rate limiter wait error: %w", err)
	}

	// Attempt request with retries and exponential backoff
	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		response, err := c.makeRequest(messages)
		if err == nil {
			// Success!
			return response, nil
		}

		// Check if error is a rate limit error
		if isRateLimitError(err) {
			lastErr = err

			// Calculate backoff duration with jitter
			backoffDuration := c.calculateBackoff(attempt)
			time.Sleep(backoffDuration)
			continue
		}

		// Not a rate limit error, return immediately
		return "", err
	}

	// All retries failed
	return "", fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRateLimitError checks if an error is a rate limit error
func isRateLimitError(err error) bool {
	// Check if error message contains rate limit indicators
	errMsg := err.Error()
	return err != nil && (strings.Contains(errMsg, "429") ||
		strings.Contains(errMsg, "rate_limit_error") ||
		strings.Contains(errMsg, "exceed the rate limit"))
}

// calculateBackoff calculates the backoff duration with jitter
func (c *Client) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: baseBackoff * 2^attempt + random jitter
	backoff := float64(c.baseBackoff) * math.Pow(2, float64(attempt))

	// Add jitter (±20%)
	jitter := backoff * (0.8 + 0.4*rand.Float64())

	return time.Duration(jitter)
}

// makeRequest makes the actual HTTP request to the Anthropic API
// This would be your existing request code
func (c *Client) makeRequest(messages []string) (string, error) {

	reqBody := anthropicRequest{
		Model:     model,
		Messages:  convertToMessages(messages),
		MaxTokens: 4096,               // Max output tokens
		System:    prompts.SDHContext, // Include the system context
	}

	headers := map[string]string{
		"x-api-key":         c.apiKey,
		"anthropic-version": "2023-06-01",
	}

	var anthropicResp anthropicResponse
	err := utils.SendJSONRequest(
		c.httpClient,
		"POST",
		apiURL,
		reqBody,
		&anthropicResp,
		headers,
	)
	if err != nil {
		return "", err
	}

	if anthropicResp.Error.Message != "" {
		return "", fmt.Errorf("anthropic API error: %s - %s", anthropicResp.Error.Type, anthropicResp.Error.Message)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("received empty content from Anthropic API")
	}

	return anthropicResp.Content[0].Text, nil
}

// convertToMessages converts an array of strings to an array of Messages with the "user" role
func convertToMessages(contents []string) []Message {
	messages := make([]Message, len(contents))

	for i, content := range contents {
		messages[i] = Message{
			Role:    "user",
			Content: content,
		}
	}

	return messages
}

// estimateTokenCount estimates the number of tokens in a slice of strings
func estimateTokenCount(texts []string) int {
	totalTokens := 0

	for _, text := range texts {
		if text == "" {
			continue
		}

		// Basic estimation based on character count
		charCount := len(text)

		// Add base tokens plus character-based estimation
		// Use math.Ceil to round up for a conservative estimate
		messageTokens := baseTokens + int(math.Ceil(float64(charCount)/avgCharsPerToken))

		totalTokens += messageTokens
	}

	return totalTokens
}
