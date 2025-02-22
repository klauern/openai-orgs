package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/klauern/openai-orgs/cmd"
	"github.com/urfave/cli/v2"
)

var Version = "dev"

func main() {
	app := &cli.App{
		Name:    "openai-orgs",
		Usage:   "CLI for OpenAI Platform Management API",
		Version: Version,
		Commands: []*cli.Command{
			cmd.OrganizationsCommand(),
			cmd.AdminAPIKeysCommand(),
			cmd.AuditLogsCommand(),
			cmd.InvitesCommand(),
			cmd.UsersCommand(),
			cmd.ProjectsCommand(),
			cmd.ProjectUsersCommand(),
			cmd.ProjectServiceAccountsCommand(),
			cmd.ProjectApiKeysCommand(),
			cmd.ProjectRateLimitsCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output format (default: pretty)",
				Value: "pretty",
				Action: func(ctx *cli.Context, s string) error {
					if s == "" {
						return nil
					}

					if _, ok := cmd.ValidOutputFormats[s]; ok {
						return nil
					}

					return fmt.Errorf("invalid output format: %s", s)
				},
			},
			&cli.StringFlag{
				Name:    "api-key",
				Usage:   "OpenAI API key (can be set via OPENAI_API_KEY environment variable)",
				EnvVars: []string{"OPENAI_API_KEY"},
				Action: func(ctx *cli.Context, s string) error {
					if !strings.HasPrefix(s, "sk-admin-") {
						return fmt.Errorf("invalid API key, must start with sk-admin-")
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
