package main

import (
	"fmt"
	"log"
	"os"

	"sdh-agent/internal/agent"
	"sdh-agent/internal/config"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Get issue number from command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: sdh-agent <issue-number>")
	}

	var issueNumber int
	_, err = fmt.Sscanf(os.Args[1], "%d", &issueNumber)
	if err != nil {
		log.Fatal("Invalid issue number")
	}

	// Initialize and run the agent
	sdhAgent := agent.NewSDHAgent(*cfg)

	log.Printf("▶️  Starting analysis for issue: #%d\n", issueNumber)
	report, err := sdhAgent.ProcessIssue(issueNumber)
	if err != nil {
		log.Fatalf("❌ An error occurred during processing: %v", err)
	}

	log.Println("✅ Successfully processed issue and generated report:")
	log.Println("===== REPORT BEGIN =====")

	// Print the final report
	fmt.Println(report)

	log.Println("===== REPORT END =====")
}
