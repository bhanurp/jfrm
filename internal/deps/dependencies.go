package deps

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/jfrog/gofrog/version"
	"golang.org/x/mod/modfile"
)

var allowedDeps = map[string]bool{
	"github.com/jfrog/jfrog-cli-core/v2":     true,
	"github.com/jfrog/jfrog-client-go":       true,
	"github.com/jfrog/jfrog-cli-artifactory": true,
	"github.com/jfrog/jfrog-cli-security":    true,
	"github.com/jfrog/build-info-go":         true,
	"github.com/jfrog/gofrog":                true,
}

var dryRunReport []string

// GetRepoName extracts the repository name from git remote
func GetRepoName() (string, error) {
	output, err := execCmd("git", "remote", "get-url", "upstream")
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`[:/]([^/]+/[^/]+)(?:\.git)?$`)
	fmt.Println("Finding git repo name")
	fmt.Printf("%s\n", output)
	fmt.Println("---------------------")
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("could not determine repo name")
}

// execCmd executes a command and returns the output
func execCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// GetDependencies reads and parses go.mod file
func GetDependencies() (map[string]string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, err
	}
	deps := make(map[string]string)
	for _, req := range modFile.Require {
		deps[req.Mod.Path] = req.Mod.Version
	}
	return deps, nil
}

// GetLatestModuleVersion fetches the latest version for a module
func GetLatestModuleVersion(module string) (string, error) {
	fmt.Printf("Fetching latest version for module: %s\n", module)
	url := fmt.Sprintf("https://proxy.golang.org/%s/@latest", module)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching latest version: %s\n", resp.Body)
		return "", fmt.Errorf("unexpected response code 3: %d", resp.StatusCode)
	}
	var data struct {
		Version string `json:"Version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Version, nil
}

// IsAllowedDependency checks if a dependency is in the allowed list
func IsAllowedDependency(module string) bool {
	return allowedDeps[module]
}

// IsNewerVersion checks if the latest version is newer than current
func IsNewerVersion(current, latest string) bool {
	currentVersion := version.NewVersion(current)
	return currentVersion.Compare(latest) > 0
}

// UpdateDependency updates a dependency to the latest version
func UpdateDependency(module, currentVersion, latestVer string, dryRun bool) error {
	if dryRun {
		log.Printf("[Dry Run] Would update: %s -> %s\n", module, latestVer)
		dryRunReport = append(dryRunReport, fmt.Sprintf("- `%s`: **%s → %s**", module, currentVersion, latestVer))
		return nil
	}
	_, err := execCmd("go", "get", fmt.Sprintf("%s@%s", module, latestVer))
	return err
}

// GitExec executes git commands
func GitExec(params ...string) error {
	cmd := exec.Command("git", params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetDryRunReport returns the dry run report
func GetDryRunReport() []string {
	return dryRunReport
}

// ClearDryRunReport clears the dry run report
func ClearDryRunReport() {
	dryRunReport = nil
}
