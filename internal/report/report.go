package report

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bhanurp/jfrm/internal/deps"
	"github.com/bhanurp/jfrm/internal/version"
)

// GenerateDryRunReport generates a dry-run report
func GenerateDryRunReport(repo string, prs []string, tag string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	report := fmt.Sprintf("# Dry-Run Report\n\n**Repository:** %s\n**Generated On:** %s\n\n", repo, timestamp)

	dryRunReport := deps.GetDryRunReport()
	if len(dryRunReport) == 0 {
		report += "✅ All dependencies are already up to date!\n"
	} else {
		report += "### Dependencies that would be updated:\n\n"
		for _, line := range dryRunReport {
			report += line + "\n"
		}
	}

	releaseType := version.DetermineReleaseType(prs)

	if len(prs) > 0 {
		report += "\n### Merged PRs since the latest release:\n\n"
		for _, pr := range prs {
			if strings.Contains(pr, ", ,") { // Check if PR has no labels
				report += pr + " (No labels)\n"
			} else {
				report += pr + "\n"
			}
		}
		report += fmt.Sprintf("\n### Decision on new release: %s\n", releaseType)
		report += fmt.Sprintf("Next possible version: %s\n", version.GetNextVersion(tag, releaseType))
	}

	err := os.WriteFile("dry-run-report.md", []byte(report), 0644)
	if err != nil {
		return fmt.Errorf("failed to write dry-run report: %w", err)
	}
	log.Println("✅ Dry-Run Report generated: dry-run-report.md")

	// Clear the dry run report after generating
	deps.ClearDryRunReport()
	return nil
}

// GenerateDependencyReport generates a comprehensive dependency report
func GenerateDependencyReport(repo string, dependencies map[string]string, prs []string, tag string, outputFile string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	report := fmt.Sprintf("# Dependency Report\n\n**Repository:** %s\n**Generated On:** %s\n**Current Version:** %s\n\n", repo, timestamp, tag)

	// Dependency Status Section
	report += "## Dependency Status\n\n"
	report += "| Module | Current Version | Latest Version | Status |\n"
	report += "|--------|----------------|----------------|--------|\n"

	updatesAvailable := 0
	for mod, currentVer := range dependencies {
		if deps.IsAllowedDependency(mod) {
			latestVer, err := deps.GetLatestModuleVersion(mod)
			if err != nil {
				report += fmt.Sprintf("| %s | %s | Error | ❌ Error |\n", mod, currentVer)
				continue
			}

			status := "✅ Up to date"
			if deps.IsNewerVersion(currentVer, latestVer) {
				status = "🔄 Update available"
				updatesAvailable++
			}
			report += fmt.Sprintf("| %s | %s | %s | %s |\n", mod, currentVer, latestVer, status)
		}
	}

	report += fmt.Sprintf("\n**Summary:** %d out of %d dependencies have updates available.\n\n", updatesAvailable, len(dependencies))

	// Recent Activity Section
	if len(prs) > 0 {
		report += "## Recent Activity\n\n"
		report += "### Merged PRs since the latest release:\n\n"
		for _, pr := range prs {
			if strings.Contains(pr, ", ,") {
				report += "- " + pr + " (No labels)\n"
			} else {
				report += "- " + pr + "\n"
			}
		}

		releaseType := version.DetermineReleaseType(prs)
		nextVersion := version.GetNextVersion(tag, releaseType)
		report += fmt.Sprintf("\n### Release Analysis\n")
		report += fmt.Sprintf("- **Recommended release type:** %s\n", releaseType)
		report += fmt.Sprintf("- **Next version:** %s\n", nextVersion)
	} else {
		report += "## Recent Activity\n\n"
		report += "No merged PRs found since the latest release.\n"
	}

	// Recommendations Section
	report += "\n## Recommendations\n\n"
	if updatesAvailable > 0 {
		report += "1. **Update Dependencies:** Consider updating the outdated dependencies to their latest versions.\n"
		report += "2. **Run Tests:** After updating dependencies, ensure all tests pass.\n"
		report += "3. **Review Changes:** Check for any breaking changes in the updated dependencies.\n"
	} else {
		report += "✅ All dependencies are up to date. No immediate action required.\n"
	}

	if len(prs) > 0 {
		report += "4. **Consider Release:** Based on the merged PRs, consider creating a new release.\n"
	}

	err := os.WriteFile(outputFile, []byte(report), 0644)
	if err != nil {
		return fmt.Errorf("failed to write dependency report: %w", err)
	}
	log.Printf("✅ Dependency Report generated: %s", outputFile)
	return nil
}
