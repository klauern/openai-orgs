package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/klauern/openai-orgs/cmd"
	"github.com/urfave/cli/v3"
)

var Version = "dev"

func main() {
	app := &cli.Command{
		Name:    "openai-orgs",
		Usage:   "CLI for OpenAI Platform Management API",
		Version: Version,
		Commands: []*cli.Command{
			cmd.AdminAPIKeysCommand(),
			cmd.AuditLogsCommand(),
			cmd.InvitesCommand(),
			cmd.UsersCommand(),
			cmd.ProjectsCommand(),
			cmd.ProjectUsersCommand(),
			cmd.ProjectServiceAccountsCommand(),
			cmd.ProjectApiKeysCommand(),
			cmd.ProjectRateLimitsCommand(),
			cmd.UsageCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output format (default: pretty)",
				Value: "pretty",
				Action: func(ctx context.Context, cmd *cli.Command, value string) error {
					if value == "" {
						return nil
					}

					validFormats := []string{"pretty", "json", "jsonl"}
					for _, format := range validFormats {
						if format == value {
							return nil
						}
					}

					return fmt.Errorf("invalid output format: %s (valid formats: %v)", value, validFormats)
				},
			},
			&cli.StringFlag{
				Name:     "api-key",
				Usage:    "OpenAI API key (can be set via OPENAI_API_KEY environment variable)",
				Sources:  cli.EnvVars("OPENAI_API_KEY"),
				Required: true,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
