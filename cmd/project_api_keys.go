package cmd

import (
	"fmt"
	"os"
	"strings"

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
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit the number of API keys returned",
			},
			&cli.StringFlag{
				Name:  "after",
				Usage: "Return API keys after this ID",
			},
		},
		Action: listProjectApiKeys,
	}
}

func retrieveProjectApiKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project API key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "api-key-id",
				Usage:    "ID of the API key to retrieve",
				Required: true,
			},
		},
		Action: retrieveProjectApiKey,
	}
}

func deleteProjectApiKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a project API key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
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
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

	projectID := c.String("project-id")
	limit := c.Int("limit")
	after := c.String("after")

	apiKeys, err := client.ListProjectApiKeys(projectID, limit, after)
	if err != nil {
		return fmt.Errorf("failed to list project API keys: %w", err)
	}

	fmt.Println("ID | Name | Redacted Value | Created At | Owner")
	fmt.Println(strings.Repeat("-", 80))
	for _, apiKey := range apiKeys.Data {
		fmt.Printf("%s | %s | %s | %s | %s (%s)\n",
			apiKey.ID,
			apiKey.Name,
			apiKey.RedactedValue,
			apiKey.CreatedAt.String(),
			apiKey.Owner.Name,
			apiKey.Owner.Type,
		)
	}

	return nil
}

func retrieveProjectApiKey(c *cli.Context) error {
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

	projectID := c.String("project-id")
	apiKeyID := c.String("api-key-id")

	apiKey, err := client.RetrieveProjectApiKey(projectID, apiKeyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project API key: %w", err)
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
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

	projectID := c.String("project-id")
	apiKeyID := c.String("api-key-id")

	err := client.DeleteProjectApiKey(projectID, apiKeyID)
	if err != nil {
		return fmt.Errorf("failed to delete project API key: %w", err)
	}

	fmt.Printf("API Key with ID %s has been deleted from project %s\n", apiKeyID, projectID)
	return nil
}
