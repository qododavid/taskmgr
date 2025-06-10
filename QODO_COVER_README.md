# Qodo Test Coverage Bot Configuration

This directory contains the configuration for `qodo-cover`, an automated test coverage bot designed to run in GitHub Actions on pull requests.

## Quick Start

1. Copy `agent.toml` and `mcp.json` to your repository root
2. Add `QODO_API_KEY` to your repository secrets
3. Copy the workflow file to `.github/workflows/qodo-cover.yml`
4. Open a PR and add the `qodo-cover` label to trigger the bot!

## What it does

The bot automatically:
1. Analyzes PR changes to identify files that need test coverage
2. Generates appropriate tests for uncovered code  
3. Creates a follow-up PR with the new tests
4. Reports coverage results back to the original PR

## How it's triggered

The bot is designed to run automatically via GitHub Actions when:
- A pull request has the `qodo-cover` label applied
- The labeled pull request is updated with new commits

To trigger the bot:
1. Open a pull request
2. Add the `qodo-cover` label to the PR
3. The bot will automatically analyze and add test coverage

It can also be run locally using the Qodo CLI for testing purposes.

## Prerequisites

### For GitHub Actions
- `QODO_API_KEY` secret configured in your repository
- Repository permissions for creating PRs and branches
- Test framework installed in your project (pytest, go test, npm test, etc.)

### For Local Usage
- Qodo CLI installed (`qodo`, `gen`, or `qli`)
- GitHub CLI (`gh`) authenticated
- Python/Node.js/Go development environment (as needed by your project)

## Setup

### GitHub Actions Setup

1. Create the `qodo-cover` label in your repository:
   - Go to Issues → Labels → New label
   - Name: `qodo-cover`
   - Description: "Trigger test coverage bot"
   - Color: Choose any color you like

2. Add the `QODO_API_KEY` secret to your repository:
   - Go to Settings → Secrets and variables → Actions
   - Add a new secret named `QODO_API_KEY` with your API key

3. Place the `agent.toml` and `mcp.json` files in your repository root

4. Create the workflow file as shown in the usage section

### Local Setup

1. Ensure you have the Qodo CLI installed and authenticated:
   ```bash
   qodo login  # or set QODO_API_KEY environment variable
   ```

2. Place the `agent.toml` and `mcp.json` files in your project root or any parent directory

## Usage

### GitHub Actions (Recommended)

Create a workflow file `.github/workflows/qodo-cover.yml`:

```yaml
name: Test Coverage Bot
on:
  pull_request:
    types: [opened, synchronize, labeled]

jobs:
  coverage:
    runs-on: ubuntu-latest
    if: contains(github.event.pull_request.labels.*.name, 'qodo-cover')
    permissions:
      contents: write
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Qodo Coverage Bot
        uses: qododavid/qodo-gen-cli@main
        with:
          prompt: "qodo-cover"
          agentfile: "${{ github.workspace }}/agent.toml"
          key-value-pairs: |
            desired_coverage=90
        env:
          QODO_API_KEY: ${{ secrets.QODO_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          PR_NUMBER: ${{ github.event.pull_request.number }}
```

### Local usage (from within a GitHub repository with an open PR):
```bash
qodo qodo-cover
```

### With custom desired coverage:
```bash
# With custom desired coverage (default is 80%)
qodo qodo-cover --desired_coverage=90

# Or as a percentage
qodo qodo-cover --desired_coverage=85
```

### Run as a webhook server for CI/CD integration:
```bash
# Start webhook server on default port 4000
qodo --webhook

# Or on custom port
qodo --webhook --port=8080

# Trigger via webhook (from another terminal or CI system)
curl -X POST http://localhost:8080/webhook/qodo-cover \
  -H "Content-Type: application/json" \
  -d '{
    "desired_coverage": 85
  }'
```

## Configuration Files

### agent.toml
Contains the bot configuration, including:
- Command definition and description
- Detailed instructions for the AI agent
- Available tools (filesystem and shell)
- Single configurable argument (`desired_coverage`)
- MCP server configurations
- Structured output schema

### mcp.json
Defines the MCP (Model Context Protocol) servers:
- **filesystem**: For reading/writing test files
- **shell**: For running GitHub CLI, git, and test commands

## How it works

1. **Setup Phase**:
   - Detects repository and PR context from GitHub Actions environment variables
   - Posts initial greeting comment on PR
   - Uses the already checked-out repository (GitHub Actions handles this)

2. **Analysis Phase**:
   - Gets list of changed files from PR diff
   - Classifies files as logic-bearing (needs tests) or non-logic
   - Runs existing tests with coverage for changed files
   - Determines which files need additional test coverage

3. **Test Generation**:
   - For files needing tests, generates comprehensive test cases
   - Writes test files in appropriate test directories
   - Runs tests to ensure they pass and meet desired coverage level

4. **PR Creation**:
   - Creates new branch `add-coverage-{PR_NUMBER}`
   - Commits only the new/updated test files
   - Opens follow-up PR targeting the original PR branch

5. **Reporting**:
   - Posts summary comment on original PR
   - Includes coverage before/after metrics
   - Links to the follow-up PR with tests

## Arguments Reference

| Argument | Type | Default | Description |
|----------|------|---------|-------------|
| `desired_coverage` | number | 80 | Desired coverage percentage for changed lines |

## Supported Test Frameworks

The bot automatically detects and uses the appropriate test framework:
- **Python**: `pytest` with coverage
- **Go**: `go test` with coverage
- **JavaScript/TypeScript**: `npm test` or `yarn test`
- **Other**: Detects from project configuration

## Important Notes

- The bot **never modifies production code**, only test files
- All git push operations use GitHub CLI for proper authentication
- The bot restarts analysis if the PR is updated during execution
- Coverage is calculated specifically for changed lines, not overall project
- Test files follow project conventions (e.g., `tests/`, `test_*.py`, `*_test.go`)

## Example Workflow

1. Developer opens PR with new feature
2. Developer adds `qodo-cover` label to the PR
3. GitHub Actions triggers the coverage bot
4. Bot analyzes changes and finds uncovered code
5. Bot generates appropriate test cases
6. Bot creates follow-up PR with tests
7. Original PR gets comment with coverage report
8. Developer reviews and merges test PR
9. Original PR now has full test coverage

## Environment Variables

When running in GitHub Actions, these variables are automatically available:
- `GITHUB_REPOSITORY`: The owner/repo format (e.g., "octocat/hello-world")
- `GITHUB_TOKEN`: Authentication token with repository permissions
- `GITHUB_EVENT_PATH`: Path to the webhook event payload
- `GITHUB_EVENT_NAME`: The name of the webhook event (e.g., "pull_request")
- `PR_NUMBER`: Set by the workflow to pass the PR number to the agent

## Troubleshooting

### GitHub Actions Issues
- **Workflow not triggering**: Ensure the PR has the `qodo-cover` label applied
- **Permission errors**: Ensure the workflow has `contents: write` and `pull-requests: write` permissions
- **QODO_API_KEY not found**: Add the secret to your repository settings
- **Cannot create PR**: Check that GITHUB_TOKEN has appropriate permissions
- **Label removed mid-run**: The workflow will complete even if the label is removed during execution

### Local Usage Issues
- **Authentication issues**: Ensure `qodo login` is completed and `gh` is authenticated
- **Test failures**: Bot will retry and fix tests automatically
- **Coverage not detected**: Ensure test framework is properly configured in project 