package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

func AdminAPIKeysCommand() *cli.Command {
	return &cli.Command{
		Name:  "admin-api-keys",
		Usage: "Manage organization API keys",
		Commands: []*cli.Command{
			listAdminAPIKeysCommand(),
			createAdminAPIKeyCommand(),
			retrieveAdminAPIKeyCommand(),
			deleteAdminAPIKeyCommand(),
		},
	}
}

func listAdminAPIKeysCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all organization API keys",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
		},
		Action: listAdminAPIKeys,
	}
}

func createAdminAPIKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new organization API key",
		Flags: []cli.Flag{
			nameFlag,
			&cli.StringSliceFlag{
				Name:     "scopes",
				Usage:    "API key scopes (comma-separated)",
				Required: true,
			},
		},
		Action: createAdminAPIKey,
	}
}

func retrieveAdminAPIKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve details of a specific organization API key",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: retrieveAdminAPIKey,
	}
}

func deleteAdminAPIKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete an organization API key",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: deleteAdminAPIKey,
	}
}

func listAdminAPIKeys(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	apiKeys, err := client.ListAdminAPIKeys(
		limit,
		cmd.String("after"),
	)
	if err != nil {
		return wrapError("list admin API keys", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Redacted Value", "Created At", "Last Used At", "Scopes"},
		Rows:    make([][]string, len(apiKeys.Data)),
	}

	for i, key := range apiKeys.Data {
		data.Rows[i] = []string{
			key.ID,
			key.Name,
			key.RedactedValue,
			key.CreatedAt.String(),
			key.LastUsedAt.String(),
			strings.Join(key.Scopes, ", "),
		}
	}

	printTableData(data)
	return nil
}

func createAdminAPIKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	name := cmd.String("name")
	scopes := cmd.StringSlice("scopes")

	apiKey, err := client.CreateAdminAPIKey(name, scopes)
	if err != nil {
		return wrapError("create admin API key", err)
	}

	fmt.Printf("API Key created:\n")
	fmt.Printf("ID: %s\nName: %s\nRedacted Value: %s\nCreated At: %s\n",
		apiKey.ID,
		apiKey.Name,
		apiKey.RedactedValue,
		apiKey.CreatedAt.String(),
	)
	fmt.Printf("Scopes: %s\n", strings.Join(apiKey.Scopes, ", "))

	return nil
}

func retrieveAdminAPIKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	apiKey, err := client.RetrieveAdminAPIKey(cmd.String("id"))
	if err != nil {
		return wrapError("retrieve admin API key", err)
	}

	fmt.Printf("API Key details:\n")
	fmt.Printf("ID: %s\nName: %s\nRedacted Value: %s\nCreated At: %s\n",
		apiKey.ID,
		apiKey.Name,
		apiKey.RedactedValue,
		apiKey.CreatedAt.String(),
	)
	fmt.Printf("Last Used At: %s\n", apiKey.LastUsedAt.String())
	fmt.Printf("Scopes: %s\n", strings.Join(apiKey.Scopes, ", "))

	return nil
}

func deleteAdminAPIKey(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	err := client.DeleteAdminAPIKey(cmd.String("id"))
	if err != nil {
		return wrapError("delete admin API key", err)
	}

	fmt.Printf("API Key %s deleted successfully\n", cmd.String("id"))
	return nil
}
