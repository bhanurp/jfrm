package main

import (
	"log"
	"os"

	"github.com/bhanurp/jfrm/internal/cli/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "jfrm",
		Usage: "Manage releases and dependencies for JFrog projects",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Usage:   "Run in dry-run mode (no changes will be made)",
			},
		},
		Commands: []*cli.Command{
			commands.UpdateDependencies(),
			commands.CheckDependencies(),
			commands.GenerateReport(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
