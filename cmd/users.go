package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func UsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Manage organization users",
		Commands: []*cli.Command{
			listUsersCommand(),
			retrieveUserCommand(),
			deleteUserCommand(),
			modifyUserRoleCommand(),
		},
	}
}

func listUsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all users",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
		},
		Action: listUsers,
	}
}

func retrieveUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific user",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: retrieveUser,
	}
}

func deleteUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a user",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: deleteUser,
	}
}

func modifyUserRoleCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify-role",
		Usage: "Modify a user's role",
		Flags: []cli.Flag{
			idFlag,
			roleFlag,
		},
		Action: modifyUserRole,
	}
}

func listUsers(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	users, err := client.ListUsers(
		limit,
		cmd.String("after"),
	)
	if err != nil {
		return wrapError("list users", err)
	}

	data := TableData{
		Headers: []string{"ID", "Email", "Name", "Role"},
		Rows:    make([][]string, len(users.Data)),
	}

	for i, user := range users.Data {
		data.Rows[i] = []string{
			user.ID,
			user.Email,
			user.Name,
			user.Role,
		}
	}

	printTableData(data)
	return nil
}

func retrieveUser(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	user, err := client.RetrieveUser(cmd.String("id"))
	if err != nil {
		return wrapError("retrieve user", err)
	}

	fmt.Printf("User details:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nRole: %s\n",
		user.ID,
		user.Email,
		user.Name,
		user.Role,
	)

	return nil
}

func deleteUser(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	err := client.DeleteUser(cmd.String("id"))
	if err != nil {
		return wrapError("delete user", err)
	}

	fmt.Printf("User %s deleted successfully\n", cmd.String("id"))
	return nil
}

func modifyUserRole(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	id := cmd.String("id")
	role := cmd.String("role")

	err := client.ModifyUserRole(id, role)
	if err != nil {
		return wrapError("modify user role", err)
	}

	// Retrieve the updated user to show the changes
	user, err := client.RetrieveUser(id)
	if err != nil {
		return wrapError("retrieve updated user", err)
	}

	fmt.Printf("User role modified:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nNew Role: %s\n",
		user.ID,
		user.Email,
		user.Name,
		user.Role,
	)

	return nil
}
