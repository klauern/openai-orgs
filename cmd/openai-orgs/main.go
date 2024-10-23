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
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
