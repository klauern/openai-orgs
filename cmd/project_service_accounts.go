package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func ProjectServiceAccountsCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-service-accounts",
		Usage: "Manage project service accounts",
		Commands: []*cli.Command{
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
			projectIDFlag,
			limitFlag,
			afterFlag,
		},
		Action: listProjectServiceAccounts,
	}
}

func createProjectServiceAccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new project service account",
		Flags: []cli.Flag{
			projectIDFlag,
			nameFlag,
		},
		Action: createProjectServiceAccount,
	}
}

func retrieveProjectServiceAccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project service account",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: retrieveProjectServiceAccount,
	}
}

func deleteProjectServiceAccountCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a project service account",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: deleteProjectServiceAccount,
	}
}

func listProjectServiceAccounts(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	serviceAccounts, err := client.ListProjectServiceAccounts(
		cmd.String("project-id"),
		limit,
		cmd.String("after"),
	)
	if err != nil {
		return wrapError("list project service accounts", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Created At"},
		Rows:    make([][]string, len(serviceAccounts.Data)),
	}

	for i, account := range serviceAccounts.Data {
		data.Rows[i] = []string{
			account.ID,
			account.Name,
			account.CreatedAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func createProjectServiceAccount(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	serviceAccount, err := client.CreateProjectServiceAccount(
		cmd.String("project-id"),
		cmd.String("name"),
	)
	if err != nil {
		return wrapError("create project service account", err)
	}

	fmt.Printf("Project Service Account created:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\n",
		serviceAccount.ID,
		serviceAccount.Name,
		serviceAccount.CreatedAt.String(),
	)

	return nil
}

func retrieveProjectServiceAccount(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	serviceAccount, err := client.RetrieveProjectServiceAccount(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return wrapError("retrieve project service account", err)
	}

	fmt.Printf("Project Service Account details:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\n",
		serviceAccount.ID,
		serviceAccount.Name,
		serviceAccount.CreatedAt.String(),
	)

	return nil
}

func deleteProjectServiceAccount(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	err := client.DeleteProjectServiceAccount(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return wrapError("delete project service account", err)
	}

	fmt.Printf("Project Service Account %s deleted successfully\n", cmd.String("id"))
	return nil
}
