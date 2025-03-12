package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func ProjectUsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-users",
		Usage: "Manage project users",
		Commands: []*cli.Command{
			listProjectUsersCommand(),
			createProjectUserCommand(),
			retrieveProjectUserCommand(),
			modifyProjectUserCommand(),
			deleteProjectUserCommand(),
		},
	}
}

func listProjectUsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all project users",
		Flags: []cli.Flag{
			projectIDFlag,
			limitFlag,
			afterFlag,
		},
		Action: listProjectUsers,
	}
}

func createProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new project user",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
			roleFlag,
		},
		Action: createProjectUser,
	}
}

func retrieveProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project user",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: retrieveProjectUser,
	}
}

func modifyProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify",
		Usage: "Modify a project user",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
			roleFlag,
		},
		Action: modifyProjectUser,
	}
}

func deleteProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a project user",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: deleteProjectUser,
	}
}

func listProjectUsers(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	projectUsers, err := client.ListProjectUsers(
		cmd.String("project-id"),
		limit,
		cmd.String("after"),
	)
	if err != nil {
		return wrapError("list project users", err)
	}

	data := TableData{
		Headers: []string{"ID", "Email", "Name", "Role", "Added At"},
		Rows:    make([][]string, len(projectUsers.Data)),
	}

	for i, user := range projectUsers.Data {
		data.Rows[i] = []string{
			user.ID,
			user.Email,
			user.Name,
			user.Role,
			user.AddedAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func createProjectUser(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	projectUser, err := client.CreateProjectUser(
		cmd.String("project-id"),
		cmd.String("id"),
		cmd.String("role"),
	)
	if err != nil {
		return wrapError("create project user", err)
	}

	fmt.Printf("Project User created:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nRole: %s\nAdded At: %s\n",
		projectUser.ID,
		projectUser.Email,
		projectUser.Name,
		projectUser.Role,
		projectUser.AddedAt.String(),
	)

	return nil
}

func retrieveProjectUser(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	projectUser, err := client.RetrieveProjectUser(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return wrapError("retrieve project user", err)
	}

	fmt.Printf("Project User details:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nRole: %s\nAdded At: %s\n",
		projectUser.ID,
		projectUser.Email,
		projectUser.Name,
		projectUser.Role,
		projectUser.AddedAt.String(),
	)

	return nil
}

func modifyProjectUser(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	projectUser, err := client.ModifyProjectUser(
		cmd.String("project-id"),
		cmd.String("id"),
		cmd.String("role"),
	)
	if err != nil {
		return wrapError("modify project user", err)
	}

	fmt.Printf("Project User modified:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nNew Role: %s\nAdded At: %s\n",
		projectUser.ID,
		projectUser.Email,
		projectUser.Name,
		projectUser.Role,
		projectUser.AddedAt.String(),
	)

	return nil
}

func deleteProjectUser(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	err := client.DeleteProjectUser(
		cmd.String("project-id"),
		cmd.String("id"),
	)
	if err != nil {
		return wrapError("delete project user", err)
	}

	fmt.Printf("Project User %s deleted successfully\n", cmd.String("id"))
	return nil
}
