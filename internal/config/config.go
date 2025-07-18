package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Configuration holds all necessary API configurations
type Configuration struct {
	GitHubToken     string
	GitHubRepoOwner string
	GitHubRepoName  string
	LlmApiKey       string
}

// Load configuration from environment variables
// It looks for a .env file for local development
func Load() (*Configuration, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Create config with values from environment
	config := &Configuration{
		GitHubToken:     os.Getenv("GITHUB_TOKEN"),
		LlmApiKey:       os.Getenv("LLM_API_KEY"),
		GitHubRepoOwner: os.Getenv("GITHUB_REPO_OWNER"),
		GitHubRepoName:  os.Getenv("GITHUB_REPO_NAME"),
	}

	// Validate the loaded configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate ensures all required configuration values are set
func (c *Configuration) Validate() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	if c.LlmApiKey == "" {
		return fmt.Errorf("LLM_API_KEY environment variable not set")
	}

	if c.GitHubRepoOwner == "" {
		return fmt.Errorf("GITHUB_REPO_OWNER environment variable not set")
	}

	if c.GitHubRepoName == "" {
		return fmt.Errorf("GITHUB_REPO_NAME environment variable not set")
	}

	return nil
}
