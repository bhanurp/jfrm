package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// githubReposBase allows tests to override the API base (default is GitHub REST)
var githubReposBase = "https://api.github.com/repos"

// GetLatestReleaseVersionAndCommitSHA fetches the latest Go module version and its commit SHA
func GetLatestReleaseVersionAndCommitSHA(module string) (string, string, time.Time, error) {
	latestVersion, releasedTime, err := fetchLatestVersion(module)
	if err != nil {
		return "", "", releasedTime, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	commitSHA, err := fetchCommitSHA(module, latestVersion)
	if err != nil {
		return "", "", releasedTime, fmt.Errorf("failed to fetch commit SHA for version %s: %w", latestVersion, err)
	}

	return latestVersion, commitSHA, releasedTime, nil
}

// fetchLatestVersion retrieves the latest tagged release version
func fetchLatestVersion(module string) (string, time.Time, error) {
	fmt.Printf("Fetching latest release for module: %s\n", module)
	url := fmt.Sprintf("%s/%s/releases/latest", githubReposBase, module)
	fmt.Println("Fetching latest release version using ", url)
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching latest version: %s\n", resp.Status)
		return "", time.Time{}, fmt.Errorf("unexpected response code 2: %d", resp.StatusCode)
	}

	var data struct {
		TagName     string    `json:"tag_name"`
		PublishedAt time.Time `json:"published_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", time.Time{}, err
	}

	fmt.Printf("Fetched version: %s, Time: %s\n", data.TagName, data.PublishedAt)

	return data.TagName, data.PublishedAt, nil
}

// fetchCommitSHA retrieves the commit SHA for a given version
func fetchCommitSHA(module, version string) (string, error) {
	url := fmt.Sprintf("%s/%s/git/refs/tags/%s", githubReposBase, module, version)
	fmt.Printf("Fetching commit SHA for %s\n", url)
	client := &http.Client{Timeout: 30 * time.Second} // Increased timeout to 30 seconds

	var resp *http.Response
	for i := 0; i < 5; i++ { // Retry up to 5 times
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err = client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to fetch commit SHA: %w", err)
		}
		if resp.StatusCode == http.StatusOK {
			break
		}

		if i < 4 {
			backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(backoff + jitter)
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("module or version not found: %s@%s", module, version)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching latest version: %s\n", resp.Body)
		return "", fmt.Errorf("unexpected response code 1: %d", resp.StatusCode)
	}

	var data struct {
		Ref    string `json:"ref"`
		URL    string `json:"url"`
		Object struct {
			SHA  string `json:"sha"`
			Type string `json:"type"`
		} `json:"object"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	fmt.Printf("Received data for commit sha: %v\n", data)

	if data.Object.SHA == "" {
		return "", fmt.Errorf("unexpected: sha is empty")
	}

	return data.Object.SHA, nil
}

// GetAllMergedPRs fetches all merged PRs since the last release
func GetAllMergedPRs(repo string, lastReleaseDate time.Time) ([]string, error) {
	base := "dev"
	if repo == "jfrog/jfrog-cli-artifactory" {
		base = "main"
	}
	url := fmt.Sprintf("%s/%s/pulls?state=closed&base=%s", githubReposBase, repo, base)
	fmt.Printf("Fetching all closed PRs for repo: %s URL used : %s\n", repo, url)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch PRs: %s", resp.Status)
	}

	var prs []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
		ClosedAt *time.Time `json:"closed_at"`
		MergedAt *time.Time `json:"merged_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, err
	}

	var prList []string
	for _, pr := range prs {
		if pr.ClosedAt.After(lastReleaseDate) {
			var tags []string
			for _, label := range pr.Labels {
				tags = append(tags, label.Name)
			}
			if pr.MergedAt != nil {
				prList = append(prList, fmt.Sprintf("PR #%d, %s, %s, %s, %s", pr.Number, pr.Title, pr.User.Login, strings.Join(tags, ", "), pr.MergedAt))
			}
		}
	}

	return prList, nil
}

// CreatePullRequest creates a pull request
func CreatePullRequest(branch, base, repo, token string) (string, error) {
	prBody := map[string]string{
		"title": "Update dependencies",
		"head":  branch,
		"base":  base,
		"body":  "This PR updates Go dependencies to the latest versions.",
	}
	jsonBody, _ := json.Marshal(prBody)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s/pulls", githubReposBase, repo), bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to create PR: %s", resp.Status)
	}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	prID, _ := result["number"].(float64)
	log.Printf("PR created successfully: #%d\n", int(prID))
	return fmt.Sprintf("%d", int(prID)), nil
}

// GetPullRequestStatus fetches the status of a pull request
func GetPullRequestStatus(prID, repo, token string) error {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s/pulls/%s", githubReposBase, repo, prID), nil)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to fetch PR status: %s", resp.Status)
	}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}
	log.Printf("PR Status: %s\n", result["state"].(string))
	return nil
}
