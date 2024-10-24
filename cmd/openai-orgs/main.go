package main

import (
	"log"
	"os"

	"github.com/klauern/openai-orgs/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "openai-orgs",
		Usage: "CLI for OpenAI Platform Management API",
		Commands: []*cli.Command{
			cmd.AuditLogsCommand(),
			cmd.InvitesCommand(),
			cmd.UsersCommand(),
			cmd.ProjectsCommand(),
			cmd.ProjectUsersCommand(),
			cmd.ProjectServiceAccountsCommand(),
			cmd.ProjectApiKeysCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-key",
				Usage:   "OpenAI API key (can be set via OPENAI_API_KEY environment variable)",
				EnvVars: []string{"OPENAI_API_KEY"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
