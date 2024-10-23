package cmd

import (
	"fmt"
	"os"
	"strings"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

func ProjectServiceAccountsCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-service-accounts",
		Usage: "Manage project service accounts",
		Subcommands: []*cli.Command{
			listProjectServiceAccountsCommand(),
			createProjectServiceAccountCommand(),
			retrieveProjectServiceAccountCommand(),
			deleteProjectServiceAccountCommand(),
		},
	}
}

func listProjectServiceAccountsCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all project service accounts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit the number of service accounts returned",
			},
			&cli.StringFlag{
				Name:  "after",
				Usage: "Return service accounts after this ID",
			},
		},
		Action: listProjectServiceAccounts,
	}
}

func createProjectServiceAccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new project service account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Name of the service account",
				Required: true,
			},
		},
		Action: createProjectServiceAccount,
	}
}

func retrieveProjectServiceAccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project service account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID of the service account to retrieve",
				Required: true,
			},
		},
		Action: retrieveProjectServiceAccount,
	}
}

func deleteProjectServiceAccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a project service account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID of the service account to delete",
				Required: true,
			},
		},
		Action: deleteProjectServiceAccount,
	}
}

func listProjectServiceAccounts(c *cli.Context) error {
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

	projectID := c.String("project-id")
	limit := c.Int("limit")
	after := c.String("after")

	accounts, err := client.ListProjectServiceAccounts(projectID, limit, after)
	if err != nil {
		return fmt.Errorf("failed to list project service accounts: %w", err)
	}

	fmt.Println("ID | Name | Role | Created At")
	fmt.Println(strings.Repeat("-", 80))
	for _, account := range accounts.Data {
		fmt.Printf("%s | %s | %s | %s\n",
			account.ID,
			account.Name,
			account.Role,
			account.CreatedAt.String(),
		)
	}

	return nil
}

func createProjectServiceAccount(c *cli.Context) error {
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

	projectID := c.String("project-id")
	name := c.String("name")

	account, err := client.CreateProjectServiceAccount(projectID, name)
	if err != nil {
		return fmt.Errorf("failed to create project service account: %w", err)
	}

	fmt.Printf("Project service account created:\n")
	fmt.Printf("ID: %s\nName: %s\nRole: %s\nCreated At: %s\n",
		account.ID,
		account.Name,
		account.Role,
		account.CreatedAt.String(),
	)

	return nil
}

func retrieveProjectServiceAccount(c *cli.Context) error {
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

	projectID := c.String("project-id")
	id := c.String("id")

	account, err := client.RetrieveProjectServiceAccount(projectID, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve project service account: %w", err)
	}

	fmt.Printf("Project service account details:\n")
	fmt.Printf("ID: %s\nName: %s\nRole: %s\nCreated At: %s\n",
		account.ID,
		account.Name,
		account.Role,
		account.CreatedAt.String(),
	)

	return nil
}

func deleteProjectServiceAccount(c *cli.Context) error {
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

	projectID := c.String("project-id")
	id := c.String("id")

	err := client.DeleteProjectServiceAccount(projectID, id)
	if err != nil {
		return fmt.Errorf("failed to delete project service account: %w", err)
	}

	fmt.Printf("Project service account with ID %s has been deleted\n", id)
	return nil
}
