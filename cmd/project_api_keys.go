package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func ProjectApiKeysCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-api-keys",
		Usage: "Manage project API keys",
		Commands: []*cli.Command{
			listProjectApiKeysCommand(),
			retrieveProjectApiKeyCommand(),
			deleteProjectApiKeyCommand(),
		},
	}
}

func listProjectApiKeysCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all project API keys",
		Flags: []cli.Flag{
			projectIDFlag,
			limitFlag,
			afterFlag,
		},
		Action: listProjectApiKeys,
	}
}

func retrieveProjectApiKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project API key",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: retrieveProjectApiKey,
	}
}

func deleteProjectApiKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a project API key",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: deleteProjectApiKey,
	}
}

func listProjectApiKeys(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	projectApiKeys, err := client.ListProjectApiKeys(
		cmd.String("project-id"),
		limit,
		cmd.String("after"),
	)
	if err != nil {
		return wrapError("list project API keys", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Redacted Value", "Created At", "Owner"},
		Rows:    make([][]string, len(projectApiKeys.Data)),
	}

	for i, key := range projectApiKeys.Data {
		data.Rows[i] = []string{
			key.ID,
			key.Name,
			key.RedactedValue,
			key.CreatedAt.String(),
			fmt.Sprintf("%s (%s)", key.Owner.Name, key.Owner.Type),
		}
	}

	printTableData(data)
	return nil
}

func retrieveProjectApiKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	projectApiKey, err := client.RetrieveProjectApiKey(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return wrapError("retrieve project API key", err)
	}

	fmt.Printf("Project API Key details:\n")
	fmt.Printf("ID: %s\nName: %s\nRedacted Value: %s\nCreated At: %s\n",
		projectApiKey.ID,
		projectApiKey.Name,
		projectApiKey.RedactedValue,
		projectApiKey.CreatedAt.String(),
	)
	fmt.Printf("Owner: %s (%s)\n", projectApiKey.Owner.Name, projectApiKey.Owner.Type)

	return nil
}

func deleteProjectApiKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	err := client.DeleteProjectApiKey(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return wrapError("delete project API key", err)
	}

	fmt.Printf("Project API Key %s deleted successfully\n", cmd.String("id"))
	return nil
}
