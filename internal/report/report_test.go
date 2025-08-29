package report

import (
	"os"
	"testing"
)

func TestGenerateDryRunReport(t *testing.T) {
	repo := "owner/repo"
	prs := []string{"PR #1, test, user, bug, 2025-07-17"}
	tag := "v1.2.3"
	if err := GenerateDryRunReport(repo, prs, tag); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// File should exist
	if _, err := os.Stat("dry-run-report.md"); err != nil {
		t.Fatalf("expected report file, got error: %v", err)
	}
	_ = os.Remove("dry-run-report.md")
}

