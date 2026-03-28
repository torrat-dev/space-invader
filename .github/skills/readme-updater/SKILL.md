---
name: readme-updater
description: '**WORKFLOW SKILL** — Automatically create or update README.md files for Go and HTMX projects with current repository context, ensuring documentation stays synchronized with code changes that programmers often forget to document.

USE FOR: Initial README generation for new Go/HTMX projects, updating existing READMEs after significant code changes, maintaining up-to-date project documentation.

DO NOT USE FOR: General documentation writing, non-README files, projects in other languages/frameworks.

INVOKES: file system tools (read/write README.md), semantic_search for codebase analysis, subagents for content generation and review.'
---

## Workflow Steps

1. Analyze repository structure and key files (main.go, go.mod, templates/, etc.)
2. Check for existing README.md; if present, read current content
3. Extract project details: name, description, dependencies (Go modules, HTMX), build/run instructions
4. Generate content for standard README sections:
   - Project title and description
   - Installation instructions
   - Usage examples
   - Build and run commands
   - Project structure overview
   - Contributing guidelines
   - License information
5. If README exists, merge new information with existing content, preserving manual additions
6. Write or update README.md file
7. Validate the generated content for accuracy and completeness

## Quality Criteria

- README should include all standard sections for a Go/HTMX project
- Commands should be tested and functional
- Content should be clear and concise
- Merging should not overwrite user-added content
- No outdated information from previous versions