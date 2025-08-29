package version

import "testing"

func TestIncrementPatchVersion(t *testing.T) {
	if got := IncrementPatchVersion("1.2.3"); got != "1.2.4" {
		t.Fatalf("expected 1.2.4, got %s", got)
	}
}

func TestIncrementMinorVersion(t *testing.T) {
	if got := IncrementMinorVersion("1.2.3"); got != "1.3.0" {
		t.Fatalf("expected 1.3.0, got %s", got)
	}
}

func TestGetNextVersion(t *testing.T) {
	if got := GetNextVersion("1.2.3", "next patch"); got != "1.2.4" {
		t.Fatalf("expected 1.2.4, got %s", got)
	}
	if got := GetNextVersion("1.2.3", "next minor"); got != "1.3.0" {
		t.Fatalf("expected 1.3.0, got %s", got)
	}
}

