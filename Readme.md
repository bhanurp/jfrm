# JFrogCLI Release Manager (jfrm)

A Go CLI tool for managing releases and dependencies for JFrog projects. This tool automates the process of updating Go dependencies, analyzing repository changes, and generating comprehensive reports.

## Features

- **Dependency Management**: Check and update Go dependencies to their latest versions
- **GitHub Integration**: Fetch latest releases, commit SHAs, and merged PRs
- **Dry Run Mode**: Preview changes without making actual modifications
- **Report Generation**: Generate detailed dependency and release reports
- **Pull Request Creation**: Automatically create PRs for dependency updates
- **Version Analysis**: Determine appropriate release types based on changes

## Upcoming features

- **GitHub Action - Auto PR**: Automatically create a PR with all required dependency updates
- **GitHub Action - Release chain**: Optionally create a PR for the next repository in the release chain
- **Automated Releases**: Generate release notes and request approval before publishing a new release

## Installation

### Homebrew (Recommended)
```bash
brew install bhanurp/jfrm-tap/jfrm
```

### Manual Installation
```bash
# Clone and build
git clone https://github.com/bhanurp/jfrm.git
cd jfrm
go build -o jfrm ./cmd/jfrm
sudo mv jfrm /usr/local/bin/
```

### Go Install (Public Repository Only)
```bash
go install github.com/bhanurp/jfrm/cmd/jfrm@latest
```

## Usage

### Check Dependencies

Check the current status of all dependencies:

```bash
jfrm check-dependencies
```

### Update Dependencies

Update dependencies to their latest versions:

```bash
# Dry run mode (preview changes)
jfrm update-dependencies --dry-run

# Update dependencies and create a pull request
jfrm update-dependencies --create-pr

# Update dependencies without creating PR
jfrm update-dependencies
```

### Generate Reports

Generate comprehensive dependency reports:

```bash
# Generate default report
jfrm generate-report

# Generate report with custom output file
jfrm generate-report --output custom-report.md
```

## Configuration

### Environment Variables

- `GITHUB_TOKEN`: GitHub API token for authenticated requests (optional but recommended)

### Allowed Dependencies

The tool only manages dependencies from the following JFrog modules:
- `github.com/jfrog/jfrog-cli-core/v2`
- `github.com/jfrog/jfrog-client-go`
- `github.com/jfrog/jfrog-cli-artifactory`
- `github.com/jfrog/jfrog-cli-security`
- `github.com/jfrog/build-info-go`
- `github.com/jfrog/gofrog`

## Supported repositories

- github.com/jfrog/jfrog-cli
- github.com/jfrog/jfrog-cli-core/v2
- github.com/jfrog/jfrog-client-go
- github.com/jfrog/jfrog-cli-artifactory
- github.com/jfrog/jfrog-cli-security
- github.com/jfrog/build-info-go
- github.com/jfrog/gofrog

Notes:
- Default base branch is `upstream/dev`.
- For `jfrog/jfrog-cli-artifactory`, default base is `upstream/main`.
- You can override via `--remote <remote>/<branch>` (e.g., `--remote upstream/master`).

## Project Structure

```
jfrm/
├── cmd/
│   └── jfrm/
│       └── main.go              # CLI entry point
├── internal/
│   ├── cli/
│   │   └── commands/            # CLI commands
│   │       ├── update_dependencies.go
│   │       ├── check_dependencies.go
│   │       └── generate_report.go
│   ├── deps/
│   │   └── dependencies.go      # Dependency management
│   ├── github/
│   │   └── github.go           # GitHub API integration
│   ├── version/
│   │   └── version.go          # Version management
│   └── report/
│       └── report.go           # Report generation
├── go.mod
├── go.sum
└── README.md
```

## Development

### Prerequisites

- Git

### Building

```bash
make build
# or
go build -o jfrm cmd/jfrm/main.go
```

### Running Tests

```bash
make test
# or
go test ./...
```

### Creating Releases of jfrm

The release process is fully automated:

```bash
# Create a new release (this will trigger GitHub Actions)
make release
# or
./scripts/create-release.sh v0.1.0
```

This will automatically:
1. Create a Git tag
2. Trigger GitHub Actions to build binaries
3. Create a GitHub release with assets
4. Update the Homebrew tap

### Automated Workflows

- **Test**: Runs on every push and PR to main branch
- **Release**: Runs when you push a tag (v*)
- **Homebrew Update**: Automatically updates the tap after release

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
