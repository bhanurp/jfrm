# JFrog Release Manager (jfrm)

A Go CLI tool for managing releases and dependencies for JFrog projects. This tool automates the process of updating Go dependencies, analyzing repository changes, and generating comprehensive reports.

## Features

- **Dependency Management**: Check and update Go dependencies to their latest versions
- **GitHub Integration**: Fetch latest releases, commit SHAs, and merged PRs
- **Dry Run Mode**: Preview changes without making actual modifications
- **Report Generation**: Generate detailed dependency and release reports
- **Pull Request Creation**: Automatically create PRs for dependency updates
- **Version Analysis**: Determine appropriate release types based on changes

## Installation

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

- Go 1.23.6 or later
- Git

### Building

```bash
go build -o jfrm cmd/jfrm/main.go
```

### Running Tests

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.