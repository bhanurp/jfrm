package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetLatestReleaseVersionAndCommitSHA(t *testing.T) {
	mux := http.NewServeMux()
	// latest release endpoint
	mux.HandleFunc("/owner/repo/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"tag_name":     "v1.2.3",
			"published_at": time.Now().Format(time.RFC3339),
		})
	})
	// tag ref endpoint
	mux.HandleFunc("/owner/repo/git/refs/tags/v1.2.3", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"ref":    "refs/tags/v1.2.3",
			"object": map[string]string{"sha": "abc123", "type": "tag"},
		})
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	oldBase := githubReposBase
	githubReposBase = ts.URL
	defer func() { githubReposBase = oldBase }()

	tag, sha, tms, err := GetLatestReleaseVersionAndCommitSHA("owner/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag != "v1.2.3" || sha != "abc123" || tms.IsZero() {
		t.Fatalf("unexpected result tag=%s sha=%s time=%v", tag, sha, tms)
	}
}

