# CLAUDE.md

<br/>

<br/>

## Project Structure

- Go-based GitHub Action (Docker container action)
- Generates changelogs from Conventional Commits
- Groups changes by type (Features, Bug Fixes, etc.), highlights breaking changes
- Features: PR/issue link extraction, author filtering (exclude bots), custom type mapping, dry-run mode

<br/>

## Build & Test

```bash
make build       # Build binary
make test        # Unit tests with coverage
make cover       # Generate coverage report (90% threshold in CI)
make fmt         # Format code
make lint        # Run go vet
```

<br/>

## Key Directories

- `cmd/` — Entry point (main.go)
- `internal/changelog/` — Core logic: parser.go (Conventional Commits), generator.go (changelog generation)
- `internal/config/` — Configuration loading from env vars (INPUT_*)
- `internal/git/` — Git operations wrapper
- `internal/output/` — GitHub Actions output helpers

<br/>

## Action Inputs

Key inputs: `output_file`, `tag_pattern`, `exclude_types`, `include_breaking`, `date_format`, `unreleased`, `since_tag`, `until_tag`, `custom_type_mapping`, `exclude_authors`, `dry_run`

Default excluded authors: `GitHub Action, GitHub Actions, dependabot[bot], renovate[bot], github-actions[bot]`

<br/>

## CI

- `ci.yml` — Unit tests (90% coverage threshold), Docker build & dry-run, action integration test
- Docker: multi-stage build (golang:1.26-alpine → alpine:3.23)

<br/>
