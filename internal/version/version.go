package version

import (
	"strings"

	"github.com/blang/semver/v4"
)

// IncrementMinorVersion increments the minor version and resets patch to 0
func IncrementMinorVersion(tag string) string {
	v, err := semver.ParseTolerant(tag)
	if err != nil {
		return ""
	}
	v.Minor++
	v.Patch = 0
	return v.String()
}

// IncrementPatchVersion increments the patch version
func IncrementPatchVersion(version string) string {
	v, err := semver.ParseTolerant(version)
	if err != nil {
		return ""
	}
	v.Patch++
	return v.String()
}

// GetNextVersion determines the next version based on release type
func GetNextVersion(tag string, releaseType string) string {
	if releaseType == "next patch" {
		return IncrementPatchVersion(tag)
	}
	return IncrementMinorVersion(tag)
}

// DetermineReleaseType determines the release type based on PR content
func DetermineReleaseType(prs []string) string {
	releaseType := "next patch"
	for _, pr := range prs {
		if strings.Contains(pr, "new feature") || strings.Contains(pr, "feature request") {
			releaseType = "next minor"
			break
		}
	}
	return releaseType
}
