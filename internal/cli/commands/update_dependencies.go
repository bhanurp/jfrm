package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/bhanurp/jfrm/internal/deps"
	"github.com/bhanurp/jfrm/internal/github"
	"github.com/bhanurp/jfrm/internal/report"
	"github.com/urfave/cli/v2"
)

// UpdateDependencies creates the update-dependencies command
func UpdateDependencies() *cli.Command {
	return &cli.Command{
		Name:    "update-dependencies",
		Aliases: []string{"ud"},
		Usage:   "Update Go dependencies to latest versions",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Usage:   "Run in dry-run mode (no changes will be made)",
			},
			&cli.BoolFlag{
				Name:    "create-pr",
				Aliases: []string{"p"},
				Usage:   "Create a pull request with the changes",
			},
		},
		Action: func(c *cli.Context) error {
			dryRun := c.Bool("dry-run")
			createPR := c.Bool("create-pr")

			if dryRun {
				log.Println("Running in Dry Run mode (No changes will be made)")
			}

			// Get repository information
			repo, err := deps.GetRepoName()
			if err != nil {
				return fmt.Errorf("failed to detect repository: %w", err)
			}
			log.Printf("Detected repository: %s\n", repo)

			// Get current dependencies
			dependencies, err := deps.GetDependencies()
			if err != nil {
				return fmt.Errorf("failed to read go.mod: %w", err)
			}

			// Update dependencies
			updates := make(map[string]string)
			for mod, currentVer := range dependencies {
				if !deps.IsAllowedDependency(mod) {
					continue
				}
				fmt.Printf("Fetching latest version for: %s\n", mod)
				latestVer, err := deps.GetLatestModuleVersion(mod)
				if err != nil {
					log.Printf("Skipping %s: %v", mod, err)
					continue
				}
				if deps.IsNewerVersion(currentVer, latestVer) {
					log.Printf("Updating %s from %s -> %s", mod, currentVer, latestVer)
					updates[mod] = latestVer
				}
				if err := deps.UpdateDependency(mod, currentVer, latestVer, dryRun); err != nil {
					log.Printf("Failed to update %s: %v", mod, err)
				}
			}

			if len(updates) == 0 {
				log.Println("No dependencies to update!")
				return nil
			}

			// Get latest release information
			tag, lastReleaseSHA, releasedTime, err := github.GetLatestReleaseVersionAndCommitSHA(repo)
			if err != nil {
				return fmt.Errorf("failed to get latest release: %w", err)
			}
			log.Printf("Latest release: %s (Commit: %s) released on [%s]\n", tag, lastReleaseSHA, releasedTime.GoString())

			// Get merged PRs since last release
			prs, err := github.GetAllMergedPRs(repo, releasedTime)
			if err != nil {
				log.Printf("Error fetching merged PRs: %v\n", err)
			}

			if len(prs) == 0 {
				fmt.Println("No merged PRs found since the latest release.")
			} else {
				fmt.Println("Merged PRs since the latest release:")
				for _, pr := range prs {
					fmt.Println(pr)
				}
			}

			// Generate report if in dry-run mode
			if dryRun {
				return report.GenerateDryRunReport(repo, prs, tag)
			}

			// Create PR if requested
			if createPR {
				return createPullRequestWithUpdates(repo, updates)
			}

			return nil
		},
	}
}

func createPullRequestWithUpdates(repo string, updates map[string]string) error {
	branchName := "update-dependencies"
	if err := deps.GitExec("checkout", "-b", branchName); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	if err := deps.GitExec("add", "go.mod", "go.sum"); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}
	if err := deps.GitExec("commit", "-m", "chore: update dependencies to latest versions"); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	if err := deps.GitExec("push", "origin", branchName); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	prID, err := github.CreatePullRequest(branchName, repo, token)
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	err = github.GetPullRequestStatus(prID, repo, token)
	if err != nil {
		log.Printf("Failed to get PR status: %v", err)
	}

	return nil
}
