package cmd

import (
	"fmt"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

func ProjectApiKeysCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-api-keys",
		Usage: "Manage project API keys",
		Subcommands: []*cli.Command{
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
			&cli.StringFlag{
				Name:     "api-key-id",
				Usage:    "ID of the API key to delete",
				Required: true,
			},
		},
		Action: deleteProjectApiKey,
	}
}

func listProjectApiKeys(c *cli.Context) error {
	client := newClient(c)

	apiKeys, err := client.ListProjectApiKeys(
		c.String("project-id"),
		c.Int("limit"),
		c.String("after"),
	)
	if err != nil {
		return wrapError("list project API keys", err)
	}

	headers := []string{"ID", "Name", "Redacted Value", "Created At", "Owner"}
	rows := make([][]string, len(apiKeys.Data))

	for i, apiKey := range apiKeys.Data {
		rows[i] = []string{
			apiKey.ID,
			apiKey.Name,
			apiKey.RedactedValue,
			apiKey.CreatedAt.String(),
			fmt.Sprintf("%s (%s)", apiKey.Owner.Name, apiKey.Owner.Type),
		}
	}

	printTable(headers, rows)
	return nil
}

func retrieveProjectApiKey(c *cli.Context) error {
	client := newClient(c)

	apiKey, err := client.RetrieveProjectApiKey(
		c.String("project-id"),
		c.String("api-key-id"),
	)
	if err != nil {
		return wrapError("retrieve project API key", err)
	}

	fmt.Printf("API Key details:\n")
	fmt.Printf("ID: %s\nName: %s\nRedacted Value: %s\nCreated At: %s\n",
		apiKey.ID,
		apiKey.Name,
		apiKey.RedactedValue,
		apiKey.CreatedAt.String(),
	)
	fmt.Printf("Owner: %s (%s)\n", apiKey.Owner.Name, apiKey.Owner.Type)

	return nil
}

func deleteProjectApiKey(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	projectID := c.String("project-id")
	apiKeyID := c.String("api-key-id")

	err := client.DeleteProjectApiKey(projectID, apiKeyID)
	if err != nil {
		return fmt.Errorf("failed to delete project API key: %w", err)
	}

	fmt.Printf("API Key with ID %s has been deleted from project %s\n", apiKeyID, projectID)
	return nil
}
