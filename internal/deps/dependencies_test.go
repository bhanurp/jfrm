package deps

import (
	"os"
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	if !IsNewerVersion("1.2.3", "1.2.4") {
		t.Fatalf("expected 1.2.4 to be newer than 1.2.3")
	}
	if IsNewerVersion("1.2.3", "1.2.3") {
		t.Fatalf("expected equal versions to not be newer")
	}
	if IsNewerVersion("1.3.0", "1.2.9") {
		t.Fatalf("expected 1.2.9 not newer than 1.3.0")
	}
}

func TestIsAllowedDependency(t *testing.T) {
	if !IsAllowedDependency("github.com/jfrog/gofrog") {
		t.Fatalf("expected gofrog to be allowed")
	}
	if IsAllowedDependency("example.com/not-allowed") {
		t.Fatalf("unexpected allowed dependency")
	}
}

func TestGetDependencies_Parse(t *testing.T) {
	data := []byte("module example.com/x\n\nrequire (\n\tgithub.com/jfrog/gofrog v1.2.3\n\tgithub.com/blang/semver/v4 v4.0.0\n)\n")
	if err := os.WriteFile("go.mod", data, 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove("go.mod") })
	m, err := GetDependencies()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["github.com/jfrog/gofrog"] != "v1.2.3" {
		t.Fatalf("expected parsed version v1.2.3, got %s", m["github.com/jfrog/gofrog"])
	}
}

