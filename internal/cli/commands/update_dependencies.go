package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bhanurp/jfrm/internal/deps"
	"github.com/bhanurp/jfrm/internal/github"
	"github.com/bhanurp/jfrm/internal/report"
	"github.com/bhanurp/jfrm/internal/version"
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
			&cli.StringFlag{
				Name:  "remote",
				Usage: "Base in form <remote>/<branch> (default: upstream/dev; upstream/main for jfrog/jfrog-cli-artifactory)",
			},
			&cli.StringFlag{
				Name:  "new-branch",
				Usage: "Override the generated branch name (e.g., update-dependencies-1.2.3)",
			},
		},
		Action: func(c *cli.Context) error {
			dryRun := c.Bool("dry-run")
			createPR := c.Bool("create-pr")

			// Determine default base remote/branch
			repo, err := deps.GetRepoName()
			if err != nil {
				return fmt.Errorf("failed to detect repository: %w", err)
			}
			baseRemote, baseBranch := resolveDefaultBase(repo)

			userBase := strings.TrimSpace(c.String("remote"))
			if userBase != "" {
				parts := strings.SplitN(userBase, "/", 2)
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return fmt.Errorf("invalid --remote value; expected <remote>/<branch>")
				}
				baseRemote, baseBranch = parts[0], parts[1]
			}

			// Preflight validation before any changes
			if err := runPreflightChecks(createPR); err != nil {
				return err
			}

			// Validate remote exists; if missing and user did not specify --remote, fallback to origin
			remotesOut, err := exec.Command("git", "remote").CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to list remotes: %w", err)
			}
			remotes := string(remotesOut)
			if !strings.Contains(remotes, baseRemote) {
				if userBase != "" {
					return fmt.Errorf("remote '%s' not found; configure it first", baseRemote)
				}
				// fallback to origin
				baseRemote = "origin"
				if !strings.Contains(remotes, baseRemote) {
					return fmt.Errorf("remote '%s' not found; configure it first", baseRemote)
				}
				if detected := detectDefaultRemoteBranch(baseRemote); detected != "" {
					if parts := strings.SplitN(detected, "/", 2); len(parts) == 2 {
						baseBranch = parts[1]
					}
				}
			}

			// Fetch the base branch refs (best-effort) and verify
			_ = exec.Command("git", "fetch", baseRemote, baseBranch).Run()
			if err := exec.Command("git", "rev-parse", "--verify", fmt.Sprintf("refs/remotes/%s/%s", baseRemote, baseBranch)).Run(); err != nil {
				return fmt.Errorf("base '%s/%s' not found after fetch", baseRemote, baseBranch)
			}

			if dryRun {
				log.Println("Running in Dry Run mode (No changes will be made)")
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

			// Get latest release information and merged PRs (for next version prediction)
			tag, lastReleaseSHA, releasedTime, err := github.GetLatestReleaseVersionAndCommitSHA(repo)
			if err != nil {
				return fmt.Errorf("failed to get latest release: %w", err)
			}
			log.Printf("Latest release: %s (Commit: %s) released on [%s]\n", tag, lastReleaseSHA, releasedTime.GoString())

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

			// If there are no dependency updates and no newly merged PRs, no need to release
			if len(updates) == 0 && len(prs) == 0 {
				log.Println("No dependency updates and no merged changes since last release — no new release needed.")
				return nil
			}

			// Generate report if in dry-run mode
			if dryRun {
				return report.GenerateDryRunReport(repo, prs, tag)
			}

			// Ensure go.sum is updated after any changes
			if len(updates) > 0 {
				if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
					log.Printf("warning: failed running 'go mod tidy': %v", err)
				}
			}

			// Create PR if requested
			if createPR {
				releaseType := version.DetermineReleaseType(prs)
				nextVersion := version.GetNextVersion(tag, releaseType)
				if strings.TrimSpace(nextVersion) == "" {
					nextVersion = "next"
				}
				branchName := buildBranchName(c.String("new-branch"), nextVersion)

				// Create local branch from the remote base
				if err := exec.Command("git", "checkout", "-B", branchName, fmt.Sprintf("%s/%s", baseRemote, baseBranch)).Run(); err != nil {
					return fmt.Errorf("failed to create branch from %s/%s: %w", baseRemote, baseBranch, err)
				}
				if err := deps.GitExec("add", "go.mod", "go.sum"); err != nil {
					return fmt.Errorf("failed to add files: %w", err)
				}
				if err := deps.GitExec("commit", "-m", fmt.Sprintf("chore(%s): update dependencies to latest versions", nextVersion)); err != nil {
					return fmt.Errorf("failed to commit: %w", err)
				}
				if err := deps.GitExec("push", "origin", branchName, "--force-with-lease"); err != nil {
					return fmt.Errorf("failed to push: %w", err)
				}

				token := os.Getenv("GITHUB_TOKEN")
				prID, err := github.CreatePullRequest(branchName, baseBranch, repo, token)
				if err != nil {
					return fmt.Errorf("failed to create PR: %w", err)
				}
				if err := github.GetPullRequestStatus(prID, repo, token); err != nil {
					log.Printf("Failed to get PR status: %v", err)
				}
			}

			return nil
		},
	}
}

