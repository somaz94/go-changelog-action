# Go Changelog Action

[![Continuous Integration](https://github.com/somaz94/go-changelog-action/actions/workflows/ci.yml/badge.svg)](https://github.com/somaz94/go-changelog-action/actions/workflows/ci.yml)
[![Lint Codebase](https://github.com/somaz94/go-changelog-action/actions/workflows/linter.yml/badge.svg)](https://github.com/somaz94/go-changelog-action/actions/workflows/linter.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go-based GitHub Action that generates changelogs from
[Conventional Commits](https://www.conventionalcommits.org/).

## Features

- Parses git history using Conventional Commits format
- Groups changes by type (Features, Bug Fixes, etc.)
- Highlights BREAKING CHANGES
- Auto-detects PR links (`(#123)`) and issue references (`closes #456`)
- Generates compare links between versions (`v1.0.0...v1.1.0`)
- Contributors section per release
- Includes non-conventional commits in "Other Changes" (optional)
- Tag range filtering with `since_tag` / `until_tag`
- Custom type-to-section mapping via JSON
- Supports unreleased changes section
- Configurable tag patterns, date formats, and excluded types
- Dry-run mode for previewing output
- Auto-detects repository URL for commit links

## Usage

### Basic

```yaml
steps:
  - name: Checkout
    uses: actions/checkout@v6
    with:
      fetch-depth: 0

  - name: Generate Changelog
    uses: somaz94/go-changelog-action@v1
    with:
      output_file: CHANGELOG.md
```

### With Custom Options

```yaml
steps:
  - name: Checkout
    uses: actions/checkout@v6
    with:
      fetch-depth: 0

  - name: Generate Changelog
    id: changelog
    uses: somaz94/go-changelog-action@v1
    with:
      output_file: CHANGELOG.md
      tag_pattern: 'v[0-9]*.[0-9]*.[0-9]*'
      exclude_types: 'chore,style'
      include_breaking: true
      include_non_conventional: true
      date_format: '2006-01-02'
      unreleased: true
      skip_commits: '^Merge'

  - name: Commit Changelog
    uses: somaz94/go-git-commit-action@v1
    with:
      branch: main
      commit_message: 'docs: update changelog'
      user_email: 'actions@github.com'
      user_name: 'GitHub Actions'
      file_pattern: 'CHANGELOG.md'
      github_token: ${{ secrets.GITHUB_TOKEN }}
```

### Tag Range (Partial Changelog)

```yaml
- name: Generate Changelog for specific range
  uses: somaz94/go-changelog-action@v1
  with:
    since_tag: v1.0.0
    until_tag: v2.0.0
```

### Custom Type Mapping

```yaml
- name: Generate Changelog with custom sections
  uses: somaz94/go-changelog-action@v1
  with:
    custom_type_mapping: '{"feat": "New Features", "fix": "Bugfixes", "perf": "Optimizations"}'
```

### Dry Run (Preview Only)

```yaml
- name: Preview Changelog
  uses: somaz94/go-changelog-action@v1
  with:
    dry_run: true
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `output_file` | Output file path for the changelog | No | `CHANGELOG.md` |
| `tag_pattern` | Glob pattern to match version tags | No | `v[0-9]*.[0-9]*.[0-9]*` |
| `exclude_types` | Comma-separated commit types to exclude | No | `` |
| `include_breaking` | Include BREAKING CHANGE section | No | `true` |
| `date_format` | Date format (Go time format) | No | `2006-01-02` |
| `header` | Custom header for changelog | No | `# Changelog` |
| `unreleased` | Include unreleased changes section | No | `true` |
| `unreleased_title` | Title for unreleased section | No | `Unreleased` |
| `skip_commits` | Regex pattern to skip commits | No | `^Merge` |
| `repository_url` | Repository URL for links (auto-detected) | No | `` |
| `dry_run` | Preview without writing to file | No | `false` |
| `include_non_conventional` | Include non-conventional commits in "Other Changes" | No | `false` |
| `since_tag` | Generate changelog starting from this tag (inclusive) | No | `` |
| `until_tag` | Generate changelog up to this tag (inclusive) | No | `` |
| `custom_type_mapping` | JSON mapping of commit types to section names | No | `` |

## Outputs

| Output | Description |
|--------|-------------|
| `changelog_file` | Path to the generated changelog file |
| `changelog_content` | Generated changelog content |
| `entries_count` | Number of changelog entries generated |
| `latest_version` | Latest version tag found |

## Generated Changelog Format

```markdown
# Changelog

## [v1.1.0](https://github.com/owner/repo/compare/v1.0.0...v1.1.0) (2024-03-15)

### Features

- **auth:** add OAuth2 support ([#42](https://github.com/owner/repo/pull/42)) ([abc1234](https://github.com/owner/repo/commit/abc1234))

### Bug Fixes

- fix memory leak ([def5678](https://github.com/owner/repo/commit/def5678)), closes [#99](https://github.com/owner/repo/issues/99)

### Contributors

- alice
- bob
```

## Conventional Commits

This action parses commits following the
[Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Supported Types

| Type | Section |
|------|---------|
| `feat` | Features |
| `fix` | Bug Fixes |
| `docs` | Documentation |
| `style` | Styles |
| `refactor` | Code Refactoring |
| `perf` | Performance Improvements |
| `test` | Tests |
| `build` | Builds |
| `ci` | Continuous Integration |
| `chore` | Chores |
| `revert` | Reverts |

### Breaking Changes

Breaking changes are detected via:

- `!` after the type/scope: `feat!: remove deprecated API`
- `BREAKING CHANGE` in the commit body

### PR and Issue Links

- PR numbers are auto-detected from `(#123)` in commit messages
- Issue references are detected from `closes #123`, `fixes #456`,
  `resolves #789` in commit messages and bodies

## Project Structure

```
.
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── changelog/
│   │   ├── generator.go     # Changelog generation logic
│   │   ├── generator_test.go
│   │   ├── parser.go        # Conventional commit parser
│   │   └── parser_test.go
│   ├── config/
│   │   ├── config.go        # Configuration from env vars
│   │   └── config_test.go
│   ├── git/
│   │   ├── git.go           # Git operations
│   │   └── git_test.go
│   └── output/
│       ├── output.go        # GitHub Actions output helpers
│       └── output_test.go
├── action.yml
├── Dockerfile
├── Makefile
└── go.mod
```

## Development

### Prerequisites

- Go 1.26+
- Docker (for container builds)

### Build

```bash
make build
```

### Test

```bash
make test
```

### Coverage

```bash
make cover
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.
