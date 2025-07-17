package commands

import (
	"fmt"
	"log"

	"github.com/bhanurp/jfrm/internal/deps"
	"github.com/bhanurp/jfrm/internal/github"
	"github.com/bhanurp/jfrm/internal/report"
	"github.com/urfave/cli/v2"
)

// GenerateReport creates the generate-report command
func GenerateReport() *cli.Command {
	return &cli.Command{
		Name:    "generate-report",
		Aliases: []string{"gr"},
		Usage:   "Generate a dependency update report",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path for the report",
				Value:   "dependency-report.md",
			},
		},
		Action: func(c *cli.Context) error {
			outputFile := c.String("output")

			// Get repository information
			repo, err := deps.GetRepoName()
			if err != nil {
				return fmt.Errorf("failed to detect repository: %w", err)
			}

			// Get current dependencies
			dependencies, err := deps.GetDependencies()
			if err != nil {
				return fmt.Errorf("failed to read go.mod: %w", err)
			}

			// Get latest release information
			tag, _, releasedTime, err := github.GetLatestReleaseVersionAndCommitSHA(repo)
			if err != nil {
				return fmt.Errorf("failed to get latest release: %w", err)
			}

			// Get merged PRs since last release
			prs, err := github.GetAllMergedPRs(repo, releasedTime)
			if err != nil {
				log.Printf("Error fetching merged PRs: %v\n", err)
			}

			// Generate the report
			return report.GenerateDependencyReport(repo, dependencies, prs, tag, outputFile)
		},
	}
}
