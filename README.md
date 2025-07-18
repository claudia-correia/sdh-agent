# AI Agent for SDH Issue Resolution

This Go application implements an AI agent that analyzes a GitHub issue, finds similar resolved issues, and posts a comprehensive analysis report as a comment.

## Setup

1.  **Prerequisites:**
    * Go 1.21 or later.
    * A GitHub account with a Personal Access Token (PAT) that has read access to repository issues.
    * An API key for the LLM of choice.
        * Currently, only Anthropic Claude is supported.

2.  **Configuration:**
    * Clone the repository.
    * Create a `.env` file in the root directory by copying the `.env.example` file:
        ```bash
        cp .env.example .env
        ```
    * Edit the `.env` file and fill in your actual credentials and repository details.

3.  **Install Dependencies:**
    Open your terminal in the project root and run:
    ```bash
    go mod tidy
    ```

## What to Expect

The agent will:

1. Load configuration from your `.env` file
2. Connect to GitHub and retrieve the specified issue
3. Process the issue (summarize, find similar issues, analyze)
4. Generate a report
5. Print the report to the console

## Usage

To run the agent, execute the `cmd/sdh-agent/main.go` program from your terminal, passing the number of the target GitHub issue as a command-line argument.

```bash
go run main.go <issue-number>
```

Replace <issue-number> with the GitHub issue number you want to analyze. For example, use `123` for issue https://github.com/your-github-username/your-repo-name/issues/123.

### Building an Executable

If you want to build an executable:

```bash
go build -o sdh-agent cmd/sdh-agent/main.go
```

Then run:

```bash
./sdh-agent <issue-number>
```
