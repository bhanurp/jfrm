package commands

import (
	"fmt"
	"log"

	"github.com/bhanurp/jfrm/internal/deps"
	"github.com/urfave/cli/v2"
)

// CheckDependencies creates the check-dependencies command
func CheckDependencies() *cli.Command {
	return &cli.Command{
		Name:    "check-dependencies",
		Aliases: []string{"cd"},
		Usage:   "Check current dependency status",
		Action: func(c *cli.Context) error {
			dependencies, err := deps.GetDependencies()
			if err != nil {
				return fmt.Errorf("failed to read go.mod: %w", err)
			}

			fmt.Println("Current Dependencies:")
			fmt.Println("=====================")
			for mod, currentVer := range dependencies {
				if deps.IsAllowedDependency(mod) {
					latestVer, err := deps.GetLatestModuleVersion(mod)
					if err != nil {
						log.Printf("Failed to get latest version for %s: %v", mod, err)
						continue
					}
					status := "✅ Up to date"
					if deps.IsNewerVersion(currentVer, latestVer) {
						status = fmt.Sprintf("🔄 Update available: %s → %s", currentVer, latestVer)
					}
					fmt.Printf("%s: %s (%s)\n", mod, currentVer, status)
				}
			}
			return nil
		},
	}
}
