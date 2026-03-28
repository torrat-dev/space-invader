---
name: go-linting
description: '**WORKFLOW SKILL** — Run Go linting tools (go vet, staticcheck) on the project, report issues, trigger game-refactor skill if issues are found, and optionally create a GitHub issue. USE FOR: checking Go code quality and potential bugs in the space-invader project. DO NOT USE FOR: non-Go projects, general refactoring without game context. INVOKES: run_in_terminal for linting commands, subagents for refactoring, ask-questions for user options.'
---

# Go Linting Skill

## Workflow

This skill performs a full multi-step Go linting process on the entire project:

1. **Check Dependencies**: Ensure staticcheck is installed; install via `go install` if missing.
2. **Run go vet**: Checks for suspicious constructs in the Go code.
3. **Run staticcheck**: Performs static analysis to find bugs and performance issues.
4. **Run gofmt**: Check for formatting issues.
5. **Parse and Summarize Issues**: Collect output from tools and summarize for the AI agent.
6. **Report Issues**: Display the summary of any linting issues found.
7. **Trigger Refactor**: If issues are detected, invoke the game-refactor skill with the error summary as context to fix them.
8. **Optional GitHub Issue**: If issues are found, ask the user if they want to create a GitHub issue with the linting summary.
9. **Error Handling**: If tools fail due to missing dependencies, suggest installation steps.

## Usage

Invoke this skill when you need to ensure the Go codebase is clean and follows best practices, particularly in the context of the space-invader game project. It can also help track issues via GitHub for team collaboration.