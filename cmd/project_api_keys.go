package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

func ProjectAPIKeysCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-api-keys",
		Usage: "Manage project API keys",
		Commands: []*cli.Command{
			listProjectAPIKeysCommand(),
			createProjectAPIKeyCommand(),
			retrieveProjectAPIKeyCommand(),
			deleteProjectAPIKeyCommand(),
		},
	}
}

func listProjectAPIKeysCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all project API keys",
		Flags: []cli.Flag{
			projectIDFlag,
			limitFlag,
			afterFlag,
		},
		Action: listProjectAPIKeys,
	}
}

func createProjectAPIKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new project API key",
		Flags: []cli.Flag{
			projectIDFlag,
			nameFlag,
		},
		Action: createProjectAPIKey,
	}
}

func retrieveProjectAPIKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project API key",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: retrieveProjectAPIKey,
	}
}

func deleteProjectAPIKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a project API key",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: deleteProjectAPIKey,
	}
}

func listProjectAPIKeys(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	apiKeys, err := client.ListProjectApiKeys(
		cmd.String("project-id"),
		int(cmd.Int("limit")),
		cmd.String("after"),
	)
	if err != nil {
		return fmt.Errorf("failed to list project API keys: %w", err)
	}

	switch cmd.String("output") {
	case "json":
		data, err := json.Marshal(apiKeys)
		if err != nil {
			return fmt.Errorf("failed to marshal API keys: %w", err)
		}
		fmt.Println(string(data))
	default:
		data := TableData{
			Headers: []string{"ID", "Name", "Created At", "Owner"},
			Rows:    make([][]string, len(apiKeys.Data)),
		}
		for i, key := range apiKeys.Data {
			data.Rows[i] = []string{
				key.ID,
				key.Name,
				key.CreatedAt.String(),
				key.Owner.String(),
			}
		}
		printTableData(data)
	}
	return nil
}

func createProjectAPIKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	apiKey, err := client.CreateProjectApiKey(
		cmd.String("project-id"),
		cmd.String("name"),
	)
	if err != nil {
		return fmt.Errorf("failed to create project API key: %w", err)
	}

	switch cmd.String("output") {
	case "json":
		data, err := json.Marshal(apiKey)
		if err != nil {
			return fmt.Errorf("failed to marshal API key: %w", err)
		}
		fmt.Println(string(data))
	default:
		data := TableData{
			Headers: []string{"ID", "Name", "Created At", "Owner"},
			Rows: [][]string{{
				apiKey.ID,
				apiKey.Name,
				apiKey.CreatedAt.String(),
				apiKey.Owner.String(),
			}},
		}
		printTableData(data)
	}
	return nil
}

func retrieveProjectAPIKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	apiKey, err := client.RetrieveProjectApiKey(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return fmt.Errorf("failed to retrieve project API key: %w", err)
	}

	switch cmd.String("output") {
	case "json":
		data, err := json.Marshal(apiKey)
		if err != nil {
			return fmt.Errorf("failed to marshal API key: %w", err)
		}
		fmt.Println(string(data))
	default:
		data := TableData{
			Headers: []string{"ID", "Name", "Created At", "Owner"},
			Rows: [][]string{{
				apiKey.ID,
				apiKey.Name,
				apiKey.CreatedAt.String(),
				apiKey.Owner.String(),
			}},
		}
		printTableData(data)
	}
	return nil
}

func deleteProjectAPIKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	err := client.DeleteProjectApiKey(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return fmt.Errorf("failed to delete project API key: %w", err)
	}

	fmt.Printf("Successfully deleted project API key %s\n", cmd.String("id"))
	return nil
}