// Helpers kept unexported for testing
func resolveDefaultBase(repo string) (remote, branch string) {
	remote, branch = "upstream", "dev"
	if repo == "jfrog/jfrog-cli-artifactory" {
		branch = "main"
	}
	return
}

func buildBranchName(override, next string) string {
	if strings.TrimSpace(override) != "" {
		return override
	}
	return fmt.Sprintf("update-dependencies-%s", next)
}

func detectDefaultRemoteBranch(remote string) string {
	// Try symbolic-ref to find remote HEAD like "origin/main"
	out, err := exec.Command("git", "symbolic-ref", "--quiet", "--short", fmt.Sprintf("refs/remotes/%s/HEAD", remote)).CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	// Fallback: git remote show <remote> and parse "HEAD branch: <name>"
	showOut, err := exec.Command("git", "remote", "show", remote).CombinedOutput()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(showOut), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "HEAD branch:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				if name != "" {
					return fmt.Sprintf("%s/%s", remote, name)
				}
			}
		}
	}
	return ""
}

// runPreflightChecks validates environment and repository state before any changes are attempted.
func runPreflightChecks(requirePR bool) error {
	var issues []string

	// go.mod must exist
	if _, err := os.Stat("go.mod"); err != nil {
		issues = append(issues, "missing go.mod in project root")
	}

	// Ensure git is available
	if _, err := exec.LookPath("git"); err != nil {
		issues = append(issues, "git not found in PATH")
	}

	// Ensure go is available
	if _, err := exec.LookPath("go"); err != nil {
		issues = append(issues, "go not found in PATH")
	}

	// Working tree must be clean
	if out, err := exec.Command("git", "status", "--porcelain").CombinedOutput(); err != nil {
		issues = append(issues, fmt.Sprintf("failed to check git status: %v", err))
	} else if strings.TrimSpace(string(out)) != "" {
		issues = append(issues, "working tree not clean; commit or stash changes first")
	}

	// Remote origin must exist
	if out, err := exec.Command("git", "remote", "get-url", "origin").CombinedOutput(); err != nil || strings.TrimSpace(string(out)) == "" {
		issues = append(issues, "git remote 'origin' not configured")
	}

	// If PR creation requested, validate token and remote push access (best-effort)
	if requirePR {
		if os.Getenv("GITHUB_TOKEN") == "" {
			issues = append(issues, "GITHUB_TOKEN is not set (required for PR creation)")
		}
		// Quick GitHub API reachability check
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get("https://api.github.com/rate_limit")
		if err != nil || resp.StatusCode >= 400 {
			issues = append(issues, "cannot reach GitHub API (network/auth issue)")
		} else {
			_ = resp.Body.Close()
		}
	}

	if len(issues) > 0 {
		return fmt.Errorf("preflight checks failed:\n- %s", strings.Join(issues, "\n- "))
	}
	return nil
}
