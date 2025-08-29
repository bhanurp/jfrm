package commands

import "testing"

func TestResolveDefaultBase(t *testing.T) {
	r, b := resolveDefaultBase("some/other")
	if r != "upstream" || b != "dev" {
		t.Fatalf("expected upstream/dev, got %s/%s", r, b)
	}
	r, b = resolveDefaultBase("jfrog/jfrog-cli-artifactory")
	if r != "upstream" || b != "main" {
		t.Fatalf("expected upstream/main for artifactory, got %s/%s", r, b)
	}
}

func TestBuildBranchName(t *testing.T) {
	if got := buildBranchName("", "1.2.3"); got != "update-dependencies-1.2.3" {
		t.Fatalf("unexpected branch name: %s", got)
	}
	if got := buildBranchName("custom-name", "1.2.3"); got != "custom-name" {
		t.Fatalf("override not respected: %s", got)
	}
}
